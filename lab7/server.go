package main

import "net"
import "fmt"
import "bufio"
import "strings"

func main() {

	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":3029")

	// Открываем порт
	conn, _ := ln.Accept()

	// Запускаем цикл
	for {
		// Будем прослушивать все сообщения разделенные \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", string(message))
		// Процесс выборки для полученной строки
		newmessage := strings.ToUpper(message)
		// Отправить новую строку обратно клиенту
		conn.Write([]byte(newmessage + "\n"))
	}
}