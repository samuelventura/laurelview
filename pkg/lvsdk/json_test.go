package lvsdk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSdkJson_getValueSuccess(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	val, err := getValue(mm, "key")
	PanicIfError(err)
	assert.Equal(t, true, val)
}

func TestSdkJson_getValueNotFound(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := getValue(mm, "nf")
	assert.Equal(t, "key not found `nf`", err.Error())
}

func TestSdkJson_parseStringSuccess(t *testing.T) {
	mm := make(Map)
	mm["key"] = "string"
	val, err := ParseString(mm, "key")
	PanicIfError(err)
	assert.Equal(t, "string", val)
}

func TestSdkJson_parseStringNotFound(t *testing.T) {
	mm := make(Map)
	mm["key"] = "string"
	_, err := ParseString(mm, "nf")
	assert.Equal(t, "key not found `nf`", err.Error())
}

func TestSdkJson_parseStringInvalidType(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := ParseString(mm, "key")
	assert.Equal(t, "invalid type `key:bool`", err.Error())
}

func TestSdkJson_parseUintSuccessFloat(t *testing.T) {
	mm := make(Map)
	mm["key"] = float64(1)
	val, err := ParseUint(mm, "key")
	PanicIfError(err)
	assert.Equal(t, uint(1), val)
}

func TestSdkJson_parseUintSuccessString(t *testing.T) {
	mm := make(Map)
	mm["key"] = "1"
	val, err := ParseUint(mm, "key")
	PanicIfError(err)
	assert.Equal(t, uint(1), val)
}

func TestSdkJson_parseUintNotFound(t *testing.T) {
	mm := make(Map)
	mm["key"] = float64(1)
	_, err := ParseUint(mm, "nf")
	assert.Equal(t, "key not found `nf`", err.Error())
}

func TestSdkJson_parseUintInvalidType(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := ParseUint(mm, "key")
	assert.Equal(t, "invalid type `key:bool`", err.Error())
}

func TestSdkJson_maybeStringSuccess(t *testing.T) {
	mm := make(Map)
	mm["key"] = "value"
	val, err := MaybeString(mm, "key", "default")
	PanicIfError(err)
	assert.Equal(t, "value", val)
}

func TestSdkJson_maybeStringDefault(t *testing.T) {
	mm := make(Map)
	mm["key"] = "value"
	val, err := MaybeString(mm, "nf", "default")
	PanicIfError(err)
	assert.Equal(t, "default", val)
}

func TestSdkJson_maybeStringInvalidType(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := MaybeString(mm, "key", "default")
	assert.Equal(t, "invalid type `key:bool`", err.Error())
}

func TestSdkJson_maybeUintSuccessFloat(t *testing.T) {
	mm := make(Map)
	mm["key"] = float64(1)
	val, err := MaybeUint(mm, "key", 2)
	PanicIfError(err)
	assert.Equal(t, uint(1), val)
}

func TestSdkJson_maybeUintSuccessString(t *testing.T) {
	mm := make(Map)
	mm["key"] = "1"
	val, err := MaybeUint(mm, "key", 2)
	PanicIfError(err)
	assert.Equal(t, uint(1), val)
}

func TestSdkJson_maybeUintDefault(t *testing.T) {
	mm := make(Map)
	mm["key"] = float64(1)
	val, err := MaybeUint(mm, "nf", 2)
	PanicIfError(err)
	assert.Equal(t, uint(2), val)
}

func TestSdkJson_maybeUintInvalidType(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := MaybeUint(mm, "key", 2)
	assert.Equal(t, "invalid type `key:bool`", err.Error())
}

func TestSdkJson_encodeMutationSuccess(t *testing.T) {
	mut := &Mutation{}
	mut.Name = "name"
	mut.Sid = "sid"
	bytes, err := EncodeMutation(mut)
	PanicIfError(err)
	assert.Equal(t, `{"name":"name","sid":"sid"}`, string(bytes))
	mut.Args = true
	bytes, err = EncodeMutation(mut)
	PanicIfError(err)
	assert.Equal(t, `{"args":true,"name":"name","sid":"sid"}`, string(bytes))
	mut.Args = make(Map)
	bytes, err = EncodeMutation(mut)
	PanicIfError(err)
	assert.Equal(t, `{"args":{},"name":"name","sid":"sid"}`, string(bytes))
}

func TestSdkJson_decodeMutationSuccess(t *testing.T) {
	mut, err := DecodeMutation([]byte(`{"name":"name"}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "", mut.Sid)
	assert.Equal(t, nil, mut.Args)
	mut, err = DecodeMutation([]byte(`{"name":"name","sid":"sid"}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "sid", mut.Sid)
	assert.Equal(t, nil, mut.Args)
	mut, err = DecodeMutation([]byte(`{"name":"name","args":true}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "", mut.Sid)
	assert.Equal(t, true, mut.Args)
	mut, err = DecodeMutation([]byte(`{"name":"name","args":{}}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "", mut.Sid)
	//empty maps are always equal
	assert.Equal(t, make(Map), make(Map))
	assert.Equal(t, make(Map), mut.Args)
}

func TestSdkJson_unmarshalInteger(t *testing.T) {
	var val Any
	err := json.Unmarshal([]byte("1"), &val)
	PanicIfError(err)
	assert.Equal(t, float64(1), val)
}

func TestSdkJson_unmarshalFloating(t *testing.T) {
	var val Any
	err := json.Unmarshal([]byte("1.2"), &val)
	PanicIfError(err)
	assert.Equal(t, float64(1.2), val)
}

func TestSdkJson_unmarshalBoolean(t *testing.T) {
	var val Any
	err := json.Unmarshal([]byte("true"), &val)
	PanicIfError(err)
	assert.Equal(t, true, val)
}
