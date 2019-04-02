package routers

import (
	"flightgo/controllers/api"
	"github.com/astaxie/beego"
)

func init() {
	// ----------------------api----------------------------------
	//
	// 航班数据
	beego.Router("/api/flights/:year:int/:month:int/:day:int/",
		&api.AllDayFlight{}) // 航班数据,查询条件是年月日,注意,此api不存在限制,返回数据可能过于庞大
	beego.Router("/api/flights/:ident:string/",
		&api.ThereMonthsFlight{}) // 航班数据,查询条件航班号,最多返回90天(过去三个月)
	beego.Router("/api/flights/:ident:string/:year:int/:month:int/:day:int/",
		&api.OneFlight{}) // 航班数据,查询条件航班号、年月日
	beego.Router("/api/ontime/:ident:string/", &api.OnTime{}) // 准点率

	beego.Router("/api/flights/chain/:year:int/:month:int/:day:int/",
		&api.Chain{})
	beego.Router("/api/flights/:airline:string/:year:int/:month:int/:day:int/:page:int/",
		&api.PageFlight{})

	// 机场
	beego.Router("/api/airports/",
		&api.Airports{}) // 机场名称

	// 数据聚合分析
	beego.Router("/api/airports/:airport([\u4e00-\u9fa5]+)/"+
		":year:int/:month:int/:day:int/flights/", &api.AirportAnalysis{})

	beego.Router("/api/airlines/:airline([\u4e00-\u9fa5]+)/"+
		":year:int/:month:int/:day:int/flights/", &api.AirlineAnalysis{})

	beego.Router("/api/country/:country:string/"+
		":year:int/:month:int/:day:int/flights/", &api.CountryAnalysis{})

}
