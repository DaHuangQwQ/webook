package domain

type Target struct {
	// 因为什么而打赏
	Biz   string
	BizId int64
	// 作为一个可选的东西
	// 也就是你要打赏的东西是什么
	BizName string

	// 打赏的目标用户
	Uid int64
}

type Reward struct {
	Id     int64
	Uid    int64
	Target Target
	// 同样不着急引入货币。
	Amt    int64
	Status RewardStatus
}

// Completed 是否已经完成
// 目前来说，也就是是否处理了支付回调
func (r Reward) Completed() bool {
	return r.Status == RewardStatusFailed || r.Status == RewardStatusPayed
}

type RewardStatus uint8

func (r RewardStatus) AsUint8() uint8 {
	return uint8(r)
}

const (
	RewardStatusUnknown = iota
	RewardStatusInit
	RewardStatusPayed
	RewardStatusFailed
)

type CodeURL struct {
	Rid int64
	URL string
}
