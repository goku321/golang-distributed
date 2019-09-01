package main

import (
	"flag"
	"fmt"
	"net"
	"encoding/json"
)

type NodeInfo struct {
	NodeId int
	NodeIpAddr string
	Port string
}

type data struct {
	Source NodeInfo
	Dest NodeInfo
	Message []string
}

func main() {
	clusterIp := flag.String("clusterip", "127.0.0.1:8001", "ip address of slave node")
  port := flag.String("port", "8001", "port to use")
  flag.Parse()
}

func createNode(node NodeInfo) {
}
