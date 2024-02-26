# poseidon

A Golang and Gnark implementation of the Poseidon hash function.
The Golang version is an adaption of the [iden3](https://github.com/iden3/go-iden3-crypto/tree/master/poseidon)
implementation, but with support for multiple curves.
The Gnark implementation is an adaption of the [circom](https://github.com/iden3/circomlib/blob/master/circuits/poseidon.circom)
implementation, available in both native field and emulated versions.

### Usage

Standard:
```golang
poseidon.Hash[*fr.Element]([]*big.Int{in1, in2})
```

Gnark circuit:
```golang
poseidon.Hash(api, []frontend.Variable{in1, in2})
```

Gnark circuit using emulated field:
```golang
bnField, _ := emulated.NewField[sw_bn254.ScalarField](api)
poseidon.Hash(bnField, []*emulated.Element[sw_bn254.ScalarField]{in1, in2})
```

### Constants

The [constants](./constants/) were generated using a combination of a version of the
[poseidon sage script](https://extgit.iaik.tugraz.at/krypto/hadeshash/-/blob/master/code/generate_params_poseidon.sage)
from the [hadeshash](https://extgit.iaik.tugraz.at/krypto/hadeshash) project, and
[triplewz's](https://github.com/triplewz/poseidon) generator implementation forked to support multiple field elements
provided by [gnark-crypto](https://github.com/Consensys/gnark-crypto). The hadeshash script has a
[minor modification](https://github.com/mdehoog/poseidon/commit/39c59b520c44d1c94ac95fc0789db4910a39f25e)
to round up the `Rp` value to the nearest multiple of `t`. The generated constants match the
[constants in the circom library](https://github.com/iden3/circomlib/blob/master/circuits/poseidon_constants.circom)
for the BN254 curve.

You can regenerate the constants using `make constants`. The repo currently has constants generated for:

| Curve     | Alpha | Constants                                       |
|-----------|-------|-------------------------------------------------|
| BN254     | 5     | [constants/bn254.go](./constants/bn254.go)      |
| BLS12-381 | 5     | [constants/bls12_381.go](./constants/12_381.go) |
| BW6-761   | 5     | [constants/bw6_761.go](./constants/bw6_761.go)  |

Note that other alpha values are not yet supported in the hash implementations.
