package rank

import (
	"context"
	"sort"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/core/cache"
	"xr-game-server/core/lambda"
	"xr-game-server/dao/rankdao"
	"xr-game-server/entity"
	"xr-game-server/gameevent"
)

type CommonRank struct {
	SortData []*ObjModel
	//榜单最大数
	Len int
	//榜单版本
	Version int64
	DB      *cache.CacheMgr
	TypeId  uint32
	Asc     bool
	//结算时间
	SettlementTime *time.Time
}

type ByDescVal []*ObjModel

func (a ByDescVal) Len() int {
	return len(a)
}
func (a ByDescVal) Less(i, j int) bool {
	//升序排序,如果值相等，则按照时间排序,时间小的,排在前面
	return a[i].Val > a[j].Val || (a[i].Val == a[j].Val && a[i].UpdateTime.Before(a[j].UpdateTime))
}
func (a ByDescVal) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ByAscVal []*ObjModel

func (a ByAscVal) Len() int {
	return len(a)
}
func (a ByAscVal) Less(i, j int) bool {
	//降序排序,如果值相等，则按照时间排序,时间小的,排在前面
	return a[i].Val < a[j].Val || (a[i].Val == a[j].Val && a[i].UpdateTime.After(a[j].UpdateTime))
}
func (a ByAscVal) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func NewCommonRank(len int, typeId uint32, asc bool) *CommonRank {
	return &CommonRank{
		Len:    len,
		TypeId: typeId,
		Asc:    asc,
	}
}

// Init 初始化
func (receiver *CommonRank) Init() {
	receiver.SortData = make([]*ObjModel, common.Zero)
	receiver.DB = cache.NewCacheMgr()
	receiver.Version = time.Now().UnixMilli()
}

func (receiver *CommonRank) getDb(id uint64) *entity.PlayerRank {
	rankDb := receiver.DB.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		one := rankdao.GetRankDataBy(receiver.TypeId, id)
		if one == nil {
			one = entity.NewPlayerRank(id, receiver.TypeId)
		}
		return one, nil
	})
	//数据库加
	return rankDb.(*entity.PlayerRank)
}

func (receiver *CommonRank) Add(data *gameevent.AddRankValEventData) {
	//更新数据库
	cacheDb := receiver.getDb(data.Id)
	if cacheDb.LockTime != nil && cacheDb.LockTime.After(time.Now()) && cacheDb.Lock {
		return
	}
	ret := lambda.Filter(receiver.SortData, func(obj *ObjModel) bool {
		return obj.Id == data.Id
	})
	if len(ret) > common.Zero {
		ret[common.Zero].Add(data.Val)
	} else {
		obj := NewObjModel(data.Id, data.Val, time.Now())
		receiver.SortData = append(receiver.SortData, obj)
	}
	//加事件
	cacheDb.AddVal(data.Val)
	cacheDb.SetUpdatedAt(time.Now())
}

func (receiver *CommonRank) Up(id uint64) {
	rankDb := receiver.getDb(id)
	rankDb.SetLock(false)
	rankDb.SetLockTime(nil)
	rankDb.SetUpdatedAt(time.Now())
	receiver.Sort()
}
func (receiver *CommonRank) Down(id uint64, lockTime *time.Time) {
	rankDb := receiver.getDb(id)
	rankDb.SetLock(true)
	rankDb.SetLockTime(lockTime)
	rankDb.SetUpdatedAt(time.Now())
	rankDb.SetVal(common.Zero)
	receiver.SortData = lambda.Filter(receiver.SortData, func(obj *ObjModel) bool {
		return obj.Id != id
	})
	receiver.Sort()
}

func (receiver *CommonRank) Reduce(reduceData *gameevent.ReduceRankEventData) {
	//更新数据库
	cacheDb := receiver.getDb(reduceData.Id)
	if cacheDb.LockTime != nil && cacheDb.LockTime.After(time.Now()) && cacheDb.Lock {
		return
	}
	ret := lambda.Filter(receiver.SortData, func(obj *ObjModel) bool {
		return obj.Id == reduceData.Id
	})
	if len(ret) > common.Zero {
		ret[common.Zero].Reduce(reduceData.Val)
	}
	//减事件
	cacheDb.ReduceVal(reduceData.Val)
	cacheDb.SetUpdatedAt(time.Now())
	receiver.Sort()
}

func (receiver *CommonRank) Sort() {

	if receiver.Asc {
		sort.Sort(ByAscVal(receiver.SortData))
	} else {
		sort.Sort(ByDescVal(receiver.SortData))
	}
	receiver.Version = time.Now().UnixMilli()
	//开始检查榜单是否超标
	if len(receiver.SortData) > receiver.Len {
		receiver.SortData = receiver.SortData[:receiver.Len]
	}
}

// GetValById 获取某个id的分数
func (receiver *CommonRank) GetValById(id uint64) *ObjModel {
	//查询榜单是否存在
	ret := lambda.Filter(receiver.SortData, func(obj *ObjModel) bool {
		return obj.Id == id
	})
	if len(ret) > common.Zero {
		return ret[common.Zero]
	}
	//查询数据库
	cacheData := receiver.DB.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		one := rankdao.GetRankDataBy(receiver.TypeId, id)
		if one == nil {
			one = entity.NewPlayerRank(id, receiver.TypeId)
		}
		return one, nil
	})
	cacheDb, _ := cacheData.(*entity.PlayerRank)
	return NewObjModel(id, cacheDb.Val, cacheDb.UpdatedAt)
}
