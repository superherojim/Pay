package service

import (
	"context"
	"strconv"
	"time"

	"bk/internal/repository"
	"bk/pkg/enum"
)

type StatsService interface {
	GetMerchantCount(ctx context.Context) (int64, error)
	GetTotalOrders(ctx context.Context) (int64, error)
	GetRecentSuccessOrders(ctx context.Context, days string) (int64, error)
	GetDashboardStats(ctx context.Context) (map[string]interface{}, error)
}

type statsService struct {
	merchantRepo repository.MerchantsRepository
	orderRepo    repository.OrderRepository
}

func NewStatsService(merchantRepo repository.MerchantsRepository, orderRepo repository.OrderRepository) StatsService {
	return &statsService{
		merchantRepo: merchantRepo,
		orderRepo:    orderRepo,
	}
}

func (s *statsService) GetMerchantCount(ctx context.Context) (int64, error) {
	return s.merchantRepo.TotalCount(ctx)
}

func (s *statsService) GetTotalOrders(ctx context.Context) (int64, error) {
	return s.orderRepo.TotalCount(ctx)
}

func (s *statsService) GetRecentSuccessOrders(ctx context.Context, days string) (int64, error) {
	day, _ := strconv.Atoi(days)
	startTime := time.Now().AddDate(0, 0, -day)
	return s.orderRepo.CountByStatusAndTime(ctx, enum.OrderStatusSuccess, startTime)
}

func (s *statsService) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	// 并行获取所有统计
	type result struct {
		key   string
		value interface{}
		err   error
	}

	ch := make(chan result, 3)

	go func() {
		total, err := s.merchantRepo.TotalCount(ctx)
		ch <- result{"merchantTotal", total, err}
	}()

	go func() {
		total, err := s.orderRepo.TotalCount(ctx)
		ch <- result{"orderTotal", total, err}
	}()

	go func() {
		start := time.Now().AddDate(0, 0, -7)
		total, err := s.orderRepo.CountByStatusAndTime(ctx, enum.OrderStatusSuccess, start)
		ch <- result{"success7Days", total, err}
	}()

	stats := make(map[string]interface{})
	for i := 0; i < 3; i++ {
		res := <-ch
		if res.err != nil {
			return nil, res.err
		}
		stats[res.key] = res.value
	}

	return stats, nil
}
