'use client'

import React, { useState } from 'react'
import {
  Box,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Heading,
  Button,
  HStack,
  Badge,
  useToast,
  Spinner,
  Text,
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogContent,
  AlertDialogOverlay,
  useDisclosure,
} from '@chakra-ui/react'
import { FiEdit2, FiTrash2, FiRefreshCw } from 'react-icons/fi'
import Link from 'next/link'
import useSWR from 'swr'
import { fetcher } from '@/lib/api'

type Tenant = {
  metadata: {
    name: string
    creationTimestamp: string
  }
  spec: {
    displayName: string
    description: string
  }
  status: {
    phase: string
    namespace: string
    serverStatus?: {
      phase: string
      readyReplicas: number
      totalReplicas: number
    }
    redisStatus?: {
      phase: string
      readyReplicas: number
      totalReplicas: number
    }
  }
}

export function TenantList() {
  const { data, error, isLoading, mutate } = useSWR<{ items: Tenant[] }>('/api/tenants', fetcher)
  const toast = useToast()
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [tenantToDelete, setTenantToDelete] = useState<string | null>(null)
  const cancelRef = React.useRef<HTMLButtonElement>(null)

  const handleDelete = async () => {
    if (!tenantToDelete) return

    try {
      await fetch(`/api/tenants/${tenantToDelete}`, {
        method: 'DELETE',
      })

      toast({
        title: 'Tenant deleted',
        description: `Tenant ${tenantToDelete} has been deleted successfully.`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      })

      mutate()
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to delete tenant. Please try again.',
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    } finally {
      setTenantToDelete(null)
      onClose()
    }
  }

  const confirmDelete = (name: string) => {
    setTenantToDelete(name)
    onOpen()
  }

  if (isLoading) {
    return (
      <Box textAlign="center" py={10}>
        <Spinner size="xl" />
        <Text mt={4}>Loading tenants...</Text>
      </Box>
    )
  }

  if (error) {
    return (
      <Box textAlign="center" py={10}>
        <Text color="red.500">Error loading tenants. Please try again.</Text>
        <Button mt={4} leftIcon={<FiRefreshCw />} onClick={() => mutate()}>
          Retry
        </Button>
      </Box>
    )
  }

  const tenants = data?.items || []

  return (
    <Box bg="white" p={5} shadow="md" borderRadius="md">
      <HStack justify="space-between" mb={4}>
        <Heading as="h2" size="md">Tenants</Heading>
        <Button size="sm" leftIcon={<FiRefreshCw />} onClick={() => mutate()}>
          Refresh
        </Button>
      </HStack>

      {tenants.length === 0 ? (
        <Text py={4}>No tenants found. Create your first tenant to get started.</Text>
      ) : (
        <Table variant="simple">
          <Thead>
            <Tr>
              <Th>Name</Th>
              <Th>Display Name</Th>
              <Th>Namespace</Th>
              <Th>Status</Th>
              <Th>Server</Th>
              <Th>Redis</Th>
              <Th>Actions</Th>
            </Tr>
          </Thead>
          <Tbody>
            {tenants.map((tenant) => (
              <Tr key={tenant.metadata.name}>
                <Td>{tenant.metadata.name}</Td>
                <Td>{tenant.spec.displayName || '-'}</Td>
                <Td>{tenant.status.namespace || '-'}</Td>
                <Td>
                  <Badge
                    colorScheme={
                      tenant.status.phase === 'Running'
                        ? 'green'
                        : tenant.status.phase === 'Provisioning'
                        ? 'blue'
                        : tenant.status.phase === 'Failed'
                        ? 'red'
                        : 'gray'
                    }
                  >
                    {tenant.status.phase}
                  </Badge>
                </Td>
                <Td>
                  {tenant.status.serverStatus ? (
                    <Badge
                      colorScheme={
                        tenant.status.serverStatus.phase === 'Running'
                          ? 'green'
                          : tenant.status.serverStatus.phase === 'Provisioning'
                          ? 'blue'
                          : tenant.status.serverStatus.phase === 'Degraded'
                          ? 'yellow'
                          : tenant.status.serverStatus.phase === 'Failed'
                          ? 'red'
                          : 'gray'
                      }
                    >
                      {tenant.status.serverStatus.readyReplicas}/{tenant.status.serverStatus.totalReplicas}
                    </Badge>
                  ) : (
                    '-'
                  )}
                </Td>
                <Td>
                  {tenant.status.redisStatus ? (
                    <Badge
                      colorScheme={
                        tenant.status.redisStatus.phase === 'Running'
                          ? 'green'
                          : tenant.status.redisStatus.phase === 'Provisioning'
                          ? 'blue'
                          : tenant.status.redisStatus.phase === 'Degraded'
                          ? 'yellow'
                          : tenant.status.redisStatus.phase === 'Failed'
                          ? 'red'
                          : 'gray'
                      }
                    >
                      {tenant.status.redisStatus.readyReplicas}/{tenant.status.redisStatus.totalReplicas}
                    </Badge>
                  ) : (
                    '-'
                  )}
                </Td>
                <Td>
                  <HStack spacing={2}>
                    <Link href={`/tenants/${tenant.metadata.name}/edit`} passHref>
                      <Button
                        size="sm"
                        colorScheme="blue"
                        variant="ghost"
                        leftIcon={<FiEdit2 />}
                      >
                        Edit
                      </Button>
                    </Link>
                    <Button
                      size="sm"
                      colorScheme="red"
                      variant="ghost"
                      leftIcon={<FiTrash2 />}
                      onClick={() => confirmDelete(tenant.metadata.name)}
                    >
                      Delete
                    </Button>
                  </HStack>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      )}

      <AlertDialog isOpen={isOpen} onClose={onClose} leastDestructiveRef={cancelRef}>
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Tenant
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to delete tenant "{tenantToDelete}"? This action cannot be undone.
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>Cancel</Button>
              <Button colorScheme="red" onClick={handleDelete} ml={3}>
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </Box>
  )
}
