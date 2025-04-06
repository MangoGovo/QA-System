package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"QA-System/internal/model"
	"QA-System/internal/pkg/redis" // 保留你自己项目中的 Redis 包
	"github.com/gin-gonic/gin"
	redisPkg "github.com/redis/go-redis/v9" // 添加 Redis 库
)

// GetUserLimit 获取用户的对该问卷的访问次数
func GetUserLimit(c context.Context, stu_id string, sid string, durationType string) (uint, error) {
	// 从 redis 中获取用户的对该问卷的访问次数, durationtype为dailyLimit或sumLimit
	item := "survey:" + sid + ":duration_type:" + durationType + ":stu_id:" + stu_id
	var limit uint
	err := redis.RedisClient.Get(c, item).Scan(&limit)
	return limit, err
}

// SetUserLimit 设置用户的对该问卷的单日访问次数
func SetUserLimit(c context.Context, stuId string, sid string, limit int, durationType string) error {
	// 设置用户的对该问卷的访问次数, durationtype为dailyLimit或sumLimit
	item := "survey:" + sid + ":duration_type:" + durationType + ":stu_id:" + stuId
	// 获取当前时间和第二天零点的时间
	now := time.Now()
	tomorrow := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	).Add(24 * time.Hour)
	duration := time.Until(tomorrow) // 计算当前时间到第二天零点的时间间隔
	err := redis.RedisClient.Set(c, item, limit, duration).Err()
	return err
}

// InscUserLimit 更新用户的对该问卷的访问次数+1
func InscUserLimit(c context.Context, stuId string, sid string, durationType string) error {
	// 更新用户的对该问卷的访问次数,durationtype为dailyLimit或sumLimit
	item := "survey:" + sid + ":duration_type:" + durationType + ":stu_id:" + stuId
	err := redis.RedisClient.Incr(c, item).Err()
	return err
}

// SetUserSumLimit 设置用户对该问卷的总访问次数
func SetUserSumLimit(c context.Context, stuId string, sid string, sumLimit int, durationType string) error {
	// 设置用户的对该问卷的访问次数, durationtype为dailyLimit或sumLimit
	item := "survey:" + sid + ":duration_type:" + durationType + ":stu_id:" + stuId
	// 获取当前时间到问卷截止的时间
	var survey *model.Survey
	survey, err := GetSurveyByUUID(sid)
	if err != nil {
		return err
	}
	endTime := survey.Deadline
	sumDuration := time.Until(endTime)
	err = redis.RedisClient.Set(c, item, sumLimit, sumDuration).Err()
	return err
}

// CheckLimit 判断限制次数
func CheckLimit(c *gin.Context, stuId string, survey *model.Survey, key string, limitVal uint) (bool, error) {
	if limitVal == 0 {
		return false, nil
	}
	limit, err := GetUserLimit(c, stuId, survey.UUID, key)
	if err != nil && !errors.Is(err, redisPkg.Nil) {
		return false, err
	}
	if err == nil && limit >= limitVal {
		return false, fmt.Errorf("%s已达上限", key)
	}
	return errors.Is(err, redisPkg.Nil), nil
}
