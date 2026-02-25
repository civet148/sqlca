package sqlca

func (e *Engine) addPreload(query string, args ...any) *Engine {
	if e.preloads == nil {
		e.preloads = make(map[string][]any)
	}
	e.preloads[query] = args
	return e
}

func (e *Engine) handlePreloads() (err error) {
	for query, args := range e.preloads {
		if err = e.execPreload(query, args...); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) execPreload(query string, args ...any) (err error) {
	return nil
}
