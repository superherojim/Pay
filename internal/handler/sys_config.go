package handler

import (
	v1 "bk/api/v1"
	"bk/internal/model"
	"bk/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SysConfigHandler struct {
	*Handler
	sysConfigService service.SysConfigService
}

func NewSysConfigHandler(
	handler *Handler,
	sysConfigService service.SysConfigService,
) *SysConfigHandler {
	return &SysConfigHandler{
		Handler:          handler,
		sysConfigService: sysConfigService,
	}
}

func (h *SysConfigHandler) GetSysConfig(ctx *gin.Context) {
	sysConfig, err := h.sysConfigService.GetSysConfig(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, sysConfig)
}

func (h *SysConfigHandler) CreateSysConfig(ctx *gin.Context) {
	sysConfig := &model.SysConfig{}
	if err := ctx.ShouldBindJSON(sysConfig); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	err := h.sysConfigService.CreateSysConfig(ctx, sysConfig)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, sysConfig)
}

func (h *SysConfigHandler) UpdateSysConfig(ctx *gin.Context) {
	sysConfig := &model.SysConfig{}
	if err := ctx.ShouldBindJSON(sysConfig); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	err := h.sysConfigService.UpdateSysConfig(ctx, sysConfig)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, sysConfig)
}
