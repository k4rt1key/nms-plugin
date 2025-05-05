package test

import (
	"context"
	"fmt"
	"nms-plugin/commands"
	"nms-plugin/winrm"
	"time"
)

func Test(c string) {

	command := commands.GetCommand(c)

	client := winrm.NewClient(
		"172.16.8.128",
		5985,
		"Administrator",
		"Mind@123",
	)

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	data, err := client.ExecuteCommand(ctx, command)

	if len(data) != 0 {
		fmt.Println("DATA" + string(data))
	}

	if err != nil {
		fmt.Println("ERROR" + err.Error())
	}

}
