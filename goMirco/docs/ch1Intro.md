### 微服务开发的框架Go Micro

#### 总览

Go Micro提供了分布式系统开发的核心要求，包括RPC和事件驱动的通信。Microde 哲学是提供一种完整的寄语可插拔体系结构的默认配置。我们提供了默认设置，可帮助您快速入门、使用。

#### 特征

Go Micro提取了分布式系统的详细信息。以下主要功能。

##### 服务发现 

自动服务注册和名称解析。服务发现是微服务开发的核心。当服务A需要与服务B通话时，它需要该服务的位置。默认发现机制是多播DNS（mdns），一种零配置系统。

##### 负载平衡 

基于服务发现的客户端负载平衡。一旦获得了服务的任意数量的实例的地址，我们现在需要一种方法来确定要路由到的节点。我们使用随机散列负载平衡来提供服务之间的平均分配，并在出现问题时重试其他节点。

##### 消息编码

基于内容类型的动态消息编码。客户端和服务器将使用编解码器以及content-type来为您无缝编码和解码Go类型。各种消息可以被编码并从不同的客户端发送。客户端和服务器默认情况下会处理此问题。默认情况下，这包括protobuf和json。

##### 请求/响应 

基于RPC的请求/响应，支持双向流。我们为同步通信提供了一个抽象。对服务的请求将被自动解决，负载均衡，拨号和流式传输。默认的传输使用是gRPC。

##### 异步消息传递 

PubSub内置为异步通信和事件驱动的体系结构的一等公民。事件通知是微服务开发中的核心模式。默认消息系统是嵌入式NATS 服务器。

##### 可插拔接口 

Go Micro对每个分布式系统抽象都使用Go接口。因此，这些接口是可插入的，并允许Go Micro与运行时无关。您可以插入任何底层技术。在github.com/micro/go-plugins中找到插件 。

#### 入门

##### 依赖

Go Micro默认使用protobuf。这样一来，我们就可以编写代码以生成样板代码，并提供一种有效的线路格式来在服务之间来回传输数据。

我们还需要某种形式的服务发现，以将服务名称以及元数据和端点信息解析为它们的地址。请参阅下面的更多信息。

**Protobuf**
安装protobuf才能代码生成API接口

protoc-gen-micro

**服务发现**

服务发现用于将服务名称解析为地址。默认情况下，我们提供使用多播DNS的zeroconf发现系统。这是大多数操作系统内置的。如果您需要更具弹性和多主机的功能，请使用etcd。

**Etcd**
Etcd可用作替代服务发现系统。

- 下载并运行etcd

- 传递`--registry=etcd` 给任何命令或环境变量 `MICRO_REGISTRY=etcd`

```
MICRO_REGISTRY=etcd go run main.go
```

服务发现是可拔插的结构。在micro / go-plugins存储库中找到consul，kubernetes，zookeeper等插件。

#### 安装

Go Micro是用于基于Go的开发的框架。可以使用go工具链轻松获得此功能。

在服务中导入go-micro

```
import "github.com/micro/go-micro/v2"
```


我们提供发行标签，建议您坚持使用最新的稳定发行版。利用go modules将启用此功能。

```shell
#enable go modules
export GO111MODULE=on
#initialise go modules in your app
go mod init
#now go get
go get ./...
```

#### 编写服务

去看看hello world示例开始