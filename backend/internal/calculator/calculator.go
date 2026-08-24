package calculator

import (
	"errors"
	"math"
)

var (
	ErrInvalidBinaryOperandCount = errors.New("exactly two operands are required")
	ErrInvalidUnaryOperandCount  = errors.New("exactly one operand is required")
	ErrDivisionByZero            = errors.New("cannot divide by zero")
	ErrNegativeSquareRoot        = errors.New("cannot calculate the square root of a negative number")
	ErrInvalidResult             = errors.New("calculation produced an invalid numeric result")
)

func Add(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidBinaryOperandCount
	}

	return finiteResult(operands[0] + operands[1])
}

func Subtract(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidBinaryOperandCount
	}

	return finiteResult(operands[0] - operands[1])
}

func Multiply(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidBinaryOperandCount
	}

	return finiteResult(operands[0] * operands[1])
}

func Divide(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidBinaryOperandCount
	}

	if operands[1] == 0 {
		return 0, ErrDivisionByZero
	}

	return finiteResult(operands[0] / operands[1])
}

func SquareRoot(operands []float64) (float64, error) {
	if len(operands) != 1 {
		return 0, ErrInvalidUnaryOperandCount
	}

	if operands[0] < 0 {
		return 0, ErrNegativeSquareRoot
	}

	return finiteResult(math.Sqrt(operands[0]))
}

func Percentage(operands []float64) (float64, error) {
	if len(operands) != 1 {
		return 0, ErrInvalidUnaryOperandCount
	}

	return finiteResult(operands[0] / 100)
}

func Power(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidBinaryOperandCount
	}

	return finiteResult(math.Pow(operands[0], operands[1]))
}

func finiteResult(result float64) (float64, error) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrInvalidResult
	}

	return result, nil
}
