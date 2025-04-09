# NeuralLog Admin E2E Tests

This directory contains end-to-end tests for the NeuralLog Admin web interface using Playwright. These tests interact with a real Kubernetes API to test the full functionality of the application.

## Test Structure

The tests are organized into the following files:

- `test-setup.ts` - Test fixtures and utilities for setting up and tearing down test data
- `setup-kind.sh` - Script to set up a local Kubernetes environment using Kind
- `homepage.spec.ts` - Tests for the homepage and tenant listing
- `tenant-creation.spec.ts` - Tests for tenant creation
- `tenant-editing.spec.ts` - Tests for tenant editing
- `tenant-deletion.spec.ts` - Tests for tenant deletion
- `responsive-and-a11y.spec.ts` - Tests for responsive design and accessibility

## Running Tests

### Prerequisites

Before running the tests, you need to set up a local Kubernetes environment:

1. Install Kind (Kubernetes in Docker)
2. Install kubectl
3. Run the setup script:

```bash
# Make the script executable
chmod +x setup-kind.sh

# Run the setup script
./setup-kind.sh
```

### Running the Tests

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

### Environment Variables

You can configure the tests using the following environment variables:

- `API_URL`: The URL of the Kubernetes API (default: http://localhost:30000)

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

## Real API Integration

These tests interact with a real Kubernetes API to provide true end-to-end testing. The `test-setup.ts` file contains utilities for setting up and tearing down test data in the Kubernetes cluster.

The tests perform the following operations on the real API:

- Creating test tenants before tests run
- Reading tenant data during tests
- Updating tenant configurations
- Deleting test tenants after tests complete

This approach ensures that the tests verify the full functionality of the application, including its integration with the Kubernetes API.

## Adding New Tests

To add new tests:

1. Create a new test file in the `e2e` directory
2. Import the test fixtures from `test-setup.ts`
3. Use the testSetup fixture to ensure test data is available
4. Write your tests using Playwright's API

Example:

```typescript
import { test, expect } from './test-setup';

test.describe('My New Feature', () => {
  test('should do something', async ({ page, testSetup }) => {
    // The testSetup fixture ensures test tenants are created

    // Navigate to the page
    await page.goto('/my-feature');

    // Perform actions and assertions
    await page.click('button:has-text("Do Something")');
    await expect(page.locator('text=Success')).toBeVisible();
  });
});
```
