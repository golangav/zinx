package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
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
	// 给当前的Server添加一个router，server注册的链接对应的处理业务
	Router ziface.IRouter
}

// 启动服务器
func (s *Server) Start() {

	go func() {
		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen error", err)
		}
		fmt.Printf("%s starting... listen_ip=%s:%d\n", s.Name, s.IP, s.Port)

		var cid uint32
		cid = 0
		// 阻塞的等待客户端的链接，处理客户端链接业务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept err", err)
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()

		}
	}()

}

// 停止服务器
func (s *Server) Stop() {

}

// 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	//TODO 做一些服务器启动之后的额外操作

	// 阻塞状态
	select {}
}

// 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
}

func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.ConfigObj.Name,
		IPVersion: "tcp4",
		IP:        utils.ConfigObj.Host,
		Port:      utils.ConfigObj.Port,
		Router:    nil,
	}
	return s
}
