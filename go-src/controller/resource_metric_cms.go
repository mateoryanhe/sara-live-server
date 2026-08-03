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

func (c *ResourceMetricController) GetResourceMetricMemoryTrend(ctx context.Context, req *resourcemetricdto.CMSResourceMetricMemoryTrendReq) (res *resourcemetricdto.CMSResourceMetricMemoryTrendRes, err error) {
	return resourcemonitor.GetCMSResourceMetricMemoryTrend(ctx, req)
}

func (c *ResourceMetricController) GetResourceMetricHeapTrend(ctx context.Context, req *resourcemetricdto.CMSResourceMetricHeapTrendReq) (res *resourcemetricdto.CMSResourceMetricHeapTrendRes, err error) {
	return resourcemonitor.GetCMSResourceMetricHeapTrend(ctx, req)
}

func (c *ResourceMetricController) GetResourceMetricRatioTrend(ctx context.Context, req *resourcemetricdto.CMSResourceMetricRatioTrendReq) (res *resourcemetricdto.CMSResourceMetricRatioTrendRes, err error) {
	return resourcemonitor.GetCMSResourceMetricRatioTrend(ctx, req)
}

func (c *ResourceMetricController) GetResourceMetricCpuTrend(ctx context.Context, req *resourcemetricdto.CMSResourceMetricCpuTrendReq) (res *resourcemetricdto.CMSResourceMetricCpuTrendRes, err error) {
	return resourcemonitor.GetCMSResourceMetricCpuTrend(ctx, req)
}

func (c *ResourceMetricController) GetResourceMetricOnlineTrend(ctx context.Context, req *resourcemetricdto.CMSResourceMetricOnlineTrendReq) (res *resourcemetricdto.CMSResourceMetricOnlineTrendRes, err error) {
	return resourcemonitor.GetCMSResourceMetricOnlineTrend(ctx, req)
}
