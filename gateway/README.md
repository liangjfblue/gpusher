
gateway(网关)
- 维护与client的连接, 转发消息到logic处理
- 接收logic的下发推送消息, 发送给client
- 保存连接信息, 负载情况到redis, 提供web模块http查询