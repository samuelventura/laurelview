package lvnbe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaoCrud(t *testing.T) {
	var dao = NewDao(":memory:")

	assert.Equal(t, 0, len(dao.All()))

	item1 := dao.Create("name1", "json1")
	assert.Equal(t, "name1", item1.Name)
	assert.Equal(t, "json1", item1.Json)
	assert.Equal(t, 1, len(dao.All()))

	item2 := dao.Create("name2", "json2")
	assert.Equal(t, "name2", item2.Name)
	assert.Equal(t, "json2", item2.Json)
	assert.Equal(t, 2, len(dao.All()))

	item3 := dao.Update(item2.ID, "name3", "json3")
	assert.Equal(t, "json3", item3.Json)

	all := dao.All()
	assert.Equal(t, "json1", all[0].Json)
	assert.Equal(t, "name3", all[1].Name)
	assert.Equal(t, "json3", all[1].Json)

	dao.Delete(item1.ID)
	assert.Equal(t, 1, len(dao.All()))

	dao.Delete(item2.ID)
	assert.Equal(t, 0, len(dao.All()))
}
