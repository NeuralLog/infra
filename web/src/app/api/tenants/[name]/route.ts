import { NextRequest, NextResponse } from 'next/server'
import * as k8s from '@kubernetes/client-node'

// Check if we're in test mode
const isTestMode = process.env.TEST_MODE === 'true';

// Get mock tenants from the parent route
let mockTenants: any[] = [];
if (isTestMode) {
  // Import mock tenants dynamically in test mode
  import('../route').then(module => {
    mockTenants = (module as any).mockTenants;
  }).catch(error => {
    console.error('Error importing mock tenants:', error);
  });
}

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

export async function GET(
  request: NextRequest,
  { params }: { params: { name: string } }
) {
  try {
    if (isTestMode) {
      // Find tenant in mock data
      const tenant = mockTenants.find(t => t.metadata.name === params.name);

      if (!tenant) {
        return NextResponse.json(
          { error: `Tenant ${params.name} not found` },
          { status: 404 }
        );
      }

      return NextResponse.json(tenant);
    }

    // Use Kubernetes API in production mode
    const response = await k8sApi.getClusterCustomObject(
      GROUP,
      VERSION,
      PLURAL,
      params.name
    )

    return NextResponse.json(response.body)
  } catch (error) {
    console.error(`Error fetching tenant ${params.name}:`, error)
    return NextResponse.json(
      { error: `Failed to fetch tenant ${params.name}` },
      { status: 500 }
    )
  }
}

export async function PUT(
  request: NextRequest,
  { params }: { params: { name: string } }
) {
  try {
    const body = await request.json()

    // Ensure the name in the path matches the name in the body
    if (body.metadata.name !== params.name) {
      return NextResponse.json(
        { error: 'Tenant name in path does not match name in body' },
        { status: 400 }
      )
    }

    if (isTestMode) {
      // Find tenant in mock data
      const tenantIndex = mockTenants.findIndex(t => t.metadata.name === params.name);

      if (tenantIndex === -1) {
        return NextResponse.json(
          { error: `Tenant ${params.name} not found` },
          { status: 404 }
        );
      }

      // Update tenant
      mockTenants[tenantIndex] = {
        ...mockTenants[tenantIndex],
        spec: body.spec,
      };

      return NextResponse.json(mockTenants[tenantIndex]);
    }

    // Use Kubernetes API in production mode
    const response = await k8sApi.replaceClusterCustomObject(
      GROUP,
      VERSION,
      PLURAL,
      params.name,
      body
    )

    return NextResponse.json(response.body)
  } catch (error) {
    console.error(`Error updating tenant ${params.name}:`, error)
    return NextResponse.json(
      { error: `Failed to update tenant ${params.name}` },
      { status: 500 }
    )
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: { name: string } }
) {
  try {
    if (isTestMode) {
      // Find tenant in mock data
      const tenantIndex = mockTenants.findIndex(t => t.metadata.name === params.name);

      if (tenantIndex === -1) {
        return NextResponse.json(
          { error: `Tenant ${params.name} not found` },
          { status: 404 }
        );
      }

      // Remove tenant
      mockTenants.splice(tenantIndex, 1);

      return NextResponse.json({ message: `Tenant ${params.name} deleted successfully` });
    }

    // Use Kubernetes API in production mode
    const response = await k8sApi.deleteClusterCustomObject(
      GROUP,
      VERSION,
      PLURAL,
      params.name
    )

    return NextResponse.json(response.body)
  } catch (error) {
    console.error(`Error deleting tenant ${params.name}:`, error)
    return NextResponse.json(
      { error: `Failed to delete tenant ${params.name}` },
      { status: 500 }
    )
  }
}
