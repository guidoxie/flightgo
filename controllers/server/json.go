package server

import (
	"encoding/json"
	"flightgo/redis"
	"github.com/astaxie/beego/context"
)

// 从Redis读取key的数据,如果有,返回true并把数据返回给用户
//
// 字符串
func RedisString(ctx *context.Context, key string) bool {
	cli := redis.RedisClient()
	data, err := cli.Get(key).Bytes()
	// redis没有数据
	if err != nil {
		return false
	}
	ServeJson(ctx, data)
	return true
}

// 从Redis读取key,field的数据,如果有,返回true并把数据返回给用户
//
// 哈希表
func RedisHex(ctx *context.Context, key string, field string) bool {
	cli := redis.RedisClient()
	data, err := cli.HGet(key, field).Bytes()
	// redis没有数据
	if err != nil {
		return false
	}
	ServeJson(ctx, data)
	return true
}

// json 数据返回
//
// 重新定义json数据返回,当数据类型为[]byte时,假设已经进行json格式的序列化,不再进行编码
func ServeJson(ctx *context.Context, data interface{}) {

	var (
		resp []byte
		ok   bool
	)
	ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	if resp, ok = data.([]byte); ok {
		ctx.Output.Body(resp)
	} else {
		resp, _ := json.Marshal(data)
		ctx.Output.Body(resp)
	}
}
