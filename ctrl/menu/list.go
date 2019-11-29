package menu

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-admin/conf"
	"go-admin/models"
	"go-admin/modules/response"
	"net/http"
)
type Role struct {
	Key string `form:"key" json:"key"`
	Name string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	Routes []interface{} `form:"routes" json:"routes"`
}
func List(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get(conf.Cfg.Token)
	if v == nil {
		response.ShowError(c, "fail")
		return
	}
	uid := session.Get(v)
	user := models.SystemUser{Id: uid.(int)}
	has := user.GetRow()
	if !has {
		response.ShowError(c, "fail")
		return
	}

	menu := models.SystemMenu{}
	if user.Nickname == "admin" {
		menuArr, err := menu.GetAll()
		if err != nil {
			response.ShowError(c, "fail")
			return
		}
		jsonArr :=tree(menuArr)
		response.ShowData(c,jsonArr)
		return
	} else {
		menuArr:=menu.GetRouteByUid(uid)
		jsonArr :=tree(menuArr)
		response.ShowData(c,jsonArr)
		return
	}
}
func tree(menuArr []models.SystemMenu) ([]interface{}) {
	role := models.SystemRole{}
	mrArr := role.GetRowMenu()
	var menuMap = make(map[int][]models.SystemMenu, 0)
	for _, value := range menuArr {
		menuMap[value.Pid] = append(menuMap[value.Pid], value)
	}
	var jsonArr []interface{}

	mainMenu, ok := menuMap[0]
	if !ok {
		return nil
	}
	for _, value := range mainMenu {
		var item = make(map[string]interface{})
		item["path"] = value.Path
		item["component"] = value.Component
		if value.Redirect != "" {
			item["redirect"] = value.Redirect
		}
		if value.Alwaysshow ==1 {
			item["alwaysShow"] = true
		}
		if value.Hidden == 1 {
			item["hidden"] = true
		}
		var meta=make(map[string]interface{})
		_,ok:=mrArr[value.Id]
		if ok {
			meta["roles"]=mrArr[value.Id]
		}
		if value.MetaTitle!=""{
			meta["title"]=value.MetaTitle
		}
		if value.MetaIcon!="" {
			meta["icon"]=value.MetaIcon
		}
		if value.MetaAffix==1 {
			meta["affix"] = true
		}
		if value.MetaNocache==1 {
			meta["noCache"] = true
		}
		if len(meta)>0 {
			item["meta"]=meta
		}
		if _,ok:=menuMap[value.Id] ;ok{
			item["children"]=treeChilden(menuMap[value.Id],mrArr)
		}
		jsonArr = append(jsonArr,item)
	}
	return jsonArr

}
func treeChilden(menuArr []models.SystemMenu, mrArr map[int][]string)[]interface{} {
	var jsonArr []interface{}
	for _,value:=range menuArr  {
		var item = make(map[string]interface{})
		item["path"] = value.Path
		item["component"] = value.Component
		if value.Redirect != "" {
			item["redirect"] = value.Redirect
		}
		if value.Alwaysshow ==1 {
			item["alwaysShow"] = true
		}
		if value.Hidden == 1 {
			item["hidden"] = true
		}
		var meta=make(map[string]interface{})
		_,ok:=mrArr[value.Id]
		if ok {
			meta["roles"]=mrArr[value.Id]
		}
		if value.MetaTitle!=""{
			meta["title"]=value.MetaTitle
		}
		if value.MetaIcon!="" {
			meta["icon"]=value.MetaIcon
		}
		if value.MetaAffix==1 {
			meta["affix"] = true
		}
		if value.MetaNocache==1 {
			meta["noCache"] = true
		}
		if len(meta)>0 {
			item["meta"]=meta
		}
		jsonArr = append(jsonArr,item)
	}
	return jsonArr
}
func Roles(c *gin.Context){
	model:=models.SystemRole{}
	menu:=models.SystemMenu{}
	roleArr :=model.GetAll()
	var roleMenu []Role
	for _,value:=range roleArr {
		r:=Role{}
		r.Key=value.Name
		r.Name=value.Name
		r.Description=value.Description
		menuArr:=menu.GetRouteByRole(value.Id)
		r.Routes=tree(menuArr)
		roleMenu = append(roleMenu,r)
	}
	response.ShowData(c,roleMenu)
	return
}

func Dashboard(c *gin.Context){
	roleMenu:="{\"menuList\":[{\"create_time\":\"2018-03-1611:33:00\",\"menu_type\":\"M\",\"children\":[{\"create_time\":\"2018-03-1611:33:00\",\"menu_type\":\"C\",\"children\":[],\"parent_id\":1,\"menu_name\":\"用户管理\",\"icon\":\"#\",\"perms\":\"system:user:index\",\"order_num\":1,\"menu_id\":4,\"url\":\"/system/user\"},{\"create_time\":\"2018-12-2810:36:20\",\"menu_type\":\"M\",\"children\":[{\"create_time\":\"2018-12-2810:50:28\",\"menu_type\":\"C\",\"parent_id\":73,\"menu_name\":\"人员通讯录\",\"icon\":null,\"perms\":\"system:person:index\",\"order_num\":1,\"menu_id\":74,\"url\":\"/system/book/person\"}],\"parent_id\":1,\"menu_name\":\"通讯录管理\",\"icon\":\"fafa-address-book-o\",\"perms\":null,\"order_num\":1,\"menu_id\":73,\"url\":\"#\"}],\"parent_id\":0,\"menu_name\":\"系统管理\",\"icon\":\"fafa-adjust\",\"perms\":null,\"order_num\":2,\"menu_id\":1,\"url\":\"#\"},{\"create_time\":\"2018-03-1611:33:00\",\"menu_type\":\"M\",\"children\":[{\"create_time\":\"2018-03-1611:33:00\",\"menu_type\":\"C\",\"parent_id\":2,\"menu_name\":\"数据监控\",\"icon\":\"#\",\"perms\":\"monitor:data:view\",\"order_num\":3,\"menu_id\":15,\"url\":\"/system/druid/monitor\"}],\"parent_id\":0,\"menu_name\":\"系统监控\",\"icon\":\"fafa-video-camera\",\"perms\":null,\"order_num\":5,\"menu_id\":2,\"url\":\"#\"}],\"user\":{\"login_name\":\"admin\",\"user_id\":1,\"user_name\":\"管理员\",\"dept_id\":1}}"
	var data map[string]interface{}
	err := json.Unmarshal([]byte(roleMenu), &data)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data":  data,
	})
	return
}