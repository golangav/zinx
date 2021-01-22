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
	fmt.Println("call router handle...")

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil{
		fmt.Println(err)
	}

}


func main() {
	s := znet.NewServer()

	// 给当前server增加一个自定义的router
	s.AddRouter(&PingRouter{})
	s.Serve()
}
