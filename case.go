package sqlca

import "fmt"

type when struct {
	strWhen string
	strThen string
}
type CaseWhen struct {
	e       *Engine
	whens   []*when
	strElse string
	strEnd  string
}

func (c *CaseWhen) Case(strThen string, strWhen string, args ...interface{}) *CaseWhen {
	c.whens = append(c.whens, &when{
		strThen: strThen,
		strWhen: c.e.formatString(strWhen, args...),
	})
	return c
}

func (c *CaseWhen) Else(strElse string) *CaseWhen {
	c.strElse = strElse
	return c
}

func (c *CaseWhen) End(strName string) *Engine {
	var e *Engine

	e = c.e
	c.strEnd = strName

	e.strCaseWhen = DATABASE_KEY_NAME_CASE
	for _, v := range c.whens {
		e.strCaseWhen += fmt.Sprintf(" %s %s %s %s ", DATABASE_KEY_NAME_WHEN, v.strWhen, DATABASE_KEY_NAME_THEN, e.getQuoteColumnValue(v.strThen))
	}
	if c.strElse != "" {
		e.strCaseWhen += fmt.Sprintf(" %s %s ", DATABASE_KEY_NAME_ELSE, e.getQuoteColumnValue(c.strElse))
	}
	e.strCaseWhen += fmt.Sprintf(" %s %s ", DATABASE_KEY_NAME_END, c.strEnd)
	return c.e
}
