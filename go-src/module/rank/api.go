package rank

import (
	"context"
	"xr-game-server/constants/common"
	"xr-game-server/dto/rankdto"
)

func GetRank(context context.Context, req *rankdto.GetRankReq) (res *rankdto.GetRankRes, err error) {
	//roleId := httpserver.GetAuthId(context)
	//roleData := roledao.GetRoleByRole(gconv.Uint64(roleId))
	//commonRank := GetRankBy(req.TypeId, roleData.ServerId)
	//res = rankdto.NewGetRankRes(req.TypeId)
	//if commonRank == nil {
	//	return res, nil
	//}
	//for i, val := range commonRank.SortData {
	//	rankVal := rankdto.NewRankVal()
	//	gconv.Struct(val, rankVal)
	//	targetRole := roledao.GetRoleByRole(val.Id)
	//	res.Data = append(res.Data, rankdto.NewRankDto(rankVal, rankdto.NewRankInfo(targetRole), i+1))
	//}
	//myRank, ok := lambda.Find(res.Data, func(item *rankdto.RankDto) bool {
	//	return item.RankInfo.ID == roleId
	//})
	//if ok {
	//	res.My = myRank
	//} else {
	//	myData := commonRank.GetValById(roleId)
	//	rankVal := rankdto.NewRankVal()
	//	gconv.Struct(myData, rankVal)
	//	res.My = rankdto.NewRankDto(rankVal, rankdto.NewRankInfo(roleData), common.Zero)
	//}
	////排行榜倒计时,可以根据类型扩展一下,建议使用corn表达式，计算出下次时间
	//rankCfg := gamecfg.GetRankCfgById(req.TypeId)
	//if rankCfg != nil {
	//	switch rankCfg.SettlementType {
	//	case gamecfg.DaySettlement:
	//		{
	//			nextDay, _ := xrcorn.GetNextTime(xrcorn.DayPattern)
	//			res.SettlementTime = gconv.String(nextDay.Sub(time.Now()).Milliseconds())
	//		}
	//		break
	//	case gamecfg.WeekSettlement:
	//		{
	//			nextMonday, _ := xrcorn.GetNextTime(xrcorn.WeekPattern)
	//			res.SettlementTime = gconv.String(nextMonday.Sub(time.Now()).Milliseconds())
	//		}
	//		break
	//	case gamecfg.MonthSettlement:
	//		{
	//			nextMonth, _ := xrcorn.GetNextTime(xrcorn.MonthPattern)
	//			res.SettlementTime = gconv.String(nextMonth.Sub(time.Now()).Milliseconds())
	//		}
	//		break
	//	}
	//}
	return res, nil
}

func GetRankVersionReq(ctx context.Context, req *rankdto.GetRankVersionReq) (int64, error) {
	//roleId := httpserver.GetAuthId(ctx)
	//roleData := roledao.GetRoleByRole(gconv.Uint64(roleId))
	//commonRank := GetRankBy(req.TypeId, roleData.ServerId)
	//if commonRank != nil {
	//	return commonRank.Version, nil
	//} else {
	//	return common.Zero, nil
	//}
	return common.Zero, nil
}

func UpRankReq(ctx context.Context, req *rankdto.UpRankReq) (bool, error) {
	//rankCfgs := gamecfg.GetAllRank()
	//for _, rankCfg := range rankCfgs {
	//	rankActor.Send(NewUpEventData(req.RoleId, rankCfg.Id), func(msg any) {
	//		val := msg.(*UpEventData)
	//		mgr := mgrMap[val.TypeId]
	//		roleData := roledao.GetRoleByRole(val.Id)
	//		commonRank := mgr.GetRankByServerId(roleData.ServerId)
	//		commonRank.Up(val.Id)
	//
	//	})
	//}
	return true, nil
}

func DownRankReq(ctx context.Context, req *rankdto.DownRankReq) (bool, error) {
	//if req.LockTime == nil {
	//	//封100年
	//	targetTime := time.Now().Add(time.Hour * 100 * 24)
	//	req.LockTime = &targetTime
	//}
	//rankCfgs := gamecfg.GetAllRank()
	//for _, rankCfg := range rankCfgs {
	//	rankActor.Send(NewDownEventData(req.RoleId, req.LockTime, rankCfg.Id), func(msg any) {
	//		val := msg.(*DownEventData)
	//		mgr := mgrMap[val.TypeId]
	//		roleData := roledao.GetRoleByRole(val.Id)
	//		commonRank := mgr.GetRankByServerId(roleData.ServerId)
	//		commonRank.Down(val.Id, val.LockTime)
	//	})
	//}

	return true, nil
}
