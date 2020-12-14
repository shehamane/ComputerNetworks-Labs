package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "localhost:3029")
	for {
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// Отправляем в socket
		fmt.Fprintf(conn, text + "\n")
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
	}
}
