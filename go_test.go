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
		{ // 1
			[]byte(emptyFile),
			[]byte(baseFile),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 0, Inserted: true, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 0, Inserted: true, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 0, Inserted: true, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 2
			[]byte(baseFile),
			[]byte(emptyFile),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 0, DelLines: 1, Inserted: false, Deleted: true, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 0, DelLines: 1, Inserted: false, Deleted: true, IsOther: false},
				&Change{Name: "func main", InsLines: 0, DelLines: 1, Inserted: false, Deleted: true, IsOther: false},
			},
			error(nil),
		},
		{ // 3
			[]byte(baseFile),
			[]byte(moveMain),
			[]*Change{},
			error(nil),
		},
		{ // 4
			[]byte(baseFile),
			[]byte(moveMainDeleteBbbb),
			[]*Change{
				&Change{Name: "func bbbb", InsLines: 0, DelLines: 1, Inserted: false, Deleted: true, IsOther: false},
			},
			error(nil),
		},
		{ // 5
			[]byte(baseFile),
			[]byte(moveMainDeleteBbbbAddLine),
			[]*Change{
				&Change{Name: "func bbbb", InsLines: 0, DelLines: 1, Inserted: false, Deleted: true, IsOther: false},
				&Change{Name: "other", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{ // 6
			[]byte(baseFile),
			[]byte(moveAndChangeMain),
			[]*Change{
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 7
			[]byte(baseFile),
			[]byte(moveAll),
			[]*Change{},
			error(nil),
		},
		{ // 8
			[]byte(baseFile),
			[]byte(moveAndChangeTwo),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 9
			[]byte(baseFile),
			[]byte(changeAll),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 10
			[]byte(baseFile),
			[]byte(moveAndChangeAll),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 11
			[]byte(baseFile),
			[]byte(moveMainAndAddLine),
			[]*Change{
				&Change{Name: "other", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{ // 12
			[]byte(baseFile),
			[]byte(moveMainAndAdd3Lines),
			[]*Change{
				&Change{Name: "other", InsLines: 3, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{ // 13
			[]byte(baseFile),
			[]byte(baseFileMultiLine),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 14
			[]byte(baseFile),
			[]byte(baseFileMultiLineMoveMain),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 15
			[]byte(baseFileMultiLine),
			[]byte(baseFileMultiLineMoveMain),
			[]*Change{},
			error(nil),
		},
		{ // 16
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
		{ // 17
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
		{ // 18
			[]byte(baseFileMultiLine),
			[]byte(baseFileMultiLineMoveMainAddLineAndReturn),
			[]*Change{
				&Change{Name: "func main", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "other", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{ // 19
			[]byte(baseFileMultiLine),
			[]byte(baseFileMultiLineMoveMainAddLineAndReturnAll),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 1, DelLines: 0, Inserted: false, Deleted: false, IsOther: false},
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

const emptyFile = "package main"

const baseFile = `package main
func main() {}
func aaaa() {}
func bbbb() {}`

const moveMain = `package main
func aaaa() {}
func main() {}
func bbbb() {}`

const moveMainDeleteBbbb = `package main
func aaaa() {}
func main() {}`

const moveMainDeleteBbbbAddLine = `package main
func aaaa() {}

func main() {}`

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

const baseFileMultiLineMoveMainAddLineAndReturnAll = `package main

func aaaa() {
	return
}
func main() {
	return
}
func bbbb() {
	return
}`
