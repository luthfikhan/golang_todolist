package helper

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"todolist/model"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var TaskCacheKeyPrefix = "CacheTask:"

type Cache struct {
	Ctx context.Context
	Rdb *redis.Client
}

func (c *Cache) logCache(id int, err error, procces string) {
	if err == nil {
		Log.Info(gin.H{"taskid": id}, procces+" successfully")
	} else {
		Log.Error(gin.H{"taskid": id, "err": err.Error()}, procces+" error")
	}
}

func (c *Cache) SetTaskCache(id int, task *model.Task) {
	stringValue, _ := json.Marshal(task)

	err := c.Rdb.SetEx(c.Ctx, TaskCacheKeyPrefix+strconv.Itoa(id), stringValue, time.Duration(5*time.Minute)).Err()

	c.logCache(id, err, "Set task cache")
}

func (c *Cache) GetTaskCache(id int) *model.Task {
	var task *model.Task

	cache, err := c.Rdb.Get(c.Ctx, TaskCacheKeyPrefix+strconv.Itoa(id)).Result()

	if err != nil {
		c.logCache(id, err, "Get task cache")
		return nil
	}
	err = json.Unmarshal([]byte(cache), &task)

	if err != nil {
		c.logCache(id, err, "Get task cache (Unmarshal)")
		return nil
	}

	c.logCache(id, nil, "Get task cache")
	return task
}

func (c *Cache) DeleteTaskCache(id int) {
	err := c.Rdb.Del(c.Ctx, TaskCacheKeyPrefix+strconv.Itoa(id)).Err()

	c.logCache(id, err, "Delete task cache")
}
