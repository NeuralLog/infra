import { test, expect } from './test-setup';
import { devices } from '@playwright/test';

test.describe('Responsive Design and Accessibility', () => {
  test('should be responsive on mobile devices', async ({ page }) => {
    // Set the viewport to mobile size
    await page.setViewportSize(devices['iPhone 12'].viewport);

    // Navigate to the homepage
    await page.goto('/');

    // Check that the navbar is properly displayed
    await expect(page.locator('h1').filter({ hasText: 'NeuralLog Admin' })).toBeVisible();
    await expect(page.locator('button:has-text("New Tenant")')).toBeVisible();

    // Check that the tenant list is properly displayed
    await expect(page.locator('h2:has-text("Tenants")')).toBeVisible();
    await expect(page.locator('table')).toBeVisible();

    // Navigate to the new tenant page
    await page.click('button:has-text("New Tenant")');

    // Check that the form is properly displayed
    await expect(page.locator('h1:has-text("Create New Tenant")')).toBeVisible();
    await expect(page.locator('input[placeholder="my-tenant"]')).toBeVisible();

    // Check that the form controls are properly sized for mobile
    const inputWidth = await page.locator('input[placeholder="my-tenant"]').evaluate(el => {
      const rect = el.getBoundingClientRect();
      return rect.width;
    });

    // The input should take up most of the viewport width
    expect(inputWidth).toBeGreaterThan(devices['iPhone 12'].viewport.width * 0.7);
  });

  test('should be responsive on tablet devices', async ({ page }) => {
    // Set the viewport to tablet size
    await page.setViewportSize(devices['iPad (gen 7)'].viewport);

    // Navigate to the homepage
    await page.goto('/');

    // Check that the layout adjusts for tablet
    await expect(page.locator('h1').filter({ hasText: 'NeuralLog Admin' })).toBeVisible();
    await expect(page.locator('button:has-text("New Tenant")')).toBeVisible();

    // Navigate to the new tenant page
    await page.click('button:has-text("New Tenant")');

    // Check that the form layout adjusts for tablet
    // On tablets, we should have a two-column layout for some form sections
    const formGrid = page.locator('div').filter({ hasText: /Tenant Name/ }).first();
    const gridStyle = await formGrid.evaluate(el => {
      return window.getComputedStyle(el).getPropertyValue('grid-template-columns');
    });

    // Should have multiple columns
    expect(gridStyle).not.toBe('1fr');
  });

  test('should have proper focus management', async ({ page }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Press Tab to focus on the first focusable element
    await page.keyboard.press('Tab');

    // Check that the "New Tenant" button is focused
    await expect(page.locator('button:has-text("New Tenant")')).toBeFocused();

    // Press Tab again to focus on the next element
    await page.keyboard.press('Tab');

    // Check that the "Refresh" button is focused
    await expect(page.locator('button:has-text("Refresh")')).toBeFocused();

    // Continue tabbing to reach the first tenant's Edit button
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');

    // Check that the Edit button is focused
    const editButton = page.locator('table tbody tr').nth(0).locator('button:has-text("Edit")');
    await expect(editButton).toBeFocused();

    // Press Enter to navigate to the edit page
    await page.keyboard.press('Enter');

    // Check that we're on the edit page
    await expect(page).toHaveURL(/\/tenants\/test-tenant-1\/edit/);
  });

  test('should handle keyboard navigation in dialogs', async ({ page }) => {
    // Navigate to the homepage
    await page.goto('/');

    // Click the delete button for the first tenant
    await page.locator('table tbody tr').nth(0).locator('text=Delete').click();

    // Check that the dialog is displayed
    await expect(page.locator('div[role="alertdialog"]')).toBeVisible();

    // Check that the Cancel button is focused (as it should be the least destructive action)
    await expect(page.locator('button:has-text("Cancel")')).toBeFocused();

    // Press Tab to focus on the Delete button
    await page.keyboard.press('Tab');

    // Check that the Delete button is focused
    await expect(page.locator('button:has-text("Delete")')).toBeFocused();

    // Press Escape to close the dialog
    await page.keyboard.press('Escape');

    // Check that the dialog is closed
    await expect(page.locator('div[role="alertdialog"]')).not.toBeVisible();
  });
});
