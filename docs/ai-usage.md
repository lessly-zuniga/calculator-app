# AI Usage

## AI Tooling Used

Codex and AI-assisted development were used during implementation and review of Calculator App. AI supported narrowly scoped coding, testing, review, infrastructure, and documentation tasks. It did not independently define the requirements, architecture, or acceptance criteria.

## Development Workflow

Development was intentionally incremental:

1. Define a small task with explicit boundaries.
2. Give AI a narrowly scoped instruction.
3. Inspect the generated changes.
4. Run the relevant application or service.
5. Manually exercise the behavior.
6. Correct issues before continuing.
7. Commit the completed unit of work.
8. Move to the next task.

The prompts below are concise, representative reconstructions of the actual instructions, not verbatim transcripts or the complete conversation history.

## Incremental Implementation and Validation

## Backend API Architecture

Backend requests follow a direct flow:

```text
HTTP request
    ↓
Handler
  - method enforcement
  - JSON decoding and transport validation
    ↓
Calculator domain
  - operand cardinality
  - arithmetic and operation-specific validation
  - domain errors and finite-result validation
    ↓
Handler
  - domain-error mapping
  - HTTP status and JSON response
```

Handlers own HTTP methods, JSON decoding, transport-boundary validation, domain-error mapping, and response encoding. The calculator package owns operand cardinality, arithmetic behavior, operation-specific validation, reusable domain errors, and finite-result checks. The domain package does not depend on `net/http`.

### 1. Backend foundation and health endpoint

**Endpoint**

`GET /api/v1/health`

**Contract**

This endpoint has no request body and is independent of calculator domain logic. A healthy service returns HTTP 200 with:

```json
{
  "status": "ok"
}
```

The current implementation accepts only GET. Other methods return HTTP 405, set `Allow: GET`, and return `METHOD_NOT_ALLOWED` with the message `Only GET requests are allowed`. The initial implementation did not enforce GET; that defect was found and corrected during the later backend consistency review.

**Representative Prompt**

> Implement the initial Go HTTP server and health endpoint only.
>
> Add `GET /api/v1/health`.
>
> Return JSON and HTTP 200 when the service is running. Keep the endpoint independent from calculator domain logic. Keep the implementation minimal and idiomatic. Do not add dependencies, calculator operations, or unnecessary abstractions.

**Manual Validation**

The API was run on port 8080 and exercised before calculator operations were added:

```http
GET /api/v1/health
```

The request returned HTTP 200 and:

```json
{
  "status": "ok"
}
```

### 2. Addition endpoint

**Endpoint**

`POST /api/v1/add`

**Contract**

Request:

```json
{
  "operands": [2, 3]
}
```

Success, HTTP 200:

```json
{
  "result": 5
}
```

Addition is a binary operation requiring exactly two numeric operands. Malformed JSON, unknown fields, missing operands, nonnumeric operands, extra JSON values, and incorrect cardinality are invalid requests. The calculator package computes the sum; the handler owns method enforcement, JSON decoding, error mapping, and JSON responses. This endpoint established the binary-operation structure reused by later handlers.

**Representative Prompt**

> Implement only the addition endpoint at `POST /api/v1/add`.
>
> Accept JSON and require exactly two numeric operands. Calculate the sum in the calculator/domain layer; keep HTTP parsing and response handling in the handler layer. Return the established JSON result or structured invalid-request error. Preserve `/api/v1/health`. Do not implement subtraction, multiplication, division, or advanced operations. Do not modify the frontend, add dependencies, or refactor unrelated code.
>
> Stop after addition so it can be manually verified before continuing.

**Manual Validation**

`POST /api/v1/add` with `{"operands":[2,3]}` returned HTTP 200 with `{"result":5}`. Negative, decimal, and zero cases were subsequently covered by automated domain tests.

### 3. Subtraction endpoint

**Endpoint**

`POST /api/v1/subtract`

**Contract**

Request:

```json
{
  "operands": [10, 4]
}
```

Success, HTTP 200:

```json
{
  "result": 6
}
```

Subtraction requires exactly two numeric operands. `operands[0]` is the minuend and `operands[1]` is the subtrahend, so order is preserved. The domain package performs subtraction; its handler enforces POST and the shared JSON contract.

**Representative Prompt**

> Implement only subtraction at `POST /api/v1/subtract`, following the architecture and request/response contract established by addition.
>
> Require exactly two numeric operands. Treat the first as the minuend and the second as the subtrahend. Keep arithmetic in the domain layer and HTTP concerns in the handler. Preserve health and addition behavior. Do not implement later operations, modify the frontend, add dependencies, or refactor unrelated code.

