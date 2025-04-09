# Getting Started with NeuralLog Admin

This guide will help you set up and start using the NeuralLog Admin interface.

## Prerequisites

Before you begin, ensure you have the following:

- A Kubernetes cluster (v1.19+)
- kubectl configured to access your cluster
- Node.js (v16+) and npm (v7+)
- Docker (for local development)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/NeuralLog/infra.git
cd infra/web
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Configure Environment Variables

Create a `.env.local` file in the `web` directory:

```
# Kubernetes API URL (leave empty to use in-cluster config)
KUBERNETES_API_URL=

# Namespace for the operator (default: default)
OPERATOR_NAMESPACE=default

# Base domain for tenant URLs (e.g., example.com)
BASE_DOMAIN=example.com
```

### 4. Install the Tenant Operator

The NeuralLog Admin interface requires the Tenant Operator to be installed in your Kubernetes cluster.

```bash
# Navigate to the operator directory
cd ../operator

# Install the CRD
kubectl apply -f config/crd/bases/neurallog.io_tenants.yaml

# Install the operator
kubectl apply -f config/samples/operator.yaml
```

### 5. Build and Run the Admin Interface

```bash
# Navigate back to the web directory
cd ../web

# Build the application
npm run build

# Start the application
npm start
```

The application will be available at http://localhost:3000.

## Quick Start

### 1. Create Your First Tenant

1. Open the NeuralLog Admin interface in your browser (http://localhost:3000).
2. Click the "New Tenant" button.
3. Fill in the tenant details:
   - **Tenant Name**: A unique identifier for the tenant (e.g., `my-tenant`)
   - **Display Name**: A human-readable name (e.g., `My Tenant`)
   - **Description**: A brief description of the tenant
4. Configure the server and Redis settings as needed.
5. Click "Create Tenant".

### 2. View Tenant Details

1. On the homepage, you'll see a list of all tenants.
2. Click on a tenant name to view its details.
3. The details page shows the tenant's status, server and Redis configurations, and other information.

### 3. Edit a Tenant

1. On the homepage, click the "Edit" button for the tenant you want to modify.
2. Update the tenant's configuration as needed.
3. Click "Update Tenant" to save your changes.

### 4. Delete a Tenant

1. On the homepage, click the "Delete" button for the tenant you want to remove.
2. Confirm the deletion when prompted.

## Next Steps

- Read the [User Manual](./user-manual.md) for detailed information on using the NeuralLog Admin interface.
- Explore the [Architecture Overview](./architecture.md) to understand how the system works.
- Check out the [API Documentation](./api-docs.md) if you want to integrate with the NeuralLog API.
