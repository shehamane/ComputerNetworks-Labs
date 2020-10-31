package main

import (
	"bufio"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os/exec"
	"strings"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		term := terminal.NewTerminal(s, "> ")
		reader := bufio.NewReader(s)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			line = line[:len(line)-1]
			log.Printf("command> %s", line)
			words := strings.Split(line, " ")
			result := exec.Command(words[0], words[1:]...)
			out, _ := result.CombinedOutput()
			if out != nil {
				term.Write([]byte(line + "> \n"))
				term.Write(append(out, '\n'))
			}
		}
		log.Println("terminal closed")
	})

	log.Println("starting ssh server on port 3029...")
	log.Fatal(ssh.ListenAndServe(":3029", nil))
}