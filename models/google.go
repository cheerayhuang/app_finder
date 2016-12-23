package models

import (
	"app_finder/models/db"

	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

const (
	REG_TITLE          = `<div class="id-app-title".*?>(.*?)</div>`
	REG_ARTIST         = `<span itemprop="name">(.*?)</span>`
	REG_GENRE          = `<span itemprop="genre">(.*?)</span>`
	REG_CONTENT_RATING = `<div class="content" itemprop="contentRating">(.*?)</div>`
	REG_USER_RATING    = `<div class="score".*?>(.*?)</div>`
	REG_RATING_COUNT   = `<span class="reviews-num".*?>(.*?)</span>`
	REG_NOT_FOUND      = `<title>Not Found</title>`

	GOOGLE_PLAY_URL = `https://play.google.com/store/apps/details?id=`
)

func GoogleLookup(bundleId string) map[string]string {

	logs.Debug("bundleId:", bundleId)

	if _, ok := mysql.(*db.DBase); !ok {
		Init()
	}

	mysql.(*db.DBase).SetDefaultTable("google_play_app")
	fields := []string{
		"bundleId",
		"trackCensoredName",
		"trackViewUrl",
		"genre",
		"artistName",
		"trackContentRating",
		"averageUserRating",
		"userRatingCount",
		"id",
		"currency",
		"price",
	}
	mysql.(*db.DBase).SetDefaultFields(fields)

	isExist, err := mysql.Exist("bundleId", bundleId)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	if isExist {
		logs.Debug("The google app is already exist in mysql")
		return map[string]string{bundleId: "ok"}
	}

	// dont found in mysql, lockup via apple API and then insert to mysql.
	logs.Debug("find this app via crawling google play...")
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:1085")
	}
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	r, err := client.Get(GOOGLE_PLAY_URL + bundleId)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	if _CheckNotFound(body) {
		return map[string]string{bundleId: "not found"}
	}

	info := make([]interface{}, 0)
	info = append(info, bundleId)
	info = append(info, _GraspContent(REG_TITLE, body, "str"))
	info = append(info, GOOGLE_PLAY_URL+bundleId)
	info = append(info, _GraspContent(REG_GENRE, body, "str"))
	info = append(info, _GraspContent(REG_ARTIST, body, "str"))
	info = append(info, _GraspContent(REG_CONTENT_RATING, body, "str"))
	info = append(info, _GraspContent(REG_USER_RATING, body, "float"))
	info = append(info, _GraspContent(REG_RATING_COUNT, body, "bigint"))
	info = append(info, 0)
	info = append(info, "")
	info = append(info, 0.0)

	_, err = mysql.Insert(info...)
	if err != nil {
		return map[string]string{"bundleId": bundleId, "err": err.Error()}
	}

	return map[string]string{bundleId: "ok"}
}

func _GraspContent(reg string, body []byte, t string) interface{} {
	re := regexp.MustCompile(reg)
	res := re.FindSubmatch(body)
	if res == nil {
		logs.Error("can not match regexp: %q.", reg)
		return string("")
	}

	switch t {
	case "str":
		return res[1]

	case "bigint":
		numStr := strings.Replace(string(res[1]), ",", "", -1)
		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return 0
		}
		return num

	case "float":
		num, err := strconv.ParseFloat(string(res[1]), 32)
		if err != nil {
			return 0.0
		}
		return num
	}

	// can't go here.
	return string("")
}

func _CheckNotFound(body []byte) bool {
	re := regexp.MustCompile(REG_NOT_FOUND)
	return re.Match(body)
}
