package middleware

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewCors,
	NewJWTAuth,
	NewRecovery,
)
