package main

import (
	"flag"
	"fmt"
	"hp-server-lib/config"
	"hp-server-lib/log"
	"hp-server-lib/net/acme"
	"hp-server-lib/net/http"
	"hp-server-lib/net/server"
	"hp-server-lib/service"
	"hp-server-lib/task"
	"hp-server-lib/version"
	"hp-server-lib/web"
	"os"

	daemon "github.com/kardianos/service"
)

var serviceAction string
var configPath string

type program struct {
	stopChan chan struct{}
}

func (p *program) Start(s daemon.Service) error {
	if daemon.Interactive() {
		log.Info("服务以交互模式启动")
	} else {
		log.Info("服务启动成功")
	}
	go p.run()
	return nil
}

func (p *program) Stop(s daemon.Service) error {
	if daemon.Interactive() {
		log.Info("服务以交互模式停止")
	} else {
		log.Info("服务正在停止")
	}
	close(p.stopChan)
	return nil
}

func (p *program) run() {
	if err := config.LoadConfig(configPath); err != nil {
		log.Error(fmt.Sprintf("配置加载失败：%v", err))
		return
	}
	log.Info("配置文件加载成功")

	go p.starServer()

	<-p.stopChan
	log.Info("核心业务逻辑已停止")
}

func (p *program) starServer() {
	tcpServer := server.NewCmdServer()
	go tcpServer.StartServer(config.ConfigData.Cmd.Port)

	quicServer := server.NewHpQuicServer(server.NewHPHandler())
	go quicServer.StartServer(config.ConfigData.Tunnel.Port)

	hpTcpServer := server.NewHPTcpServer(server.NewHPHandler())
	go hpTcpServer.StartServer(config.ConfigData.Tunnel.Port)

	go web.StartWebServer(config.ConfigData.Admin.Port)

	go service.InitForward()

	if config.ConfigData.Tunnel.OpenDomain {
		go http.StartHttpServer()
		go http.StartHttpsServer()

		go service.InitDomainCache()
		go service.InitReverseECache()

		go func() {
			err2 := acme.StartAcmeServer(config.ConfigData.Acme.Email, config.ConfigData.Acme.HttpPort)
			if err2 != nil {
				log.Error("证书申请服务启动失败..." + err2.Error())
			} else {
				task.StartSslTask()
			}
		}()
	}

	<-p.stopChan
	log.Info("服务正在关闭...")
}

func init() {
	flag.StringVar(&serviceAction, "action", "", "服务操作：install/start/stop/uninstall/status")
	flag.StringVar(&configPath, "conf", "app.yml", "配置文件路径（默认：app.yml）")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方法：%s [参数]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "参数说明：")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\n示例：")
		fmt.Fprintln(os.Stderr, "  交互模式运行：", os.Args[0], "-conf app.yml")
		fmt.Fprintln(os.Stderr, "  安装服务：", os.Args[0], "-conf app.yml -action install")
		fmt.Fprintln(os.Stderr, "  启动服务：", os.Args[0], "-action start")
		fmt.Fprintln(os.Stderr, "  停止服务：", os.Args[0], "-action stop")
		fmt.Fprintln(os.Stderr, "  查看状态：", os.Args[0], "-action status")
		fmt.Fprintln(os.Stderr, "  卸载服务：", os.Args[0], "-action uninstall")
	}
}

func main() {
	flag.Parse()

	prg := &program{
		stopChan: make(chan struct{}),
	}

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取当前工作目录失败：%v\n", err)
		os.Exit(1)
	}

	serviceConfig := &daemon.Config{
		Name:             "hp-lite-server",
		DisplayName:      "hp-lite-server",
		Description:      "hp-lite-server 核心服务（含隧道、管理后台、证书管理等功能）",
		Arguments:        []string{"-conf", configPath},
		WorkingDirectory: workDir,
	}

	s, err := daemon.New(prg, serviceConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "服务创建失败：%v\n", err)
		os.Exit(1)
	}

	if !daemon.Interactive() {
		log.DisableColor()
	}

	switch serviceAction {
	case "install":
		if err := s.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "服务安装失败：%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ 服务安装成功！服务名称：%s\n", serviceConfig.Name)
		fmt.Printf("📌 配置文件路径：%s\n", configPath)
		fmt.Println("💡 后续操作：")
		fmt.Println("   启动服务：", os.Args[0], "-action start")
		fmt.Println("   停止服务：", os.Args[0], "-action stop")
		fmt.Println("   查看状态：", os.Args[0], "-action status")
		fmt.Println("   卸载服务：", os.Args[0], "-action uninstall")
		fmt.Println("   [注意事项]：")
		fmt.Println("   1、安装服务后请不要删除当前文件，如果需要删除当前文件，请先停止服务、然后在卸载服务、最后在删除文件")
		fmt.Println("   2、更新程序前请先停止服务再替换程序然后再启动服务")

	case "start":
		if err := s.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "服务启动失败：%v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 服务启动成功！可通过 -action status 查看状态")

	case "stop":
		if err := s.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "服务停止失败：%v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 服务停止成功！")

	case "uninstall":
		if err := s.Uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "服务卸载失败：%v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 服务卸载成功！")

	case "status":
		status, err := s.Status()
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取服务状态失败：%v\n", err)
			os.Exit(1)
		}
		switch status {
		case daemon.StatusRunning:
			fmt.Println("🟢 服务状态：运行中")
		case daemon.StatusStopped:
			fmt.Println("🔴 服务状态：已停止")
		default:
			fmt.Printf("🟡 服务状态：%v\n", status)
		}

	case "":
		printBanner(configPath, workDir)
		if err := s.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "交互模式运行失败：%v\n", err)
			os.Exit(1)
		}
	}
}

func printBanner(cfg string, workDir string) {
	banner := `
 ____                      ____
|  _ \ _ __ _____  ___   _|___ \
| |_) | '__/ _ \ \/ / | | | __) |
|  __/| | | (_) >  <| |_| |/ __/
|_|   |_|  \___/_/\_\\__, |_____|
                     |___/
`
	fmt.Println("\x1b[36m" + banner + "\x1b[0m")
	fmt.Printf("\x1b[36m::\x1b[0m \x1b[37mProxy2 Server\x1b[0m \x1b[90m(v%s)\x1b[0m\n", version.Version)
	fmt.Printf("\x1b[36m::\x1b[0m \x1b[90m配置文件: \x1b[37m%s\x1b[0m\n", cfg)
	fmt.Printf("\x1b[36m::\x1b[0m \x1b[90m工作目录: \x1b[37m%s\x1b[0m\n", workDir)
	fmt.Println()
}
