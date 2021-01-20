package znet

import (
	"zinx/ziface"
	"fmt"
	"net"
	"errors"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

// 定义当前客户端链接所绑定handle api
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error{
	// 回显示功能
	fmt.Println("callback...")
	if _, err := conn.Write(data[:cnt]); err != nil{
		fmt.Println("callback err", err)
		return errors.New("CallBack error")
	}
	return nil
}

// 启动服务器
func (s *Server) Start()  {

	go func() {
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil{
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil{
			fmt.Println("listen error", err)
		}
		fmt.Printf("%s starting... listen_ip=%s:%d\n", s.Name,s.IP,s.Port)

		var cid uint32
		cid = 0
		// 阻塞的等待客户端的链接，处理客户端链接业务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil{
				fmt.Println("accept err", err)
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn,cid, CallBackToClient)
			cid ++
			go dealConn.start()

		}
	}()


}

// 停止服务器
func (s *Server) Stop()  {

}

// 运行服务器
func (s *Server) Serve()  {
	// 启动server的服务功能
	s.Start()

	//TODO 做一些服务器启动之后的额外操作

	// 阻塞状态
	select {}
}

func NewServer(name string) ziface.IServer  {
	s := &Server{
		Name:name,
		IPVersion:"tcp4",
		IP:"0.0.0.0",
		Port:8999,
	}
	return s
}