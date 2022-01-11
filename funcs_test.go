package main

import (
	"reflect"
	"testing"
)

func TestReverseComplement(t *testing.T) {
	test1 := "ACTG"
	test2 := "ATACGC"
	test3 := "ACGTGCA"

	v1 := ReverseComplement(test1)
	v2 := ReverseComplement(test2)
	v3 := ReverseComplement(test3)

	if v1 != "CAGT" {
		t.Errorf("ReverseComplement(test1) = %s; wants TGAC", v1)
	}
	if v2 != "GCGTAT" {
		t.Errorf("ReverseComplement(test1) = %s; wants TATGCG", v2)
	}
	if v3 != "TGCACGT" {
		t.Errorf("ReverseComplement(test1) = %s; wants TGCACGT", v3)
	}
}

// Returns if two lists contain the same elements
func ListsEqual(l1, l2 []string) bool {
	l1Map, l2Map := make(map[string]int, 0), make(map[string]int, 0)
	for _, v1 := range l1 {
		l1Map[v1]++
	}
	for _, v2 := range l2 {
		l2Map[v2]++
	}
	return reflect.DeepEqual(l1Map, l2Map)
}

func TestGenerateReadLTuples(t *testing.T) {
	test1 := "ACTG"
	test2 := "ATACGC"
	test3 := "ACGTGCA"

	l := 3
	v1 := GenerateReadLTuples(test1, l)
	v2 := GenerateReadLTuples(test2, l)
	v3 := GenerateReadLTuples(test3, l)

	if !ListsEqual(v1, []string{"ACT", "CTG"}) {
		t.Error("GenerateLTuples(test1,l) =", v1)
	}
	if !ListsEqual(v2, []string{"ATA", "TAC", "ACG", "CGC"}) {
		t.Error("GenerateLTuples(test2,l) =", v2)
	}
	if !ListsEqual(v3, []string{"ACG", "CGT", "GTG", "TGC", "GCA"}) {
		t.Error("GenerateLTuples(test3,l) =", v3)
	}

}

func TestGenerateSampleLTuples(t *testing.T) {
	// ACGCGTCG
	t1reads := []string{
		"ACGC",
		"GCGTC",
		"CGCGT",
		"GCGTCG",
		"ACGCGT",
	}

	v1 := GenerateSamleLTuples(t1reads, 3, "")

	if !ListsEqual(v1, []string{"ACG", "CGC", "GCG", "CGT", "GTC", "TCG"}) {
		t.Error("GenerateSampleLTuples(test1,3,false) =", v1)
	}

}
