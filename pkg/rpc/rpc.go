package rpc

import (
	"bufio"
	"io/ioutil"
	"log"
)

func ReadLine(reader *bufio.Reader) (line string, err error) {
	return reader.ReadString('\n')
}

func WriteLine(line string, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	return
}

const fileServerDir = "serverFile/"//const fileServerDir = "cmd/server/serverFile/"

func ListServerFile() (listServerFile string, err error) {
	files, err := ioutil.ReadDir(fileServerDir)
	if err != nil {
		log.Println("can't read dir: ", err)
		return "", err
	}
	listServerFile = ""
	for _, file := range files {
		if listServerFile == "" {
			listServerFile = file.Name()
		} else {
			listServerFile = listServerFile + " " + file.Name()
		}
	}
	listServerFile = listServerFile + "\n"
	return listServerFile, err
}

const Download  = "download"
const Upload    = "upload"
const List  = "list"
