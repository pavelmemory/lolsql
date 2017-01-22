package generator

import (
	"text/template"
	"strings"
	"github.com/less-leg/utils"
)

var Package, _ = template.New("").Parse(`package {{.}}`)

var Imports, _ = template.New("").Parse(`
import (
{{.}}
)
// This code was auto-generated by LOLSQL code-generation tool.
// Please do not modify it manually. All changes will be deleted after regeneration.
`)

var Column_interface, _ = template.New("").Parse(`
type column interface {
	Column() string
}`)

var Lol_struct, _ = template.New("").Parse(`
type lol struct {
	selectColumns []column
	whereInited   bool
}

func (this *lol) Render() string {
	if (len(this.selectColumns) == 0) {
		return "select {{.Value2}} from {{.Value1}}"
	}
	cols := make([]string, 0, len(this.selectColumns))
	for _, selectColumn := range this.selectColumns {
		cols = append(cols, selectColumn.Column())
	}
	return "select " + strings.Join(cols, ", ") + " from {{.Value1}}"
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

var ColumnStub_struct, _ = template.New("").Funcs(
	template.FuncMap{
		"Title": strings.Title,
		"ToLower": strings.ToLower,
		"DotToUnderscore": utils.DotToUnderscore}).Parse(`
{{range .}}
type {{DotToUnderscore .Value1 | ToLower}}Stub struct { column }
var {{DotToUnderscore .Value1 | ToLower}}StubConst {{DotToUnderscore .Value1 | ToLower}}Stub
func {{DotToUnderscore .Value1 | Title}}() *{{DotToUnderscore .Value1 | ToLower}}Stub {return &{{DotToUnderscore .Value1 | ToLower}}StubConst}
func (*{{DotToUnderscore .Value1 | ToLower}}Stub) Column() string {return "{{.Value2}}"}
{{end}}
`)

var Conditions = map[string]*template.Template{
	"int": ConditionInt,
	"int8": ConditionInt,
	"int16": ConditionInt,
	"int32": ConditionInt,
	"int64": ConditionInt,

	"float32": ConditionFloat,
	"float64": ConditionFloat,

	"string": ConditionString,

	"time.Time": ConditionTime,
}

var ConditionInt, _ = template.New("").Funcs(template.FuncMap{"Title": strings.Title, "ToLower": strings.ToLower}).Parse(`
type {{ToLower .FieldToColumn.Value1}}{{Title .StructName}} struct {
	HasNext
	values    []{{.IsNullable}}{{.TypeName}}
	checkNull bool
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) NotNull() LolCondition {
	if len(this.values) > 0 {
		panic("Invalid usage: cannot check both equality and NULL check in one condition")
	}
	this.checkNull = true
	return this
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) render() string {
	{{if .IsNullable}}
	if this.values == nil || len(this.values) == 0 {
		if this.checkNull {
			return "{{.FieldToColumn.Value2}} is not null"
		} else {
			return "{{.FieldToColumn.Value2}} is null"
		}
	}
	{{end}}
	if (len(this.values) == 1) {
		return "{{.FieldToColumn.Value2}} = " + strconv.Itoa({{.IsNullable}}this.values[0])
	}
	vstr := make([]string, 0, len(this.values))
	for _, vptr := range this.values {
		{{if .IsNullable}}
		if vptr != nil {
			vstr = append(vstr, strconv.Itoa({{.IsNullable}}vptr))
		}
		{{else}}
			vstr = append(vstr, strconv.Itoa(vptr))
		{{end}}
	}
	return "{{.FieldToColumn.Value2}} in (" + strings.Join(vstr, ", ") + ")"
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}

{{if .IsNullable}}
func {{Title .FieldToColumn.Value1}}Is(value ...{{.IsNullable}}{{.TypeName}}) LolCondition {
	return &{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}{values:value}
}
{{else}}
func {{Title .FieldToColumn.Value1}}Is(v0 {{.TypeName}}, vnext ...{{.IsNullable}}{{.TypeName}}) LolCondition {
	return &{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}{values:append([]{{.TypeName}}{v0}, vnext...)}
}
{{end}}
`)

var ConditionFloat, _ = template.New("").Funcs(template.FuncMap{"Title": strings.Title, "ToLower": strings.ToLower}).Parse(`
type {{ToLower .FieldToColumn.Value1}}{{Title .StructName}} struct {
	HasNext
	values    []{{.IsNullable}}float64
	checkNull bool
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) NotNull() LolCondition {
	if len(this.values) > 0 {
		panic("Invalid usage: cannot check both equality and NULL check in one condition")
	}
	this.checkNull = true
	return this
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) render() string {
	{{if .IsNullable}}
	if this.values == nil || len(this.values) == 0 {
		if this.checkNull {
			return "{{ToLower .FieldToColumn.Value1}} is not null"
		} else {
			return "{{ToLower .FieldToColumn.Value1}} is null"
		}
	}
	{{end}}
	if (len(this.values) == 1) {
		return "{{ToLower .FieldToColumn.Value1}} = " + strconv.FormatFloat({{.IsNullable}}this.values[0], 'f', -1, 64)
	}
	vstr := make([]string, 0, len(this.values))
	for _, vptr := range this.values {
		{{if .IsNullable}}
		if vptr != nil {
			vstr = append(vstr, strconv.FormatFloat({{.IsNullable}}vptr, 'f', -1, 64))
		}
		{{else}}
		vstr = append(vstr, strconv.FormatFloat(vptr, 'f', -1, 64))
		{{end}}
	}
	return "{{ToLower .FieldToColumn.Value1}} in (" + strings.Join(vstr, ", ") + ")"
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}

func {{Title .FieldToColumn.Value1}}Is(value ...{{.IsNullable}}float64) LolCondition {
	return &{{ToLower .FieldToColumn.Value1}}{{Title .StructName}} {values:value}
}
`)

var ConditionString, _ = template.New("").Funcs(template.FuncMap{"Title": strings.Title, "ToLower": strings.ToLower}).Parse(`
type {{ToLower .FieldToColumn.Value1}}{{Title .StructName}} struct {
	HasNext
	values    []{{.IsNullable}}{{.TypeName}}
	checkNot bool
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) render() string {
	{{if .IsNullable}}
	if this.values == nil || len(this.values) == 0 {
		if this.checkNot {
			return "{{.FieldToColumn.Value2}} is not null"
		} else {
			return "{{.FieldToColumn.Value2}} is null"
		}
	}
	{{end}}
	if len(this.values) == 1 {
		if this.checkNot {
			return "{{.FieldToColumn.Value2}} <> '" + {{.IsNullable}}(this.values[0]) + "'"
		} else {
			return "{{.FieldToColumn.Value2}} = '" + {{.IsNullable}}(this.values[0]) + "'"
		}
	}
	vstr := make([]string, 0, len(this.values))
	for _, vptr := range this.values {
		{{if .IsNullable}}
		if vptr != nil {
			vstr = append(vstr, "'" + *vptr + "'")
		} else {
			panic("NULL can't be used as one of the arguments in the list")
		}
		{{else}}
		vstr = append(vstr, "'" + vptr + "'")
		{{end}}
	}
	if this.checkNot {
		return "{{.FieldToColumn.Value2}} not in (" + strings.Join(vstr, ", ") + ")"
	} else {
		return "{{.FieldToColumn.Value2}} in (" + strings.Join(vstr, ", ") + ")"
	}
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) Render() string {
	if this.Next() != nil {
		return this.render() + " " + this.Next().Render()
	}
	return this.render()
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) And(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionAnd{condition:cond})
	return this
}

func (this *{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}) Or(cond LolCondition) LolCondition {
	this.SetNext(&lolConditionOr{condition:cond})
	return this
}

// TODO: if column is not nullable, then field condition must has at least 1 parameter
func {{Title .FieldToColumn.Value1}}Is(values ...{{.IsNullable}}string) LolCondition {
	return &{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}{values:values}
}

// TODO: if column is not nullable, then field condition must has at least 1 parameter
func {{Title .FieldToColumn.Value1}}IsNot(values ...{{.IsNullable}}string) LolCondition {
	return &{{ToLower .FieldToColumn.Value1}}{{Title .StructName}}{values:values, checkNot: true}
}
`)

var ConditionTime, _ = template.New("").Funcs(template.FuncMap{"Title": strings.Title, "ToLower": strings.ToLower}).Parse(`
panic("Time to fuck you asshole!")
`)
