package main

import (
	"fmt"
	"os"
	"time"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codeskyblue/go-sh"
)

func main() {
	for {
		node, _ := sh.Command("kubectl", "-s", os.Getenv("KUBE_MASTER"), "describe", "pod", os.Getenv("HOSTNAME")).Command("grep", "Node").Command("awk", "{print $2}").Command("sed", "s@/.*@@").Output()
		fmt.Println(string(node))
		time.Sleep(time.Duration(5) * time.Second)
	}
}
