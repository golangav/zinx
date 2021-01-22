package znet

import (
	"testing"
	"net"
	"fmt"
	"io"
)

// 封包 拆包的单元测试
func TestDataPack(t *testing.T)  {

	/* 模拟服务器 */
	// 1. 创建一个socketTCP
	listenner,err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil{
		fmt.Println("server listen err", err)
		return
	}

	go func() {
		// 2. 从客户端读取数据，拆包处理
		for {
			conn,err := listenner.Accept()
			if err != nil{
				fmt.Println("server accept error", err)
			}

			go func(con net.Conn) {
				/* 拆包处理 */
				dp := NewDataPack()
				for {
					// 1. 第一次从conn读，把包的head读出来，读取8个字节
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn,headData)
					if err != nil{
						fmt.Println("read head error")
						break
					}
					// 将8字节流的消息封装为消息头信息
					msgHead, err := dp.Unpack(headData)
					if err != nil{
						fmt.Println("server unpack err", err)
					}

					if msgHead.GetMsgLen() > 0{
						// 2. msg 是有数据的，读取msg
						msg := msgHead.(*Message)
						msg.Msg = make([]byte, msg.GetMsgLen())
						// 根据msglen的长度再次从io流中读取
						_, err := io.ReadFull(conn,msg.Msg)
						if err != nil{
							fmt.Println("server unpack data err", err)
							return
						}

						// 完整的一个消息已经读取完毕
						fmt.Println("MsgID=", msg.MsgId, "MsgLen=", msg.MsgLen, "Msg=", string(msg.Msg))
					}


				}

			}(conn)

		}
	}()

	/* 模拟客户端 */

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil{
		fmt.Println("client err", err)
		return
	}
	dp := NewDataPack()

	// 模拟粘包过程，封装两个包一同发送

	msg1 := &Message{
		MsgId:1,
		MsgLen:4,
		Msg:[]byte{'a','b','c','d'},
	}

	sendData1,err := dp.Pack(msg1)
	if err != nil{
		fmt.Println("client1 pack msg err", err)
		return
	}

	msg2 := &Message{
		MsgId:2,
		MsgLen:7,
		Msg:[]byte{'a','b','c','d','a','b','c'},
	}

	sendData2,err := dp.Pack(msg2)
	if err != nil{
		fmt.Println("client2 pack msg err", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}