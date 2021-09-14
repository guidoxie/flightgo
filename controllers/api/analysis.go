package api

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/guidoxie/flightgo/controllers/server"
	"github.com/guidoxie/flightgo/models"
	"github.com/guidoxie/flightgo/redis"
	"github.com/guidoxie/flightgo/util"
)

// 聚合分析
type AirportAnalysis struct {
	beego.Controller
}

//获取机场聚合数据
func (c *AirportAnalysis) Get() {

	shortName := c.Ctx.Input.Param(":airport")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	date := year + "-" + month + "-" + day
	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "airport_analysis", shortName+date) {
		return
	}
	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	data := make(map[string]interface{})

	inData := make(map[string]interface{})

	// 初始化为空slice [],指定长度为0
	inOnTime := make([]map[string]string, 0)
	inCancel := make([]map[string]string, 0)
	inDelay := make([]map[string]string, 0)
	inNotReported := make([]map[string]string, 0)

	outData := make(map[string]interface{})

	// 初始化为空slice [],指定长度为0
	outOnTime := make([]map[string]string, 0)
	outCancel := make([]map[string]string, 0)
	outDelay := make([]map[string]string, 0)
	outNotReported := make([]map[string]string, 0)

	var airport models.Airports
	airport.ShortName = shortName

	o := orm.NewOrm()
	o.Read(&airport, "short_name")
	data["airport"] = airport.Iata
	outCount := 0
	inCount := 0
	code := 0

	qsFD := o.QueryTable("flight_datas").Limit(-1)

	var outFD []orm.Params
	var inFD []orm.Params

	_, err := qsFD.Filter("date", date).Filter("origin", airport.Iata).Values(&outFD,
		"ident", "destination", "state", "a_take_off", "a_landing", "s_landing")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = data
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	_, err = qsFD.Filter("date", date).Filter("destination", airport.Iata).Values(&inFD,
		"ident", "origin", "state", "a_take_off", "a_landing", "s_landing")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = data
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	for _, fd := range outFD {
		outCount += 1

		if code == 0 {
			code = 1
		}

		dataMap := make(map[string]string)
		destination, _ := fd["Destination"].(string)
		ident, _ := fd["Ident"].(string)
		dataMap[destination] = ident
		state, _ := fd["State"].(int64)

		if state == 1 {
			sLanding, _ := fd["SLanding"].(uint64)
			aLanding, _ := fd["ALanding"].(uint64)
			if sLanding != 0 && aLanding != 0 {

				aLanding, _ := fd["ALanding"].(uint64)
				sLanding, _ := fd["SLanding"].(uint64)
				time := aLanding - sLanding

				if time > 900 {
					outDelay = append(outDelay, dataMap)
				} else {
					outOnTime = append(outOnTime, dataMap)
				}
			}

		} else if state == 2 || state == 3 {
			outNotReported = append(outNotReported, dataMap)

		} else if state == 4 {
			outCancel = append(outCancel, dataMap)
		}

	}

	for _, fd := range inFD {
		inCount += 1

		if code == 0 {
			code = 1
		}

		dataMap := make(map[string]string)
		origin, _ := fd["Origin"].(string)
		ident, _ := fd["Ident"].(string)
		dataMap[origin] = ident
		state, _ := fd["State"].(int64)

		if state == 1 {
			sLanding, _ := fd["SLanding"].(uint64)
			aLanding, _ := fd["ALanding"].(uint64)

			if sLanding != 0 && aLanding != 0 {

				time := aLanding - sLanding

				if time > 900 {
					inDelay = append(inDelay, dataMap)
				} else {
					inOnTime = append(inOnTime, dataMap)
				}
			}

		} else if state == 2 || state == 3 {
			inNotReported = append(inNotReported, dataMap)

		} else if state == 4 {
			inCancel = append(inCancel, dataMap)
		}

	}

	outData["准点"] = outOnTime
	outData["取消"] = outCancel
	outData["延误"] = outDelay
	outData["未上报"] = outNotReported

	inData["准点"] = inOnTime
	inData["取消"] = inCancel
	inData["延误"] = inDelay
	inData["未上报"] = inNotReported

	data["out"] = outData
	data["in"] = inData
	data["out_count"] = outCount
	data["in_count"] = inCount
	data["code"] = code

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
	redis.RedisClient().HSetNX("airport_analysis", shortName+date, b)
	server.ServeJson(c.Ctx, b)

}

// 获取航司聚合数据
type AirlineAnalysis struct {
	beego.Controller
}

