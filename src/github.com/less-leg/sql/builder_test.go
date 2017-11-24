package sql

import (
	"bytes"
	"testing"
	"text/template"
)

func TestTypeSelectTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.New("").Funcs(TemplateFunctions)
	tmpl = template.Must(tmpl.Parse(TypeSelect))
	err := tmpl.Execute(&buffer, struct {
		TypeName     string
		TypeSelector string
	}{
		TypeName:     "Order",
		TypeSelector: "test_model",
	})
	if err != nil {
		t.Fatal(err)
	}
	if buffer.String() !=
		`func Order(fields ...OrderField) orderBuilder {
	return orderBuilder{}
}

func (b orderBuilder) Where(...sql.Condition) orderBuilder {
	return b
}

func (b orderBuilder) GroupBy(column ...OrderField) orderBuilder {
	return b
}

func (b orderBuilder) Having(...sql.Condition) orderBuilder {
	return b
}

func (b orderBuilder) OrderBy(...sql.SortOrder) orderBuilder {
	return b
}

func (b orderBuilder) Get() ([]test_model.Order, error) {
	return nil, nil
}

func (b orderBuilder) GetPtr() ([]*test_model.Order, error) {
	return nil, nil
}` {
		t.Fatal(buffer.String())
	}
}
