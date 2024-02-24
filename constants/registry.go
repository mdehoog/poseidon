package constants

import "math/big"

func C(modulus *big.Int, t int) ([]*big.Int, bool) {
	p, ok := fields[modulus.Text(16)]
	if !ok {
		return nil, false
	}
	return p.C[t], true
}

func S(modulus *big.Int, t int) ([]*big.Int, bool) {
	p, ok := fields[modulus.Text(16)]
	if !ok {
		return nil, false
	}
	return p.S[t], true
}

func M(modulus *big.Int, t int) ([][]*big.Int, bool) {
	p, ok := fields[modulus.Text(16)]
	if !ok {
		return nil, false
	}
	return p.M[t], true
}

func P(modulus *big.Int, t int) ([][]*big.Int, bool) {
	p, ok := fields[modulus.Text(16)]
	if !ok {
		return nil, false
	}
	return p.P[t], true
}

type Strings struct {
	F string
	C [][]string
	S [][]string
	M [][][]string
	P [][][]string
}

type parsed struct {
	C [][]*big.Int
	S [][]*big.Int
	M [][][]*big.Int
	P [][][]*big.Int
}

var fields = make(map[string]*parsed)

func register(c *Strings) {
	p := &parsed{
		C: make([][]*big.Int, len(c.C)),
		S: make([][]*big.Int, len(c.S)),
		M: make([][][]*big.Int, len(c.M)),
		P: make([][][]*big.Int, len(c.P)),
	}
	modulus, _ := new(big.Int).SetString(c.F, 16)
	fields[modulus.Text(16)] = p

	for i := range c.C {
		p.C[i] = make([]*big.Int, len(c.C[i]))
		for j := range c.C[i] {
			p.C[i][j], _ = new(big.Int).SetString(c.C[i][j], 16)
		}
	}
	for i := range c.S {
		p.S[i] = make([]*big.Int, len(c.S[i]))
		for j := range c.S[i] {
			p.S[i][j], _ = new(big.Int).SetString(c.S[i][j], 16)
		}
	}
	for i := range c.M {
		p.M[i] = make([][]*big.Int, len(c.M[i]))
		for j := range c.M[i] {
			p.M[i][j] = make([]*big.Int, len(c.M[i][j]))
			for k := range c.M[i][j] {
				p.M[i][j][k], _ = new(big.Int).SetString(c.M[i][j][k], 16)
			}
		}
	}
	for i := range c.P {
		p.P[i] = make([][]*big.Int, len(c.P[i]))
		for j := range c.P[i] {
			p.P[i][j] = make([]*big.Int, len(c.P[i][j]))
			for k := range c.P[i][j] {
				p.P[i][j][k], _ = new(big.Int).SetString(c.P[i][j][k], 16)
			}
		}
	}
}
