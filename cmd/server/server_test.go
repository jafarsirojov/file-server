package main

import (
	"bufio"
	"bytes"
	"file-server/cmd/rpc"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"
)

func Test_uploadFileToServer(t *testing.T) {
	const addr = "localhost:7777"

	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatalf("can't listen on %s: %v", addr, err)
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatalf("can't accept %v\n", err)
			}
			go handleConn(conn)
		}
	}()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to %s: %v", addr, err)
	}
	writer := bufio.NewWriter(conn)
	options := "Go01.pdf"
	cmd := "upload#" + options
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		t.Fatal("can't write", err)
	}
	file, err := ioutil.ReadFile("./testFile/"+options)
	if err != nil {
		t.Fatalf("can't read file: %v\n", err)
	}
	_, err = writer.Write(file)
	if err != nil {
		t.Fatalf("can't copy the file %s\n", file)
	}
	err = writer.Flush()
	if err != nil {
		t.Fatalf("can't flush %v\n", err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("can't conn close %v\n", err)
	}
	uploadFile, err := ioutil.ReadFile("./serverFile/"+options)
	if err != nil {
		t.Fatalf("can't read file upload error:  %v\n",err)
	}
	if  !bytes.Equal(file, uploadFile){
		t.Fatalf("files are not equal: %v", err)
	}
}


func Test_downloadFileInServer(t *testing.T) {
	logOutput, err := os.OpenFile("logServer.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logOutput)
	defer logOutput.Close()
	const addr = "localhost:8888"
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatalf("can't listen on %s: %v", addr, err)
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatalf("can't accept %v\n", err)
			}
			go handleConn(conn)
		}
	}()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to %s: %v", addr, err)
	}
	writer := bufio.NewWriter(conn)
	options := "Go01.pdf"
	cmd := "download#" + options
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		t.Fatalf("can't write command %v\n", err)
	}
	reader := bufio.NewReader(conn)
	downloadBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("can't reader file error: %v\n", err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("can't conn close %v\n", err)
	}
	err = ioutil.WriteFile("./Downloads/"+options, downloadBytes, 0666)
	if err != nil {
		t.Fatalf("can't write file: %v\n", err)
	}
	downloadFile, err := ioutil.ReadFile("./Downloads/"+options)
	if err != nil {
		t.Fatalf("can't read file upload error:  %v\n",err)
	}
	if !bytes.Equal(downloadBytes,downloadFile) {
		t.Fatalf("files are not equal: %v", err)
	}
}


func Test_listFileForServer(t *testing.T) {
	logOutput, err := os.OpenFile("logServer.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logOutput)
	defer logOutput.Close()
	const addr = "localhost:9999"
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatalf("can't listen on %s: %v", addr, err)
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatalf("can't accept %v\n", err)
			}
			go handleConn(conn)
		}
	}()
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("can't connect to %s: %v", addr, err)
	}
	writer := bufio.NewWriter(conn)
	cmd := "list# "
	err = rpc.WriteLine(cmd, writer)
	if err != nil {
		t.Fatal("can't write", err)
	}
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		if err != io.EOF {
			t.Fatalf("can't read: %v", err)
		}
	}
	err = conn.Close()
	if err != nil {
		t.Fatalf("can't conn close %v\n", err)
	}
	if line!="Go01.pdf\n"{
		t.Fatal("list not found")
	}
}

