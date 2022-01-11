package main

import (
	"bufio"
	"fmt"
	"os"
)

//Keys returns the keys from a map of strings to ints
func Keys(m map[string]int) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

//GenerateReadLTuples returns a list of unique l-tuples in a string
func GenerateReadLTuples(s string, l int) []string {
	lTupMap := make(map[string]int)
	for i := 0; i <= len(s)-l; i++ {
		lTup := string(s[i : i+l])
		if _, ok := lTupMap[lTup]; !ok {
			lTupMap[lTup]++
		}
	}
	return Keys(lTupMap)
}

//SaveLTuples writes all unique lTuple for a collection of reads to file
// File will be created in pwd unless a full or partial path is provided
// All folders in savepath must already exist
func SaveLTuples(lTups []string, savepath string) {
	openFile, err := os.Create(savepath)
	if err != nil {
		panic("Could not create file from given file path")
	}
	defer openFile.Close()
	writer := bufio.NewWriter(openFile)

	for _, tup := range lTups {
		fmt.Fprintln(writer, tup)
	}
	writer.Flush()
}

//GenerateSampleLTuples returns a list of all unique Ltuples from a collection of reads
// If save is a non empty string the set of unique l-tuples will be saved to file
func GenerateSamleLTuples(reads []string, l int, save string) []string {
	lTupMap := make(map[string]int)
	for _, read := range reads {
		readLTups := GenerateReadLTuples(read, l)
		for _, lTup := range readLTups {
			if _, ok := lTupMap[lTup]; !ok {
				lTupMap[lTup]++
			}
		}
	}

	lTups := Keys(lTupMap)
	if save != "" {
		SaveLTuples(lTups, save+".txt")
	}

	return lTups
}

// GenerateForwardRevLTuples Returns a list of unique L tups for forward and reverse complement of a reads
func GenerateForwardRevLTuples(fwReads, revReads []string, l int, save string) ([]string, []string) {
	fwLtups := GenerateSamleLTuples(fwReads, l, save)
	revLtups := GenerateSamleLTuples(revReads, l, save)
	return fwLtups, revLtups
}
