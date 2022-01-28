package sqlca

import "time"

type MgoExecutor struct {

}

func newMgoExecutor(strUrl string, options ...interface{}) (executor, error) {
	return &MgoExecutor{
	}, nil
}

func (m *MgoExecutor) Ping() (err error) {
	return nil
}


func (m *MgoExecutor) SetMaxOpenConns(n int) {

}

func (m *MgoExecutor) SetMaxIdleConns(n int) {

}

func (m *MgoExecutor) SetConnMaxLifetime(d time.Duration) {

}

func (m *MgoExecutor) SetConnMaxIdleTime(d time.Duration) {

}

func (m *MgoExecutor) Exec(e *Engine, strSQL string) (rowsAffected, lastInsertId int64, err error) {
	return
}

func (m *MgoExecutor) Query(e *Engine, strSQL string) (count int64, err error){
	return
}

func (m *MgoExecutor) QueryEx(e *Engine, strSQL string) (rowsAffected, total int64, err error) {
	return
}

func (m *MgoExecutor) QueryRaw(e *Engine, strSQL string) (rowsAffected int64, err error) {
	return
}

func (m *MgoExecutor) QueryMap(e *Engine, strSQL string) (rowsAffected int64, err error) {

	return
}

func (m *MgoExecutor) Update(e *Engine, strSQL string) (rowsAffected int64, err error) {

	return
}

func (m *MgoExecutor) Insert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	return
}

func (m *MgoExecutor) Upsert(e *Engine, strSQL string) (lastInsertId int64, err error) {

	return
}

func (m *MgoExecutor) Delete(e *Engine, strSQL string) (rowsAffected int64, err error) {


	return
}

func (m *MgoExecutor) txBegin() (tx executor, err error) {
	return
}

func (m *MgoExecutor) txGet(e *Engine, dest interface{}, strQuery string, args ...interface{}) (count int64, err error) {
	return
}

func (m *MgoExecutor) txExec(e *Engine, strQuery string, args ...interface{}) (lastInsertId, rowsAffected int64, err error) {
	return
}

func (m *MgoExecutor) txRollback() (err error) {
	return
}

func (m *MgoExecutor) txCommit() (err error) {
	return
}
