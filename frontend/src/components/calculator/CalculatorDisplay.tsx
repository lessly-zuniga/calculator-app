type CalculatorDisplayProps = {
  expression: string;
  value: string;
  loading: boolean;
  error: string | null;
};

export function CalculatorDisplay({ expression, value, loading, error }: CalculatorDisplayProps) {
  return (
    <section className="calculator-display" aria-label="Calculator display">
      <div className="calculator-display__expression">{expression || '\u00a0'}</div>
      <output className="calculator-display__value" aria-live="polite">
        {loading ? 'Calculating…' : value}
      </output>
      {error && (
        <p className="calculator-display__error" role="alert">
          {error}
        </p>
      )}
    </section>
  );
}
