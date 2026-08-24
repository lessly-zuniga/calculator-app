import { useCalculator } from '../../hooks/useCalculator';
import { CalculatorDisplay } from './CalculatorDisplay';
import { CalculatorKeypad } from './CalculatorKeypad';

export function Calculator() {
  const { displayValue, expression, loading, error, handleAction } = useCalculator();

  return (
    <section className="calculator" aria-label="Calculator">
      <CalculatorDisplay
        expression={expression}
        value={displayValue}
        loading={loading}
        error={error}
      />
      <CalculatorKeypad disabled={loading} onAction={handleAction} />
    </section>
  );
}
