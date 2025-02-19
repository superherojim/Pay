package handler

import (
	v1 "bk/api/v1"
	"bk/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MerchantsMetaHandler struct {
	*Handler
	merchantsMetaService service.MerchantsMetaService
}

func NewMerchantsMetaHandler(
	handler *Handler,
	merchantsMetaService service.MerchantsMetaService,
) *MerchantsMetaHandler {
	return &MerchantsMetaHandler{
		Handler:              handler,
		merchantsMetaService: merchantsMetaService,
	}
}

func (h *MerchantsMetaHandler) GetMerchantsMeta(ctx *gin.Context) {
	id := ctx.Param("mid")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	mer, err := h.merchantsMetaService.GetMerchantsMeta(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, mer)
}

func (h *MerchantsMetaHandler) UpdateMerchantsMeta(ctx *gin.Context) {
	id := ctx.Param("mid")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	var req v1.MerchantsMeta
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err = h.merchantsMetaService.UpdateMerchantsMeta(ctx, mid, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
