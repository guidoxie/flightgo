package routers

import (
	"github.com/astaxie/beego"
	"github.com/guidoxie/flightgo/controllers/view"
)

func init() {
	// ---------------------views----------------------------
	beego.Router("/", &view.Index{})                                                       // 主页
	beego.Router("/details/:year:int/:month:int/:day:int/", &view.Details{})               // 某天全球数据详情
	beego.Router("/analysis/", &view.Analysis{})                                           // 数据聚合分析
	beego.Router("/flight/:isIdent:string/:year:int/:month:int/:day:int/", &view.Flight{}) // 特定航班数据详情

	// ---------------------搜索相关-------------------
	beego.Router("/search/:o([\u4e00-\u9fa5]+)/:d([\u4e00-\u9fa5]+)/:year:int/:month:int/:day:int/",
		&view.Search{}) // 根据起降点搜索显示页
	beego.Router("/search/", &view.Search{}, "get:Search")            // 搜索
	beego.Router("/search/date/", &view.Search{}, "get:SearchByDate") // 根据日期搜索

}
