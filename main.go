package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	// "sort"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var masterKey bool = false
var mutex = &sync.Mutex{}
var result [][]string

type NodeInfo struct {
	NodeId     int    `json:"nodeId"`
	NodeIpAddr string `json:"nodeIpAddr"`
	Port       string `json:"port"`
	IsMaster   bool   `json:"isMaster"`
}

var nodes = make(map[NodeInfo]string)
var masterNode NodeInfo

/* A Request/Response format to transfer between nodes
   `Message` is the sorted/unsorted slice */
type data struct {
	Source  NodeInfo
	Dest    NodeInfo
	Type    string
	Message []string
}

func main() {
	numberOfNodes := flag.Int("numberOfNodes", 3, "number of slaves to use")
	clusterIp := flag.String("clusterIp", "127.0.0.1", "ip address of slave node")
	port := flag.String("port", "3000", "port to use")
	flag.Parse()

	parsedPortInInt, err := strconv.ParseInt(*port, 10, 64)
	if err != nil {
		fmt.Println("Error parsing port number")
	}

	// sampleData := []string{"Sah", "Deepak", "Abhishek", "Sharma", "Zathura", "Harsh", "Jay", "Eight", "Nine"}

	wg.Add(*numberOfNodes)
	for i := 0; i < *numberOfNodes; i++ {
		parsedPortInInt++

		node := createNode(*clusterIp, strconv.Itoa(int(parsedPortInInt)))
		go selectMasterNode(node)
	}
	wg.Wait()

	wg.Add(*numberOfNodes)
	for k, v := range nodes {
		if v == "master" {
			masterNode = k
			go listenOnPort(k)
		} else {
			go connectToNode(k)
		}
	}
	wg.Wait()
}

func getRequestObject(source NodeInfo, dest NodeInfo, dataType string, dataToSort []string) data {
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
		Type:    dataType,
		Message: dataToSort,
	}
}

func createNode(ipAddr string, port string) NodeInfo {
	return NodeInfo{
		NodeId:     1,
		NodeIpAddr: ipAddr,
		Port:       port,
		IsMaster:   false,
	}
}

func connectToNode(node NodeInfo) {
	defer wg.Done()
	laddr, _ := net.ResolveTCPAddr("tcp", node.NodeIpAddr+":"+node.Port)
	raddr, _ := net.ResolveTCPAddr("tcp", masterNode.NodeIpAddr+":"+masterNode.Port)

	for {
		conn, err := net.DialTCP("tcp", laddr, raddr)
		if err == nil {
			request := getRequestObject(node, masterNode, "getData", []string{"a"})
			json.NewEncoder(conn).Encode(request)
			handleResponseFromSlave(conn)
			conn.Close()
			break
		}
		fmt.Println("There is no Master node available. Waiting...", err)
		// break
	}
}

func listenOnPort(node NodeInfo) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":"+node.Port)
	if err != nil {
		fmt.Printf("Unable to create server at port: %s\n", node.Port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection.")
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	var request data
	json.NewDecoder(conn).Decode(&request)
	if request.Type == "getData" {
		fmt.Println("getData")
	}
	// sort.Strings(request.Message)
	var response data
	response = getRequestObject(request.Dest, request.Source, "sorted", request.Message)
	json.NewEncoder(conn).Encode(&response)
	conn.Close()
}

func handleResponseFromSlave(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	var response data
	decoder.Decode(&response)
	result = append(result, response.Message)
	fmt.Println("This is the result: ", result)
}

func divideWork([]string) {}

func selectMasterNode(node NodeInfo) {
	mutex.Lock()
	if masterKey {
		nodes[node] = "slave"
		wg.Done()
		mutex.Unlock()
		return
	}
	fmt.Println(node.Port)
	masterKey = true
	masterNode = node
	nodes[node] = "master"
	wg.Done()
	mutex.Unlock()
}
