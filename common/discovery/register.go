/**
 *
 * @author liangjf
 * @create on 2020/5/28
 * @version 1.0
 */
package discovery

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"

	"go.etcd.io/etcd/clientv3"
)

//ServiceDesc 服务描述信息
type ServiceDesc struct {
	//服务名称
	ServiceName string
	//ip地址
	Host string
	//端口
	Port int
	//心跳间隔 秒
	TTL time.Duration
}

//IRegister 服务注册和下线接口
type IRegister interface {
	//服务注册
	Register(context.Context, ServiceDesc) error
	//服务下线
	UnRegister(context.Context, ServiceDesc) error
}

type register struct {
	addrs       []string
	dialTimeout time.Duration
	client      *clientv3.Client
	onceDo      sync.Once
}

func NewRegister(addrs []string, dialTimeout time.Duration) IRegister {
	return &register{
		addrs:       addrs,
		dialTimeout: dialTimeout,
	}
}

func (r *register) Register(ctx context.Context, serviceDesc ServiceDesc) (err error) {
	select {
	case <-ctx.Done():
		err = errors.New("register ctx done")
		return
	default:
	}

	r.onceDo.Do(func() {
		go func() {
			c := clientv3.Config{
				Endpoints: r.addrs,
			}

			if r.dialTimeout > 0 {
				c.DialTimeout = r.dialTimeout
			}

			r.client, err = clientv3.New(c)

			var lgr *clientv3.LeaseGrantResponse
			lgr, err = r.client.Grant(ctx, int64(serviceDesc.TTL.Seconds()))
			if err != nil {
				return
			}

			serviceAddr := net.JoinHostPort(serviceDesc.Host, fmt.Sprint(serviceDesc.Port))

			//服务注册kv格式 key:/scheme/serviceName/ip:port 		value:ip:port
			//客户端监听/scheme/serviceName/, 得到/scheme/serviceName/下的可用服务地址列表, 从而得到ip:port
			serviceKey := fmt.Sprintf("/%s/%s/%s",
				scheme,
				serviceDesc.ServiceName,
				serviceAddr,
			)

			if _, err = r.client.Get(context.Background(), serviceDesc.ServiceName); err != nil {
				if err == rpctypes.ErrKeyNotFound {
					//首次设置, key不存在, 先设置
					err = nil
					if _, err = r.client.Put(ctx, serviceKey, serviceAddr, clientv3.WithLease(lgr.ID)); err != nil {
						return
					}
				} else {
					return
				}
			} else {
				if _, err = r.client.Put(ctx, serviceKey, serviceAddr, clientv3.WithLease(lgr.ID)); err != nil {
					return
				}
			}

			//续期的一半时间就续租
			ticker := time.NewTicker(time.Duration(int(serviceDesc.TTL.Seconds())/2) * time.Second)
			for {
				select {
				case <-ticker.C:
					if _, err := r.client.KeepAlive(context.TODO(), lgr.ID); err != nil {
						return
					}
				}
			}
		}()
	})
	return
}

func (r *register) UnRegister(ctx context.Context, serviceDesc ServiceDesc) (err error) {
	select {
	case <-ctx.Done():
		return errors.New("ctx done")
	default:
	}

	_, err = r.client.Delete(ctx, serviceDesc.ServiceName)
	return
}
