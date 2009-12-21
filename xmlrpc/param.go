package xmlrpc

import (
	"reflect"
	"fmt"
)

func Params(params ...) []ParamValue {
	pStruct := reflect.NewValue(params).(*reflect.StructValue)
	par := make([]ParamValue, pStruct.NumField())

	for n := 0; n < len(par); n++ {
		par[n] = param(pStruct.Field(n))
	}
	return par
}

func structParams(v *reflect.StructValue) StructValue {
	p := make(StructValue, v.NumField())
	for n := 0; n < v.NumField(); n++ {
		key := v.Type().(*reflect.StructType).Field(n).Name
		p[key] = param(v.Field(n))
	}
	return p
}

func byteParams(params *reflect.SliceValue) Base64Value {
	par := make([]byte, params.Len())

	for n := 0; n < len(par); n++ {
		par[n] = params.Elem(n).(*reflect.Uint8Value).Get()
	}
	return Base64Value(par)
}

func arrayParams(params *reflect.SliceValue) ArrayValue {
	par := make([]ParamValue, params.Len())

	for n := 0; n < len(par); n++ {
		par[n] = param(params.Elem(n))
	}
	return par
}

func param(param interface{}) ParamValue {
	switch v := param.(type) {
	case *reflect.IntValue:
		return IntValue(v.Get())
	case *reflect.BoolValue:
		return BooleanValue(v.Get())
	case *reflect.StringValue:
		return StringValue(v.Get())
	case *reflect.FloatValue:
		return DoubleValue(v.Get())
	case *reflect.SliceValue:
		if _, ok := v.Type().(*reflect.SliceType).Elem().(*reflect.Uint8Type); ok { // A []byte is really a Base64Type
			return byteParams(v)
		}
		return arrayParams(v)
	case *reflect.StructValue:
		return structParams(v)
	}
	return StringValue(fmt.Sprintf("Error: Unknown Param Type (%T)\n", param))
}
