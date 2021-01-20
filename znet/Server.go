package znet

import (
	"zinx/ziface"
	"fmt"
	"net"
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

// 启动服务器
func (s *Server) Start()  {
	fmt.Printf("start server listen ip: %s:%d",s.IP,s.Port)

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
		fmt.Println("start Zinx server successful", s.Name)

		// 阻塞的等待客户端的链接，处理客户端链接业务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil{
				fmt.Println("accept err", err)
				continue
			}

			// 已经和客户端建立连接，做一些业务，做一个最大512字节长度的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil{
						fmt.Println("recv buf err",err)
						continue
					}

					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)
					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
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