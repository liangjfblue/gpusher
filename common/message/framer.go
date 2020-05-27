/**
 *
 * @author liangjf
 * @create on 2020/5/26
 * @version 1.0
 */
package message

import (
	"errors"
	"io"
	"net"

	"github.com/liangjfblue/gpusher/common/codec"
)

var (
	ErrMagicIsError     = errors.New("header magic is error")
	ErrMaxPayloadLength = errors.New("payload is too max")
	ErrReadLessThenNeed = errors.New("read less then need length")
)

const (
	DefaultPayloadLength = 1024
	MaxPayloadLength     = 4 * 1024 * 1024
)

type IFramer interface {
	ReadFramer() ([]byte, error)
}

type framer struct {
	rawConn net.Conn
}

func NewFramer(rawConn net.Conn) IFramer {
	return &framer{
		rawConn: rawConn,
	}
}

func (f *framer) ReadFramer() ([]byte, error) {
	var (
		err error
		n   int
	)

	framerHead := make([]byte, codec.FrameHeadLen)
	if n, err = io.ReadFull(f.rawConn, framerHead); err != nil {
		return nil, err
	} else if n != codec.FrameHeadLen {
		return nil, ErrReadLessThenNeed
	}

	if !codec.CheckMagic(framerHead) {
		return nil, ErrMagicIsError
	}

	dl := codec.GetDataLen(framerHead)

	if dl > MaxPayloadLength {
		return nil, ErrMaxPayloadLength
	}

	payload := make([]byte, dl)
	if n, err = io.ReadFull(f.rawConn, payload); err != nil {
		return nil, err
	} else if uint32(n) != dl {
		return nil, ErrReadLessThenNeed
	}

	dataPack := append(framerHead, payload...)
	return dataPack, nil
}
