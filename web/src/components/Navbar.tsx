'use client'

import { Box, Flex, Heading, Button, HStack, useColorModeValue } from '@chakra-ui/react'
import { FiPlus } from 'react-icons/fi'
import Link from 'next/link'

export function Navbar() {
  return (
    <Box bg={useColorModeValue('white', 'gray.800')} px={4} boxShadow="sm">
      <Flex h={16} alignItems="center" justifyContent="space-between">
        <HStack spacing={8} alignItems="center">
          <Heading as="h1" size="md">NeuralLog Admin</Heading>
        </HStack>
        <Flex alignItems="center">
          <Link href="/tenants/new" passHref>
            <Button
              variant="solid"
              colorScheme="blue"
              size="sm"
              mr={4}
              leftIcon={<FiPlus />}
            >
              New Tenant
            </Button>
          </Link>
        </Flex>
      </Flex>
    </Box>
  )
}
