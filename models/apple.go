package models

import (
	"app_finder/models/db"

	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
)

func init() {

}

const (
	SEARCH_URL = "https://itunes.apple.com/search"
	LOOKUP_URL = "https://itunes.apple.com/lookup"
)

type SearchParams struct {
	Term    string `json:"term"`
	Country string `json:"country"`
	Limit   int    `json:"limit"`
}

type SearchReply struct {
}

func AppleSearch(params SearchParams) (r map[string]string) {

	return map[string]string{"ok": "tbd"}

}

func AppleLookup(bundleId string) map[string]string {

	mysql, err := db.New("mysql",
		"sdkbox",
		"1234",
		"apple_google_apps",
		"localhost",
		"3306",
		nil,
	)
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}

	mysql.(*db.DBase).SetDefaultTable("apple_store_app")
	fields := []string{"id",
		"bundleId",
		"trackCensoredName",
		"trackViewUrl",
		"description",
		"currency",
		"price",
		"artistId",
		"artistName",
		"sellerName",
		"trackContentRating",
		"averageUserRating",
		"userRatingCount",
		"blob",
	}
	mysql.(*db.DBase).SetDefaultFields(fields)

	isExist, err := mysql.Exist("bundleId", bundleId)
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}

	if isExist {
		id, err := mysql.QueryID("bundleId", bundleId)
		if err != nil {
			return map[string]string{"bundleId": "0", "err": err.Error()}
		}
		logs.Debug("The apple app is already exist in mysql")
		return map[string]string{bundleId: strconv.FormatInt(id, 10)}
	}

	// dont found in mysql, lockup via apple API and then insert to mysql.
	r, err := http.Get(LOOKUP_URL + "?bundleId=" + bundleId + "&limit=1")
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}

	if js.Get("resultCount").MustInt() == 0 {
		return map[string]string{"bundleId": "0", "err": "not found"}
	}

	c, err := _GetPropertFromJson(js.Get("results").GetIndex(0))
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}
	c = append(c, body)

	_, err = mysql.Insert(c...)
	if err != nil {
		return map[string]string{"bundleId": "0", "err": err.Error()}
	}

	return map[string]string{bundleId: strconv.FormatInt(c[0].(int64), 10)}
}

func _GetPropertFromJson(js *simplejson.Json) ([]interface{}, error) {
	r := make([]interface{}, 0)

	r = append(r, js.Get("trackId").MustInt64())
	r = append(r, js.Get("bundleId").MustString())
	r = append(r, js.Get("trackCensoredName").MustString())
	r = append(r, js.Get("trackViewUrl").MustString())
	r = append(r, js.Get("description").MustString())
	r = append(r, js.Get("currency").MustString())
	r = append(r, js.Get("price").MustFloat64())
	r = append(r, js.Get("artistId").MustInt64())
	r = append(r, js.Get("artistName").MustString())
	r = append(r, js.Get("sellerName").MustString())
	r = append(r, js.Get("trackContentRating").MustString())
	r = append(r, js.Get("averageUserRating").MustFloat64())
	r = append(r, js.Get("userRatingCount").MustInt64())

	return r, nil
}
