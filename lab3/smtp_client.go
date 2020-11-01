package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/smtp"
	"strings"
)

func sendMessage(addr string, from string, psw string, to string, subject string, msg string) error {
	host, _, _ := net.SplitHostPort(addr)
	tlsconfig := &tls.Config{
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", from, psw, host)
	err = client.Auth(auth)
	if err != nil {
		return err
	}

	err = client.Mail(from)
	if err != nil {
		return err
	}

	err = client.Rcpt(to)
	if err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(wc,
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n", to, subject, msg)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}

var (
	to = flag.String("to", "", "To")
	subject = flag.String("sub", "", "Subject")
	msg = flag.String("msg", "", "Message")
)

func main() {
	flag.Parse()

	configFile ,err := ioutil.ReadFile("config.txt")
	if err!=nil{
		panic(err)
	}

	lines := strings.Split(string(configFile), "\n")

	err = sendMessage(lines[0], lines[1], lines[2], *to, *subject, *msg)
	if err!=nil{
		panic(err)
	}
}
