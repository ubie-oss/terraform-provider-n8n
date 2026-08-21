package provider

import (
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func dynamicObjectToMap(v types.Dynamic) (map[string]any, error) {
	if v.IsNull() || v.IsUnknown() {
		return nil, fmt.Errorf("credential data must be a known object")
	}
	if v.IsUnderlyingValueNull() || v.IsUnderlyingValueUnknown() {
		return nil, fmt.Errorf("credential data must be a known object")
	}
	raw, err := attrValueToGo(v.UnderlyingValue())
	if err != nil {
		return nil, err
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("credential data must be an object or map, got %T", raw)
	}
	return m, nil
}

func attrValueToGo(v attr.Value) (any, error) {
	if v == nil || v.IsNull() {
		return nil, nil
	}
	if v.IsUnknown() {
		return nil, fmt.Errorf("credential data contains unknown values")
	}
	switch t := v.(type) {
	case types.Dynamic:
		return attrValueToGo(t.UnderlyingValue())
	case types.String:
		return t.ValueString(), nil
	case types.Bool:
		return t.ValueBool(), nil
	case types.Int64:
		return t.ValueInt64(), nil
	case types.Float64:
		return t.ValueFloat64(), nil
	case types.Number:
		return numberToGo(t.ValueBigFloat())
	case types.List:
		return collectionToGo(t.Elements())
	case types.Tuple:
		return collectionToGo(t.Elements())
	case types.Set:
		return collectionToGo(t.Elements())
	case types.Map:
		return mapAttrToGo(t.Elements())
	case types.Object:
		return mapAttrToGo(t.Attributes())
	default:
		return nil, fmt.Errorf("unsupported credential data type %T", v)
	}
}

func numberToGo(f *big.Float) (any, error) {
	if f == nil {
		return nil, nil
	}
	if f.IsInt() {
		i, acc := f.Int64()
		if acc == big.Exact {
			return i, nil
		}
	}
	fl, _ := f.Float64()
	return fl, nil
}

func collectionToGo(elems []attr.Value) ([]any, error) {
	out := make([]any, 0, len(elems))
	for _, e := range elems {
		v, err := attrValueToGo(e)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func mapAttrToGo(elems map[string]attr.Value) (map[string]any, error) {
	out := make(map[string]any, len(elems))
	for k, e := range elems {
		v, err := attrValueToGo(e)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}
