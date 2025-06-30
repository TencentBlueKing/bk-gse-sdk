# Data-Report 数据管道上报

## 快速接入
### 1. 确定数据上报的data-id(channel-id)
`data-id`决定了你的数据上报后的路由，最终转发到哪个数据仓库中（如kafka）。
联系GSE管理员获取你的`data-id`。

### 2. 部署你的进程
通过Agent提供的数据上报能力（Unix下为IPC，Windows下为Local STREAM Socket），建立链接连到本机的Agent上。
SDK提供全套的通信链接管理、上报数据接口等。

```golang
func main() {
    // 初始化client SDK
    // 此处需要填写:
    // - DomainSocketPath: unix(linux + macos)下, 使用的是本地ipc文件连接, 需要指定文件路径
    //      若是Windows下, 使用的是本地监听的TCP端口, 改用agentmessage.WithLocalSocketPort(uint)
    // - Logger: 可指定用户自己的日志实现, 默认是打stdout日志, 也可以使用types.NewEmptyLogger关闭日志
    client, err := agentmessage.New(
		agentreport.WithDomainSocketPath(config.DomainSocketPath),
		agentreport.WithLogger(types.NewDefaultLogger(1)),
    )
    if err != nil {
        panic(err)
    }

    // 开始连接Agent, Launch会阻塞至首次成功连接为止, 可以通过context来控制尝试连接的超时时间
    // 一旦成功Launch之后, client会自动管理与Agent之间的连接, 并对异常情况进行自动重连
	// launch client, it will try to connect to agent and keep the connection.
    if err = client.Launch(ctx); err != nil {
        panic(err)
    }

    // 上报数据到指定的data-id
    // content则为数据内容
    if err = client.ReportData(ctx, dataID, []byte("[this is the data reported from client]")); err != nil {
        panic(err)
    }
}
```

### 3. 验证
在data-id路由落地的数据仓库中验证是否收到上报的数据