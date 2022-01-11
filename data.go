package main

//MakeDeBruijnGraph returns a de bruijn graph from a given list of l-tuples
func MakeDeBruijnGraph(lTupMap []string) *Graph {
	deBruijn := NewGraph()
	// For each l-tuple add an edge and nodes to a graph
	for _, lTup := range lTupMap {
		u, v := &Node{lTup[:len(lTup)-1]}, &Node{lTup[1:]}
		e := &Edge{start: u, end: v, value: lTup}
		deBruijn.AddEdge(e)
	}
	return deBruijn
}
