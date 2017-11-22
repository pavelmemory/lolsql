package sql

import (
	"testing"
	"text/template"
	"bytes"
	"time"
)

type OrderStartField interface {
	Equal(time.Time) Condition
	NotEqual(time.Time) Condition
	In(...time.Time) Condition
	NotIn(...time.Time) Condition
}

func TestFieldInterfaceTemplate(t *testing.T) {
	var buffer bytes.Buffer
	tmpl := template.Must(template.New("FieldInterfaceDeclaration").Parse(FieldInterfaceDeclaration))
	err := tmpl.ExecuteTemplate(&buffer, "FieldInterfaceDeclaration", struct{
		TypeName string
		FieldName string
		FieldType string
		FieldTypePackageAlias string
		FieldNullable bool
		FieldLikable bool
		FieldBetweenable bool
	}{
		TypeName:"Order",
		FieldName:"Start",
		FieldType:"Time",
		FieldTypePackageAlias:"time",
		FieldNullable:true,
		FieldLikable:true,
		FieldBetweenable:true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buffer.String())
}
