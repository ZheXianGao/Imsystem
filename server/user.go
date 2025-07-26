package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string //channel
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn,server *Server) *User{
	useraddr := conn.RemoteAddr().String()
	user:= &User{
		Name: useraddr,
		Addr: useraddr,
		C: make(chan string),
		conn: conn,
		server : server,
	}

	//启动监听当前User channel消息的goroutine
	go user.ListenMessage()


	return user

}

func (this *User) Online(){
	this.server.maplock.Lock()
	this.server.OnlineMap[this.Name]=this
	this.server.maplock.Unlock()

	//广播上线消息
	this.server.BroadCast(this,"已上线")

}

func (this *User) Offline(){
	//用户下线，将用户从onlinemao中移除
	this.server.maplock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.maplock.Unlock()

	//广播下线消息
	this.server.BroadCast(this,"下线")
}

//用户处理消息
func (this *User) DoMessage(msg string){
	if msg=="who"{
		this.server.maplock.Lock()
	    for _,user := range this.server.OnlineMap{
			onlineMsg := "[" + user.Addr + "]"+user.Name+":"+"在线...\n"
			this.SendMsg(onlineMsg)
		}
	    this.server.maplock.Unlock()

	}else if len(msg)>7 && msg[:7]=="rename|"{
		//消息格式：rename|张三
		newName := msg[7:]

		//判断这个名字是否存在
		_,ok := this.server.OnlineMap[newName]
		if ok{
			this.SendMsg("The name was used")
		}else{
			this.server.maplock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName] = this
			this.server.maplock.Unlock()

			this.Name = newName
			this.SendMsg("您已经更新用户名："+this.Name+"\n")
		}

	} else if len(msg)>4 && msg[:3] == "to|"{
		//消息：to|张三|消息内容

		//1 获取对方的用户名 
		remoteName := strings.Split(msg, "|")[1]
		if remoteName ==""{
			this.SendMsg("消息格式不正确，请使用 to|张三|消息内容  的格式")
			return
		}
		//2 根据用户名 得到对方User对象
		remoteUser ,ok := this.server.OnlineMap[remoteName]
		if !ok{
			this.SendMsg("该用户名不存在")
			return
		}

		//3 获取消息内容，通过对方User对象将消息发送出去
		content := strings.Split(msg,"|")[2]
		if content == ""{
			this.SendMsg("无消息内容，请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对您说:" + content+"\n")

	}else{
		this.server.BroadCast(this,msg)
	}
	
}

func(this *User) SendMsg(msg string){
	this.conn.Write([]byte(msg))
}

//监听当前Userchannel，有消息则发送给对端客户端
func(this *User) ListenMessage(){
	for{
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))      //[]byte(...)：将字符串类型的消息转换成字节切片，因为conn.Write方法的参数要求是[]byte类型。
	}
} 