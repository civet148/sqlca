package sqlca

import (
	"database/sql"
	"fmt"
	"github.com/civet148/log"
	"github.com/jmoiron/sqlx"
	"time"
)

type SqlxExecutor struct {
	db *sqlx.DB
	tx *sql.Tx
}

func newSqlxExecutor(strDriverName, strDSN string) (Executor, error) {
	var err error
	var db *sqlx.DB
	if db, err = sqlx.Open(strDriverName, strDSN); err != nil {
		log.Errorf("open driver name [%v] DSN [%v] error [%v]", strDriverName, strDSN, err.Error())
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Errorf("ping driver name [%v] DSN [%v] error [%v]", strDriverName, strDSN, err.Error())
		panic(err.Error())
		return nil, err
	}
	return &SqlxExecutor{
		db: db,
	}, nil
}

func (m *SqlxExecutor) Ping() (err error) {
	return nil
}

func (m *SqlxExecutor) SetMaxOpenConns(n int) {

}

func (m *SqlxExecutor) SetMaxIdleConns(n int) {

}

func (m *SqlxExecutor) SetConnMaxLifetime(d time.Duration) {

}

func (m *SqlxExecutor) SetConnMaxIdleTime(d time.Duration) {

}

func (m *SqlxExecutor) Exec(e *Engine, strSQL string) (rowsAffected, lastInsertId int64, err error) {
	return
}

func (m *SqlxExecutor) Query(e *Engine, strSQL string) (count int64, err error) {
	var rows *sql.Rows
	if rows, err = m.db.Query(strSQL); err != nil {
		log.Errorf("query [%v] error [%v]", strSQL, err.Error())
		return
	}
	defer rows.Close()
	return e.fetchRows(rows)
}

func (m *SqlxExecutor) QueryEx(e *Engine, strSQL string) (rowsAffected, total int64, err error) {
	var rowsQuery, rowsCount *sql.Rows
	if rowsQuery, err = m.db.Query(strSQL); err != nil {
		log.Errorf("query [%v] error [%v]", strSQL, err.Error())
		return
	}

	defer rowsQuery.Close()
	if rowsAffected, err = e.fetchRows(rowsQuery); err != nil {
		return
	}

	strCountSql := e.makeSqlxQueryCount()
	if rowsCount, err = m.db.Query(strCountSql); err != nil {
		log.Errorf("query [%v] error [%v]", strCountSql, err.Error())
		return
	}

	defer rowsCount.Close()
	for rowsCount.Next() {
		total++
	}
	return
}

func (m *SqlxExecutor) QueryRaw(e *Engine, strSQL string) (rowsAffected int64, err error) {
	var rows *sqlx.Rows
	log.Debugf("query [%v]", strSQL)
	if rows, err = m.db.Queryx(strSQL); err != nil {
		log.Errorf("query [%v] error [%v]", strSQL, err.Error())
		return
	}
	defer rows.Close()
	return e.fetchRows(rows.Rows)
}

func (m *SqlxExecutor) QueryMap(e *Engine, strSQL string) (rowsAffected int64, err error) {
	var rows *sqlx.Rows
	if rows, err = m.db.Queryx(strSQL); err != nil {
		log.Errorf("SQL [%v] query error [%v]", strSQL, err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		rowsAffected++
		fetcher, _ := e.getFetcher(rows.Rows)
		*e.model.(*[]map[string]string) = append(*e.model.(*[]map[string]string), fetcher.mapValues)
	}
	return
}

func (m *SqlxExecutor) Insert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	switch e.adapterType {
	case AdapterType_Mssql:
		{
			if e.isPkInteger() && e.isPkValueNil() {
				lastInsertId, err = m.mssqlQueryInsert(e, strSQL)
			}
		}
	case AdapterType_Postgres:
		{
			if e.isPkInteger() && e.isPkValueNil() {
				lastInsertId, err = m.postgresQueryInsert(e, strSQL)
			}
		}
	default:
		{
			var r sql.Result
			r, err = m.db.Exec(strSQL)
			if err != nil {
				log.Errorf("error %v model %+v", err, e.model)
				return
			}

			lastInsertId, _ = r.LastInsertId() //MSSQL Server not support last insert id
		}
	}
	return
}

func (m *SqlxExecutor) Update(e *Engine, strSQL string) (rowsAffected int64, err error) {
	var r sql.Result
	r, err = m.db.Exec(strSQL)
	if err != nil {
		log.Errorf("error %v model %+v", err, e.model)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		log.Errorf("get rows affected error [%v] query [%v] model [%+v]", err, strSQL, e.model)
		return
	}
	log.Debugf("RowsAffected [%v] query [%v]", rowsAffected, strSQL)
	return
}

func (m *SqlxExecutor) Upsert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	switch e.adapterType {
	case AdapterType_Mssql:
		{
			lastInsertId, err = m.mssqlUpsert(e, e.makeSqlxInsert())
		}
	case AdapterType_Postgres:
		{
			lastInsertId, err = m.postgresQueryUpsert(e, strSQL)
		}
	default:
		{
			var r sql.Result
			r, err = m.db.Exec(strSQL)
			if err != nil {
				log.Errorf("error %v model %+v", err, e.model)
				return
			}
			lastInsertId, err = r.LastInsertId()
			if err != nil {
				log.Errorf("get last insert id error %v model %+v", err, e.model)
				return
			}
		}
	}
	return
}

