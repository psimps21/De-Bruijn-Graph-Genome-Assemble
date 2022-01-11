package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReverseComplement returns the reverse compleiment of a string s
func ReverseComplement(s string) string {
	compMap := map[string]string{
		"A": "T",
		"C": "G",
		"G": "C",
		"T": "A",
	}

	revComp := make([]string, len(s))
	lastInx := len(s) - 1
	for i := len(s) - 1; i >= 0; i-- {
		revComp[lastInx-i] = compMap[string(s[i])]
	}

	return strings.Join(revComp, "")
}

// GenerateReadRevComps returns a list of reverese complements for a given list a reads
func GenerateReadRevComps(reads []string) []string {
	// Make reverse complement for each read
	readRevComps := make([]string, len(reads))
	for i, read := range reads {
		readRevComps[i] = ReverseComplement(read)
	}
	return readRevComps
}

// ReadFastq Returns two lists for a given fastq file, one of the reads and another of their reverse complements
func ReadFastq(filename string) ([]string, []string) {
	fastqFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: Problem opening the given file")
	}
	defer fastqFile.Close()

	scanner := bufio.NewScanner(fastqFile)

	var reads, tempReadList []string
	startSeq := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			if string(line[0]) == "@" {
				startSeq = true
				continue
			}
			if string(line[0]) == "+" {
				startSeq = false
				reads = append(reads, strings.Join(tempReadList, ""))
				tempReadList = []string{}
				continue
			}
			if startSeq == true {
				tempReadList = append(tempReadList, scanner.Text())
			}
		}

	}

	revReads := GenerateReadRevComps(reads)

	return reads, revReads
}

// GenerateReadPath
func GenerateReadPath(read string, l int) *ReadPath {
	rp := &ReadPath{}
	lastNode, currNode := &PathNode{}, &PathNode{}
	for i := 0; i < len(read)-l+2; i++ {
		if i == 0 {
			currNode.value = string(read[i : i+l-1])
			rp.head = currNode
		} else {
			currNode.value = string(read[i : i+l-1])
			lastNode.next = currNode
			lastNode.edge = &PathEdge{value: string(lastNode.value) + string(currNode.value[len(currNode.value)-1])}

		}
		lastNode = currNode
		currNode = &PathNode{}
		rp.len++
	}
	return rp
}

// GenerateReadPathSet returns a point to a PathSet with *ReadPaths for every read
func GenerateReadPathSet(reads []string, l int) *PathSet {
	var ps PathSet
	ps = make([]*ReadPath, len(reads))
	for i, read := range reads {
		ps[i] = GenerateReadPath(read, l)
	}
	return &ps
}
