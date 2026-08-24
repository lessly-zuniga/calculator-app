import { act, renderHook, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type {
  BinaryOperation,
  CalculatorButtonAction,
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
import { useCalculator } from './useCalculator';

vi.mock('../services/calculatorApi', () => ({
  add: vi.fn(),
  subtract: vi.fn(),
  multiply: vi.fn(),
  divide: vi.fn(),
  squareRoot: vi.fn(),
  percentage: vi.fn(),
  power: vi.fn(),
}));

type CalculatorHook = ReturnType<typeof useCalculator>;
type HookResult = { current: CalculatorHook };

function press(result: HookResult, action: CalculatorButtonAction) {
  act(() => result.current.handleAction(action));
}

function enterNumber(result: HookResult, value: string) {
  for (const digit of value) {
    press(result, digit === '.' ? { type: 'decimal' } : { type: 'digit', value: digit });
  }
}

function selectOperation(result: HookResult, operation: BinaryOperation) {
  press(result, { type: 'operation', operation });
}

async function expectResult(result: HookResult, expected: number) {
  press(result, { type: 'equals' });
  await waitFor(() => expect(result.current.result).toBe(expected));
}

describe('useCalculator', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('accepts number input', () => {
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '123');

    expect(result.current.displayValue).toBe('123');
  });

  it('accepts decimal input', () => {
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '12.5');

    expect(result.current.displayValue).toBe('12.5');
  });

  it.each([
    {
      name: 'addition',
      operation: 'add' as const,
      calculate: add,
      left: '2',
      right: '3',
      answer: 5,
    },
    {
      name: 'subtraction',
      operation: 'subtract' as const,
      calculate: subtract,
      left: '10',
      right: '4',
      answer: 6,
    },
    {
      name: 'multiplication',
      operation: 'multiply' as const,
      calculate: multiply,
      left: '6',
      right: '7',
      answer: 42,
    },
    {
      name: 'division',
      operation: 'divide' as const,
      calculate: divide,
      left: '20',
      right: '4',
      answer: 5,
    },
    {
      name: 'exponentiation',
      operation: 'power' as const,
      calculate: power,
      left: '2',
      right: '3',
      answer: 8,
    },
  ])('performs $name through the API service', async (testCase) => {
    vi.mocked(testCase.calculate).mockResolvedValue(testCase.answer);
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, testCase.left);
    selectOperation(result, testCase.operation);
    enterNumber(result, testCase.right);
    await expectResult(result, testCase.answer);

    if (testCase.operation === 'add') {
      expect(add).toHaveBeenCalledWith(Number(testCase.left), Number(testCase.right));
    } else {
      expect(testCase.calculate).toHaveBeenCalledWith([
        Number(testCase.left),
        Number(testCase.right),
      ]);
    }
  });

  it('shows a division-by-zero error from the API service', async () => {
    vi.mocked(divide).mockRejectedValue(new Error('Cannot divide by zero'));
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '20');
    selectOperation(result, 'divide');
    enterNumber(result, '0');
    press(result, { type: 'equals' });

    await waitFor(() => expect(result.current.error).toBe('Cannot divide by zero'));
    expect(result.current.result).toBeNull();
  });

  it('performs square root through the API service', async () => {
    vi.mocked(squareRoot).mockResolvedValue(9);
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '81');
    press(result, { type: 'operation', operation: 'squareRoot' });

    await waitFor(() => expect(result.current.result).toBe(9));
    expect(squareRoot).toHaveBeenCalledWith([81]);
  });

  it('converts a standalone percentage through the API service', async () => {
    vi.mocked(percentage).mockResolvedValue(0.1);
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '10');
    press(result, { type: 'operation', operation: 'percentage' });

    await waitFor(() => expect(result.current.result).toBe(0.1));
    expect(percentage).toHaveBeenCalledWith([10]);
  });

  it.each([
    { operation: 'add' as const, answer: 220 },
    { operation: 'subtract' as const, answer: 180 },
  ])('applies a relative percentage for $operation', async ({ operation, answer }) => {
    vi.mocked(percentage).mockResolvedValue(0.1);
    vi.mocked(multiply).mockResolvedValue(20);
    vi.mocked(operation === 'add' ? add : subtract).mockResolvedValue(answer);
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '200');
    selectOperation(result, operation);
    enterNumber(result, '10');
    press(result, { type: 'operation', operation: 'percentage' });
    await waitFor(() => expect(multiply).toHaveBeenCalledWith([200, 0.1]));
    await waitFor(() => expect(result.current.loading).toBe(false));
    await expectResult(result, answer);

    expect(result.current.expression).toBe(`200 ${operation === 'add' ? '+' : '−'} 10%`);
    if (operation === 'add') {
      expect(add).toHaveBeenCalledWith(200, 20);
    } else {
      expect(subtract).toHaveBeenCalledWith([200, 20]);
    }
  });

  it.each([
    { operation: 'multiply' as const, calculate: multiply, answer: 20 },
    { operation: 'divide' as const, calculate: divide, answer: 2000 },
  ])('uses a decimal percentage for $operation', async ({ operation, calculate, answer }) => {
    vi.mocked(percentage).mockResolvedValue(0.1);
    vi.mocked(calculate).mockResolvedValue(answer);
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '200');
    selectOperation(result, operation);
    enterNumber(result, '10');
    press(result, { type: 'operation', operation: 'percentage' });
    await waitFor(() => expect(percentage).toHaveBeenCalledWith([10]));
    await waitFor(() => expect(result.current.loading).toBe(false));
    await expectResult(result, answer);

    expect(result.current.expression).toBe(`200 ${operation === 'multiply' ? '×' : '÷'} 10%`);
    expect(calculate).toHaveBeenCalledWith([200, 0.1]);
  });

  it('replaces a pending operator without changing the first operand', () => {
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '10');
    selectOperation(result, 'add');
    selectOperation(result, 'multiply');

    expect(result.current.displayValue).toBe('10');
    expect(result.current.expression).toBe('10 ×');
    expect(result.current.selectedOperation).toBe('multiply');
    expect(add).not.toHaveBeenCalled();
    expect(multiply).not.toHaveBeenCalled();
  });

  it('clears calculator state', () => {
    const { result } = renderHook(() => useCalculator());

    enterNumber(result, '10');
    selectOperation(result, 'add');
    enterNumber(result, '4');
    press(result, { type: 'clear' });

    expect(result.current.displayValue).toBe('0');
    expect(result.current.expression).toBe('');
    expect(result.current.selectedOperation).toBeNull();
    expect(result.current.result).toBeNull();
    expect(result.current.error).toBeNull();
  });
});
