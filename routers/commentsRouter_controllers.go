package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["app_finder/controllers:AppleController"] = append(beego.GlobalControllerRouter["app_finder/controllers:AppleController"],
		beego.ControllerComments{
			Method: "Search",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:AppleController"] = append(beego.GlobalControllerRouter["app_finder/controllers:AppleController"],
		beego.ControllerComments{
			Method: "Lookup",
			Router: `/:id(.+)`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:ObjectController"] = append(beego.GlobalControllerRouter["app_finder/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:ObjectController"] = append(beego.GlobalControllerRouter["app_finder/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:ObjectController"] = append(beego.GlobalControllerRouter["app_finder/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:ObjectController"] = append(beego.GlobalControllerRouter["app_finder/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:ObjectController"] = append(beego.GlobalControllerRouter["app_finder/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["app_finder/controllers:UserController"] = append(beego.GlobalControllerRouter["app_finder/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
