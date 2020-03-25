错误处理
Go Micro的错误处理和错误
Go Micro为分布式系统中发生的大多数错误（包括错误）提供了抽象和类型。通过提供一组核心错误以及定义详细错误类型的能力，我们可以一致地了解典型的Go错误字符串之外的情况。

总览
我们定义以下错误类型：

```go
type Error struct {
    Id     string `json:"id"`
    Code   int32  `json:"code"`
    Detail string `json:"detail"`
    Status string `json:"status"`
}
```


在系统中要求您从处理程序返回错误或从客户端接收错误的任何地方，您都应该假定它是微错误或应该产生一个错误。默认情况下，我们返回 `errors.InternalServerError`（内部错误）和`errors.Timeout`（超时）。

##### 用法示例

假设您的处理程序中发生了一些错误。然后，您应该确定返回哪种错误并执行以下操作。

假设提供的某些数据无效

```go
return errors.BadRequest("com.example.srv.service", "invalid field")
```


如果发生内部错误

```go
if err != nil {
	return errors.InternalServerError("com.example.srv.service", "failed to read db: %v", err.Error())
}
```


同样的，如果从客户端收到错误

```go
pbClient := pb.NewGreeterService("go.micro.srv.greeter", service.Client())
rsp, err := pb.Client(context, req)
if err != nil {
	// parse out the error
    e := errors.Parse(err.Error())
    // inspect the value
	if e.Code == 401 {
	// unauthorised...
}
```

}