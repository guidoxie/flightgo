package api

import (
	"encoding/json"
	"flightgo/controllers/server"
	"flightgo/models"
	"flightgo/redis"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Airports struct {
	beego.Controller
}

// 获取全部机场名
func (c *Airports) Get() {

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisString(c.Ctx, "airports") {
		return
	}

	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	// 直接生成固定大小的slice
	shortNames := make([]string, 9518)

	var datas []models.Airports
	o := orm.NewOrm()
	qs := o.QueryTable("airports")

	qs.Limit(-1).All(&datas, "short_name", "iata")

	for i, data := range datas {
		shortNames[i] = data.Iata + " - " + data.ShortName
	}
	resp.Data = shortNames

	b, err := json.Marshal(resp)
	if err != nil {
		resp.Code = -1
		resp.Message = err.Error()
		resp.Data = b
		server.ServeJson(c.Ctx, resp)
		return
	}
	redis.RedisClient().SetNX("airports", b, 0)
	server.ServeJson(c.Ctx, b)
}
