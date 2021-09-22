package lvndb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbStateDispose(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		disp := NewState(ctx)
		disp(Mn(":dispose"))
		to.MatchWait(t, 200, "trace", "state", "{:dispose,,<nil>,<nil>}")
		to.MatchWait(t, 200, "trace", "hub", "{:dispose,,<nil>,<nil>}")
		disp(Mn(":dispose"))
		to.MatchWait(t, 200, "debug", "state", "{:dispose,,<nil>,<nil>}")
	})
}

func TestDbStateStartupEmpty(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		disp := NewState(ctx)
		disp(Mnsa(":add", "tid", to.Dispatch("cb")))
		to.MatchWait(t, 200, "trace", "hub", "{:add,tid,")
		to.MatchWait(t, 200, "trace", "hub", "{all,tid,..lvndb.OneArgs,..}")
		disp(Mnsa("create", "tid", OneArgs{Name: "name1", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "hub", "{create,tid,lvndb.OneArgs,{1 name1 json1}}")
	})
}

func TestDbStateStartupWithData(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		dao.Create("name1", "json1")
		disp := NewState(ctx)
		disp(Mnsa(":add", "tid", to.Dispatch("cb")))
		to.MatchWait(t, 200, "trace", "hub", "{all,tid,..lvndb.OneArgs,.{1 name1 json1}.}")
	})
}

func TestDbStateCreateSuccess(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		disp := NewState(ctx)
		disp(Mnsa("create", "tid", OneArgs{Name: "name1", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "state", "{create,tid,lvndb.OneArgs,{0 name1 json1}}")
		to.MatchWait(t, 200, "trace", "hub", "{create,tid,lvndb.OneArgs,{1 name1 json1}}")
		item := dao.All()[0]
		assert.Equal(t, uint(1), item.ID)
		assert.Equal(t, "name1", item.Name)
		assert.Equal(t, "json1", item.Json)
	})
}

func TestDbStateUpdateSuccess(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		dao.Create("name1", "json1")
		disp := NewState(ctx)
		disp(Mnsa("update", "tid", OneArgs{Id: 1, Name: "name2", Json: "json2"}))
		to.MatchWait(t, 200, "trace", "state", "{update,tid,lvndb.OneArgs,{1 name2 json2}}")
		to.MatchWait(t, 200, "trace", "hub", "{update,tid,lvndb.OneArgs,{1 name2 json2}}")
		item := dao.All()[0]
		assert.Equal(t, uint(1), item.ID)
		assert.Equal(t, "name2", item.Name)
		assert.Equal(t, "json2", item.Json)
	})
}

func TestDbStateDeleteSuccess(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		dao.Create("name1", "json1")
		disp := NewState(ctx)
		disp(Mnsa("delete", "tid", uint(1)))
		to.MatchWait(t, 200, "trace", "state", "{delete,tid,uint,1}")
		to.MatchWait(t, 200, "trace", "hub", "{delete,tid,uint,1}")
		assert.Equal(t, 0, len(dao.All()))
	})
}

func TestDbStateCreateError(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		disp := NewState(ctx)
		disp(Mnsa("create", "tid", OneArgs{Name: " ", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "state", "{create,tid,lvndb.OneArgs,{0   json1}}")
		to.MatchWait(t, 200, "debug", "state", "{create,tid,lvndb.OneArgs,{0   json1}}", "name cannot be empty")
		assert.Equal(t, 0, len(dao.All()))
	})
}

func TestDbStateUpdateErrorNotFound(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		disp := NewState(ctx)
		disp(Mnsa("update", "tid", OneArgs{Id: 1, Name: "name1", Json: "json1"}))
		to.MatchWait(t, 200, "trace", "state", "{update,tid,lvndb.OneArgs,{1 name1 json1}}")
		to.MatchWait(t, 200, "debug", "state", "{update,tid,lvndb.OneArgs,{1 name1 json1}}", "item not found", "1")
	})
}

func TestDbStateUpdateErrorEmptyName(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		dao.Create("name1", "json1")
		disp := NewState(ctx)
		disp(Mnsa("update", "tid", OneArgs{Id: 1, Name: " ", Json: "json2"}))
		to.MatchWait(t, 200, "trace", "state", "{update,tid,lvndb.OneArgs,{1   json2}}")
		to.MatchWait(t, 200, "debug", "state", "{update,tid,lvndb.OneArgs,{1   json2}}", "name cannot be empty")
	})
}

func TestDbStateDeleteErrorNotFound(t *testing.T) {
	testSetupState(func(to TestOutput, ctx Context, dao Dao, log Logger) {
		disp := NewState(ctx)
		disp(Mnsa("delete", "tid", uint(1)))
		to.MatchWait(t, 200, "trace", "state", "{delete,tid,uint,1}")
		to.MatchWait(t, 200, "debug", "state", "{delete,tid,uint,1}", "item not found", "1")
	})
}

func testSetupState(callback func(to TestOutput, ctx Context, dao Dao, log Logger)) {
	var dao = NewDao(":memory:")
	defer dao.Close()
	to := NewTestOutput()
	defer to.Close()
	log := to.Logger()
	ctx := NewContext(to.Log)
	defer WaitClose(ctx.Close)
	ctx.SetValue("dao", dao)
	ctx.SetDispatch("hub", to.Dispatch("hub"))
	callback(to, ctx, dao, log)
}
