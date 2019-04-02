package view

import "github.com/astaxie/beego"

type Index struct {
	beego.Controller
}

func (c *Index) Get() {

	c.TplName = "index.html"
}
