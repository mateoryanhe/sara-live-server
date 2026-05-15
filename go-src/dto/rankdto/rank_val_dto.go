package rankdto

import "time"

type RankVal struct {
	Id         string    `json:"id"`
	UpdateTime time.Time `json:"updateTime"`
	Val        string    `json:"val"`
}

func NewRankVal() *RankVal {
	return &RankVal{}
}
