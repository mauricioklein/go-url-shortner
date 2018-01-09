package store

import (
	"strconv"
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore_QueryLinkByID(t *testing.T) {
	client, cleanup := createRedisClient()
	defer cleanup()

	redisStore := NewRedisStore(client)

	// test 1: query with no hash
	_, err := redisStore.QueryLinkByID(1)
	assert.Equal(t, errKeyDoesNotExists, err)

	// test 2: hash exists, but id isn't present
	client.HSet(Key, "2", "http://foo.bar/")

	_, err = redisStore.QueryLinkByID(1)
	assert.Error(t, err)

	// test 3: with link stored
	client.HSet(Key, "1", "http://www.google.com/")

	link, err := redisStore.QueryLinkByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "http://www.google.com/", link)
}

func TestRedisStore_StoreLink(t *testing.T) {
	client, cleanup := createRedisClient()
	defer cleanup()

	redisStore := NewRedisStore(client)
	link := Link{ID: 1, URL: "http://www.google.com/"}

	ok, err := redisStore.StoreLink(&link)
	assert.True(t, ok)
	assert.NoError(t, err)

	// check if the link was actualy stored
	storedLink, err := client.HGet(Key, strconv.Itoa(link.ID)).Result()
	assert.NoError(t, err)
	assert.Equal(t, link.URL, storedLink)
}

func createRedisClient() (*redis.Client, func() *redis.StatusCmd) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return client, client.FlushAll
}
