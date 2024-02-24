package main

import (
	"fmt"

	bls12381fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6761fr "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"

	"github.com/mdehoog/poseidon/generator"
)

const levels = 16
const startingLevel = 2
const field = 1
const sbox = 0

type generateFunc func(path string, levels, startingLevel, field, sbox, alpha int) error

type config struct {
	generateFunc generateFunc
	alpha        int
}

// for alpha values, see https://eprint.iacr.org/2019/458.pdf (page 6)
var fields = map[string]config{
	"bn254": {
		generateFunc: generator.GenerateConstantsFile[*bn254fr.Element],
		alpha:        5,
	},
	"bls12_381": {
		generateFunc: generator.GenerateConstantsFile[*bls12381fr.Element],
		alpha:        5,
	},
	"bw6_761": {
		generateFunc: generator.GenerateConstantsFile[*bw6761fr.Element],
		alpha:        5,
	},
}

func main() {
	for name, cfg := range fields {
		if err := cfg.generateFunc(fmt.Sprintf("constants/%s.go", name), levels, startingLevel, field, sbox, cfg.alpha); err != nil {
			panic(err)
		}
	}
}
