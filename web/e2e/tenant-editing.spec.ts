import { test, expect } from './test-setup';

test.describe('Tenant Editing', () => {
  test('should load tenant data for editing', async ({ page }) => {
    // Navigate to the edit page for the first tenant
    await page.goto('/tenants/test-tenant-1/edit');

    // Check that the page title is correct
    await expect(page.locator('h1')).toContainText('Edit Tenant: test-tenant-1');

    // Check that the form is pre-filled with tenant data
    await expect(page.locator('input[placeholder="my-tenant"]')).toHaveValue('test-tenant-1');
    await expect(page.locator('input[placeholder="my-tenant"]')).toBeDisabled();
    await expect(page.locator('input[placeholder="My Tenant"]')).toHaveValue('Test Tenant 1');
    await expect(page.locator('textarea[placeholder="Description of the tenant"]')).toHaveValue('A test tenant for e2e testing');

    // Check server configuration
    await expect(page.locator('input[name="serverReplicas"]')).toHaveValue('1');
    await expect(page.locator('input[placeholder="neurallog/server:latest"]')).toHaveValue('neurallog/server:latest');

    // Check Redis configuration
    await expect(page.locator('input[name="redisReplicas"]')).toHaveValue('1');
    await expect(page.locator('input[placeholder="redis:7-alpine"]')).toHaveValue('redis:7-alpine');
    await expect(page.locator('input[placeholder="1Gi"]')).toHaveValue('1Gi');

    // Check network policy configuration
    await expect(page.locator('input[id="networkPolicyEnabled"]')).toBeChecked();
  });

  test('should update tenant configuration', async ({ page }) => {
    // Navigate to the edit page for the first tenant
    await page.goto('/tenants/test-tenant-1/edit');

    // Update the display name and description
    await page.fill('input[placeholder="My Tenant"]', 'Updated Test Tenant');
    await page.fill('textarea[placeholder="Description of the tenant"]', 'This tenant has been updated');

    // Update server configuration
    await page.fill('input[name="serverReplicas"]', '3');

    // Update Redis configuration
    await page.fill('input[placeholder="1Gi"]', '10Gi');

    // Submit the form
    await page.click('button:has-text("Update Tenant")');

    // Check that we're redirected to the homepage
    await expect(page).toHaveURL('/');

    // Check that the success toast is displayed
    await expect(page.locator('div[role="status"]')).toContainText('Tenant updated');

    // Check that the tenant is updated in the list
    const tenantRows = page.locator('table tbody tr');
    const updatedTenantRow = tenantRows.filter({ hasText: 'test-tenant-1' });
    await expect(updatedTenantRow.locator('td').nth(1)).toContainText('Updated Test Tenant');

    // Navigate back to the edit page to verify the changes were saved
    await page.goto('/tenants/test-tenant-1/edit');

    // Check that the form reflects the updated values
    await expect(page.locator('input[placeholder="My Tenant"]')).toHaveValue('Updated Test Tenant');
    await expect(page.locator('textarea[placeholder="Description of the tenant"]')).toHaveValue('This tenant has been updated');
    await expect(page.locator('input[name="serverReplicas"]')).toHaveValue('3');
    await expect(page.locator('input[placeholder="1Gi"]')).toHaveValue('10Gi');
  });

  test('should handle API errors during tenant loading', async ({ page }) => {
    // Set up the mock API to return an error for a non-existent tenant
    await page.route('/api/tenants/non-existent-tenant', async (route) => {
      await route.fulfill({
        status: 404,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Tenant not found' }),
      });
    });

    // Navigate to the edit page for a non-existent tenant
    await page.goto('/tenants/non-existent-tenant/edit');

    // Check that the error message is displayed
    await expect(page.locator('text=Failed to load tenant data')).toBeVisible();
    await expect(page.locator('button:has-text("Retry")')).toBeVisible();
  });

  test('should handle API errors during tenant update', async ({ page }) => {
    // Set up the mock API to return an error for the update
    await page.route('/api/tenants/test-tenant-1', async (route) => {
      const method = route.request().method();

      if (method === 'PUT') {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal server error' }),
        });
      } else {
        await route.continue();
      }
    });

    // Navigate to the edit page for the first tenant
    await page.goto('/tenants/test-tenant-1/edit');

    // Update the display name
    await page.fill('input[placeholder="My Tenant"]', 'Error Test');

    // Submit the form
    await page.click('button:has-text("Update Tenant")');

    // Check that we're still on the edit page
    await expect(page).toHaveURL('/tenants/test-tenant-1/edit');

    // Check that the error toast is displayed
    await expect(page.locator('div[role="alert"]')).toContainText('Failed to update tenant');
  });
});
