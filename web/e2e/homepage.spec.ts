import { test, expect } from './test-setup';

test.describe('Homepage', () => {
  test('should display the homepage with tenant list', async ({ page, testSetup }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Check that the page title is correct
    await expect(page).toHaveTitle(/NeuralLog Admin/);

    // Check that the header is displayed
    await expect(page.locator('h1').filter({ hasText: 'NeuralLog Admin Dashboard' })).toBeVisible();

    // Check that the tenant list is displayed
    await expect(page.locator('h2')).toContainText('Tenants');

    // Check that the tenants are displayed in the table
    const tenantRows = page.locator('table tbody tr');
    await expect(tenantRows).toHaveCount(2);

    // Check the first tenant details
    const firstRow = tenantRows.nth(0);
    await expect(firstRow.locator('td').nth(0)).toContainText('test-tenant-1');
    await expect(firstRow.locator('td').nth(1)).toContainText('Test Tenant 1');
    await expect(firstRow.locator('td').nth(3).locator('span')).toContainText('Running');

    // Check the second tenant details
    const secondRow = tenantRows.nth(1);
    await expect(secondRow.locator('td').nth(0)).toContainText('test-tenant-2');
    await expect(secondRow.locator('td').nth(1)).toContainText('Test Tenant 2');
    await expect(secondRow.locator('td').nth(3).locator('span')).toContainText('Degraded');
  });

  test('should have working navigation', async ({ page, testSetup }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Click on the "New Tenant" button
    await page.click('text=New Tenant');

    // Check that we're on the new tenant page
    await expect(page).toHaveURL(/\/tenants\/new/);

    // Go back to the homepage
    await page.goto('/');

    // Click on the "Edit" button for the first tenant
    await page.locator('table tbody tr').nth(0).locator('text=Edit').click();

    // Check that we're on the edit tenant page
    await expect(page).toHaveURL(/\/tenants\/test-tenant-1\/edit/);
  });

  test('should show empty state when no tenants exist', async ({ page }) => {
    // This test would require deleting all tenants, which we don't want to do in a real environment
    // Instead, we'll skip this test in real E2E mode
    test.skip(process.env.CI === 'true', 'Skipping in CI environment');

    // Navigate to the homepage
    await page.goto('/');

    // Check that the empty state message is displayed
    await expect(page.locator('text=No tenants found')).toBeVisible();
    await expect(page.locator('text=Create your first tenant to get started')).toBeVisible();
  });

  test('should refresh tenant list', async ({ page, testSetup }) => {
    // Navigate to the homepage
    await page.goto('/');

    // We'll rely on the test setup to have created the tenants

    // Click the refresh button
    await page.click('button:has-text("Refresh")');

    // Check that the new tenant is displayed
    const tenantRows = page.locator('table tbody tr');
    await expect(tenantRows).toHaveCount(3);

    // Check the new tenant details
    const newRow = tenantRows.nth(2);
    await expect(newRow.locator('td').nth(0)).toContainText('new-tenant');
    await expect(newRow.locator('td').nth(1)).toContainText('New Tenant');
  });
});
