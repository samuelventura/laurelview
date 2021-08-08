package lvnbe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateCrud(t *testing.T) {
	var dao = NewDao(":memory:")
	var state = NewState(dao)

	assert.Equal(t, 0, len(state.All().Items))

	i1 := &CreateArgs{Name: "name1", Json: "json1"}
	err := state.Apply(&Mutation{Name: "create", Args: i1})
	PanicIfError(err)
	assert.Equal(t, uint(1), i1.Id)
	assert.Equal(t, "name1", i1.Name)
	assert.Equal(t, "json1", i1.Json)
	assert.Equal(t, 1, len(state.All().Items))

	i2 := &CreateArgs{Name: "name2", Json: "json2"}
	err = state.Apply(&Mutation{Name: "create", Args: i2})
	PanicIfError(err)
	assert.Equal(t, uint(2), i2.Id)
	assert.Equal(t, "name2", i2.Name)
	assert.Equal(t, "json2", i2.Json)
	assert.Equal(t, 2, len(state.All().Items))

	err = state.Apply(&Mutation{Name: "update",
		Args: &UpdateArgs{i2.Id, "name3", "json3"}})
	PanicIfError(err)

	all := state.All()
	assert.Equal(t, "name1", all.Items[0].Name)
	assert.Equal(t, "json1", all.Items[0].Json)
	assert.Equal(t, "name3", all.Items[1].Name)
	assert.Equal(t, "json3", all.Items[1].Json)

	err = state.Apply(&Mutation{Name: "delete", Args: &DeleteArgs{i2.Id}})
	PanicIfError(err)
	assert.Equal(t, 1, len(state.All().Items))
}

func TestStateDeleteError(t *testing.T) {
	var dao = NewDao(":memory:")
	var state = NewState(dao)
	err := state.Apply(&Mutation{Name: "delete", Args: &DeleteArgs{0}})
	assert.Equal(t, "unknown item 0", err.Error())
}

func TestStateUpdateError(t *testing.T) {
	var dao = NewDao(":memory:")
	var state = NewState(dao)
	err := state.Apply(&Mutation{Name: "update", Args: &UpdateArgs{0, "name1", "json1"}})
	assert.Equal(t, "unknown item 0", err.Error())
}
