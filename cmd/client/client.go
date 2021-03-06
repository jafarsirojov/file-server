package main

import (
	"bufio"
	"file-server/pkg/rpc"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var download = flag.String("download", "default", "download")
var upload   = flag.String("upload", "default", "upload")
var list     = flag.Bool("list", false, "list")

func main() {
	logOutput, err := os.OpenFile("logClient.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logOutput)
	defer logOutput.Close()

	flag.Parse()
	var cmd,operations, options string
	if *download != "default" {
		options = *download
		operations = rpc.Download
	} else if *upload != "default" {
		options = *upload
		operations = rpc.Upload
	} else if *list != false {
		operations = rpc.List
	} else {
		return
	}


	const addr = "localhost:7777"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("can't connect to %s: %v", addr, err)
	}
	defer conn.Close()
	log.Println("client connected")
	writer := bufio.NewWriter(conn)

	cmd = operations + "#" + options
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		log.Fatal(err)
	}

	if operations == "upload" {
		file, err := ioutil.ReadFile(options)
		if err != nil {
			log.Printf("can't not find the file %s\n", file)
		}
		_, err = writer.Write(file)
		if err != nil {
			log.Printf("can't copy the file %s\n", file)
			return
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalf("can't flush %v\n", err)
		}
		fmt.Println("Файл успешно загружен в сервер)")
		return
	}

	reader := bufio.NewReader(conn)

	if operations == "download" {
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("can't read data: %v", err)
				fmt.Println("can't download!!!")
				return
			}
		}
		options="Downloads/"+options
		err = ioutil.WriteFile(options, bytes, 0666)
		if err != nil {
			log.Printf("can't write file: %v", err)
		}
	}

	line, err := rpc.ReadLine(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read: %v", err)
			return
		}
	}
	if operations == "download" {
		line = "Файл успешно скачен ;)"
	}
	fmt.Println(line)

}
