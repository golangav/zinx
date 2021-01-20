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


func (r *PingRouter) PreHandle(request ziface.IRequest){
	fmt.Println("pre handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil{
		fmt.Println("before ping error")
	}
}

func (r *PingRouter) Handle(request ziface.IRequest){
	fmt.Println("handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("handle... ping...\n"))
	if err != nil{
		fmt.Println("handle... ping error")
	}
}

func (r *PingRouter) PostHandle(request ziface.IRequest){
	fmt.Println("after handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil{
		fmt.Println("after handle...")
	}
}


func main() {
	s := znet.NewServer("[zinx v3.0]")

	// 给当前server增加一个自定义的router
	s.AddRouter(&PingRouter{})
	s.Serve()
}
