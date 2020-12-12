package main

import "net"
import "fmt"
import "bufio"
import "strings"
import "strconv"

type Point struct {
	x, y int
}

func handler(message string, points *[]Point) {
	words := strings.Split(message, " ")
	if words[0] == "create" {
		*points = (*points)[0:(len(words)-1)/2]
		for i := 0; i < (len(words)-1)/2; i++ {
			x, _ := strconv.Atoi(words[2*i+1])
			y, _ := strconv.Atoi(words[2*i+2])
			(*points)[i].x = x
			(*points)[i].y = y
		}
	}
	if words[0] == "edit" {
		k, _ := strconv.Atoi(words[1])
		x, _ := strconv.Atoi(words[2])
		y, _ := strconv.Atoi(words[3])
		(*points)[k - 1].x = x
		(*points)[k - 1].y = y
	}
}

func main() {
	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":3029")

	// Открываем порт
	conn, _ := ln.Accept()

	var points []Point = make([]Point, 1000)

	// Запускаем цикл
	for {
		// Будем прослушивать все сообщения разделенные \n
		reader := bufio.NewReader(conn)
		message, _ := reader.ReadString('\n')
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", message)
		// Процесс выборки для полученной строки
		handler(message[0:len(message)-1], &points)
		// Отправить новую строку обратно клиенту
		//conn.Write([]byte(newmessage + "\n"))
	}
}
