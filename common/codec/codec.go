/**
 *
 * @author liangjf
 * @create on 2020/5/26
 * @version 1.0
 */
package codec

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrDataEmpty = errors.New("data is empty")
)

//ICodec 编解码
type ICodec interface {
	Encode(*FrameHeader, []byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

const (
	FrameHeadLen = 15
	Magic        = 0x45
	Version      = 0
)

const (
	GeneralMsg   = 0x00
	HeartbeatMsg = 0x01
)

//FrameHeader 帧头
type FrameHeader struct {
	Magic        uint8  //1 magic
	Version      uint8  //1 version
	MsgType      uint8  //1 msg type, 0x0: general req,  0x1: heartbeat
	ReqType      uint8  //1 request type, 0x0: send and receive,   0x1: send but not receive
	CompressType uint8  //1 compression or not,  0x0: not compression,  0x1: compression
	StreamID     uint16 //2 stream ID
	Length       uint32 //4 total packet length
	Reserved     uint32 //4 4 bytes reserved
}

var (
	_codecMap    map[string]ICodec
	DefaultCodec = NewDefaultCodec()
)

const (
	Default = "default"
)

func init() {
	RegisterCodec("default", DefaultCodec)
}

func RegisterCodec(name string, codec ICodec) {
	if _codecMap == nil {
		_codecMap = make(map[string]ICodec)
	}
	_codecMap[name] = codec
}

func GetCodec(name string) ICodec {
	if codec, ok := _codecMap[name]; ok {
		return codec
	}
	return DefaultCodec
}

//defaultCodec 默认编解码器
type defaultCodec struct {
}

func NewDefaultCodec() ICodec {
	return &defaultCodec{}
}

func (c *defaultCodec) Encode(header *FrameHeader, d []byte) ([]byte, error) {
	if d == nil {
		return nil, ErrDataEmpty
	}

	var (
		err error
	)

	framerSize := FrameHeadLen + len(d)
	b := make([]byte, 0, framerSize)
	buf := bytes.NewBuffer(b)

	framerHeader := FrameHeader{
		Magic:        Magic,
		Version:      Version,
		MsgType:      0x0,
		ReqType:      0x0,
		CompressType: 0x0,
		Length:       uint32(len(d)),
		Reserved:     0,
	}

	if header.MsgType != 0 {
		framerHeader.MsgType = 0x01
	}

	if header.ReqType != 0 {
		framerHeader.ReqType = 0x01
	}

	if header.CompressType != 0 {
		framerHeader.CompressType = 0x01
	}

	if header.Reserved != 0 {
		framerHeader.Reserved = 0x01
	}

	if err = binary.Write(buf, binary.BigEndian, framerHeader.Magic); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.BigEndian, framerHeader.Version); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.BigEndian, framerHeader.MsgType); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.BigEndian, framerHeader.ReqType); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.BigEndian, framerHeader.CompressType); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.BigEndian, framerHeader.Length); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, framerHeader.Reserved); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.BigEndian, d); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *defaultCodec) Decode(framer []byte) ([]byte, error) {
	return framer[FrameHeadLen:], nil
}

func CheckMagic(framer []byte) bool {
	return framer[0] == Magic
}

func GetDataLen(framer []byte) uint32 {
	return binary.BigEndian.Uint32(framer[7:11])
}

func IsHeartBeatMsg(framer []byte) bool {
	return true //framer[2] == HeartbeatMsg
}

func GetVersion(framer []byte) uint32 {
	return Version // binary.BigEndian.Uint32(framer[1:2])
}
