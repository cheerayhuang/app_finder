package controllers

import (
	"app_finder/models"
	"github.com/astaxie/beego"
)

type NotfoundController struct {
	beego.Controller
}

// @router /:id(.*) [post]
func (this *NotfoundController) Notfound() {
	bundleId := this.Ctx.Input.Param(":id")

	this.Data["json"] = models.Notfound(bundleId)
	this.ServeJSON()
}
