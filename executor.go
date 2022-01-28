package sqlca

import "time"

type executor interface {
	Ping() (err error)
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	Exec(e *Engine, strSQL string)  (rowsAffected, lastInsertId int64, err error)
	Query(e *Engine, strSQL string) (count int64, err error)
	QueryEx(e *Engine, strSQL string) (rowsAffected, total int64, err error)
	QueryRaw(e *Engine, strSQL string) (rowsAffected int64, err error)
	QueryMap(e *Engine, strSQL string) (rowsAffected int64, err error)
	Update(e *Engine, strSQL string) (rowsAffected int64, err error)
	Insert(e *Engine, strSQL string) (lastInsertId int64, err error)
	Upsert(e *Engine, strSQL string) (lastInsertId int64, err error)
	Delete(e *Engine, strSQL string) (rowsAffected int64, err error)
	//tx methods
	txBegin() (executor, error)
	txGet(e *Engine, dest interface{}, strQuery string, args ...interface{}) (count int64, err error)
	txExec(e *Engine, strQuery string, args ...interface{}) (lastInsertId, rowsAffected int64, err error)
	txRollback() error
	txCommit() error
}
