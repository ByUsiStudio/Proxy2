package main

import (
	"flag"
	"fmt"
	"hp-lib/log"
	"hp-lib/net/cmd"
	"hp-lib/util"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kardianos/service"
)

type program struct {
	serverIp   string
	serverPort int
	deviceId   string
	cmdClient  *cmd.CmdClient
	stopChan   chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		log.Info("服务以交互模式启动")
	} else {
		log.Info("服务启动成功")
	}
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	if service.Interactive() {
		log.Info("服务以交互模式停止")
	} else {
		log.Info("服务正在停止")
	}
	close(p.stopChan)
	return nil
}

func (p *program) run() {
	p.cmdClient = cmd.NewCmdClient(func(message string) {
		log.Info(message)
	})

	p.cmdClient.Connect(p.serverIp, p.serverPort, p.deviceId)
	log.Infof("已连接到服务器 %s:%d (设备ID: %s)", p.serverIp, p.serverPort, p.deviceId)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !p.cmdClient.GetStatus() {
					log.Info("与服务器断开连接，正在重连...")
					p.cmdClient.Connect(p.serverIp, p.serverPort, p.deviceId)
				}
			case <-p.stopChan:
				log.Info("重连循环已停止")
				return
			}
		}
	}()

	<-p.stopChan
	log.Info("服务核心逻辑已停止")
}

func parseConnectionCode(c string) (serverIp string, serverPort int, deviceId string, err error) {
	if c == "" {
		return "", 0, "", fmt.Errorf("连接码为空")
	}

	base32 := util.DecodeFromLowerCaseBase32(strings.TrimSpace(c))
	conn := strings.Split(base32, ",")

	if len(conn) != 2 {
		return "", 0, "", fmt.Errorf("连接码格式错误：分割后长度不为2（实际：%d）", len(conn))
	}

	server := conn[0]
	deviceId = conn[1]

	split := strings.Split(server, ":")
	if len(split) != 2 {
		return "", 0, "", fmt.Errorf("服务器地址格式错误：%s（应为 IP:端口）", server)
	}

	serverIp = split[0]
	port, err := strconv.Atoi(split[1])
	if err != nil {
		return "", 0, "", fmt.Errorf("端口号不是有效数字：%s，错误：%v", split[1], err)
	}

	if port < 1 || port > 65535 {
		return "", 0, "", fmt.Errorf("端口号无效：%d（必须在 1-65535 之间）", port)
	}

	return serverIp, port, deviceId, nil
}

func parseServer(server string) (serverIp string, serverPort int, err error) {
	if server == "" {
		return "", 0, fmt.Errorf("服务器地址为空")
	}
	err, _, host, port := util.ProtocolInfo(server)
	if err != nil {
		return "", 0, fmt.Errorf("服务器地址解析失败：%v", err)
	}
	return host, port, nil
}

