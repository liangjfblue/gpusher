/**
 *
 * @author liangjf
 * @create on 2020/6/8
 * @version 1.0
 */
package pool

import (
	"context"
	"errors"
	"net"
	"sync"
)

var (
	ErrPoolParam       = errors.New("param is error")
	ErrPoolBilderNil   = errors.New("builder is nil")
	ErrPoolBuilder     = errors.New("pool builder is error")
	ErrPoolConnIsNil   = errors.New("pool conn is nil")
	ErrPoolBuilderConn = errors.New("pool builder conn error")
	ErrPoolNoIdeConn   = errors.New("pool no ide conn")
)

type Builder func(context.Context) (net.Conn, error)

type poolChannel struct {
	opts Options
	sync.RWMutex
	conns chan net.Conn
	// net.Conn builder
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewPool(opts ...Option) (IPool, error) {
	p := new(poolChannel)
	p.opts = defaultOptions

	for _, o := range opts {
		o(&p.opts)
	}

	p.conns = make(chan net.Conn, p.opts.initCap)
	p.ctx, p.cancelFunc = context.WithCancel(context.TODO())

	if err := p.initOther(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *poolChannel) initOther() error {
	if p.opts.initCap < 0 || p.opts.maxCap < 1 || p.opts.maxIdle > p.opts.maxCap {
		return ErrPoolParam
	}

	if p.opts.builder == nil {
		return ErrPoolBilderNil
	}

	for i := 0; i < p.opts.initCap; i++ {
		conn, err := p.opts.builder(p.ctx)
		if err != nil {
			p.Close()
			return ErrPoolBuilder
		}
		p.conns <- conn
	}

	return nil
}

func (p *poolChannel) put(conn net.Conn) {
	if conn == nil {
		return
	}

	if p.conns == nil {
		conn.Close()
		return
	}

	select {
	case p.conns <- conn:
	default:
		conn.Close()
		return
	}
}

func (p *poolChannel) Get() (net.Conn, error) {
	p.RLock()
	defer p.RUnlock()

	if p.conns == nil {
		return nil, ErrPoolConnIsNil
	}

	select {
	case c := <-p.conns:
		return p.WrapConn(c), nil
	default:
		conn, err := p.opts.builder(p.ctx)
		if err != nil {
			return nil, ErrPoolBuilderConn
		}
		if len(p.conns) >= p.opts.maxCap {
			return nil, ErrPoolNoIdeConn
		}
		p.conns <- conn
		return conn, nil
	}
}

func (p *poolChannel) Close() {
	p.Lock()
	defer p.Unlock()

	p.cancelFunc()
	p.opts.builder = nil

	if p.conns == nil {
		return
	}

	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
}
