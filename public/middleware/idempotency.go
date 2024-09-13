package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CacheEntry struct {
	sync.Mutex
	lockeAt time.Time
}

type IdempotencyCache struct {
	sync.RWMutex
	CacheEntrys map[string]*CacheEntry
	lockTimeout time.Duration
}

// 初始化幂等性缓存
func NewIdempotencyCache(lockTimeout time.Duration) *IdempotencyCache {
	return &IdempotencyCache{
		CacheEntrys: make(map[string]*CacheEntry),
		lockTimeout: lockTimeout,
	}
}

func (c *IdempotencyCache) Get(key string) (*CacheEntry, bool) {
	c.RLocker()
	defer c.RUnlock()
	entry, exist := c.CacheEntrys[key]
	if !exist || time.Since(entry.lockeAt) > c.lockTimeout {
		return nil, false
	}
	return entry, true
}

// 如果相同请求对应的锁还存在但是过期了,对锁进行续约
// 如果请求对应的锁不存在则创建对应的请求锁
func (c *IdempotencyCache) Set(key string) *CacheEntry {
	c.Lock()
	defer c.Unlock()
	_, exist := c.CacheEntrys[key]
	if exist {
		c.CacheEntrys[key].lockeAt = time.Now()
	} else {
		c.CacheEntrys[key] = &CacheEntry{
			lockeAt: time.Now(),
		}
	}
	return c.CacheEntrys[key]
}

// 限流:同一时间最多只能由100个协程释放锁,防止内存泄漏
var LimitCh = make(chan struct{}, 100)

func (c *IdempotencyCache) ClearUp() {
	c.Lock()
	defer c.Unlock()
	for key, value := range c.CacheEntrys {
		LimitCh <- struct{}{}
		go func(key string, entry *CacheEntry) {
			if time.Since(c.CacheEntrys[key].lockeAt) >= c.lockTimeout {
				entry.Lock()
				delete(c.CacheEntrys, key)
				<-LimitCh
			}
		}(key, value)
	}
}

// 生成全局缓存,处理接口幂等性
var IdempotencyMiddlewareCache = NewIdempotencyCache(12 * time.Hour)

// 生成缓存的key,由url+请求方法+请求内容
func GenerateCacheKey(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(strings.NewReader(string(body)))
	rawData := req.URL.String() + req.Method + string(body)
	hash := md5.Sum([]byte(rawData))
	return hex.EncodeToString(hash[:]), nil
}

// 在进入控制器之前调用
func IdempotencyMiddlewareBefore(ctx *context.Context) {
	if ctx.Input.Method() == http.MethodGet {
		return
	}
	req := ctx.Request
	//获取hashKey
	cacheKey, err := GenerateCacheKey(req)
	if err != nil {
		logs.Debug("failed to get request hashKey:", err)
		ctx.Output.SetStatus(400)
		jsonString := `{
			"err_code":400,
			"message":{
				"zh-CN": "获取请求的HashKey失败",
				"en-US":	"failed to get request hashKey"
			}
		}`
		var data map[string]interface{}
		json.Unmarshal([]byte(jsonString), &data)
		ctx.Output.JSON(data, false, false)
		panic(errors.New("exit"))
	}
	entry, ok := IdempotencyMiddlewareCache.Get(cacheKey)
	if ok {
		entry.Lock()
	} else {
		entry = IdempotencyMiddlewareCache.Set(cacheKey)
		entry.Lock()
	}
}

// 在进入控制器之后调用
func IdempotencyMiddlewareAfter(ctx *context.Context) {
	if ctx.Request.Method == http.MethodGet {
		return
	}
	req := ctx.Request
	cacheKey, _ := GenerateCacheKey(req)
	EntryUnlock(cacheKey)
}

func EntryUnlock(cacheKey string) {
	entry, exist := IdempotencyMiddlewareCache.Get(cacheKey)
	if exist {
		entry.Unlock()
	}
}

// 前提:请求在controller层进行失败,抛出错误中断请求
// 捕获中断错误,然后进行对请求锁的释放流程
var ProcessRequest = func(ctx *context.Context, config *beego.Config) {
	if r := recover(); r != nil {
		if ctx.Request.Method == http.MethodGet {
			return
		}
		if err, ok := r.(error); ok && err == beego.ErrAbort {
			req := ctx.Request
			cacheKey, _ := GenerateCacheKey(req)
			EntryUnlock(cacheKey)
		}
	}
}