func main() {
	var (
		c             string
		server        string
		deviceId      string
		serviceAction string
	)

	flag.StringVar(&c, "c", "", "连接码（与-server/-deviceId二选一）")
	flag.StringVar(&server, "server", util.DefaultServer, "服务器地址（默认：https://pro2.cdifit.cn）")
	flag.StringVar(&deviceId, "deviceId", "", "设备ID（与-server配合使用）")
	flag.StringVar(&serviceAction, "action", "", "服务操作：install/start/stop/uninstall/status")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方法：%s [参数]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "参数说明：")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\n示例：")
		fmt.Fprintln(os.Stderr, "  使用连接码：", os.Args[0], "-c \"你的连接码\"")
		fmt.Fprintln(os.Stderr, "  使用服务器和设备ID：", os.Args[0], "-server \"https://pro2.cdifit.cn:443\" -deviceId \"设备ID\"")
		fmt.Fprintln(os.Stderr, "  安装服务：", os.Args[0], "-c \"你的连接码\" -action install")
		fmt.Fprintln(os.Stderr, "  启动服务：", os.Args[0], "-action start")
		fmt.Fprintln(os.Stderr, "  停止服务：", os.Args[0], "-action stop")
		fmt.Fprintln(os.Stderr, "  查看状态：", os.Args[0], "-action status")
		fmt.Fprintln(os.Stderr, "  卸载服务：", os.Args[0], "-action uninstall")
	}

	flag.Parse()

	var (
		serverIp   string
		serverPort int
		err        error
	)

	if serviceAction != "" {
		switch serviceAction {
		case "start", "stop", "status", "uninstall":
			break

		case "install":
			if c != "" {
				serverIp, serverPort, deviceId, err = parseConnectionCode(c)
				if err != nil {
					log.Fatalf("连接码解析失败：%v", err)
				}
			} else if server != "" && deviceId != "" {
				serverIp, serverPort, err = parseServer(server)
				if err != nil {
					log.Fatalf("服务器地址解析失败：%v", err)
				}
			} else {
				log.Fatal("错误：安装服务必须通过 -c 参数或 -server/-deviceId 参数指定连接信息")
			}

		default:
			log.Fatalf("无效的操作：%s，支持的操作：install/start/stop/uninstall/status", serviceAction)
		}
	} else {
		if c != "" {
			serverIp, serverPort, deviceId, err = parseConnectionCode(c)
			if err != nil {
				log.Fatalf("连接码解析失败：%v", err)
			}
		} else if server != "" && deviceId != "" {
			serverIp, serverPort, err = parseServer(server)
			if err != nil {
				log.Fatalf("服务器地址解析失败：%v", err)
			}
		} else {
			log.Fatal("错误：必须通过 -c 参数或 -server/-deviceId 参数指定连接信息")
		}
	}

	var serviceArgs []string
	if c != "" {
		serviceArgs = []string{"-c", c}
	} else {
		serviceArgs = []string{"-server", server, "-deviceId", deviceId}
	}
	serviceConfig := &service.Config{
		Name:        "hp-lite",
		DisplayName: "hp-lite",
		Description: "hp-lite 命令行客户端服务，用于与中心服务器通信",
		Arguments:   serviceArgs,
	}
	if serviceAction == "install" {
		if c != "" {
			serviceConfig.Description += "（连接码：" + c + "）"
		} else {
			serviceConfig.Description += "（服务器：" + server + "，设备ID：" + deviceId + "）"
		}
	}

	prg := &program{
		serverIp:   serverIp,
		serverPort: serverPort,
		deviceId:   deviceId,
		stopChan:   make(chan struct{}),
	}

	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Fatalf("服务创建失败：%v", err)
	}

	if !service.Interactive() {
		log.DisableColor()
	}

	switch serviceAction {
	case "install":
		err = s.Install()
		if err != nil {
			log.Fatalf("服务安装失败：%v", err)
		}
		log.Printf("✅ 服务安装成功！服务名称：%s", serviceConfig.Name)
		log.Printf("📌 连接码：%s", c)
		log.Println("💡 后续操作：")
		log.Println("   启动服务：", os.Args[0], "-action start")
		log.Println("   停止服务：", os.Args[0], "-action stop")
		log.Println("   查看状态：", os.Args[0], "-action status")
		log.Println("   卸载服务：", os.Args[0], "-action uninstall")
		log.Println("   [注意事项]：")
		log.Println("   1、安装服务后请不要删除当前文件，如果需要删除当前文件，请先停止服务、然后在卸载服务、最后在删除文件")
		log.Println("   2、更新程序前请先停止服务再替换程序然后再启动服务")

	case "start":
		err = s.Start()
		if err != nil {
			log.Fatalf("服务启动失败：%v", err)
		}
		log.Println("✅ 服务启动成功！可通过 -action status 查看状态")

	case "stop":
		err = s.Stop()
		if err != nil {
			log.Fatalf("服务停止失败：%v", err)
		}
		log.Println("✅ 服务停止成功！")

	case "uninstall":
		err = s.Uninstall()
		if err != nil {
			log.Fatalf("服务卸载失败：%v", err)
		}
		log.Println("✅ 服务卸载成功！")

	case "status":
		status, err := s.Status()
		if err != nil {
			log.Fatalf("获取服务状态失败：%v", err)
		}
		switch status {
		case service.StatusRunning:
			log.Println("🟢 服务状态：运行中")
		case service.StatusStopped:
			log.Println("🔴 服务状态：已停止")
		default:
			log.Printf("🟡 服务状态：%v", status)
		}

	case "":
		log.Printf("🚀 以交互模式启动（服务器：%s:%d，设备ID：%s）", serverIp, serverPort, deviceId)
		err = s.Run()
		if err != nil {
			log.Fatalf("交互模式运行失败：%v", err)
		}
	}
}
