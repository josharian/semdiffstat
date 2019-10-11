package semdiffstat

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGo(t *testing.T) {
	tests := []struct {
		fileA   []byte
		fileB   []byte
		changes []*Change
		err     error
	}{
		{
			[]byte(baseFile),
			[]byte(moveMain),
			[]*Change{},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(moveAndChangeMain),
			[]*Change{
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(moveAll),
			[]*Change{},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(moveAndChangeTwo),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(changeAll),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(moveAndChangeAll),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(moveMainAndAddLine),
			[]*Change{
				&Change{Name: "other", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(moveMainAndAdd3Lines),
			[]*Change{
				&Change{Name: "other", InsLines: 3, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(baseFileMultiLine),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(baseFileMultiLineMoveMain),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{
			[]byte(baseFileMultiLine),
			[]byte(baseFileMultiLineMoveMain),
			[]*Change{},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(baseFileMultiLineMoveMainAddLine),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "other", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{
			[]byte(baseFile),
			[]byte(baseFileMultiLineMoveMainAddLineAndReturn),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 3, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "other", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("TEST %d/%d", i+1, len(tests)), func(t *testing.T) {
			changes, err := Go(test.fileA, test.fileB)
			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}
			lenc := len(changes)
			lent := len(test.changes)
			if lenc != lent {
				t.Fatalf("expected %d changes, got %d changes", lent, lenc)
			}
			for i, c := range changes {
				if !reflect.DeepEqual(*c, *test.changes[i]) {
					t.Fatalf("expected %v, got %v", *test.changes[i], *c)
				}
			}
		})
	}
}

//------------------------------------------------------
// test files

const baseFile = `package main
func main() {}
func aaaa() {}
func bbbb() {}`

const moveMain = `package main
func aaaa() {}
func main() {}
func bbbb() {}`

const moveAndChangeMain = `package main
func aaaa() {}
func main() { }
func bbbb() {}`

const moveAll = `package main
func aaaa() {}
func bbbb() {}
func main() {}`

const moveAndChangeTwo = `package main
func aaaa() { }
func bbbb() {}
func main() { }`

const changeAll = `package main
func main() { }
func aaaa() { }
func bbbb() { }`

const moveAndChangeAll = `package main
func aaaa() { }
func bbbb() { }
func main() { }`

const moveMainAndAddLine = `package main
func aaaa() {}

func main() {}
func bbbb() {}`

const moveMainAndAdd3Lines = `package main

func aaaa() {}

func main() {}

func bbbb() {}`

const baseFileMultiLine = `package main
func main() {
}
func aaaa() {
}
func bbbb() {
}`

const baseFileMultiLineMoveMain = `package main
func aaaa() {
}
func main() {
}
func bbbb() {
}`

const baseFileMultiLineMoveMainAddLine = `package main

func aaaa() {
}
func main() {
}
func bbbb() {
}`

const baseFileMultiLineMoveMainAddLineAndReturn = `package main

func aaaa() {
}
func main() {
	return
}
func bbbb() {
}`
