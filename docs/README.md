# nudgepod

## The problem

Updating a ConfigMap in Kubernetes does not automatically update running pods. After editing the ConfigMap, you must locate each deployment or stateful set that uses it and restart it by hand. This manual workflow is cumbersome, error prone, and easy to forget. NudgePod removes this burden by detecting ConfigMap changes and nudging pods to reload on their own.


## Description
NudgePod keeps your workloads in sync by watching for any changes to ConfigMap resources and nudging pods to pick up the latest configuration. When a ConfigMap is edited, NudgePod finds the deployments or stateful sets that rely on it and triggers a smooth restart.

