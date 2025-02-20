package handler

import (
	"bytes"
	"cheemshappy_pay/internal/service"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sort"

	v1 "cheemshappy_pay/api/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	*Handler
	orderService        service.OrderService
	walletService       service.WalletService
	merchantsApiService service.MerchantsApiService
}

func NewOrderHandler(
	handler *Handler,
	orderService service.OrderService,
	walletService service.WalletService,
	merchantsApiService service.MerchantsApiService,
) *OrderHandler {
	return &OrderHandler{
		Handler:             handler,
		orderService:        orderService,
		walletService:       walletService,
		merchantsApiService: merchantsApiService,
	}
}

func (h *OrderHandler) GetOrderList(ctx *gin.Context) {
	req := &v1.OrderListReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	orders, err := h.orderService.GetOrderList(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, orders)
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	req := new(v1.OrderParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	apikey := ctx.GetHeader("cheemshappy_pays_token")
	if apikey == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	order, err := h.orderService.CreateOrder(ctx, req, apikey)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, order)
}

func (h *OrderHandler) CancelOrder(ctx *gin.Context) {
	req := new(v1.OrderParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	apikey := ctx.GetHeader("cheemshappy_pays_token")
	if apikey == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	o, err := h.orderService.CancelOrder(ctx, req, apikey)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, o)
}

func (h *OrderHandler) SuccessOrder(ctx *gin.Context) {
	req := new(v1.OrderParam)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	o, err := h.orderService.SuccessOrder(ctx, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, o)
}

func (h *OrderHandler) GetOrderPay(ctx *gin.Context) {
	no := ctx.Param("no")
	order, err := h.orderService.GetOrderPay(ctx, no)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		v1.HandleError(ctx, http.StatusNotFound, fmt.Errorf("err"), nil)
		return
	}
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, order)
}

func (h *OrderHandler) GetOrderPayTx(ctx *gin.Context) {
	no := ctx.Param("no")
	req := new(v1.OrderPayTxOut)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	_, err := h.orderService.ListenOrder(ctx, no, req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, "ok")
}

func (h *OrderHandler) TestCall(ctx *gin.Context) {
	type TestCallbackReq struct {
		MerchantOrderNo string `json:"merchant_order_no" binding:"required"`
		OrderNo         string `json:"order_no" binding:"required"`
		Status          string `json:"status" binding:"required,oneof=success failed"`
		Sign            string `json:"sign" binding:"required"`
		Timestamp       int64  `json:"timestamp" binding:"required"`
	}

	req := new(TestCallbackReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	sign := generateSignature(map[string]interface{}{
		"merchant_order_no": req.MerchantOrderNo,
		"order_no":          req.OrderNo,
		"status":            req.Status,
		"timestamp":         req.Timestamp,
	}, "1e059821-1896-43f5-a227-9b1857e95517")
	if sign != req.Sign {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, "ok")
}

// 生成HMAC-SHA256签名
func generateSignature(data map[string]interface{}, secret string) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%s=%v&", k, data[k]))
	}
	buf.Truncate(buf.Len() - 1) // 移除最后一个&

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}
