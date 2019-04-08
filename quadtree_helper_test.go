package quadtree

import (
	"testing"
)

/*
 * test for random pattern helper funcitons
 */
func TestRandomPattern(t *testing.T) {
	qt, randomNumber := treeWithRandomPattern(5)
	treeCorrectness(t, qt)
	qt.assertRandomPattern(t, randomNumber)
}
