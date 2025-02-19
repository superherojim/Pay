package handler

import (
	v1 "bk/api/v1"
	"bk/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	*Handler
	walletService service.WalletService
}

func NewWalletHandler(
	handler *Handler,
	walletService service.WalletService,
) *WalletHandler {
	return &WalletHandler{
		Handler:       handler,
		walletService: walletService,
	}
}
func (h *WalletHandler) GetWallets(ctx *gin.Context) {
	req := new(v1.WalletList)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	p, err := h.walletService.GetWallet(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, p)
}

func (h *WalletHandler) GetWalletByMID(ctx *gin.Context) {
	id := ctx.Param("mid")
	mid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	wallet, err := h.walletService.GetWalletByMID(ctx, mid)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, wallet)
}

func (h *WalletHandler) DeleteWallet(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err = h.walletService.DeleteWallet(ctx, idInt)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *WalletHandler) CreateWallet(ctx *gin.Context) {
	req := new(v1.Wallet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err := h.walletService.CreateWallet(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *WalletHandler) AddWallet(ctx *gin.Context) {
	req := new(v1.AddWallet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err := h.walletService.AddWallet(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

func (h *WalletHandler) UpdateWallet(ctx *gin.Context) {
	req := new(v1.AddWallet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	err := h.walletService.UpdateWallet(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
