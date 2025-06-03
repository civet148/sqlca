package sqlca

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v3/types"
	"net/url"
	"strconv"
	"strings"
)

const (
	urlSchemeSep       = "://"
	urlQueryMax        = "max"
	urlQueryIdle       = "idle"
	urlQueryCharset    = "charset"
	urlQuerySchema     = "schema"
	urlQuerySearchPath = "search_path"
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

func (ui *UrlInfo) Url() string {
	var queries []string
	for k, v := range ui.Queries {
		queries = append(queries, fmt.Sprintf("%s=%s", k, v))
	}
	strQueries := strings.Join(queries, "&")
	return fmt.Sprintf("%s://%s:%s@%s%s?%s", ui.Scheme, ui.User, ui.Password, ui.Host, ui.Path, strQueries)
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

func (d *dsnParameter) parseUrlInfo(ui *UrlInfo) {
	var ok bool
	var val string

	d.user = ui.User
	d.host = ui.Host
	d.ip, d.port = getHostPort(d.host)
	d.password = ui.Password
	d.db = parseDatabaseName(ui.Path)
	d.charset = ui.Queries[urlQueryCharset]
	d.queries = ui.Queries

	if val, ok = ui.Queries[urlQueryMax]; ok {
		if val != "" {
			d.max, _ = strconv.Atoi(val)
		} else {
			d.max = 100
		}
		delete(ui.Queries, urlQueryMax)
	}

	if val, ok = ui.Queries[urlQueryIdle]; ok {
		if val != "" {
			d.idle, _ = strconv.Atoi(val)
		} else {
			d.idle = 1
		}
		delete(ui.Queries, urlQueryIdle)
	}
}

// URL have some special characters in password(支持URL中密码包含特殊字符)
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

		index := strings.LastIndex(strUrl, urlSchemeSep)
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
				strUrl = strScheme + urlSchemeSep
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

// DSN="root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4"
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

// DSN="host=127.0.0.1 port=5432 user=root password=123456 dbname=mydb sslmode=disable"
func (e *Engine) parsePostgresUrl(strUrl string) (parameter dsnParameter) {
	ui := ParseUrl(strUrl)
	parameter.parseUrlInfo(ui)
	e.setDatabaseName(parseDatabaseName(ui.Path))
	strDatabase := e.getDatabaseName()
	strIP, strPort := getHostPort(ui.Host)
	parameter.strDSN = buildPostgresDSN(strIP, strPort, ui.User, ui.Password, strDatabase, ui.Queries)
	return
}

func buildPostgresDSN(strIP, strPort, strUser, strPassword, strDatabase string, queries map[string]string) string {
	var kvs []string
	var strExtras string
	for k, v := range queries {
		if k == urlQuerySchema {
			k = urlQuerySearchPath
		}
		kvs = append(kvs, fmt.Sprintf("%s=%s", k, v))
	}
	if len(kvs) > 0 {
		strExtras = strings.Join(kvs, " ")
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", strIP, strPort, strUser, strPassword, strDatabase, strExtras)
}

// DSN: "/var/lib/my.db"
func (e *Engine) parseSqliteUrl(strUrl string) (parameter dsnParameter) {

	s := strings.Split(strUrl, urlSchemeSep)
	assert(len(s) == 2, "invalid url [%v] of sqlite, eg. 'sqlite:///var/lib/my.db'", strUrl)
	parameter.strDSN = s[1]
	return
}

// DSN no windows authentication: "Provider=SQLOLEDB;port=1433;server=127.0.0.1\SQLEXPRESS;database=test;user id=sa;password=123456"
// DSN with windows authentication: "Provider=SQLOLEDB;integrated security=SSPI;port=1433;Data Source=127.0.0.1;database=mydb"
func (e *Engine) parseMssqlUrl(strUrl string) (parameter dsnParameter) {
	ui := ParseUrl(strUrl)
	parameter.parseUrlInfo(ui)
	e.setDatabaseName(parseDatabaseName(ui.Path))
	parameter.strDSN = buildMssqlDSN(parameter.ip, parameter.port, parameter.user, parameter.password, parameter.db, parameter.queries)
	return
}

func buildMssqlDSN(strIP, strPort, strUser, strPassword, strDatabase string, queries map[string]string) (strDSN string) {
	var isWindowsAuth bool
	var dsnArgs []string
	dsnArgs = append(dsnArgs, "Provider=SQLOLEDB") //set driver provider
	if isWindowsAuth {                             //windows authentication
		dsnArgs = append(dsnArgs, "integrated security=SSPI") //set security mode
	}

	strDataSource := fmt.Sprintf("server=%s", strIP)           // set data source (host ip or domain)
	dsnArgs = append(dsnArgs, fmt.Sprintf("port=%s", strPort)) //set port to connect
	if strInst, ok := queries["instance"]; ok {
		if strInst != "" {
			strDataSource += "\\" + strInst //set instance name if not null
		}
	}
	dsnArgs = append(dsnArgs, strDataSource)
	dsnArgs = append(dsnArgs, fmt.Sprintf("database=%s", strDatabase)) //database name
	if !isWindowsAuth {
		dsnArgs = append(dsnArgs, fmt.Sprintf("user id=%s", strUser))
		dsnArgs = append(dsnArgs, fmt.Sprintf("password=%s", strPassword))
	}
	return strings.Join(dsnArgs, ";")
}

//root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4
func (e *Engine) parseMysqlDSN(adapterType types.AdapterType, strMySQLDSN string) (dsn dsnDriver) {
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
			if qv[0] == urlQueryCharset {
				dsn.parameter.charset = qv[1]
			}
			dsn.parameter.queries[qv[0]] = qv[1]
		}
	}
	e.setDatabaseName(strDatabaseName)
	return
}

// rawMySql2Url convert raw mysql data source name to url, e.g "root:123456@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4"
func rawMySql2Url(strRawDSN string) string {
	strRawDSN = strings.ReplaceAll(strRawDSN, "tcp(", "")
	strRawDSN = strings.ReplaceAll(strRawDSN, ")", "")
	return fmt.Sprintf("%s://%s", types.DRIVER_NAME_MYSQL, strRawDSN)
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

//root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4
func Url2MySql(strUrl string) (string, error) {
	if strings.Index(strUrl, urlSchemeSep) == -1 {
		return strUrl, nil
	}
	ui := ParseUrl(strUrl)
	var params []string
	strUser := ui.User
	strPasswd := ui.Password
	strHost := ui.Host
	strPath := ui.Path
	for k, v := range ui.Queries {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}
	return fmt.Sprintf("%s:%s@tcp(%s)%s?%s", strUser, strPasswd, strHost, strPath, strings.Join(params, "&")), nil
}
