package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type user struct {
	conn net.Conn
	name string
}

var online map[string]*user

var msg chan string

func manager() {

	online = make(map[string]*user)
	msg = make(chan string)
	for {
		m := <-msg

		var name = ""
		for _, ss := range m {
			if ss != ' ' {
				name += string(ss)
			} else {
				break
			}
		}
		fmt.Println(name)
		if name == "online" {
			use := strings.Split(m, " ")[1]
			for _, u := range online {
				if u.name == use {
					var joint string = "当前在线人数为:"
					joint += strconv.Itoa(len(online)) + "人\n"
					u.conn.Write([]byte(joint))
				}
				//u.conn.Write([]byte("当前在线人数为:" + string(rune(len(online))) + "人"))

			}
			continue
		}

		for _, u := range online {
			if u.name != name {
				u.conn.Write([]byte(m))
			}
		}
	}
}

func userHander(conn net.Conn) {

	defer conn.Close()

	addr := conn.RemoteAddr().String() + "login" + "\n"
	msg <- addr
	var uu = user{conn, conn.RemoteAddr().String()}

	online[conn.RemoteAddr().String()] = &uu

	for {
		by := make([]byte, 1024)
		n, _ := conn.Read(by)
		line := string(by[:n])
		//处理online
		if line == "online\r\n" {
			msg <- "online " + uu.name
			continue

		}

		if line == "exit\r\n" {

			delete(online, conn.RemoteAddr().String())
			msg <- uu.name + " 已成功退出"

			return
		}
		//处理修改名字
		var name = ""
		for _, ss := range line {

			if ss != ' ' {
				name += string(ss)
			} else {
				break
			}
		}
		if name == "change" {

			newName := strings.Split(line, " ")
			oldName := uu.name
			uu.name = newName[1]
			online[conn.RemoteAddr().String()].name = newName[1]
			msg <- "用户" + oldName + "已更名为" + newName[1] + "\n"
			continue
		}

		msg <- uu.name + ":  " + string(by[:n])

	}
}

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {

	}
	go manager()
	for {
		conn, err := listener.Accept()
		if err != nil {

		}
		go userHander(conn)

	}

}
