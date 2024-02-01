package dto

import "latipe-transaction-service/pkgs/util/pagable"

type GetAllTransactionLogRequest struct {
	AuthorizationHeader
	Status    int    `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Query     *pagable.Query
}

type GetAllTransactionLogResponse struct {
	pagable.ListResponse
}
