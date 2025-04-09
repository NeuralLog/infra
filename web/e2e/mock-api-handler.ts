import { Page } from '@playwright/test';

export interface Tenant {
  metadata: {
    name: string;
    creationTimestamp: string;
  };
  spec: {
    displayName: string;
    description: string;
    server?: {
      replicas: number;
      image: string;
      resources?: {
        cpu?: {
          request?: string;
          limit?: string;
        };
        memory?: {
          request?: string;
          limit?: string;
        };
      };
      env?: Array<{
        name: string;
        value: string;
      }>;
    };
    redis?: {
      replicas: number;
      image: string;
      storage: string;
      resources?: {
        cpu?: {
          request?: string;
          limit?: string;
        };
        memory?: {
          request?: string;
          limit?: string;
        };
      };
    };
    networkPolicy?: {
      enabled: boolean;
    };
  };
  status: {
    phase: string;
    namespace: string;
    serverStatus?: {
      phase: string;
      readyReplicas: number;
      totalReplicas: number;
      message?: string;
    };
    redisStatus?: {
      phase: string;
      readyReplicas: number;
      totalReplicas: number;
      message?: string;
    };
  };
}

export class MockApiHandler {
  private page: Page;
  private tenants: Tenant[] = [];

  constructor(page: Page) {
    this.page = page;
    this.initializeMockTenants();
  }

  private initializeMockTenants() {
    // Create some mock tenants for testing
    this.tenants = [
      {
        metadata: {
          name: 'test-tenant-1',
          creationTimestamp: new Date(Date.now() - 86400000).toISOString(), // 1 day ago
        },
        spec: {
          displayName: 'Test Tenant 1',
          description: 'A test tenant for e2e testing',
          server: {
            replicas: 1,
            image: 'neurallog/server:latest',
          },
          redis: {
            replicas: 1,
            image: 'redis:7-alpine',
            storage: '1Gi',
          },
          networkPolicy: {
            enabled: true,
          },
        },
        status: {
          phase: 'Running',
          namespace: 'tenant-test-tenant-1',
          serverStatus: {
            phase: 'Running',
            readyReplicas: 1,
            totalReplicas: 1,
            message: 'Server is running',
          },
          redisStatus: {
            phase: 'Running',
            readyReplicas: 1,
            totalReplicas: 1,
            message: 'Redis is running',
          },
        },
      },
      {
        metadata: {
          name: 'test-tenant-2',
          creationTimestamp: new Date(Date.now() - 43200000).toISOString(), // 12 hours ago
        },
        spec: {
          displayName: 'Test Tenant 2',
          description: 'Another test tenant for e2e testing',
          server: {
            replicas: 2,
            image: 'neurallog/server:latest',
          },
          redis: {
            replicas: 1,
            image: 'redis:7-alpine',
            storage: '2Gi',
          },
          networkPolicy: {
            enabled: true,
          },
        },
        status: {
          phase: 'Degraded',
          namespace: 'tenant-test-tenant-2',
          serverStatus: {
            phase: 'Degraded',
            readyReplicas: 1,
            totalReplicas: 2,
            message: 'Server is degraded: 1/2 replicas ready',
          },
          redisStatus: {
            phase: 'Running',
            readyReplicas: 1,
            totalReplicas: 1,
            message: 'Redis is running',
          },
        },
      },
    ];
  }

