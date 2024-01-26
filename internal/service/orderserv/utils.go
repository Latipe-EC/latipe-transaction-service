package orderserv

import "latipe-transaction-service/internal/domain/entities"

func MappingTxStatus(trans *entities.TransactionLog, serviceType int, status int) *entities.TransactionLog {
	for _, i := range trans.Commits {
		if i.ServiceName == MappingServiceName(serviceType) {
			i.TxStatus = status
		}
	}
	return trans
}

func MappingServiceName(serviceType int) string {
	switch serviceType {
	case PRODUCT_SERVICE:
		return PRODUCT_SERVICE_NAME
	case PROMOTION_SERVICE:
		return PROMOTION_SERVICE_NAME
	case PAYMENT_SERVICE:
		return PAYMENT_SERVICE_NAME
	case DELIVERY_SERVICE:
		return DELIVERY_SERVICE_NAME
	}
	return ""
}

func IsCommitSuccess(commits []entities.Commits) bool {
	for _, i := range commits {
		if i.TxStatus != entities.TX_SUCCESS {
			return false
		}
	}
	return true
}
