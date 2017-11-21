package parser

type TypeIdentity struct {
	Name    string
	Package string
}

type TypeIdentityRef struct {
	Name     string
	Selector string
}

type Field interface {
	GetName() string
	GetTypeName() string
	GetSelector() string
	GetTag() string
	IsItSlice() bool
	IsItReference() bool
	IsItEmbedded() bool
	GetTypeSpecifications() []TypeSpecification
}

type TypeSpecification string

const (
	Reference TypeSpecification = "*"
	Slice     TypeSpecification = "[]"
)

type UserDefinedTypeField struct {
	Name                      string
	TypeName                  string
	Selector                  string
	Tag                       string
	IsSlice                   bool
	IsReference               bool
	IsEmbedded                bool
	OrderedTypeSpecifications []TypeSpecification
}

var _ Field = UserDefinedTypeField{}

func (ufield UserDefinedTypeField) GetName() string {
	return ufield.Name
}

func (ufield UserDefinedTypeField) GetTypeName() string {
	return ufield.TypeName
}

func (ufield UserDefinedTypeField) GetSelector() string {
	return ufield.Selector
}

func (ufield UserDefinedTypeField) GetTag() string {
	return ufield.Tag
}

func (ufield UserDefinedTypeField) IsItSlice() bool {
	return ufield.IsSlice
}

func (ufield UserDefinedTypeField) IsItReference() bool {
	return ufield.IsReference
}

func (ufield UserDefinedTypeField) IsItEmbedded() bool {
	return ufield.IsEmbedded
}

func (ufield UserDefinedTypeField) GetTypeSpecifications() []TypeSpecification {
	return ufield.OrderedTypeSpecifications
}

type Type interface {
	GetIdentity() TypeIdentity
	GetFields() map[string]Field
	IsItFromStdlib() bool
	IsItAlias() bool
}

type StdlibType struct {
	TypeIdentity
}

var _ Type = StdlibType{}

func (std StdlibType) GetIdentity() TypeIdentity {
	return std.TypeIdentity
}

func (std StdlibType) GetFields() map[string]Field {
	return nil
}

func (std StdlibType) IsItFromStdlib() bool {
	return true
}

func (std StdlibType) IsItAlias() bool {
	return false
}

type UserDefinedType struct {
	TypeIdentity
	Fields map[string]Field
}

var _ Type = UserDefinedType{}

func (udef UserDefinedType) GetIdentity() TypeIdentity {
	return udef.TypeIdentity
}

func (udef UserDefinedType) GetFields() map[string]Field {
	return udef.Fields
}

func (udef UserDefinedType) IsItFromStdlib() bool {
	return false
}

func (udef UserDefinedType) IsItAlias() bool {
	return false
}

type UserDefinedAlias struct {
	TypeIdentity
	ActualType TypeIdentityRef
}

var _ Type = UserDefinedAlias{}

func (udefa UserDefinedAlias) GetIdentity() TypeIdentity {
	return udefa.TypeIdentity
}

func (udefa UserDefinedAlias) GetFields() map[string]Field {
	return nil
}

func (udefa UserDefinedAlias) IsItFromStdlib() bool {
	return false
}

func (udefa UserDefinedAlias) IsItAlias() bool {
	return true
}

type Import struct {
	Alias string
	Path  string
}
