package view

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/guidoxie/flightgo/models"
	"github.com/guidoxie/flightgo/util"
	"strconv"
	"strings"
)

type Search struct {
	beego.Controller
}

func (c *Search) Get() {

	o := c.Ctx.Input.Param(":o")
	d := c.Ctx.Input.Param(":d")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	data := searchByOd(o, d, year, month, day)
	c.Data["data"] = data
	if len(data) > 6 {
		c.TplName = "search.html"
	} else {
		c.TplName = "search_none.html"
	}

}

func searchByOd(o, d, year, month, day string) map[string]interface{} {

	res := make(map[string]interface{})

	res["year"] = year
	res["month"] = month
	res["day"] = day

	if len(month) < 2 {
		month = "0" + month
	}
	if len(day) < 2 {
		day = "0" + day
	}
	res["date"] = year + "-" + month + "-" + day
	res["origin"] = o
	res["destination"] = d

	var origin, destination models.Airports

	or := orm.NewOrm()
	origin.ShortName = o
	destination.ShortName = d
	or.Read(&origin, "short_name")
	or.Read(&destination, "short_name")

	if origin.Iata != "" && destination.Iata != "" {
		var flightDatas []orm.Params
		var list []map[string]interface{}

		qsFD := or.QueryTable("flight_datas").Limit(-1)
		qsFD.Filter("date", res["date"]).Filter("origin", origin.Iata).
			Filter("destination", destination.Iata).Values(&flightDatas)

		for _, fd := range flightDatas {

			var airline models.Airlines
			var state models.State

			iata, _ := fd["Airline"].(string)
			id, _ := fd["State"].(int64)

			airline.Iata = iata
			or.Read(&airline, "iata")

			state.Id = int(id)
			or.Read(&state)

			dataMap := make(map[string]interface{})
			dataMap["log"] = airline.Logo
			dataMap["airline"] = airline.ShortName
			dataMap["ident"] = fd["Ident"]
			dataMap["aircraft"] = fd["Aircraft"]
			dataMap["origin"] = o
			dataMap["destination"] = d
			dataMap["state"] = state.Name
			// 飞行时间
			aLanding, okL := fd["ALanding"].(uint64)
			aTakeOff, okT := fd["ATakeOff"].(uint64)
			if okL && okT {
				dataMap["fly_time"] = util.CountFlyTime(aLanding - aTakeOff)

			} else {
				dataMap["fly_time"] = "--"
			}
			list = append(list, dataMap)

		}
		if len(list) > 0 {
			res["length"] = len(list)
			res["result"] = list
		}

	}
	return res
}

func (c *Search) Search() {

	searchType := c.GetString("search-type")
	date := c.GetString("date")
	slice := strings.Split(date, "-")
	year := slice[0]
	month := slice[1]
	day := slice[2]

	if searchType == "0" {
		ident := strings.ToUpper(c.GetString("flight_no"))
		url := fmt.Sprintf("/flight/%s/%s/%s/%s/", ident, year, month, day)
		c.Redirect(url, 302)

	} else if searchType == "1" {
		o := c.GetString("flight_od_o")
		d := c.GetString("flight_od_d")
		url := fmt.Sprintf("/search/%s/%s/%s/%s/%s/", o, d, year, month, day)
		c.Redirect(url, 302)
	} else {
		c.Ctx.WriteString("error")
	}

}

func (c *Search) SearchByDate() {
	date := c.GetString("date")

	d := strings.Split(date, "-")

	if len(d) < 3 || len(d[1]) > 2 || len(d[2]) > 2 {
		c.Ctx.WriteString("error")
		return
	}

	_, err0 := strconv.Atoi(d[0])
	_, err1 := strconv.Atoi(d[1])
	_, err2 := strconv.Atoi(d[2])

	if err0 == nil && err1 == nil && err2 == nil {

		details := Details{}
		// /details/:year:int/:month:int/:day:int/
		url := details.URLFor("Details.Get", ":year",
			d[0], ":month", d[1], ":day", d[2])
		c.Redirect(url, 302)
		return
	}

	c.Ctx.WriteString("error")

}
