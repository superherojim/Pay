package handler

import (
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SysWalletHandler struct {
	*Handler
	sysWalletService service.SysWalletService
	walletService    service.WalletService
}

func NewSysWalletHandler(
	handler *Handler,
	sysWalletService service.SysWalletService,
	walletService service.WalletService,
) *SysWalletHandler {
	return &SysWalletHandler{
		Handler:          handler,
		sysWalletService: sysWalletService,
		walletService:    walletService,
	}
}

func (h *SysWalletHandler) GetSysWallet(ctx *gin.Context) {
	wallet, err := h.sysWalletService.GetSysWallet(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, wallet)
}

func (h *SysWalletHandler) UpdateSysWallet(ctx *gin.Context) {
	req := new(v1.SysWalletUpdateReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.sysWalletService.UpdateSysWallet(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *SysWalletHandler) CreateSysWallet(ctx *gin.Context) {
	if err := h.sysWalletService.CreateSysWallet(ctx); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *SysWalletHandler) DeriveChildWallet(ctx *gin.Context) {
	// 获取商户ID
	req := new(v1.Wallet)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	// 调用服务派生钱包
	childWallet, err := h.sysWalletService.DeriveChildWallet(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	// 关联商户ID
	childWallet.MID = req.MID
	// 保存到wallet表
	if err := h.walletService.CreateChildWallet(ctx, childWallet); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, childWallet)
}
