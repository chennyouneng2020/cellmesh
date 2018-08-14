package main

import (
	"bufio"
	"fmt"
	"github.com/davyxu/cellmesh/demo/proto"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellmesh/svcfx"
	"os"
	"strings"
)

func login() (agentAddr string) {

	loginReq, err := service.CreateConnection("demo.login", service.NewMsgRequestor)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer loginReq.Stop()

	proto.Login(loginReq, &proto.LoginREQ{
		Version:  "1.0",
		Platform: "demo",
		UID:      "1234",
	}, func(ack *proto.LoginACK) {

		agentAddr = fmt.Sprintf("%s:%d", ack.Server.IP, ack.Server.Port)
	})

	return
}

func getAgentRequestor(agentAddr string) service.Requestor {
	waitGameReady := make(chan service.Requestor)
	go service.KeepConnection(service.NewMsgRequestor(agentAddr), "", waitGameReady)
	return <-waitGameReady
}

func ReadConsole(callback func(string)) {

	for {

		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)

		callback(text)

	}
}

func main() {

	svcfx.Init()

	agentAddr := login()

	fmt.Println("agent:", agentAddr)

	agentReq := getAgentRequestor(agentAddr)

	proto.Verify(agentReq, &proto.VerifyREQ{
		GameToken: "verify",
	}, func(ack *proto.VerifyACK) {

		fmt.Println(ack)
	})

	ReadConsole(func(s string) {

		proto.Chat(agentReq, &proto.ChatREQ{
			Content: s,
		}, func(ack *proto.ChatACK) {

			fmt.Println(ack)
		})

	})

}
