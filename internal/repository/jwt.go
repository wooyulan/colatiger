package repository

import (
	"colatiger/internal/service"
	"colatiger/pkg/helper/hash"
	"colatiger/pkg/log"
	"context"
	"strconv"
	"time"
)

type jwtRepo struct {
	data *Repository
	log  *log.Logger
}

func NewJwtRepo(data *Repository, log *log.Logger) service.JwtRepo {
	return &jwtRepo{
		data: data,
		log:  log,
	}
}

func (r *jwtRepo) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + hash.MD5Byte([]byte(tokenStr))
}

func (r *jwtRepo) JoinBlackList(ctx context.Context, tokenStr string, joinUnix int64, expires time.Duration) error {
	return r.data.rdb.SetNX(ctx, r.getBlackListKey(tokenStr), joinUnix, expires).Err()
}

func (r *jwtRepo) GetBlackJoinUnix(ctx context.Context, tokenStr string) (int64, error) {
	joinUnixStr, err := r.data.rdb.Get(ctx, r.getBlackListKey(tokenStr)).Result()
	if err != nil {
		return 0, err
	}

	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if joinUnixStr == "" || err != nil {
		return 0, err
	}

	return joinUnix, nil
}
