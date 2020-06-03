/**
 *
 * @author liangjf
 * @create on 2020/6/2
 * @version 1.0
 */
package push

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var (
	ErrChannelIsClosed = errors.New("channel is closed")
)

//KafkaSender kafka发送者
type KafkaSender struct {
	_syncProducer  sarama.SyncProducer
	_asyncProducer sarama.AsyncProducer

	isSync      bool
	brokerAddrs []string
	sync.RWMutex

	stopChan    chan struct{}
	sendChannel chan *PushMsg
}

func NewKafkaSender(brokerAddrs []string, isSync bool) IQueueSender {
	return &KafkaSender{
		isSync:      isSync,
		brokerAddrs: brokerAddrs,
		sendChannel: make(chan *PushMsg),
		stopChan:    make(chan struct{}, 1),
	}
}

func (q *KafkaSender) Init() error {
	err := q.initKafkaLog()

	go q.sendRun()

	return err
}

func (q *KafkaSender) Send(msg *PushMsg) error {
	q.toChannel(msg)
	return nil
}

func (q *KafkaSender) Stop() {
	q.Lock()
	defer q.Unlock()

	if len(q.stopChan) > 0 {
		return
	}

	q.stopChan <- struct{}{}
}

func (q *KafkaSender) initKafkaLog() error {
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Flush.Frequency = 100 * time.Millisecond

	if q.isSync {
		q._syncProducer, err = sarama.NewSyncProducer(q.brokerAddrs, config)
		if err != nil {
			return err
		}
	} else {
		q._asyncProducer, err = sarama.NewAsyncProducer(q.brokerAddrs, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *KafkaSender) sendRun() {
	if q.isSync {
		q.sendSync()
	} else {
		q.sendAsync()
	}
}

func (q *KafkaSender) sendSync() {
	defer func() {
		close(q.sendChannel)
		close(q.stopChan)

		if err := q._syncProducer.Close(); err != nil {
			log.Fatalf("gpusher:  close producer fail, because is %s\n", err.Error())
		}
	}()

	for {
		select {
		case msg, ok := <-q.sendChannel:
			if !ok {
				log.Fatal("gpusher: chan is closed")
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("gpusher: json body error")
				continue
			}

			fmt.Println("------------------")
			fmt.Println(string(data))
			fmt.Println("------------------")
			kafkaMsg := sarama.ProducerMessage{
				Topic: msg.Tag,
				Value: sarama.StringEncoder(data),
				//Key:   sarama.StringEncoder(""),
			}

			if _, _, err := q._syncProducer.SendMessage(&kafkaMsg); err != nil {
				log.Printf("gpusher: err:%s\n", err.Error())
				return
			}
		case <-q.stopChan:
			return
		}
	}
}

func (q *KafkaSender) sendAsync() {
	defer func() {
		close(q.sendChannel)
		close(q.stopChan)

		if err := q._asyncProducer.Close(); err != nil {
			log.Fatalf("gpusher: close producer fail, because is %s\n", err.Error())
		}
	}()

	for {
		select {
		case msg, ok := <-q.sendChannel:
			if !ok {
				log.Fatal("gpusher: chan is closed")
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("gpusher: json body error")
				continue
			}

			kafkaMsg := sarama.ProducerMessage{
				Topic: msg.Tag,
				Value: sarama.StringEncoder(data),
				//Key:   sarama.StringEncoder(""),
			}

			q._asyncProducer.Input() <- &kafkaMsg

			select {
			case <-q._asyncProducer.Successes():
				continue
			case err := <-q._asyncProducer.Errors():
				//TODO retry to send log to kafka
				log.Fatal("gpusher: produced message error: ", err)
				return
			default:
			}
		case <-q.stopChan:
			log.Println("gpusher: kafka stop")
			return
		}
	}
}

func (q *KafkaSender) toChannel(msg *PushMsg) {
	select {
	case q.sendChannel <- msg:
	default:
	}
}

//KafkaReceiver kafka接收者
type KafkaReceiver struct {
	brokerAddrs []string
	sync.RWMutex

	stopChan    chan struct{}
	recvChannel chan []byte
}

func NewKafkaReceiver(brokerAddrs []string) IQueueReceiver {
	return &KafkaReceiver{
		brokerAddrs: brokerAddrs,
		recvChannel: make(chan []byte),
		stopChan:    make(chan struct{}, 1),
	}
}

func (q *KafkaReceiver) Init() error {
	q.recvRun()
	return nil
}

func (q *KafkaReceiver) Recv(f func([]byte)) error {
	defer func() {
		close(q.recvChannel)
		close(q.stopChan)
	}()

	for {
		select {
		case msg, ok := <-q.recvChannel:
			if !ok {
				log.Fatal("gpusher: channel is closed")
				return ErrChannelIsClosed
			}

			f(msg)

		case <-q.stopChan:
			return nil
		}
	}
}

func (q *KafkaReceiver) Stop() {
	q.Lock()
	defer q.Unlock()

	if len(q.stopChan) > 0 {
		return
	}

	q.stopChan <- struct{}{}
}

func (q *KafkaReceiver) recvRun() {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Version = sarama.V2_4_0_0

	consumerT, err := sarama.NewConsumer(q.brokerAddrs, config)
	if err != nil {
		log.Fatal("gpusher: new kafka consumer err: ", err.Error())
		return
	}

	topics, err := consumerT.Topics()
	if err != nil {
		log.Fatal("gpusher: kafka Topics err: ", err.Error())
		return
	}

	topicP := make(map[string]int)
	for _, topic := range topics {

		//推送主题必须是 app_AppName的格式
		if !strings.HasPrefix(topic, "app_") {
			continue
		}

		ps, err := consumerT.Partitions(topic)
		if err != nil {
			continue
		}
		topicP[topic] = len(ps)
	}
	consumerT.Close()

	for topic := range topicP {
		consumer, err := sarama.NewConsumer(q.brokerAddrs, config)
		if err != nil {
			log.Printf("gpusher: Create consumer for %s topic failed: %s\n", topic, err.Error())
			return
		}

		partitions, err := consumer.Partitions(topic)
		if err != nil {
			log.Printf("Get %s topic's paritions failed: %s\n", topic, err.Error())

		}

		for _, p := range partitions {
			pc, e := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
			if e != nil {
				log.Printf("Start consumer for partition failed %d: %s\n", p, topic)
				return
			}
			go func(pc sarama.PartitionConsumer) {
				for msg := range pc.Messages() {
					q.recvChannel <- msg.Value
				}

				q.Stop()
			}(pc)
		}
	}
}
