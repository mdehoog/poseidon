package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"math"
	"math/big"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/triplewz/poseidon"

	"github.com/mdehoog/poseidon/constants"
)

func GenerateConstantsFile[E poseidon.Element[E]](path string, levels, startingLevel, field, sbox, alpha int) error {
	fieldSize := poseidon.Bits[E]()
	securityLevel := poseidon.SecurityLevel
	modulus := poseidon.Modulus[E]()

	cs := constants.Strings{
		F: modulus.Text(16),
		C: make([][]string, levels),
		S: make([][]string, levels),
		M: make([][][]string, levels),
		P: make([][][]string, levels),
	}

	for level := 0; level < levels; level++ {
		width := level + startingLevel
		args := []string{
			"./generator/generate_params_poseidon.sage",
			strconv.Itoa(field),
			strconv.Itoa(sbox),
			strconv.Itoa(fieldSize),
			strconv.Itoa(width),
			strconv.Itoa(alpha),
			strconv.Itoa(securityLevel),
			fmt.Sprintf("0x%s", modulus.Text(16)),
		}
		fmt.Printf("Executing 'sage %s'\n", strings.Join(args, " "))
		out, err := exec.Command("sage", args...).Output()
		if err != nil {
			return fmt.Errorf("sage command failed: %w", err)
		}

		rf, err := strconv.Atoi(string(regexp.MustCompile(`R_F=(\d+)`).FindSubmatch(out)[1]))
		if err != nil {
			return fmt.Errorf("could not parse Rf: %w", err)
		}
		rp, err := strconv.Atoi(string(regexp.MustCompile(`R_P=(\d+)`).FindSubmatch(out)[1]))
		if err != nil {
			return fmt.Errorf("could not parse Rp: %w", err)
		}

		mdsMatrixStart := bytes.Index(out, []byte("MDS matrix:"))
		mdsMatrixEnd := mdsMatrixStart + bytes.Index(out[mdsMatrixStart:], []byte("]]"))
		mdsMatrixString := string(out[mdsMatrixStart+12 : mdsMatrixEnd+2])
		hexStringRegexp := regexp.MustCompile(`'0x[0-9a-fA-F]+'`)

		var mdsMatrix poseidon.Matrix[E]
		mdsMatrixStrings := hexStringRegexp.FindAllString(mdsMatrixString, -1)
		mdsWidth := int(math.Round(math.Sqrt(float64(len(mdsMatrixStrings)))))
		for i := 0; i < mdsWidth; i++ {
			mdsMatrix = append(mdsMatrix, make([]E, mdsWidth))
			for j := 0; j < mdsWidth; j++ {
				match := mdsMatrixStrings[j*mdsWidth+i]
				mdsValue, ok := new(big.Int).SetString(match[3:len(match)-1], 16)
				if !ok {
					return fmt.Errorf("could not parse hex value: %s", match)
				}
				mdsMatrix[i][j] = poseidon.NewElement[E]().SetBigInt(mdsValue)
			}
		}

		pc, err := poseidon.GenCustomPoseidonConstants[E](width, field, sbox, rf, rp, mdsMatrix)
		if err != nil {
			return fmt.Errorf("generate constants error: %w", err)
		}

		cs.C[level] = make([]string, len(pc.CompRoundConsts))
		for i, e := range pc.CompRoundConsts {
			cs.C[level][i] = e.BigInt(new(big.Int)).Text(16)
		}

		for _, e := range pc.Sparse {
			for _, w := range e.WHat {
				cs.S[level] = append(cs.S[level], w.BigInt(new(big.Int)).Text(16))
			}
			for _, v := range e.V {
				cs.S[level] = append(cs.S[level], v.BigInt(new(big.Int)).Text(16))
			}
		}

		cs.M[level] = make([][]string, len(mdsMatrix))
		for i, e := range mdsMatrix {
			cs.M[level][i] = make([]string, len(e))
			for j, f := range e {
				cs.M[level][i][j] = f.BigInt(new(big.Int)).Text(16)
			}
		}

		cs.P[level] = make([][]string, len(pc.PreSparse))
		for i, e := range pc.PreSparse {
			cs.P[level][i] = make([]string, len(e))
			for j, f := range e {
				cs.P[level][i][j] = f.BigInt(new(big.Int)).Text(16)
			}
		}
	}

	var b bytes.Buffer
	err := generateTemplate(&b, &cs)
	if err != nil {
		return fmt.Errorf("generate template error: %w", err)
	}
	source, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("format source error: %w", err)
	}
	err = os.WriteFile(path, source, 0644)
	if err != nil {
		return fmt.Errorf("write file error: %w", err)
	}
	return nil
}
