# p2p-demo

P2P Golang 内网穿透实现

NAT 类型: 端口限制锥形 (NAT3)

## 一、示例说明:


Server 中包含服务端代码,它的主要作用是交换两个内网 P2P 节点临时生成的公网地址


Client 中包含客户端代码,当它接收到另一个客户端的临时公网地址后将会尝试循环发送消息。


## 二、注意事项
请注意,server 端暂时只能支持两个 Client ,当两个 NAT 后的客户端建立连接后,就可以将服务器关闭了

两个 Client 就实现了内网穿透通信。

## 三、运行说明

1. 请在公网服务器上执行 `go run server.go` 它将启动 P2P 转发服务,启动在 0.0.0.0:6999
2. 在 client.go 中配置你 server的 地址,在第 99 行 修改为 IP:6999 的形式
3. 准备两台NAT内网主机启动执行 `go run client.go` 
4. 关闭 P2P 转发服务,你将能看到两台客户端内网穿透通信




