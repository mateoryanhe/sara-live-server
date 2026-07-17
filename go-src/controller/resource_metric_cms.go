package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/resourcemetricdto"
	"xr-game-server/module/resourcemonitor"
)

const ResourceMetricUrl = "/resourceMetric"

type ResourceMetricController struct{}

func initResourceMetricController() {
	httpserver.RegCMS(ResourceMetricUrl, &ResourceMetricController{})
}

// GetResourceMetricTrend CMS获取系统资源趋势
func (c *ResourceMetricController) GetResourceMetricTrend(ctx context.Context, req *resourcemetricdto.CMSResourceMetricTrendReq) (res *resourcemetricdto.CMSResourceMetricTrendRes, err error) {
	return resourcemonitor.GetCMSResourceMetricTrend(ctx, req)
}