func (m *SqlxExecutor) Delete(e *Engine, strSQL string) (rowsAffected int64, err error) {
	var r sql.Result
	r, err = m.db.Exec(strSQL)
	if err != nil {
		log.Errorf("error %v model %+v", err, e.model)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		log.Errorf("get rows affected error [%v] query [%v] model [%+v]", err, strSQL, e.model)
		return
	}
	log.Debugf("RowsAffected [%v] query [%v]", rowsAffected, strSQL)
	return
}

func (m *SqlxExecutor) TxBegin() (tx Executor, err error) {
	m.tx, err = m.db.Begin()
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return &SqlxExecutor{
		db: m.db,
		tx: m.tx,
	}, nil
}

func (m *SqlxExecutor) TxGet(e *Engine, dest interface{}, strQuery string, args ...interface{}) (count int64, err error) {
	var rows *sql.Rows
	rows, err = m.tx.Query(strQuery)
	if err != nil {
		log.Errorf("TxGet sql [%v] args %v query error [%v] auto rollback [%v]", strQuery, args, err.Error(), e.bAutoRollback)
		e.autoRollback()
		return
	}
	e.setModel(dest)

	defer rows.Close()
	if count, err = e.fetchRows(rows); err != nil {
		log.Errorf("TxGet sql [%v] args %v fetch row error [%v] auto rollback [%v]", strQuery, args, err.Error(), e.bAutoRollback)
		e.autoRollback()
		return
	}
	return
}

func (m *SqlxExecutor) TxExec(e *Engine, strQuery string, args ...interface{}) (lastInsertId, rowsAffected int64, err error) {
	return
}

func (m *SqlxExecutor) TxRollback() (err error) {
	return
}

func (m *SqlxExecutor) TxCommit() (err error) {
	return
}

func (m *SqlxExecutor) postgresQueryInsert(e *Engine, strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	strSQL += fmt.Sprintf(" RETURNING \"%v\"", e.GetPkName())
	log.Debugf("[%v]", strSQL)
	if rows, err = m.db.Query(strSQL); err != nil {
		log.Errorf("tx.Query error [%v]", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&lastInsertId); err != nil {
			log.Warnf("rows.Scan warning [%v]", err.Error())
			return
		}
	}
	return
}

func (m *SqlxExecutor) postgresQueryUpsert(e *Engine, strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	log.Debugf("[%v]", strSQL)
	if rows, err = m.db.Query(strSQL); err != nil {
		log.Errorf("tx.Query error [%v]", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&lastInsertId); err != nil {
			log.Warnf("rows.Scan warning [%v]", err.Error())
			return
		}
	}
	return
}

func (m *SqlxExecutor) mssqlQueryInsert(e *Engine, strSQL string) (lastInsertId int64, err error) {
	var rows *sql.Rows
	strSQL += " SELECT SCOPE_IDENTITY() AS last_insert_id"
	log.Debugf("[%v]", strSQL)
	if rows, err = m.db.Query(strSQL); err != nil {
		log.Errorf("tx.Query error [%v]", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&lastInsertId); err != nil {
			log.Warnf("rows.Scan warning [%v]", err.Error())
			return
		}
	}
	return
}

func (m *SqlxExecutor) mssqlUpsert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	var db Executor
	var query = e.makeSqlxQueryPrimaryKey()
	if db, err = m.TxBegin(); err != nil {
		log.Errorf("TxBegin error [%v]", err.Error())
		return
	}
	var count int64
	if count, err = db.TxGet(e, &lastInsertId, query); err != nil {
		log.Errorf("TxGet [%v] error [%v]", query, err.Error())
		_ = db.TxRollback()
		return
	}
	if count == 0 {
		// INSERT INTO users(...) values(...)  SELECT SCOPE_IDENTITY() AS last_insert_id
		//if _, _, err = db.TxExec(strSQL); err != nil
		if lastInsertId, err = e.mssqlQueryInsert(strSQL); err != nil {
			log.Errorf("mssqlQueryInsert [%v] error [%v]", strSQL, err.Error())
			_ = db.TxRollback()
			return
		}
	} else {
		// UPDATE users SET xxx=yyy WHERE id=nnn
		strUpdates := fmt.Sprintf("%v %v %v %v %v %v=%v",
			DATABASE_KEY_NAME_UPDATE, e.getTableName(),
			DATABASE_KEY_NAME_SET, e.getOnConflictDo(),
			DATABASE_KEY_NAME_WHERE, e.GetPkName(), lastInsertId)
		log.Debugf("%v", strUpdates)
		if _, _, err = db.TxExec(e, strUpdates); err != nil {
			log.Errorf("TxExec [%v] error [%v]", strSQL, err.Error())
			_ = db.TxRollback()
			return
		}
	}

	if err = db.TxCommit(); err != nil {
		log.Errorf("TxCommit [%v] error [%v]", strSQL, err.Error())
		return
	}
	return
}