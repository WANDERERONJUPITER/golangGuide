### Go Config

动态的可拔插的配置库

应用中的多数配置是静态的，并从多个源中用较为复杂的方式整合在一起。GO-config让这些变得更简单，可拔插并且很容易合并。

#### 特征

##### 动态加载 

根据需要从多个源加载配置。Go Config在后台管理观看配置源，并自动合并和更新内存视图。

##### 可插拔源 

从任意数量的源中进行选择以加载和合并配置。后端源被抽象为内部使用的标准格式，并通过编码器进行解码。源可以是环境变量，标志，文件，etcd，k8s configmap等。

##### 可合并的配置

如果您指定多个配置源，而不论其格式如何，它们将被合并并在单个视图中显示。这极大地简化了优先顺序的加载和基于环境的更改。

##### 观察更改（可选）

观看配置以查看对特定值的更改。使用Go Config的观察程序热重新加载您的应用。您不必处理临时的hup重新加载或其他任何操作，只需继续阅读配置并观察是否需要通知即可。

##### 安全恢复 

万一配置加载不佳或由于某些未知原因而被完全删除，您可以在直接访问任何配置值时指定回退值。这样可以确保在出现问题时，您始终可以阅读一些理智的默认值。

#### 快速开始

Source源 -从中加载配置的后端
Encoder编码器 -处理编码/解码源配置
Reader-使用统一格式合并多个编码源
Config-配置管理器，可管理多个源
Usage用法 -go-config的用法示例
常见问题解答 - 一般问题和解答
TODO -TODO任务/功能

##### Sources

A Source是从中加载配置的后端。可以同时使用多个来源。

支持以下来源：

- cli-从解析的CLI标志中读取
- consul-从consul那里读
- env-从环境变量读取
- etcd-从etcd v3读取
- file文件 -从文件读取
- flag标志 -从标志读取
- memory内存 -从内存读取

还有一些社区支持的插件，它们支持以下来源：

- configmap-从k8s configmap中读取
- grpc-从grpc服务器读取
- runtimevar-从Go Cloud Development Kit运行时变量中读取
- url-从URL读取
- vault保管库 -从保管库服务器读取

todo：

- 支持从git url中读取

###### changeset变更集

源将config作为ChangeSet返回。这是多个后端的单个内部抽象。

```go
type ChangeSet struct {
	// Raw encoded config data
	Data      []byte
	// MD5 checksum of the data
	Checksum  string
	// Encoding format e.g json, yaml, toml, xml
	Format    string
	// Source of the config e.g file, consul, etcd
	Source    string
	// Time of loading or update
	Timestamp time.Time
}
```

##### Encoder

一个Encoder手柄源配置编码/解码。后端源可能以许多不同的格式存储配置。编码器使我们能够处理任何格式。如果未指定编码器，则默认为json。

支持以下编码格式：

- json
- yaml
- toml
- xml
- hcl

##### Reader

reader将多个变更集表示为单个合并且可查询的值集。

```go
type Reader interface {
	// Merge multiple changeset into a single format
	Merge(...*source.ChangeSet) (*source.ChangeSet, error)
	// Return return Go assertable values
	Values(*source.ChangeSet) (Values, error)
	// Name of the reader e.g a json reader
	String() string
}
```


读取器利用编码器将变更集解码为`map[string]interface{}`然后将其合并为一个变更集。它查看“格式”字段以确定编码器。然后，变更集被表示为一组，Values具有检索Go类型和回退无法加载值的能力。

```go
// Values is returned by the reader
type Values interface {
	// Return raw data
        Bytes() []byte
	// Retrieve a value
        Get(path ...string) Value
	// Return values as a map
        Map() map[string]interface{}
	// Scan config into a Go type
        Scan(v interface{}) error
}
```


该Value接口允许强制类型转换/类型断言转换为具有后备默认值的类型。

```go
type Value interface {
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	StringSlice(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
	Bytes() []byte
}
```

##### Config

Config 管理所有配置，抽象出源，编码器和读取器。

它管理来自多个后端源的读取，同步和监视，并将它们表示为单个合并的可查询源。

```go
// Config is an interface abstraction for dynamic configuration
type Config interface {
        // provide the reader.Values interface
        reader.Values
	// Stop the config loader/watcher
	Close() error
	// Load config sources
	Load(source ...source.Source) error
	// Force a source changeset sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
}
```

#### 示例

- 配置示例
- 新配置
- 读取配置
- 读取值
- 监视目录
- 从多个源中读取
- 设置源编码器
- 添加Reader Encoder

