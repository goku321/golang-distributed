package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"runtime"
	// "sort"
	"sync"
)

var wg sync.WaitGroup

// var sampleData []string = ["Sah", "Deepak", "Abhishek", "Sharma", "Zathura", "Harsh", "Jay"]

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

	sampleData := []string{"Sah", "Deepak", "Abhishek", "Sharma", "Zathura", "Harsh", "Jay"}

	masterNode := NodeInfo{
		NodeId:     0,
		NodeIpAddr: *clusterIp,
		Port:       *port,
	}

	if *nodeType == "master" {
		wg.Add(3)
		slaveNode1 := NodeInfo{NodeId: 1, NodeIpAddr: *clusterIp, Port: "3002"}
		slaveNode2 := NodeInfo{NodeId: 2, NodeIpAddr: *clusterIp, Port: "3003"}
		slaveNode3 := NodeInfo{NodeId: 3, NodeIpAddr: *clusterIp, Port: "3004"}
		requestObject1 := getRequestObject(masterNode, slaveNode1, sampleData)
		requestObject2 := getRequestObject(masterNode, slaveNode2, sampleData)
		requestObject3 := getRequestObject(masterNode, slaveNode3, sampleData)
		go connectToNode(slaveNode1, requestObject1)
		go connectToNode(slaveNode2, requestObject2)
		go connectToNode(slaveNode3, requestObject3)
		wg.Wait()
	} else {
		slaveNode := createNode(*clusterIp, *port)
		listenOnPort(slaveNode)
	}
}

func getRequestObject(source NodeInfo, dest NodeInfo, dataToSort []string) data {
	return data{
		Source: NodeInfo{
			NodeId:     source.NodeId,
			NodeIpAddr: source.NodeIpAddr,
			Port:       source.Port,
		},
		Dest: NodeInfo{
			NodeId:     dest.NodeId,
			NodeIpAddr: dest.NodeIpAddr,
			Port:       dest.Port,
		},
		Message: dataToSort,
	}
}

func createNode(ipAddr string, port string) NodeInfo {
	return NodeInfo{
		NodeId:     1,
		NodeIpAddr: ipAddr,
		Port:       port,
	}
}

func connectToNode(node NodeInfo, request data) {
	defer wg.Done()
	conn, _ := net.Dial("tcp", node.NodeIpAddr+":"+node.Port)
	json.NewEncoder(conn).Encode(request)
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

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var request data
	json.NewDecoder(conn).Decode(&request)
	fmt.Println("Formatted Data: ", request)
}
