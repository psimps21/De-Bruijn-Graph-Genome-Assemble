package main

import (
	"fmt"
	"os"
	"strings"
)

// DebruinizeFile Returns two De Bruijn graphs made from a given fastq file, one for the reads and another for their reverse complements
func DebruinizeFile(filename string, l int, save string) (*Graph, *Graph, *PathSet) {
	fwReads, revReads := ReadFastq(filename)
	fwReadLTups, revReadLTups := GenerateForwardRevLTuples(fwReads, revReads, l, save)
	G_fw, G_rev := MakeDeBruijnGraph(fwReadLTups), MakeDeBruijnGraph(revReadLTups)
	fwPathSet := GenerateReadPathSet(fwReads, l)
	return G_fw, G_rev, fwPathSet
}

func main() {
	var file, test string
	if len(os.Args) == 1 {
		file = "1"
	} else {
		file = os.Args[1]
	}

	if file == "0" {
		test = "small_test.fastq"
	} else {
		test = "small_test_2.fastq"
	}
	G, _, fwPathSet := DebruinizeFile(test, 3, "")

	startNode := G.FindStartNode()
	P := G.FindEulerianPath(startNode, []*Node{})
	str := []string{P[len(P)-1].value}
	for i := len(P) - 2; i >= 0; i-- {
		str = append(str, string(P[i].value[len(P[i].value)-1]))
	}
	fmt.Println("Eulerian Walk:", strings.Join(str, ""), "\n")

	fmt.Println("Original Read Path Set")
	for i, path := range *fwPathSet {
		fmt.Printf("Read Path %d Nodes: ", i)
		PrintReadPathNodes(path)
		fmt.Printf("Read Path %d Sequence: ", i)
		PrintReadPathEdges(path)
		fmt.Println()
	}

	G.SetInOutDegree()

	_, redFwPathSet := ReducePaths(G, fwPathSet)

	fmt.Println()

	fmt.Println("Reduced Read Path Set")
	for i, path := range *redFwPathSet {
		fmt.Printf("Read Path %d Nodes: ", i)
		PrintReadPathNodes(path)
		fmt.Printf("Read Path %d sequence: ", i)
		PrintReadPathEdges(path)
		fmt.Println()
	}
}
