package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/hfdend/cxz/cli"
)

const (
	// TokenAccessKey 用户token缓存key
	TokenAccessKey = "token_access_"
	// TokenRefreshKey 刷新token缓存key
	TokenRefreshKey = "token_refresh_"
	// TokenUserIDKey TokenUserIDKey
	TokenUserIDKey = "token_user_id_"
)

// Token Token
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

// TokenDefault TokenDefault
var TokenDefault Token

func NewToken(expiry time.Time) *Token {
	t := new(Token)
	t.AccessToken = uuid.New().String()
	t.RefreshToken = uuid.New().String()
	t.Expiry = expiry
	return t
}

func (t Token) EncodeJSON() []byte {
	b, _ := json.Marshal(t)
	return b
}

func (t *Token) SaveUser(userId int) error {
	accessKey := fmt.Sprintf("%s%s", TokenAccessKey, t.AccessToken)
	refreshKey := fmt.Sprintf("%s%s", TokenRefreshKey, t.RefreshToken)
	userIdKey := fmt.Sprintf("%s%d", TokenUserIDKey, userId)
	if err := cli.Redis.Set(accessKey, userId, t.Expiry.Sub(time.Now())).Err(); err != nil {
		return err
	}
	if err := cli.Redis.Set(refreshKey, t.AccessToken, t.Expiry.Add(7*24*time.Hour).Sub(time.Now())).Err(); err != nil {
		return err
	}
	if err := cli.Redis.Set(userIdKey, t.EncodeJSON(), t.Expiry.Add(7*24*time.Hour).Sub(time.Now())).Err(); err != nil {
		return err
	}
	return nil
}

func (Token) Clean(userId int) error {
	userIdKey := fmt.Sprintf("%s%d", TokenUserIDKey, userId)
	data, err := cli.Redis.Get(userIdKey).Bytes()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}
	var t Token
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	accessKey := fmt.Sprintf("%s%s", TokenAccessKey, t.AccessToken)
	refreshKey := fmt.Sprintf("%s%s", TokenRefreshKey, t.RefreshToken)
	if err := cli.Redis.Del(userIdKey).Err(); err != nil {
		return err
	}
	if err := cli.Redis.Del(accessKey).Err(); err != nil {
		return err
	}
	if err := cli.Redis.Del(refreshKey).Err(); err != nil {
		return err
	}
	return nil
}

func (Token) GetUserId(accessToken string) (int, error) {
	key := fmt.Sprintf("%s%s", TokenAccessKey, accessToken)
	id, err := cli.Redis.Get(key).Int64()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return int(id), nil
}
