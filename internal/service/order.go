package service

import (
	v1 "bk/api/v1"
	"bk/internal/model"
	"bk/internal/repository"
	"bk/pkg/enum"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderService interface {
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	GetOrderList(ctx context.Context, req *v1.OrderListReq) (*v1.Paginator, error)
	CreateOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error)
	CancelOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error)
	SuccessOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error)
	GetOrderPay(ctx context.Context, no string) (*v1.OrderPayOut, error)
	ListenOrder(ctx context.Context, no string, req *v1.OrderPayTxOut) (*model.Order, error)
}

func NewOrderService(
	service *Service,
	orderRepository repository.OrderRepository,
	merchantsApiService MerchantsApiService,
	merchantsService MerchantsService,
	walletService WalletService,
	sysConfigService SysConfigService,
) OrderService {
	return &orderService{
		Service:             service,
		orderRepository:     orderRepository,
		merchantsApiService: merchantsApiService,
		merchantsService:    merchantsService,
		walletService:       walletService,
		sysConfigService:    sysConfigService,
	}
}

type orderService struct {
	*Service
	orderRepository     repository.OrderRepository
	merchantsApiService MerchantsApiService
	walletService       WalletService
	merchantsService    MerchantsService
	sysConfigService    SysConfigService
}

func (s *orderService) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	return s.orderRepository.GetOrder(ctx, id)
}

func (s *orderService) GetOrderList(ctx context.Context, req *v1.OrderListReq) (*v1.Paginator, error) {
	return s.orderRepository.GetOrderList(ctx, req)
}

func (s *orderService) CreateOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error) {
	mapi, err := s.merchantsApiService.CheckApikey(ctx, apikey)
	if err != nil {
		return nil, err
	}
	if mapi == nil {
		return nil, errors.New("apikey not found")
	}
	if mapi.MID != req.MID {
		return nil, errors.New("unknown merchant")
	}
	if req.Amount == "" {
		return nil, errors.New("amount is required")
	}
	if req.Coin == "" {
		return nil, errors.New("coin is required")
	}
	if req.OrderNo == "" {
		return nil, errors.New("order no is required")
	}
	if req.Chain == "" {
		return nil, errors.New("chain is required")
	}
	pattern := `^(0|([1-9]\d*))(\.\d{1,18})?$`
	matched, _ := regexp.MatchString(pattern, req.Amount)
	if !matched {
		return nil, errors.New("amount is invalid")
	}

	supportedChains := map[string]bool{
		// 主网
		"1":     true, // Ethereum
		"56":    true, // BSC
		"137":   true, // Polygon
		"43114": true, // Avalanche
		"10":    true, // Optimism
		"42161": true, // Arbitrum

		// 测试网
		"5":        true, // Goerli
		"11155111": true, // Sepolia
		"97":       true, // BSC Testnet
		"80001":    true, // Polygon Mumbai
		"43113":    true, // Avalanche Fuji
		"420":      true, // Optimism Goerli
		"421613":   true, // Arbitrum Goerli
		"17000":    true, // Holesky
	}

	if !supportedChains[req.Chain] {
		validChains := []string{
			"1 (Mainnet)", "56 (BSC)", "137 (Polygon)",
			"43114 (Avalanche)", "10 (Optimism)", "42161 (Arbitrum)",
			"5 (Goerli)", "11155111 (Sepolia)", "97 (BSC Test)",
			"80001 (Mumbai)", "43113 (Fuji)", "420 (Optimism Test)",
			"421613 (Arbitrum Test)", "17000 (Holesky)",
		}
		return nil, fmt.Errorf("不支持的区块链网络，支持列表：%s", strings.Join(validChains, ", "))
	}

	order := &model.Order{
		MID:       mapi.MID,
		CNo:       req.OrderNo,
		NotifyURL: mapi.CallbackURL,
		ReturnURL: req.ReturnURL,
		Status:    enum.OrderStatusPending,
		Coin:      req.Coin,
		Amount:    req.Amount,
		Remark:    req.Remark,
	}
	sysConfig, err := s.sysConfigService.GetSysConfig(ctx)
	if err != nil {
		return nil, err
	}
	order.PayURL = sysConfig.Domain + "cheems/happy/pay?no=" + req.OrderNo
	return s.orderRepository.CreateOrder(ctx, order, mapi.MID)
}

func (s *orderService) CancelOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error) {
	mapi, err := s.merchantsApiService.CheckApikey(ctx, apikey)
	if err != nil {
		return nil, err
	}
	if mapi == nil {
		return nil, errors.New("apikey not found")
	}
	if mapi.MID != req.MID {
		return nil, errors.New("unknown merchant")
	}
	order := &model.Order{
		MID:    mapi.MID,
		CNo:    req.OrderNo,
		Status: enum.OrderStatusCanceled,
	}
	order, err = s.orderRepository.CancelOrder(ctx, order, mapi.MID)
	if err != nil {
		return nil, err
	}
	s.callbackMerchant(order, apikey, enum.OrderStatusCanceled)
	return order, nil
}

