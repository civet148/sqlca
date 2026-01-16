package sqlca

import (
	"context"
	"fmt"

	"github.com/civet148/sqlca/v3/types"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type MigrateAfterCB func(ctx context.Context, db *Engine)

// createDialector 根据DSN连接串自动创建并返回Dialector对象
func createDialector(adapterType types.AdapterType, dsn dsnParameter) (gorm.Dialector, error) {
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
func createMySQLDialector(dsn dsnParameter) (gorm.Dialector, error) {
	return mysql.Open(dsn.DSN), nil
}

// createPostgresDialector 创建PostgreSQL的Dialector
func createPostgresDialector(dsn dsnParameter) (gorm.Dialector, error) {
	return postgres.Open(dsn.DSN), nil
}

// createSQLServerDialector 创建SQL Server的Dialector
func createSQLServerDialector(dsn dsnParameter) (gorm.Dialector, error) {
	return sqlserver.Open(dsn.DSN), nil
}

// createSQLiteDialector 创建SQLite的Dialector
func createSQLiteDialector(dsn dsnParameter) (gorm.Dialector, error) {
	return sqlite.Open(dsn.DSN), nil
}
