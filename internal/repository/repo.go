package repository

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRedis,
	NewDB,
	NewRepository,
	NewUserRepository,
	NewJwtRepo,
)
