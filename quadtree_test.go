package quadtree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyTree(t *testing.T) {
	qt := EmptyTree(0)
	assert.Equal(t, uint(0), qt.Level)

	qt = EmptyTree(qt.Level - 1)
	assert.Equal(t, uint(0), qt.Level)

	qt = EmptyTree(7)
	treeCorrectness(t, qt)
}

func TestGrowToFit(t *testing.T) {
	qt := EmptyTree(1)
	qt = qt.GrowToFit(63, 63)

	assert.Equal(t, uint(7), qt.Level)
	treeCorrectness(t, qt)

}

func TestSetCellPanic(t *testing.T) {
	qt := EmptyTree(1)
	qt = qt.GrowToFit(3, 3)

	assert.Panics(t, func() { qt.SetCell(8, 8, 1) })
}

func TestSetCell(t *testing.T) {
	qt := EmptyTree(1)
	for counter := 0; counter < 10; counter++ {
		//x := dim(uint64(rand.Uint32())<<1 | uint64(rand.Uint32()))
		//y := dim(uint64(rand.Uint32())<<1 | uint64(rand.Uint32()))
		x := Dim((counter - 5) * 3)
		y := Dim((counter - 5) * counter)
		qt = qt.GrowToFit(x, y)
		qt = qt.SetCell(x, y, 1)
		assert.Equal(t, Dim(1), qt.Cell(x, y))
		qt = qt.SetCell(x, y, 0)
		assert.Equal(t, Dim(0), qt.Cell(x, y))
	}

	// check that not all cells get set
	qt = qt.SetCell(1, 1, 1)
	assert.Equal(t, Dim(0), qt.Cell(2, 2))
}

func TestCell(t *testing.T) {
	qt := EmptyTree(1)
	qt = qt.GrowToFit(55, 233)
	assert.Equal(t, Dim(0), qt.Cell(55, 233))
	qt = qt.SetCell(55, 233, 1)
	assert.Equal(t, Dim(1), qt.Cell(55, 233))
	treeCorrectness(t, qt)
}

func TestFindLifeCells(t *testing.T) {
	qt := EmptyTree(1)
	qt = qt.GrowToFit(55, 233)
	qt = qt.SetCell(55, 232, 1)
	qt = qt.SetCell(55, 233, 1)
	qt.FindLifeCells(-(1 << (qt.Level - 1)), -(1 << (qt.Level - 1)), func(x, y Dim) { fmt.Println(x, y) })
}

func TestOneGen(t *testing.T) {
	// dying overpopulation
	var bitmask uint16 = 0xFFFF
	assert.Equal(t, int64(0), oneGen(bitmask).Population)

	// liveless
	bitmask = 0x0000
	assert.Equal(t, int64(0), oneGen(bitmask).Population)

	// 3 live neighbours
	// 0b0111 0000 0000
	bitmask = 0x0700
	assert.Equal(t, int64(1), oneGen(bitmask).Population)

	// 2 live neighbours and self is live
	// 0b0011 0010 0000
	bitmask = 0x0320
	assert.Equal(t, int64(1), oneGen(bitmask).Population)

	// 1 live neighbours and self is live
	// 0b0010 0010 0000
	bitmask = 0x0220
	assert.Equal(t, int64(0), oneGen(bitmask).Population)

	// 3 live neighbours below
	// 0b0000 0000 0111
	bitmask = 0x0007
	assert.Equal(t, int64(1), oneGen(bitmask).Population)
}

func TestCenteredSubnode(t *testing.T) {
	qt := EmptyTree(3) //(-4,3)
	qt.SetCell(1, 1, 1)
	qt.SetCell(-1, -1, 1)
	center := qt.centeredSubnode()
	center = center.grow()
	assert.Equal(t, qt, center)
}

func TestCenteredHorizontal(t *testing.T) {
	w := EmptyTree(2)
	e := EmptyTree(2)
	w.SetCell(1, -1, 1)
	e.SetCell(-2, 0, 1)
	centerH := centeredHorizontal(w, e)

	expect := backslashLevelOne()
	assert.Equal(t, expect, centerH)

	w = EmptyTree(2)
	e = EmptyTree(2)
	w.SetCell(1, -0, 1)
	e.SetCell(-2, -1, 1)
	centerH = centeredHorizontal(w, e)

	expect = slashLevelOne()
	assert.Equal(t, expect, centerH)
}

func TestCenteredVertical(t *testing.T) {
	n := EmptyTree(2)
	s := EmptyTree(2)
	n.SetCell(-1, 1, 1)
	s.SetCell(0, -2, 1)
	centerV := centeredVertical(n, s)

	expect := backslashLevelOne()
	assert.Equal(t, expect, centerV)

	n = EmptyTree(2)
	s = EmptyTree(2)
	n.SetCell(0, 1, 1)
	s.SetCell(-1, -2, 1)
	centerV = centeredVertical(n, s)

	expect = slashLevelOne()
	assert.Equal(t, expect, centerV)
}

