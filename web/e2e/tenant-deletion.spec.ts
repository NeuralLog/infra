import { test, expect } from './test-setup';

test.describe('Tenant Deletion', () => {
  test('should delete a tenant', async ({ page }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Get the initial count of tenants
    const initialTenantCount = await page.locator('table tbody tr').count();

    // Click the delete button for the first tenant
    await page.locator('table tbody tr').nth(0).locator('text=Delete').click();

    // Check that the confirmation dialog is displayed
    await expect(page.locator('div[role="alertdialog"]')).toBeVisible();
    await expect(page.locator('div[role="alertdialog"]')).toContainText('Delete Tenant');
    await expect(page.locator('div[role="alertdialog"]')).toContainText('Are you sure you want to delete tenant "test-tenant-1"?');

    // Confirm the deletion
    await page.click('button:has-text("Delete")');

    // Check that the success toast is displayed
    await expect(page.locator('div[role="status"]')).toContainText('Tenant deleted');

    // Check that the tenant is removed from the list
    await expect(page.locator('table tbody tr')).toHaveCount(initialTenantCount - 1);
    await expect(page.locator('table tbody tr:has-text("test-tenant-1")')).not.toBeVisible();
  });

  test('should cancel tenant deletion', async ({ page }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Get the initial count of tenants
    const initialTenantCount = await page.locator('table tbody tr').count();

    // Click the delete button for the first tenant
    await page.locator('table tbody tr').nth(0).locator('text=Delete').click();

    // Check that the confirmation dialog is displayed
    await expect(page.locator('div[role="alertdialog"]')).toBeVisible();

    // Cancel the deletion
    await page.click('button:has-text("Cancel")');

    // Check that the dialog is closed
    await expect(page.locator('div[role="alertdialog"]')).not.toBeVisible();

    // Check that the tenant count remains the same
    await expect(page.locator('table tbody tr')).toHaveCount(initialTenantCount);
  });

  test('should handle API errors during tenant deletion', async ({ page }) => {
    // Set up the mock API to return an error for deletion
    await page.route('/api/tenants/test-tenant-1', async (route) => {
      const method = route.request().method();

      if (method === 'DELETE') {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal server error' }),
        });
      } else {
        await route.continue();
      }
    });

    // Navigate to the homepage
    await page.goto('/');

    // Click the delete button for the first tenant
    await page.locator('table tbody tr').nth(0).locator('text=Delete').click();

    // Confirm the deletion
    await page.click('button:has-text("Delete")');

    // Check that the error toast is displayed
    await expect(page.locator('div[role="alert"]')).toContainText('Failed to delete tenant');

    // Check that the tenant is still in the list
    await expect(page.locator('table tbody tr:has-text("test-tenant-1")')).toBeVisible();
  });
});