  public async initialize() {
    // Intercept API calls and return mock responses
    await this.page.route('/api/tenants', async (route) => {
      const method = route.request().method();
      
      if (method === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ items: this.tenants }),
        });
      } else if (method === 'POST') {
        const requestBody = JSON.parse(await route.request().postData() || '{}');
        const newTenant = this.createTenantFromRequest(requestBody);
        this.tenants.push(newTenant);
        
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify(newTenant),
        });
      } else {
        await route.continue();
      }
    });

    // Handle individual tenant routes
    await this.page.route('/api/tenants/:name', async (route) => {
      const method = route.request().method();
      const url = route.request().url();
      const tenantName = url.split('/').pop() || '';
      const tenantIndex = this.tenants.findIndex(t => t.metadata.name === tenantName);
      
      if (method === 'GET') {
        if (tenantIndex === -1) {
          await route.fulfill({
            status: 404,
            contentType: 'application/json',
            body: JSON.stringify({ error: `Tenant ${tenantName} not found` }),
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify(this.tenants[tenantIndex]),
          });
        }
      } else if (method === 'PUT') {
        if (tenantIndex === -1) {
          await route.fulfill({
            status: 404,
            contentType: 'application/json',
            body: JSON.stringify({ error: `Tenant ${tenantName} not found` }),
          });
        } else {
          const requestBody = JSON.parse(await route.request().postData() || '{}');
          this.tenants[tenantIndex] = {
            ...this.tenants[tenantIndex],
            spec: requestBody.spec,
          };
          
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify(this.tenants[tenantIndex]),
          });
        }
      } else if (method === 'DELETE') {
        if (tenantIndex === -1) {
          await route.fulfill({
            status: 404,
            contentType: 'application/json',
            body: JSON.stringify({ error: `Tenant ${tenantName} not found` }),
          });
        } else {
          this.tenants.splice(tenantIndex, 1);
          
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({ message: `Tenant ${tenantName} deleted successfully` }),
          });
        }
      } else {
        await route.continue();
      }
    });
  }

  private createTenantFromRequest(requestBody: any): Tenant {
    const name = requestBody.metadata.name;
    const now = new Date().toISOString();
    
    return {
      metadata: {
        name,
        creationTimestamp: now,
      },
      spec: requestBody.spec,
      status: {
        phase: 'Provisioning',
        namespace: `tenant-${name}`,
        serverStatus: {
          phase: 'Provisioning',
          readyReplicas: 0,
          totalReplicas: requestBody.spec.server?.replicas || 1,
          message: 'Server is being provisioned',
        },
        redisStatus: {
          phase: 'Provisioning',
          readyReplicas: 0,
          totalReplicas: requestBody.spec.redis?.replicas || 1,
          message: 'Redis is being provisioned',
        },
      },
    };
  }

  public async addTenant(tenant: Partial<Tenant>) {
    const name = tenant.metadata?.name || `test-tenant-${this.tenants.length + 1}`;
    const now = new Date().toISOString();
    
    const newTenant: Tenant = {
      metadata: {
        name,
        creationTimestamp: now,
      },
      spec: {
        displayName: tenant.spec?.displayName || `Test Tenant ${this.tenants.length + 1}`,
        description: tenant.spec?.description || 'A test tenant',
        server: tenant.spec?.server || {
          replicas: 1,
          image: 'neurallog/server:latest',
        },
        redis: tenant.spec?.redis || {
          replicas: 1,
          image: 'redis:7-alpine',
          storage: '1Gi',
        },
        networkPolicy: tenant.spec?.networkPolicy || {
          enabled: true,
        },
      },
      status: {
        phase: 'Running',
        namespace: `tenant-${name}`,
        serverStatus: {
          phase: 'Running',
          readyReplicas: 1,
          totalReplicas: 1,
          message: 'Server is running',
        },
        redisStatus: {
          phase: 'Running',
          readyReplicas: 1,
          totalReplicas: 1,
          message: 'Redis is running',
        },
      },
    };
    
    this.tenants.push(newTenant);
    return newTenant;
  }

  public async updateTenant(name: string, updates: Partial<Tenant['spec']>) {
    const tenantIndex = this.tenants.findIndex(t => t.metadata.name === name);
    if (tenantIndex === -1) {
      throw new Error(`Tenant ${name} not found`);
    }
    
    this.tenants[tenantIndex].spec = {
      ...this.tenants[tenantIndex].spec,
      ...updates,
    };
    
    return this.tenants[tenantIndex];
  }

  public async deleteTenant(name: string) {
    const tenantIndex = this.tenants.findIndex(t => t.metadata.name === name);
    if (tenantIndex === -1) {
      throw new Error(`Tenant ${name} not found`);
    }
    
    this.tenants.splice(tenantIndex, 1);
  }

  public async cleanup() {
    // Reset the mock tenants
    this.initializeMockTenants();
  }
}
