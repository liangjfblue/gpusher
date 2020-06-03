/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package connect

import (
	"errors"
	"hash/crc32"
	"sync"

	"github.com/liangjfblue/gpusher/common/logger/log"

	"github.com/liangjfblue/gpusher/gateway/common"
	"github.com/liangjfblue/gpusher/gateway/config"
)

var (
	ErrChannelBucketChannelNotExist = errors.New("channel not exist")
	ErrAppNotExist                  = errors.New("app not exist")
)

var (
	onceDo         sync.Once
	_clientChannel *ChannelBucketList
)

func InitClientChannel(conf *config.Config) {
	onceDo.Do(func() {
		_clientChannel = NewBucketList(conf)
	})
}

func GetClientChannel() *ChannelBucketList {
	return _clientChannel
}

//ChannelBucket a bucket
type ChannelBucket struct {
	Data  map[string]IChannel
	mutex *sync.RWMutex
}

func (c *ChannelBucket) Lock() {
	c.mutex.Lock()
}

func (c *ChannelBucket) UnLock() {
	c.mutex.Unlock()
}

func (c *ChannelBucket) RLock() {
	c.mutex.RLock()
}

func (c *ChannelBucket) RUnlock() {
	c.mutex.RUnlock()
}

//ChannelBucketList channel桶
type ChannelBucketList struct {
	bucketNum int
	Buckets   map[int][]*ChannelBucket //每个app分桶, 提高查询效率
}

//NewBucketList 创建channel桶
func NewBucketList(conf *config.Config) *ChannelBucketList {
	c := new(ChannelBucketList)
	c.bucketNum = conf.Channel.BucketNum

	if c.bucketNum <= 0 {
		c.bucketNum = 16
	}

	c.Buckets = make(map[int][]*ChannelBucket)

	return c
}

//initAppBucket 初始化桶
func (c *ChannelBucketList) initAppBucket(appId int) {
	if _, ok := c.Buckets[appId]; !ok {
		c.Buckets[appId] = make([]*ChannelBucket, 0, c.bucketNum)

		for i := 0; i < c.bucketNum; i++ {
			c.Buckets[appId] = append(c.Buckets[appId], &ChannelBucket{
				Data:  map[string]IChannel{},
				mutex: &sync.RWMutex{},
			})
		}
	}
}

//CountAll 计算gateway的客户端总数
func (c *ChannelBucketList) CountAll() (count int64) {
	for _, apps := range c.Buckets {
		for _, app := range apps {
			count += int64(len(app.Data))
		}
	}
	return
}

//CountApp 计算gateway的某app客户端总数
func (c *ChannelBucketList) CountApp(appId int) (count int64) {
	if app, ok := c.Buckets[appId]; ok {
		for _, b := range app {
			count += int64(len(b.Data))
		}
	}
	return
}

//calBucket 计算bucket的位置
func (c *ChannelBucketList) calBucket(appId int, key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % len(c.Buckets[appId])
}

//New 创建新的客户端channel
func (c *ChannelBucketList) New(appId int, uuid string) (IChannel, *ChannelBucket, error) {
	c.initAppBucket(appId)

	index := c.calBucket(appId, uuid)
	bb := c.Buckets[appId][index]

	bb.Lock()
	defer bb.UnLock()

	if cc, ok := bb.Data[uuid]; !ok {
		nc := NewConnChannel()
		bb.Data[uuid] = nc

		log.GetLogger(common.GatewayLog).Debug("new conn: uuid:%s", uuid)
		return nc, bb, nil
	} else {
		log.GetLogger(common.GatewayLog).Debug("uuid:%s is existed", uuid)
		return cc, bb, nil
	}
}

//Get 获取客户端通道
func (c *ChannelBucketList) Get(appId int, uuid string, newFlag bool) (IChannel, error) {
	c.initAppBucket(appId)

	index := c.calBucket(appId, uuid)
	bb := c.Buckets[appId][index]

	bb.Lock()
	defer bb.UnLock()

	if cc, ok := bb.Data[uuid]; ok {
		return cc, nil
	} else {
		if newFlag {
			log.GetLogger(common.GatewayLog).Debug("new conn: uuid:%s", uuid)
			nc := NewConnChannel()
			bb.Data[uuid] = nc
			return nc, nil
		} else {
			return nil, ErrChannelBucketChannelNotExist
		}
	}
}

//GetApp 获取某app客户端通道
func (c *ChannelBucketList) GetApp(appId int) ([]IChannel, error) {
	if _, ok := c.Buckets[appId]; !ok {
		return nil, ErrAppNotExist
	}

	cc := make([]IChannel, 0)
	for _, bc := range c.Buckets[appId] {
		bc.RLock()
		for _, cs := range bc.Data {
			cc = append(cc, cs)
		}
		bc.RUnlock()
	}
	return cc, nil
}

//GetApp 获取某app客户端通道
func (c *ChannelBucketList) GetAll() ([]IChannel, error) {
	cc := make([]IChannel, 0)
	for _, app := range c.Buckets {
		for _, bc := range app {
			bc.RLock()
			for _, cs := range bc.Data {
				cc = append(cc, cs)
			}
			bc.RUnlock()
		}
	}
	return cc, nil
}

//Remove 删除channel本地缓存
func (c *ChannelBucketList) Remove(appId int, uuid string) error {
	if _, ok := c.Buckets[appId]; !ok {
		log.GetLogger(common.GatewayLog).Warn("appId not exist, appId:%d, uuid:%s", appId, uuid)
		return ErrChannelBucketChannelNotExist
	}

	index := c.calBucket(appId, uuid)
	bb := c.Buckets[appId][index]

	bb.Lock()
	defer bb.UnLock()

	if _, ok := bb.Data[uuid]; ok {
		log.GetLogger(common.GatewayLog).Debug("remove channel, appId:%d, uuid:%s", appId, uuid)
		delete(bb.Data, uuid)
		return nil
	} else {
		log.GetLogger(common.GatewayLog).Warn("channel not exist, appId:%d, uuid:%s", appId, uuid)
		return ErrChannelBucketChannelNotExist
	}
}

//Close 关闭本地缓存
func (c *ChannelBucketList) Close() {
	channels := make([]IChannel, 0, c.CountAll())

	for _, app := range c.Buckets {
		for _, bs := range app {
			bs.RLock()
			for _, cs := range bs.Data {
				channels = append(channels, cs)
			}
			bs.RUnlock()
		}
	}

	for _, channel := range channels {
		if err := channel.Close(); err != nil {
			log.GetLogger(common.GatewayLog).Error("channel close err:%v", err.Error())
		}
	}
}
