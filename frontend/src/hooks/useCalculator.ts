import { useState } from 'react';
import type {
  BinaryOperation,
  CalculatorButtonAction,
  CalculatorOperation,
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

export function useCalculator() {
  const [currentInput, setCurrentInput] = useState('0');
  const [selectedOperation, setSelectedOperation] = useState<BinaryOperation | null>(null);
  const [firstOperand, setFirstOperand] = useState<number | null>(null);
  const [secondOperandStarted, setSecondOperandStarted] = useState(false);
  const [resolvedSecondOperand, setResolvedSecondOperand] = useState<number | null>(null);
  const [expression, setExpression] = useState('');
  const [result, setResult] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function enterDigit(digit: string) {
    const shouldStartSecondOperand = selectedOperation !== null && !secondOperandStarted;
    const shouldReplacePercentageOperand = resolvedSecondOperand !== null;
    setCurrentInput((value) =>
      shouldStartSecondOperand || shouldReplacePercentageOperand || value === '0' || result !== null
        ? digit
        : value + digit,
    );
    if (shouldReplacePercentageOperand && selectedOperation !== null && firstOperand !== null) {
      setExpression(`${firstOperand} ${operationLabels[selectedOperation]}`);
    }
    setResolvedSecondOperand(null);
    if (selectedOperation !== null) {
      setSecondOperandStarted(true);
    }
    setResult(null);
    setError(null);
  }

  function enterDecimal() {
    const shouldStartSecondOperand = selectedOperation !== null && !secondOperandStarted;
    const shouldReplacePercentageOperand = resolvedSecondOperand !== null;
    setCurrentInput((value) => {
      if (result !== null || shouldStartSecondOperand || shouldReplacePercentageOperand) {
        return '0.';
      }
      return value.includes('.') ? value : `${value}.`;
    });
    if (shouldReplacePercentageOperand && selectedOperation !== null && firstOperand !== null) {
      setExpression(`${firstOperand} ${operationLabels[selectedOperation]}`);
    }
    setResolvedSecondOperand(null);
    if (selectedOperation !== null) {
      setSecondOperandStarted(true);
    }
    setResult(null);
    setError(null);
  }

  function clear() {
    setCurrentInput('0');
    setSelectedOperation(null);
    setFirstOperand(null);
    setSecondOperandStarted(false);
    setResolvedSecondOperand(null);
    setExpression('');
    setResult(null);
    setError(null);
  }

  function selectBinaryOperation(operation: BinaryOperation) {
    if (selectedOperation !== null && firstOperand !== null && !secondOperandStarted) {
      setSelectedOperation(operation);
      setExpression(`${firstOperand} ${operationLabels[operation]}`);
      setCurrentInput(String(firstOperand));
      setResult(null);
      setError(null);
      return;
    }

    const operand = Number(currentInput);
    setFirstOperand(operand);
    setSelectedOperation(operation);
    setExpression(`${currentInput} ${operationLabels[operation]}`);
    setSecondOperandStarted(false);
    setResolvedSecondOperand(null);
    if (selectedOperation !== null) {
      setCurrentInput('0');
    }
    setResult(null);
    setError(null);
  }

  async function runSquareRoot() {
    const operand = Number(currentInput);
    setExpression(`${operationLabels.squareRoot}(${currentInput})`);
    await runCalculation(() => squareRoot([operand]));
  }

  async function applyPercentage() {
    const percentageInput = Number(currentInput);

    if (selectedOperation === null || firstOperand === null) {
      setExpression(`${currentInput}%`);
      await runCalculation(() => percentage([percentageInput]));
      return;
    }

    setExpression(`${firstOperand} ${operationLabels[selectedOperation]} ${currentInput}%`);
    setLoading(true);
    setError(null);

    try {
      const decimalPercentage = await percentage([percentageInput]);
      const effectiveOperand =
        selectedOperation === 'add' || selectedOperation === 'subtract'
          ? await multiply([firstOperand, decimalPercentage])
          : decimalPercentage;

      setResolvedSecondOperand(effectiveOperand);
      setSecondOperandStarted(true);
      setResult(null);
    } catch (calculationError) {
      setError(
        calculationError instanceof Error ? calculationError.message : 'Unable to calculate result',
      );
    } finally {
      setLoading(false);
    }
  }

  async function calculateResult() {
    if (selectedOperation === null || firstOperand === null) {
      return;
    }

    const secondOperand =
      resolvedSecondOperand ?? (secondOperandStarted ? Number(currentInput) : 0);
    if (resolvedSecondOperand === null) {
      setExpression(`${firstOperand} ${operationLabels[selectedOperation]} ${secondOperand}`);
    }
    await runCalculation(() => binaryCalculations[selectedOperation](firstOperand, secondOperand));
    setSelectedOperation(null);
    setFirstOperand(null);
    setSecondOperandStarted(false);
    setResolvedSecondOperand(null);
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
        if (action.operation === 'squareRoot') {
          void runSquareRoot();
        } else if (action.operation === 'percentage') {
          void applyPercentage();
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
