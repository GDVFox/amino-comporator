# 🧪Amino Acid Comparator

## Description

Calculates comparison matrix for amino acids from `.fasta` input file. Each element of matrix is in the range [0; 1], where 1 means a complete match.

### Algorithm

Mamaev's Pair - pair of size and weight that are used to calculate similarity between acids

Algorithms concept:

1. Make substrings of some size from acid code.
2. Calculate Hashes of substrings. Make set of this hashes.
3. Repeat for different sizes. Make Set of sets for different sizes.
4. Compare two acids:
   1. Calculate similarity for all pairs of sets with same Mamaevs size
   2. Using Mamaevs weight calculate similarity

### Example

For amino acids `MHSKVT` and `THSKVM`:

Calculate Mamaevs sizes:

1. `MHSKVT`:
   |Mamaevs size | substrings | set |
   |----|----|----|
   | 3  |{MHS, HSK, SKV, KVT}| {1,2,3,4}|
   | 4  |{MHSK, HSKV, SKVT}|  {5,6,7} |
2. `THSKVM`:
   |Mamaevs size | substrings | set |
   |----|----|----|
   | 3  |{THS, HSK, SKV, KVM}| {10,2,3,11}|
   | 4  |{THSK, HSKV, SKVM}|  {12,6,13} |

Compare hash sets:

| Mamaevs size | first | second | similarity | Mamaevs weight |  weighted similarity |
|------|-------|--------|------------|--------|----------------------|
|  3   | {1,2,3,4} | {10,2,3,11} | 2 | 3      | 2*3 = 6 |
|  4   | {5,6,7} |  {12,6,13} | 1 | 7      | 1*7 = 7 |

Result:

Total similarity : 15, after normalization: 0,63

*Note: Not actual values and formulas, used only as example. Real formulas can be found in code*

## Build

```bash
mkdir _build && go build -o _build/amino-comporator *.go
```

## Run

```bash
cd _build && ./amino-comporator --config=config.conf --data=<your_fasta_file> --output=human
```

Possible `output` modes:

* `--output=matrix` - output as math matrix;
* `--output=python` - output as python two dim array;
* `--output=matlab` - output as MATLAB matrix;
* `--output=human` - output as human readable text;

(c) BowWow Team + 🐓
