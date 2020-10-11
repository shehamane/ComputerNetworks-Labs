package main

import (
	"flag"
	"github.com/goftp/file-driver"
	"github.com/goftp/server"
	"log"
)

func main(){
	var (
		root = flag.String("root", "", "Root directory to serve")
		user = flag.String("user", "admin", "Username for login")
		pass = flag.String("pass", "12345", "Password for login")
		port = flag.Int("port", 2121, "Port")
		host = flag.String("host", "localhost", "Host")
	)
	flag.Parse()
	if *root == "" {
		log.Fatalf("Please set a root to serve with -root")
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
		WelcomeMessage: "Добро пожаловать на сервер",
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *user, *pass)
	ftpServer := server.NewServer(opts)
	err := ftpServer.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
