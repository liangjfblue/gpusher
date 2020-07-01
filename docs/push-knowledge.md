- [x] 消息完整性
- [x] 断线重连
- [x] 连接保持
- [ ] 消息可靠性
- [ ] 消息推送速度
- [ ] 消息推送的实时性
- [ ] 消息持久化
- [x] 服务高可用
- [ ] 消息安全
- [ ] 状态监控


## 消息完整性
- 定长消息
- 特定字符分割
- 消息头+消息体(消息头为固定长度，可包含版本号，包体长度等信息，而消息体长度由所发的消息动态决定)
- 自定义协议(二进制协议等)

一般是使用[消息头+消息体]方案, 通过读取等待完整数据帧来保证数据完整性


## 断线重连
客户端支持断线重连


## 连接保持
心跳包检测连接的健康性(定时/智能时间间隔(根据消息的频率来调整心跳包的时间间隔)), 一定时间未收到服务端的ack尝试重复n次, 超过就gg


## 消息可靠性
引入ack机制, 推送消息给客户端, 客户端回复ack, 服务端一定时间未收到ack重发

### 消息发送
1.客户端处于离线状态，则直接将消息放入离线消息队列，待客户端上线后，服务端再推送消息或由客户端主动拉取消息

2.客户端在线，为了确保客户端能尽可能地收到，发送端在仅收到客户端ACK或超过规定重发次数后才停止继续发送

如果消息发送失败：

A. 如果发送超时，放入重发消息队列；
B. 如果发送异常，比如Socket异常，此时需要关闭连接，待客户端重连之后，推送再消息或由客户端主动拉取消息

如果收到客户端地ACK响应，则标记该消息已发送且被客户端成功消费

如果在指定时间内，未收到客户端地ACK，则将消息放入重发消息队列

对于重发消息队列的消息，如果发送指定次数或超过指定时间，仍未收到ACK，则将消息放入死信队列DLQ（Dead Letter Queue），并通知应用服务处理此消息

收到DLQ消息后，应用服务可以对该消息做一些处理或修复，重新放入发送队列中；也可以直接丢弃

### 消息接收
- 接收失败：比如客户端异常或App退出。此时需要在App重启并再次连接服务端之后，主动拉取消息或等待服务端推送
- 接收成功，但消费该消息时失败：此时无法向服务器发送ACK消息。在这种情况下，客户端可能会收到多条相同的消息，客户端去重便可
- 为了解决客户端接收了多条消息，在还未消费完时客户端挂掉丢消息的情况。可以在客户端加入接收消息队列（可持久化在数据库中,split3?），保证客户端消费能力较弱的情况下，也能正确消费完所有消息


## 消息推送速度    
- 1.协议的优化
- 2.实时监控推送性能
    - 1.间隔时间推送, 缓和系统的压力
    - 2.消息拆分
- 3.减小消息体的体积, 预留压缩开关, 比如超过128K就打开压缩
- 4.每次推送, 只推送消息ID,等客户端来拉取(半推半拉模式)
    

## 消息推送的实时性
- 1.消息推送失败后，立即重发。失败指定次数之后，放入重发消息队列，由其单独调度，如果超过特定时间（比如一个消息）还未发出，则交给应用服务决定如何处理该消息
- 2.采用离线消息队列：即客户端如果不在线，则将消息存储于离线消息队列。客户端上线之后，从该队列取出与其相关的消息，通过单独的服务推送出去，以避免同其他调度器发生干扰


## 消息持久化
- 需要持久化的消息主要包括:
    - 待推送的消息
    - 推送失败，需要重推的消息
    - 推送失败，需要应用服务处理的消息
    - 推送失败，需要人工处理的消息
    - 离线消息
    - 从客户端接收到的，但还没来得及处理，或处理失败的消息

- 服务重启了，按以下方式处理:
    - 在客户端连接之后，由客户端主动拉取
    - 拒绝客户端初次连接时的数据拉取请求，由服务端调度器推送。因为服务重启之后，会有大量的客户端请求连接，如果响应其拉取请求，容易造成网络拥堵


## 服务高可用
- 1.监控机器内存, cpu, 网络等指标, 随时扩容缩容
- 2.守护进程
- 3.使用集群


## 消息安全
- 1.消息加密
- 2.消息签名

> 加密、签名的顺序: 只能先签名后加密。不能对加密后的消息进行签名，这没有意义


## 状态监控
- 日志手机系统elk
- 性能指标监控prometheus+grafana

