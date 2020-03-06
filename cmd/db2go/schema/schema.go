package schema

type SchemaInfo struct {
	ConnUrl     string
	Databases   []string
	Tables      []string
	Scheme      string
	Host        string
	User        string
	Password    string
	Charset     string
	OutDir      string
	Prefix      string
	Suffix      string
	PackageName string
	Tags        string
}
