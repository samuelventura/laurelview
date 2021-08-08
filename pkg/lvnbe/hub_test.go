package lvnbe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHubAddDuplicatedSid(t *testing.T) {
	var err error
	to := newTestOutput()
	dao := NewDao(":memory:")
	state := NewState(dao)
	hub := NewHub(state)
	err = hub.Add("0", func(mut *Mutation) { to.out(mut.Name) })
	PanicIfError(err)
	to.matchNext(t, "all")
	err = hub.Add("0", NopCallback)
	assert.Equal(t, "duplicated sid 0", err.Error())
	to.assertEmpty(t)
}

func TestHubRemoveNonExistentSid(t *testing.T) {
	var err error
	to := newTestOutput()
	dao := NewDao(":memory:")
	state := NewState(dao)
	hub := NewHub(state)
	err = hub.Remove("0")
	assert.Equal(t, "non-existent sid 0", err.Error())
	to.assertEmpty(t)
}

func TestHubLifeCycle(t *testing.T) {
	var err error
	to := newTestOutput()
	dao := NewDao(":memory:")
	dao.Create("name1", "json1")
	state := NewState(dao)
	hub := NewHub(state)

	cbg := func(cid string) Callback {
		return func(mut *Mutation) {
			to.out(cid, mut.Name)
			switch mut.Name {
			case "all":
				args := mut.Args.(*AllArgs)
				assert.Equal(t, uint(1), args.Items[0].Id)
				assert.Equal(t, "json1", args.Items[0].Json)
			case "create":
				args := mut.Args.(*CreateArgs)
				assert.Equal(t, uint(2), args.Id)
				assert.Equal(t, "name2", args.Name)
				assert.Equal(t, "json2", args.Json)
			case "update":
				args := mut.Args.(*UpdateArgs)
				assert.Equal(t, uint(2), args.Id)
				assert.Equal(t, "name3", args.Name)
				assert.Equal(t, "json3", args.Json)
			case "delete":
				args := mut.Args.(*DeleteArgs)
				assert.Equal(t, uint(2), args.Id)
			}
		}
	}

	err = hub.Add("0", cbg("callback0"))
	PanicIfError(err)
	to.matchNext(t, "callback0", "all")
	to.assertEmpty(t)
	hub.Add("1", cbg("callback1"))
	to.matchNext(t, "callback1", "all")
	to.assertEmpty(t)

	create := &CreateArgs{0, "name2", "json2"}
	err = hub.Apply(&Mutation{Name: "create", Args: create})
	PanicIfError(err)
	to.matchNext(t, "callback0", "create")
	to.matchNext(t, "callback1", "create")
	to.assertEmpty(t)

	update := &UpdateArgs{2, "name3", "json3"}
	err = hub.Apply(&Mutation{Name: "update", Args: update})
	PanicIfError(err)
	to.matchNext(t, "callback0", "update")
	to.matchNext(t, "callback1", "update")
	to.assertEmpty(t)

	delete := &DeleteArgs{2}
	err = hub.Apply(&Mutation{Name: "delete", Args: delete})
	PanicIfError(err)
	to.matchNext(t, "callback0", "delete")
	to.matchNext(t, "callback1", "delete")
	to.assertEmpty(t)

	err = hub.Apply(&Mutation{Name: "delete", Args: delete})
	assert.Equal(t, "unknown item 2", err.Error())
	to.assertEmpty(t)

	err = hub.Apply(&Mutation{Name: "update", Args: update})
	assert.Equal(t, "unknown item 2", err.Error())
	to.assertEmpty(t)
}
