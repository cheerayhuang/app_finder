package models

import (
	"app_finder/models/db"

	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
)

var mysql db.DB

func Init() {
	var err error
	mysql, err = db.New("postgres",
		"awsuser",
		"678D9ed29db6",
		"events",
		"sdkbox.crzhtaa6feqd.us-west-2.redshift.amazonaws.com",
		"5439",
		nil,
	)
	if err != nil {
		logs.Error("create mysql object failed: ", err.Error())
	}
}

const (
	SEARCH_URL = "https://itunes.apple.com/search"
	LOOKUP_URL = "https://itunes.apple.com/lookup"
)

func AppleLookup(bundleId string) map[string]string {

	logs.Debug("bundleId:", bundleId)

	if _, ok := mysql.(*db.DBase); !ok {
		Init()
	}

	mysql.(*db.DBase).SetDefaultTable("apple_store_app")
	fields := []string{"id",
		"bundleId",
		"trackCensoredName",
		"trackViewUrl",
		"genre1",
		"genre2",
		"genre3",
		"genre4",
		"currency",
		"price",
		"artistId",
		"artistName",
		"sellerName",
		"trackContentRating",
		"averageUserRating",
		"userRatingCount",
		//"blob",
	}
	mysql.(*db.DBase).SetDefaultFields(fields)

	isExist, err := mysql.Exist("bundleId", bundleId)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	if isExist {
		id, err := mysql.QueryID("bundleId", bundleId)
		if err != nil {
			return map[string]string{"bundleId": bundleId, "err": err.Error()}
		}
		logs.Debug("The apple app is already exist in mysql")
		return map[string]string{bundleId: strconv.FormatInt(id, 10), "not_access_apple_api": "1"}
	}

	// dont found in mysql, lockup via apple API and then insert to mysql.
	logs.Debug("find this app via apple api...")
	r, err := http.Get(LOOKUP_URL + "?bundleId=" + bundleId + "&limit=1")
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	if js.Get("resultCount").MustInt() == 0 {
		return map[string]string{bundleId: "not found"}
	}

	c, err := _GetPropertFromJson(js.Get("results").GetIndex(0))
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}
	// field "blob"
	//c = append(c, body)

	_, err = mysql.Insert(c...)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	return map[string]string{bundleId: strconv.FormatInt(c[0].(int64), 10)}
}

func _GetPropertFromJson(js *simplejson.Json) ([]interface{}, error) {
	r := make([]interface{}, 0)

	r = append(r, js.Get("trackId").MustInt64())
	r = append(r, js.Get("bundleId").MustString())
	r = append(r, js.Get("trackCensoredName").MustString())
	r = append(r, js.Get("trackViewUrl").MustString())

	g := []string{"", "", "", ""}
	genres := js.Get("genres").MustArray()
	for k, v := range genres {
		g[k] = v.(string)
	}
	for _, v := range g {
		r = append(r, v)
	}

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
