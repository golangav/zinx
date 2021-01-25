package znet

import (
	"zinx/ziface"
	"strconv"
	"fmt"
	"zinx/utils"
)

/*
  消息处理模块的实现
*/

type MsgHandel struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32] ziface.IRouter
	// 负责work取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Work池的work数量
	WorkerPoolSize uint32
}

func NewMsgHandel() ziface.IMsgHandle{
	return &MsgHandel{
		Apis:make(map[uint32] ziface.IRouter),
		WorkerPoolSize: utils.ConfigObj.WorkerPoolSize,
		TaskQueue: make([]chan ziface.IRequest, utils.ConfigObj.WorkerPoolSize),
	}
}

// 调度、执行对应的Router消息处理方法
func (mh *MsgHandel) DoMsgHandler(request ziface.IRequest){
	// 1. 从request中找点msgID
	handle,ok := mh.Apis[request.GetMsgId()]
	if !ok{
		fmt.Println("api not found... msgID=",request.GetMsgId())
	}
	// 2. 根据MsgID 调度对应的router业务即可
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}


// 为消息增加具体的处理逻辑
func (mh *MsgHandel) AddRouter(msgID uint32, router ziface.IRouter){
	// 1. 判断 当前msg绑定的api处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok{
		// id 已经存在
		panic("repeat api msgID=" + strconv.Itoa(int(msgID)))
	}

	// 2. 增加msg与api的绑定关系
	mh.Apis[msgID] = router
}

// 启动一个work工作池
func (mh *MsgHandel) StartWorkerPool(){
	for i :=0; i < int(mh.WorkerPoolSize); i++{
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.ConfigObj.MaxWorkerTaskLen)
		go mh.StartWorkerOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个work工作流程
func (mh *MsgHandel) StartWorkerOneWorker(workID int, taskQueue chan ziface.IRequest){
	fmt.Println("work ID=", workID, "is started...")

	for{
		select {
			case request := <- taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

// 将消息发送给消息任务队列处理
func (mh *MsgHandel) SendMsgToTaskQueque(request ziface.IRequest) {
	workID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workID] <- request
}





