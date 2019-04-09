[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/noctilu/quadtree)
[![GoDoc](https://godoc.org/github.com/noctilu/quadtree?status.svg)](https://godoc.org/github.com/noctilu/quadtree)
# quadtree
Go implementation of a quadtree with Game of Life's hashlife algorithm.

## Usage
```go
// empty tree with level 2
qt := quadtree.EmptyTree(2)

// quadtrees are immutable, so each change gives you a new quadtree
qt = qt.SetCell(1, 1, 1)

// print tree to console, don't do that for bigger trees
qt.Print()

// calculates next generation
qtNext := qt.NextGen()

```






