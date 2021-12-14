# json

仿照`gin`框架实现，根据编译时指定的`tags`指定对应的`json`解析方式。

## 标准库（encoding/json）

`go`语言原生支持的`json`解析器。

### 编译

```shell
go build -tags=nomsgpack -o web-server .
```

## jsoniter

一个滴滴研发的`json`解析器，[Jsoniter 的 Golang 版本可以比标准库（encoding/json）快 6 倍之多。而且这个性能是在不使用代码生成的前提下获得的。](http://jsoniter.com/index.cn.html)

### 编译

```shell
go build -tags=jsoniter,nomsgpack -o web-server .
```
