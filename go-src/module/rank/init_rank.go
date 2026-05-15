package rank

func initRank(event any) {
	//rankCfg := gamecfg.GetAllRank()
	////如果项目需要跨服榜单,需要项目自己处理,可以加一个类型,判断榜单
	//for _, cfg := range rankCfg {
	//	//检查服务器列表
	//	servers := serverdao.GetAllServer()
	//	//检查榜单管理器是否存在
	//	mgr, mgrOk := mgrMap[cfg.Id]
	//	if !mgrOk {
	//		mgr = NewMgr()
	//		mgrMap[cfg.Id] = mgr
	//	}
	//	for _, server := range servers {
	//		//出现合服标志
	//		if server.ZoneId > common.Zero {
	//			//检查榜单是否存在
	//			zoneRank, zoneOk := mgr.ZoneRank[server.ZoneId]
	//			serverRank, serverOk := mgr.ServerRank[server.ID]
	//			if (!serverOk) && (!zoneOk) {
	//				//生成一个区服榜单
	//				commonRank := NewCommonRank(cfg.Len, cfg.Id, cfg.Asc)
	//				commonRank.Init()
	//				dbData := rankdao.GetZoneData(server.ZoneId, cfg.Id, cfg.Asc, cfg.Len)
	//				for _, list := range dbData {
	//					commonRank.SortData = append(commonRank.SortData, NewObjModel(list.RoleId, list.Val, list.UpdatedAt))
	//				}
	//				mgr.ZoneRank[server.ZoneId] = commonRank
	//			} else {
	//				//如果存在服务器列表,需要移动榜单
	//				if serverOk {
	//					//检查区服是否存在
	//					if zoneOk {
	//						//区服榜已经存在,合并榜单
	//						zoneRank.SortData = append(zoneRank.SortData, serverRank.SortData...)
	//						zoneRank.Sort()
	//						//合并db缓存数据
	//						ctx := gctx.New()
	//						keys, _ := serverRank.DB.Cache.Keys(ctx)
	//						for _, key := range keys {
	//							serverDbCacheVal, _ := serverRank.DB.Cache.Get(ctx, key)
	//							zoneRank.DB.FlushCache(key, serverDbCacheVal)
	//						}
	//					} else {
	//						//移动榜单变成区服榜
	//						mgr.ZoneRank[server.ZoneId] = serverRank
	//						//清掉旧榜单
	//						delete(mgr.ServerRank, server.ID)
	//					}
	//				}
	//			}
	//		} else {
	//			//判断榜单是否存在
	//			_, ok := mgr.ServerRank[server.ID]
	//			if !ok {
	//				commonRank := NewCommonRank(cfg.Len, cfg.Id, cfg.Asc)
	//				commonRank.Init()
	//				mgr.ServerRank[server.ID] = commonRank
	//				dbData := rankdao.GetServerRankData(server.ID, cfg.Id, cfg.Asc, cfg.Len)
	//				for _, list := range dbData {
	//					commonRank.SortData = append(commonRank.SortData, NewObjModel(list.RoleId, list.Val, list.UpdatedAt))
	//				}
	//			}
	//		}
	//
	//	}
	//}
}
