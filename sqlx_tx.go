package sqlca

import "github.com/civet148/gotools/log"

type sqlcaTx struct {
	strTxSQL string
	kvs      []*cacheKeyValue
}

func newTx(strTxSQL string, kvs ...*cacheKeyValue) *sqlcaTx {

	return &sqlcaTx{
		strTxSQL: strTxSQL, //tx sql to exec
		kvs:      kvs,      //redis cache update key-value slice
	}
}

func (e *Engine) txExec(args ...*sqlcaTx) error {

	assert(args, "tx args is nil")

	tx, err := e.db.Begin()
	if tx == nil || err != nil {
		log.Errorf("invoke Begin return tx is nil or error [%v]", err)
		return err
	}

	for _, v := range args {
		if _, err = tx.Exec(v.strTxSQL); err != nil {
			log.Errorf("invoke tx.Exec error [%v]", err.Error())
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		log.Errorf("invoke tx.Commit error [%v], rollback...", err.Error())
		return err
	}

	if e.isDebug() {
		for i, v := range args {
			log.Debugf("tx[%v] = '%v'", i, v)
		}
		log.Debugf("commit ok")
	}

	return nil
}