**Manual Validation**

`POST /api/v1/subtract` with `{"operands":[10,4]}` returned HTTP 200 with `{"result":6}`. This ordered case was subsequently covered by domain and HTTP tests. The documentation does not claim that a separate reversed-operand request was manually executed.

### 4. Multiplication endpoint

**Endpoint**

`POST /api/v1/multiply`

**Contract**

Request:

```json
{
  "operands": [6, 7]
}
```

Success, HTTP 200:

```json
{
  "result": 42
}
```

Multiplication is POST-only and requires exactly two numeric operands. The calculator package multiplies them; the handler parses JSON and writes the structured response. Positive, negative, zero, and decimal values are supported when the result is finite.

**Representative Prompt**

> Implement only multiplication at `POST /api/v1/multiply` using the existing binary-operation architecture.
>
> Accept JSON, enforce POST, and require exactly two numeric operands. Keep multiplication in the calculator/domain package and HTTP parsing and JSON responses in the handler. Preserve every existing endpoint. Do not modify the frontend, add dependencies, implement other operations, or refactor unrelated code.

**Manual Validation**

`POST /api/v1/multiply` with `{"operands":[6,7]}` returned HTTP 200 with `{"result":42}`. Negative, zero, and decimal cases were subsequently covered by automated domain tests.

### 5. Division endpoint

**Endpoint**

`POST /api/v1/divide`

**Contract**

Request:

```json
{
  "operands": [20, 4]
}
```

Success, HTTP 200:

```json
{
  "result": 5
}
```

Division requires exactly two numeric operands. `operands[0]` is the dividend and `operands[1]` is the divisor; their order is preserved. A zero divisor produces the reusable domain error `ErrDivisionByZero`, which the HTTP handler maps to HTTP 422 and a structured `DIVISION_BY_ZERO` response.

**Representative Prompt**

> Implement only division at `POST /api/v1/divide`.
>
> Require exactly two numeric operands. Treat the first as the dividend and the second as the divisor, preserving order. Handle division by zero explicitly: it must not produce infinity, NaN, a crash, or invalid JSON. Represent it as a domain error and translate that error into the structured HTTP response at the handler boundary. Preserve existing operations. Do not modify the frontend, add dependencies, implement unrelated operations, or refactor working code.

**Manual Validation**

`POST /api/v1/divide` with `{"operands":[20,4]}` returned HTTP 200 with `{"result":5}`. The documented error-path check used `{"operands":[20,0]}` and returned HTTP 422:

```json
{
  "error": {
    "code": "DIVISION_BY_ZERO",
    "message": "Cannot divide by zero"
  }
}
```

This confirmed that the failure crossed the domain/HTTP boundary correctly instead of leaking a non-finite value.

### 6. Exponentiation endpoint

**Endpoint**

`POST /api/v1/power`

**Contract**

Request:

```json
{
  "operands": [2, 3]
}
```

Success, HTTP 200:

```json
{
  "result": 8
}
```

Power requires exactly two numeric operands. `operands[0]` is the base and `operands[1]` is the exponent, so order matters. `math.Pow` is used in the domain layer, and non-finite results are rejected before the handler writes HTTP 200.

**Representative Prompt**

> Implement only exponentiation at `POST /api/v1/power`.
>
> Require exactly two numeric operands. Treat the first as the base and the second as the exponent. Preserve order, reuse the domain/handler separation, and protect the API from non-finite results. Preserve existing operations. Do not modify frontend behavior, add dependencies, implement other operations, or refactor unrelated code.

**Manual Validation**

`POST /api/v1/power` with `{"operands":[2,3]}` returned HTTP 200 with `{"result":8}`. Zero, negative, and valid decimal exponent behavior was subsequently covered by automated domain tests.

### 7. Square root endpoint

**Endpoint**

`POST /api/v1/square-root`

**Contract**

Request:

```json
{
  "operands": [81]
}
```

Success, HTTP 200:

```json
{
  "result": 9
}
```

Square root is unary and requires exactly one numeric operand. A negative operand produces `ErrNegativeSquareRoot` in the domain package. The handler translates it into HTTP 422 and `NEGATIVE_SQUARE_ROOT`, without placing HTTP concepts in the calculator package.

**Representative Prompt**

> Implement only square root at `POST /api/v1/square-root`.
>
> Treat it as a unary operation and require exactly one numeric operand. Reject negative values using a domain error. Keep HTTP behavior out of the calculator package and translate the domain error into a structured handler response. Preserve all endpoints. Do not modify the frontend, add dependencies, implement other operations, or refactor unrelated code.

