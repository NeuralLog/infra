import { NextRequest, NextResponse } from 'next/server'
import * as k8s from '@kubernetes/client-node'

// Initialize Kubernetes client
const kc = new k8s.KubeConfig()
kc.loadFromDefault()

const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

const GROUP = 'neurallog.io'
const VERSION = 'v1'
const PLURAL = 'tenants'

export async function GET(
  request: NextRequest,
  { params }: { params: { name: string } }
) {
  try {
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
