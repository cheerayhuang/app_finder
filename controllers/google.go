package controllers

import (
	"app_finder/models"
	"github.com/astaxie/beego"
)

type GoogleController struct {
	beego.Controller
}

// @router /:id(.+) [get]
func (this *GoogleController) Lookup() {
	bundleId := this.Ctx.Input.Param(":id")

	this.Data["json"] = models.GoogleLookup(bundleId)
	this.ServeJSON()
}
