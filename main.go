package main

import (
	_ "flightgo/routers"
	"github.com/astaxie/beego"
	_ "flightgo/redis"
)

func main() {
	beego.Run()

}
