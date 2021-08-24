package lvsdk

import (
	"encoding/json"
	"fmt"
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
	val, err := parseString(mm, "key")
	PanicIfError(err)
	assert.Equal(t, "string", val)
}

func TestSdkJson_parseStringNotFound(t *testing.T) {
	mm := make(Map)
	mm["key"] = "string"
	_, err := parseString(mm, "nf")
	assert.Equal(t, "key not found `nf`", err.Error())
}

func TestSdkJson_parseStringInvalidType(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := parseString(mm, "key")
	assert.Equal(t, "invalid type `key:bool`", err.Error())
}

func TestSdkJson_maybeMapSuccess(t *testing.T) {
	mm := make(Map)
	km := make(Map)
	mm["key"] = km
	val, err := maybeMap(mm, "key", nil)
	PanicIfError(err)
	assert.Equal(t, testMapPointer(km), testMapPointer(val))
}

func TestSdkJson_maybeMapInvalidType(t *testing.T) {
	mm := make(Map)
	mm["key"] = true
	_, err := maybeMap(mm, "key", nil)
	assert.Equal(t, "invalid type `key:bool`", err.Error())
}

func TestSdkJson_maybeMapDefault(t *testing.T) {
	mm := make(Map)
	km := make(Map)
	dm := make(Map)
	mm["key"] = km
	val, err := maybeMap(mm, "nf", dm)
	PanicIfError(err)
	assert.NotEqual(t, testMapPointer(dm), testMapPointer(km))
	assert.NotEqual(t, testMapPointer(km), testMapPointer(val))
	assert.Equal(t, testMapPointer(dm), testMapPointer(val))
}

func testMapPointer(mm Map) string {
	return fmt.Sprintf("%p", mm)
}

func TestSdkJson_encodeMutationSuccess(t *testing.T) {
	mut := &Mutation{}
	mut.Name = "name"
	mut.Sid = "sid"
	bytes, err := encodeMutation(mut)
	PanicIfError(err)
	assert.Equal(t, `{"name":"name","sid":"sid"}`, string(bytes))
	mut.Args = true
	bytes, err = encodeMutation(mut)
	PanicIfError(err)
	assert.Equal(t, `{"args":true,"name":"name","sid":"sid"}`, string(bytes))
	mut.Args = make(Map)
	bytes, err = encodeMutation(mut)
	PanicIfError(err)
	assert.Equal(t, `{"args":{},"name":"name","sid":"sid"}`, string(bytes))
}

func TestSdkJson_decodeMutationSuccess(t *testing.T) {
	mut, err := decodeMutation([]byte(`{"name":"name"}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "", mut.Sid)
	assert.Equal(t, nil, mut.Args)
	mut, err = decodeMutation([]byte(`{"name":"name","sid":"sid"}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "sid", mut.Sid)
	assert.Equal(t, nil, mut.Args)
	mut, err = decodeMutation([]byte(`{"name":"name","args":true}`))
	PanicIfError(err)
	assert.Equal(t, "name", mut.Name)
	assert.Equal(t, "", mut.Sid)
	assert.Equal(t, true, mut.Args)
	mut, err = decodeMutation([]byte(`{"name":"name","args":{}}`))
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
