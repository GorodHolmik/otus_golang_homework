package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type RuneType string

var (
	Digit     = RuneType("DIGIT")
	Letter    = RuneType("LETTER")
	BackSlash = RuneType("BACKSLASH")
	Empty     = RuneType("EMPTY")
)

type UnpackState struct {
	char        rune
	runeType    RuneType
	backSlashed bool
}

func newBackSlashedCopy(charState UnpackState) UnpackState {
	return UnpackState{charState.char, charState.runeType, true}
}

func defineRuneType(char rune) RuneType {
	switch {
	case unicode.IsDigit(char):
		return Digit
	case char == '\\':
		return BackSlash
	default:
		return Letter
	}
}

func buildUnpacked(
	previousState UnpackState,
	currentState UnpackState,
	stringBuilder *strings.Builder,
) (UnpackState, error) {
	if previousState.runeType == BackSlash {
		if currentState.runeType == Letter || currentState.runeType == Empty {
			return UnpackState{runeType: Empty}, ErrInvalidString
		} else if !previousState.backSlashed {
			return newBackSlashedCopy(currentState), nil
		}
	}
	switch currentState.runeType {
	case Digit:
		if previousState.runeType == Empty || (previousState.runeType == Digit && !previousState.backSlashed) {
			return UnpackState{runeType: Empty}, ErrInvalidString
		}
		count, err := strconv.Atoi(string(currentState.char))
		if err != nil {
			return UnpackState{runeType: Empty}, err
		}
		repeatedChar := strings.Repeat(string(previousState.char), count)
		stringBuilder.WriteString(repeatedChar)
		return UnpackState{runeType: Empty}, nil
	default:
		if previousState.runeType != Empty {
			stringBuilder.WriteRune(previousState.char)
		}
		return currentState, nil
	}
}

func Unpack(input string) (string, error) {
	var err error
	previousState := UnpackState{runeType: Empty}
	stringBuilder := strings.Builder{}
	for _, currentRune := range input {
		previousState, err = buildUnpacked(
			previousState,
			UnpackState{char: currentRune, runeType: defineRuneType(currentRune)},
			&stringBuilder,
		)
		if err != nil {
			return "", err
		}
	}
	_, err = buildUnpacked(previousState, UnpackState{runeType: Empty}, &stringBuilder)
	if err != nil {
		return "", err
	}
	return stringBuilder.String(), nil
}
