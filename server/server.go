package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// 创建接口
type Server struct {
	Ip   string
	Port int

	//在线用户列表、
	OnlineMap map[string]*User
	maplock   sync.RWMutex //锁

	//消息广播
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string)}

	return server
}

// 启动服务器接口
func (this *Server) Start() error {
	//listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer listener.Close() //防止忘记关闭

	//启动监听
	go this.ListenMessager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		//do handle
		go this.Handler(conn)

	}

}

// 监听Message广播消息channel的goroutine
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		this.maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.maplock.Unlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn, this)
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				log.Printf("User %s disconnected", user.Name) // 更加详细的日志
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				log.Printf("Error reading from user %s: %v", user.Name, err) // 更加详细的日志
				user.Offline()                                               // 确保在发生错误时也离线
				return
			}

			msg := string(buf[:n-1]) // 提取消息去除“\n”
			user.DoMessage(msg)      // 用户针对msg进行消息处理

			select {
			case isLive <- true: // 用户的任意消息，代表当前活跃
			default:
				// 防止阻塞，即使select没有准备好
			}
		}
	}()

	// 超时定时器
	timeout := time.Second * 100
	timer := time.NewTimer(timeout) // 创建一个定时器

	for {
		select {
		case <-isLive:
			// 活跃，重置定时器
			timer.Reset(timeout) // 重置定时器

		case <-timer.C:
			// 已经超时，将当前user关闭
			log.Printf("Kicking out user %s due to timeout", user.Name) // 更加详细的日志
			user.SendMsg("你被踢了\n")
			close(user.C)
			conn.Close()
			this.maplock.Lock()
			delete(this.OnlineMap, user.Name) // 确保从在线用户列表中删除
			this.maplock.Unlock()
			return
		}
	}
}
