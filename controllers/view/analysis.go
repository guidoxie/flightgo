package view

import (
	"github.com/astaxie/beego"
)

type Analysis struct {
	beego.Controller
}

func (c *Analysis) Get() {
	c.TplName = "analysis.html"
}