func TestCenteredSubSubnode(t *testing.T) {
	qt, _ := treeWithRandomPattern(1)
	grown := qt.grow().grow()
	centeredSubSubnode := grown.centeredSubSubnode()

	assert.Equal(t, qt, centeredSubSubnode)
}

func TestSlowSimulation(t *testing.T) {
	qt := EmptyTree(2)

	// empty stays empty
	emptyResult := qt.slowSimulation()
	assert.Equal(t, EmptyTree(1), emptyResult)

	// 1 | 1
	// 0 | 1
	qt = EmptyTree(2)
	qt.SetCell(-1, -1, 1)
	qt.SetCell(0, -1, 1)
	qt.SetCell(0, 0, 1)

	fullResult := qt.slowSimulation()
	expect := EmptyTree(1)
	expect.SetCell(0, 0, 1)
	expect.SetCell(-1, 0, 1)
	expect.SetCell(-1, -1, 1)
	expect.SetCell(0, -1, 1)
	assert.Equal(t, expect, fullResult)

	// next genartion should be full as well
	fullResult = fullResult.grow().slowSimulation()
	assert.Equal(t, expect, fullResult)

	// 1 | 1| 1| 1
	// 1 | 1| 1| 1
	// 1 | 1| 1| 1
	// 1 | 1| 1| 1
	qt = EmptyTree(2)
	for x := Dim(-2); x < 2; x++ {
		for y := Dim(-2); y < 2; y++ {
			qt.SetCell(x, y, 1)
		}
	}
	emptyResult2 := qt.slowSimulation()
	assert.Equal(t, EmptyTree(1), emptyResult2)
}

// trivial case of empty tree
// more testing should happen on universe level
func TestNextGeneration(t *testing.T) {
	qt := EmptyTree(4)
	qt = qt.grow()
	qtNext := qt.NextGeneration()
	qtNext = qtNext.grow()
	assert.Equal(t, qt, qtNext)
}

func TestString(t *testing.T) {
	qt, _ := treeWithRandomPattern(3)
	fmt.Sprint(qt)
}

/*
 * Benchmarks
 */
var result Dim

func benchmarkAddAndReadCells(size Dim, b *testing.B) {
	qt := EmptyTree(1)
	qt = qt.GrowToFit(Dim(size), Dim(size))
	//b.ResetTimer()
	for n := 0; n < b.N; n++ {
		qt.SetCell(2, 2, 1)
		result = qt.Cell(2, 2)
	}
}

func BenchmarkAddAndReadCells3(b *testing.B)  { benchmarkAddAndReadCells(Dim(1)<<3, b) }
func BenchmarkAddAndReadCells16(b *testing.B) { benchmarkAddAndReadCells(Dim(1)<<16, b) }
func BenchmarkAddAndReadCells32(b *testing.B) { benchmarkAddAndReadCells(Dim(1)<<32, b) }

func benchmarkGrowToFit(size Dim, b *testing.B) {
	for n := 0; n < b.N; n++ {
		qt := EmptyTree(1)
		qt = qt.GrowToFit(Dim(size), Dim(size))
	}
}

func BenchmarkGrowToFit3(b *testing.B)  { benchmarkGrowToFit(Dim(1)<<3, b) }
func BenchmarkGrowToFit8(b *testing.B)  { benchmarkGrowToFit(Dim(1)<<8, b) }
func BenchmarkGrowToFit16(b *testing.B) { benchmarkGrowToFit(Dim(1)<<16, b) }
func BenchmarkGrowToFit32(b *testing.B) { benchmarkGrowToFit(Dim(1)<<32, b) }

/*
* Helper
 */
// treeCorrectness recursivly checks Level of each node and that leaf nodes have no childs
func treeCorrectness(t *testing.T, qt *Quadtree) {
	if qt.Level == 0 {
		for _, child := range qt.childs() {
			assert.Nil(t, child, "Leafe nodes shouldn't have child nodes")
		}
		return
	}

	for _, child := range qt.childs() {
		if child == nil {
			continue
		}
		assert.Equal(t, qt.Level-1, child.Level)
		treeCorrectness(t, child)
	}
}

// slashLevelOne returns a level one tree with the following pattern
// 0 | 1
// 1 | 0
func slashLevelOne() (qt *Quadtree) {
	qt = EmptyTree(1)
	qt.SetCell(0, -1, 1)
	qt.SetCell(-1, 0, 1)
	return
}

// backslashLevelOne returns a level one tree with the following pattern
// 1 | 0
// 0 | 1
func backslashLevelOne() (qt *Quadtree) {
	qt = EmptyTree(1)
	qt.SetCell(0, 0, 1)
	qt.SetCell(-1, -1, 1)
	return
}
