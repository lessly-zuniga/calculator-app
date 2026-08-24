import { CalculatorButton } from './CalculatorButton';
import type { CalculatorButtonAction, CalculatorButtonDefinition } from './calculator.types';

type CalculatorKeypadProps = {
  disabled: boolean;
  onAction: (action: CalculatorButtonAction) => void;
};

const buttons: CalculatorButtonDefinition[] = [
  { label: 'C', accessibleLabel: 'Clear', variant: 'function', action: { type: 'clear' } },
  {
    label: '√',
    accessibleLabel: 'Square root',
    variant: 'function',
    action: { type: 'operation', operation: 'squareRoot' },
  },
  {
    label: '%',
    accessibleLabel: 'Percentage',
    variant: 'function',
    action: { type: 'operation', operation: 'percentage' },
  },
  {
    label: '÷',
    accessibleLabel: 'Divide',
    variant: 'operator',
    action: { type: 'operation', operation: 'divide' },
  },
  { label: '7', variant: 'number', action: { type: 'digit', value: '7' } },
  { label: '8', variant: 'number', action: { type: 'digit', value: '8' } },
  { label: '9', variant: 'number', action: { type: 'digit', value: '9' } },
  {
    label: '×',
    accessibleLabel: 'Multiply',
    variant: 'operator',
    action: { type: 'operation', operation: 'multiply' },
  },
  { label: '4', variant: 'number', action: { type: 'digit', value: '4' } },
  { label: '5', variant: 'number', action: { type: 'digit', value: '5' } },
  { label: '6', variant: 'number', action: { type: 'digit', value: '6' } },
  {
    label: '−',
    accessibleLabel: 'Subtract',
    variant: 'operator',
    action: { type: 'operation', operation: 'subtract' },
  },
  { label: '1', variant: 'number', action: { type: 'digit', value: '1' } },
  { label: '2', variant: 'number', action: { type: 'digit', value: '2' } },
  { label: '3', variant: 'number', action: { type: 'digit', value: '3' } },
  {
    label: '+',
    accessibleLabel: 'Add',
    variant: 'operator',
    action: { type: 'operation', operation: 'add' },
  },
  {
    label: 'xʸ',
    accessibleLabel: 'Power',
    variant: 'operator',
    action: { type: 'operation', operation: 'power' },
  },
  { label: '0', variant: 'number', action: { type: 'digit', value: '0' } },
  { label: '.', accessibleLabel: 'Decimal point', variant: 'number', action: { type: 'decimal' } },
  { label: '=', accessibleLabel: 'Equals', variant: 'equals', action: { type: 'equals' } },
];

export function CalculatorKeypad({ disabled, onAction }: CalculatorKeypadProps) {
  return (
    <div className="calculator-keypad" aria-label="Calculator keypad">
      {buttons.map((button) => (
        <CalculatorButton
          key={button.label}
          label={button.label}
          accessibleLabel={button.accessibleLabel}
          variant={button.variant}
          disabled={disabled}
          onPress={() => onAction(button.action)}
        />
      ))}
    </div>
  );
}
