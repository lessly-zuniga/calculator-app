export type CalculatorRequest = {
  operands: number[];
};

export type CalculatorSuccessResponse = {
  result: number;
};

export type CalculatorError = {
  code: string;
  message: string;
};

export type CalculatorErrorResponse = {
  error: CalculatorError;
};
