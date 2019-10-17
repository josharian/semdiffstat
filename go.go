package semdiffstat

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"

	"github.com/pkg/diff"
)

// A Change represents a modification to some element of source code.
type Change struct {
	// Name is the name of the modified element.
	Name string

	// InsLines (DelLines) is the number of inserted (deleted) lines within the changed element.
	// For the other/unidentified change, these are only estimates.
	InsLines, DelLines int

	// Inserted (Deleted) indicates whether the element was inserted (deleted) entirely.
	Inserted, Deleted bool

	// IsOther indicates whether this change represents the catch-all other/unidentified change.
	IsOther bool
}

// Go calculates a semantic diffstat between Go source code a and b.
// It returns an error iff either a or b has parse errors.
// Identifiable changes are sorted by Name; the "other" change is always last.
func Go(a, b []byte) (changes []*Change, err error) {
	afset := token.NewFileSet()
	afile, err := parser.ParseFile(afset, "a", a, 0)
	if err != nil {
		return nil, err
	}

	bfset := token.NewFileSet()
	bfile, err := parser.ParseFile(bfset, "b", b, 0)
	if err != nil {
		return nil, err
	}

	var x bySplits
	x.asrc = a
	x.bsrc = b

	x.asplit = append(x.asplit, 0)
	for _, d := range afile.Decls {
		x.adecls = append(x.adecls, nil)
		x.asplit = append(x.asplit, afset.Position(d.Pos()).Offset)
		x.adecls = append(x.adecls, d)
		x.asplit = append(x.asplit, afset.Position(d.End()).Offset)
	}
	x.asplit = append(x.asplit, afset.Position(afile.End()).Offset)
	x.adecls = append(x.adecls, nil)

	x.bsplit = append(x.bsplit, 0)
	for _, d := range bfile.Decls {
		x.bdecls = append(x.bdecls, nil)
		x.bsplit = append(x.bsplit, bfset.Position(d.Pos()).Offset)
		x.bdecls = append(x.bdecls, d)
		x.bsplit = append(x.bsplit, bfset.Position(d.End()).Offset)
	}
	x.bsplit = append(x.bsplit, bfset.Position(bfile.End()).Offset)
	x.bdecls = append(x.bdecls, nil)

	es := diff.Myers(context.Background(), &x)

	ins := make(map[string]int)
	del := make(map[string]int)
	var other *Change
	for _, r := range es.IndexRanges {
		switch {
		case r.IsDelete():
			for i, d := range x.adecls[r.LowA:r.HighA] {
				if isFuncDecl(d) {
					del[describeFuncDecl(d)] = r.LowA + i
				} else if other == nil {
					other = &Change{Name: "other", IsOther: true}
				}
			}
		case r.IsInsert():
			for i, d := range x.bdecls[r.LowB:r.HighB] {
				if isFuncDecl(d) {
					ins[describeFuncDecl(d)] = r.LowB + i
				} else if other == nil {
					other = &Change{Name: "other", IsOther: true}
				}
			}
		}
	}

	for s, bi := range ins {
		if ai, ok := del[s]; ok {
			ds := x.diffstat(ai, bi)
			changes = append(changes, &Change{Name: s, DelLines: ds.del, InsLines: ds.ins})
		} else {
			changes = append(changes, &Change{Name: s, Inserted: true, InsLines: x.bLines(bi)})
		}
	}
	for s, ai := range del {
		if _, ok := ins[s]; !ok {
			changes = append(changes, &Change{Name: s, Deleted: true, DelLines: x.aLines(ai)})
		}
	}

	if other != nil {
		// Calculate diffstats between non-declarations.
		aIsLarger := false
		lens := len(x.asplit) - 1
		if len(x.asplit) > len(x.bsplit) {
			lens = len(x.bsplit) - 1
			aIsLarger = true
		}
		// If df is increased, it means that an index has been displaced by some change.
		var df int64
		for i := 1; i < lens; i++ {
			if int64(x.asplit[i]) != (int64(x.bsplit[i]) - df) {
				df = abs(int64(x.asplit[i] - x.bsplit[i]))
				// "other" changes may happen at even indices only.
				j := i-1
				if j%2 == 0 {
					ds := x.diffstat(j, j)
					other.InsLines += ds.ins
					other.DelLines += ds.del
				}
			}
		}
		// If x.asplit and x.bsplit are not equal in length,
		// the remaining changes must be either insertions or deletions.
		if len(x.asplit) != len(x.bsplit) {
			splitL := len(x.bsplit)
			splitF := func(i int) {
				other.InsLines += x.bLinesNonEmpty(i)
			}
			if aIsLarger {
				splitL = len(x.asplit)
				splitF = func(i int) {
					other.DelLines += x.aLinesNonEmpty(i)
				}
			}
			for i := lens - 1; i < splitL; i++ {
				if i%2 == 0 {
					splitF(i)
				}
			}
		}
		//TODO check for "other" changes between last index and EOF
		changes = append(changes, other)
	}
	sort.Slice(changes, func(i, j int) bool {
		cj := changes[j]
		// the "other" change always sorted last
		if cj.IsOther {
			return true
		}
		ci := changes[i]
		if ci.IsOther {
			return false
		}
		return ci.Name < cj.Name
	})
	return changes, nil
}

