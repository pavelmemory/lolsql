package sql

import "time"

type Condition interface {
	AppendValues([]interface{}) []interface{}
	And(condition Condition) Condition
	Or(condition Condition) Condition
}

type ComparatorOperation byte

const (
	Equal ComparatorOperation = iota
	NotEqual
	Greater
	GreaterOrEqual
	Lesser
	LesserOrEqual
	In
	NotIn
	IsNull
	IsNotNull
	Like
	NotLike
	LikeOr
	NotLikeOr
	Between
	NotBetween
)

type LogicalOperator byte

const (
	Conjunction LogicalOperator = iota
	Disjunction
)

type MultiCondition struct {
	Conditions []Condition
	LogicalOperator
}

func (c MultiCondition) AppendValues(vals []interface{}) []interface{} {
	for _, cond := range c.Conditions {
		vals = cond.AppendValues(vals)
	}
	return vals
}

func (c MultiCondition) And(condition Condition) Condition {
	return MultiCondition{Conditions: append(c.Conditions, condition), LogicalOperator: Conjunction}
}

func (c MultiCondition) Or(condition Condition) Condition {
	return MultiCondition{Conditions: append(c.Conditions, condition), LogicalOperator: Disjunction}
}

type SingleCondition struct {
	Field
	ComparatorOperation
	Values []interface{}
}

func MultiTimes(times ...time.Time) (vals []interface{}) {
	for _, v := range times {
		vals = append(vals, v)
	}
	return
}

func MultiBytes(bytes ...byte) (vals []interface{}) {
	for _, v := range bytes {
		vals = append(vals, v)
	}
	return
}

func MultiStrings(strs ...string) (vals []interface{}) {
	for _, v := range strs {
		vals = append(vals, v)
	}
	return
}

func MultiInts(ints ...int) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiInt8s(ints ...int8) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiInt16s(ints ...int16) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiInt32s(ints ...int32) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiInt64s(ints ...int64) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiUints(ints ...uint) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiUint8s(ints ...uint8) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiUint16s(ints ...uint16) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiUint32s(ints ...uint32) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiUint64s(ints ...uint64) (vals []interface{}) {
	for _, v := range ints {
		vals = append(vals, v)
	}
	return
}

func MultiFloat32s(floats ...float32) (vals []interface{}) {
	for _, v := range floats {
		vals = append(vals, v)
	}
	return
}

func MultiFloat64s(floats ...float64) (vals []interface{}) {
	for _, v := range floats {
		vals = append(vals, v)
	}
	return
}

func (c SingleCondition) GetFields() []Field {
	return []Field{c.Field}
}

func (c SingleCondition) AppendValues(vals []interface{}) []interface{} {
	return append(vals, c.Values...)
}

func (c SingleCondition) And(condition Condition) Condition {
	return MultiCondition{Conditions: append([]Condition{c}, condition), LogicalOperator: Conjunction}
}

func (c SingleCondition) Or(condition Condition) Condition {
	return MultiCondition{Conditions: append([]Condition{c}, condition), LogicalOperator: Disjunction}
}

func And(condition Condition, conditions ...Condition) Condition {
	return MultiCondition{Conditions: append([]Condition{condition}, conditions...), LogicalOperator: Conjunction}
}

func Or(condition Condition, conditions ...Condition) Condition {
	return MultiCondition{Conditions: append([]Condition{condition}, conditions...), LogicalOperator: Disjunction}
}

//
//type ByteEquable interface {
//	Equal(byte) Condition
//}
//
//type StringEquable interface {
//	Equal(string) Condition
//}
//
//type IntEquable interface {
//	Equal(int) Condition
//}
//
//type Int8Equable interface {
//	Equal(int8) Condition
//}
//
//type Int16Equable interface {
//	Equal(int16) Condition
//}
//
//type Int32Equable interface {
//	Equal(int32) Condition
//}
//
//type Int64Equable interface {
//	Equal(int64) Condition
//}
//
//type UintEquable interface {
//	Equal(uint) Condition
//}
//
//type Uint8Equable interface {
//	Equal(uint8) Condition
//}
//
//type Uint16Equable interface {
//	Equal(uint16) Condition
//}
//
//type Uint32Equable interface {
//	Equal(uint32) Condition
//}
//
//type Uint64Equable interface {
//	Equal(uint64) Condition
//}
//
//type Float32Equable interface {
//	Equal(float32) Condition
//}
//
//type Float64Equable interface {
//	Equal(float64) Condition
//}
//
//type Likable interface {
//	Like(v string, vs ...string) Condition
//	NotLike(v string, vs ...string) Condition
//	LikeOr(v string, vs ...string) Condition
//	NotLikeOr(v string, vs ...string) Condition
//}

var FieldInterfaceDeclaration = `
{{- define "FieldInterfaceDeclaration" -}}
type {{.TypeName}}{{.FieldName}}Field interface {
	{{.TypeName}}Field
	{{template "CommonFieldConditionOperations" . }}
	{{if .FieldNullable}}{{template "NullableFieldConditionOperations" . }}{{end}}
	{{if .FieldLikable}}{{template "LikeFieldConditionOperations" . }}{{end}}
	{{if .FieldBetweenable}}{{template "BetweenFieldConditionOperations" . }}{{end}}
}
{{- end -}}

{{- define "CommonFieldConditionOperations" -}}
	Equal({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	NotEqual({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	Greater({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	GreaterOrEqual({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	Lesser({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	LesserOrEqual({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	In(...{{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	NotIn(...{{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
{{- end -}}

{{- define "BetweenFieldConditionOperations" -}}
	Between({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}, {{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	NotBetween({{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}, {{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
{{- end -}}

{{- define "LikeFieldConditionOperations" -}}
	Like(...{{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	NotLike(...{{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	LikeOr(...{{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
	NotLikeOr(...{{.FieldTypePackageAlias}}{{if .FieldTypePackageAlias}}.{{end}}{{.FieldType}}) sql.Condition
{{- end -}}

{{- define "NullableFieldConditionOperations" -}}
	IsNull() sql.Condition
	IsNotNull() sql.Condition
{{- end -}}`
