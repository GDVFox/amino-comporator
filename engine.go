package main

import (
	"hash/fnv"

	"github.com/scylladb/go-set/u64set"
	"gonum.org/v1/gonum/mat"
)

// AminoSequence represents description of amino acid
type AminoSequence struct {
	ID          string
	Description string
	Value       string
}

func hashPart(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func hashSequence(s string, l int) *u64set.Set {
	set := u64set.NewWithSize(len(s)/l + 1)
	beg, end := 0, l
	for end <= len(s) {
		set.Add(hashPart(s[beg:end]))
		beg++
		end++
	}
	return set
}

func buildSequencesHashes(seqs []*AminoSequence) map[string][]*u64set.Set {
	hashSets := make(map[string][]*u64set.Set, len(seqs))
	for _, seq := range seqs {
		hashSets[seq.ID] = make([]*u64set.Set, 0, len(conf.MamaevsPairs))
		for _, pair := range conf.MamaevsPairs {
			hashSets[seq.ID] = append(hashSets[seq.ID], hashSequence(seq.Value, pair.Size))
		}
	}

	return hashSets
}

func mamaevsValue(hs1, hs2 []*u64set.Set) float64 {
	mamaevsValue, maxValue := 0., 0.
	for i, pair := range conf.MamaevsPairs {
		maxValue += float64(MaxInt(hs1[i].Size(), hs2[i].Size()) * pair.Weight)
		mamaevsValue += float64(u64set.Intersection(hs1[i], hs2[i]).Size() * pair.Weight)
	}

	// normalize the value so that longer sequences don't score more points
	return mamaevsValue / maxValue
}

// BuildComparisonMatrix builds comparison matrix for seqs
func BuildComparisonMatrix(seqs []*AminoSequence) *mat.Dense {
	comparisonMatrix := mat.NewDense(len(seqs), len(seqs), nil)
	hashes := buildSequencesHashes(seqs)

	for i := 0; i < len(seqs)-1; i++ {
		for j := i + 1; j < len(seqs); j++ {
			comparisonMatrix.Set(i, j, mamaevsValue(hashes[seqs[i].ID], hashes[seqs[j].ID]))
		}
	}

	return comparisonMatrix
}