**Manual Validation**

`POST /api/v1/square-root` with `{"operands":[81]}` returned HTTP 200 with `{"result":9}`. The documented negative-input check returned HTTP 422:

```json
{
  "error": {
    "code": "NEGATIVE_SQUARE_ROOT",
    "message": "Cannot calculate the square root of a negative number"
  }
}
```

This endpoint introduced unary operand cardinality, whose error wording was corrected during the later consistency review.

### 8. Percentage endpoint

**Endpoint**

`POST /api/v1/percentage`

**Contract**

Request:

```json
{
  "operands": [25]
}
```

Success, HTTP 200:

```json
{
  "result": 0.25
}
```

Percentage is unary and requires exactly one numeric operand. The domain operation returns `operands[0] / 100`. It is intentionally independent from pending operators and frontend state.

**Representative Prompt**

> Implement only the backend percentage operation at `POST /api/v1/percentage`.
>
> Treat it as unary, require exactly one numeric operand, and convert it with `value / 100`. Keep it independent from UI state and pending calculator operators; do not implement contextual percentage semantics in the backend. Preserve the API architecture and every existing endpoint. Do not modify the frontend yet, add dependencies, implement unrelated behavior, or refactor working code.

**Manual Validation**

`POST /api/v1/percentage` with `{"operands":[25]}` returned HTTP 200 with `{"result":0.25}`. Negative, decimal, and zero behavior was subsequently covered by automated domain tests.

The endpoint does not know whether the user is entering `200 + 10%`, `200 - 10%`, `200 × 10%`, or `200 ÷ 10%`. That context was added later in frontend orchestration while final arithmetic remained backend-driven.

#### Design Boundary

The responsibilities remain explicit:

```text
Backend percentage:       percentage(10) → 0.1
Frontend orchestration:   for 200 - 10%, determine that the contextual amount is 20
Backend arithmetic:       subtract(200, 20) → 180
```

The frontend coordinates existing API operations; it does not replace the backend calculator domain with a second arithmetic implementation.

## Shared API Conventions Established Incrementally

Implementing one endpoint at a time gradually established a consistent contract:

- Health is GET-only; all calculator endpoints are POST-only.
- Handlers return `Content-Type: application/json` for success and error responses.
- Calculator requests use `{"operands":[...]}` and successes use `{"result":number}`.
- Errors use `{"error":{"code":string,"message":string}}`.
- Add, subtract, multiply, divide, and power require exactly two operands.
- Square root and percentage require exactly one operand.
- JSON decoding rejects malformed input, unknown fields, nonnumeric operands, and extra JSON values.
- Domain code validates operand cardinality and calculation-specific failures without importing HTTP concepts.
- Handlers map domain errors to HTTP responses and validate methods, requests, and response encoding.
- All operations reject NaN and positive or negative infinity through `ErrInvalidResult` before success is written.

| Scenario                                    | HTTP status | API behavior                                                                                                                                              |
| ------------------------------------------- | ----------: | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Successful health or calculation request    |         200 | `{"status":"ok"}` or `{"result":number}`                                                                                                                  |
| Unsupported method                          |         405 | `METHOD_NOT_ALLOWED`; calculator message is `Only POST requests are allowed`, health message is `Only GET requests are allowed`; an `Allow` header is set |
| Malformed or otherwise invalid JSON request |         400 | `INVALID_REQUEST` with the operation's operand-count message                                                                                              |
| Invalid binary operand count                |         400 | `INVALID_REQUEST`: `Exactly two numeric operands are required`                                                                                            |
| Invalid unary operand count                 |         400 | `INVALID_REQUEST`: `Exactly one numeric operand is required`                                                                                              |
| Division by zero                            |         422 | `DIVISION_BY_ZERO`: `Cannot divide by zero`                                                                                                               |
| Negative square root                        |         422 | `NEGATIVE_SQUARE_ROOT`: `Cannot calculate the square root of a negative number`                                                                           |
| Non-finite calculation result               |         422 | `INVALID_RESULT`: `Calculation produced an invalid numeric result`                                                                                        |

## Endpoint Implementation Discipline

The API was deliberately not generated as one large change. The progression was:

```text
Health → manual verification
Addition → manual verification → commit
Subtraction → manual verification → commit
Multiplication → manual verification → commit
Division → happy-path and error-path verification → commit
Exponentiation → manual verification → commit
Square root → happy-path and domain-error verification → commit
Percentage → manual verification → commit
Backend consistency review → fixes → automated tests → handler coverage expansion
```

