package userserv

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2/log"
	"latipe-transaction-service/config"
	"latipe-transaction-service/internal/adapter/userserv/dto"
)

type UserService struct {
	restyClient *resty.Client
	cfg         *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	restyClient := resty.New().SetDebug(true)

	return &UserService{
		restyClient: restyClient,
		cfg:         cfg,
	}
}

func (us UserService) Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error) {
	resp, err := us.restyClient.
		SetBaseURL(us.cfg.AdapterService.AuthService.BaseURL).
		R().
		SetBody(req).
		SetContext(ctx).
		SetDebug(false).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Authorize token]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Authorize token]: %s", resp.Body())
		return nil, err
	}

	var regResp *dto.AuthorizationResponse

	if err := json.Unmarshal(resp.Body(), &regResp); err != nil {
		log.Errorf("[Authorize token]: %s", err)
		return nil, err
	}

	return regResp, nil
}
