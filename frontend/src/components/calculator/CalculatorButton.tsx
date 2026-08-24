import type { CalculatorButtonVariant } from './calculator.types';

type CalculatorButtonProps = {
  label: string;
  accessibleLabel?: string;
  variant: CalculatorButtonVariant;
  disabled?: boolean;
  onPress: () => void;
};

export function CalculatorButton({
  label,
  accessibleLabel,
  variant,
  disabled = false,
  onPress,
}: CalculatorButtonProps) {
  return (
    <button
      type="button"
      className={`calculator-button calculator-button--${variant}`}
      aria-label={accessibleLabel}
      disabled={disabled}
      onClick={onPress}
    >
      {label}
    </button>
  );
}
