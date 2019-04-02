package view

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Details struct {
	beego.Controller
}

func (c *Details) Get() {
	year := c.Ctx.Input.Param(":year")
	month := c.Ctx.Input.Param(":month")
	day := c.Ctx.Input.Param(":day")
	date := year + "-" + month + "-" + day

	c.Data["date"] = date
	c.Data["year"] = year
	c.Data["month"] = month
	c.Data["day"] = day

	if CheckDay(date) {
		c.TplName = "details.html"
	} else {
		c.TplName = "details_none.html"
	}

}

// 根据特定日期,检查是否有数据
func CheckDay(date string) bool {

	o := orm.NewOrm()
	qs := o.QueryTable("flight_datas")
	if count, err := qs.Filter("date", date).Count(); err == nil && count > 0 {
		return true
	}
	return false
}
