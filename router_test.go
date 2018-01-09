package urlshortner

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"

	"github.com/mauricioklein/go-url-shortner/store"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Homepage(t *testing.T) {
	linkStore, closeConn, err := createLinkStore()
	if err != nil {
		t.Error(err)
	}
	defer closeConn()

	redisStore, cleanRedis := createRedisStore()
	defer cleanRedis()

	router := NewRouter(linkStore, redisStore)
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)

	assert.Contains(t, w.Body.String(), "Welcome to the URL shortner")
}

func TestRouter_RegisterURL_WithValidURL(t *testing.T) {
	linkStore, closeConn, err := createLinkStore()
	if err != nil {
		t.Error(err)
	}
	defer closeConn()

	redisStore, cleanRedis := createRedisStore()
	defer cleanRedis()

	router := NewRouter(linkStore, redisStore)
	r := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: "/"},
		Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencode"}},
		Form:   url.Values{"url": []string{"http://www.google.com/"}},
	}
	w := httptest.NewRecorder()

	r.ParseForm()
	router.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Contains(t, w.Body.String(), "Your url is")
}

func TestRouter_RegisterURL_WithInvalidURL(t *testing.T) {
	linkStore, closeConn, err := createLinkStore()
	if err != nil {
		t.Error(err)
	}
	defer closeConn()

	redisStore, cleanRedis := createRedisStore()
	defer cleanRedis()

	router := NewRouter(linkStore, redisStore)
	r := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Path: "/"},
		Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencode"}},
		Form:   url.Values{"url": []string{"foobar"}},
	}
	w := httptest.NewRecorder()

	r.ParseForm()
	router.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestRouter_ProxyURLCode(t *testing.T) {
	linkStore, closeConn, err := createLinkStore()
	if err != nil {
		t.Error(err)
	}
	defer closeConn()

	redisStore, cleanRedis := createRedisStore()
	defer cleanRedis()

	// insert link in database
	link, _ := linkStore.PersistURL("http://www.google.com")

	router := NewRouter(linkStore, redisStore)
	r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", encode(link.ID)), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusMovedPermanently)
	assert.Contains(t, w.Body.String(), link.URL)
}

func createLinkStore() (*store.LinkStore, func() error, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", "db", "postgres", "test")
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	return store.NewLinkStore(db), db.Close, nil
}

func createRedisStore() (*store.RedisStore, func() *redis.StatusCmd) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return store.NewRedisStore(client), client.FlushAll
}
