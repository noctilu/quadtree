package quadtree_test

import (
	"github/noctilu/quadtree"
)

func Example() {
	// empty tree with level 2
	qt := quadtree.EmptyTree(2)

	// quadtrees are immutable, so each change gives you a new quadtree
	qt = qt.SetCell(1, 1, 1)

	// print tree to console, don't do that for bigger trees
	qt.Print()

	// calculates next generation
	qtNext := qt.NextGen()

	qtNext.Print()
}
