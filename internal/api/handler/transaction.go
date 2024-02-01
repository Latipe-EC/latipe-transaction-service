package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"latipe-transaction-service/internal/domain/dto"
	"latipe-transaction-service/internal/service/transactionserv"
	"latipe-transaction-service/pkgs/util/pagable"
	responses "latipe-transaction-service/pkgs/util/response"
)

type TransactionApiHandler struct {
	transactionServ transactionserv.TransactionService
}

func NewTransactionHandler(service transactionserv.TransactionService) *TransactionApiHandler {
	return &TransactionApiHandler{transactionServ: service}
}
func (h TransactionApiHandler) GetAllTransaction(ctx *fiber.Ctx) error {
	var request dto.GetAllTransactionLogRequest

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return responses.ErrInvalidParameters
	}

	request.Query = query
	dataResp, err := h.transactionServ.GetAllTransaction(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (h TransactionApiHandler) GetTransactionById(ctx *fiber.Ctx) error {
	var request dto.GetTransactionLogByIdRequest

	if err := ctx.ParamsParser(&request); err != nil {
		return responses.ErrInvalidParameters
	}

	dataResp, err := h.transactionServ.GetDetailTransaction(ctx.Context(), &request)
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		case mongo.ErrNoDocuments:
			return responses.ErrNotFoundRecord
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (h TransactionApiHandler) GetTransactionInQueue(ctx *fiber.Ctx) error {

	dataResp, err := h.transactionServ.CheckTransactionPending(ctx.Context())
	if err != nil {
		log.Errorf("%v", err)
		switch err {
		default:
			return responses.ErrInternalServer
		}
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}
