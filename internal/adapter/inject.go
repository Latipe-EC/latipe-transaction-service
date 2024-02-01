package adapter

import (
	"github.com/google/wire"
	"latipe-transaction-service/internal/adapter/userserv"
)

var Set = wire.NewSet(userserv.NewUserService)
