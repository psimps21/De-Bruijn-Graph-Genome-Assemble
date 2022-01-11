package main

import (
	"fmt"
	"strings"
)

// PathEdge is an edge in a read path
type PathEdge struct {
	value string
}

// PathNode is a node in the LinkedList for a ReadPath
type PathNode struct {
	value string
	next  *PathNode
	edge  *PathEdge
}

// ReadPath is a LinkedList of PathNodes representing a read Path in the graph G
type ReadPath struct {
	head *PathNode
	len  int // Number of nodes in path
}

// RPNode is a node in the LinkedList for the Queue data structure
type RPNode struct {
	value *ReadPath
	next  *RPNode
}

// Queue is a LinkedList implementation of the queue data structure
type Queue struct {
	head, tail *RPNode
	len        int
}

type PathSet []*ReadPath

// PrintReadPath Prints a read path
func PrintReadPathNodes(rp *ReadPath) {
	node := rp.head
	str := make([]string, 0)
	for node != nil {
		str = append(str, node.value)
		node = node.next
	}
	fmt.Println(strings.Join(str, " "))
}

func PrintReadPathEdges(rp *ReadPath) {
	node := rp.head
	str := make([]string, 0)
	for node.next != nil {
		if node == rp.head {
			str = append(str, node.edge.value)
		} else {
			str = append(str, string(node.edge.value[len(node.value):]))
		}
		node = node.next
	}
	fmt.Println(strings.Join(str, ""))
}

/*
	ReadPath Methods
*/

// NumEdges returns the number of edges in a read path
func (rp *ReadPath) NumEdges() int {
	if rp.len == 0 {
		return 0
	}
	return rp.len - 1
}

// FindXYInPath returns a list of indexes for the PathNode at the start of an XY subpath
func (rp *ReadPath) FindXYInPath(xVal, yVal string) []*PathNode {
	if rp.len < 3 {
		return nil
	}
	vInList := make([]*PathNode, 0)
	var count int
	currNode := rp.head

	for count < rp.len-2 {
		if currNode.edge.value == xVal {
			if currNode.next.edge.value == yVal {
				vInList = append(vInList, currNode)
			}
		}
		currNode = currNode.next
		count++
	}

	if len(vInList) > 0 {
		return vInList
	}
	return nil
}

//IsEndEdge returns pointer to the penultimate node if a path ends with the given edge value
func (rp *ReadPath) IsEndEdge(value string) *PathNode {
	if rp.len < 2 {
		return nil
	}
	// Find to penultimate node
	currNode := rp.head
	for currNode.next.next != nil {
		currNode = currNode.next
	}
	if currNode.edge.value == value {
		return currNode
	}
	return nil
}

//IsStartEdge returns true if a path starts with the given edge value
func (rp *ReadPath) IsStartEdge(value string) bool {
	if rp.head.edge.value == value {
		return true
	}
	return false
}

// StartYSub substitute z for y in a path starting with edge y
func (rp *ReadPath) StartYSub(zVal, startNodeVal string) {
	// Update value for first node and first edge
	rp.head.value = startNodeVal
	rp.head.edge.value = zVal
}

// EndXSub substitute edge z for x in a path ending in edge x
// xStart is a pointer to the node at the head of edge x
func (rp *ReadPath) EndXSub(zVal, endNodeVal string, xStart *PathNode) {
	// Update value for last edge and last node
	xStart.edge.value = zVal
	xStart.next.value = endNodeVal
}

// XYSub substitutes z for the edges x and y in a graph with an XY subpath
// subInx is a list indicies that start edge X
func (rp *ReadPath) XYSub(zVal string, xHeads []*PathNode) {
	// if the currNode index is the start of an X edge insert edge z and update value of neighbor node
	for _, node := range xHeads {
		node.edge.value = zVal
		node.next = node.next.next
		rp.len--
	}
}

// XYDetachPath Perform an xy-detachment on a read path
func (rp *ReadPath) XYDetachPath(x, y, z *Edge) {
	var (
		startY  bool
		endX    *PathNode
		innerXY []*PathNode
	)
	startY = rp.IsStartEdge(y.value)
	endX = rp.IsEndEdge(x.value)
	innerXY = rp.FindXYInPath(x.value, y.value)

	// Substitute z for edges x and y in path
	if startY {
		rp.StartYSub(z.value, x.start.value)
	}
	if endX != nil {
		rp.EndXSub(z.value, z.end.value, endX)
	}
	if innerXY != nil {
		rp.XYSub(z.value, innerXY)
	}
}

/*
	Queue Methods
*/

// Enqueue adds a node to the queue for a given read path and increments the length of the queue by 1
func (q *Queue) Enqueue(rp *ReadPath) {
	node := &RPNode{value: rp}
	if q.len == 0 {
		q.head, q.tail = node, node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.len++
}

// Dequeue rmoves the head of the queue and decrements the length of the queue by 1
func (q *Queue) Dequeue() *ReadPath {
	if q.len == 0 {
		return nil
	}

	path := q.head.value
	if q.len == 1 {
		q.head, q.tail = nil, nil
	} else {
		q.head = q.head.next
	}

	q.len--
	return path
}

func (q *Queue) Len() int {
	return q.len
}

/*
	PathSet Methods
*/

// PathsToDetach returns a queue conatining every ReadPath in a PathSet with more than one edge
func (ps *PathSet) PathsToReduce() *Queue {
	queue := &Queue{}
	for _, path := range *ps {
		if path.NumEdges() > 1 {
			queue.Enqueue(path)
		}
	}
	return queue
}

func (ps *PathSet) XYDetchAllPaths(x, y, z *Edge) {
	for _, path := range *ps {
		path.XYDetachPath(x, y, z)
	}
}

// ReducePaths returns a graph that has undergone x,y detachements until every read path contains one edge
func ReducePaths(g *Graph, ps *PathSet) (*Graph, *PathSet) {
	// Get queue of ReadPaths to reduce
	queue := ps.PathsToReduce()

	// Reduce paths until all paths have length 1
	for queue.Len() != 0 {
		path := queue.Dequeue()
		if path.len > 2 {
			vIn, vMid, vOut := path.head, path.head.next, path.head.next.next

			// Define x, y, and z
			x := g.GetEdgeFromUV(vIn.value, vMid.value)
			y := g.GetEdgeFromUV(vMid.value, vOut.value)
			zVal := vIn.edge.value + string(vMid.edge.value[len(vMid.value):])
			z := &Edge{start: g.GetNodeFromValue(vIn.value), end: g.GetNodeFromValue(vOut.value), value: zVal}

			// Perform xy-detchment for all paths
			ps.XYDetchAllPaths(x, y, z)

			// Remove x and y from G and add z to G
			g.AddEdge(z)
			g.RemoveEdge(x)
			g.RemoveEdge(y)

			// Remove vMid if inDegree and outDegree are 0
			vMidNode := g.GetNodeFromValue(vMid.value)
			if g.inDegree[vMidNode] == 0 && g.outDegree[vMidNode] == 0 {
				g.RemoveNode(vMidNode)
			}

			// Reduce path and queue path again if more than one edge
			if path.len > 2 {
				queue.Enqueue(path)
			}
		}
	}

	return g, ps
}
