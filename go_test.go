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
			[]byte(fileA),
			[]byte(fileB),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 3, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 3, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 3, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "other", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{ // 2
			[]byte(fileB),
			[]byte(fileBRename),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 0, DelLines: 3, Inserted: false, Deleted: true, IsOther: false},
				&Change{Name: "func xxxx", InsLines: 3, DelLines: 0, Inserted: true, Deleted: false, IsOther: false},
			},
			error(nil),
		},
		{ // 3
			[]byte(fileB),
			[]byte(fileBDelete),
			[]*Change{
				&Change{Name: "func bbbb", InsLines: 0, DelLines: 3, Inserted: false, Deleted: true, IsOther: false},
				&Change{Name: "other", InsLines: 0, DelLines: 1, Inserted: false, Deleted: false, IsOther: true},
			},
			error(nil),
		},
		{ // 4
			[]byte(fileA),
			[]byte(fileBRename),
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 0, DelLines: 1, Inserted: false, Deleted: true, IsOther: false},
				&Change{Name: "func bbbb", InsLines: 3, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func main", InsLines: 3, DelLines: 1, Inserted: false, Deleted: false, IsOther: false},
				&Change{Name: "func xxxx", InsLines: 3, DelLines: 0, Inserted: true, Deleted: false, IsOther: false},
				&Change{Name: "other", InsLines: 2, DelLines: 1, Inserted: false, Deleted: false, IsOther: true},
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

const fileA = `package main
//comment
func main() {}
func aaaa() {}
func bbbb() {}`

const fileB = `package main
func main() {
	return
}
//comment
func aaaa() {
	return
}
//comment
func bbbb() {
	return
}`

const fileBRename = `package main
func main() {
	return
}
//comment
func xxxx() {
	return
}
//comment
func bbbb() {
	return
}`

const fileBDelete = `package main
func main() {
	return
}
//comment
func aaaa() {
	return
}`
