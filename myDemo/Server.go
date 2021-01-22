package main

import (
	"zinx/znet"
	"zinx/ziface"
	"fmt"
)

// 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) Handle(request ziface.IRequest){
	fmt.Println("call ping router handle...")

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil{
		fmt.Println(err)
	}

}

type TelnetRouter struct {
	znet.BaseRouter
}

func (r *TelnetRouter) Handle(request ziface.IRequest){
	fmt.Println("call telnet router handle...")

	err := request.GetConnection().SendMsg(400, []byte("telnet...telnet...telnet"))
	if err != nil{
		fmt.Println(err)
	}

}


func main() {
	s := znet.NewServer()

	// 给当前server增加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &TelnetRouter{})
	s.Serve()
}
