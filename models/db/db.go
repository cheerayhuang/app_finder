package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"shiftred/error"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

type DB interface {
	Query([]string, int64, int64) (*sql.Rows, error)
	QueryID(string, string) (int64, error)
	Insert(...interface{}) (int64, error)
	Delete()
	Update()
	Count() (int64, error)
	Exist(string, string) (bool, error)
	Close() error
}

type DBase struct {
	d        string
	user     string
	pwd      string
	database string
	host     string
	port     string
	table    string
	fields   []string
	db       *sql.DB
	tm_start int64
	tm_end   int64
}

func New(db, user, pwd, database, host, port string, now *time.Time) (DB, error) {
	if db == "" || user == "" || pwd == "" || database == "" || host == "" || port == "" {
		return nil, MyErr.New(MyErr.DB_CONN_MISS_PARAMS, "miss Database Connection Paramters")
	}

	r := new(DBase)
	r.user = user
	r.pwd = pwd
	r.database = database
	r.host = host
	r.port = port
	r.d = db

	var err error
	if db == "postgres" {
		conn_str := []string{"user=" + r.user, "password=" + r.pwd, "dbname=" + r.database, "host=" + r.host, "port=" + r.port, "sslmode=disable"}
		r.db, err = sql.Open(db, strings.Join(conn_str, " "))
		if err != nil {
			return nil, err
		}
	} else {
		conn_str := r.user + ":" + r.pwd + "@tcp(" + r.host + ":" + r.port + ")/" + database
		r.db, err = sql.Open(db, conn_str)
		if err != nil {
			return nil, err
		}
	}

	r.table = "event_log"
	r.fields = []string{"sdkbox_app_package_id", "sdkbox_platform", "sdkbox_country_code", "sdkbox_sdkboxversion", "sdkbox_app", "sdkbox_event", "timestamp"}

	/*if r.d == "postgres" {
		err = r._CalcDuration(now)
		if err != nil {
			return nil, err
		}
	}*/

	return r, nil
}

func (this *DBase) SetDefaultTable(name string) {
	if name == "" {
		return
	}

	this.table = name
}

func (this *DBase) SetDefaultFields(fields []string) {
	if fields == nil || len(fields) == 0 {
		return
	}

	this.fields = make([]string, len(fields))
	copy(this.fields, fields)
}

func (this *DBase) Query(fields []string, limit, offset int64) (*sql.Rows, error) {
	if fields == nil || len(fields) == 0 {
		if this.fields != nil && len(this.fields) != 0 {
			fields = this.fields
		} else {
			return nil, MyErr.New(MyErr.DB_QUERY_MISS_PARAMS, "query fields are empty")
		}
	}

	stat, err := this.QueryStat(fields, limit, offset)
	if err != nil {
		return nil, err
	}

	rows, err := this.db.Query(stat)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (this *DBase) QueryID(where_f, where_v string) (int64, error) {
	var id int64
	s := "SELECT id FROM " + this.table + " WHERE " + where_f + " = '" + where_v + "'"
	logs.Debug(s)
	err := this.db.QueryRow(s).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (this *DBase) Count() (int64, error) {
	var count int64
	err := this.db.QueryRow("SELECT COUNT(*) FROM " + this.table).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (this *DBase) Exist(field, value string) (bool, error) {
	var count int
	s := "SELECT COUNT(*) FROM " + this.table + " WHERE " + field + " = '" + value + "'"
	logs.Debug(s)
	err := this.db.QueryRow(s).Scan(&count)
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	}

	return false, nil
}

func (this *DBase) Delete() {

}

func (this *DBase) Update() {

}

func (this *DBase) Insert(args ...interface{}) (int64, error) {
	l := len(args)
	if l == 0 {
		return 0, MyErr.New(MyErr.DB_INSERT_MISS_VALUES, "miss values in insert statement.")
	}

	var marks []string
	for i := 1; i <= l; i++ {
		if this.d == "postgres" {
			marks = append(marks, "$"+strconv.Itoa(i))
		} else {
			marks = append(marks, "?")
		}
	}
	marks_str := strings.Join(marks, ",")
	stmt_str := "INSERT INTO " + this.table + " VALUES(" + marks_str + ")"
	//logs.Debug(stmt_str)
	stmt, err := this.db.Prepare(stmt_str)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	row_count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return row_count, nil
}

func (this *DBase) Close() error {
	return this.db.Close()
}

func (this *DBase) QueryStat(fields []string, limit, offset int64) (string, error) {
	if len(fields) == 0 {
		return "", MyErr.New(MyErr.DB_QUERY_MISS_PARAMS, "query fields are empty")
	}

	var fields_str []string
	for _, f := range fields {
		fields_str = append(fields_str, f)
	}

	select_stat := "SELECT " + strings.Join(fields_str, ",") + " FROM " + this.table
	where_stat := ""
	limit_stat := ""
	if this.d == "postgres" {
		//where_stat = " WHERE timestamp >= " + strconv.FormatInt(this.tm_start, 10) + " AND timestamp <= " + strconv.FormatInt(this.tm_end, 10)
		limit_stat = " LIMIT " + strconv.FormatInt(limit, 10) + " OFFSET " + strconv.FormatInt(offset, 10)
	}

	return select_stat + where_stat + limit_stat, nil
}

func (this *DBase) _CalcDuration(now *time.Time) error {
	if now == nil {
		return MyErr.New(MyErr.DB_QUERY_MISS_PARAMS, "can't get timestamp")
	}

	hour := time.Duration(now.Hour())
	min := time.Duration(now.Minute())
	sec := time.Duration(now.Second() + 1)

	end := now.Add(-(24*5*time.Hour + hour*time.Hour + min*time.Minute + sec*time.Second))
	start := end.Add(-(24 * time.Hour))

	this.tm_start = start.Unix()
	this.tm_end = end.Unix()

	return nil
}