func (s *orderService) SuccessOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error) {
	order := &model.Order{
		CNo:    req.OrderNo,
		Status: enum.OrderStatusSuccess,
	}
	order, err := s.orderRepository.SuccessOrder(ctx, order, 0)
	if err != nil {
		return nil, err
	}
	go s.callbackMerchant(order, apikey, enum.OrderStatusSuccess)
	return order, nil
}

func (s *orderService) GetOrderPay(ctx context.Context, no string) (*v1.OrderPayOut, error) {
	order, err := s.orderRepository.GetOrderPay(ctx, no)
	if err != nil {
		return nil, err
	}
	wallet, err := s.walletService.GetWalletByMID(ctx, order.MID)
	if err != nil {
		return nil, err
	}
	mer, err := s.merchantsService.GetMerchants(ctx, wallet.MID)
	if err != nil {
		return nil, err
	}
	return &v1.OrderPayOut{
		MerName:   mer.Nickname,
		OrderNo:   order.CNo,
		Chain:     order.Chain,
		Coin:      order.Coin,
		Amount:    order.Amount,
		Ac:        wallet.Ac,
		ReturnURL: order.ReturnURL,
	}, nil
}

func (s *orderService) ListenOrder(ctx context.Context, no string, req *v1.OrderPayTxOut) (*model.Order, error) {
	order, err := s.orderRepository.ListenOrder(ctx, no, req)
	if err != nil {
		return nil, err
	}
	go s.callbackMerchant(order, "1e059821-1896-43f5-a227-9b1857e95517", enum.OrderStatusSuccess)
	return order, nil
}

func (s *orderService) callbackMerchant(order *model.Order, apikey, status string) {
	go func() {
		err := s.syncCallbackMerchant(order, apikey, status)
		order.NotifyStatus = "success"
		if err != nil {
			s.logger.Error("回调失败", zap.Error(err))
			order.NotifyStatus = "failed"
		}
		s.orderRepository.UpdateOrder(context.Background(), order)
	}()
}

func (s *orderService) syncCallbackMerchant(order *model.Order, apikey, status string) error {
	// 获取商户API配置
	mapi, err := s.merchantsApiService.GetMerchantsApiByApikey(context.Background(), apikey)
	if err != nil || mapi == nil {
		s.logger.Error("获取商户API配置失败", zap.Error(err))
		return err
	}

	if mapi.CallbackURL == "" {
		s.logger.Error("商户API配置回调URL为空", zap.String("apikey", apikey))
		return errors.New("商户API配置回调URL为空")
	}

	if !strings.Contains(order.NotifyURL, mapi.CallbackURL) {
		s.logger.Error("商户API配置回调URL不匹配", zap.String("apikey", apikey), zap.String("notify_url", order.NotifyURL))
		return errors.New("商户API配置回调URL不匹配")
	}

	// 构造回调请求体
	callbackReq := map[string]interface{}{
		"merchant_order_no": order.CNo,
		"order_no":          order.No,
		"status":            status,
		"amount":            order.Amount,
		"coin":              order.Coin,
		"timestamp":         time.Now().Unix(),
	}

	// 生成签名
	sign := generateSignature(callbackReq, mapi.SecretKey)

	// 设置请求头
	headers := map[string]string{
		"Content-Type":           "application/json",
		"cheems-happy-key":       mapi.Apikey,
		"cheems-happy-sign":      sign,
		"cheems-happy-timestamp": fmt.Sprintf("%d", callbackReq["timestamp"]),
		"cheems-happy-nonce":     uuid.New().String(),
	}
	jsonBody, _ := json.Marshal(callbackReq)
	// 发送POST请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", mapi.CallbackURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.logger.Error("创建请求失败", zap.Error(err))
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 重试逻辑（3次重试）
	maxRetries := 3
	for i := 0; i <= maxRetries; i++ {
		resp, err := client.Post(mapi.CallbackURL, "application/json", bytes.NewBuffer(jsonBody))
		if err == nil && resp.StatusCode == http.StatusOK {
			s.logger.Info("回调成功",
				zap.String("url", order.NotifyURL),
				zap.String("order_no", order.No),
				zap.String("status", status))
			return nil
		}

		if i < maxRetries {
			time.Sleep(time.Duration(i+1) * 2 * time.Second) // 指数退避
		}
	}

	s.logger.Error("回调失败",
		zap.String("order_no", order.No),
		zap.String("url", order.NotifyURL),
		zap.String("status", status))
	return errors.New("回调失败" + order.NotifyURL + " " + order.No + " " + status)
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
