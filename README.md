# gpusher

## 模块
- gateway
- logic
- message
- web

## 模块功能
g### ateway(网关)
- 维护与client的连接, 转发消息到logic处理
- 接收logic的下发推送消息, 发送给client
- 保存连接信息, 负载情况到redis, 提供web模块http查询

### logic(逻辑服)
- 接收网络的客户端转发消息, 处理消息, 判断转发关系(转发给另外哪个client)
- 监听web推送到消息队列的消息, 判断推送消息的下发路由关系(下发给哪个网关)
- 保存client与网关的路由关系(rpc调用message的接口) 

### message(数据层)
- 提供rpc接口
- 保存client与网关的路由关系
- 提供获取client与网关的路由关系
- 保存离线消息
- 提供获取离线消息

### web(web服务)
- 提供RESTful接口
- 提供推送消息接口
- 获取系统信息接口(读取redis)


### 服务启动流程
- 1.启动gateway
- 2.启动message服
- 3.启动logic服
- 4.启动web服务


### 客户端连接网关
    client-->gateway


### logic grpc长连接所有网关
          -->gateway1
    logic -->gateway2
          -->gateway3


### 消息推送流程
    web-->kafka-->logic         -->gateway-->client
                       ->message
                        










