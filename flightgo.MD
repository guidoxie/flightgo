# 数据库表结构明细

### aircraft：飞机型号数据
* type 飞机型号简称
* friendly_type 飞机型号全称



### airlines： 航空公司数据
* iata 国际航空运输协会指定的航空公司代码
* icao 国际民航组织指定的航空公司代码
* full_name 航空公司全称
* short_name 航空公司简称
* type 航空公司类型， out: 国外航空公司， in：国内航空公司
* logo 航空公司logo链接
* url 航空公司官网链接

### airpots： 机场数据
* iata 国际航空运输协会指定的机场代码
* icao 国际民航组织指定的机场代码
* lat 机场所在纬度
* lon 机场所在经度
* short_name 机场简称
* country 机场所属国家代码

### flight_datas: 航班数据
* ident 航班号简称
* date 航班起飞日期
* s_take_off 计划起飞时间
* e_take_off 预计起飞时间
* a_take_off 实际起飞时间
* s_landing 计划降落时间
* e_landing 预计降落时间
* a_landing 实际降落时间
* state 航班状态
* aircraft 飞机机型
* airline 所属航空公司代码
* origin 起飞地机场代码
* destination 降落地机场代码
* distance 飞行距离
* friendly_ident 航班号全称


### flight_number: 航班号数据
* ident 航班号简称
* iataIdent 国际航空运输协会指定的航班号

### plane_img： 飞机图片数据
* aircraft 飞机信号
* url 飞机图片链接

### state：航班状态类型
* type 状态类型
* name 状态对应的中文名称