This produced small diffs, explicit scope, simpler review, easier regression isolation, clear commit boundaries, lower AI context requirements, and less risk of unrelated generated changes. Each capability was exercised before the next one was added.

### 9. Backend review before testing

**Representative Prompt**

> Review the backend without modifying files. Inspect route registration, HTTP methods, request and operand validation, error consistency, JSON responses, status codes, and edge cases. Report only concrete defects.

**Validation**

The review identified three real issues:

- `/api/v1/health` did not enforce GET.
- Unary operations reused binary operand-count wording.
- NaN or infinity could reach JSON serialization after HTTP 200 had already been selected.

These were corrected before testing expanded. Non-finite results now return HTTP 422 with code `INVALID_RESULT` and message `Calculation produced an invalid numeric result`.

### 10. Backend automated tests

**Representative Prompt**

> Add calculator domain tests after manual verification, then add HTTP handler tests using the Go standard library and `net/http/httptest`. Cover valid requests, unsupported methods, malformed JSON, invalid operand counts, division by zero, negative square root, and non-finite results. Use table-driven tests where useful.

**Validation**

```sh
gofmt -w .
go vet ./...
go test -v ./...
go test -cover ./...
```

Current measured statement coverage is 72.1% for `cmd/server`, 100.0% for `internal/calculator`, and 84.8% for `internal/handler`. These are package-level values, not a claim of 100% overall coverage.

### 11. Frontend API client

**Representative Prompt**

> Implement only the TypeScript API layer using native `fetch` and relative `/api/...` URLs. Centralize request handling, preserve structured backend errors, do not perform final arithmetic locally, do not modify the calculator UI, and do not add dependencies.

**Validation**

The client was connected to the running API through Vite's proxy. A temporary `Test API` button called `add(2, 3)` and displayed the result or API error, allowing the browser request to be functionally inspected before the full UI was built.

### 12. Frontend component architecture

**Representative Prompt**

> Structure the frontend around `AppShell`, `Calculator`, `CalculatorDisplay`, `CalculatorKeypad`, `CalculatorButton`, and `useCalculator`. Keep components small, keep networking outside presentation components, use local state rather than a global state library, and avoid unnecessary dependencies.

**Validation**

The application was built and rendered after the structural change. Component composition and API connectivity were inspected before interaction behavior and visual styling were expanded.

### 13. Calculator UI

**Representative Prompt**

> Implement a polished responsive calculator inspired by the simplicity, proportions, hierarchy, and interaction quality of Apple's Calculator while keeping the design original. Preserve the architecture and logic. Support desktop, tablet, and mobile with semantic controls, visible focus states, large touch targets, and restrained feedback.

**Validation**

The UI was manually inspected on mobile, tablet, desktop, and constrained-height tablet and desktop viewports. Checks included control visibility, touch targets, focus, display overflow, grid alignment, and horizontal scrolling. This exposed vertical clipping on shorter viewports.

### 14. Short-viewport responsive fix

**Representative Prompt**

> Improve responsive behavior for short tablet and desktop viewports. Preserve mobile, prevent vertical clipping, adapt spacing, display height, buttons, and gaps to available height, and prefer CSS `clamp()`, `100dvh`, and height-based media queries over JavaScript calculations.

**Validation**

The viewport was repeatedly resized across short laptop, landscape tablet, split-screen, and non-maximized desktop dimensions until all controls remained visible and usable. Extremely constrained layouts scroll instead of clipping.

### 15. Removal of ± and keypad adjustment

**Representative Prompt**

> Remove the `±` feature completely, including interaction logic and dead code. Preserve the layout and do not modify backend behavior.

The visual follow-up was:

> Move `√x` into the main keypad and make the first row `AC | % | √x | ÷`. Give AC the same circular dimensions as other buttons, remove the standalone square-root control, and preserve behavior.

**Validation**

The keypad was checked across responsive layouts. The final grid has one square-root control, no empty slot, no sign-toggle control, and a balanced four-column first row.

### 16. Operator replacement bug

**Representative Prompt**

> If an operator is pending and no second operand has been entered, selecting another binary operator must replace it, preserve the first operand, avoid calculation, and avoid an API request merely because the operator changed.

**Validation**

Manual testing reproduced the defect: `10`, `+`, `×` changed the display to `0`. After the fix, changing `+` to `×` leaves it at `10`, and `10 + × 4 =` evaluates as `10 × 4`, returning `40` through the API.

### 17. Contextual percentage behavior

**Representative Prompt**

