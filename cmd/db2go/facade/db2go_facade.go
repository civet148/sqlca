package facade

import (
	"github.com/civet148/sqlca"
	"github.com/civet148/sqlca/cmd/db2go/schema"
)

type Exporter interface {
	Export(cmd schema.Commander, e *sqlca.Engine) *schema.TableResult
}
