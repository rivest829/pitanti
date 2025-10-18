package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/v3/pkg/session"
	"log"
	"os"
	"pitanti/utils"
	"strings"
	"time"

	"github.com/topfreegames/pitaya/v3/pkg/client"
)

func main() {
	// 创建Pitaya客户端
	cli := client.New(logrus.InfoLevel, 1*time.Second)
	defer cli.Disconnect()

	// 启动交互式命令行界面
	runCLI(cli)
}

func listenMsg(cli *client.Client) {
	for message := range cli.IncomingMsgChan {
		fmt.Printf("新消息: %s : %d\n", message.Route, message.ID)
		fmt.Printf("Type: %d\n", message.Type)
		fmt.Printf("Data: %s\n", string(message.Data))
		fmt.Printf("Err: %v\n", message.Err)
	}
}

func runCLI(cli *client.Client) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n=== Pitaya 调试客户端 ===")
	fmt.Println("可用命令:")
	fmt.Println("  login             - 登录")
	fmt.Println("  req             - 请求")
	fmt.Println("  quit             - 退出客户端")
	fmt.Println("========================\n")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Split(input, " ")
		command := strings.ToLower(parts[0])
		switch command {
		case "req":
			router := strings.ToLower(parts[1])
			var data []byte
			if len(parts) > 2 {
				data = []byte(parts[2])
			}
			_, err := cli.SendRequest(router, []byte(data))
			utils.Must(err)
		case "login":
			roleId := strings.ToLower(parts[1])
			login(roleId, cli)
		case "quit", "exit":
			fmt.Println("退出客户端...")
			return
		default:
			fmt.Printf("未知命令: %s\n", command)
		}
	}
}

func login(roleId string, cli *client.Client) {
	cli.SetClientHandshakeData(&session.HandshakeData{
		User: map[string]interface{}{
			"roleId": roleId,
		},
	})
	// 连接到服务器 (默认端口3100)
	err := cli.ConnectTo("localhost:3100")
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	fmt.Println("成功连接到服务器!")
	go listenMsg(cli)

	_, err = cli.SendRequest("game.comp.login", []byte(fmt.Sprintf("{\"roleId\":\"%s\"}", roleId)))
	utils.Must(err)
}
