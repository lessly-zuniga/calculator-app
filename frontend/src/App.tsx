import { useState } from 'react';
import { add } from './services/calculatorApi';

function App() {
  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function testApi() {
    setResult(null);
    setError(null);

    try {
      setResult(await add(2, 3));
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Unable to call the API');
    }
  }

  return (
    <main>
      <button type="button" onClick={testApi}>
        Test API
      </button>
      {result !== null && <p>Result: {result}</p>}
      {error !== null && <p role="alert">{error}</p>}
    </main>
  );
}

export default App;
