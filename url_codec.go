package sqlca

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	URL_SCHEME_SEP = "://"
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
func parseUrl(strUrl string) (ui *UrlInfo) {

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

func getDatabaseName(strPath string) string {
	idx := strings.LastIndex(strPath, "/")
	if idx == -1 {
		assert(false, "[%v] invalid database path", strPath)
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
func (e *Engine) parseMysqlUrl(strUrl string) (strDSN string) {

	ui := parseUrl(strUrl)
	strDSN = fmt.Sprintf("%s:%s@tcp(%s)%s", ui.User, ui.Password, ui.Host, ui.Path)
	var queries []string
	for k, v := range ui.Queries {
		queries = append(queries, fmt.Sprintf("%v=%v", k, v))
	}
	if len(queries) > 0 {
		strDSN += fmt.Sprintf("?%s", strings.Join(queries, "&"))
	}
	return
}

//DSN="host=127.0.0.1 port=5432 user=root password=123456 dbname=mydb sslmode=disable"
func (e *Engine) parsePostgresUrl(strUrl string) (strDSN string) {

	ui := parseUrl(strUrl)
	strDatabase := getDatabaseName(ui.Path)
	strIP, strPort := getHostPort(ui.Host)

	var ok bool
	var strSSLMode string
	if strSSLMode, ok = ui.Queries["sslmode"]; !ok {
		strSSLMode = "disable"
	}
	strDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", strIP, strPort, ui.User, ui.Password, strDatabase, strSSLMode)
	return
}

//DSN: "/var/lib/my.db"
func (e *Engine) parseSqliteUrl(strUrl string) (strDSN string) {

	s := strings.Split(strUrl, URL_SCHEME_SEP)
	assert(len(s) == 2, "invalid url [%v] of sqlite, eg. 'sqlite:///var/lib/my.db'", strUrl)
	strDSN = s[1]
	return
}

//DSN no windows authentication: "Provider=SQLOLEDB;port=1433;Data Source=127.0.0.1;Initial Catalog=mydb;user id=sa;password=123456"
//DSN with windows authentication: "Provider=SQLOLEDB;integrated security=SSPI;port=1433;Data Source=127.0.0.1;Initial Catalog=mydb;user id=sa;password=123456"
func (e *Engine) parseMssqlUrl(strUrl string) (strDSN string) {

	var isWindowsAuth bool
	var dsnArgs []string

	ui := parseUrl(strUrl)
	if strWindowsAuth, ok := ui.Queries["windows"]; ok {
		if strWindowsAuth == "true" {
			isWindowsAuth = true
		}
	}

	dsnArgs = append(dsnArgs, "Provider=SQLOLEDB") //set driver provider
	if isWindowsAuth {                             //windows authentication
		dsnArgs = append(dsnArgs, "integrated security=SSPI") //set security mode
	}

	strIP, strPort := getHostPort(ui.Host)
	strDataSource := fmt.Sprintf("Data Source=%s", strIP)      // set data source (host ip or domain)
	dsnArgs = append(dsnArgs, fmt.Sprintf("port=%s", strPort)) //set port to connect
	if strInst, ok := ui.Queries["instance"]; ok {
		if strInst != "" {
			strDataSource += "\\" + strInst //set instance name if not null
		}
	}
	dsnArgs = append(dsnArgs, strDataSource)
	dsnArgs = append(dsnArgs, fmt.Sprintf("Initial Catalog=%s", getDatabaseName(ui.Path))) //database name
	dsnArgs = append(dsnArgs, fmt.Sprintf("user id=%s", ui.User))
	dsnArgs = append(dsnArgs, fmt.Sprintf("password=%s", ui.Password))
	strDSN = strings.Join(dsnArgs, ";")
	return
}

func (e *Engine) parseRedisUrl(strUrl string) (strDSN string) {
	return
}

func (e *Engine) parseMemcacheUrl(strUrl string) (strDSN string) {
	return
}

func (e *Engine) parseMemoryUrl(strUrl string) (strDSN string) {
	return
}
