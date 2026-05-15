package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
)

var node *snowflake.Node

// GetId 获取雪花算法id
func GetId() uint64 {
	return uint64(node.Generate())
}

func InitSnowflake() {
	if cfg.GetServerCfg().Id == common.Zero {
		//服务器节点默认1
		node, _ = snowflake.NewNode(1)
	} else {
		node, _ = snowflake.NewNode(cfg.GetServerCfg().Id)
	}
}
