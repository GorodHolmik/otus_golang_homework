package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString   = errors.New("invalid string")
	ErrUnknownRuneType = errors.New("unknown rune type")
	ErrWrongRuneType   = errors.New("wrong rune type")
)

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

func newUnpackState(char rune) UnpackState {
	return UnpackState{
		char:     char,
		runeType: defineRuneType(char),
	}
}

func newEmptyUnpackState() UnpackState {
	return UnpackState{runeType: Empty}
}

func newBackSlashedCopy(charState UnpackState) UnpackState {
	return UnpackState{
		char:        charState.char,
		runeType:    charState.runeType,
		backSlashed: true,
	}
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
	switch previousState.runeType {
	case Empty:
		return processEmptyState(previousState, currentState)
	case Letter:
		return processLetterState(previousState, currentState, stringBuilder)
	case Digit:
		return processDigitState(previousState, currentState, stringBuilder)
	case BackSlash:
		return processBackSlashState(previousState, currentState, stringBuilder)
	}

	return newEmptyUnpackState(), nil
}

func processBackSlashState(
	previousState UnpackState,
	currentState UnpackState,
	stringBuilder *strings.Builder,
) (UnpackState, error) {
	if previousState.runeType != BackSlash {
		return newEmptyUnpackState(), ErrWrongRuneType
	}

	switch currentState.runeType {
	case Digit:
		if previousState.backSlashed {
			count, _ := strconv.Atoi(string(currentState.char))
			repeatedChar := strings.Repeat(string(previousState.char), count)
			stringBuilder.WriteString(repeatedChar)

			return newEmptyUnpackState(), nil
		}
		return newBackSlashedCopy(currentState), nil

	case BackSlash:
		if previousState.backSlashed {
			stringBuilder.WriteRune(previousState.char)

			return currentState, nil
		}
		return newBackSlashedCopy(currentState), nil

	case Letter:
		fallthrough

	case Empty:
		if previousState.backSlashed {
			stringBuilder.WriteRune(previousState.char)

			return currentState, nil
		}
		return newEmptyUnpackState(), ErrInvalidString
	default:
		return newEmptyUnpackState(), ErrUnknownRuneType
	}
}

func processDigitState(
	previousState UnpackState,
	currentState UnpackState,
	stringBuilder *strings.Builder,
) (UnpackState, error) {
	if previousState.runeType != Digit {
		return newEmptyUnpackState(), ErrWrongRuneType
	}

	if !previousState.backSlashed {
		return newEmptyUnpackState(), ErrInvalidString
	}
	switch currentState.runeType {
	case Empty:
		fallthrough

	case BackSlash:
		fallthrough

	case Letter:
		stringBuilder.WriteRune(previousState.char)

		return currentState, nil

	case Digit:
		count, _ := strconv.Atoi(string(currentState.char))
		repeatedChar := strings.Repeat(string(previousState.char), count)
		stringBuilder.WriteString(repeatedChar)

		return newEmptyUnpackState(), nil
	default:
		return newEmptyUnpackState(), ErrUnknownRuneType
	}
}

func processLetterState(
	previousState UnpackState,
	currentState UnpackState,
	stringBuilder *strings.Builder,
) (UnpackState, error) {
	if previousState.runeType != Letter {
		return newEmptyUnpackState(), ErrWrongRuneType
	}

	if previousState.backSlashed {
		return newEmptyUnpackState(), ErrInvalidString
	}

	switch currentState.runeType {
	case Empty:
		fallthrough
	case BackSlash:
		fallthrough
	case Letter:
		stringBuilder.WriteRune(previousState.char)

		return currentState, nil
	case Digit:
		count, err := strconv.Atoi(string(currentState.char))
		if err != nil {
			return newEmptyUnpackState(), err
		}

		repeatedChar := strings.Repeat(string(previousState.char), count)
		stringBuilder.WriteString(repeatedChar)

		return newEmptyUnpackState(), nil
	default:
		return newEmptyUnpackState(), ErrUnknownRuneType
	}
}

func processEmptyState(
	previousState UnpackState,
	currentState UnpackState,
) (UnpackState, error) {
	if previousState.runeType != Empty {
		return newEmptyUnpackState(), ErrWrongRuneType
	}
	return currentState, nil
}

func Unpack(input string) (string, error) {
	var err error
	previousState := newEmptyUnpackState()
	stringBuilder := strings.Builder{}

	for _, currentRune := range input {
		previousState, err = buildUnpacked(
			previousState,
			newUnpackState(currentRune),
			&stringBuilder,
		)
		if err != nil {
			return "", err
		}
	}

	_, err = buildUnpacked(previousState, newEmptyUnpackState(), &stringBuilder)
	if err != nil {
		return "", err
	}

	return stringBuilder.String(), nil
}
