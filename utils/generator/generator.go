package generator

import (
	"fmt"
	"strconv"
)

const finalInStan string = "zz999"
const firstInStan string = "AA001"

func GenerateInStan(lastInStan string) (string, error) {
	if len(lastInStan) != 6 {
		return "", fmt.Errorf("invalid inStan value: " + lastInStan)
	}

	isAlphabetReset := false
	isNumericReset := false

	stanType := lastInStan[0]
	seqAlphaNumeric := lastInStan[1:]
	seqNumeric, err := strconv.Atoi(seqAlphaNumeric[2:])
	if err != nil {
		return "", fmt.Errorf("error converting string to int")
	}
	seqFirstAlphabet := seqAlphaNumeric[0]
	seqSecondAlphabet := seqAlphaNumeric[1]

	if seqAlphaNumeric == finalInStan {
		return fmt.Sprintf("%c%s", rune(stanType), firstInStan), nil
	}

	seqNumeric++
	if seqNumeric == 1000 {
		seqNumeric = 1
		isNumericReset = true
	}

	if isNumericReset {
		seqSecondAlphabet++
	}

	if seqSecondAlphabet == 123 {
		isAlphabetReset = true
		seqSecondAlphabet = 65
	}

	if isAlphabetReset {
		seqFirstAlphabet++
	}

	newStan := fmt.Sprintf("%c%c%c%03d", rune(stanType), rune(seqFirstAlphabet), rune(seqSecondAlphabet), seqNumeric)

	return newStan, nil
}
