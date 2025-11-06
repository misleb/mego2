package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind_GeneratesCorrectSQL_WithJoinAndWhere(t *testing.T) {
	model := &TestModel{}
	join := &TestJoinModel{}
	find := Find(model).Join(join).Where(AnyMap{"name = $1": "blah"})

	expectedSQL := "SELECT test_table.id,test_table.name,test_table.age FROM test_table LEFT JOIN test_join_table ON test_table.id = test_join_table.test_id WHERE name = $1"
	actualSQL := find.SQL()
	assert.Equal(t, expectedSQL, actualSQL)
}

func TestFind_GeneratesCorrectSQL_WithJoin(t *testing.T) {
	model := &TestModel{}
	join := &TestJoinModel{}
	find := Find(model).Join(join)

	expectedSQL := "SELECT test_table.id,test_table.name,test_table.age FROM test_table LEFT JOIN test_join_table ON test_table.id = test_join_table.test_id"
	actualSQL := find.SQL()
	assert.Equal(t, expectedSQL, actualSQL)
}

func TestFind_GeneratesCorrectSQL_WithoutWhere(t *testing.T) {
	model := &TestModel{}
	find := Find(model)

	expectedSQL := "SELECT test_table.id,test_table.name,test_table.age FROM test_table"
	actualSQL := find.SQL()

	assert.Equal(t, expectedSQL, actualSQL)
}

func TestFind_GeneratesCorrectSQL_WithWhere(t *testing.T) {
	model := &TestModel{}
	find := Find(model)
	find.Where(AnyMap{"name = $1": "test"})

	expectedSQL := "SELECT test_table.id,test_table.name,test_table.age FROM test_table WHERE name = $1"
	actualSQL := find.SQL()

	assert.Equal(t, expectedSQL, actualSQL)
}

func TestFind_GeneratesCorrectSQL_WithMultipleWhere(t *testing.T) {
	model := &TestModel{}
	find := Find(model)
	find.Where(AnyMap{
		"name = ?": "test",
		"age > ?":  18,
	})

	actualSQL := find.SQL()

	assert.Equal(t, find.args, []any{"test", 18})
	assert.Contains(t, actualSQL, "SELECT test_table.id,test_table.name,test_table.age FROM test_table WHERE")
	assert.Contains(t, actualSQL, "age > $2")
	assert.Contains(t, actualSQL, "name = $1")
}

func TestFind_GeneratesCorrectSQL_WithEmptyWhere(t *testing.T) {
	model := &TestModel{}
	find := Find(model)
	find.Where(AnyMap{})

	expectedSQL := "SELECT test_table.id,test_table.name,test_table.age FROM test_table"
	actualSQL := find.SQL()

	assert.Equal(t, expectedSQL, actualSQL)
}

func TestFind_GeneratesCorrectSQL_ChainedWhere(t *testing.T) {
	model := &TestModel{}
	find := Find(model)
	find.Where(AnyMap{"name = $1": "test"})
	find.Where(AnyMap{"age > $2": 18})

	actualSQL := find.SQL()

	// Check the base structure
	assert.Contains(t, actualSQL, "SELECT test_table.id,test_table.name,test_table.age FROM test_table")
	assert.Contains(t, actualSQL, "WHERE")
	assert.Contains(t, actualSQL, "name = $1")
	assert.Contains(t, actualSQL, "age > $2")
	assert.Contains(t, actualSQL, " AND ")

	// Verify it starts correctly
	assert.True(t, len(actualSQL) > 0)
	assert.Contains(t, actualSQL, "SELECT test_table.id,test_table.name,test_table.age FROM test_table WHERE")
}
