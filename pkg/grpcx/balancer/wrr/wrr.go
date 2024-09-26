package wrr

import (
	"context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sync"
)

const name = "custom_wrr"

func init() {
	// NewBalancerBuilder 帮我们 PickerBuilder 转化为 BalanceBuilder
	balancer.Register(base.NewBalancerBuilder("custom_wrr", &PickerBuilder{}, base.Config{HealthCheck: true}))
}

type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, 0)
	for subConn, info := range info.ReadySCs {
		cc := &conn{cc: subConn, available: true}
		md, ok := info.Address.Metadata.(map[string]any)
		if ok {
			weightVal := md["weight"]
			weight, _ := weightVal.(float64)
			cc.weight = int(weight)
		}
		cc.currentWeight = cc.weight
		conns = append(conns, cc)
	}
	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*conn
	mutex sync.Mutex
}

// Pick 基于权重的负载均衡算法
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var (
		totalWeight int
		maxCC       *conn
	)

	for _, conn := range p.conns {
		if conn.available == false {
			continue
		}
		totalWeight += conn.weight
		conn.currentWeight += conn.weight
		if maxCC == nil || maxCC.currentWeight < conn.currentWeight {
			maxCC = conn
		}
	}

	maxCC.currentWeight -= totalWeight

	return balancer.PickResult{
		SubConn: maxCC.cc,
		Done: func(info balancer.DoneInfo) {
			// 很多动态算法 根据结果来 调整权重
			err := info.Err
			if err == nil {

			}
			switch err {
			case context.Canceled:
				return
			case context.DeadlineExceeded:
				return
			case io.EOF, io.ErrUnexpectedEOF:
				// 节点已经崩了
				maxCC.available = false
				return
			default:
				st, ok := status.FromError(err)
				if ok {
					code := st.Code()
					switch code {
					case codes.Unavailable:
						// 这里可能表达的是 熔断
						// 挪走该节点， 该节点已经不可用
						maxCC.available = false
						go func() {
							// 开一个额外的 goroutine 去探活
							// 借助 health check
							// for loop
							if p.healthCheck(maxCC) {
								maxCC.available = true
								// 最好加点流量控制的措施
								// maxCC.currentWeight
								// 掷骰子
							}
						}()
					case codes.ResourceExhausted:
						// 这里可能表达的是 限流
						// 可以挪走 可以留着，留着把两个权重一起调低

						// 加一个错误码 表示降级
					}
				}
			}
		},
	}, nil
}

func (p *Picker) healthCheck(cc *conn) bool {
	// 调用 GRPC 内置的 healthCheck 接口
	return true
}

type conn struct {
	weight        int
	currentWeight int
	cc            balancer.SubConn
	available     bool
	// vip 节点 非 vip 节点 VIP节点全崩了 考虑挤占非VIP 节点
	group string
}
