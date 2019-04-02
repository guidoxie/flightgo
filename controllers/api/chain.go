package api

import (
	"encoding/json"
	"flightgo/controllers/server"
	"flightgo/models"
	"flightgo/redis"
	"flightgo/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 获取中国机场热力图数据
type Chain struct {
	beego.Controller
}

func (c *Chain) Get() {

	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	date := year + month + day

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "chain", date) {
		return
	}

	resp := models.Response{
		Code:    0,
		Message: "success",
	}

	var data = make([]interface{}, 0)
	o := orm.NewOrm()
	qsFD := o.QueryTable("flight_datas")
	qsCA := o.QueryTable("airports")

	var chainAirports orm.ParamsList
	var flightDatas orm.ParamsList

	_, err := qsCA.Limit(-1).Filter("country", "CN").ValuesFlat(&chainAirports, "iata")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = data
		c.Data["json"] = resp

		c.ServeJSON()
		return
	}
	qsFD = qsFD.Limit(-1).Filter("date", date).
		Filter("destination__in", chainAirports).Filter("state__in", 1, 2, 3, 5)
	_, err = qsFD.ValuesFlat(&flightDatas, "destination")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = data
		c.Data["json"] = resp

		c.ServeJSON()
		return
	}
	count := make(map[string]int)
	for _, cn := range flightDatas {
		if cn, ok := cn.(string); ok {
			count[cn] += 1
		}
	}

	// 排序
	counts := util.SortMapByValue(count)
	for _, c := range counts {
		var s []interface{}
		s = append(s, c.Key)
		s = append(s, c.Value)
		data = append(data, s)
	}
	resp.Data = data

	b, err := json.Marshal(resp)
	if err != nil {
		resp.Code = -1
		resp.Message = err.Error()
		resp.Data = b
		server.ServeJson(c.Ctx, resp)
		return
	}
	redis.RedisClient().HSetNX("chain", date, b)
	server.ServeJson(c.Ctx, b)
}
