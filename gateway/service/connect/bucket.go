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

	"github.com/liangjfblue/gpusher/gateway/config"
	"github.com/liangjfblue/gpusher/gateway/defind"
)

var (
	ErrChannelBucketChannelNotExist = errors.New("channel not exist")
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
	Buckets   []*ChannelBucket
}

//NewBucketList 创建channel桶
func NewBucketList(conf *config.Config) *ChannelBucketList {
	c := new(ChannelBucketList)
	c.bucketNum = conf.Channel.BucketNum

	if c.bucketNum <= 0 {
		c.bucketNum = 16
	}

	c.initAppBucket()

	return c
}

//initAppBucket 初始化桶
func (c *ChannelBucketList) initAppBucket() {
	c.Buckets = make([]*ChannelBucket, 0, c.bucketNum)
	for i := 0; i < c.bucketNum; i++ {
		c.Buckets = append(c.Buckets, &ChannelBucket{
			Data:  make(map[string]IChannel),
			mutex: &sync.RWMutex{},
		})
	}
}

//CountAll 计算gateway的客户端总数
func (c *ChannelBucketList) CountAll() (count int64) {
	for _, b := range c.Buckets {
		count += int64(len(b.Data))
	}
	return
}

//calBucket 计算bucket的位置
func (c *ChannelBucketList) calBucket(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % len(c.Buckets)
}

//New 创建新的客户端channel
func (c *ChannelBucketList) New(key string) (IChannel, *ChannelBucket, error) {
	index := c.calBucket(key)
	bb := c.Buckets[index]

	bb.Lock()
	defer bb.UnLock()

	if cc, ok := bb.Data[key]; ok {
		log.GetLogger(defind.GatewayLog).Debug("key:%s is existed", key)
		return cc, bb, nil
	} else {
		log.GetLogger(defind.GatewayLog).Debug("new conn: key:%s", key)
		cc := NewConnChannel()
		bb.Data[key] = cc

		return cc, bb, nil
	}
}

//Add 添加channel到本地缓存
func (c *ChannelBucketList) Get(key string, newFlag bool) (IChannel, error) {
	//计算出hash的桶
	index := c.calBucket(key)
	bb := c.Buckets[index]

	bb.Lock()
	defer bb.UnLock()

	if cc, ok := bb.Data[key]; ok {
		log.GetLogger(defind.GatewayLog).Debug("key:%s is existed", key)
		return cc, nil
	} else {
		if newFlag {
			log.GetLogger(defind.GatewayLog).Debug("new conn key:%s", key)
			cc := NewConnChannel()
			bb.Data[key] = cc
			return cc, nil
		} else {
			return nil, ErrChannelBucketChannelNotExist
		}
	}
}

//Remove 删除channel本地缓存
func (c *ChannelBucketList) Remove(key string) error {
	//计算出hash的桶
	index := c.calBucket(key)
	bb := c.Buckets[index]

	bb.Lock()
	defer bb.UnLock()

	if _, ok := bb.Data[key]; ok {
		log.GetLogger(defind.GatewayLog).Debug("remove channel,key:%s", key)
		delete(bb.Data, key)
		return nil
	} else {
		log.GetLogger(defind.GatewayLog).Warn("key channel not exist, key:%s", key)
		return ErrChannelBucketChannelNotExist
	}
}

//Close 关闭本地缓存
func (c *ChannelBucketList) Close() {
	channelTmp := make([]IChannel, 0, c.CountAll())

	for _, b := range c.Buckets {
		b.RLock()
		for _, channel := range b.Data {
			channelTmp = append(channelTmp, channel)
		}
		b.RUnlock()
	}

	for _, channel := range channelTmp {
		if err := channel.Close(); err != nil {
			log.GetLogger(defind.GatewayLog).Error("channel close err:%v", err.Error())
		}
	}
}
