# Problem Definition

While many tools are actually available to implement compliance requirements, almost none of them actually focus on tracking the implementation themselves.

The goal of Argus is to overseer and continuously Attest Requirement Implementations across multiple Resources, and feed that information back with observability metrics. 

# Project Objective

To have something Attesting (reconciling) Implementations and Rerquirements continuously, generating metrics to be consumed with observability tools (i.e. be able to create a compliance SLO/SLI)

# Use Cases

## 1 The AWS Account
Use Case Details: [#7](https://github.com/ContainerSolutions/argus/issues/7)

## 2 The WebApp
Use Case Details: [#8](https://github.com/ContainerSolutions/argus/issues/8)

## 3 The Kubernetes
Use Case Details: [#6](https://github.com/ContainerSolutions/argus/issues/6)


# High Level overview

## General Architecture
![General Architecture](pics/arch.png)

## Objects overview
![Object View](pics/argus.drawio.png)

