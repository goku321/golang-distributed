package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var masterKey bool = false
var mutex = &sync.Mutex{}
var result [][]string

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

	nodeType := flag.String("nodeType", "master", "type of node")
	numberOfNodes := flag.Int("numberOfNodes", 3, "number of slaves to use")
	clusterIp := flag.String("clusterIp", "127.0.0.1", "ip address of slave node")
	port := flag.String("port", "3000", "port to use")
	flag.Parse()

	parsedPortInInt, err := strconv.ParseInt(*port, 10, 64)
	if err != nil {
		fmt.Println("Error parsing port number")
	}

	// sampleData := []string{"Sah", "Deepak", "Abhishek", "Sharma", "Zathura", "Harsh", "Jay", "Eight", "Nine"}
	// ip, _ := net.InterfaceAddrs()

	if *nodeType == "master" {
		wg.Add(*numberOfNodes)
		for i := 0; i < *numberOfNodes; i++ {
			parsedPortInInt++

			node := createNode(*clusterIp, strconv.Itoa(int(parsedPortInInt)))
			go selectMasterNode(node)
		}
		// for i, j := 0, 0; i < *numberOfSlaves; i, j = i+1, j+3 {
		// 	parsedPortInInt++

		// 	slaveNode := createNode(*clusterIp, strconv.Itoa(int(parsedPortInInt)))
		// 	go listenOnPort(slaveNode)

		// 	requestObject := getRequestObject(masterNode, slaveNode, sampleData[j:j+3])
		// 	go connectToNode(slaveNode, requestObject)
		// }
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
	for {
		conn, err := net.DialTimeout("tcp", node.NodeIpAddr+":"+node.Port, time.Duration(10)*time.Second)
		if err == nil {
			json.NewEncoder(conn).Encode(request)
			handleResponseFromSlave(conn)
			conn.Close()
			break
		}
		fmt.Println("There is no slave node available. Waiting...")
	}
}

func listenOnPort(node NodeInfo) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":"+node.Port)
	if err != nil {
		fmt.Println("unable to create server")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection.")
		}

		handleConnection(conn)
		conn.Close()
		break
	}
}

func handleConnection(conn net.Conn) {
	var request data
	json.NewDecoder(conn).Decode(&request)
	sort.Strings(request.Message)
	var response data
	response = getRequestObject(request.Dest, request.Source, request.Message)
	json.NewEncoder(conn).Encode(&response)
}

func handleResponseFromSlave(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	var response data
	decoder.Decode(&response)
	result = append(result, response.Message)
	fmt.Println(result)
}

func divideWork([]string) {}

func selectMasterNode(node NodeInfo) {
	mutex.Lock()
	if masterKey {
		wg.Done()
		mutex.Unlock()
		return
	}
	fmt.Println(node.Port)
	masterKey = true
	wg.Done()
	// Assign Node as Master
	mutex.Unlock()
}
