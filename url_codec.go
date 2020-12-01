package sqlca

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/gotools/log"
	"github.com/civet148/redigogo"
	"net/url"
	"strconv"
	"strings"
)

const (
	URL_SCHEME_SEP    = "://"
	URL_QUERY_SLAVE   = "slave"
	URL_QUERY_MAX     = "max"
	URL_QUERY_IDLE    = "idle"
	URL_QUERY_CHARSET = "charset"
)

const (
	//DSN no windows authentication: "Provider=SQLOLEDB;port=1433;server=127.0.0.1\SQLEXPRESS;database=test;user id=sa;password=123456"
	//DSN with windows authentication: "Provider=SQLOLEDB;integrated security=SSPI;port=1433;Data Source=127.0.0.1;database=mydb"
	WINDOWS_DSN_PROVIDER_SQLOLEDB        = "Provider=SQLOLEDB"
	WINDOWS_DSN_PORT                     = "Port"
	WINDOWS_DSN_DATA_SOURCE              = "Server"
	WINDOWS_DSN_INITIAL_CATALOG          = "Database"
	WINDOWS_DSN_USER_ID                  = "User Id"
	WINDOWS_DSN_PASSWORD                 = "Password"
	WINDOWS_DSN_INTEGRATED_SECURITY_SSPI = "Integrated Security=SSPI"
)

type UrlInfo struct {
	Scheme     string
	Host       string // host name and port like '127.0.0.1:3306'
	User       string
	Password   string
	Path       string
	Fragment   string
	Opaque     string
	ForceQuery bool
	Queries    map[string]string
}

type dsnDriver struct {
	strDriverName string
	parameter     dsnParameter
}

type dsnParameter struct {
	host     string //ip:port
	ip       string //ip
	port     string //port
	user     string
	password string
	db       string
	charset  string
	slave    bool
	max      int
	idle     int
	strDSN   string
	queries  map[string]string
}

func (d *dsnDriver) SetMax(max int) {
	if max > 0 {
		d.parameter.max = max
	}
}

func (d *dsnDriver) SetIdle(idle int) {
	if idle > 0 {
		d.parameter.idle = idle
	}
}

func (d *dsnDriver) SetSlave(slave bool) {
	if slave {
		d.parameter.slave = slave
	}
}

func (d *dsnParameter) parseUrlInfo(ui *UrlInfo) {
	var ok bool
	var val string

	d.user = ui.User
	d.host = ui.Host
	d.ip, d.port = getHostPort(d.host)
	d.password = ui.Password
	d.db = parseDatabaseName(ui.Path)
	d.charset = ui.Queries[URL_QUERY_CHARSET]
	d.queries = ui.Queries

	if val, ok = ui.Queries[URL_QUERY_SLAVE]; ok {
		if val == "true" {
			d.slave = true
		}
		delete(ui.Queries, URL_QUERY_SLAVE)
	}

	if val, ok = ui.Queries[URL_QUERY_MAX]; ok {
		if val != "" {
			d.max, _ = strconv.Atoi(val)
		} else {
			d.max = 100
		}
		delete(ui.Queries, URL_QUERY_MAX)
	}

	if val, ok = ui.Queries[URL_QUERY_IDLE]; ok {
		if val != "" {
			d.idle, _ = strconv.Atoi(val)
		} else {
			d.idle = 1
		}
		delete(ui.Queries, URL_QUERY_IDLE)
	}
}

