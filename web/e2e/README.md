# NeuralLog Admin E2E Tests

This directory contains end-to-end tests for the NeuralLog Admin web interface using Playwright.

## Test Structure

The tests are organized into the following files:

- `fixtures.ts` - Test fixtures and utilities
- `mock-api-handler.ts` - Mock API handler for intercepting API calls
- `homepage.spec.ts` - Tests for the homepage and tenant listing
- `tenant-creation.spec.ts` - Tests for tenant creation
- `tenant-editing.spec.ts` - Tests for tenant editing
- `tenant-deletion.spec.ts` - Tests for tenant deletion
- `responsive-and-a11y.spec.ts` - Tests for responsive design and accessibility

## Running Tests

You can run the tests using the following npm scripts:

```bash
# Run all tests
npm test

# Run tests with UI mode
npm run test:ui

# Run tests in headed mode (with browser visible)
npm run test:headed

# Run tests in debug mode
npm run test:debug
```

## Test Coverage

The tests cover the following aspects of the application:

### Homepage and Tenant Listing
- Display of tenant list
- Navigation to tenant creation and editing
- Empty state when no tenants exist
- Refreshing the tenant list

### Tenant Creation
- Creating a tenant with default values
- Creating a tenant with custom configuration
- Validation of required fields
- Handling API errors during tenant creation

### Tenant Editing
- Loading tenant data for editing
- Updating tenant configuration
- Handling API errors during tenant loading and update

### Tenant Deletion
- Deleting a tenant
- Canceling tenant deletion
- Handling API errors during tenant deletion

### Responsive Design and Accessibility
- Responsive layout on mobile devices
- Responsive layout on tablet devices
- Keyboard navigation and focus management
- Dialog accessibility

## Mock API

The tests use a mock API handler to intercept API calls and return mock responses. This allows the tests to run without a real backend and to simulate various scenarios, including error conditions.

The mock API handler is implemented in `mock-api-handler.ts` and provides the following functionality:

- Intercepting API calls to `/api/tenants` and `/api/tenants/:name`
- Returning mock tenant data
- Simulating tenant creation, update, and deletion
- Simulating API errors

## Adding New Tests

To add new tests:

1. Create a new test file in the `e2e` directory
2. Import the test fixtures from `fixtures.ts`
3. Use the mock API handler to set up test data
4. Write your tests using Playwright's API

Example:

```typescript
import { test, expect } from './fixtures';

test.describe('My New Feature', () => {
  test('should do something', async ({ page, mockApi }) => {
    // Set up test data
    await mockApi.addTenant({
      metadata: { name: 'test-tenant' },
      spec: { displayName: 'Test Tenant' }
    });

    // Navigate to the page
    await page.goto('/my-feature');

    // Perform actions and assertions
    await page.click('button:has-text("Do Something")');
    await expect(page.locator('text=Success')).toBeVisible();
  });
});
```
