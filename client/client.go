package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string,serverPort int) *Client{
	//创建客户端对象
	client := &Client{
		ServerIp: serverIp,
		ServerPort: serverPort,
		flag: 999,
	}

	//链接server
	conn,err := net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err != nil{
		fmt.Println("net.Dial error:",err)
		return nil
	}

	client.conn = conn

	//返回对象
	return client
}

//处理server回应消息
func(client *Client) DealResponse(){
	//一旦client.conn有数据，就直接copy到stdout标准输出上
	io.Copy(os.Stdout,client.conn)
}

func (client *Client) menu() bool {
    var flag int
    fmt.Println("1. 公聊模式")
    fmt.Println("2. 私聊模式")
    fmt.Println("3. 更新用户名")
    fmt.Println("0. 退出")
    fmt.Scan(&flag)

    if flag >= 0 && flag <= 3 {
        client.flag = flag
        return true
    } else {
        fmt.Println(">>> 请输入合法范围内的数字<<<")
        return false
    }
}

func (client *Client) PublicChat() {
    var chatMsg string
    fmt.Println(">>> 请输入公聊内容, exit退出...")
    fmt.Scan(&chatMsg)
    for chatMsg != "exit" {
        // 公聊的发送
        if len(chatMsg) != 0 {
            sendMsg := chatMsg + "\n"
            _, err := client.conn.Write([]byte(sendMsg))
            if err != nil {
                fmt.Println("conn Write err:", err)
				break
            }
            
        }
		chatMsg=""
		fmt.Println(">>> 公聊结束, exit退出...")
		fmt.Scanln(&chatMsg)
    }
    
}

//查询在线用户
func (client *Client) SelectUsers(){
	sendMsg := "who\n"
	_,err := client.conn.Write([]byte(sendMsg))
	if err !=nil{
		fmt.Println("conn Write err:",err)
		return
	}
}

func (client *Client) PrivateChat(){
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名]，exit退出：")
	fmt.Scan(&remoteName)

	for remoteName !="exit"{
		fmt.Println(">>>>请输入消息内容，exit退出")
		fmt.Scan(&chatMsg)

		for chatMsg!="exit"{
			if len(chatMsg) != 0 {
            
				sendMsg := "to|"+remoteName+"|"+chatMsg+"\n\n"
                _, err := client.conn.Write([]byte(sendMsg))           
				if err != nil {
                fmt.Println("conn Write err:", err)
				break
            }
            
        }
		chatMsg = ""
		fmt.Println(">>>>请输入消息内容，exit退出")
		fmt.Scan(&chatMsg)

		}

		client.SelectUsers()
	    fmt.Println(">>>>请输入聊天对象[用户名]，exit退出：")
	    fmt.Scan(&remoteName)

	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入用户名：")
	fmt.Scan(&client.Name)

	sendMsg := "rename|"+client.Name+"\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("conn.Write err:",err)
		return false
	}

	return true
}

func (client *Client) Run(){
	for client.flag != 0{
		for client.menu() != true{
		}

		//根据不同的模型处理不同的业务
		switch client.flag{
		case 1:
			//公聊
			client.PublicChat()
			break
		case 2:
			//私聊模式
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			client.UpdateName()
			break	

		}
	}
}

var serverIp string
var serverPort int

//    ./client -ip 127.0.0.1 -port 8888
func init(){
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器ip地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort,"port",8888,"设置服务器端口（默认是8888）")
}

func main(){
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp,serverPort)
	if client == nil{
		fmt.Println(">>>> 链接服务器失败...")
		return
	}

	//单独开一个gorotine去处理server的回执消息
	go client.DealResponse() 

	fmt.Println(">>>> 链接服务器成功...")

	//启动客户端业务
	client.Run()
}