package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var msg = make(chan string)

func read() {
	for {
		m := <-msg
		fmt.Print(m)

	}
}

func write(conn net.Conn) {
	for {
		s := make([]byte, 10, 10)

		n, err := conn.Read(s)
		if err != nil {

		}

		msg <- string(s[:n])

	}

}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {

	}
	defer conn.Close()
	go read()
	go write(conn)
	for {
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err != nil {

		}
		if s == "exit\r\n" {

			conn.Write([]byte(s))

			break
		} else if s == "change\r\n" {

			fmt.Println("请输入要更改的名字")
			s, _ := reader.ReadString('\n')
			s2 := []byte(s[0 : len(s)-2])
			s3 := append([]byte("change "), s2...)
			conn.Write([]byte(s3))
			continue

		}
		conn.Write([]byte(s))
	}

}