//URL have some special characters in password(支持URL中密码包含特殊字符)
func ParseUrl(strUrl string) (ui *UrlInfo) {

	ui = &UrlInfo{Queries: make(map[string]string, 1)}

	var encodes = map[string]string{
		"`":  "%60",
		"#":  "%23",
		"?":  "%3f",
		"<":  "%3c",
		">":  "%3e",
		"[":  "%5b",
		"]":  "%5d",
		"{":  "%7b",
		"}":  "%7d",
		"/":  "%2f",
		"|":  "%7c",
		"\\": "%5c",
		"%":  "%25",
		"^":  "%5e",
	}

	var decodes = map[string]string{
		"%60": "`",
		"%23": "#",
		"%3f": "?",
		"%3c": "<",
		"%3e": ">",
		"%5b": "[",
		"%5d": "]",
		"%7b": "{",
		"%7d": "}",
		"%2f": "/",
		"%7c": "|",
		"%5c": "\\",
		"%25": "%",
		"%5e": "^",
	}
	_ = decodes

	// scheme://[userinfo@]host:port/path[?query][#fragment]

	strUrl = strings.TrimSpace(strUrl)
	if strings.Contains(strUrl, "@") { // if a url have user+password, there must be have '@'
		// find first '://'
		var strScheme string
		_ = strScheme

		index := strings.LastIndex(strUrl, URL_SCHEME_SEP)
		if index > 0 {
			strScheme = strUrl[:index]
			strUrl = strUrl[index+3:]
		}

		// find last '@'
		index = strings.LastIndex(strUrl, "@")
		if index > 0 {
			strPrefix := strUrl[:index]
			strSuffix := strUrl[index:]
			for k, v := range encodes {
				//encode user and password special character(s) to url encode
				strPrefix = strings.ReplaceAll(strPrefix, k, v)
			}

			if strScheme != "" {
				strUrl = strScheme + URL_SCHEME_SEP
			}
			strUrl += strPrefix + strSuffix
		}
	}

	u, err := url.Parse(strUrl)
	if err != nil {
		return
	}
	ui.Path = u.Path
	ui.Host = u.Host
	ui.Scheme = u.Scheme
	ui.Fragment = u.Fragment
	ui.Opaque = u.Opaque
	ui.ForceQuery = u.ForceQuery

	if u.User != nil {
		ui.User = u.User.Username()
		ui.Password, _ = u.User.Password()
		for k, v := range decodes {
			//decode password from url encode to special character(s)
			ui.Password = strings.ReplaceAll(ui.Password, k, v)
		}
	}
	vs, _ := url.ParseQuery(u.RawQuery)
	for k, v := range vs {
		ui.Queries[k] = v[0]
	}
	return
}

func parseDatabaseName(strPath string) string {
	idx := strings.LastIndex(strPath, "/")
	if idx == -1 {
		log.Errorf("[%v] invalid database path", strPath)
	}
	return strPath[idx+1:]
}

func getHostPort(strHost string) (ip, port string) {

	ipport := strings.Split(strHost, ":")
	assert(len(ipport) == 2, "invalid host:port string [%v]", strHost)
	ip = ipport[0]
	port = ipport[1]
	return
}

//DSN="root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4"
func (e *Engine) parseMysqlUrl(strUrl string) (parameter dsnParameter) {

	ui := ParseUrl(strUrl)
	parameter.parseUrlInfo(ui)
	e.setDatabaseName(parseDatabaseName(ui.Path))
	parameter.strDSN = fmt.Sprintf("%s:%s@tcp(%s)%s", ui.User, ui.Password, ui.Host, ui.Path)
	var queries []string
	for k, v := range ui.Queries {
		queries = append(queries, fmt.Sprintf("%v=%v", k, v))
	}
	if len(queries) > 0 {
		parameter.strDSN += fmt.Sprintf("?%s", strings.Join(queries, "&"))
	}
	return
}

//DSN="host=127.0.0.1 port=5432 user=root password=123456 dbname=mydb sslmode=disable"
func (e *Engine) parsePostgresUrl(strUrl string) (parameter dsnParameter) {

	ui := ParseUrl(strUrl)
	parameter.parseUrlInfo(ui)

	e.setDatabaseName(parseDatabaseName(ui.Path))
	strDatabase := e.getDatabaseName()
	strIP, strPort := getHostPort(ui.Host)

	var ok bool
	var strSSLMode string

	if strSSLMode, ok = ui.Queries["sslmode"]; !ok {
		strSSLMode = "disable"
	}
	parameter.strDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", strIP, strPort, ui.User, ui.Password, strDatabase, strSSLMode)
	return
}

//DSN: "/var/lib/my.db"
func (e *Engine) parseSqliteUrl(strUrl string) (parameter dsnParameter) {

	s := strings.Split(strUrl, URL_SCHEME_SEP)
	assert(len(s) == 2, "invalid url [%v] of sqlite, eg. 'sqlite:///var/lib/my.db'", strUrl)
	parameter.strDSN = s[1]
	return
}

