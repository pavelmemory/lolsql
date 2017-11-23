package sql

import (
	"strings"
	"text/template"
)

var TemplateFunctions = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
	"Title":   strings.Title,
	"AsTypeName": func(FieldTypePackageAlias, FieldType string) string {
		if FieldTypePackageAlias == "" {
			return FieldType
		}
		return FieldTypePackageAlias + "." + FieldType
	},
}

var TypeSelect =
`func {{.TypeName}}(fields ...{{.TypeName}}Field) {{ToLower .TypeName}}Builder {
	return {{ToLower .TypeName}}Builder{}
}

func (b {{ToLower .TypeName}}Builder) Where(...sql.Condition) {{ToLower .TypeName}}Builder {
	return b
}

func (b {{ToLower .TypeName}}Builder) GroupBy(column ...OrderField) {{ToLower .TypeName}}Builder {
	return b
}

func (b {{ToLower .TypeName}}Builder) Having(...sql.Condition) {{ToLower .TypeName}}Builder {
	return b
}

func (b {{ToLower .TypeName}}Builder) OrderBy(...sql.SortOrder) {{ToLower .TypeName}}Builder {
	return b
}

func (b {{ToLower .TypeName}}Builder) Get() ([]test_model.Order, error) {
	return nil, nil
}

func (b {{ToLower .TypeName}}Builder) GetPtr() ([]*test_model.Order, error) {
	return nil, nil
}`

/*
type OrderField interface {
	sql.Field
}

type orderField struct {}

func (orderField) GetType() parser.TypeIdentity {
	return parser.TypeIdentity{Name:"Order", Package:"github.com/less-leg/test_model"}
}
*/
var FieldBaseDeclaration = `type {{.TypeName}}Field interface {
	sql.Field
}

type {{ToLower .TypeName}}Field struct {}

func ({{ToLower .TypeName}}Field) GetType() parser.TypeIdentity {
	return parser.TypeIdentity{Name:"{{.TypeName}}", Package:"{{.Package}}"}
}`

/*
type orderStartField struct {
	orderField
}

func (orderStartField) GetName() string {
	return "Start"
}

func Start() OrderStartField {
	return orderStartField{}
}
*/
var FieldDeclaration = `type {{ToLower .TypeName}}{{.FieldName}}Field struct {
	{{ToLower .TypeName}}Field
}

func ({{ToLower .TypeName}}{{.FieldName}}Field) GetName() string {
	return "{{.FieldName}}"
}

func {{.FieldName}}() {{.TypeName}}{{.FieldName}}Field {
	return {{ToLower .TypeName}}{{.FieldName}}Field{}
}`

/*
func MultiTimeTime(times ...time.Time) (vals []interface{}) {
	for _, v := range times {
		vals = append(vals, v)
	}
	return
}
*/
var SliceXTypeToSliceInterfaces = `func Multi{{Title .Selector}}{{Title .TypeName}}({{ToLower .TypeName}}s ...{{if .Selector}}{{.Selector}}.{{end}}{{.TypeName}}) (vals []interface{}) {
	for _, v := range {{ToLower .TypeName}}s {
		vals = append(vals, v)
	}
	return
}`

/*
func (t orderStartField) Equal(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Equal,
		Values:              []interface{}{v},
	}
}
...
*/
var FieldConditionDeclaration = `func (t {{ToLower .TypeName}}{{.FieldName}}Field) Equal(v {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Equal,
		Values:              []interface{}{v},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) NotEqual(v {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotEqual,
		Values:              []interface{}{v},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) Greater(v {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Greater,
		Values:              []interface{}{v},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) GreaterOrEqual(v {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.GreaterOrEqual,
		Values:              []interface{}{v},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) Lesser(v {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Lesser,
		Values:              []interface{}{v},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) LesserOrEqual(v {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.LesserOrEqual,
		Values:              []interface{}{v},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) In(v ...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.In,
		Values:              common.Multi{{Title .FieldSelector}}{{Title .FieldType}}(v...),
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) NotIn(v ...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotIn,
		Values:              common.Multi{{Title .FieldSelector}}{{Title .FieldType}}(v...),
	}
}

{{- if .FieldNullable}}
func (t {{ToLower .TypeName}}{{.FieldName}}Field) IsNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNull,
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) IsNotNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNotNull,
	}
}
{{end -}}

{{- if .FieldLikable}}
func (t {{ToLower .TypeName}}{{.FieldName}}Field) Like(vs ...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Conjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.Like,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) LikeOr(vs ...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Disjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.Like,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) NotLike(vs ...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Conjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.NotLike,
			Values:              []interface{}{v},
		})
	}
	return mc
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) NotLikeOr(vs ...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	mc := sql.MultiCondition{LogicalOperator: sql.Disjunction}
	for _, v := range vs {
		mc.Conditions = append(mc.Conditions, sql.SingleCondition{
			Type:                t.GetType(),
			Field:               t.GetName(),
			ComparatorOperation: sql.NotLike,
			Values:              []interface{}{v},
		})
	}
	return mc
}
{{end -}}

{{- if .FieldBetweenable}}
func (t {{ToLower .TypeName}}{{.FieldName}}Field) Between(v1 {{AsTypeName .FieldSelector .FieldType}}, v2 {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Between,
		Values:              []interface{}{v1, v2},
	}
}

func (t {{ToLower .TypeName}}{{.FieldName}}Field) NotBetween(v1 {{AsTypeName .FieldSelector .FieldType}}, v2 {{AsTypeName .FieldSelector .FieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotBetween,
		Values:              []interface{}{v1, v2},
	}
}{{end -}}`

/*
type OrderStartField interface {
	OrderField
	Equal(time.Time) sql.Condition
	NotEqual(time.Time) sql.Condition
	Greater(time.Time) sql.Condition
	...
}
*/
var FieldConditionInterfaceDeclaration = `{{- define "CommonFieldConditionOperations" -}}
	Equal({{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	NotEqual({{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	Greater({{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	GreaterOrEqual({{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	Lesser({{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	LesserOrEqual({{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	In(...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	NotIn(...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition
{{- end -}}

{{- define "BetweenFieldConditionOperations" -}}
	Between({{AsTypeName .FieldSelector .FieldType}}, {{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	NotBetween({{AsTypeName .FieldSelector .FieldType}}, {{AsTypeName .FieldSelector .FieldType}}) sql.Condition
{{- end -}}

{{- define "LikeFieldConditionOperations" -}}
	Like(...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	NotLike(...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	LikeOr(...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition
	NotLikeOr(...{{AsTypeName .FieldSelector .FieldType}}) sql.Condition
{{- end -}}

{{- define "NullableFieldConditionOperations" -}}
	IsNull() sql.Condition
	IsNotNull() sql.Condition
{{- end -}}

type {{.TypeName}}{{.FieldName}}Field interface {
	{{.TypeName}}Field
	{{template "CommonFieldConditionOperations" . }}
	{{if .FieldNullable}}{{template "NullableFieldConditionOperations" . }}{{end}}
	{{if .FieldLikable}}{{template "LikeFieldConditionOperations" . }}{{end}}
	{{if .FieldBetweenable}}{{template "BetweenFieldConditionOperations" . }}{{end}}
}`
