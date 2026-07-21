package hp_android_lib

import (
	"hp-lib/log"
	"hp-lib/net/cmd"
	"hp-lib/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Callback interface {
	SendResult(msg string)
}

var (
	cmdClient    *cmd.CmdClient
	exitChan     = make(chan struct{})
	mu           sync.Mutex
	isRunning    bool
	reconnecting bool
)

func Start(c string, callback Callback) {
	mu.Lock()
	if isRunning {
		mu.Unlock()
		callback.SendResult("服务已在运行")
		return
	}
	isRunning = true
	mu.Unlock()

	if c != "" {
		log.Printf("使用连接码模式连接")
		base32 := util.DecodeFromLowerCaseBase32(strings.TrimSpace(c))
		con := strings.Split(base32, ",")
		if len(con) != 2 {
			callback.SendResult("连接码错误")
			return
		}
		server := con[0]
		deviceId := con[1]
		split := strings.Split(server, ":")
		if len(split) != 2 {
			callback.SendResult("连接码错误")
			return
		}
		serverPort, _ := strconv.Atoi(split[1])
		cmdClient = cmd.NewCmdClient(callback.SendResult)
		cmdClient.Connect(split[0], serverPort, deviceId)

		go func() {
			for {
				select {
				case <-exitChan:
					return
				default:
					time.Sleep(time.Duration(10) * time.Second)
					mu.Lock()
					if !isRunning {
						mu.Unlock()
						return
					}
					if !reconnecting && cmdClient != nil && !cmdClient.GetStatus() {
						reconnecting = true
						mu.Unlock()
						cmdClient.Connect(split[0], serverPort, deviceId)
						callback.SendResult("中心服务器重连中")
						mu.Lock()
						reconnecting = false
					}
					mu.Unlock()
				}
			}
		}()

		<-exitChan
	} else {
		callback.SendResult("连接码错误")
	}
}

func Close() bool {
	mu.Lock()
	defer mu.Unlock()
	if cmdClient != nil {
		cmdClient.Close()
		cmdClient = nil
	}
	isRunning = false
	select {
	case <-exitChan:
	default:
		close(exitChan)
	}
	exitChan = make(chan struct{})
	return true
}

func GetStatus() bool {
	mu.Lock()
	defer mu.Unlock()
	if cmdClient != nil {
		return cmdClient.GetStatus()
	}
	return false
}
