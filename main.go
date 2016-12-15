package main

import (
	_ "app_finder/routers"

	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/logs"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	/*logs.SetLogger(logs.AdapterFile, `{"filename":"/tmp/app_finder.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":7}`)
	 */

	beego.Run()
}
