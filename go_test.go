package semdiffstat

import (
	"reflect"
	"testing"
)

func TestGo(t *testing.T) {
	tests := []struct {
		name    string
		fileA   string
		fileB   string
		changes []*Change
		err     error
	}{
		{ // 1
			"inline funcs expanded +other",
			fileA,
			fileB,
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 3, DelLines: 1},
				&Change{Name: "func bbbb", InsLines: 3, DelLines: 1},
				&Change{Name: "func main", InsLines: 3, DelLines: 1},
				&Change{Name: "other", InsLines: 2, DelLines: 1, IsOther: true},
			},
			error(nil),
		},
		{ // 2
			"func renamed +other",
			fileB,
			fileBRename,
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 0, DelLines: 3, Deleted: true},
				&Change{Name: "func xxxx", InsLines: 3, DelLines: 0, Inserted: true},
				&Change{Name: "other", InsLines: 1, DelLines: 0, IsOther: true},
			},
			error(nil),
		},
		{ // 3
			"func deleted +other",
			fileB,
			fileBDelete,
			[]*Change{
				&Change{Name: "func bbbb", InsLines: 0, DelLines: 3, Deleted: true},
				&Change{Name: "other", InsLines: 0, DelLines: 1, IsOther: true},
			},
			error(nil),
		},
		{ // 4
			"inline funcs expanded and renamed +other",
			fileA,
			fileBRename,
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 0, DelLines: 1, Deleted: true},
				&Change{Name: "func bbbb", InsLines: 3, DelLines: 1},
				&Change{Name: "func main", InsLines: 3, DelLines: 1},
				&Change{Name: "func xxxx", InsLines: 3, DelLines: 0, Inserted: true},
				&Change{Name: "other", InsLines: 3, DelLines: 1, IsOther: true},
			},
			error(nil),
		},
		{ // 5
			"insert func and inline all +other",
			fileBDelete,
			fileA,
			[]*Change{
				&Change{Name: "func aaaa", InsLines: 1, DelLines: 3},
				&Change{Name: "func bbbb", InsLines: 1, DelLines: 0, Inserted: true},
				&Change{Name: "func main", InsLines: 1, DelLines: 3},
				&Change{Name: "other", InsLines: 1, DelLines: 1, IsOther: true},
			},
			error(nil),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			changes, err := Go([]byte(test.fileA), []byte(test.fileB))
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
// main docs
func main() {}
func aaaa() {}
func bbbb() {}`

const fileB = `package main
func main() {
	return
}
// aaaa docs
func aaaa() {
	return
}
// bbbb docs
func bbbb() {
	return
}`

const fileBRename = `package main
func main() {
	return
}
// xxxx docs
func xxxx() {
	return
}
// bbbb docs
// bbbb added line
func bbbb() {
	return
}`

const fileBDelete = `package main
func main() {
	return
}
// aaaa docs
func aaaa() {
	return
}`
