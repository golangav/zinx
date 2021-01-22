package main

import (
	"fmt"
	"time"
	"net"
	"zinx/znet"
	"io"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("clint start...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil{
		fmt.Println("client start err")
		return
	}

	for {
		// 发送数据
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("zinx cliet send data")))
		if err != nil{
			fmt.Println("pack error", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil{
			fmt.Println("write error")
			return
		}

		// 接收数据
		binaryHead := make([]byte,dp.GetHeadLen())
		if _, err := io.ReadFull(conn,binaryHead); err != nil{
			fmt.Println("read head error", err)
			break
		}
		msgHead,err := dp.Unpack(binaryHead)
		if err != nil{
			fmt.Println("client unpack msgHead error", err)
			break
		}
		if msgHead.GetMsgLen() >0{
			msg := msgHead.(*znet.Message)
			msg.Msg= make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn,msg.Msg); err != nil{
				fmt.Println("read msg data error")
				return
			}

			fmt.Println("rece server msgId=", msg.MsgId, "msgLen=", msg.MsgLen, "msg=",string(msg.Msg))
		}

		time.Sleep(1 * time.Second)


	}
}
