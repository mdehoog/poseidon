module github.com/mdehoog/poseidon

go 1.21.1

require (
	github.com/consensys/gnark-crypto v0.12.1
	github.com/triplewz/poseidon v0.0.0-20230828015038-79d8165c88ed
)

require (
	github.com/bits-and-blooms/bitset v1.7.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
)

replace github.com/triplewz/poseidon => github.com/mdehoog/triplewz-poseidon v0.0.0-20240224062259-6ac0a3e35064
