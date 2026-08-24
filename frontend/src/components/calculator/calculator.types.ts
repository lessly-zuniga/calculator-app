export type CalculatorButtonVariant = 'number' | 'function' | 'operator' | 'equals';

export type BinaryOperation = 'add' | 'subtract' | 'multiply' | 'divide' | 'power';
export type UnaryOperation = 'squareRoot' | 'percentage';
export type CalculatorOperation = BinaryOperation | UnaryOperation;

export type CalculatorButtonAction =
  | { type: 'digit'; value: string }
  | { type: 'decimal' }
  | { type: 'clear' }
  | { type: 'operation'; operation: CalculatorOperation }
  | { type: 'equals' };

export type CalculatorButtonDefinition = {
  label: string;
  accessibleLabel?: string;
  variant: CalculatorButtonVariant;
  action?: CalculatorButtonAction;
  disabled?: boolean;
};
