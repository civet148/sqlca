package sqlca

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/civet148/sqlca/v3/types"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MigrateAfterCB func(ctx context.Context, db *Engine)

// createDialector 根据DSN连接串自动创建并返回Dialector对象
func createDialector(adapterType types.AdapterType, dsn string) (gorm.Dialector, error) {
	switch adapterType {
	case types.AdapterSqlx_MySQL:
		return createMySQLDialector(dsn)
	case types.AdapterSqlx_Sqlite:
		return createSQLiteDialector(dsn)
	case types.AdapterSqlx_Postgres, types.AdapterSqlx_OpenGauss:
		return createPostgresDialector(dsn)
	case types.AdapterSqlx_Mssql:
		return createSQLServerDialector(dsn)
	}
	return nil, fmt.Errorf("cannot determine database type")
}

// createMySQLDialector 创建MySQL的Dialector
func createMySQLDialector(dsn string) (gorm.Dialector, error) {
	// 解析URL格式的DSN
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MySQL DSN: %w", err)
	}

	// 获取连接信息
	username := u.User.Username()
	password, _ := u.User.Password()
	host := u.Hostname()
	port := u.Port()
	database := strings.TrimPrefix(u.Path, "/")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3306"
	}
	if database == "" {
		return nil, fmt.Errorf("database name is required")
	}

	// 处理查询参数
	query := u.Query()
	params := make(map[string]string)

	// 转换常见参数
	if charset := query.Get("charset"); charset != "" {
		params["charset"] = charset
	}
	if parseTime := query.Get("parseTime"); parseTime != "" {
		params["parseTime"] = parseTime
	}
	if loc := query.Get("loc"); loc != "" {
		params["loc"] = loc
	}
	if timeout := query.Get("timeout"); timeout != "" {
		params["timeout"] = timeout
	}

	// 添加其他参数
	for key, values := range query {
		if key != "charset" && key != "parseTime" && key != "loc" && key != "timeout" {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	}

	// 构建标准的MySQL DSN
	var standardDSN string
	if password != "" {
		standardDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	} else {
		standardDSN = fmt.Sprintf("%s@tcp(%s:%s)/%s", username, host, port, database)
	}

	// 添加参数
	if len(params) > 0 {
		var queryParams []string
		for key, value := range params {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}
		standardDSN = standardDSN + "?" + strings.Join(queryParams, "&")
	}

	return mysql.Open(standardDSN), nil
}

// createPostgresDialector 创建PostgreSQL的Dialector
func createPostgresDialector(dsn string) (gorm.Dialector, error) {
	// PostgreSQL可以直接使用URL格式的DSN
	// 但需要确保格式正确
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL DSN: %w", err)
	}

	// 确保协议正确
	if u.Scheme == "postgresql" {
		dsn = strings.Replace(dsn, "postgresql://", "postgres://", 1)
	}

	return postgres.Open(dsn), nil
}

// createSQLServerDialector 创建SQL Server的Dialector
func createSQLServerDialector(dsn string) (gorm.Dialector, error) {
	// 解析SQL Server连接字符串
	// GORM的SQL Server驱动使用标准连接字符串
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL Server DSN: %w", err)
	}

	username := u.User.Username()
	password, _ := u.User.Password()
	host := u.Hostname()
	port := u.Port()
	database := strings.TrimPrefix(u.Path, "/")

	if port == "" {
		port = "1433"
	}

	// 构建SQL Server连接字符串
	// 格式: sqlserver://username:password@host:port?database=dbname
	var connStr string
	if password != "" {
		connStr = fmt.Sprintf("sqlserver://%s:%s@%s:%s", username, password, host, port)
	} else {
		connStr = fmt.Sprintf("sqlserver://%s@%s:%s", username, host, port)
	}

	if database != "" {
		connStr = connStr + "?database=" + database
	}

	// 添加其他查询参数
	query := u.Query()
	if len(query) > 0 {
		if strings.Contains(connStr, "?") {
			connStr = connStr + "&"
		} else {
			connStr = connStr + "?"
		}

		var params []string
		for key, values := range query {
			if key != "database" && len(values) > 0 {
				params = append(params, fmt.Sprintf("%s=%s", key, values[0]))
			}
		}
		connStr = connStr + strings.Join(params, "&")
	}

	return sqlserver.Open(connStr), nil
}

// createSQLiteDialector 创建SQLite的Dialector
func createSQLiteDialector(dsn string) (gorm.Dialector, error) {
	// SQLite DSN格式: sqlite:///path/to/database.db
	path := strings.TrimPrefix(dsn, "sqlite://")
	if path == "" {
		path = ":memory:" // 内存数据库
	}

	// SQLite使用文件路径，不需要复杂的解析
	return sqlite.Open(path), nil
}

// 辅助函数：检测是否为MySQL DSN格式
func isMySQLDSN(dsn string) bool {
	// 检查常见的MySQL DSN模式
	patterns := []string{
		"@tcp(",
		"@unix(",
		"/?",
	}

	for _, pattern := range patterns {
		if strings.Contains(dsn, pattern) {
			return true
		}
	}

	// 检查是否包含MySQL常见的参数
	mysqlParams := []string{"charset=", "parseTime=", "loc="}
	for _, param := range mysqlParams {
		if strings.Contains(dsn, param) {
			return true
		}
	}

	return false
}

// 辅助函数：检测是否为PostgreSQL DSN格式
func isPostgreSQLDSN(dsn string) bool {
	// PostgreSQL常见的参数
	pgParams := []string{"sslmode=", "search_path=", "TimeZone="}
	for _, param := range pgParams {
		if strings.Contains(dsn, param) {
			return true
		}
	}

	// 检查是否为host= port= dbname= 格式
	if strings.Contains(dsn, "host=") &&
		(strings.Contains(dsn, "dbname=") || strings.Contains(dsn, "database=")) {
		return true
	}

	return false
}
