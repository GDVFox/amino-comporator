package main

var conf = NewConfig()

// Config contains settings of util
type Config struct {
	MamaevsPairs []*MamaevsPair `toml:"mamaevs-pairs"`
	Threads      int            `toml:"threads"`
}

// NewConfig returns new instance of config
func NewConfig() *Config {
	return &Config{
		MamaevsPairs: []*MamaevsPair{
			{
				Size:   4,
				Weight: 64,
			},
			{
				Size:   8,
				Weight: 32,
			},
			{
				Size:   16,
				Weight: 16,
			},
			{
				Size:   32,
				Weight: 8,
			},
			{
				Size:   64,
				Weight: 4,
			},
		},
		Threads: 4,
	}
}

// MamaevsPair size of comparable substring and weight of successful comparison
type MamaevsPair struct {
	Size   int `toml:"size"`
	Weight int `toml:"weight"`
}
