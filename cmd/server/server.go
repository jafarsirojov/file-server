package main

import (
	"bufio"
	"file-server/cmd/rpc"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

const fileServerDir = "cmd/server/serverFile/"


func main() {
	_ = os.Mkdir(fileServerDir, 0666)
	logOutput, err := os.OpenFile("logServer.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logOutput)
	defer logOutput.Close()

	const addr = "0.0.0.0:7777"
	log.Println("server starting")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can't listen on %s: %v", addr, err)
	}
	defer listener.Close()
	log.Println("server started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept %v\n", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	name := conn.RemoteAddr().String()
	log.Printf("%+v connected\n", name)
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("error while reading: %v", err)
		return
	}

	index := strings.IndexByte(line, '#')
	writer := bufio.NewWriter(conn)
	if index == -1 {
		log.Printf("invalid line received %s", line)
		err := rpc.WriteLine("error: invalid line", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
		return
	}
	cmd, options := line[:index], line[index+1:len(line)-1]

	switch cmd {
	case "list":
		listServerFile, err := rpc.ListServerFile()
		if err != nil {
			log.Printf("cant list server file %v", err)
			rpc.WriteLine("Упс! В сервере нет файлов для скачивание)", writer)
			return
		}
		rpc.WriteLine(listServerFile, writer)

	case "upload":
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("can't read data: %v", err)
			}
		}
		err = ioutil.WriteFile(fileServerDir+options, bytes, 0666)
		if err != nil {
			log.Printf("can't write file: %v", err)
		}
		msg:="Файл в сервер успешно загружен!!!\n"
		err = rpc.WriteLine(msg, writer)
		if err != nil {
			log.Printf("can't writ line error: %v\n",err)
		}

	case "download":
		var fileDownload string
		fileDownload=fileServerDir+options
		file,err:= ioutil.ReadFile(fileDownload)
		if err != nil {
			log.Printf("can't not find the file %s\n",file)
		}
		_, err = writer.Write( file)
		if err != nil {
			log.Printf("can't copy the file %s\n",file)
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalf("can't flush %v\n", err)
		}
	}
}
