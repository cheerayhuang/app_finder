package models

import (
	"app_finder/models/db"

	"github.com/astaxie/beego/logs"
)

var mysqlNotfound db.DB

func InitNotFound() {
	var err error
	mysqlNotfound, err = db.New("mysql",
		"sdkbox",
		"1234",
		"sdkbox",
		"localhost",
		"3306",
		nil,
	)
	if err != nil {
		logs.Error("create mysql object failed: ", err.Error())
	}
}

func Notfound(bundleId string) map[string]string {
	if bundleId != "" {
		logs.Debug("Add bundleId: %s to Notfound table", bundleId)
	} else {
		logs.Debug("Get Notfound list")
	}

	if _, ok := mysqlNotfound.(*db.DBase); !ok {
		InitNotFound()
	}

	mysqlNotfound.(*db.DBase).SetDefaultFields([]string{"bundleId"})
	mysqlNotfound.(*db.DBase).SetDefaultTable("not_found_app")

	switch bundleId {
	case "":
		logs.Debug("Get Notfound list")
		res, err := _ReplyNotfoundList()
		if err != nil {
			logs.Error("Get Notfound list falied: %s", err.Error())
			return map[string]string{"err": "Get Notfound list failed"}
		}
		return res

	default:
		logs.Debug("Add bundleId: %s to Notfound table", bundleId)
		err := _StoreBundleId(bundleId)
		if err != nil {
			logs.Error("Store notfound bundleId failed: %s", err.Error())
			return map[string]string{"bundleId": bundleId, "err": "Store notfound bundleId failed"}
		}
		return map[string]string{bundleId: "store to Notfound table"}
	}
}

func NotfoundDelete(bundleId string) map[string]string {
	if _, ok := mysqlNotfound.(*db.DBase); !ok {
		InitNotFound()
	}

	mysqlNotfound.(*db.DBase).SetDefaultTable("not_found_app")

	err := mysqlNotfound.Delete("bundleId", bundleId)
	if err != nil {
		logs.Error("Delete bundleId from Notfound table failed: %s", err.Error())
		return map[string]string{"bundleId": bundleId, "err": "delete from Notfound table failed"}
	}

	return map[string]string{bundleId: "delete from Notfound table"}
}

func _ReplyNotfoundList() (map[string]string, error) {
	rows, err := mysqlNotfound.Query(nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var bundleId string
		err = rows.Scan(&bundleId)
		if err != nil {
			return nil, err
		}
		result[bundleId] = "1"
	}

	return result, nil
}

func _StoreBundleId(bundleId string) error {
	isExist, err := mysqlNotfound.Exist("bundleId", bundleId)
	if err != nil {
		return err
	}
	if isExist {
		logs.Debug("bundleId has been already in Notfound table")
		return nil
	}

	info := []interface{}{bundleId}
	_, err = mysqlNotfound.Insert(info...)
	if err != nil {
		return err
	}

	return nil
}
