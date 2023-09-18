package types

import "fmt"

type PropertyDependencyMapping struct {
	IsKey        bool
	ValuePath    string
	LiteralValue string
	Reference    *Reference
}

type Reference struct {
	Label    string // e.g. "test"
	Kind     string // e.g. "data", "resource", "var"
	Name     string // e.g. "azurerm_resource_group", "azapi_resource"
	Property string // e.g. "id", "name"
}

func (r Reference) String() string {
	if r.Kind == "resource" {
		return fmt.Sprintf("%s.%s.%s", r.Name, r.Label, r.Property)
	}
	return fmt.Sprintf("%s.%s.%s.%s", r.Kind, r.Name, r.Label, r.Property)
}

func (r *Reference) IsKnown() bool {
	return r != nil && r.Kind != "" && r.Label != "" && r.Name != "" && r.Property != ""
}

type Value interface {
	String() string
	DeepCopy() Value
}

type RawValue struct {
	Raw string
}

func (v RawValue) String() string {
	return v.Raw
}

func (v RawValue) DeepCopy() Value {
	return NewRawValue(v.Raw)
}

func NewRawValue(raw string) RawValue {
	return RawValue{
		Raw: raw,
	}
}

type ReferenceValue struct {
	Reference string
}

func (v ReferenceValue) String() string {
	return v.Reference
}

func (v ReferenceValue) DeepCopy() Value {
	return NewReferenceValue(v.Reference)
}

func NewReferenceValue(reference string) ReferenceValue {
	return ReferenceValue{
		Reference: reference,
	}
}

type StringLiteralValue struct {
	Literal string
}

func (v StringLiteralValue) String() string {
	return fmt.Sprintf(`"%s"`, v.Literal)
}

func (v StringLiteralValue) DeepCopy() Value {
	return NewStringLiteralValue(v.Literal)
}

func NewStringLiteralValue(literal string) StringLiteralValue {
	return StringLiteralValue{
		Literal: literal,
	}
}

var _ Value = &RawValue{}
var _ Value = &ReferenceValue{}
var _ Value = &StringLiteralValue{}
