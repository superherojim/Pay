package handler

import (
	v1 "bk/api/v1"
	"bk/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MerchantsHandler struct {
	*Handler
	merchantsService service.MerchantsService
}

func NewMerchantsHandler(
	handler *Handler,
	merchantsService service.MerchantsService,
) *MerchantsHandler {
	return &MerchantsHandler{
		Handler:          handler,
		merchantsService: merchantsService,
	}
}
func (h *MerchantsHandler) GetMerchantsList(ctx *gin.Context) {
	req := new(v1.MerchantsListReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	merchants, err := h.merchantsService.GetMerchantsList(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, merchants)
}

func (h *MerchantsHandler) CreateMerchants(ctx *gin.Context) {
	req := new(v1.MerParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	req.CreatedBy = userId
	err := h.merchantsService.CreateMerchants(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *MerchantsHandler) GetMerchants(ctx *gin.Context) {
	id := ctx.Param("id")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	mer, err := h.merchantsService.GetMerchants(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, mer)
}

func (h *MerchantsHandler) DeleteMerchants(ctx *gin.Context) {
	id := ctx.Param("id")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	err = h.merchantsService.DeleteMerchants(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *MerchantsHandler) UpdateMerchants(ctx *gin.Context) {
	req := new(v1.MerParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, err, nil)
		return
	}
	userId := GetUserIdFromCtx(ctx)
	req.CreatedBy = userId
	err := h.merchantsService.UpdateMerchants(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *MerchantsHandler) GetMerchantsIN(ctx *gin.Context) {
	merchants, err := h.merchantsService.GetMerchantsIN(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, merchants)
}
