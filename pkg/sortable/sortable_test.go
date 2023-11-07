package sortable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testModel struct {
	Embed
	ID int
}

func getTestArray() []Sortable {
	return []Sortable{
		&testModel{ID: 1},
		&testModel{ID: 2},
		&testModel{ID: 3},
		&testModel{ID: 4},
		&testModel{ID: 5},
	}
}

func TestNormalize(t *testing.T) {
	a := getTestArray()
	Normalize(a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

}

func TestMoveFirst(t *testing.T) {
	a := getTestArray()
	Normalize(a)

	MoveFirst(100, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveFirst(0, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveFirst(4, &a)
	assert.Equal(t, 5, a[0].(*testModel).ID)
	assert.Equal(t, 1, a[1].(*testModel).ID)
	assert.Equal(t, 2, a[2].(*testModel).ID)
	assert.Equal(t, 3, a[3].(*testModel).ID)
	assert.Equal(t, 4, a[4].(*testModel).ID)

	MoveFirst(4, &a)
	assert.Equal(t, 4, a[0].(*testModel).ID)
	assert.Equal(t, 5, a[1].(*testModel).ID)
	assert.Equal(t, 1, a[2].(*testModel).ID)
	assert.Equal(t, 2, a[3].(*testModel).ID)
	assert.Equal(t, 3, a[4].(*testModel).ID)
}

func TestMoveLast(t *testing.T) {
	a := getTestArray()
	Normalize(a)

	MoveLast(100, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveLast(len(a)-1, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveLast(0, &a)
	assert.Equal(t, 2, a[0].(*testModel).ID)
	assert.Equal(t, 3, a[1].(*testModel).ID)
	assert.Equal(t, 4, a[2].(*testModel).ID)
	assert.Equal(t, 5, a[3].(*testModel).ID)
	assert.Equal(t, 1, a[4].(*testModel).ID)

	MoveLast(0, &a)
	assert.Equal(t, 3, a[0].(*testModel).ID)
	assert.Equal(t, 4, a[1].(*testModel).ID)
	assert.Equal(t, 5, a[2].(*testModel).ID)
	assert.Equal(t, 1, a[3].(*testModel).ID)
	assert.Equal(t, 2, a[4].(*testModel).ID)
}

func TestMoveBefore(t *testing.T) {
	a := getTestArray()
	Normalize(a)

	MoveBefore(100, 200, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveBefore(0, 1, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveBefore(1, 0, &a)
	assert.Equal(t, 2, a[0].(*testModel).ID)
	assert.Equal(t, 1, a[1].(*testModel).ID)
	assert.Equal(t, 3, a[2].(*testModel).ID)
	assert.Equal(t, 4, a[3].(*testModel).ID)
	assert.Equal(t, 5, a[4].(*testModel).ID)

	MoveBefore(1, 0, &a)
	assert.Equal(t, 1, a[0].(*testModel).ID)
	assert.Equal(t, 2, a[1].(*testModel).ID)
	assert.Equal(t, 3, a[2].(*testModel).ID)
	assert.Equal(t, 4, a[3].(*testModel).ID)
	assert.Equal(t, 5, a[4].(*testModel).ID)
}

func TestMoveAfter(t *testing.T) {
	a := getTestArray()
	Normalize(a)

	MoveAfter(100, 200, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveAfter(len(a)-1, len(a)-1, &a)
	for i, sortable := range a {
		assert.Equal(t, i, sortable.GetSort())
		assert.Equal(t, i, sortable.(*testModel).ID-1)
	}

	MoveAfter(4, 0, &a)
	assert.Equal(t, 1, a[0].(*testModel).ID)
	assert.Equal(t, 5, a[1].(*testModel).ID)
	assert.Equal(t, 2, a[2].(*testModel).ID)
	assert.Equal(t, 3, a[3].(*testModel).ID)
	assert.Equal(t, 4, a[4].(*testModel).ID)

	MoveAfter(4, 0, &a)
	assert.Equal(t, 1, a[0].(*testModel).ID)
	assert.Equal(t, 4, a[1].(*testModel).ID)
	assert.Equal(t, 5, a[2].(*testModel).ID)
	assert.Equal(t, 2, a[3].(*testModel).ID)
	assert.Equal(t, 3, a[4].(*testModel).ID)
}
