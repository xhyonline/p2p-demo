package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-reuseport"
	"math"
	"math/big"
	"net"
	"time"
)

type Handler struct {
	// 中继服务器的连接句柄
	ServerConn net.Conn
	// p2p 连接
	P2PConn net.Conn
	// 端口复用
	LocalPort int
}

// WaitNotify 等待远程服务器发送通知告知我们另一个用户的公网IP
func (s *Handler) WaitNotify() {
	buffer := make([]byte, 1024)
	n, err := s.ServerConn.Read(buffer)
	if err != nil {
		panic("从服务器获取用户地址失败" + err.Error())
	}
	data := make(map[string]string)
	if err := json.Unmarshal(buffer[:n], &data); err != nil {
		panic("获取用户信息失败" + err.Error())
	}
	fmt.Println("客户端获取到了对方的地址:", data["address"])
	// 断开服务器连接
	defer s.ServerConn.Close()
	// 请求用户的临时公网 IP
	go s.DailP2PAndSayHello(data["address"], data["dst_uid"])
}

// DailP2PAndSayHello 连接对方临时的公网地址,并且不停的发送数据
func (s *Handler) DailP2PAndSayHello(address, uid string) {
	var errCount = 1
	var conn net.Conn
	var err error
	for {
		// 重试三次
		if errCount > 3 {
			break
		}
		time.Sleep(time.Second)
		conn, err = reuseport.Dial("tcp", fmt.Sprintf(":%d", s.LocalPort), address)
		if err != nil {
			fmt.Println("请求第", errCount, "次地址失败,用户地址:", address)
			errCount++
			continue
		}
		break
	}
	if errCount > 3 {
		panic("客户端连接失败")
	}
	s.P2PConn = conn
	go s.P2PRead()
	go s.P2PWrite()
}

// P2PRead 读取 P2P 节点的数据
func (s *Handler) P2PRead() {
	for {
		buffer := make([]byte, 1024)
		n, err := s.P2PConn.Read(buffer)
		if err != nil {
			fmt.Println("读取失败", err.Error())
			time.Sleep(time.Second)
			continue
		}
		body := string(buffer[:n])
		fmt.Println("读取到的内容是:", body)
		fmt.Println("来自地址", s.P2PConn.RemoteAddr())
		fmt.Println("=============")
	}
}

// P2PWrite 向远程 P2P 节点写入数据
func (s *Handler) P2PWrite() {
	for {
		if _, err := s.P2PConn.Write([]byte("你好呀~")); err != nil {
			fmt.Println("客户端写入错误")
		}
		time.Sleep(time.Second)
	}
}

func main() {
	// 指定本地端口
	localPort := RandPort(10000, 50000)
	// 向 P2P 转发服务器注册自己的临时生成的公网 IP (请注意,Dial 这里拨号指定了自己临时生成的本地端口)
	serverConn, err := reuseport.Dial("tcp", fmt.Sprintf(":%d", localPort), "你自己的公网服务器IP:6999")
	if err != nil {
		panic("请求远程服务器失败" + err.Error())
	}
	h := &Handler{ServerConn: serverConn, LocalPort: int(localPort)}
	h.WaitNotify()
	time.Sleep(time.Hour)
}

// RandPort 生成区间范围内的随机端口
func RandPort(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return min + result.Int64()
}
