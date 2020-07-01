/**
 *
 * @author liangjf
 * @create on 2020/7/1
 * @version 1.0
 */
package models

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/liangjfblue/gpusher/logic/common"
)

var (
	_etcdClient *clientv3.Client
)

func InitEtcd(etcdAddr []string) error {
	var err error
	_etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints: etcdAddr,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetAllGateway() ([]string, error) {
	var (
		err         error
		gatewayAdds []string
	)
	rp, err := _etcdClient.Get(context.TODO(), common.GatewayPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, kv := range rp.Kvs {
		gatewayAdds = append(gatewayAdds, string(kv.Value))
	}
	return gatewayAdds, nil
}