//DSN no windows authentication: "Provider=SQLOLEDB;port=1433;server=127.0.0.1\SQLEXPRESS;database=test;user id=sa;password=123456"
//DSN with windows authentication: "Provider=SQLOLEDB;integrated security=SSPI;port=1433;Data Source=127.0.0.1;database=mydb"
func (e *Engine) parseMssqlUrl(strUrl string) (parameter dsnParameter) {

	var isWindowsAuth bool
	var dsnArgs []string

	ui := ParseUrl(strUrl)
	parameter.parseUrlInfo(ui)
	if strWindowsAuth, ok := ui.Queries["windows"]; ok {
		if strWindowsAuth == "true" {
			isWindowsAuth = true
		}
	}
	e.setDatabaseName(parseDatabaseName(ui.Path))

	dsnArgs = append(dsnArgs, "Provider=SQLOLEDB") //set driver provider
	if isWindowsAuth {                             //windows authentication
		dsnArgs = append(dsnArgs, "integrated security=SSPI") //set security mode
	}

	strIP, strPort := getHostPort(ui.Host)
	strDataSource := fmt.Sprintf("server=%s", strIP)           // set data source (host ip or domain)
	dsnArgs = append(dsnArgs, fmt.Sprintf("port=%s", strPort)) //set port to connect
	if strInst, ok := ui.Queries["instance"]; ok {
		if strInst != "" {
			strDataSource += "\\" + strInst //set instance name if not null
		}
	}
	dsnArgs = append(dsnArgs, strDataSource)
	dsnArgs = append(dsnArgs, fmt.Sprintf("database=%s", e.getDatabaseName())) //database name
	if !isWindowsAuth {
		dsnArgs = append(dsnArgs, fmt.Sprintf("user id=%s", ui.User))
		dsnArgs = append(dsnArgs, fmt.Sprintf("password=%s", ui.Password))
	}
	parameter.strDSN = strings.Join(dsnArgs, ";")
	return
}

//DSN: `{"password":"123456","db_index":0,"master_host":"127.0.0.1:6379","replicate_hosts":["127.0.0.1:6380","127.0.0.1:6381"]}`
func (e *Engine) parseRedisUrl(strUrl string) (parameter dsnParameter) {

	ui := ParseUrl(strUrl)
	cc := &redigogo.Config{
		Password:   ui.User, //redis have no user, just password
		MasterHost: fmt.Sprintf("%v", ui.Host),
	}

	if v, ok := ui.Queries[CACHE_DB_INDEX]; ok {
		cc.Index, _ = strconv.Atoi(v)
	}
	if v, ok := ui.Queries[CACHE_REPLICATE]; ok {
		cc.ReplicateHosts = strings.Split(v, ",")
	}

	if jsonData, err := json.Marshal(cc); err != nil {
		log.Errorf("url [%v] illegal", strUrl)
	} else {
		parameter.strDSN = string(jsonData)
	}
	return
}

//root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4
func (e *Engine) parseMysqlDSN(adapterType AdapterType, strMySQLDSN string) (dsn dsnDriver) {
	var strQueries string
	var querySlice []string
	var strUserPassword string
	dsn.parameter.queries = make(map[string]string, 0)
	strDatabaseName := trimBetween(strMySQLDSN, "/", "?")
	dsn.strDriverName = adapterType.DriverName()
	dsn.parameter.strDSN = strMySQLDSN
	dsn.parameter.db = strDatabaseName
	dsn.parameter.host = trimBetween(strMySQLDSN, "(", ")")
	dsn.parameter.ip, dsn.parameter.port = getHostPort(dsn.parameter.host)
	strQueries = cutLeft(strMySQLDSN, "?")
	strUserPassword = cutRight(strMySQLDSN, "@")
	ss := strings.Split(strUserPassword, ":")
	dsn.parameter.user = ss[0]
	dsn.parameter.password = ss[1]
	querySlice = strings.Split(strQueries, "&")
	for _, q := range querySlice {
		qv := strings.Split(q, "=")
		if len(qv) == 2 {
			if qv[0] == URL_QUERY_CHARSET {
				dsn.parameter.charset = qv[1]
			}
			dsn.parameter.queries[qv[0]] = qv[1]
		}
	}
	e.setDatabaseName(strDatabaseName)
	return
}

func cutFirst(strIn, strSep string) (strOut string) {

	return
}

func trimBetween(strIn, strLeftSep, strRightSep string) (strOut string) {
	strOut = cutLeft(strIn, strLeftSep)
	strOut = cutRight(strOut, strRightSep)
	return
}

func cutLeft(strIn, strSep string) (strOut string) {

	if strIn == "" || strSep == "" || len(strSep) != 1 {
		return strIn
	}

	nIdx := strings.LastIndex(strIn, strSep)
	if nIdx == -1 {
		strOut = strIn
	} else if nIdx == 0 {
		if len(strIn) > 1 {
			strOut = strIn[nIdx+1:]
		}
	} else {
		strOut = strIn[nIdx+1:]
	}
	return
}

func cutRight(strIn, strSep string) (strOut string) {
	if strIn == "" || strSep == "" || len(strSep) != 1 {
		return strIn
	}

	nIdx := strings.Index(strIn, strSep)
	if nIdx == -1 {
		strOut = strIn
	} else if nIdx > 0 {
		strOut = strIn[:nIdx]
	}

	return
}