func isFuncDecl(d ast.Decl) bool {
	_, ok := d.(*ast.FuncDecl)
	return ok
}

func describeFuncDecl(d ast.Decl) string {
	fn := d.(*ast.FuncDecl)
	if fn.Recv == nil {
		return "func " + fn.Name.String()
	}
	// method
	typ := fn.Recv.List[0].Type
	ptr := ""
	if star, ok := typ.(*ast.StarExpr); ok {
		ptr = "*"
		typ = star.X
	}
	return fmt.Sprintf("func (%s%v).%s", ptr, typ.(*ast.Ident).Name, fn.Name)
}

type bySplits struct {
	asrc, bsrc     []byte
	adecls, bdecls []ast.Decl
	asplit, bsplit []int
}

func (x *bySplits) LenA() int             { return len(x.asplit) - 1 }
func (x *bySplits) LenB() int             { return len(x.bsplit) - 1 }
func (x *bySplits) Equal(ai, bi int) bool { return bytes.Equal(x.aBytes(ai), x.bBytes(bi)) }

// aBytes returns the bytes from asrc at split index ai.
func (x *bySplits) aBytes(ai int) []byte { return x.asrc[x.asplit[ai]:x.asplit[ai+1]] }
func (x *bySplits) bBytes(bi int) []byte { return x.bsrc[x.bsplit[bi]:x.bsplit[bi+1]] }

func (x *bySplits) aLines(i int) int {
	return len(bytes.Split(x.aBytes(i), newline))
}

func (x *bySplits) bLines(i int) int {
	return len(bytes.Split(x.bBytes(i), newline))
}

func (x *bySplits) aLinesNonEmpty(i int) (lines int) {
	for _, b := range bytes.Split(x.aBytes(i), newline) {
		if len(b) != 0 {
			lines++
		}
	}
	return
}

func (x *bySplits) bLinesNonEmpty(i int) (lines int) {
	for _, b := range bytes.Split(x.bBytes(i), newline) {
		if len(b) != 0 {
			lines++
		}
	}
	return
}

var newline = []byte("\n")

type diffstat struct {
	ins, del int
}

// diffstat returns the number of deleted and inserted lines
// results from a traditional line-based diff of the code segments ai and bi.
func (x *bySplits) diffstat(ai, bi int) diffstat {
	a := bytes.Split(x.aBytes(ai), newline)
	b := bytes.Split(x.bBytes(bi), newline)
	ab := diff.Bytes(a, b)
	es := diff.Myers(context.Background(), ab)
	ins, del := es.Stat()
	return diffstat{ins: ins, del: del}
}

func abs(i int64) int64 {
	j := i >> 63
	return (i ^ j) - j
}
