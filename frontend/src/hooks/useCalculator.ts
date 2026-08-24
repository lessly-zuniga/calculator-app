import { useState } from 'react';
import type {
  BinaryOperation,
  CalculatorButtonAction,
  CalculatorOperation,
  UnaryOperation,
} from '../components/calculator/calculator.types';
import {
  add,
  divide,
  multiply,
  percentage,
  power,
  squareRoot,
  subtract,
} from '../services/calculatorApi';

const operationLabels: Record<CalculatorOperation, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: '^',
  squareRoot: '√',
  percentage: '%',
};

const binaryCalculations: Record<
  BinaryOperation,
  (leftOperand: number, rightOperand: number) => Promise<number>
> = {
  add,
  subtract: (leftOperand, rightOperand) => subtract([leftOperand, rightOperand]),
  multiply: (leftOperand, rightOperand) => multiply([leftOperand, rightOperand]),
  divide: (leftOperand, rightOperand) => divide([leftOperand, rightOperand]),
  power: (leftOperand, rightOperand) => power([leftOperand, rightOperand]),
};

const unaryCalculations: Record<UnaryOperation, (operand: number) => Promise<number>> = {
  squareRoot: (operand) => squareRoot([operand]),
  percentage: (operand) => percentage([operand]),
};

export function useCalculator() {
  const [currentInput, setCurrentInput] = useState('0');
  const [selectedOperation, setSelectedOperation] = useState<BinaryOperation | null>(null);
  const [firstOperand, setFirstOperand] = useState<number | null>(null);
  const [expression, setExpression] = useState('');
  const [result, setResult] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function enterDigit(digit: string) {
    setCurrentInput((value) => (value === '0' || result !== null ? digit : value + digit));
    setResult(null);
    setError(null);
  }

  function enterDecimal() {
    setCurrentInput((value) => {
      if (result !== null) {
        return '0.';
      }
      return value.includes('.') ? value : `${value}.`;
    });
    setResult(null);
    setError(null);
  }

  function clear() {
    setCurrentInput('0');
    setSelectedOperation(null);
    setFirstOperand(null);
    setExpression('');
    setResult(null);
    setError(null);
  }

  function selectBinaryOperation(operation: BinaryOperation) {
    const operand = Number(currentInput);
    setFirstOperand(operand);
    setSelectedOperation(operation);
    setExpression(`${currentInput} ${operationLabels[operation]}`);
    setCurrentInput('0');
    setResult(null);
    setError(null);
  }

  async function runUnaryOperation(operation: UnaryOperation) {
    const operand = Number(currentInput);
    setExpression(`${operationLabels[operation]}(${currentInput})`);
    await runCalculation(() => unaryCalculations[operation](operand));
  }

  async function calculateResult() {
    if (selectedOperation === null || firstOperand === null) {
      return;
    }

    const secondOperand = Number(currentInput);
    setExpression(`${firstOperand} ${operationLabels[selectedOperation]} ${currentInput}`);
    await runCalculation(() => binaryCalculations[selectedOperation](firstOperand, secondOperand));
    setSelectedOperation(null);
    setFirstOperand(null);
  }

  async function runCalculation(calculate: () => Promise<number>) {
    setLoading(true);
    setError(null);

    try {
      const nextResult = await calculate();
      setResult(nextResult);
      setCurrentInput(String(nextResult));
    } catch (calculationError) {
      setError(
        calculationError instanceof Error ? calculationError.message : 'Unable to calculate result',
      );
    } finally {
      setLoading(false);
    }
  }

  function handleAction(action: CalculatorButtonAction) {
    switch (action.type) {
      case 'digit':
        enterDigit(action.value);
        break;
      case 'decimal':
        enterDecimal();
        break;
      case 'clear':
        clear();
        break;
      case 'operation':
        if (action.operation === 'squareRoot' || action.operation === 'percentage') {
          void runUnaryOperation(action.operation);
        } else {
          selectBinaryOperation(action.operation);
        }
        break;
      case 'equals':
        void calculateResult();
        break;
    }
  }

  return {
    currentInput,
    selectedOperation,
    expression,
    result,
    displayValue: result === null ? currentInput : String(result),
    loading,
    error,
    handleAction,
  };
}
