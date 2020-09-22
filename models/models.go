package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 机型表
type Aircraft struct {
	Id           int    `json:"id"`
	Type         string `json:"type"`          // 机型
	FriendlyType string `json:"friendly_type"` // 机型全称
}

// 航空公司表
type Airlines struct {
	Id        int    `json:"id"`
	Iata      string `json:"iata"`       // 航司两字代号
	Icao      string `json:"icao"`       // 航司三字代号
	FullName  string `json:"full_name"`  // 航司全称
	ShortName string `json:"short_name"` // 航司简称
	Type      string `json:"type"`       // 国外（out）和国内（in）
	Logo      string `json:"logo"`       // 航司logo链接
	Url       string `json:"url"`        // 航司首页
}

// 机场表
type Airports struct {
	Id        int     `json:"id"`
	Iata      string  `json:"iata"`       // 机场三字代号
	Icao      string  `json:"icao"`       // 机场四字代号
	Lat       float64 `json:"lat"`        // 纬度
	Lon       float64 `json:"lon"`        // 经度
	ShortName string  `json:"short_name"` // 机场简称
	Country   string  `json:"country"`    // 机场所属国家
}

// 航班信息表
type FlightDatas struct {
	Id            int    `json:"id,omitempty"`             // id
	Ident         string `json:"ident,omitempty"`          // 航班号
	Date          string `json:"date,omitempty"`           // 日期
	STakeOff      uint64 `json:"s_take_off,omitempty"`     // 计划起飞时间
	ETakeOff      uint64 `json:"e_take_off,omitempty"`     // 预测起飞时间
	ATakeOff      uint64 `json:"a_take_off,omitempty"`     // 实际起飞时间
	SLanding      uint64 `json:"s_landing,omitempty"`      // 计划降落时间
	ELanding      uint64 `json:"e_landing,omitempty"`      // 预测降落时间
	ALanding      uint64 `json:"a_landing,omitempty"`      // 实际降落时间
	State         int    `json:"state,omitempty"`          // 航班状态
	Aircraft      string `json:"aircraft,omitempty"`       // 机型
	Airline       string `json:"airline"`                  // 航司
	Origin        string `json:"origin,omitempty"`         // 起飞点
	Destination   string `json:"destination,omitempty"`    // 降落点
	Distance      uint32 `json:"distance,omitempty"`       // 距离
	FriendlyIdent string `json:"friendly_ident,omitempty"` // 航班号全称
}

// 航班号表
type FlightNumber struct {
	Id        int    `json:"id"`
	Ident     string `json:"ident"`      // 三字航班号
	IataIdent string `json:"iata_ident"` // 两字航班号

}

type PlaneImg struct {
	Id       int    `json:"id"`
	Aircraft string `json:"aircraft"`
	Url      string `json:"url"`
}

// 航班状态表
type State struct {
	Id   int    `json:"id"`   // 状态id
	Type string `json:"type"` // 状态类型
	Name string `json:"name"` // 状态名称
}

// 返回的数据格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 分页数据
type PageData struct {
	FlightDatas
	FlyTime string `json:"fly_time"`
}

// 初始化数据模型
func init() {

	// 读取配置文件的数据库参数
	mysqlUser := beego.AppConfig.String("mysqluser")
	mysqlPass := beego.AppConfig.String("mysqlpass")
	mysqlUrl := beego.AppConfig.String("mysqlurl")
	mysqlDB := beego.AppConfig.String("mysqldb")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlUrl, mysqlDB)

	// 注册数据模型
	orm.RegisterModel(new(Aircraft))
	orm.RegisterModel(new(FlightDatas))
	orm.RegisterModel(new(State))
	orm.RegisterModel(new(Airlines))
	orm.RegisterModel(new(Airports))
	orm.RegisterModel(new(FlightNumber))
	orm.RegisterModel(new(PlaneImg))

	// 注册数据库源
	orm.RegisterDataBase("default", "mysql", dataSource)
}
