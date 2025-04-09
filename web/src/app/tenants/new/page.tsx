'use client'

import { Box, Container, Heading, Breadcrumb, BreadcrumbItem, BreadcrumbLink } from '@chakra-ui/react'
import { Navbar } from '@/components/Navbar'
import { TenantForm } from '@/components/TenantForm'
import Link from 'next/link'

export default function NewTenant() {
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
            <BreadcrumbLink>New Tenant</BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>
        
        <Heading as="h1" size="xl" mb={6}>Create New Tenant</Heading>
        
        <TenantForm />
      </Container>
    </Box>
  )
}
