package api

import (
	"encoding/json"
	"flightgo/controllers/server"
	"flightgo/models"
	"flightgo/redis"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 准点率
type OnTime struct {
	beego.Controller
}

func (c *OnTime) Get() {

	ident := c.Ctx.Input.Param(":ident")

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "on_time", ident) {
		return
	}
	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	o := orm.NewOrm()
	var fd []orm.Params
	qsFD := o.QueryTable("flight_datas").Limit(30)
	_, err := qsFD.Filter("ident", ident).Values(&fd, "state", "a_landing", "s_landing")
	var count float64
	var overdue float64

	if err == nil {

		for _, f := range fd {
			aLanding, _ := f["ALanding"].(uint64)
			sLanding, _ := f["SLanding"].(uint64)

			if count < 30 {
				if f["State"] == 4 {
					beego.Info(4)
					count += 1
					overdue += 1
				} else if sLanding != 0 && aLanding != 0 {
					if (float64(aLanding)-float64(sLanding))/60 > 15 {
						overdue += 1
					}
					count += 1
				}
			} else {
				break
			}

		}
	} else {
		resp.Code = -1
		resp.Message = "发生错误"
	}

	data := make(map[string]float64)

	if count != 0 {
		data["on_time"] = (1 - (overdue / count)) * 100
	} else {
		data["on_time"] = 0
	}

	resp.Data = data
	//c.Data["json"] = resp
	//c.ServeJSON()

	b, err := json.Marshal(resp)
	if err != nil {
		resp.Code = -1
		resp.Message = err.Error()
		resp.Data = b
		server.ServeJson(c.Ctx, resp)
		return
	}
	redis.RedisClient().HSetNX("on_time", ident, b)
	server.ServeJson(c.Ctx, b)
}
