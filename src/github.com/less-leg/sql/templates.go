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

var TypeSelect = `{{- $builderType := print (ToLower .TypeName) "Builder" -}}
func {{.TypeName}}(fields ...{{.TypeName}}Field) {{$builderType}} {
	return {{$builderType}}{}
}

func (b {{$builderType}}) Where(...sql.Condition) {{$builderType}} {
	return b
}

func (b {{$builderType}}) GroupBy(column ...OrderField) {{$builderType}} {
	return b
}

func (b {{$builderType}}) Having(...sql.Condition) {{$builderType}} {
	return b
}

func (b {{$builderType}}) OrderBy(...sql.SortOrder) {{$builderType}} {
	return b
}

func (b {{$builderType}}) Get() ([]{{AsTypeName .TypeSelector .TypeName}}, error) {
	return nil, nil
}

func (b {{$builderType}}) GetPtr() ([]*{{AsTypeName .TypeSelector .TypeName}}, error) {
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
var SliceXTypeToSliceInterfaces = `func Multi{{Title .Selector}}{{Title .TypeName}}({{ToLower .TypeName}}s ...{{AsTypeName .Selector .TypeName}}) (vals []interface{}) {
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
var FieldConditionDeclaration = `{{- $typeFieldName := print (ToLower .TypeName) .FieldName "Field" -}}
{{- $typeFieldType := AsTypeName .FieldSelector .FieldType -}}
func (t {{$typeFieldName}}) Equal(v {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Equal,
		Values:              []interface{}{v},
	}
}

func (t {{$typeFieldName}}) NotEqual(v {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotEqual,
		Values:              []interface{}{v},
	}
}

func (t {{$typeFieldName}}) Greater(v {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Greater,
		Values:              []interface{}{v},
	}
}

func (t {{$typeFieldName}}) GreaterOrEqual(v {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.GreaterOrEqual,
		Values:              []interface{}{v},
	}
}

func (t {{$typeFieldName}}) Lesser(v {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Lesser,
		Values:              []interface{}{v},
	}
}

func (t {{$typeFieldName}}) LesserOrEqual(v {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.LesserOrEqual,
		Values:              []interface{}{v},
	}
}

func (t {{$typeFieldName}}) In(v ...{{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.In,
		Values:              common.Multi{{Title .FieldSelector}}{{Title .FieldType}}(v...),
	}
}

func (t {{$typeFieldName}}) NotIn(v ...{{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotIn,
		Values:              common.Multi{{Title .FieldSelector}}{{Title .FieldType}}(v...),
	}
}

{{- if .FieldNullable}}
func (t {{$typeFieldName}}) IsNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNull,
	}
}

func (t {{$typeFieldName}}) IsNotNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNotNull,
	}
}
{{end -}}

{{- if .FieldLikable}}
func (t {{$typeFieldName}}) Like(vs ...{{$typeFieldType}}) sql.Condition {
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

func (t {{$typeFieldName}}) LikeOr(vs ...{{$typeFieldType}}) sql.Condition {
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

func (t {{$typeFieldName}}) NotLike(vs ...{{$typeFieldType}}) sql.Condition {
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

func (t {{$typeFieldName}}) NotLikeOr(vs ...{{$typeFieldType}}) sql.Condition {
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
func (t {{$typeFieldName}}) Between(v1 {{$typeFieldType}}, v2 {{$typeFieldType}}) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Between,
		Values:              []interface{}{v1, v2},
	}
}

func (t {{$typeFieldName}}) NotBetween(v1 {{$typeFieldType}}, v2 {{$typeFieldType}}) sql.Condition {
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
{{- $typeFieldType := AsTypeName .FieldSelector .FieldType -}}
	Equal({{$typeFieldType}}) sql.Condition
	NotEqual({{$typeFieldType}}) sql.Condition
	Greater({{$typeFieldType}}) sql.Condition
	GreaterOrEqual({{$typeFieldType}}) sql.Condition
	Lesser({{$typeFieldType}}) sql.Condition
	LesserOrEqual({{$typeFieldType}}) sql.Condition
	In(...{{$typeFieldType}}) sql.Condition
	NotIn(...{{$typeFieldType}}) sql.Condition
{{- end -}}

{{- define "BetweenFieldConditionOperations" -}}
{{- $typeFieldType := AsTypeName .FieldSelector .FieldType -}}
	Between({{$typeFieldType}}, {{$typeFieldType}}) sql.Condition
	NotBetween({{$typeFieldType}}, {{$typeFieldType}}) sql.Condition
{{- end -}}

{{- define "LikeFieldConditionOperations" -}}
{{- $typeFieldType := AsTypeName .FieldSelector .FieldType -}}
	Like(...{{$typeFieldType}}) sql.Condition
	NotLike(...{{$typeFieldType}}) sql.Condition
	LikeOr(...{{$typeFieldType}}) sql.Condition
	NotLikeOr(...{{$typeFieldType}}) sql.Condition
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
