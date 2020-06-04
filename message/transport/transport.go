/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package transport

import (
	"context"
)

type ITransport interface {
	Init(...Option)
	ListenServer(context.Context) error
}

type Option func(*Options)
