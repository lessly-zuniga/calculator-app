package calculator

import (
	"errors"
	"math"
	"testing"
)

type operation func([]float64) (float64, error)

func TestCalculations(t *testing.T) {
	tests := []struct {
		name      string
		calculate operation
		operands  []float64
		want      float64
	}{
		{name: "add", calculate: Add, operands: []float64{2, 3}, want: 5},
		{name: "add negative values", calculate: Add, operands: []float64{-2, -3}, want: -5},
		{name: "add decimal values", calculate: Add, operands: []float64{1.5, 2.25}, want: 3.75},
		{name: "add zero", calculate: Add, operands: []float64{0, 3}, want: 3},
		{name: "subtract", calculate: Subtract, operands: []float64{10, 4}, want: 6},
		{name: "subtract negative values", calculate: Subtract, operands: []float64{-4, -10}, want: 6},
		{name: "subtract decimal values", calculate: Subtract, operands: []float64{5.5, 2.25}, want: 3.25},
		{name: "subtract zero", calculate: Subtract, operands: []float64{4, 0}, want: 4},
		{name: "multiply", calculate: Multiply, operands: []float64{6, 7}, want: 42},
		{name: "multiply negative value", calculate: Multiply, operands: []float64{-6, 7}, want: -42},
		{name: "multiply decimal values", calculate: Multiply, operands: []float64{1.5, 2.5}, want: 3.75},
		{name: "multiply by zero", calculate: Multiply, operands: []float64{7, 0}, want: 0},
		{name: "divide", calculate: Divide, operands: []float64{20, 4}, want: 5},
		{name: "divide negative value", calculate: Divide, operands: []float64{-20, 4}, want: -5},
		{name: "divide decimal values", calculate: Divide, operands: []float64{7.5, 2.5}, want: 3},
		{name: "divide zero", calculate: Divide, operands: []float64{0, 4}, want: 0},
		{name: "power", calculate: Power, operands: []float64{2, 3}, want: 8},
		{name: "power zero exponent", calculate: Power, operands: []float64{5, 0}, want: 1},
		{name: "power negative exponent", calculate: Power, operands: []float64{2, -2}, want: 0.25},
		{name: "power decimal exponent", calculate: Power, operands: []float64{4, 0.5}, want: 2},
		{name: "square root", calculate: SquareRoot, operands: []float64{81}, want: 9},
		{name: "square root zero", calculate: SquareRoot, operands: []float64{0}, want: 0},
		{name: "square root decimal", calculate: SquareRoot, operands: []float64{2.25}, want: 1.5},
		{name: "percentage", calculate: Percentage, operands: []float64{25}, want: 0.25},
		{name: "percentage negative value", calculate: Percentage, operands: []float64{-20}, want: -0.2},
		{name: "percentage decimal value", calculate: Percentage, operands: []float64{12.5}, want: 0.125},
		{name: "percentage zero", calculate: Percentage, operands: []float64{0}, want: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.calculate(test.operands)
			if err != nil {
				t.Fatalf("calculate() error = %v", err)
			}
			if got != test.want {
				t.Errorf("calculate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name      string
		calculate operation
		operands  []float64
		wantError error
	}{
		{name: "division by zero", calculate: Divide, operands: []float64{20, 0}, wantError: ErrDivisionByZero},
		{name: "negative square root", calculate: SquareRoot, operands: []float64{-1}, wantError: ErrNegativeSquareRoot},
		{name: "add invalid operand count", calculate: Add, operands: []float64{1}, wantError: ErrInvalidBinaryOperandCount},
		{name: "subtract invalid operand count", calculate: Subtract, operands: []float64{1}, wantError: ErrInvalidBinaryOperandCount},
		{name: "multiply invalid operand count", calculate: Multiply, operands: []float64{1}, wantError: ErrInvalidBinaryOperandCount},
		{name: "divide invalid operand count", calculate: Divide, operands: []float64{1}, wantError: ErrInvalidBinaryOperandCount},
		{name: "power invalid operand count", calculate: Power, operands: []float64{1}, wantError: ErrInvalidBinaryOperandCount},
		{name: "square root invalid operand count", calculate: SquareRoot, operands: []float64{1, 2}, wantError: ErrInvalidUnaryOperandCount},
		{name: "percentage invalid operand count", calculate: Percentage, operands: []float64{1, 2}, wantError: ErrInvalidUnaryOperandCount},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.calculate(test.operands)
			if !errors.Is(err, test.wantError) {
				t.Errorf("calculate() error = %v, want %v", err, test.wantError)
			}
		})
	}
}

func TestNonFiniteResults(t *testing.T) {
	tests := []struct {
		name      string
		calculate operation
		operands  []float64
	}{
		{name: "add overflow", calculate: Add, operands: []float64{math.MaxFloat64, math.MaxFloat64}},
		{name: "subtract overflow", calculate: Subtract, operands: []float64{-math.MaxFloat64, math.MaxFloat64}},
		{name: "multiply overflow", calculate: Multiply, operands: []float64{math.MaxFloat64, 2}},
		{name: "divide overflow", calculate: Divide, operands: []float64{math.MaxFloat64, math.SmallestNonzeroFloat64}},
		{name: "power infinity", calculate: Power, operands: []float64{10, 1000}},
		{name: "power NaN", calculate: Power, operands: []float64{-1, 0.5}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.calculate(test.operands)
			if !errors.Is(err, ErrInvalidResult) {
				t.Errorf("calculate() error = %v, want %v", err, ErrInvalidResult)
			}
		})
	}
}
