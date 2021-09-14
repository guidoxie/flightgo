package main

import (
	"github.com/astaxie/beego"
	_ "github.com/guidoxie/flightgo/redis"
	_ "github.com/guidoxie/flightgo/routers"
)

func main() {
	beego.Run()

}
