package main

import (
	"net"
	"strconv"
	"strings"
)
import "fmt"
import "bufio"

type Vector struct {
	dim int
	cs  []float64
}

//createvector name x1 y1 x2 y2
//dotproduct name1 name2

//createvector a 1 2 3 4

func getVector(cs []string) (*Vector, error) {
	v := new(Vector)
	v.dim = len(cs)
	v.cs = make([]float64, v.dim)
	var err error
	for i := 0; i < v.dim; i++ {
		v.cs[i], err = strconv.ParseFloat(cs[i], 32)
		if err != nil {
			return nil, err
		}
	}
	return v, err
}

func dotProduct(v1, v2 *Vector) float64 {
	var ans float64 = 0

	for i := 0; i < v1.dim; i++ {
		ans += v1.cs[i] * v2.cs[i]
	}
	return ans
}

func handleRequest(message string, vectors *map[string]*Vector) string {
	words := strings.Split(message, " ")
	command := words[0]
	args := words[1:]

	if command == "cv" {
		var err error
		(*vectors)[args[0]], err = getVector(args[1:])
		if err != nil {
			return "Ошибка: неверные аргументы"
		}
		return "Вектор создан!"
	}
	if command == "dp" {
		if len(args) != 2 {
			return "Ошибка: неверное число аргументов"
		}
		v1 := (*vectors)[args[0]]
		v2 := (*vectors)[args[1]]

		if v1 == nil || v2 == nil {
			return "Ошибка: один из указанных векторов не существует"
		}

		if v1.dim != v2.dim {
			return "Ошибка: вектора разной размерности"
		}
		dp := fmt.Sprintf("%f", dotProduct(v1, v2))
		return dp
	}
	return ""
}

func main() {
	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":3029")

	for {
		conn, _ := ln.Accept()

		vectors := make(map[string]*Vector)

		for {
			reader := bufio.NewReader(conn)
			message, _ := reader.ReadString('\n')
			fmt.Print("Message Received:", message)
			if strings.TrimSpace(message) == "exit" {
				conn.Write([]byte("Вы покинули сервер." + "\n"))
				conn.Close()
				break
			}
			ans := handleRequest(strings.TrimSpace(message), &vectors)
			conn.Write([]byte(ans + "\n"))
		}
	}
}