> Use contextual percentage semantics: `10 % → 0.1`, `200 + 10 % → 220`, `200 - 10 % → 180`, `200 × 10 % → 20`, and `200 ÷ 10 % → 2000`. Preserve the percentage expression, keep final arithmetic backend-driven, and do not change API contracts.

**Validation**

Those exact UI sequences were checked. `/api/v1/percentage` remains a unary decimal conversion. The frontend determines context, while percentage conversion and resulting arithmetic continue to use backend calls.

### 18. Frontend testing

**Representative Prompt**

> Add minimal Vitest infrastructure and focused behavior tests rather than CSS or snapshots. Cover number input, decimals, arithmetic, API errors, square root, exponentiation, standalone and contextual percentages, operator replacement, and reset.

**Validation**

```sh
npm test
npm run test:coverage
npm run lint
npm run build
```

The current suite has 16 passing hook tests. Measured `useCalculator.ts` coverage is 92.3% statements, 77.61% branches, 100% functions, and 92.3% lines.

### 19. Backend handler coverage improvement

**Representative Prompt**

> Inspect handlers and tests, identify meaningful uncovered HTTP paths, and increase coverage with focused `net/http/httptest` cases. Do not duplicate exhaustive arithmetic tests or chase 100% only as a metric.

**Validation**

Tests covered each endpoint with representative success, method, malformed-body, operand-count, domain-error, invalid-result, JSON-shape, and content-type cases where applicable. `go test -v ./...` and `go test -cover ./...` pass; current handler coverage is 84.8%.

### 20. Backend Dockerization

**Representative Prompt**

> Dockerize only the Go backend with a multi-stage build. Build the binary in the build stage, use a small runtime image, copy only required runtime files, expose port 8080, and preserve behavior.

**Validation**

```sh
docker build -t calculator-api ./backend
docker run --rm -p 8080:8080 calculator-api
```

With the container running, the health endpoint and a calculator POST endpoint were exercised through `localhost:8080`.

### 21. Frontend Dockerization

**Representative Prompt**

> Dockerize the React frontend with a Node build stage and Nginx runtime. Serve `dist`, proxy `/api/` to `http://backend:8080`, preserve SPA fallback, and do not change calculator behavior.

**Validation**

The image build verified `npm ci` and the Vite production build. Running it without Compose exposed the expected constraint that `backend` is a Docker service hostname, resolvable only on a shared Docker network. Full proxy validation therefore used Compose.

### 22. Docker Compose

**Representative Prompt**

> Add backend and frontend Compose services on the default network. Keep the backend internal, map frontend to port 3000, add a backend healthcheck, wait for backend health before frontend startup, and add no unnecessary networks, volumes, or variables.

**Validation**

```sh
docker compose config
docker compose build
docker compose up --build
```

The app was opened at `http://localhost:3000`. `GET http://localhost:3000/api/v1/health` verified the Nginx proxy path, and calculator operations were exercised through the Dockerized UI. The stack stops with `docker compose down`.

### 23. GitHub Actions CI

**Representative Prompt**

> Create separate frontend and backend jobs. Frontend runs `npm ci`, lint, formatting, coverage tests, and build. Backend verifies `gofmt`, then runs `go vet`, tests, and coverage. Use official actions; do not deploy or add secrets.

**Validation**

The workflow was pushed and inspected in GitHub Actions. Backend passed. Frontend initially reported Prettier failures for `src/main.tsx` and `tsconfig.json`; those were corrected in a follow-up commit and CI was run again. The repository verifies the workflow and formatting follow-up but does not contain authoritative evidence of the final hosted run status.

## Developer-Directed Decisions

The developer explicitly defined and reviewed the repository structure, task order, frontend/backend separation, API contracts, arithmetic semantics, contextual percentages, status and error behavior, component responsibilities, responsive direction, testing scope, Docker topology, CI checks, and commit boundaries.

## Verification Process

AI-generated changes were not accepted based only on generated code. Each unit followed:

```text
AI instruction
→ code review
→ local execution
→ manual behavior check
→ automated validation when applicable
→ commit
```

Verification included `curl`, browser interaction and network inspection, `gofmt`, `go vet`, Go domain and HTTP handler tests, Go coverage, ESLint, Prettier, Vitest, frontend coverage, Vite builds, Docker image builds, Docker Compose, and GitHub Actions. Failures and manually discovered issues became separate, narrowly scoped follow-up tasks rather than being hidden in broad refactors.

## Final Note

AI was used as an implementation and review accelerator. Architecture, expected behavior, scope, acceptance criteria, and final engineering decisions were explicitly defined and reviewed by the developer. Generated changes were validated through executable checks and manual behavior testing before being accepted.
