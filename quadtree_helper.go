package quadtree

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// treeWithRandomPattern returns a tree with specified level and intialized with live cells where the corresponding bit in randomNumber is set.
func treeWithRandomPattern(level uint) (qt *Quadtree, randomNumber *big.Int) {
	qt = EmptyTree(1)
	for i := uint(1); i < level; i++ {
		qt = qt.grow()
	}
	edgeLength := Dim(1) << level // level = 3 => 8, upperBound = 64
	cellsInTree := uint(edgeLength * edgeLength)

	// random big Int were each bit corresponds to one cell
	upperBound := new(big.Int)
	upperBound.SetInt64(1)
	upperBound.Lsh(upperBound, cellsInTree-1)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber = new(big.Int).Rand(r, upperBound)

	for x := Dim(0); x < edgeLength; x++ {
		for y := Dim(0); y < edgeLength; y++ {
			bitPosition := x*edgeLength + y
			ux := x - edgeLength/2
			uy := y - edgeLength/2

			if randomNumber.Bit(int(bitPosition)) != 0 {
				qt = qt.SetCell(ux, uy, 1)
			}
		}
	}

	return qt, randomNumber
}

func (qt *Quadtree) FillTreeWithRandomPattern(start Dim, end Dim) *Quadtree {
	edgeLength := end - start
	cellsInTree := uint(edgeLength * edgeLength)

	// random big Int were each bit corresponds to one cell
	upperBound := new(big.Int)
	upperBound.SetInt64(1)
	upperBound.Lsh(upperBound, cellsInTree-1)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := new(big.Int).Rand(r, upperBound)

	for x := Dim(0); x < edgeLength; x++ {
		for y := Dim(0); y < edgeLength; y++ {
			bitPosition := x*edgeLength + y
			ux := x + start
			uy := y + start

			if randomNumber.Bit(int(bitPosition)) != 0 {
				qt = qt.SetCell(ux, uy, 1)
			} else {
				qt = qt.SetCell(ux, uy, 0)
			}
		}
	}
	return qt
}

//assertRandomPattern asserts that the tree has live cells were the corresponding bit position in randomNumber is set
func (qt *Quadtree) assertRandomPattern(t *testing.T, randomNumber *big.Int) {
	edgeLength := Dim(1) << qt.Level
	for x := Dim(0); x < edgeLength; x++ {
		for y := Dim(0); y < edgeLength; y++ {
			bitPosition := x*edgeLength + y
			ux := x - edgeLength/2
			uy := y - edgeLength/2
			//fmt.Println("Get", bitPosition, "random", randomNumber.Bit(int(bitPosition)), "tree", qt.cell(ux, uy))
			assert.Equal(t, randomNumber.Bit(int(bitPosition)), uint(qt.Cell(ux, uy)), "at position %d", bitPosition)
		}
	}
}
