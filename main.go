package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"sort"
	"runtime"
	"sync"
	"net"
)

var wg sync.WaitGroup

type NodeInfo struct {
	NodeId     int    `json:"nodeId"`
	NodeIpAddr string `json:"nodeIpAddr"`
	Port       string `json:"port"`
}

/* A Request/Response format to transfer between nodes
   `Message` is the sorted/unsorted slice */
type data struct {
	Source  NodeInfo
	Dest    NodeInfo
	Message []string
}

func main() {
	// Allocate one logical processor
	runtime.GOMAXPROCS(1)
	nodeType := flag.String("nodetype", "master", "type of node")
	// numberOfSlaves := flag.Int("numberofslaves", 3, "number of slaves to use")
	clusterIp := flag.String("clusterip", "127.0.0.1:8001", "ip address of slave node")
	port := flag.String("port", "8001", "port to use")
	flag.Parse()

	node1 := NodeInfo{
		NodeId:     1,
		NodeIpAddr: *clusterIp,
		Port:       *port,
	}
	fmt.Println(node1)
	if *nodeType == "master" {
		wg.Add(2)
		go connectToNode(NodeInfo{NodeId: 1, NodeIpAddr: *clusterIp, Port: "3002"})
		go connectToNode(NodeInfo{NodeId: 1, NodeIpAddr: *clusterIp, Port: "3003"})
		wg.Wait()
	} else {
		listenOnPort(node1)
	}
}

func createNode(node NodeInfo) {
}

func connectToNode(node NodeInfo) {
	defer wg.Done()
	conn, _ := net.Dial("tcp", node.NodeIpAddr+":"+node.Port)
	data := []string{"c", "b", "a"}
	json.NewEncoder(conn).Encode(data)
	status, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)
}

func listenOnPort(node NodeInfo) {
	fmt.Println(string(node.Port))
	ln, err := net.Listen("tcp", ":"+string(node.Port))
	if err != nil {
		fmt.Println("unable to create server")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}

		fmt.Println("This is the connection: ", conn)
		var data []string
		json.NewDecoder(conn).Decode(&data)
		sort.Strings(data)
		fmt.Println("Got this: ", data)
		// go handleConnection(conn)
	}
}

// func handleConnection(conn) {
// 	fmt.Println("This is the connection: ", conn)
// }
