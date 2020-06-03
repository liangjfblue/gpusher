
# gateway(网关)
- 维护与client的连接, 转发消息到logic处理
- 接收logic的下发推送消息, 发送给client
- 保存连接信息, 负载情况到redis, 提供web模块http查询



## 客户端连接相关设计
bucket-->
connChannel-->
connect

connect:
客户端连接的抽象结构, 持有net.Conn和msgChan, 监听msgChan并write消息到client的goroutine.
主要用于推送消息到通道消息的转换, 推到msgChan, 监听goroutine收到推送消息write给client

connChannel:
connect抽象的业务逻辑处理, 负责token添加和校验, connect的add/delete/close, 
启动connect goroutine和提供写消息到connect通道的接口

bucket:
connChannel的桶集合. 负责connChannel的缓存, 使用hash分片, 提高查询效率
