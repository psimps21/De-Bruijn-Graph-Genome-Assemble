package main

import (
	"math/rand"
)

type Graph struct {
	nodes        []*Node
	edges        []*Edge
	nodeValueMap map[string]*Node
	edgeValueMap map[string]map[string]*Edge
	inDegree     map[*Node]int
	outDegree    map[*Node]int
}

type Node struct {
	value string
}

type Edge struct {
	start, end *Node
	value      string
	weight     int
	traversed  int
}

// NewGraph returns a Graph with initialized attributes
func NewGraph() *Graph {
	return &Graph{make([]*Node, 0), make([]*Edge, 0), make(map[string]*Node), make(map[string]map[string]*Edge), make(map[*Node]int), make(map[*Node]int)}
}

func GetEulerianPath(g *Graph) []*Node {
	start := g.FindStartNode()
	return g.FindEulerianPath(start, []*Node{})
}

//
/* Graph Methods */
//

//FindEulerianPath Recursively builds and returns a list of nodes corresponding to the Eulerian path in the graph
func (g *Graph) FindEulerianPath(n *Node, path []*Node) []*Node {
	n = g.NodeInGraph(n)
	g.inDegree[n]--
	// How to handle multiple of the same edges
	if g.outDegree[n] == 0 {
		return append(path, n)
	}

	childMap := g.edgeValueMap[n.value]
	// Make list of children and random order to choose them in
	childKeys, counter := make([]string, len(childMap)), 0
	for i := range childMap {
		childKeys[counter] = i
		counter++
	}
	randOrder := rand.Perm(len(childMap))

	for _, i := range randOrder {
		k, v := childKeys[i], childMap[childKeys[i]]
		c := g.GetNodeFromValue(k)
		if g.outDegree[c] > 0 || g.inDegree[c] > 0 { // v.traversed < v.weight
			v.traversed++
			g.outDegree[n]--
			path = g.FindEulerianPath(c, path)
		}
	}
	return append(path, n)
}

//SetInOutDegrees Sets the in and out degree for each node
func (g *Graph) SetInOutDegree() {
	g.inDegree = make(map[*Node]int)
	g.outDegree = make(map[*Node]int)
	// Determine in and out degree of each node
	for _, e := range g.edges {
		for i := 0; i < e.weight; i++ {
			if val, ok := g.outDegree[g.GetNodeFromValue(e.start.value)]; ok {
				g.outDegree[g.GetNodeFromValue(e.start.value)] = val + 1
			} else {
				g.outDegree[g.GetNodeFromValue(e.start.value)] = 1
			}
			if val, ok := g.inDegree[g.GetNodeFromValue(e.end.value)]; ok {
				g.inDegree[g.GetNodeFromValue(e.end.value)] = val + 1
			} else {
				g.inDegree[g.GetNodeFromValue(e.end.value)] = 1
			}
		}
	}
}

//FindStartNode returns the start node for an Eulerian walk.
// If the number of odd nodes is 2 then the star node is the node with (outdegree - indegree) == 1
// If the number of odd nodes is 0 then start on any node with a nonzero degree
func (g *Graph) FindStartNode() *Node {
	// Find odd nodes
	oddNodes := make([]*Node, 0)
	for _, k := range g.nodes {
		if g.outDegree[k]-g.inDegree[k] == 1 || g.inDegree[k]-g.outDegree[k] == 1 {
			oddNodes = append(oddNodes, k)
		}
	}

	// Return start node for eulerian walk
	if len(oddNodes) == 2 {
		for _, n := range oddNodes {
			if g.outDegree[n]-g.inDegree[n] == 1 {
				return n
			}
		}
	} else if len(oddNodes) == 0 {
		for k, v := range g.outDegree {
			if v != 0 {
				return k
			}
		}
	}

	return nil
}

// NumNodes returns the number of nodes in the graph
func (g *Graph) NumNodes() int {
	return len(g.nodes)
}

