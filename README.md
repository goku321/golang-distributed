# golang-distributed
This project implements a system in which master node among the given nodes distribute the work to the slave nodes. The master node is selected from the pool nodes using election algorithm. After processing the data, the slave node returns data back to the master node. This program sorts the list of names in a distributed manner.

### Commands to run:

1. Runs the program with default parameters:</br>
- `numberOfNodes`(Number of slave nodes to use): `3`
- `clusterIp`(IP Address of slave node): `127.0.0.1`
- `port`(Port Number): `3000`

```sh
golang-distributed > go run main.go
```

2. Runs the program by customizing the input parameters
- `numberOfNodes`(Number of slave nodes to use): `100`
- `clusterIp`(IP Address of slave node): `127.0.0.1`
- `port`(Port Number): `7000`

```sh
golang-distributed > go run main.go --numberOfNodes 100 --port 7000
```