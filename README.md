# Go DNS data transfer

A simple tool for transfering data using the DNS protocol.

## Overview

DNS Data Transfer is a Go application that allows users to send and receive data using DNS queries. The client generates unique IDs and encodes data into an alphanumeric format, while the server listens for DNS requests and saves received data to files.

## Features

- **Client**: Sends data via DNS queries.
- **Server**: Listens for DNS queries and saves data to files. Can handle multiple transfers at once.
- **Data Encoding**: Converts UTF-8 text to a numeric representation, with each character's Unicode code represented as a zero-padded string of width 3 digits.
- **Chunking**: Splits data into manageable chunks for transmission.

## Prerequisites

- Docker
- Go 1.23 or later


### Run all

This will build the client, create and start the docker container for the server.
```shell
make all
```

Open another terminal and run
```shell
./client/client -file data.txt -dns 127.0.0.1:53
```
Replace `127.0.0.1` with the docker container IP.


#### Build only

The below command will build the client and will create the docker container for the server
```shell
make build
```

#### Run the server

Run the container with the desired IP address for A records.

Replace `192.168.100.1` in the `Makefile` with the IP address you want the server to respond with for A records.

The DNS server will listen on port 53 for UDP queries and respond with the specified IP address for any query.

```shell
make run
```

### Sending data

```
./server -ip 192.168.1.1
./client -file data.txt -dns 127.0.0.1:53
```