// NumEdges returns the number of edges in the graph
func (g *Graph) NumEdges() int {
	return len(g.edges)
}

// NodeInGraph returns a pointer a node if the graph contains a node with the given value. Otherwise return nil
func (g *Graph) NodeInGraph(n *Node) *Node {
	if val, ok := g.nodeValueMap[n.value]; ok {
		return val
	}
	return nil
}

// EdgeInGraph returns true if the graph contains a edge with the given value
func (g *Graph) EdgeInGraph(e *Edge) *Edge {
	if val, ok := g.edgeValueMap[e.start.value]; ok { // If start node in edge start map
		if val2, ok2 := val[e.end.value]; ok2 { // If end node in start node's end map
			return val2
		}
	}
	return nil
}

// GetEdgefromUV returns a pointer to the edge from a node with value u to a node with value v
func (g *Graph) GetEdgeFromUV(u, v string) *Edge {
	return g.edgeValueMap[u][v]
}

//GetNodeFromValue returns the address of the node with a given value
// If no node has given value returns nil
func (g *Graph) GetNodeFromValue(v string) *Node {
	if val, ok := g.nodeValueMap[v]; ok {
		return val
	}
	return nil
}

// AddEdge adds an edge to the graph. Also adds start and end nodes to to graph if they were not already present.
// If edge is already in the graph the function increments the edge weight by 1
func (g *Graph) AddEdge(e *Edge) {
	present := g.EdgeInGraph(e)
	if present == nil { // If edge is not in the graph
		g.AddNode(e.start)
		g.AddNode(e.end)
		e.weight = 1
		g.edges = append(g.edges, e)
		if endNodes, ok := g.edgeValueMap[e.start.value]; ok { // Start node is in edge map
			// if start node is in edge map but the edge isn't present, the end node must not be in the end node map
			endNodes[e.end.value] = e
		} else { // neither start node nor end node are in either edge map
			g.edgeValueMap[e.start.value] = map[string]*Edge{e.end.value: e}
		}
	} else { // If edge is already in graph
		present.weight++
	}
	g.outDegree[g.GetNodeFromValue(e.start.value)]++
	g.inDegree[g.GetNodeFromValue(e.end.value)]++
}

// AddNode adds a node to the graph.
// If node is already in graph, function does nothing
func (g *Graph) AddNode(n *Node) {
	if g.NodeInGraph(n) == nil { // Prevents different nodes with the same value in the graph
		g.nodes = append(g.nodes, n)
		g.nodeValueMap[n.value] = n
	}
}

// RemoveEdge will remove an edge from the graph.
// If an edge has weight > 1, i.e. appears multiple time in the graph, RemoveEdge will decrement the weight by 1
// If an edge has weight 1 it will be removed from the graph
func (g *Graph) RemoveEdge(e *Edge) {
	present := g.EdgeInGraph(e)
	if present != nil {
		if present.weight == 1 {
			var newEdges []*Edge
			for i, edge := range g.edges {
				if edge == present { // remove edge from list of edges
					newEdges = append(g.edges[:i], g.edges[i+1:]...)
					break
				}
			}
			g.edges = newEdges
			g.edgeValueMap[e.start.value][e.end.value] = nil // remove edge from value map
		} else if present.weight > 1 {
			e.weight--
		}
		// Decrease degrees of nodes connected to the edge
		g.outDegree[g.GetNodeFromValue(e.start.value)]--
		g.inDegree[g.GetNodeFromValue(e.end.value)]--
	}
}

// RemoveNode removes a node from the graph
func (g *Graph) RemoveNode(n *Node) {
	present := g.NodeInGraph(n)
	var newNodes []*Node
	if present != nil {
		for i, node := range g.nodes {
			if node == present { // remove node from list of nodes
				newNodes = append(g.nodes[:i], g.nodes[i+1:]...)
			}
		}

		// Remove edges that stat and end at this node

		// Remove node from node value map and set now list of nodes
		g.nodeValueMap[n.value] = nil
		g.nodes = newNodes
	}
}
