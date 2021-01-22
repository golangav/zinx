package znet

import (
	"zinx/ziface"
	"bytes"
	"encoding/binary"
	"zinx/utils"
	"errors"
)

// 封包 拆包的具体模块
type DataPack struct {

}

func NewDataPack()  ziface.IdataPack{
	d := &DataPack{}
	return d
}


// 获取包的头的长度方法
func (d *DataPack) GetHeadLen() uint32{
	// MsgLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// 封包方法
func (d *DataPack) Pack(msg ziface.IMessage)([]byte,error){
	// 创建一个存放byte字节流的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 将datalen 写进 dataBuf 中，unit32类型，占4个字节
	err := binary.Write(dataBuf,binary.LittleEndian,msg.GetMsgLen())
	if err != nil{
		return nil, err
	}

	// 将MagId 写进 dataBuf 中，unit32类型，占4个字节
	err = binary.Write(dataBuf,binary.LittleEndian,msg.GetMsgId())
	if err != nil{
		return nil, err
	}

	// Msg 写进 dataBuf 中
	err = binary.Write(dataBuf,binary.LittleEndian,msg.GetMsg())
	if err != nil{
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// 拆包方法
func (d *DataPack) Unpack(binaryData []byte)(ziface.IMessage, error){
	// 创建一个从输入二进制数据的ioReader
	dataBuf := bytes.NewReader(binaryData)

	// 只解压head信息，得到MsgLen 和 MsgID
	msg := &Message{}

	// 读MsgLen
	if err := binary.Read(dataBuf,binary.LittleEndian, &msg.MsgLen); err != nil{
		return nil, err
	}
	// 读MsgID
	if err := binary.Read(dataBuf,binary.LittleEndian, &msg.MsgId); err != nil{
		return nil, err
	}

	// 判断datalen 是否已经超出了我们允许的最大包长度
	if (utils.ConfigObj.MaxPackageSize > 0 && msg.MsgLen > utils.ConfigObj.MaxPackageSize) {
		return nil, errors.New("too large msg data")
	}

	return msg, nil
}

