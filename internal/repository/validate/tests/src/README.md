# External Website Testing Project

This project contains automated tests for external websites using Playwright and TypeScript.

## Setup

1. Install dependencies:

```bash
pnpm install
```

2. Install Playwright browsers:

```bash
pnpm exec playwright install
```

## Running Tests

- Run all tests:

```bash
pnpm test
```

- Run tests with UI mode (interactive):

```bash
pnpm test:ui
```

- Run tests in headed mode (visible browser):

```bash
pnpm test:headed
```

- View test report:

```bash
pnpm report
```

## Project Structure

- `tests/` - Contains all test files
- `playwright.config.ts` - Playwright configuration
- `tsconfig.json` - TypeScript configuration

## Adding New Tests

1. Create a new test file in the `tests` directory
2. Import the test and expect functions from '@playwright/test'
3. Write your test cases using the Playwright API

Example:

```typescript
import { test, expect } from "@playwright/test";

test("my test", async ({ page }) => {
  await page.goto("https://example.com");
  await expect(page).toHaveTitle(/Example/);
});
```

## Best Practices

1. Use meaningful test descriptions
2. Group related tests using `test.describe()`
3. Use appropriate assertions
4. Handle loading states and timeouts appropriately
5. Take screenshots and videos on failure for debugging
