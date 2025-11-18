package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) DashboardStatsHandler(ctx context.Context, c *app.RequestContext) {
	type dashboardStatsResponse struct {
		TotalNodes        int `json:"total_nodes"`
		TotalApplications int `json:"total_applications"`
		TotalServices     int `json:"total_services"`
	}

	c.JSON(http.StatusOK, SuccessResponse(dashboardStatsResponse{
		TotalNodes:        10,
		TotalApplications: 10,
		TotalServices:     10,
	}))
}
