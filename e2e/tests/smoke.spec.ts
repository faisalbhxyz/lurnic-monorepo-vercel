import { test, expect } from '@playwright/test';

test('opens the API base URL (headed friendly)', async ({ page, baseURL }) => {
  // We intentionally don't assert a 200 here because the server may not be running,
  // may redirect, or may not serve "/" at all. This test is primarily to visually
  // observe actions in headed mode when the server is available.
  const response = await page.goto(baseURL ?? '/', { waitUntil: 'domcontentloaded' });

  // When the server isn't reachable, Playwright throws before returning a response.
  expect(response).not.toBeNull();

  // Small, visible interactions (useful when running headed with slowMo).
  await page.mouse.move(120, 120);
  await page.waitForTimeout(250);
  await page.mouse.move(420, 220);
});

