'use client'

import { useEffect, useState } from 'react'
import { Box, Container, Heading, Breadcrumb, BreadcrumbItem, BreadcrumbLink, Spinner, Text, Button } from '@chakra-ui/react'
import { Navbar } from '@/components/Navbar'
import { TenantForm } from '@/components/TenantForm'
import Link from 'next/link'
import { useParams } from 'next/navigation'
import { FiRefreshCw } from 'react-icons/fi'

export default function EditTenant() {
  const params = useParams()
  const tenantName = params.name as string
  
  const [tenant, setTenant] = useState<any>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTenant = async () => {
    setIsLoading(true)
    setError(null)
    
    try {
      const response = await fetch(`/api/tenants/${tenantName}`)
      
      if (!response.ok) {
        throw new Error('Failed to fetch tenant')
      }
      
      const data = await response.json()
      setTenant(data)
    } catch (err) {
      setError('Failed to load tenant data. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchTenant()
  }, [tenantName])

  const initialData = tenant ? {
    name: tenant.metadata.name,
    displayName: tenant.spec.displayName || '',
    description: tenant.spec.description || '',
    serverReplicas: tenant.spec.server?.replicas || 1,
    serverImage: tenant.spec.server?.image || 'neurallog/server:latest',
    redisReplicas: tenant.spec.redis?.replicas || 1,
    redisImage: tenant.spec.redis?.image || 'redis:7-alpine',
    redisStorage: tenant.spec.redis?.storage || '1Gi',
    networkPolicyEnabled: tenant.spec.networkPolicy?.enabled !== false, // Default to true if not specified
  } : {}

  return (
    <Box minH="100vh" bg="gray.50">
      <Navbar />
      <Container maxW="container.xl" py={8}>
        <Breadcrumb mb={6}>
          <BreadcrumbItem>
            <Link href="/" passHref>
              <BreadcrumbLink>Home</BreadcrumbLink>
            </Link>
          </BreadcrumbItem>
          <BreadcrumbItem isCurrentPage>
            <BreadcrumbLink>Edit Tenant</BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>
        
        <Heading as="h1" size="xl" mb={6}>Edit Tenant: {tenantName}</Heading>
        
        {isLoading ? (
          <Box textAlign="center" py={10}>
            <Spinner size="xl" />
            <Text mt={4}>Loading tenant data...</Text>
          </Box>
        ) : error ? (
          <Box textAlign="center" py={10}>
            <Text color="red.500">{error}</Text>
            <Button mt={4} leftIcon={<FiRefreshCw />} onClick={fetchTenant}>
              Retry
            </Button>
          </Box>
        ) : (
          <TenantForm initialData={initialData} isEditing={true} />
        )}
      </Container>
    </Box>
  )
}
