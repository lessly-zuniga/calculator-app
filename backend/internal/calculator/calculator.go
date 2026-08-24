package calculator

import (
	"errors"
	"math"
)

var (
	ErrInvalidOperandCount = errors.New("exactly two operands are required")
	ErrDivisionByZero      = errors.New("cannot divide by zero")
	ErrNegativeSquareRoot  = errors.New("cannot calculate the square root of a negative number")
)

func Add(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidOperandCount
	}

	return operands[0] + operands[1], nil
}

func Subtract(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidOperandCount
	}

	return operands[0] - operands[1], nil
}

func Multiply(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidOperandCount
	}

	return operands[0] * operands[1], nil
}

func Divide(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidOperandCount
	}

	if operands[1] == 0 {
		return 0, ErrDivisionByZero
	}

	return operands[0] / operands[1], nil
}

func SquareRoot(operands []float64) (float64, error) {
	if len(operands) != 1 {
		return 0, ErrInvalidOperandCount
	}

	if operands[0] < 0 {
		return 0, ErrNegativeSquareRoot
	}

	return math.Sqrt(operands[0]), nil
}
