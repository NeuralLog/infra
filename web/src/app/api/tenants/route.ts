import { NextRequest, NextResponse } from 'next/server'
import * as k8s from '@kubernetes/client-node'

// Mock data for testing
const mockTenants = [
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

// Check if we're in test mode
const isTestMode = process.env.TEST_MODE === 'true';

// Initialize Kubernetes client if not in test mode
let k8sApi;
if (!isTestMode) {
  const kc = new k8s.KubeConfig()
  kc.loadFromDefault()
  k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)
}

const GROUP = 'neurallog.io'
const VERSION = 'v1'
const PLURAL = 'tenants'

export async function GET(request: NextRequest) {
  try {
    if (isTestMode) {
      // Return mock data in test mode
      return NextResponse.json({ items: mockTenants })
    }

    // Use Kubernetes API in production mode
    const response = await k8sApi.listClusterCustomObject(
      GROUP,
      VERSION,
      PLURAL
    )

    return NextResponse.json(response.body)
  } catch (error) {
    console.error('Error fetching tenants:', error)
    return NextResponse.json(
      { error: 'Failed to fetch tenants' },
      { status: 500 }
    )
  }
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    if (isTestMode) {
      // Create a new mock tenant
      const name = body.metadata.name;
      const now = new Date().toISOString();

      const newTenant = {
        metadata: {
          name,
          creationTimestamp: now,
        },
        spec: body.spec,
        status: {
          phase: 'Provisioning',
          namespace: `tenant-${name}`,
          serverStatus: {
            phase: 'Provisioning',
            readyReplicas: 0,
            totalReplicas: body.spec.server?.replicas || 1,
            message: 'Server is being provisioned',
          },
          redisStatus: {
            phase: 'Provisioning',
            readyReplicas: 0,
            totalReplicas: body.spec.redis?.replicas || 1,
            message: 'Redis is being provisioned',
          },
        },
      };

      // Add to mock tenants
      mockTenants.push(newTenant);

      return NextResponse.json(newTenant, { status: 201 });
    }

    // Use Kubernetes API in production mode
    const response = await k8sApi.createClusterCustomObject(
      GROUP,
      VERSION,
      PLURAL,
      body
    )

    return NextResponse.json(response.body, { status: 201 })
  } catch (error) {
    console.error('Error creating tenant:', error)
    return NextResponse.json(
      { error: 'Failed to create tenant' },
      { status: 500 }
    )
  }
}