###### 配置示例

配置文件可以是任何格式，只要我们有支持它的编码器即可。

示例json配置：

```json
{
    "hosts": {
        "database": {
            "address": "10.0.0.1",
            "port": 3306
        },
        "cache": {
            "address": "10.0.0.2",
            "port": 6379
        }
    }
}
```

###### 新配置

创建一个新的配置（或仅使用默认实例）

```go
import "github.com/micro/go-micro/v2/config"
//TODO 提交bug
conf, _ := config.NewConfig()
```

###### 从文件中加载

从文件中加载配置。它使用文件扩展名来确定配置格式。

```go
import "github.com/micro/go-micro/v2/config"
// Load json config file
config.LoadFile("/tmp/config.json")
```


通过指定带有适当文件扩展名的文件来加载yaml，toml或xml文件

```go
import "github.com/micro/go-micro/v2/config"
// Load yaml config file
config.LoadFile("/tmp/config.yaml")
```


如果扩展名不存在，请指定编码器

```go
import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
)
enc := toml.NewEncoder()
// Load toml file with encoder
config.Load(file.NewSource(
        file.WithPath("/tmp/config"),
	source.WithEncoder(enc),
))
```

###### 读取配置

读取整个配置作为 `map`

```go
// retrieve map[string]interface{}
conf := config.Map()

// map[cache:map[address:10.0.0.2 port:6379] database:map[address:10.0.0.1 port:3306]]
fmt.Println(conf["hosts"])
```


将配置扫描到结构体中

```go
type Host struct {
        Address string `json:"address"`
        Port int `json:"port"`
}

type Config struct{
	Hosts map[string]Host `json:"hosts"`
}

var conf Config

config.Scan(&conf)

// 10.0.0.1 3306
fmt.Println(conf.Hosts["database"].Address, conf.Hosts["database"].Port)
```

###### 读取值

将配置中的值扫描到结构中

```go
type Host struct {
	Address string `json:"address"`
	Port int `json:"port"`
}

var host Host

config.Get("hosts", "database").Scan(&host)

// 10.0.0.1 3306
fmt.Println(host.Address, host.Port)
```


以Go类型读取单个值

```go
// Get address. Set default to localhost as fallback
address := config.Get("hosts", "database", "address").String("localhost")

// Get port. Set default to 3000 as fallback
port := config.Get("hosts", "database", "port").Int(3000)
```

###### 监视路径

观察变化的路径。文件更改时，新值将变得可用。

```go
w, err := config.Watch("hosts", "database")
if err != nil {
	// do something
}

// wait for next value
v, err := w.Next()
if err != nil {
	// do something
}

var host Host

v.Scan(&host)
```

###### 多种来源

可以加载和合并多个源。合并优先级的顺序相反。

```go
config.Load(
	// base config from env
	env.NewSource(),
	// override env with flags
	flag.NewSource(),
	// override flags with file
	file.NewSource(
		file.WithPath("/tmp/config.json"),
	),
)
```

###### 设置源编码器

源需要编码器来编码/解码数据并指定变更集格式。

默认编码器是json。要将编码器更改为yaml，xml，toml，请指定为选项。

```
e := yaml.NewEncoder()

s := consul.NewSource(
	source.WithEncoder(e),
)
```

###### 添加阅读器编码器

读取器使用编码器解码来自不同格式源的数据。

默认阅读器支持json，yaml，xml，toml和hcl。它将合并的配置表示为json。

通过将其指定为选项来添加新的编码器。

```
e := yaml.NewEncoder()

r := json.NewReader(
	reader.WithEncoder(e),
)
```

常问问题
这和viper有什么不同？
Viper和go-config正在解决相同的问题。Go-config提供了不同的界面，并且是更大的工具微观生态系统的一部分。

编码器和阅读器有什么区别？
后端源使用编码器对数据进行编码/解码。读取器使用编码器对来自具有不同格式的多个源的数据进行解码，然后将它们合并为单个编码格式。

对于文件源，我们使用文件扩展名来确定配置格式，因此不使用编码器。

对于consul，etcd或类似的键值源，我们可以从包含多个键的前缀加载，这意味着该源需要了解编码，以便它可以返回单个变更集。

对于环境变量和标志，我们还需要一种将值编码为字节并指定格式的方法，以便以后读者可以将其合并。

为什么变更集数据未表示为map [string] interface {}？
在某些情况下，源数据实际上可能不是键值，因此更容易将其表示为字节并将解码延迟给阅读器。