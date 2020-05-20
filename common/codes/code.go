/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package codes

import "fmt"

var (
	_codes = map[int32]struct{}{}
)

func New(e int32) int32 {
	if e <= 0 {
		panic("business code must greater than zero")
	}
	return add(e)
}

func add(e int32) int32 {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("code: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return e
}

type Errno struct {
	Code int32       `json:"Code"`
	Msg  string      `json:"Msg"`
	Data interface{} `json:"Data,omitempty"`
}

func (e Errno) Error() string {
	return fmt.Sprintf("code:%d msg:%s", e.Code, e.Msg)
}

func NewErr(code int32, msg string) *Errno {
	return &Errno{
		Code: New(code),
		Msg:  msg,
	}
}

var (
	Success = NewErr(CodeSuccess, "ok")

	ErrNetworkNotSupported = NewErr(CodeNetworkNotSupported, "network type not supported")
)
