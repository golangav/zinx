package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

// 链接的模块
type Connection struct {
	// 当前链接的socket TCP 套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前链接的状态
	isClosed bool

	// 告知当前链接已经退出的 channel(由Reader告知Writer退出)
	ExitChan chan bool

	// 无缓冲的管道，用于读 写Groutine之间的消息通信
	msgChan chan []byte

	// 消息的管理MsgID 和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, ConnID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     ConnID,
		isClosed:   false,
		MsgHandler: msgHandle,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	return c
}

// 链接读取数据的业务
func (c *Connection) StartReader() {
	fmt.Println("reader goroutine is running... connID=", c.ConnID)
	defer fmt.Println("reader goroutine is exiting... connID=", c.ConnID, "remoteAddr=", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包对象

		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read head msg err", err)
			break
		}

		// 拆包 得到msgID 和MsgDalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		// 根据datalen 再次读取Data 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetMsg(data)
		fmt.Println("receive clint msgId=", msg.GetMsgId(), "msgLen=", msg.GetMsgLen(), "msg=", string(msg.GetMsg()))

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中，找到注册绑定的Conn对应的router调用
		// 根据绑定好的MsgID 找到对应处理api业务 执行
		go c.MsgHandler.DoMsgHandler(&req)
	}

}

// 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msd id=", msgId)
		return errors.New("pack error msg")
	}

	// 将数据发送给msgChan
	c.msgChan <- binaryMsg

	return nil
}

// 写消息的goroutine 专门发送给客户端消息的模块
func (c *Connection) StartWriter()  {
	fmt.Println("write goroutine is running...")
	defer fmt.Println("write", c.RemoteAddr().String(), "exit...")

	// 不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <- c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil{
				fmt.Println("send data error", err)
				return
			}
		case <- c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// 启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("conn start... connID=", c.ConnID)

	// 启动从当前链接的 读取 数据的业务
	go c.StartReader()

	// 启动从当前链接的 写 数据的业务
	go c.StartWriter()
}

// 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("conn stop... connID=", c.ConnID)

	if c.isClosed == true {
		return
	}

	// 关闭socker链接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接的 ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
