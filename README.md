# gpusher

## 模块
- gateway
- logic
- message
- web

## 模块功能
### gateway(网关)
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
- 管理员操作(推送消息, 删除消息, 查看分析等)


### 服务启动流程
- 1.启动message服
- 2.启动logic服
- 3.启动web服务
- 4.启动gateway


### gpusher使用
#### 本地启动
- 1.启动 message, logic, web, gateway
- 2.启动 test/client, 测试连接 gpusher
- 3.推送消息 http://ip:port/v1/push

#### 容器化启动
##### 镜像打包
./script/build-docker.sh

##### 运行各个服务的容器
docker run --rm -t -p 7780:7780 --net=host gpusher-message
docker run --rm -t -p 7771:7771 -p 8881:8881 --net=host gpusher-gateway
docker run --rm -t -p 7772:7772 --net=host gpusher-logic
docker run --rm -t -p 7030:7030 --net=host gpusher-web


### 客户端连接网关
    client-->gateway


### logic grpc长连接所有网关
          -->gateway1
    logic -->gateway2
          -->gateway3


### 消息推送流程
    web-->kafka-->logic         -->gateway-->client
                       ->message
     
     
### watch etcd service list
                 
    etcdctl --endpoints "http://172.16.7.16:9002,http://172.16.7.16:9004,http://172.16.7.16:9006" watch /etcd/gpusher --prefix



### 消息格式
```http://ip:port/v1/push```

    {
        "tag":"app_gpusher",
        "body": {
            "type":1,
            "uuid":"liangjf",
            "content":"hello world 123",
            "expireTime":3600,
            "offlinePush":false
        }
    }

- tag: appName, 也作为topic
- body: 推送体
    - type: 推送类型(1-个体推送, 2-app推送, 3-全体推送)
    - uuid: 推送消息接收者
    - content: 推送消息内容
    - expireTime: 推送消息过期时间
    - offlinePush: 是否离线推送(用户登陆时会拉取离线消息)



