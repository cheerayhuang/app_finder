package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["app_finder/controllers:AppleController"] = append(beego.GlobalControllerRouter["app_finder/controllers:AppleController"],
		beego.ControllerComments{
			Method: "Lookup",
			Router: `/:id(.+)`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:GoogleController"] = append(beego.GlobalControllerRouter["app_finder/controllers:GoogleController"],
		beego.ControllerComments{
			Method: "Lookup",
			Router: `/:id(.+)`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
