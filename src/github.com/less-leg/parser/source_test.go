package parser

import (
	"testing"
)

func TestReadPackageInfo(t *testing.T) {
	pkg, err := ReadPackageInfo("github.com/less-leg/test_model")
	if err != nil {
		t.Fatal(err)
	}

	types := pkg.Types
	if len(types) != 3 {
		t.Fatal(types)
	}

	myTimeType := types[TypeIdentity{Name: "MyTime", Package: "github.com/less-leg/test_model"}]
	if !myTimeType.IsItAlias() {
		t.Fatal("it is must be an alias", myTimeType)
	}
	myTimeAliasType := myTimeType.(UserDefinedAlias)
	if myTimeAliasType.Package != "github.com/less-leg/test_model" {
		t.Fatal(myTimeAliasType.Package)
	}
	if myTimeAliasType.Name != "MyTime" {
		t.Fatal(myTimeAliasType.Name)
	}
	if myTimeAliasType.ActualType != (TypeIdentityRef{Name: "Time", Selector: "time"}) {
		t.Fatal(myTimeAliasType.Name)
	}

	//userType := types["User"].(UserDefinedType)
	//if len(userType.Fields) != 3 {
	//	t.Fatal(userType.Fields)
	//}
	//
	//operationType := types["Operation"].(UserDefinedType)
	//if len(operationType.Fields) != 3 {
	//	t.Fatal(operationType.Fields)
	//}

	for _, typ := range types {
		t.Logf("%#v\n", typ)
	}
}
