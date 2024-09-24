package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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
		cc := &conn{cc: subConn}
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
		},
	}, nil
}

type conn struct {
	weight        int
	currentWeight int
	cc            balancer.SubConn
}
