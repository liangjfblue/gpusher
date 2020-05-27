/**
 *
 * @author liangjf
 * @create on 2020/5/20
 * @version 1.0
 */
package service

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

func GetUserChannel() *ChannelBucketList {
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
	AppBucket map[int][]*ChannelBucket
}

//NewBucketList 创建channel桶
func NewBucketList(conf *config.Config) *ChannelBucketList {
	c := new(ChannelBucketList)
	c.bucketNum = conf.Channel.BucketNum

	if c.bucketNum <= 0 {
		c.bucketNum = 16
	}

	c.AppBucket = make(map[int][]*ChannelBucket)
	return c
}

//CountAll 计算gateway的客户端总数
func (c *ChannelBucketList) CountAll() (count int64) {
	for _, app := range c.AppBucket {
		for _, bucket := range app {
			count += int64(len(bucket.Data))
		}
	}
	return
}

//CountApp 计算app的客户端总数
func (c *ChannelBucketList) CountApp(appId int) (count int64) {
	if val, ok := c.AppBucket[appId]; ok {
		for _, bucket := range val {
			count += int64(len(bucket.Data))
		}
	}
	return
}

//initAppBucket 初始化桶
func (c *ChannelBucketList) initAppBucket(appId int) {
	if _, ok := c.AppBucket[appId]; !ok {
		for i := 0; i < c.bucketNum; i++ {
			c.AppBucket[appId] = append(c.AppBucket[appId], &ChannelBucket{
				Data:  make(map[string]IChannel),
				mutex: &sync.RWMutex{},
			})
		}
	}
}

//calBucket 计算bucket的位置
func (c *ChannelBucketList) calBucket(appId int, key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % len(c.AppBucket[appId])
}

//New 创建新的客户端channel
func (c *ChannelBucketList) New(appId int, key string) (IChannel, *ChannelBucket, error) {
	c.initAppBucket(appId)

	index := c.calBucket(appId, key)
	bb := c.AppBucket[appId][index]

	bb.Lock()
	defer bb.UnLock()

	if cc, ok := bb.Data[key]; ok {
		log.GetLogger(defind.GatewayLog).Debug("appId:%d, key:%s is existed", appId, key)
		return cc, bb, nil
	} else {
		log.GetLogger(defind.GatewayLog).Debug("new conn: appId:%d, key:%s", appId, key)
		cc := NewUserChannel()
		bb.Data[key] = cc

		return cc, bb, nil
	}
}

//Add 添加channel到本地缓存
func (c *ChannelBucketList) Get(appId int, key string, newFlag bool) (IChannel, error) {
	c.initAppBucket(appId)

	//计算出hash的桶
	index := c.calBucket(appId, key)
	bb := c.AppBucket[appId][index]

	bb.Lock()
	defer bb.UnLock()

	if cc, ok := bb.Data[key]; ok {
		log.GetLogger(defind.GatewayLog).Debug("appId:%d, key:%s is existed", appId, key)
		return cc, nil
	} else {
		if newFlag {
			log.GetLogger(defind.GatewayLog).Debug("new conn appId:%d, key:%s", appId, key)
			cc := NewUserChannel()
			bb.Data[key] = cc
			return cc, nil
		} else {
			return nil, ErrChannelBucketChannelNotExist
		}
	}
}

//Remove 删除channel本地缓存
func (c *ChannelBucketList) Remove(appId int, key string) error {
	c.initAppBucket(appId)

	//计算出hash的桶
	index := c.calBucket(appId, key)
	bb := c.AppBucket[appId][index]

	bb.Lock()
	defer bb.UnLock()

	if _, ok := bb.Data[key]; ok {
		log.GetLogger(defind.GatewayLog).Debug("remove channel, appId:%d, key:%s", appId, key)
		delete(bb.Data, key)
		return nil
	} else {
		log.GetLogger(defind.GatewayLog).Warn("key channel not exist, appId:%d, key:%s", appId, key)
		return ErrChannelBucketChannelNotExist
	}
}

//Close 关闭本地缓存
func (c *ChannelBucketList) Close() {
	channelTmp := make([]IChannel, 0, c.CountAll())

	for _, app := range c.AppBucket {
		for _, bucket := range app {
			bucket.RLock()
			for _, channel := range bucket.Data {
				channelTmp = append(channelTmp, channel)
			}
			bucket.RUnlock()
		}
	}

	for _, channel := range channelTmp {
		if err := channel.Close(); err != nil {
			log.GetLogger(defind.GatewayLog).Error("channel close err:%v", err.Error())
		}
	}
}
