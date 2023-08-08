package configs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"load-balancer/commons"
	"load-balancer/configs/database"
	"time"
)

// Proxy Pattern
type LoadBalancer struct {
	rateLimiter              map[string]int
	maxAllowRequestPerSecond int
}

type CounterLimitReq struct {
	Counters  int    `json:"counters"`
	TimeLimit string `json:"time_limit"`
}

func NewLoadBalancerServer(c *commons.Config) *LoadBalancer {
	return &LoadBalancer{
		rateLimiter:              make(map[string]int),
		maxAllowRequestPerSecond: c.RateLimit,
	}
}

func (lb *LoadBalancer) HandleRequest(ip string) error {
	allowed := lb.checkRateLimiter(ip)

	if !allowed {
		return errors.New("rate limit exceeded")
	}

	return nil
}

func (lb *LoadBalancer) checkRateLimiter(ip string) bool {
	if lb.rateLimiter[ip] == 0 {
		lb.rateLimiter[ip] = 1
	}

	if lb.rateLimiter[ip] > lb.maxAllowRequestPerSecond {
		return false
	}

	lb.rateLimiter[ip]++

	var data CounterLimitReq

	key := string(fmt.Sprintf("caches:ips:%s", ip))
	value, err := database.RedisClient.Get(context.Background(), key).Bytes()

	json.Unmarshal(value, &data)

	if err != nil || data.TimeLimit < time.Now().String() {
		data.Counters = 0
		data.TimeLimit = time.Now().Add(time.Second * 60).String()
	}
	data.Counters++

	if data.Counters > lb.maxAllowRequestPerSecond && data.TimeLimit > time.Now().String() {
		return false
	} else {
		jsonData, err := json.Marshal(&data)
		if err != nil {
			panic(err)
		}
		err = database.RedisClient.SetEx(context.Background(), key, jsonData, time.Minute*10).Err()
		if err != nil {
			panic(err)
		}
	}

	return true
}
