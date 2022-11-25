package omodel

type AnyValue struct {
	Val AnyValueDatum
}

type AnyValueDatum interface{}

type AnyValueString struct {
	StringValue string
}
type AnyValueBool struct {
	BoolValue bool
}
type AnyValueInt struct {
	IntValue int64
}
type AnyValueDouble struct {
	DoubleValue float64
}
type AnyValueArray struct {
	ArrayValues []AnyValueDatum
}
type AnyValueKV struct {
	KVValue map[string]AnyValueDatum
}
type AnyValueBytes struct {
	BytesValue []byte
}

func (m *AnyValue) GetStringValue() string {
	if x, ok := m.Val.(*AnyValueString); ok {
		return x.StringValue
	}
	return ""
}

func (m *AnyValue) GetBoolValue() bool {
	if x, ok := m.Val.(*AnyValueBool); ok {
		return x.BoolValue
	}
	return false
}

func (m *AnyValue) GetIntValue() int64 {
	if x, ok := m.Val.(*AnyValueInt); ok {
		return x.IntValue
	}
	return 0
}

func (m *AnyValue) GetDoubleValue() float64 {
	if x, ok := m.Val.(*AnyValueDouble); ok {
		return x.DoubleValue
	}
	return 0
}

func (m *AnyValue) GetArrayValue() []AnyValueDatum {
	if x, ok := m.Val.(*AnyValueArray); ok {
		return x.ArrayValues
	}
	return nil
}

func (m *AnyValue) GetKV() map[string]AnyValueDatum {
	if x, ok := m.Val.(*AnyValueKV); ok {
		return x.KVValue
	}
	return nil
}

func (m *AnyValue) GetBytesValue() []byte {
	if x, ok := m.Val.(*AnyValueBytes); ok {
		return x.BytesValue
	}
	return nil
}
