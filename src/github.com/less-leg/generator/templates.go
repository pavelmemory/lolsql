package generator

import (
	"text/template"
	"strings"
)

var Package, _ = template.New("").Parse(`package {{.}}`)

var Imports, _ = template.New("").Parse(`
import ({{.}})
// This code was auto-generated by LOLSQL code-generation tool.
// Please do not modify it manually. All changes will be deleted after regeneration.
`)

var Column_interface, _ = template.New("").Parse(`type column interface {Column() string}`)

var Lol_struct, _ = template.New("").Parse(`
type lol struct {
	selectColumns []column
	whereInited   bool
}

func (this *lol) Render() string {
	if (len(this.selectColumns) == 0) {
		return "select {{.Columns}} from {{.TableName}}"
	}
	cols := make([]string, 0, len(this.selectColumns))
	for _, selectColumn := range this.selectColumns {
		cols = append(cols, selectColumn.Column())
	}
	return "select " + strings.Join(cols, ", ") + " from {{.TableName}}"
}

func (this *lol) Where(cond LolCondition) *lolWhere {
	if this.whereInited {
		panic("Invalid usage of WHERE statement: double usage not supported.")
	}
	this.whereInited = true
	return &lolWhere{retrieval:this, condition:cond}
}
`)

var Select_func, _ = template.New("").Parse(`
func Select(selects ...column) *lol {
	return &lol{selectColumns:selects}
}
`)

var LolWhere_struct, _ = template.New("").Parse(`
type lolWhere struct {
	retrieval *lol
	condition LolCondition
	next      []LolCondition
}
func (this *lolWhere) Render() string {
	if len(this.next) > 0 {
		conds := make([]string, 0, len(this.next))
		for _, cond := range this.next {
			conds = append(conds, cond.Render())
		}
		return this.retrieval.Render() + " where (" + this.condition.Render() + ") " + strings.Join(conds, " ")
	}
	return this.retrieval.Render() + " where " + this.condition.Render()
}
func (this *lolWhere) And(cond LolCondition) *lolWhere {
	if this.next == nil {
		this.next = make([]LolCondition, 0, 1)
	}
	this.next = append(this.next, &lolConditionAnd{condition: cond})
	return this
}
func (this *lolWhere) Or(cond LolCondition) *lolWhere {
	if this.next == nil {
		this.next = make([]LolCondition, 0, 1)
	}
	this.next = append(this.next, &lolConditionOr{condition: cond})
	return this
}`)


var LolConditionAnd_struct, _ = template.New("").Parse(`
type lolConditionAnd struct {
	HasNext
	condition LolCondition
}

func (this *lolConditionAnd) render() string {
	return "and (" + this.condition.Render() + ")"
}

func (this *lolConditionAnd) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *lolConditionAnd) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}

func (this *lolConditionAnd) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}`)

var LolConditionOr_struct, _ = template.New("").Parse(`
type lolConditionOr struct {
	HasNext
	condition LolCondition
}
func (this *lolConditionOr) render() string {
	return "or (" + this.condition.Render() + ")"
}
func (this *lolConditionOr) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}
func (this *lolConditionOr) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}
func (this *lolConditionOr) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}`)

var ColumnStub, _ = template.New("").Funcs(template.FuncMap{"Title": strings.Title, "ToLower": strings.ToLower}).Parse(`
{{ range . }}
type {{ ToLower .FieldName }}Stub struct { column }
var {{ ToLower .FieldName }}StubConst {{ ToLower .FieldName }}Stub
func {{ Title .FieldName }}() *{{ ToLower .FieldName }}Stub {return &{{ ToLower .FieldName }}StubConst}
func (*{{ ToLower .FieldName }}Stub) Column() string {return "{{ .ColumnName }}"}
{{ end }}
`)
