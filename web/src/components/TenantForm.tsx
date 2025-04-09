'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Input,
  VStack,
  Heading,
  Textarea,
  SimpleGrid,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Switch,
  FormHelperText,
  useToast,
} from '@chakra-ui/react'
import { useRouter } from 'next/navigation'

type TenantFormData = {
  name: string
  displayName: string
  description: string
  serverReplicas: number
  serverImage: string
  redisReplicas: number
  redisImage: string
  redisStorage: string
  networkPolicyEnabled: boolean
}

type TenantFormProps = {
  initialData?: Partial<TenantFormData>
  isEditing?: boolean
}

export function TenantForm({ initialData, isEditing = false }: TenantFormProps) {
  const {
    handleSubmit,
    register,
    formState: { errors, isSubmitting },
  } = useForm<TenantFormData>({
    defaultValues: {
      name: initialData?.name || '',
      displayName: initialData?.displayName || '',
      description: initialData?.description || '',
      serverReplicas: initialData?.serverReplicas || 1,
      serverImage: initialData?.serverImage || 'neurallog/server:latest',
      redisReplicas: initialData?.redisReplicas || 1,
      redisImage: initialData?.redisImage || 'redis:7-alpine',
      redisStorage: initialData?.redisStorage || '1Gi',
      networkPolicyEnabled: initialData?.networkPolicyEnabled ?? true,
    },
  })

  const [serverReplicas, setServerReplicas] = useState(initialData?.serverReplicas || 1)
  const [redisReplicas, setRedisReplicas] = useState(initialData?.redisReplicas || 1)
  
  const toast = useToast()
  const router = useRouter()

  const onSubmit = async (data: TenantFormData) => {
    try {
      const tenantData = {
        apiVersion: 'neurallog.io/v1',
        kind: 'Tenant',
        metadata: {
          name: data.name,
        },
        spec: {
          displayName: data.displayName,
          description: data.description,
          server: {
            replicas: data.serverReplicas,
            image: data.serverImage,
          },
          redis: {
            replicas: data.redisReplicas,
            image: data.redisImage,
            storage: data.redisStorage,
          },
          networkPolicy: {
            enabled: data.networkPolicyEnabled,
          },
        },
      }

      const url = isEditing ? `/api/tenants/${data.name}` : '/api/tenants'
      const method = isEditing ? 'PUT' : 'POST'

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(tenantData),
      })

      if (!response.ok) {
        throw new Error(`Failed to ${isEditing ? 'update' : 'create'} tenant`)
      }

      toast({
        title: isEditing ? 'Tenant updated' : 'Tenant created',
        description: isEditing
          ? `Tenant ${data.name} has been updated successfully.`
          : `Tenant ${data.name} has been created successfully.`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      })

      router.push('/')
    } catch (error) {
      toast({
        title: 'Error',
        description: `Failed to ${isEditing ? 'update' : 'create'} tenant. Please try again.`,
        status: 'error',
        duration: 5000,
        isClosable: true,
      })
    }
  }

  return (
    <Box as="form" onSubmit={handleSubmit(onSubmit)} bg="white" p={6} borderRadius="md" shadow="md">
      <VStack spacing={6} align="stretch">
        <Heading as="h2" size="md">
          {isEditing ? 'Edit Tenant' : 'Create New Tenant'}
        </Heading>

        <SimpleGrid columns={{ base: 1, md: 2 }} spacing={6}>
          <FormControl isInvalid={!!errors.name} isRequired isDisabled={isEditing}>
            <FormLabel>Tenant Name</FormLabel>
            <Input
              {...register('name', {
                required: 'Tenant name is required',
                pattern: {
                  value: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
                  message:
                    'Tenant name must consist of lowercase alphanumeric characters or "-", and must start and end with an alphanumeric character',
                },
              })}
              placeholder="my-tenant"
            />
            <FormHelperText>
              This will be used as the tenant identifier and cannot be changed later.
            </FormHelperText>
            <FormErrorMessage>{errors.name?.message}</FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={!!errors.displayName}>
            <FormLabel>Display Name</FormLabel>
            <Input
              {...register('displayName')}
              placeholder="My Tenant"
            />
            <FormHelperText>
              A user-friendly name for the tenant.
            </FormHelperText>
            <FormErrorMessage>{errors.displayName?.message}</FormErrorMessage>
          </FormControl>
        </SimpleGrid>

        <FormControl isInvalid={!!errors.description}>
          <FormLabel>Description</FormLabel>
          <Textarea
            {...register('description')}
            placeholder="Description of the tenant"
            rows={3}
          />
          <FormErrorMessage>{errors.description?.message}</FormErrorMessage>
        </FormControl>

        <Heading as="h3" size="sm" mt={4}>
          Server Configuration
        </Heading>

        <SimpleGrid columns={{ base: 1, md: 2 }} spacing={6}>
          <FormControl isInvalid={!!errors.serverReplicas} isRequired>
            <FormLabel>Server Replicas</FormLabel>
            <NumberInput
              min={1}
              max={10}
              value={serverReplicas}
              onChange={(_, value) => setServerReplicas(value)}
            >
              <NumberInputField
                {...register('serverReplicas', {
                  required: 'Server replicas is required',
                  min: {
                    value: 1,
                    message: 'Minimum 1 replica is required',
                  },
                  max: {
                    value: 10,
                    message: 'Maximum 10 replicas are allowed',
                  },
                })}
              />
              <NumberInputStepper>
                <NumberIncrementStepper />
                <NumberDecrementStepper />
              </NumberInputStepper>
            </NumberInput>
            <FormErrorMessage>{errors.serverReplicas?.message}</FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={!!errors.serverImage} isRequired>
            <FormLabel>Server Image</FormLabel>
            <Input
              {...register('serverImage', {
                required: 'Server image is required',
              })}
              placeholder="neurallog/server:latest"
            />
            <FormErrorMessage>{errors.serverImage?.message}</FormErrorMessage>
          </FormControl>
        </SimpleGrid>

        <Heading as="h3" size="sm" mt={4}>
          Redis Configuration
        </Heading>

        <SimpleGrid columns={{ base: 1, md: 3 }} spacing={6}>
          <FormControl isInvalid={!!errors.redisReplicas} isRequired>
            <FormLabel>Redis Replicas</FormLabel>
            <NumberInput
              min={1}
              max={3}
              value={redisReplicas}
              onChange={(_, value) => setRedisReplicas(value)}
            >
              <NumberInputField
                {...register('redisReplicas', {
                  required: 'Redis replicas is required',
                  min: {
                    value: 1,
                    message: 'Minimum 1 replica is required',
                  },
                  max: {
                    value: 3,
                    message: 'Maximum 3 replicas are allowed',
                  },
                })}
              />
              <NumberInputStepper>
                <NumberIncrementStepper />
                <NumberDecrementStepper />
              </NumberInputStepper>
            </NumberInput>
            <FormErrorMessage>{errors.redisReplicas?.message}</FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={!!errors.redisImage} isRequired>
            <FormLabel>Redis Image</FormLabel>
            <Input
              {...register('redisImage', {
                required: 'Redis image is required',
              })}
              placeholder="redis:7-alpine"
            />
            <FormErrorMessage>{errors.redisImage?.message}</FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={!!errors.redisStorage} isRequired>
            <FormLabel>Redis Storage</FormLabel>
            <Input
              {...register('redisStorage', {
                required: 'Redis storage is required',
                pattern: {
                  value: /^[0-9]+[KMGTPEkmgtpe]i$/,
                  message: 'Invalid storage format (e.g., 1Gi, 500Mi)',
                },
              })}
              placeholder="1Gi"
            />
            <FormErrorMessage>{errors.redisStorage?.message}</FormErrorMessage>
          </FormControl>
        </SimpleGrid>

        <Heading as="h3" size="sm" mt={4}>
          Network Policy
        </Heading>

        <FormControl display="flex" alignItems="center">
          <FormLabel htmlFor="networkPolicyEnabled" mb="0">
            Enable Network Policies
          </FormLabel>
          <Switch
            id="networkPolicyEnabled"
            {...register('networkPolicyEnabled')}
            defaultChecked={initialData?.networkPolicyEnabled ?? true}
          />
        </FormControl>

        <Button
          mt={6}
          colorScheme="blue"
          isLoading={isSubmitting}
          type="submit"
          size="lg"
        >
          {isEditing ? 'Update Tenant' : 'Create Tenant'}
        </Button>
      </VStack>
    </Box>
  )
}
