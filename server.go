package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
	wg        sync.WaitGroup
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string, 1024),
	}
}

func (server *Server) ListenMessager() {
	for msg := range server.Message {
		server.mapLock.RLock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.RUnlock()
	}
}

func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	defer func() {
		conn.Close()
		server.wg.Done()
	}()
	user := NewUser(conn, server)
	user.Online()

	isLive := make(chan bool)

	go func() {
		defer func() {
			user.Offline()
			isLive <- false
		}()
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 300):
			user.SendMsg("你被踢了")
			close(user.C)
			return
		}
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	go server.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		server.wg.Add(1)
		go server.Handler(conn)
	}
}
