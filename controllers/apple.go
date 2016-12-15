package controllers

import (
	"app_finder/models"
	"encoding/json"

	//"net/http"

	"github.com/astaxie/beego"
)

type AppleController struct {
	beego.Controller
}

// @router / [get]
func (this *AppleController) Search() {
	var params models.SearchParams
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)

	this.Data["json"] = models.AppleSearch(params)
	this.ServeJSON()
}

// @router /:id(.+) [get]
func (this *AppleController) Lookup() {
	bundleId := this.Ctx.Input.Param(":id")

	this.Data["json"] = models.AppleLookup(bundleId)
	this.ServeJSON()
}
