package main

import (
	"bufio"
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"os/exec"
	"strings"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))
		reader := bufio.NewReader(s)
		for{
			text, _ := reader.ReadString('\n')
			if text == "" || text == "\n"{
				continue
			}
			fmt.Println(">"+text)
			words := strings.Split(text, " ")
			result := exec.Command(words[0], words[1:]...)
			result.Stdout = s
			result.Run()
			io.WriteString(s, fmt.Sprintf("lol"))
			//io.WriteString(s, fmt.Sprintf("%s\n", string(output)))
		}
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}