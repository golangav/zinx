package znet

import "zinx/ziface"

// 实现routers时，先嵌入这个BaseRouter基类，然后根据需求对这个基类的方法进行重写
type BaseRouter struct {}

// 在处理conn业务之前的钩子方法Hook
func (r *BaseRouter) PreHandle(request ziface.IRequest){}

// 在处理conn业务的主方法Hook
func (r *BaseRouter) Handle(request ziface.IRequest){}

// 在处理conn业务之后的钩子防范Hook
func (r *BaseRouter) PostHandle(request ziface.IRequest){}