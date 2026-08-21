package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDynamicObjectToMap(t *testing.T) {
	obj, diags := types.ObjectValue(
		map[string]attr.Type{"name": types.StringType, "value": types.StringType},
		map[string]attr.Value{"name": types.StringValue("X-Test"), "value": types.StringValue("placeholder")},
	)
	if diags.HasError() {
		t.Fatalf("object: %v", diags)
	}
	got, err := dynamicObjectToMap(types.DynamicValue(obj))
	if err != nil {
		t.Fatal(err)
	}
	if got["name"] != "X-Test" || got["value"] != "placeholder" {
		t.Fatalf("got %#v", got)
	}
}

func TestDynamicObjectToMapRejectsList(t *testing.T) {
	list, diags := types.ListValue(types.StringType, []attr.Value{types.StringValue("a")})
	if diags.HasError() {
		t.Fatalf("list: %v", diags)
	}
	_, err := dynamicObjectToMap(types.DynamicValue(list))
	if err == nil {
		t.Fatal("expected error for list")
	}
}

func TestDynamicObjectToMapRejectsBool(t *testing.T) {
	_, err := dynamicObjectToMap(types.DynamicValue(types.BoolValue(true)))
	if err == nil {
		t.Fatal("expected error for bool")
	}
}

func TestDynamicObjectToMapRejectsNull(t *testing.T) {
	_, err := dynamicObjectToMap(types.DynamicNull())
	if err == nil {
		t.Fatal("expected error for null")
	}
}

func TestOptionalBoolPointer(t *testing.T) {
	if optionalBoolPointer(types.BoolNull()) != nil {
		t.Fatal("null should be nil")
	}
	p := optionalBoolPointer(types.BoolValue(true))
	if p == nil || !*p {
		t.Fatalf("got %v", p)
	}
}
