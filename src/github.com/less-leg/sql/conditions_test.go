package sql

import (
	"bytes"
	"testing"
	"text/template"
)

func TestFieldConditionInterfaceDeclarationTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.New("").Funcs(TemplateFunctions)
	tmpl = template.Must(tmpl.Parse(FieldConditionInterfaceDeclaration))
	err := tmpl.Execute(&buffer, struct {
		TypeName         string
		FieldName        string
		FieldType        string
		FieldSelector    string
		FieldNullable    bool
		FieldLikable     bool
		FieldBetweenable bool
	}{
		TypeName:         "Order",
		FieldName:        "Start",
		FieldType:        "Time",
		FieldSelector:    "time",
		FieldNullable:    true,
		FieldLikable:     true,
		FieldBetweenable: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() != `type OrderStartField interface {
	OrderField
	Equal(time.Time) sql.Condition
	NotEqual(time.Time) sql.Condition
	Greater(time.Time) sql.Condition
	GreaterOrEqual(time.Time) sql.Condition
	Lesser(time.Time) sql.Condition
	LesserOrEqual(time.Time) sql.Condition
	In(...time.Time) sql.Condition
	NotIn(...time.Time) sql.Condition
	IsNull() sql.Condition
	IsNotNull() sql.Condition
	Like(...time.Time) sql.Condition
	NotLike(...time.Time) sql.Condition
	LikeOr(...time.Time) sql.Condition
	NotLikeOr(...time.Time) sql.Condition
	Between(time.Time, time.Time) sql.Condition
	NotBetween(time.Time, time.Time) sql.Condition
}` {
		t.Fatal(buffer.String())
	}
}

func TestFieldBaseDeclarationTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.New("").Funcs(TemplateFunctions)
	tmpl = template.Must(tmpl.Parse(FieldBaseDeclaration))
	err := tmpl.Execute(&buffer, struct {
		TypeName string
		Package  string
	}{
		TypeName: "Order",
		Package:  "github.com/less-leg/test_model",
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() != `type OrderField interface {
	sql.Field
}

type orderField struct {}

func (orderField) GetType() parser.TypeIdentity {
	return parser.TypeIdentity{Name:"Order", Package:"github.com/less-leg/test_model"}
}` {
		t.Fatal(buffer.String())
	}
	t.Log(buffer.String())
}

func TestFieldDeclarationTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.New("").Funcs(TemplateFunctions)
	tmpl = template.Must(tmpl.Parse(FieldDeclaration))
	err := tmpl.Execute(&buffer, struct {
		TypeName  string
		FieldName string
	}{
		TypeName:  "Order",
		FieldName: "Start",
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() != `type orderStartField struct {
	orderField
}

func (orderStartField) GetName() string {
	return "Start"
}

func Start() OrderStartField {
	return orderStartField{}
}` {
		t.Fatal(buffer.String())
	}
	t.Log(buffer.String())
}

func TestSliceXTypeToSliceInterfacesTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.New("").Funcs(TemplateFunctions)
	tmpl = template.Must(tmpl.Parse(SliceXTypeToSliceInterfaces))
	err := tmpl.Execute(&buffer, struct {
		TypeName string
		Selector string
	}{
		TypeName: "Time",
		Selector: "time",
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() != `func MultiTimeTime(times ...time.Time) (vals []interface{}) {
	for _, v := range times {
		vals = append(vals, v)
	}
	return
}` {
		t.Fatal(buffer.String())
	}
	buffer.Reset()

	err = tmpl.Execute(&buffer, struct {
		TypeName string
		Selector string
	}{
		TypeName: "string",
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() != `func MultiString(strings ...string) (vals []interface{}) {
	for _, v := range strings {
		vals = append(vals, v)
	}
	return
}` {
		t.Fatal(buffer.String())
	}
}

func TestFieldConditionDeclarationTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.New("").Funcs(TemplateFunctions)
	tmpl = template.Must(tmpl.Parse(FieldConditionDeclaration))
	err := tmpl.Execute(&buffer, struct {
		TypeName         string
		FieldName        string
		FieldSelector    string
		FieldType        string
		FieldNullable    bool
		FieldLikable     bool
		FieldBetweenable bool
	}{
		TypeName:         "Order",
		FieldName:        "Start",
		FieldSelector:    "time",
		FieldType:        "Time",
		FieldNullable:    true,
		FieldLikable:     true,
		FieldBetweenable: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() != `func (t orderStartField) Equal(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Equal,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) NotEqual(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotEqual,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) Greater(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Greater,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) GreaterOrEqual(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.GreaterOrEqual,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) Lesser(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Lesser,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) LesserOrEqual(v time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.LesserOrEqual,
		Values:              []interface{}{v},
	}
}

func (t orderStartField) In(v ...time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.In,
		Values:              common.MultiTimeTime(v...),
	}
}

func (t orderStartField) NotIn(v ...time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotIn,
		Values:              common.MultiTimeTime(v...),
	}
}
func (t orderStartField) IsNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNull,
	}
}

func (t orderStartField) IsNotNull() sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.IsNotNull,
	}
}

func (t orderStartField) Like(vs ...time.Time) sql.Condition {
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

func (t orderStartField) LikeOr(vs ...time.Time) sql.Condition {
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

func (t orderStartField) NotLike(vs ...time.Time) sql.Condition {
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

func (t orderStartField) NotLikeOr(vs ...time.Time) sql.Condition {
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

func (t orderStartField) Between(v1 time.Time, v2 time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.Between,
		Values:              []interface{}{v1, v2},
	}
}

func (t orderStartField) NotBetween(v1 time.Time, v2 time.Time) sql.Condition {
	return sql.SingleCondition{
		Type:                t.GetType(),
		Field:               t.GetName(),
		ComparatorOperation: sql.NotBetween,
		Values:              []interface{}{v1, v2},
	}
}` {
		t.Fatal(buffer.String())
	}
}
