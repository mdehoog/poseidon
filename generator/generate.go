package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/triplewz/poseidon"
)

func main() {
	field := 1
	sbox := 0
	fieldSize := fr.Bits
	alpha := poseidon.Alpha // TODO check value for different curves
	securityLevel := poseidon.SecurityLevel

	for width := 2; width <= 2; width++ {
		args := []string{
			"./generator/generate_params_poseidon.sage",
			strconv.Itoa(field),
			strconv.Itoa(sbox),
			strconv.Itoa(fieldSize),
			strconv.Itoa(width),
			strconv.Itoa(alpha),
			strconv.Itoa(securityLevel),
			fmt.Sprintf("0x%s", fr.Modulus().Text(16)),
		}
		fmt.Printf("Executing 'sage %s'\n", strings.Join(args, " "))
		out, err := exec.Command("sage", args...).Output()
		if err != nil {
			panic(err)
		}
		roundConstantsStart := bytes.Index(out, []byte("Round constants for GF(p):"))
		roundConstantsEnd := roundConstantsStart + bytes.Index(out[roundConstantsStart:], []byte("]"))
		roundConstantsString := string(out[roundConstantsStart+27 : roundConstantsEnd+1])
		mdsMatrixStart := bytes.Index(out, []byte("MDS matrix:"))
		mdsMatrixEnd := mdsMatrixStart + bytes.Index(out[mdsMatrixStart:], []byte("]]"))
		mdsMatrixString := string(out[mdsMatrixStart+12 : mdsMatrixEnd+2])
		hexStringRegexp := regexp.MustCompile(`'0x[0-9a-fA-F]+'`)

		var roundConstants []*fr.Element
		for _, match := range hexStringRegexp.FindAllString(roundConstantsString, -1) {
			roundConstant, ok := new(big.Int).SetString(match[3:len(match)-1], 16)
			if !ok {
				panic(fmt.Sprintf("could not parse hex value: %s", match))
			}
			roundConstants = append(roundConstants, new(fr.Element).SetBigInt(roundConstant))
		}

		var mdsMatrix poseidon.Matrix[*fr.Element]
		mdsMatrixStrings := hexStringRegexp.FindAllString(mdsMatrixString, -1)
		mdsWidth := int(math.Round(math.Sqrt(float64(len(mdsMatrixStrings)))))
		for i := 0; i < mdsWidth; i++ {
			mdsMatrix = append(mdsMatrix, make([]*fr.Element, mdsWidth))
			for j := 0; j < mdsWidth; j++ {
				match := mdsMatrixStrings[j*mdsWidth+i]
				mdsValue, ok := new(big.Int).SetString(match[3:len(match)-1], 16)
				if !ok {
					panic(fmt.Sprintf("could not parse hex value: %s", match))
				}
				mdsMatrix[i][j] = new(fr.Element).SetBigInt(mdsValue)
			}
		}

		constants, err := poseidon.GenPoseidonConstants[*fr.Element](width, field, sbox, true, mdsMatrix)
		if err != nil {
			panic(err)
		}

		for _, rc := range constants.CompRoundConsts {
			fmt.Printf("0x%s\n", rc.Text(16))
		}
	}
}
