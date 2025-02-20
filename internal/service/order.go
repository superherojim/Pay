package service

import (
	"bytes"
	v1 "cheemshappy_pay/api/v1"
	"cheemshappy_pay/internal/model"
	"cheemshappy_pay/internal/repository"
	"cheemshappy_pay/pkg/chain"
	"cheemshappy_pay/pkg/enum"
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
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type OrderService interface {
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	GetOrderList(ctx context.Context, req *v1.OrderListReq) (*v1.Paginator, error)
	CreateOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error)
	CancelOrder(ctx context.Context, req *v1.OrderParam, apikey string) (*model.Order, error)
	SuccessOrder(ctx context.Context, req *v1.OrderParam) (*model.Order, error)
	GetOrderPay(ctx context.Context, no string) (*v1.OrderPayOut, error)
	ListenOrder(ctx context.Context, no string, req *v1.OrderPayTxOut) (*model.Order, error)
	ListenOrderPay(ctx context.Context) error
}

func NewOrderService(
	service *Service,
	orderRepository repository.OrderRepository,
	merchantsApiService MerchantsApiService,
	merchantsService MerchantsService,
	walletService WalletService,
	sysConfigService SysConfigService,
	conf *viper.Viper,
	verifierFactory *chain.VerifierFactory,
) OrderService {
	return &orderService{
		Service:             service,
		orderRepository:     orderRepository,
		merchantsApiService: merchantsApiService,
		merchantsService:    merchantsService,
		walletService:       walletService,
		sysConfigService:    sysConfigService,
		conf:                conf,
		verifierFactory:     verifierFactory,
	}
}

