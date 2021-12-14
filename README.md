# go-web-example

一个基于 Go 语言的 Web 项目示例

## 压测性能

### 编译

```shell
go build -tags=jsoniter,nomsgpack -o web-server .
```

### 压测机器配置

```text
MacBook Pro (13-inch, 2020, Four Thunderbolt 3 ports)
处理器  2 GHz 四核Intel Core i5
内存   16 GB 3733 MHz LPDDR4X
```

### 使用`wrk`压测

```shell
# prod 角色启动
Running 30s test @ http://localhost:8000/ping
  8 threads and 500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.47ms    2.76ms  51.51ms   88.83%
    Req/Sec    13.00k     2.99k   39.60k    63.86%
  3106226 requests in 30.07s, 432.50MB read
  Socket errors: connect 253, read 96, write 0, timeout 0
Requests/sec: 103304.12
Transfer/sec:     14.38MB

# dev 角色启动
Running 30s test @ http://localhost:8000/ping
  8 threads and 500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.47ms    4.96ms  77.69ms   86.82%
    Req/Sec     7.19k     2.18k   18.09k    67.88%
  1719375 requests in 30.06s, 239.40MB read
  Socket errors: connect 253, read 94, write 0, timeout 0
Requests/sec:  57200.47
Transfer/sec:      7.96MB
```

## 部署

1. 使用`go`

```shell
# jsoniter: 外部引用的json序列化方式
# nomsgpack: 禁用Gin的默认渲染
go build -tags=jsoniter,nomsgpack -o web-server .
```

2. 使用`docker`
   1. 编译名为`web-server`的镜像
      ```shell
      docker build . -t web-server
      ```
   2. 无参数运行：只输出帮助信息
      ```shell
      docker run web-server
      # Output:
      # Usage of ./web-server:
      # ... ...
      ```
   3. 启动服务：通过附加应用参数
      ```shell
      docker run -p 8000:8000 web-server -appMode=dev
      ```
