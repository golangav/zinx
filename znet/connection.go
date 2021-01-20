package znet

import (
	"net"
	"zinx/ziface"
	"fmt"
)

// 链接的模块

type Connection struct {
	// 当前链接的socket TCP 套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前链接的状态
	isClosed bool

	// 当前链接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc

	// 告知当前链接已经退出的 channel
	ExitChan chan bool
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, callback_api ziface.HandleFunc) *Connection  {
	c := &Connection{
		Conn:conn,
		ConnID:ConnID,
		isClosed:false,
		handleAPI:callback_api,
		ExitChan:make(chan bool, 1),
	}
	return c
}

// 链接读取数据的业务
func (c *Connection) StartReader()  {
	fmt.Println("reader goroutine is running... connID=", c.ConnID)
	defer fmt.Println("reader goroutine is exiting... connID=", c.ConnID, "remoteAddr=",c.Conn.RemoteAddr().String())
	defer c.stop()

	for {
		// 读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil{
			fmt.Println("recv buf err", err)
			continue
		}

		if err := c.handleAPI(c.Conn,buf,cnt); err != nil {
			fmt.Println("handle is error connID=", c.ConnID)
			break
		}
	}

}

// 启动链接，让当前的链接准备开始工作
func(c *Connection) start(){
	fmt.Println("conn start... connID=", c.ConnID)

	// 启动从当前链接的读取数据的业务
	go c.StartReader()

	// TODO 启动从当前链接的写数据的业务

}

// 停止链接，结束当前链接的工作
func(c *Connection) stop(){
	fmt.Println("conn stop... connID=", c.ConnID)

	if c.isClosed == true {
		return
	}

	// 关闭socker链接
	c.Conn.Close()
	close(c.ExitChan)
}

// 获取当前链接绑定的socket conn
func(c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接的 ID
func(c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP PORT
func(c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func(c *Connection) Send(data []byte) error {
	return nil
}