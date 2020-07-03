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
	URL_SCHEME_SEP  = "://"
	URL_QUERY_SLAVE = "slave"
	URL_QUERY_MAX   = "max"
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

func parseSlaveQueries(ui *UrlInfo) (ok bool) {
	var val string
	if val, ok = ui.Queries[URL_QUERY_SLAVE]; ok {
		if val != "true" {
			ok = false
		}
		delete(ui.Queries, URL_QUERY_SLAVE)
	}
	return
}

func parseMaxQueries(ui *UrlInfo) (max int) {

	var ok bool
	var val string
	if val, ok = ui.Queries[URL_QUERY_MAX]; ok {
		if val != "" {
			max, _ = strconv.Atoi(val)
		} else {
			max = 100
		}
		delete(ui.Queries, URL_QUERY_MAX)
	}
	return
}

//DSN="root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4"
func (e *Engine) parseMysqlUrl(strUrl string) (parameter dsnParameter) {

	ui := ParseUrl(strUrl)
	e.setDatabaseName(parseDatabaseName(ui.Path))
	parameter.strDSN = fmt.Sprintf("%s:%s@tcp(%s)%s", ui.User, ui.Password, ui.Host, ui.Path)
	var queries []string

	parameter.slave = parseSlaveQueries(ui)
	parameter.maxConnections = parseMaxQueries(ui)
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
	e.setDatabaseName(parseDatabaseName(ui.Path))
	strDatabase := e.getDatabaseName()
	strIP, strPort := getHostPort(ui.Host)

	var ok bool
	var strSSLMode string

	parameter.slave = parseSlaveQueries(ui)
	parameter.maxConnections = parseMaxQueries(ui)
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

//DSN no windows authentication: "Provider=SQLOLEDB;port=1433;Data Source=127.0.0.1;Initial Catalog=mydb;user id=sa;password=123456"
//DSN with windows authentication: "Provider=SQLOLEDB;integrated security=SSPI;port=1433;Data Source=127.0.0.1;Initial Catalog=mydb;user id=sa;password=123456"
func (e *Engine) parseMssqlUrl(strUrl string) (parameter dsnParameter) {

	var isWindowsAuth bool
	var dsnArgs []string

	ui := ParseUrl(strUrl)
	if strWindowsAuth, ok := ui.Queries["windows"]; ok {
		if strWindowsAuth == "true" {
			isWindowsAuth = true
		}
	}
	e.setDatabaseName(parseDatabaseName(ui.Path))

	//dsnArgs = append(dsnArgs, "Provider=SQLOLEDB") //set driver provider
	//if isWindowsAuth {                             //windows authentication
	//	dsnArgs = append(dsnArgs, "integrated security=SSPI") //set security mode
	//}
	//
	//strIP, strPort := getHostPort(ui.Host)
	//strDataSource := fmt.Sprintf("Data Source=%s", strIP)      // set data source (host ip or domain)
	//dsnArgs = append(dsnArgs, fmt.Sprintf("port=%s", strPort)) //set port to connect
	//if strInst, ok := ui.Queries["instance"]; ok {
	//	if strInst != "" {
	//		strDataSource += "\\" + strInst //set instance name if not null
	//	}
	//}
	//dsnArgs = append(dsnArgs, strDataSource)
	//dsnArgs = append(dsnArgs, fmt.Sprintf("Initial Catalog=%s", e.getDatabaseName())) //database name
	//dsnArgs = append(dsnArgs, fmt.Sprintf("user id=%s", ui.User))
	//dsnArgs = append(dsnArgs, fmt.Sprintf("password=%s", ui.Password))
	//strDSN = strings.Join(dsnArgs, ";")

	//dsnArgs = append(dsnArgs, "Provider=SQLOLEDB") //set driver provider
	if isWindowsAuth { //windows authentication
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
	dsnArgs = append(dsnArgs, fmt.Sprintf("user id=%s", ui.User))
	dsnArgs = append(dsnArgs, fmt.Sprintf("password=%s", ui.Password))
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
	e.strDatabaseName = trimBetween(strMySQLDSN, "/", "?")
	dsn.strDriverName = adapterType.DriverName()
	dsn.parameter.strDSN = strMySQLDSN
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
