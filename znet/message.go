package znet

import "zinx/ziface"

type Message struct {
	MsgId uint32		// 消息的ID
	MsgLen uint32		// 消息的长度
	Msg []byte			// 消息的内容
}

func NewMessage(id uint32, data []byte) ziface.IMessage  {
	m := &Message{
		MsgId:id,
		MsgLen:uint32(len(data)),
		Msg:data,
	}
	return m
}


// 获取消息的ID
func (m *Message) GetMsgId() uint32{
	return m.MsgId
}
// 获取消息的长度
func (m *Message) GetMsgLen() uint32{
	return m.MsgLen
}
// 获取消息的内容
func (m *Message) GetMsg() []byte{
	return m.Msg
}

// 设置消息的ID
func (m *Message) SetMsgId(id uint32){
	m.MsgId = id
}
// 设置消息的内容
func (m *Message) SetMsg(Msg []byte){
	 m.Msg = Msg
}
// 设置消息的长度
func (m *Message) SetMsgLen(len uint32){
	m.MsgLen = len
}