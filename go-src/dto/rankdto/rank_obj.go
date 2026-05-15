package rankdto

type RankDto struct {
	ValModel *RankVal  `json:"valModel"`
	RankInfo *RankInfo `json:"rankInfo"`
	Rank     int       `json:"rank" dc:"排名"`
}

func NewRankDto(valModel *RankVal, rankInfo *RankInfo, rank int) *RankDto {
	return &RankDto{
		ValModel: valModel,
		RankInfo: rankInfo,
		Rank:     rank,
	}
}
