import { test, expect } from "@playwright/test";

test.describe("Google Website Tests", () => {
  test("should load Google homepage", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveTitle(/Google/);
  });

  test.skip("should perform a search", async ({ page }) => {
    await page.goto("/");
    await page.fill('textarea[name="q"]', "Playwright testing");
    await page.press('textarea[name="q"]', "Enter");

    // Wait for search results
    await page.waitForSelector("#search");

    // Verify search results are displayed
    const searchResults = await page.locator("#search").count();
    expect(searchResults).toBeGreaterThan(0);
  });

  test.skip("should have working navigation links", async ({ page }) => {
    await page.goto("/");

    // Test Gmail link
    await page.click('a[href="https://mail.google.com"]');
    await expect(page).toHaveURL(/mail\.google\.com/);

    // Go back to Google
    await page.goto("/");

    // Test Images link
    await page.click('a[href="https://www.google.com/imghp"]');
    await expect(page).toHaveURL(/google\.com\/imghp/);
  });
});
