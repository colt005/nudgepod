# nudgepod

## The problem

Updating a ConfigMap in Kubernetes does not automatically update running pods. After editing the ConfigMap, you must locate each deployment or stateful set that uses it and restart it by hand. This manual workflow is cumbersome, error prone, and easy to forget. NudgePod removes this burden by detecting ConfigMap changes and nudging pods to reload on their own.


## Description
NudgePod keeps your workloads in sync by watching for any changes to ConfigMap resources and nudging pods to pick up the latest configuration. When a ConfigMap is edited, NudgePod finds the deployments or stateful sets that rely on it and triggers a smooth restart.


## Roadmap

- [X] Figure out how to package this and deploy and install using a single command
- [ ] ConfigMap selection – add an annotation or flag to control which ConfigMaps to watch.
- [ ] Reload controls – allow app label specifications, introduce a reload delay, and other fine-grained options.
- [ ] Namespace filtering – enable users to include or exclude specific namespaces.


## Installation

### Using Helm (Recommended)

The easiest way to install NudgePod is using Helm. NudgePod is available as a Helm chart.

#### 1. Add the Helm Repository

```bash
helm repo add nudgepod https://colt005.github.io/nudgepod
helm repo update
```

#### 2. Install NudgePod

```bash
helm upgrade --install nudgepod nudgepod/nudgepod
```

This will install NudgePod in the default namespace. To install in a specific namespace:

```bash
helm upgrade --install nudgepod nudgepod/nudgepod --namespace nudgepod --create-namespace
```

#### 3. Verify Installation

Check that the NudgePod controller is running:

```bash
kubectl get pods -l app.kubernetes.io/name=nudgepod
```

You should see the NudgePod controller pod in a `Running` state.

#### 4. Uninstall NudgePod

To remove NudgePod from your cluster:

```bash
helm uninstall nudgepod
```

If you installed in a specific namespace:

```bash
helm uninstall nudgepod -n nudgepod
```

### Configuration Options

You can customize the installation by overriding the default values:

```bash
helm upgrade --install nudgepod nudgepod/nudgepod \
  --set controllerManager.manager.image.tag=v0.1.1 \
  --set controllerManager.replicas=2
```

#### Available Configuration Options

- `controllerManager.manager.image.repository`: Container image repository (default: `rohan005/nudgepod`)
- `controllerManager.manager.image.tag`: Container image tag (default: `latest`)
- `controllerManager.replicas`: Number of controller replicas (default: `1`)
- `controllerManager.manager.resources`: Resource limits and requests for the controller
- `metricsService.type`: Service type for metrics endpoint (default: `ClusterIP`)

### Manual Installation

If you prefer to install manually without Helm, YOU CANNOT, just use helm :)


### Prerequisites

- Kubernetes cluster v1.11.3+
- kubectl configured to communicate with your cluster
- Helm v3.0+ (for Helm installation)

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/nudgepod:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/nudgepod:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/nudgepod:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/nudgepod/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.


**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
