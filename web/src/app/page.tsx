'use client'

import { Box, Container, Heading, Text, VStack } from '@chakra-ui/react'
import { Navbar } from '@/components/Navbar'
import { TenantList } from '@/components/TenantList'

export default function Home() {
  return (
    <Box minH="100vh" bg="gray.50">
      <Navbar />
      <Container maxW="container.xl" py={8}>
        <VStack spacing={8} align="stretch">
          <Box>
            <Heading as="h1" size="xl" mb={2}>NeuralLog Admin Dashboard</Heading>
            <Text color="gray.600">Manage your NeuralLog tenants and resources</Text>
          </Box>
          
          <TenantList />
        </VStack>
      </Container>
    </Box>
  )
}
