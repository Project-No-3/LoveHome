package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
	_ "strings"
	"time"
)

type AreaController struct {
	beego.Controller
}

/*
	type AreaResp struct{

		Errno string      `json:"errno"`
		Reemsg string     `json:"errmsg"`
		Data interface{}  `json:"data"`


	}

	map[string]interface{}
*/

func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *AreaController) GetAreaInfo() {
	beego.Info("==========/api/v1.0/area get succ!!!=========")

	//返回给前端的map结构体
	resp := make(map[string]interface{})
	resp["errno"] = models.RECODE_OK
	resp["errmsg"] = models.RecodeText(models.RECODE_OK)

	defer this.RetData(resp)

	//1 从缓存中redis读数据
	cache_conn, err := cache.NewCache("redis", `{"key":"lovehome","conn":"127.0.0.1:6379","dbNum":"0"}`)
	if err != nil {
		beego.Info("cache redis conn err,err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//2 如果redis有 之前的json字符串数据那么直接返回给前段
	value := cache_conn.Get("xixi")
	if value != nil {
		beego.Info("cache get value = ", value)
		fmt.Printf("value = %s \n", value)
	}
	//3 如果redis没有之前的json字符串数据， 从mysql查

	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value != nil {
		beego.Info("====== get area_info from cache !!!=======")

		var area_info interface{}

		json.Unmarshal(areas_info_value.([]byte), &area_info)
		resp["data"] = area_info
		return
	}

	o := orm.NewOrm()

	//得到查到的areas数据
	var areas []models.Area //[{aid,aname},{aid,aname},{aid,aname}]

	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		//返回错误信息给前端
		resp["errno"] = models.RECODE_DBERR
		resp["errmsg"] = models.RecodeText(models.RECODE_DBERR)

		return
	}
	if num == 0 {
		resp["errno"] = models.RECODE_NODATA
		resp["errmsg"] = models.RecodeText(models.RECODE_NODATA)

		return
	}

	//succ
	resp["data"] = areas

	areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info --> redis fail err = ", err)
		resp["errno"] = models.RECODE_DBERR
		resp["errno"] = models.RecodeText(models.RECODE_DBERR)
		return
	}

	//将封装好的返回结构体map 发送给前段
	return

}
