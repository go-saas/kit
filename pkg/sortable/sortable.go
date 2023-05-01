package sortable

import (
	"golang.org/x/exp/slices"

	"sort"
)

type Sortable interface {
	GetSort() int
	SetSort(int)
}

// Embed implements Sortable interface. resource which supports sorting like "move up", "move to" can embed this struct
type Embed struct {
	Sort int
}

var _ Sortable = (*Embed)(nil)

func (e *Embed) GetSort() int {
	return e.Sort
}

func (e *Embed) SetSort(i int) {
	e.Sort = i
}

// Normalize set sort by index
func Normalize[T Sortable](l []T) {
	if len(l) == 0 {
		return
	}
	for i, sortable := range l {
		sortable.SetSort(i)
	}
}

func Sort[T Sortable](s []T) []T {
	sort.Slice(s, func(i, j int) bool {
		return s[i].GetSort() < s[j].GetSort()
	})
	return s
}

func MoveFirst(opSort int, l *[]Sortable) {
	s, valid := validateAndSort(opSort, l)
	if !valid {
		return
	}

	opIndex := slices.IndexFunc(s, func(sortable Sortable) bool {
		return sortable.GetSort() == opSort
	})
	if opIndex < 0 {
		return
	}
	s = move(s, opIndex, 0)
	Normalize(s)
	*l = s
}

func MoveLast(opSort int, l *[]Sortable) {
	s, valid := validateAndSort(opSort, l)
	if !valid {
		return
	}

	opIndex := slices.IndexFunc(s, func(sortable Sortable) bool {
		return sortable.GetSort() == opSort
	})
	if opIndex < 0 {
		return
	}
	s = move(s, opIndex, len(s)-1)

	Normalize(s)
	*l = s
}

func MoveBefore(opSort, targetSort int, l *[]Sortable) {
	s, valid := validateAndSort(opSort, l)
	if !valid {
		return
	}
	opIndex := slices.IndexFunc(s, func(sortable Sortable) bool {
		return sortable.GetSort() == opSort
	})
	if opIndex < 0 {
		return
	}
	targetIndex := slices.IndexFunc(s, func(sortable Sortable) bool {
		return sortable.GetSort() == targetSort
	})
	if targetIndex < 0 {
		return
	}

	s = moveSafe(s, opIndex, targetIndex-1)

	Normalize(s)
	*l = s
}

func MoveAfter(opSort, targetSort int, l *[]Sortable) {
	s, valid := validateAndSort(opSort, l)
	if !valid {
		return
	}
	opIndex := slices.IndexFunc(s, func(sortable Sortable) bool {
		return sortable.GetSort() == opSort
	})
	if opIndex < 0 {
		return
	}
	targetIndex := slices.IndexFunc(s, func(sortable Sortable) bool {
		return sortable.GetSort() == targetSort
	})
	if targetIndex < 0 {
		return
	}

	s = moveSafe(s, opIndex, targetIndex+1)

	Normalize(s)
	*l = s
}

func validateAndSort(opSort int, l *[]Sortable) ([]Sortable, bool) {
	s := *l
	if len(s) == 0 || opSort < 0 || opSort >= len(s) {
		return s, false
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].GetSort() < s[j].GetSort()
	})
	return s, true
}

func insert[T any](array []T, value T, index int) []T {
	return append(array[:index], append([]T{value}, array[index:]...)...)
}

func remove[T any](array []T, index int) []T {
	return append(array[:index], array[index+1:]...)
}

func move[T any](array []T, srcIndex int, dstIndex int) []T {
	value := array[srcIndex]
	return insert(remove(array, srcIndex), value, dstIndex)
}

func moveSafe[T any](array []T, srcIndex int, dstIndex int) []T {
	if dstIndex < 0 {
		dstIndex = 0
	}
	if dstIndex >= len(array) {
		dstIndex = len(array) - 1
	}
	return move(array, srcIndex, dstIndex)
}
