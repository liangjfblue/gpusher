#!/bin/bash

#start zookeeper
echo "start zookeeper"
cd /home/liangjf/app/zookeeper/apache-zookeeper-3.6.0-bin
/home/liangjf/app/zookeeper/apache-zookeeper-3.6.0-bin/start-cluster.sh
echo "start zookeeper ok"


#start kafka
echo "start kafka"
cd /home/liangjf/app/kafka/kafka_2.11-2.4.1
/home/liangjf/app/kafka/kafka_2.11-2.4.1/start-cluster.sh
echo "start kafka ok"


#start redis
echo "start redis"
cd /home/liangjf/app/redis-5.0.5
/home/liangjf/app/redis-5.0.5/script/start-cluster.sh
echo "start redis ok"


#start etcd
echo "start etcd"
cd /home/liangjf/app/etcd-v3.3.13-linux-amd64
/home/liangjf/app/etcd-v3.3.13-linux-amd64/start-cluster.sh
echo "start etcd ok"
