package handler

import (
	v1 "bk/api/v1"
	"bk/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MerchantsApiHandler struct {
	*Handler
	merchantsApiService service.MerchantsApiService
}

func NewMerchantsApiHandler(
	handler *Handler,
	merchantsApiService service.MerchantsApiService,
) *MerchantsApiHandler {
	return &MerchantsApiHandler{
		Handler:             handler,
		merchantsApiService: merchantsApiService,
	}
}

func (h *MerchantsApiHandler) GetMerchantsApiList(ctx *gin.Context) {
	req := new(v1.MerApiListReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	mer, err := h.merchantsApiService.GetMerchantsApiList(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, mer)
}

func (h *MerchantsApiHandler) GetMerchantsApi(ctx *gin.Context) {
	id := ctx.Param("mid")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	mer, err := h.merchantsApiService.GetMerchantsApi(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, mer)
}

func (h *MerchantsApiHandler) CreateMerchantsApi(ctx *gin.Context) {
	req := new(v1.MerApiParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err := h.merchantsApiService.CreateMerchantsApi(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *MerchantsApiHandler) UpdateMerchantsApi(ctx *gin.Context) {
	req := new(v1.MerApiParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err := h.merchantsApiService.UpdateMerchantsApi(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *MerchantsApiHandler) DeleteMerchantsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err = h.merchantsApiService.DeleteMerchantsApi(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *MerchantsApiHandler) TestCallback(ctx *gin.Context) {
	id := ctx.Param("id")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err = h.merchantsApiService.TestCallback(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
