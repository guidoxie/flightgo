package tests

import (
	"flightgo/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestDB(t *testing.T) {

	o := orm.NewOrm()
	aircraft := models.Aircraft{}
	planeImg := models.PlaneImg{}
	flightNumber := models.FlightNumber{}
	airpost := models.Airports{}
	airlines := models.Airlines{}
	flightDatas := models.FlightDatas{}
	state := models.State{}

	aircraft.Id = 1
	if err := o.Read(&aircraft); err != nil {
		t.Error("aircraft 查询失败", err)
	}

	planeImg.Id = 1
	if err := o.Read(&planeImg); err != nil {
		t.Error("planeImg查询失败", err)
	}

	flightNumber.Id = 1
	if err := o.Read(&flightNumber); err != nil {
		t.Error("flightNumber 查询失败", err)
	}

	airpost.Id = 1
	if err := o.Read(&airpost); err != nil {
		t.Error("airpost 查询失败", err)
	}

	airlines.Id = 1
	if err := o.Read(&airlines); err != nil {
		t.Error("airlines 查询失败", err)
	}

	flightDatas.Id = 1
	if err := o.Read(&flightDatas); err != nil {
		t.Error("flightDatas 查询失败", err)
	}

	state.Id = 1
	if err := o.Read(&state); err != nil {
		t.Error("state 查询失败", err)
	}

}

// 初始化数据模型
func init() {

	// 读取配置文件的数据库参数
	mysqlUser := "root"
	mysqlPass := "123456789"
	mysqlUrl := "127.0.0.1:3306"
	mysqlDB := "flightgo"

	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		mysqlUser, mysqlPass, mysqlUrl, mysqlDB)

	// 注册数据模型
	orm.RegisterModel(new(models.Aircraft))
	orm.RegisterModel(new(models.FlightDatas))
	orm.RegisterModel(new(models.State))
	orm.RegisterModel(new(models.Airlines))
	orm.RegisterModel(new(models.Airports))
	orm.RegisterModel(new(models.FlightNumber))
	orm.RegisterModel(new(models.PlaneImg))

	// 注册数据库源
	orm.RegisterDataBase("default", "mysql", dataSource)
}
