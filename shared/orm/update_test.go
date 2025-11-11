package orm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateSQL(t *testing.T) {
	model := &TestModel{}
	update := Update(model).Set([]string{"name", "age", "password"})

	expectedSQL := "UPDATE test_table SET name = :name,age = :age,password = crypt(:password, gen_salt('bf')) WHERE id = :id RETURNING id"
	assert.Equal(t, expectedSQL, update.SQL())
}

func TestUpdateQueryRequiresSet(t *testing.T) {
	model := &TestModel{}
	update := Update(model)

	err := update.Query(context.Background(), nil)
	assert.EqualError(t, err, "columns to set are required")
}
