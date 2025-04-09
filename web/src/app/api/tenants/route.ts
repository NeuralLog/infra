import { NextRequest, NextResponse } from 'next/server'
import * as k8s from '@kubernetes/client-node'

// Initialize Kubernetes client
const kc = new k8s.KubeConfig()
kc.loadFromDefault()

const k8sApi = kc.makeApiClient(k8s.CustomObjectsApi)

const GROUP = 'neurallog.io'
const VERSION = 'v1'
const PLURAL = 'tenants'

export async function GET(request: NextRequest) {
  try {
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
