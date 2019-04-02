package view

import (
	"flightgo/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

type Flight struct {
	beego.Controller
}

func (c *Flight) Get() {

	ident := c.Ctx.Input.Param(":isIdent")
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	data := getdatas(ident, year, month, day)
	c.Data["data"] = data

	c.TplName = "live.html"
	if len(data) > 6 {
		c.TplName = "flight.html"
	} else {
		c.TplName = "flight_none.html"
	}

}

// 取得航班数据
func getdatas(isIdent, year, month, day string) (res map[string]interface{}) {

	res = make(map[string]interface{})

	if len(month) < 2 {
		month = "0" + month
	}
	if len(day) < 2 {
		day = "0" + day
	}

	res["in_ident"] = isIdent
	res["year"] = year
	res["month"] = month
	res["day"] = day
	res["date"] = year + "-" + month + "-" + day

	var ident string

	o := orm.NewOrm()
	var r orm.RawSeter
	r = o.Raw("select ident from flight_number where iataIdent=?", isIdent)
	r.QueryRow(&ident)

	if ident != "" {
		isIdent = ident
	}
	res["ident"] = isIdent
	var flightDatas models.FlightDatas
	qsFD := o.QueryTable("flight_datas").Limit(-1)
	err := qsFD.Filter("ident", isIdent).Filter("date", res["date"]).One(&flightDatas)

	if err == nil {
		var iataIdent string
		var r orm.RawSeter
		r = o.Raw("select iataIdent from flight_number where ident=?", isIdent)
		r.QueryRow(&iataIdent)
		res["iataident"] = iataIdent
		if flightDatas.STakeOff != 0 {
			t := strings.Split(time.Unix(int64(flightDatas.STakeOff), 0).String(), ":")
			res["s_take_off"] = t[0] + ":" + t[1]

		} else {
			res["s_take_off"] = "--"
		}

		if flightDatas.ATakeOff != 0 {
			t := strings.Split(time.Unix(int64(flightDatas.ATakeOff), 0).String(), ":")
			res["a_take_off"] = t[0] + ":" + t[1]

		} else {
			res["a_take_off"] = "--"
		}

		if flightDatas.SLanding != 0 {
			t := strings.Split(time.Unix(int64(flightDatas.SLanding), 0).String(), ":")
			res["s_landing"] = t[0] + ":" + t[1]

		} else {
			res["s_landing"] = "--"
		}

		if flightDatas.ALanding != 0 {
			t := strings.Split(time.Unix(int64(flightDatas.ALanding), 0).String(), ":")
			res["a_landing"] = t[0] + ":" + t[1]

		} else {
			res["a_landing"] = "--"
		}

		var state models.State
		state.Id = flightDatas.State
		o.Read(&state)
		res["state"] = state.Name

		var airline models.Airlines
		airline.Iata = flightDatas.Airline
		o.Read(&airline, "iata")
		dataMap := make(map[string]string)
		dataMap["name"] = airline.ShortName
		dataMap["href"] = airline.Url
		dataMap["log"] = airline.Logo
		res["airline"] = dataMap

		var origin models.Airports
		var destination models.Airports
		origin.Iata = flightDatas.Origin
		destination.Iata = flightDatas.Destination
		o.Read(&origin, "iata")
		o.Read(&destination, "iata")
		res["origin"] = origin.ShortName
		res["destination"] = destination.ShortName

		if flightDatas.Aircraft != "" {

			var planeImgs orm.ParamsList
			qsPI := o.QueryTable("plane_img").Limit(-1)
			qsPI.Filter("aircraft", flightDatas.Aircraft).ValuesFlat(&planeImgs, "url")
			res["plane_imgs"] = planeImgs

			var aircratf models.Aircraft
			aircratf.Type = flightDatas.Aircraft
			o.Read(&aircratf, "type")
			res["aircraft"] = flightDatas.Aircraft + "(" + aircratf.FriendlyType + ")"

		} else {
			res["aircraft"] = "--"
			str := make([]string, 1)
			str[0] = "https://photos-e1.flightcdn.com/photos/retriever/784d99ab4023fe7c6d6328d4116de938c6d3d04d"
			res["plane_imgs"] = str
		}
		res["on_time"] = 0

	}

	return

}
