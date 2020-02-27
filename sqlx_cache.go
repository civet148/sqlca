package sqlca

// clone engine
func (e *Engine) clone(model interface{}) *Engine {

	return &Engine{
		db:        e.db,
		cache:     e.cache,
		debug:     e.debug,
		model:     model,
		strPkName: e.strPkName,
	}
}

func (e *Engine) checkModel() bool {

	if e.model == nil {
		e.panic("orm model is nil, please call Model() method before query or update")
		return false
	}
	return true
}

func (e *Engine) getTableName() string {
	return e.strTableName
}

func (e *Engine) setTableName(strName string) {
	e.strTableName = strName
}

func (e *Engine) getPkName() string {
	return e.strPkName
}

func (e *Engine) setPkName(strName string) {
	e.strPkName = strName
}

func (e *Engine) getWhere() string {
	return e.strWhere
}

func (e *Engine) setWhere(strWhere string) {
	e.strWhere = strWhere
}

func (e *Engine) makeMySQL() (strSQL string) {

	//TODO: @libin make SQL query string (mysql)
	if e.debug {
		e.debugf(strSQL)
	}
	return
}

func (e *Engine) makePostgreSQL() (strSQL string) {

	//TODO: @libin make SQL query string (postgresql)
	if e.debug {
		e.debugf(strSQL)
	}
	return
}
