import { test, expect } from './test-setup';

test.describe('Tenant Creation', () => {
  test('should create a new tenant with default values', async ({ page }) => {
    // Navigate to the new tenant page
    await page.goto('/tenants/new');

    // Check that the page title is correct
    await expect(page.locator('h1').filter({ hasText: 'Create New Tenant' })).toBeVisible();

    // Fill in the required fields
    await page.fill('input[placeholder="my-tenant"]', 'e2e-test-tenant');
    await page.fill('input[placeholder="My Tenant"]', 'E2E Test Tenant');
    await page.fill('textarea[placeholder="Description of the tenant"]', 'A tenant created during E2E testing');

    // Submit the form
    await page.click('button:has-text("Create Tenant")');

    // Check that we're redirected to the homepage
    await expect(page).toHaveURL('/');

    // Check that the success toast is displayed
    await expect(page.locator('div[role="status"]')).toContainText('Tenant created');

    // Check that the new tenant is in the list
    const tenantRows = page.locator('table tbody tr');
    const newTenantRow = tenantRows.filter({ hasText: 'e2e-test-tenant' });
    await expect(newTenantRow).toBeVisible();
    await expect(newTenantRow.locator('td').nth(1)).toContainText('E2E Test Tenant');
  });

  test('should create a tenant with custom configuration', async ({ page }) => {
    // Navigate to the new tenant page
    await page.goto('/tenants/new');

    // Fill in the basic fields
    await page.fill('input[placeholder="my-tenant"]', 'custom-tenant');
    await page.fill('input[placeholder="My Tenant"]', 'Custom Config Tenant');
    await page.fill('textarea[placeholder="Description of the tenant"]', 'A tenant with custom configuration');

    // Customize server configuration
    await page.fill('input[name="serverReplicas"]', '2');
    await page.fill('input[placeholder="neurallog/server:latest"]', 'neurallog/server:v1.2.3');

    // Customize Redis configuration
    await page.fill('input[name="redisReplicas"]', '2');
    await page.fill('input[placeholder="redis:7-alpine"]', 'redis:6-alpine');
    await page.fill('input[placeholder="1Gi"]', '5Gi');

    // Disable network policies
    await page.click('input[id="networkPolicyEnabled"]');

    // Submit the form
    await page.click('button:has-text("Create Tenant")');

    // Check that we're redirected to the homepage
    await expect(page).toHaveURL('/');

    // Check that the success toast is displayed
    await expect(page.locator('div[role="status"]')).toContainText('Tenant created');

    // Check that the new tenant is in the list
    const tenantRows = page.locator('table tbody tr');
    const newTenantRow = tenantRows.filter({ hasText: 'custom-tenant' });
    await expect(newTenantRow).toBeVisible();
    await expect(newTenantRow.locator('td').nth(1)).toContainText('Custom Config Tenant');
  });

  test('should validate required fields', async ({ page }) => {
    // Navigate to the new tenant page
    await page.goto('/tenants/new');

    // Submit the form without filling any fields
    await page.click('button:has-text("Create Tenant")');

    // Check that validation errors are displayed
    await expect(page.locator('text=Tenant name is required')).toBeVisible();

    // Fill in the tenant name with an invalid value
    await page.fill('input[placeholder="my-tenant"]', 'Invalid Name!');

    // Check that the validation error changes
    await expect(page.locator('text=Tenant name must consist of lowercase alphanumeric characters')).toBeVisible();

    // Fix the tenant name
    await page.fill('input[placeholder="my-tenant"]', 'valid-name');

    // Check that the validation error disappears
    await expect(page.locator('text=Tenant name is required')).not.toBeVisible();
    await expect(page.locator('text=Tenant name must consist of lowercase alphanumeric characters')).not.toBeVisible();
  });

  test('should handle API errors during tenant creation', async ({ page, mockApi }) => {
    // Set up the mock API to return an error
    await page.route('/api/tenants', async (route) => {
      const method = route.request().method();

      if (method === 'POST') {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal server error' }),
        });
      } else {
        await route.continue();
      }
    });

    // Navigate to the new tenant page
    await page.goto('/tenants/new');

    // Fill in the required fields
    await page.fill('input[placeholder="my-tenant"]', 'error-tenant');
    await page.fill('input[placeholder="My Tenant"]', 'Error Tenant');

    // Submit the form
    await page.click('button:has-text("Create Tenant")');

    // Check that we're still on the new tenant page
    await expect(page).toHaveURL('/tenants/new');

    // Check that the error toast is displayed
    await expect(page.locator('div[role="alert"]')).toContainText('Failed to create tenant');
  });
});
