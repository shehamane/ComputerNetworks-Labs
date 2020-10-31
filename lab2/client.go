package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

var (
	user     = flag.String("u", "", "User name")
	host     = flag.String("h", "", "Host")
	port     = flag.Int("p", 22, "Port")
	password = flag.String("psw", "", "Password")
)

func main() {
	flag.Parse()

	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 3 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := ssh.Dial("tcp", addr, config)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	session, err := conn.NewSession()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Shell()
	if err != nil {
		panic(err)
	}

	commands := []string{"pwd", "ls", "whoami", "mkdir lol", "ls", "rm -d lol", "ls"}

	for _, cmd := range commands{
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil{
			panic(err)
		}
	}

	err = session.Wait()
	if err != nil{
		panic(err)
	}
}
