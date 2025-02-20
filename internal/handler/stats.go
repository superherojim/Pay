package handler

import (
	"net/http"

	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

// @Summary 获取商户总数
// @Router /admin/stats/merchant-count [get]
func (h *StatsHandler) GetMerchantCount(c *gin.Context) {
	count, err := h.statsService.GetMerchantCount(c)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(c, count)
}

// @Summary 获取总订单数
// @Router /admin/stats/total-orders [get]
func (h *StatsHandler) GetTotalOrders(c *gin.Context) {
	count, err := h.statsService.GetTotalOrders(c)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(c, count)
}

// @Summary 获取近期成功订单
// @Param days query int false "天数（默认7天）" default(7)
// @Router /admin/stats/recent-success [get]
func (h *StatsHandler) GetRecentSuccessOrders(c *gin.Context) {
	days := c.DefaultQuery("days", "7")
	count, err := h.statsService.GetRecentSuccessOrders(c, days)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(c, count)
}

// @Summary 仪表盘统计
// @Router /admin/dashboard/stats [get]
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats(c)
	if err != nil {
		v1.HandleError(c, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(c, stats)
}
