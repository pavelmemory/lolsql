package parser

import (
	"testing"
)

func TestReadPackageInfo_Import(t *testing.T) {
	pkg, err := ReadPackageInfo("github.com/less-leg/test_model")
	if err != nil {
		t.Fatal(err)
	}

	if len(pkg.Imports) != 2 {
		t.Fatal(pkg.Imports)
	}

	assertImport(t, pkg, Import{Path: "time", Alias: "time"})
	assertImport(t, pkg, Import{Path: "github.com/less-leg/dbmodel", Alias: "dbmodel"})
}

func TestReadPackageInfo_UserDefinedAlias(t *testing.T) {
	pkg, err := ReadPackageInfo("github.com/less-leg/test_model")
	if err != nil {
		t.Fatal(err)
	}

	assertUserDefinedAlias(t, pkg, UserDefinedAlias{
		TypeIdentity: TypeIdentity{Name: "MyTime", Package: "github.com/less-leg/test_model"},
		ActualType:   TypeIdentityRef{Name: "Time", Selector: "time"},
	})
}

func TestReadPackageInfo_UserDefinedType(t *testing.T) {
	pkg, err := ReadPackageInfo("github.com/less-leg/test_model")
	if err != nil {
		t.Fatal(err)
	}

	assertUserDefinedType(t, pkg, UserDefinedType{
		TypeIdentity: TypeIdentity{Name: "User", Package: "github.com/less-leg/test_model"},
		Fields: map[string]Field{
			"Name":    UserDefinedTypeField{Name: "Name", TypeName: "string"},
			"Father":  UserDefinedTypeField{Name: "Father", TypeName: "User", IsReference: true, Tag: `json:"omitempty"`, OrderedTypeSpecifications: []TypeSpecification{Reference}},
			"Time":    UserDefinedTypeField{Name: "Time", TypeName: "Time", Selector: "time", IsEmbedded: true},
			"TaxFree": UserDefinedTypeField{Name: "TaxFree", TypeName: "Confirmation", Selector: "dbmodel"},
		},
	})

	assertUserDefinedType(t, pkg, UserDefinedType{
		TypeIdentity: TypeIdentity{Name: "Operation", Package: "github.com/less-leg/test_model"},
		Fields: map[string]Field{
			"Begin":   UserDefinedTypeField{Name: "Begin", TypeName: "Time", Selector: "time"},
			"End":     UserDefinedTypeField{Name: "End", TypeName: "Time", Selector: "time"},
			"User":    UserDefinedTypeField{Name: "User", TypeName: "User", IsEmbedded: true, IsReference: true, OrderedTypeSpecifications: []TypeSpecification{Reference}},
			"Actions": UserDefinedTypeField{Name: "Actions", TypeName: "byte", IsSlice: true, IsReference: true, OrderedTypeSpecifications: []TypeSpecification{Slice, Reference}},
		},
	})
}

func assertImport(t *testing.T, pkg Package, imp Import) {
	if fImp, found := pkg.Imports[imp.Path]; !found {
		t.Fatal(pkg.Imports)
	} else {
		if fImp.Path != imp.Path {
			t.Fatal(fImp.Path)
		}
		if fImp.Alias != imp.Alias {
			t.Fatal(fImp)
		}
	}
}

func assertUserDefinedAlias(t *testing.T, pkg Package, uda UserDefinedAlias) {
	fUda := pkg.Types[uda.GetIdentity()].(UserDefinedAlias)
	if !fUda.IsItAlias() {
		t.Fatal("it is must be an alias", fUda)
	}
	if fUda.Package != uda.Package {
		t.Fatal(fUda.Package)
	}
	if fUda.Name != uda.Name {
		t.Fatal(fUda.Name)
	}
	if fUda.ActualType != uda.ActualType {
		t.Fatal(fUda.ActualType)
	}
}

func assertUserDefinedType(t *testing.T, pkg Package, udt UserDefinedType) {
	fType := pkg.Types[udt.GetIdentity()]
	if fType.GetIdentity() != udt.GetIdentity() {
		t.Fatal(fType.GetIdentity(), udt.GetIdentity())
	}
	if fType.IsItAlias() != udt.IsItAlias() {
		t.Fatal(fType.IsItAlias(), udt.IsItAlias())
	}
	if fType.IsItFromStdlib() != udt.IsItFromStdlib() {
		t.Fatal(fType.IsItFromStdlib(), udt.IsItFromStdlib())
	}

	fFields := fType.GetFields()
	fFieldsCopy := make(map[string]Field)
	for fName, fField := range fFields {
		udtField, found := udt.Fields[fName]
		if !found {
			t.Fatal(fName)
		}
		assertUserDefinedTypeField(t, udtField, fField)
		fFieldsCopy[fName] = fField
	}
}

func assertUserDefinedTypeField(t *testing.T, expectedUdf Field, actualUdf Field) {
	if expectedUdf.GetName() != actualUdf.GetName() {
		t.Fatal(expectedUdf.GetName(), actualUdf.GetName())
	}
	if expectedUdf.GetSelector() != actualUdf.GetSelector() {
		t.Fatal(expectedUdf.GetSelector(), actualUdf.GetSelector())
	}
	if expectedUdf.GetTypeName() != actualUdf.GetTypeName() {
		t.Fatal(expectedUdf.GetTypeName(), actualUdf.GetTypeName())
	}
	if expectedUdf.GetTag() != actualUdf.GetTag() {
		t.Fatal(expectedUdf.GetTag(), actualUdf.GetTag())
	}
	if expectedUdf.IsItEmbedded() != actualUdf.IsItEmbedded() {
		t.Fatal(expectedUdf.IsItEmbedded(), actualUdf.IsItEmbedded())
	}
	if expectedUdf.IsItReference() != actualUdf.IsItReference() {
		t.Fatal(expectedUdf.IsItReference(), actualUdf.IsItReference())
	}
	if expectedUdf.IsItSlice() != actualUdf.IsItSlice() {
		t.Fatal(expectedUdf.IsItSlice(), actualUdf.IsItSlice())
	}
	if len(expectedUdf.GetTypeSpecifications()) != len(actualUdf.GetTypeSpecifications()) {
		t.Fatal(len(expectedUdf.GetTypeSpecifications()), len(actualUdf.GetTypeSpecifications()))
	}
	expectedUdfSpecs := expectedUdf.GetTypeSpecifications()
	actualUdfSpecs := actualUdf.GetTypeSpecifications()
	for i := 0; i < len(expectedUdfSpecs); i++ {
		if expectedUdfSpecs[i] != actualUdfSpecs[i] {
			t.Fatal(i, expectedUdfSpecs[i], actualUdfSpecs[i])
		}
	}
}
