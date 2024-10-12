package api

type GetRewardReq struct {
	Meta `path:"/reward/detail"`
	Rid  int64
}
