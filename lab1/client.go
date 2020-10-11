package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func fileContains(files *[]string, filename string) bool {
	for _, name := range *files {
		if name == filename {
			return true
		}
	}
	return false
}

func sendFile(connection *ftp.ServerConn, urlSrc string, urlDest string) error {
	file, _ := os.Open(urlSrc)
	defer file.Close()
	data := make([]byte, 100)
	file.Read(data)
	err := connection.Stor(urlDest, bytes.NewBuffer(data))
	return err
}

func fetchFile(connection *ftp.ServerConn, urlSrc string, urlDest string) error {
	r, err := connection.Retr(urlSrc)
	defer r.Close()
	img, _ := os.Create(urlDest)
	defer img.Close()
	imageBytes, _ := ioutil.ReadAll(r)
	img.Write(imageBytes)
	return err
}

func handleRequest(conn *ftp.ServerConn, request string) error {
	var err error
	request = request[:len(request)-1]
	requestSplit := strings.Split(request, " ")
	commandName := requestSplit[0]
	commandArgs := requestSplit[1:len(requestSplit)]
	switch commandName {
	case "login":
		err = conn.Login(commandArgs[0], commandArgs[1])
	case "exit":
		err = conn.Quit()
		os.Exit(0)
	case "delete":
		err = conn.Delete(commandArgs[0])
	case "pull":
		err = fetchFile(conn, commandArgs[0], commandArgs[1])
	case "push":
		err = sendFile(conn, commandArgs[0], commandArgs[1])
	case "mkdir":
		err = conn.MakeDir(commandArgs[0])
	case "rmdir":
		err = conn.RemoveDir(commandArgs[0])
	case "ls":
		nameList, _ := conn.NameList("/")
		for _, name := range nameList {
			println(name)
		}
	}
	return err
}

func main() {
	var (
		host = flag.String("host", "localhost", "Host's IP")
		port = flag.String("port", "21", "Port")
	)
	flag.Parse()

	connection, err := ftp.Dial(*host+":"+*port, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		request, _ := reader.ReadString('\n')
		err := handleRequest(connection, request)
		if err != nil {
			log.Fatal(err)
		}
	}
}
