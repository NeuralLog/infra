import { test as base } from '@playwright/test';
import { MockApiHandler } from './mock-api-handler';

// Extend the base test fixture with our custom fixtures
export const test = base.extend<{
  mockApi: MockApiHandler;
}>({
  // Define the mockApi fixture
  mockApi: async ({ page }, use) => {
    const mockApiHandler = new MockApiHandler(page);
    await mockApiHandler.initialize();
    await use(mockApiHandler);
    await mockApiHandler.cleanup();
  },
});

export { expect } from '@playwright/test';
