### HelloWorld

Go Micro的helloworld示例

#### 编写服务

使用Go Micro编写服务非常简单。它提供了快速构建的框架，而无需首先了解所有内容。下面是一个简单的问候服务示例，我们将对其进行处理。

##### 服务协议

微服务的关键要求之一是严格定义接口。Micro使用protobuf来实现这一目标。

在这里，我们使用Hello方法定义了Greeter处理程序。它需要带有两个字符串参数的Request和Response。

```protobuf
syntax = "proto3";

service Greeter {
	rpc Hello(Request) returns (Response) {}
}

message Request {
	string name = 1;
}

message Response {
	string greeting = 2;
}
```

#### 生成原型

编写原型定义后，我们必须使用带有微型插件的protoc对其进行编译。

```shell
protoc --micro_out=. --go_out=. *.proto
```

实施服务
现在，我们已经定义了服务接口，我们需要实现服务。

以下是迎宾服务的代码。它执行以下操作：

实现为Greeter处理程序定义的接口
初始化微服务
注册Greeter处理程序
运行服务

```go
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro/v2"
	proto "golangGuide/goMicro/src/helloworld/proto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("HelloWorld"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
```

运行服务
现在使用Go运行示例。如果使用示例代码，请执行以下操作：

```shell
go run examples/service/main.go
```

写客户
一旦获得服务，我们实际上需要一种查询它的方法。这是微服务的核心，因为我们不仅提供服务，而且还使用其他服务。下面是查询欢迎服务的客户端代码。

生成的原型包括一个问候客户端，以减少样板代码。

```go
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro/v2"
	proto "golangGuide/goMicro/src/helloworld/proto"
)

func main() {
	// Create a new service
	service := micro.NewService(micro.Name("label your service"))
	//service := micro.NewService()

	// Initialise the client and parse command line flags
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("HelloWorld", service.Client())

	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.Request{Name: "aaron"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)
}
```

现在运行客户端

```
go run client.go
```

现在运行客户端

```shell
go run client.go
```

输出应该只是打印响应

```
Hello aaron
```

