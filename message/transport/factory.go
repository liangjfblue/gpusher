/**
 *
 * @author liangjf
 * @create on 2020/6/4
 * @version 1.0
 */
package transport

type IFactory interface {
	CreateTransport() ITransport
}

type FactoryRPCTransport struct{}

func NewFactoryRPCTransport(opts ...Option) ITransport {
	return NewRpcTransport(opts...)
}
