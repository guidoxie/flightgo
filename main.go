package main

import (
	_ "flightgo/redis"
	_ "flightgo/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()

}
