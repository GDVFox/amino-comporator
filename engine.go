package main

import (
	"hash/fnv"
	"sync"

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

func buildSequenceHashes(seq *AminoSequence) []*u64set.Set {
	hashes := make([]*u64set.Set, 0, len(conf.MamaevsPairs))
	for _, pair := range conf.MamaevsPairs {
		hashes = append(hashes, hashSequence(seq.Value, pair.Size))
	}
	return hashes
}

func buildSequencesHashes(seqs []*AminoSequence) map[string][]*u64set.Set {
	hashSets := make(map[string][]*u64set.Set, len(seqs))
	for _, seq := range seqs {
		hashSets[seq.ID] = buildSequenceHashes(seq)
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

type concBldHashes struct {
	seqNum int
	hashes []*u64set.Set
}

func runCalculateHashes(seqs []*AminoSequence, tasks *Semaphore) <-chan concBldHashes {
	hashesCh := make(chan concBldHashes, tasks.Len())
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(seqs))
		for i, seq := range seqs {
			tasks.Accuire()
			go func(i int, seq *AminoSequence) {
				hashes := concBldHashes{
					seqNum: i,
					hashes: buildSequenceHashes(seq),
				}
				tasks.Release()
				hashesCh <- hashes
				wg.Done()
			}(i, seq)
		}
		wg.Wait()
		close(hashesCh)
	}()
	return hashesCh
}

type concBldValue struct {
	seqNum1 int
	seqNum2 int
	value   float64
}

func runCalculateValues(hashesCh <-chan concBldHashes, tasks *Semaphore) <-chan concBldValue {
	valueCh := make(chan concBldValue, tasks.Len())
	go func() {
		wg := sync.WaitGroup{}
		accepted := make([]concBldHashes, 0)
		for hashes := range hashesCh {
			accepted = append(accepted, hashes)
			lastAccepted := len(accepted) - 1
			wg.Add(len(accepted) - 1)
			for i := 0; i < lastAccepted; i++ {
				tasks.Accuire()
				go func(h1, h2 concBldHashes) {
					value := concBldValue{
						seqNum1: MinInt(h1.seqNum, h2.seqNum),
						seqNum2: MaxInt(h1.seqNum, h2.seqNum),
						value:   mamaevsValue(h1.hashes, h2.hashes),
					}
					tasks.Release()
					valueCh <- value
					wg.Done()
				}(accepted[i], accepted[lastAccepted])
			}
		}
		wg.Wait()
		close(valueCh)
	}()
	return valueCh
}

// BuildComparisonMatrixConcurrent builds comparison matrix for seqs using cap amount of threads for calculating
func BuildComparisonMatrixConcurrent(seqs []*AminoSequence, cap int) *mat.Dense {
	if cap <= 1 {
		return BuildComparisonMatrix(seqs)
	}

	tasks := NewSemaphore(uint(cap))

	hashesCh := runCalculateHashes(seqs, tasks)
	valueCh := runCalculateValues(hashesCh, tasks)

	comparisonMatrix := mat.NewDense(len(seqs), len(seqs), nil)
	for v := range valueCh {
		comparisonMatrix.Set(v.seqNum1, v.seqNum2, v.value)
	}
	return comparisonMatrix
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