func (c *AirlineAnalysis) Get() {

	shortName := c.Ctx.Input.Param(":airline")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	date := year + "-" + month + "-" + day

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "airline_analysis", shortName+date) {
		return
	}

	resp := models.Response{
		Code:    0,
		Message: "success",
	}
	data := make(map[string]interface{})
	var airline models.Airlines
	code := 0
	count := 0
	result := make(map[string]interface{})
	// 初始化为空slice [],指定长度为0
	onTime := make([]map[string][]string, 0)
	cancel := make([]map[string][]string, 0)
	delay := make([]map[string][]string, 0)
	notReported := make([]map[string][]string, 0)

	o := orm.NewOrm()
	airline.ShortName = shortName
	o.Read(&airline, "short_name")

	if airline.ShortName != "" {
		data["airline"] = shortName

		var flightData []orm.Params

		qsFD := o.QueryTable("flight_datas").Limit(-1)
		_, err := qsFD.Filter("date", date).Filter("airline", airline.Iata).Values(&flightData,
			"ident", "origin", "destination", "state",
			"a_take_off", "a_landing", "s_landing")
		if err != nil {
			resp.Code = -1
			resp.Message = "发生错误"
			resp.Data = data
			c.Data["json"] = resp
			c.ServeJSON()
			return
		}
		for _, fd := range flightData {
			count += 1
			if code == 0 {
				code = 1
			}

			ident, _ := fd["Ident"].(string)
			origin, _ := fd["Origin"].(string)
			destination, _ := fd["Destination"].(string)
			list := make([]string, 2)
			list[0] = origin
			list[1] = destination

			dataMap := make(map[string][]string)
			dataMap[ident] = list
			state, _ := fd["State"].(int64)

			if state == 1 {
				sLanding, _ := fd["SLanding"].(uint64)
				aLanding, _ := fd["ALanding"].(uint64)

				if sLanding != 0 && aLanding != 0 {

					time := aLanding - sLanding

					if time > 900 {
						delay = append(delay, dataMap)
					} else {
						onTime = append(onTime, dataMap)
					}
				}

			} else if state == 2 || state == 3 {
				notReported = append(notReported, dataMap)

			} else if state == 4 {
				cancel = append(cancel, dataMap)
			}

		}

	}

	result["准点"] = onTime
	result["取消"] = cancel
	result["延误"] = delay
	result["未上报"] = notReported
	data["result"] = result
	data["code"] = code
	data["count"] = count

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
	redis.RedisClient().HSetNX("airline_analysis", shortName+date, b)
	server.ServeJson(c.Ctx, b)
}

//获取国家聚合数据
type CountryAnalysis struct {
	beego.Controller
}

func (c *CountryAnalysis) Get() {

	country := c.Ctx.Input.Param(":country")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	date := year + "-" + month + "-" + day

	// 如果redis有数据则不查询数据库,直接返回
	if server.RedisHex(c.Ctx, "country_analysis", country+date) {
		return
	}

	resp := models.Response{
		Code:    0,
		Message: "success",
	}

	data := make(map[string]interface{})
	var airport orm.ParamsList
	o := orm.NewOrm()
	qsAP := o.QueryTable("airports").Limit(-1)
	_, err := qsAP.Filter("country", country).ValuesFlat(&airport, "iata")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = data
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	//不跨境
	var flightData []orm.Params
	qsFD := o.QueryTable("flight_datas").Limit(-1).Filter("date", date)
	_, err = qsFD.Values(&flightData, "ident", "origin", "destination", "state",
		"a_take_off", "a_landing", "s_landing")
	if err != nil {
		resp.Code = -1
		resp.Message = "发生错误"
		resp.Data = data
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	inSliceMap := util.InSliceMap(airport)

	nCBCount := 0
	defCount := 0
	code := 0

	notCBFD := make(map[string][]map[string][]string)
	defFD := make(map[string][]map[string][]string)

	// 初始化为空slice [],指定长度为0
	notCBOnTime := make([]map[string][]string, 0)
	notCBCancel := make([]map[string][]string, 0)
	notCBDelay := make([]map[string][]string, 0)
	notCBNotReported := make([]map[string][]string, 0)

	// 初始化为空slice [],指定长度为0
	defOnTime := make([]map[string][]string, 0)
	defCancel := make([]map[string][]string, 0)
	defDelay := make([]map[string][]string, 0)
	defNotReported := make([]map[string][]string, 0)

	for _, fd := range flightData {

		ident, _ := fd["Ident"].(string)
		origin, _ := fd["Origin"].(string)
		destination, _ := fd["Destination"].(string)
		state, _ := fd["State"].(int64)

		list := make([]string, 2)
		list[0] = origin
		list[1] = destination

		dataMap := make(map[string][]string)
		dataMap[ident] = list

		if inSliceMap[fd["Origin"]] && inSliceMap[fd["Destination"]] {
			nCBCount += 1
			if code == 0 {
				code = 1
			}

			if state == 1 {
				sLanding, _ := fd["SLanding"].(uint64)
				aLanding, _ := fd["ALanding"].(uint64)

				if sLanding != 0 && aLanding != 0 {

					time := aLanding - sLanding

					if time > 900 {
						notCBDelay = append(notCBDelay, dataMap)
					} else {
						notCBOnTime = append(notCBOnTime, dataMap)
					}
				}

			} else if state == 2 || state == 3 {
				notCBNotReported = append(notCBNotReported, dataMap)

			} else if state == 4 {
				notCBCancel = append(notCBCancel, dataMap)
			}

		} else if inSliceMap[fd["Origin"]] || inSliceMap[fd["Destination"]] {
			defCount += 1
			if code == 0 {
				code = 1
			}

			if state == 1 {
				sLanding, _ := fd["SLanding"].(uint64)
				aLanding, _ := fd["ALanding"].(uint64)

				if sLanding != 0 && aLanding != 0 {

					time := aLanding - sLanding

					if time > 900 {
						defDelay = append(defDelay, dataMap)
					} else {
						defOnTime = append(defOnTime, dataMap)
					}
				}

			} else if state == 2 || state == 3 {
				defNotReported = append(defNotReported, dataMap)

			} else if state == 4 {
				defCancel = append(defCancel, dataMap)
			}
		}
	}

	notCBFD["准点"] = notCBOnTime
	notCBFD["取消"] = notCBCancel
	notCBFD["延误"] = notCBDelay
	notCBFD["未上报"] = notCBNotReported

	defFD["准点"] = defOnTime
	defFD["取消"] = defCancel
	defFD["延误"] = defDelay
	defFD["未上报"] = defNotReported

	data["country"] = country
	data["not_c_b"] = notCBFD
	data["default"] = defFD
	data["n_c_b_count"] = nCBCount
	data["def_count"] = defCount
	data["code"] = code

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
	redis.RedisClient().HSetNX("country_analysis", country+date, b)
	server.ServeJson(c.Ctx, b)
}
