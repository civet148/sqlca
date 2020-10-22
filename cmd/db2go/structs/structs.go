package structs

import (
	"fmt"
	"github.com/civet148/sqlca/cmd/db2go/cmder"
	"github.com/civet148/sqlca/cmd/db2go/schema"
	"strings"
)

const (
	STRUCT_PROMPT    = "STRUCT> "
	STRUCT_DELIMITER = '}'
	STRUCT_SHORT     = "s"
)

const (
	METHOD_ARGS_NULL     = ""
	METHOD_NAME_STRING   = "String"
	METHOD_NAME_GOSTRING = "GoString"
	METHOD_NAME_GET      = "Get"
	METHOD_NAME_SET      = "Set"
)

type GoMember struct {
	Name string
	Type string
	Tags string
}

type GoStruct struct {
	Name     string
	Members  []GoMember
	Contents []string
	cmd      *schema.Commander
}

func ExportStruct(cmd *schema.Commander) {
	var strInputs []string
	strInputs = cmder.Prompt(STRUCT_PROMPT, STRUCT_DELIMITER)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("")

	if gs, err := NewGoStruct(cmd, strInputs); err != nil {
		fmt.Printf("fmt.Sscanf error %s\n", err)
		return
	} else {
		if err = gs.parseMembers(); err != nil {
			fmt.Printf("parse struct members error %s\n", err)
			return
		}
		if err = gs.generateMethods(); err != nil {
			fmt.Printf("generate struct methods error %s\n", err)
			return
		}
	}
}

func NewGoStruct(cmd *schema.Commander, strInputs []string) (gs *GoStruct, err error) {

	if len(strInputs) == 0 {
		fmt.Printf("nil input\n")
		return nil, fmt.Errorf("nil input")
	}

	var strInput = strInputs[0]
	var strPrefix, strStructName, strSuffix string
	if _, err = fmt.Sscanf(strInput, "%s %s %s", &strPrefix, &strStructName, &strSuffix); err != nil {
		fmt.Printf("fmt.Sscanf error %s\n", err)
		return
	}
	if strStructName == "" {
		panic("struct syntax no illegal, go struct must begin with [type xxx struct {]")
	}
	contents := strInputs[1:]
	if len(contents) == 0 {
		panic("struct syntax no illegal, go struct have no any member")
	}
	gs = &GoStruct{
		Name:     strStructName,
		Contents: contents,
		cmd:      cmd,
	}
	return
}

func (g *GoStruct) parseMembers() (err error) {
	for _, v := range g.Contents {

		var strName, strType, strTags string
		fmt.Sscanf(v, "%s %s %s", &strName, &strType, &strTags)
		if strName == "" || strType == "" {
			continue
		}
		g.Members = append(g.Members, GoMember{
			Name: strings.TrimSpace(strName),
			Type: strings.TrimSpace(strType),
			Tags: strings.TrimSpace(strTags),
		})
	}
	return
}

func (g *GoStruct) generateMethods() (err error) {
	var strMethods string
	strMethods += g.generateString()
	strMethods += g.generateGoString()
	for _, v := range g.Members {
		strMethods += g.generateGetter(v)
		strMethods += g.generateSetter(v)
	}

	fmt.Println(strMethods)
	return
}

func (g *GoStruct) generateMethodDeclare(strMethodName, strArgs, strReturn, strLogic string) (strFunc string) {
	if strReturn == "" {
		strFunc = fmt.Sprintf("func (%s *%s) %s(%s) {\n", STRUCT_SHORT, g.Name, strMethodName, strArgs)
	} else {
		strFunc = fmt.Sprintf("func (%s *%s) %s(%s) %s {\n", STRUCT_SHORT, g.Name, strMethodName, strArgs, strReturn)
	}
	strFunc += strLogic
	strFunc += fmt.Sprintf("}\n\n")
	return
}

func (g *GoStruct) generateString() (strMethod string) {
	var strLogic string
	strLogic = fmt.Sprintf(`    data, _ := json.Marshal(%s)
    return string(data)
`, STRUCT_SHORT)
	return g.generateMethodDeclare(METHOD_NAME_STRING, METHOD_ARGS_NULL, "string", strLogic)
}

func (g *GoStruct) generateGoString() (strMethod string) {
	var strLogic string
	strLogic = fmt.Sprintf(`    return %s.String()
`, STRUCT_SHORT)
	return g.generateMethodDeclare(METHOD_NAME_GOSTRING, METHOD_ARGS_NULL, "string", strLogic)
}

func (g *GoStruct) generateGetter(m GoMember) (strMethod string) {

	var strLogic string
	strLogic = fmt.Sprintf(`    return %s.%v
`, STRUCT_SHORT, m.Name)

	var strMethodName = METHOD_NAME_GET + schema.CamelCaseConvert(m.Name)
	return g.generateMethodDeclare(strMethodName, METHOD_ARGS_NULL, m.Type, strLogic)
}

func (g *GoStruct) generateSetter(m GoMember) (strMethod string) {

	var strLogic string
	strLogic = fmt.Sprintf(`    %s.%v = v
`, STRUCT_SHORT, m.Name)
	var strArgs = fmt.Sprintf("v %s", m.Type)
	var strMethodName = METHOD_NAME_SET + schema.CamelCaseConvert(m.Name)
	return g.generateMethodDeclare(strMethodName, strArgs, METHOD_ARGS_NULL, strLogic)
}
