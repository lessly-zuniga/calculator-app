package calculator

import "errors"

var ErrInvalidOperandCount = errors.New("exactly two operands are required")

func Add(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, ErrInvalidOperandCount
	}

	return operands[0] + operands[1], nil
}
