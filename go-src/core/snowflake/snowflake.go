package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
)

var node *snowflake.Node

// GetId 获取雪花算法id
func GetId() uint64 {
	return uint64(node.Generate())
}

// InitSnowflake 按 server.id 初始化雪花 node(各环境必须唯一,合法范围 1-1023).
func InitSnowflake() {
	nodeID := cfg.GetServerCfg().Id
	if nodeID == common.Zero {
		nodeID = 1
	}
	n, err := snowflake.NewNode(nodeID)
	if err != nil {
		g.Log().Fatalf(gctx.New(), "init snowflake node failed id=%d err=%v", nodeID, err)
	}
	node = n
	g.Log().Infof(gctx.New(), "snowflake node ready id=%d", nodeID)
}
