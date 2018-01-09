package store

import (
	"errors"
	"strconv"

	"github.com/go-redis/redis"
)

const Key = "links"

var errKeyDoesNotExists = errors.New("key does not exists")

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

func (rs *RedisStore) QueryLinkByID(id int) (string, error) {
	if !rs.hashExists() {
		return "", errKeyDoesNotExists
	}

	return rs.client.HGet(Key, strconv.Itoa(id)).Result()
}

func (rs *RedisStore) StoreLink(link *Link) (bool, error) {
	return rs.client.HSet(Key, strconv.Itoa(link.ID), link.URL).Result()
}

func (rs *RedisStore) hashExists() bool {
	exists, err := rs.client.Exists(Key).Result()

	if err != nil || exists == 0 {
		return false
	}

	return true
}
