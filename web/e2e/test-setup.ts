import { test as base } from '@playwright/test';
import { exec } from 'child_process';
import { promisify } from 'util';

const execAsync = promisify(exec);

// Define a custom fixture for test setup and teardown
export const test = base.extend({
  // Setup before tests
  testSetup: [async ({}, use) => {
    // Create test tenants
    await createTestTenant('test-tenant-1', 'Test Tenant 1', 'A test tenant for e2e testing');
    await createTestTenant('test-tenant-2', 'Test Tenant 2', 'Another test tenant for e2e testing', 2);
    
    // Use the fixture
    await use();
    
    // Clean up after tests
    await deleteTestTenant('test-tenant-1');
    await deleteTestTenant('test-tenant-2');
  }, { scope: 'worker' }],
});

export { expect } from '@playwright/test';

// Helper function to create a test tenant
async function createTestTenant(name: string, displayName: string, description: string, replicas: number = 1) {
  const tenantYaml = `
apiVersion: neurallog.io/v1
kind: Tenant
metadata:
  name: ${name}
spec:
  displayName: ${displayName}
  description: ${description}
  server:
    replicas: ${replicas}
    image: neurallog/server:latest
  redis:
    replicas: 1
    image: redis:7-alpine
    storage: 1Gi
  networkPolicy:
    enabled: true
`;

  try {
    await execAsync(`kubectl apply -f - <<EOF\n${tenantYaml}\nEOF`);
    console.log(`Created test tenant: ${name}`);
  } catch (error) {
    console.error(`Failed to create test tenant ${name}:`, error);
    throw error;
  }
}

// Helper function to delete a test tenant
async function deleteTestTenant(name: string) {
  try {
    await execAsync(`kubectl delete tenant ${name} --ignore-not-found`);
    console.log(`Deleted test tenant: ${name}`);
  } catch (error) {
    console.error(`Failed to delete test tenant ${name}:`, error);
    // Don't throw error on cleanup
  }
}
