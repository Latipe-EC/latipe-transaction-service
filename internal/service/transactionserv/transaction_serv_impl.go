package transactionserv

import (
	"context"
	"latipe-transaction-service/internal/domain/dto"
	"latipe-transaction-service/internal/domain/repos"
	redisCache "latipe-transaction-service/pkgs/cache/redis"
	"latipe-transaction-service/pkgs/util/mapper"
)

type transactionService struct {
	transactionRepos repos.TransactionRepository
	cacheEngine      *redisCache.CacheEngine
}

func NewTransactionService(transRepos repos.TransactionRepository, cacheEngine *redisCache.CacheEngine) TransactionService {
	return &transactionService{transactionRepos: transRepos, cacheEngine: cacheEngine}
}

func (t transactionService) GetAllTransaction(ctx context.Context, req *dto.GetAllTransactionLogRequest) (*dto.GetAllTransactionLogResponse, error) {
	trans, total, err := t.transactionRepos.FindAll(ctx, req.Query)
	if err != nil {
		return nil, err
	}
	var dataResp []dto.GetTransactionLogByIdResponse
	if err := mapper.BindingStruct(trans, &dataResp); err != nil {
		return nil, err
	}

	resp := dto.GetAllTransactionLogResponse{}
	resp.Size = req.Query.Size
	resp.Page = req.Query.Page
	resp.Total = total
	resp.Items = dataResp
	resp.HasMore = req.Query.GetHasMore(total)

	return &resp, nil
}

func (t transactionService) GetDetailTransaction(ctx context.Context, req *dto.GetTransactionLogByIdRequest) (*dto.GetTransactionLogByIdResponse, error) {
	trans, err := t.transactionRepos.FindByTransactionId(ctx, req.TransactionID)
	if err != nil {
		return nil, err
	}

	var dataResp dto.GetTransactionLogByIdResponse
	if err := mapper.BindingStruct(trans, &dataResp); err != nil {
		return nil, err
	}

	return &dataResp, nil
}

func (t transactionService) CheckTransactionPending(ctx context.Context) (*dto.TransactionPendingResponse, error) {
	orderIds, total, err := t.cacheEngine.GetAllKey()
	if err != nil {
		return nil, err
	}

	resp := dto.TransactionPendingResponse{
		Count:    total,
		OrderIds: orderIds,
	}
	return &resp, nil
}
