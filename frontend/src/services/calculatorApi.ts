import type {
  CalculatorErrorResponse,
  CalculatorRequest,
  CalculatorSuccessResponse,
} from '../types/calculator';

export class CalculatorApiError extends Error {
  readonly code: string;
  readonly status: number;

  constructor(code: string, message: string, status: number) {
    super(message);
    this.name = 'CalculatorApiError';
    this.code = code;
    this.status = status;
  }
}

async function calculate(endpoint: string, operands: number[]): Promise<number> {
  const request: CalculatorRequest = { operands };
  const response = await fetch(endpoint, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });

  const payload = await parseJSON(response);

  if (!response.ok) {
    if (isErrorResponse(payload)) {
      throw new CalculatorApiError(payload.error.code, payload.error.message, response.status);
    }

    throw invalidResponseError(response.status);
  }

  if (!isSuccessResponse(payload)) {
    throw invalidResponseError(response.status);
  }

  return payload.result;
}

async function parseJSON(response: Response): Promise<unknown> {
  try {
    return await response.json();
  } catch {
    throw invalidResponseError(response.status);
  }
}

function isSuccessResponse(value: unknown): value is CalculatorSuccessResponse {
  return isRecord(value) && typeof value.result === 'number' && Number.isFinite(value.result);
}

function isErrorResponse(value: unknown): value is CalculatorErrorResponse {
  if (!isRecord(value) || !isRecord(value.error)) {
    return false;
  }

  return typeof value.error.code === 'string' && typeof value.error.message === 'string';
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function invalidResponseError(status: number): CalculatorApiError {
  return new CalculatorApiError(
    'INVALID_RESPONSE',
    'The server returned an invalid response',
    status,
  );
}

export function add(leftOperand: number, rightOperand: number): Promise<number> {
  return calculate('/api/v1/add', [leftOperand, rightOperand]);
}

export function subtract(operands: number[]): Promise<number> {
  return calculate('/api/v1/subtract', operands);
}

export function multiply(operands: number[]): Promise<number> {
  return calculate('/api/v1/multiply', operands);
}

export function divide(operands: number[]): Promise<number> {
  return calculate('/api/v1/divide', operands);
}

export function squareRoot(operands: number[]): Promise<number> {
  return calculate('/api/v1/square-root', operands);
}

export function percentage(operands: number[]): Promise<number> {
  return calculate('/api/v1/percentage', operands);
}

export function power(operands: number[]): Promise<number> {
  return calculate('/api/v1/power', operands);
}
