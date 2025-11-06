package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertSQL(t *testing.T) {
	user := &TestModel{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}
	insert := Insert(user)
	assert.Equal(t, "INSERT INTO test_table (name,age,password) VALUES ($1,$2,crypt($3, gen_salt('bf'))) RETURNING id", insert.SQL())
}
