/**
 *
 * @author liangjf
 * @create on 2020/6/2
 * @version 1.0
 */
package push

type IQueueSender interface {
	Init() error
	Send(*PushMsg) error
	Stop()
}

type IQueueReceiver interface {
	Init() error
	Recv(func([]byte)) error
	Stop()
}
