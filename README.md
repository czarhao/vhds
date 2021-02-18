这是一个 xds 下发 vhds 的简单例子

1. 安装 envoy 和 go 环境
2. go run example/demo_server.go 来启动模拟的下游服务
3. 启动本 demo
4. envoy --config-path example/vhds.yaml --bootstrap-version 3