type orderService struct {
	*Service
	orderRepository     repository.OrderRepository
	merchantsApiService MerchantsApiService
	walletService       WalletService
	merchantsService    MerchantsService
	sysConfigService    SysConfigService
	conf                *viper.Viper
	verifierFactory     *chain.VerifierFactory
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

	if !chain.IsSupportedChain(req.Chain) {
		validChains := chain.GetChainList(false)
		return nil, fmt.Errorf("不支持的区块链网络，支持列表：%s", strings.Join(validChains, ", "))
	}
	wallet, err := s.walletService.GetWalletByMID(ctx, mapi.MID)
	if err != nil {
		return nil, err
	}
	to := 15
	if req.TimeOut > 0 {
		to = req.TimeOut
	}
	order := &model.Order{
		MID:       mapi.MID,
		CNo:       req.OrderNo,
		NotifyURL: mapi.CallbackURL,
		ReturnURL: req.ReturnURL,
		Status:    enum.OrderStatusPending,
		Coin:      req.Coin,
		APIKey:    mapi.Apikey,
		Ac:        wallet.Ac,
		Amount:    req.Amount,
		TimeOut:   int32(to),
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
	or, err := s.orderRepository.GetOrderByNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if or == nil {
		return nil, errors.New("order not found")
	}
	if or.APIKey != apikey {
		return nil, errors.New("unknown apikey")
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
	s.callbackMerchant(order, enum.OrderStatusCanceled)
	return order, nil
}

func (s *orderService) SuccessOrder(ctx context.Context, req *v1.OrderParam) (*model.Order, error) {
	order := &model.Order{
		CNo:    req.OrderNo,
		Status: enum.OrderStatusSuccess,
	}
	order, err := s.orderRepository.SuccessOrder(ctx, order, 0)
	if err != nil {
		return nil, err
	}
	go s.callbackMerchant(order, enum.OrderStatusSuccess)
	return order, nil
}

func (s *orderService) GetOrderPay(ctx context.Context, no string) (*v1.OrderPayOut, error) {
	order, err := s.orderRepository.GetOrderPay(ctx, no)
	if err != nil {
		return nil, err
	}
	mer, err := s.merchantsService.GetMerchants(ctx, order.MID)
	if err != nil {
		return nil, err
	}
	return &v1.OrderPayOut{
		MerName:   mer.Nickname,
		OrderNo:   order.CNo,
		Chain:     order.Chain,
		Coin:      order.Coin,
		Amount:    order.Amount,
		Ac:        order.Ac,
		ReturnURL: order.ReturnURL,
	}, nil
}

func (s *orderService) ListenOrder(ctx context.Context, no string, req *v1.OrderPayTxOut) (*model.Order, error) {
	order, err := s.orderRepository.ListenOrder(ctx, no, req)
	if err != nil {
		return nil, err
	}
	go s.callbackMerchant(order, enum.OrderStatusSuccess)
	return order, nil
}

func (s *orderService) callbackMerchant(order *model.Order, status string) {
	go func() {
		err := s.syncCallbackMerchant(order, status)
		order.NotifyStatus = "success"
		if err != nil {
			s.logger.Error("回调失败", zap.Error(err))
			order.NotifyStatus = "failed"
		}
		s.orderRepository.UpdateOrder(context.Background(), order)
	}()
}

func (s *orderService) syncCallbackMerchant(order *model.Order, status string) error {
	// 获取商户API配置
	mapi, err := s.merchantsApiService.GetMerchantsApiByApikey(context.Background(), order.APIKey)
	if err != nil || mapi == nil {
		s.logger.Error("获取商户API配置失败", zap.Error(err))
		return err
	}

	if mapi.CallbackURL == "" {
		s.logger.Error("商户API配置回调URL为空", zap.String("apikey", order.APIKey))
		return errors.New("商户API配置回调URL为空")
	}

	if !strings.Contains(order.NotifyURL, mapi.CallbackURL) {
		s.logger.Error("商户API配置回调URL不匹配", zap.String("apikey", order.APIKey), zap.String("notify_url", order.NotifyURL))
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

func (s *orderService) ListenOrderPay(ctx context.Context) error {
	orders, err := s.orderRepository.GetOrderByStatus(ctx, enum.OrderStatusListening)
	if err != nil {
		return err
	}

	// 使用带缓冲的channel控制并发数
	sem := make(chan struct{}, 10) // 同时处理10个订单
	var wg sync.WaitGroup
	for _, or := range orders {
		wg.Add(1)
		sem <- struct{}{}
		go func(o *model.Order) {
			defer func() {
				<-sem
				wg.Done()
			}()

			// 检查超时
			if time.Since(o.CreatedAt).Minutes() > float64(o.TimeOut) {
				o.Status = enum.OrderStatusTimeout
				_, err = s.orderRepository.UpdateOrder(ctx, o)
				if err != nil {
					s.logger.Error("订单超时状态更新失败",
						zap.String("order_no", o.No),
						zap.Error(err))
				}
				return
			}

			// 检查交易状态
			status, err := s.checkTransactionStatus(ctx, o)
			if err != nil {
				s.logger.Error("交易状态检查失败",
					zap.String("order_no", o.No),
					zap.String("tx_hash", o.TxHash),
					zap.Error(err))
				return
			}

			// 根据状态更新订单
			if status == enum.OrderStatusSuccess {
				_, err = s.SuccessOrder(ctx, &v1.OrderParam{OrderNo: o.CNo})
				if err != nil {
					s.logger.Error("订单状态更新失败",
						zap.String("order_no", o.No),
						zap.Error(err))
				}
			}
		}(or)
	}

	wg.Wait()
	return nil
}

// 支持多链的验证器接口
func (s *orderService) checkTransactionStatus(ctx context.Context, order *model.Order) (string, error) {
	chainInfo, exists := chain.SupportedChains[order.Chain]
	if !exists {
		return "", fmt.Errorf("unsupported chain: %s", order.Chain)
	}

	rpcURL, ok := s.conf.GetStringMapString("rpc_endpoints")[order.Chain]
	if !ok {
		return "", fmt.Errorf("rpc endpoint not configured for chain: %s", order.Chain)
	}

	verifier, err := s.verifierFactory.GetVerifier(chainInfo.Type)
	if err != nil {
		return "", fmt.Errorf("get verifier failed: %v", err)
	}

	confirmations, status, err := verifier.VerifyTransaction(ctx, rpcURL, order.TxHash)
	if err != nil {
		return "", err
	}

	// 检查交易状态
	if status == 0 {
		return enum.OrderStatusFailed, nil
	}

	// 确认数检查
	if confirmations < s.getRequiredConfirmations(order.Chain) {
		return enum.OrderStatusListening, nil
	}

	return enum.OrderStatusSuccess, nil
}

// 获取不同链的确认数要求
func (s *orderService) getRequiredConfirmations(cn string) int64 {
	if info, exists := chain.SupportedChains[cn]; exists {
		return info.Confirmations
	}
	return 6 // 默认值
}
