# push

## 典型的Reactor网络模型 
- 1. 主线程master负责Acceptor监听客户端连接, 并分发fd到worker工作线程(为每个worker创建一个pipe管道通知处理函数(添加事件到event_base_loop), 通过管道来通知新fd到来
- 2. worker从管道处理函数中获取fd, 并向event_base注册户端fd的读写事件
    - 1. 设置读写对应的回调函数
    - 2. 利用客户端心跳超时机制处理半开连接
    - 3. 托管给event_base

master:

- 一个event_base
- 一个event_base_dispatch

每个worker:

- 一个event_base
- 一个event_base_dispatch

master步骤:

1. 初始化配置, 日志, 
2. 初始化redis, 刷新本地缓存逻辑服列表
3. 初始化Acceptor接收客户端连接(绑定处理分发fd到worker的回调函数)
3. 创建worker线程，用来处理来自客户端的连接
2. 初始化主线程
2. 启动启动任务线程
3. 启动监控负载线程