package api

import (
	"encoding/json"
	"flightgo/controllers/server"
	"flightgo/models"
	"flightgo/redis"
	"flightgo/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
)

// 根据年月日返回一天的航班数据
// "ident"
// "airline"
// "origin"
// "destination"
type AllDayFlight struct {
	beego.Controller
}

func (c *AllDayFlight) Get() {

	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")

	date := year + month + day

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "all_day_flight", date) {
		return
	}
	resp := models.Response{
		Code:    0,
		Message: "success",
	}

	flightDatas := make(map[string][]models.FlightDatas, 0)
	var datas = make([]models.FlightDatas, 0)

	o := orm.NewOrm()

	qs := o.QueryTable("flight_datas")
	// Limit()查询返回最大行数 -1 不限制
	_, err := qs.Limit(-1).Filter("date", date).All(&datas, "ident",
		"airline", "origin", "destination")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
	}

	for _, data := range datas {
		if data.Origin != "" && data.Destination != "" {
			if data.Airline == "" {
				flightDatas["Null"] = append(flightDatas["Null"], data)
			} else {

				flightDatas[data.Airline] = append(flightDatas[data.Airline], data)

			}
		}
	}
	resp.Data = flightDatas
	//c.Data["json"] = resp
	//
	//c.ServeJSON()

	b, err := json.Marshal(resp)
	if err != nil {
		resp.Code = -1
		resp.Message = err.Error()
		resp.Data = b
		server.ServeJson(c.Ctx, resp)
		return
	}
	redis.RedisClient().HSetNX("all_day_flight", date, b)
	server.ServeJson(c.Ctx, b)
}

// 获取一个航班的3个月的数据
type ThereMonthsFlight struct {
	beego.Controller
}

func (c *ThereMonthsFlight) Get() {

	ident := c.Ctx.Input.Param(":ident")

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "there_months_flight", ident) {
		return
	}

	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	o := orm.NewOrm()
	var datas = make([]*models.FlightDatas, 0)

	qs := o.QueryTable("flight_datas")

	// 限制90个,3个月
	_, err := qs.Limit(90).Filter("ident", ident).All(&datas)
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
	}
	resp.Data = datas
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
	redis.RedisClient().HSetNX("there_months_flight", ident, b)
	server.ServeJson(c.Ctx, b)
}

// 特定航班号的数据
type OneFlight struct {
	beego.Controller
}

func (c *OneFlight) Get() {

	ident := c.Ctx.Input.Param(":ident")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	date := year + "-" + month + "-" + day

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "one_flight", ident+date) {
		return
	}
	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	var flightData = models.FlightDatas{}
	o := orm.NewOrm()
	err := o.QueryTable("flight_datas").Filter("ident", ident).
		Filter("date", date).One(&flightData)
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
	}
	resp.Data = flightData
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
	redis.RedisClient().HSetNX("one_flight", ident+date, b)
	server.ServeJson(c.Ctx, b)
}

// 航班分页数据
type PageFlight struct {
	beego.Controller
}

func (c *PageFlight) Get() {

	airline := c.Ctx.Input.Param(":airline")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	start := c.Ctx.Input.Param(":page")
	date := year + month + day
	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "page_flight", airline+date+start) {
		return
	}
	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	pageData := make([]models.PageData, 0)

	var fd = make([]*models.FlightDatas, 0)

	page, err := strconv.Atoi(start)
	if err != nil || page < 0 {
		resp.Code = -1
		resp.Message = err.Error()
		resp.Data = pageData
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	qsFD := o.QueryTable("flight_datas").Limit(10, page*10)

	if airline != "ONTHER" {

		_, err = qsFD.Filter("date", date).Filter("airline", airline).All(&fd,
			"ident", "aircraft", "origin", "destination",
			"state", "a_take_off", "a_landing")

	} else if airline == "ONTHER" {

		_, err = qsFD.Filter("date", date).Exclude("airline", "CA").Exclude("airline", "AA").
			Exclude("airline", "BA").Exclude("airline", "AK").Exclude("airline", "MU").
			Exclude("airline", "UA").All(&fd, "ident", "aircraft", "origin",
			"destination", "state", "a_take_off", "a_landing")

	}
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = pageData
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	qsAP := o.QueryTable("airports")

	for _, f := range fd {
		var origin models.Airports
		var destination models.Airports
		qsAP.Filter("iata", f.Origin).One(&origin, "short_name")
		f.Origin = origin.ShortName

		qsAP.Filter("iata", f.Destination).One(&destination, "short_name")
		f.Destination = destination.ShortName

		pg := models.PageData{}
		pg.Ident = f.Ident
		pg.Aircraft = f.Aircraft
		pg.Origin = f.Origin
		pg.Destination = f.Destination
		pg.State = f.State
		pg.ATakeOff = f.ATakeOff
		pg.ALanding = f.ALanding
		// 飞行时间

		if f.ALanding > f.ATakeOff && f.ATakeOff != 0 && f.ALanding != 0 {
			pg.FlyTime = util.CountFlyTime(f.ALanding - f.ATakeOff)

		} else {
			pg.FlyTime = "--:--"
		}
		pageData = append(pageData, pg)
	}

	resp.Data = pageData
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
	redis.RedisClient().HSetNX("page_flight", airline+date+start, b)
	server.ServeJson(c.Ctx, b)

}
