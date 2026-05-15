package namedao

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity"
)

func InitName() {

}

func GetNameBy(typeId entity.DbNameType, val string) *entity.Name {
	var db *entity.Name
	g.Model(string(entity.TbName)).Unscoped().Where(g.Map{
		string(entity.NameTypeId): uint8(typeId),
		string(entity.NameVal):    val,
	}).Scan(&db)
	return db
}

func GetValBy(typeId entity.DbNameType) []gdb.Value {
	names, _ := g.Model(string(entity.TbName)).Unscoped().Where(entity.NameTypeId, uint8(typeId)).Fields(string(entity.NameVal)).Array()
	return names
}
