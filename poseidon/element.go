package poseidon

import (
	"math/big"
	"reflect"

	"github.com/consensys/gnark-crypto/field/pool"
)

type Element[E any] interface {
	SetUint64(uint64) E
	SetBigInt(*big.Int) E
	BigInt(*big.Int) *big.Int
	SetOne() E
	SetZero() E
	Inverse(E) E
	Set(E) E
	Square(E) E
	Mul(E, E) E
	Add(E, E) E
	Sub(E, E) E
	Marshal() []byte
}

func NewElement[E Element[E]]() E {
	typ := reflect.TypeOf((*E)(nil)).Elem()
	val := reflect.New(typ.Elem())
	return val.Interface().(E)
}

func Modulus[E Element[E]]() *big.Int {
	e := NewElement[E]()
	modulus := e.Sub(e, NewElement[E]().SetOne()).BigInt(new(big.Int))
	modulus.Add(modulus, big.NewInt(1))
	return modulus
}

func BigIntsToElements[E Element[E]](bi []*big.Int) []E {
	o := make([]E, len(bi))
	for i := range bi {
		o[i] = NewElement[E]().SetBigInt(bi[i])
	}
	return o
}

func BigInts2DtoElements2D[E Element[E]](bi [][]*big.Int) [][]E {
	o := make([][]E, len(bi))
	for i := range bi {
		o[i] = BigIntsToElements[E](bi[i])
	}
	return o
}

// Exp is a copy of gnark-crypto's implementation, but takes a pointer argument
func Exp[E Element[E]](z, x E, k *big.Int) {
	if k.IsUint64() && k.Uint64() == 0 {
		z.SetOne()
	}

	e := k
	if k.Sign() == -1 {
		// negative k, we invert
		// if k < 0: xᵏ (mod q) == (x⁻¹)ᵏ (mod q)
		x.Inverse(x)

		// we negate k in a temp big.Int since
		// Int.Bit(_) of k and -k is different
		e = pool.BigInt.Get()
		defer pool.BigInt.Put(e)
		e.Neg(k)
	}

	z.Set(x)

	for i := e.BitLen() - 2; i >= 0; i-- {
		z.Square(z)
		if e.Bit(i) == 1 {
			z.Mul(z, x)
		}
	}
}
