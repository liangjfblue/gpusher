/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

type IFactory interface {
	CreateTransport() ITransport
}

type FactoryTcpTransport struct{}

func NewFactoryTcpTransport(opts ...Option) ITransport {
	return NewTcpTransport(opts...)
}

type FactoryWSTransport struct{}

func NewFactoryWSTransport(opts ...Option) ITransport {
	return NewWSTransport(opts...)
}
