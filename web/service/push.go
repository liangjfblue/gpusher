/**
 *
 * @author liangjf
 * @create on 2020/6/3
 * @version 1.0
 */
package service

import (
	"log"
	"sync"

	"github.com/liangjfblue/gpusher/common/push"
)

type IPush interface {
	Run()
	Push(*push.PushMsg)
	Stop()
}

type DefaultPush struct {
	msg      chan *push.PushMsg
	stop     chan struct{}
	stopFlag bool
	queue    push.IQueueSender
	sync.RWMutex
}

func NewDefaultPush(queue push.IQueueSender) IPush {
	return &DefaultPush{
		msg:      make(chan *push.PushMsg),
		stop:     make(chan struct{}),
		stopFlag: false,
		queue:    queue,
	}
}

func (p *DefaultPush) Run() {
	var (
		err error
	)
	defer close(p.msg)

	for {
		select {
		case <-p.stop:
			p.stopFlag = true
			return
		case m, ok := <-p.msg:
			if !ok {
				p.stopFlag = true
				p.stop <- struct{}{}
				return
			}

			if err = p.queue.Send(&push.PushMsg{
				Tag:  m.Tag,
				Body: m.Body,
			}); err != nil {
				log.Fatal("gpusher: send msg to push err:", err.Error())
				return
			}
		}
	}
}

func (p *DefaultPush) Stop() {
	p.RLock()
	defer p.RUnlock()
	if p.stopFlag {
		return
	}

	p.stop <- struct{}{}
}

func (p *DefaultPush) Push(m *push.PushMsg) {
	p.RLock()
	defer p.RUnlock()
	if p.stopFlag {
		return
	}

	p.msg <- m
}

var (
	_pushM map[string]IPush
)

func RegisterPush(name string, push IPush) {
	if _pushM == nil {
		_pushM = make(map[string]IPush)
	}

	_pushM[name] = push
}

func GetPush(name string) IPush {
	return _pushM[name]
}
