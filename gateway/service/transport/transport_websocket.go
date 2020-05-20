/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package transport

import "context"

type wsTransport struct {
	opts Options
}

func (t *wsTransport) Init(opts ...Option) {
	for _, o := range opts {
		o(&t.opts)
	}
}
func (t *wsTransport) ListenServer(ctx context.Context) error {

	return nil
}

func NewWSTransport(opts ...Option) ITransport {
	t := new(wsTransport)
	t.opts = defaultOptions

	for _, o := range opts {
		o(&t.opts)
	}

	return t
}
