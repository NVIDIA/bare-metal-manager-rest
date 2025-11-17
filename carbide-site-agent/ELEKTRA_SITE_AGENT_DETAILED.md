# Elektra Site Agent - Comprehensive Technical Documentation

## Table of Contents

1. Executive Summary
2. Architecture Overview
3. Core Components
4. Bootstrap Process
5. Workflow Orchestration
6. Carbide gRPC Communication
7. Resource Management
8. Health and Monitoring
9. Configuration Management
10. Deployment Architecture
11. Integration with Cloud Services
12. Operational Considerations

---

## Executive Summary

The Elektra Site Agent is a sophisticated site management framework that operates at the edge of NVIDIA's Carbide distributed cloud infrastructure. It serves as the on-premises orchestrator responsible for managing data center resources, provisioning infrastructure, and maintaining bidirectional communication with the centralized cloud control plane.

### Key Responsibilities

The Site Agent fulfills several critical functions in the distributed architecture. First, it manages site workflows by orchestrating complex, multi-step operations for resource provisioning and lifecycle management. Second, it maintains the site data store, persisting workflow state and site-specific information in a Postgres database. Third, it provides northbound integration through Temporal workflows, enabling reliable communication with Forge Cloud services. Fourth, it handles southbound integration with site controllers and infrastructure components through gRPC interfaces.

The agent must be extensible to support new resource types and workflows, highly available to ensure continuous site operations, and scalable to handle growing infrastructure demands.

### Technology Stack

The foundation of the Site Agent is built on Go as the primary programming language, providing excellent concurrency support and performance. Temporal serves as the workflow orchestration engine, ensuring reliable execution of long-running operations with built-in retry and failure handling. Communication with local infrastructure uses gRPC with Protocol Buffers for efficient, type-safe remote procedure calls.

State persistence relies on Postgres as the primary database, storing workflow execution history and site state. Kubernetes provides the deployment platform with StatefulSets for ordered deployment and stable network identities. Observability is achieved through Prometheus metrics for monitoring, OpenTelemetry for distributed tracing, and structured logging with zerolog.

Security implementation includes mutual TLS for all external communications, certificate-based authentication for Temporal connections, and secure credential storage using Kubernetes secrets.

---

## Architecture Overview

### High-Level Design

The Site Agent operates as the central orchestration point at each edge site. It sits between the cloud control plane above and the local site infrastructure below. The agent receives workflow requests from the cloud through Temporal, executes those workflows locally, and communicates with infrastructure components through gRPC.

The architecture consists of several key layers. At the top level, Temporal provides the workflow orchestration interface, connecting to both cloud and site namespaces. The agent subscribes to workflows in the site-specific namespace and publishes workflow results and updates to the cloud namespace.

Within the agent itself, multiple manager components handle different resource types and operations. The Bootstrap Manager handles initial site registration and credential acquisition. The Carbide Manager maintains gRPC connections to local infrastructure services. The Workflow Orchestrator coordinates all Temporal workflow execution. Additional resource managers handle specific infrastructure components such as VPCs, subnets, instances, machines, SSH keys, InfiniBand partitions, and network security groups.

Below the agent, the local infrastructure includes Carbide API services for infrastructure management, site controllers for hardware orchestration, and the actual physical or virtual resources being managed. All communication uses secure channels with mutual TLS authentication.

### Component Interaction Flow

The typical flow begins when cloud administrators or automation systems trigger workflows through the Temporal cloud namespace. These workflows are distributed to the appropriate site based on site UUID and workflow routing rules. The Site Agent, continuously polling its subscribe queue in Temporal, picks up the workflow execution request.

The Workflow Orchestrator examines the workflow type and routes it to the appropriate resource manager. Each manager implements workflows and activities specific to its resource domain. Activities perform actual operations such as creating virtual networks, provisioning compute instances, or configuring network rules.

Operations requiring infrastructure changes flow through the Carbide gRPC client, which maintains persistent connections to local infrastructure services. These services translate high-level requests into low-level infrastructure operations. Results and state updates flow back through the same path, with the Site Agent publishing status updates to the Temporal cloud namespace.

Throughout execution, the agent maintains state in its Postgres database, ensuring workflow history persists across agent restarts. Metrics flow to Prometheus for monitoring and alerting. Distributed traces propagate through OpenTelemetry, providing end-to-end visibility into multi-hop operations.

---

## Core Components

### The Elektra Framework

The Elektra framework serves as the main organizational structure containing all subsystems. At its core is the Elektra struct, which aggregates all manager components, configuration, health status, version information, and the structured logger. This centralized design ensures consistent access to shared resources and simplifies dependency management.

Initialization follows a specific sequence. First, the framework creates the Elektra data structures, allocating memory for all manager components. Second, it loads configuration from environment variables and configuration files, validating all required settings. Third, it initializes the Elektra API, which in turn initializes all individual managers. Fourth, each manager performs its specific initialization, registering metrics, creating clients, and preparing for operation. Finally, the Start phase begins asynchronous operations such as Temporal workers, gRPC connection management, and file watchers.

The framework uses a hierarchical manager pattern where each manager is responsible for a specific domain but can access shared resources through the central Elektra structure. This enables loose coupling between components while maintaining coherent state management.

### Manager Components

The Site Agent implements sixteen specialized managers, each responsible for specific aspects of site operation.

The Bootstrap Manager handles the critical initial site registration process. When the agent first starts, it reads bootstrap credentials from a Kubernetes secret containing the site UUID, one-time password, credentials URL, and CA certificate. It uses these credentials to download client certificates from the Cloud Site Manager, enabling authenticated communication with cloud services. The manager watches the bootstrap secret for changes, automatically re-bootstrapping if credentials are updated. It stores downloaded certificates in another Kubernetes secret or local files for consumption by other components.

The Carbide Manager maintains gRPC connections to local infrastructure services. It initializes the Carbide gRPC client with mutual TLS configuration using certificates and keys from the file system or Kubernetes secrets. The client supports three security modes: insecure gRPC without encryption, server TLS with one-way authentication, and mutual TLS with bidirectional authentication. The manager implements automatic certificate reload by monitoring certificate files for changes and recreating the gRPC client when updates are detected. This enables certificate rotation without agent restart.

The Workflow Orchestrator coordinates all Temporal operations. It creates two separate Temporal clients: a publisher client for the cloud namespace where results are sent, and a subscriber client for the site namespace where work requests arrive. The orchestrator initializes a Temporal worker that polls the subscribe queue, executing workflows and activities as they arrive. It registers all workflow and activity implementations from the various resource managers, ensuring the worker can handle any workflow type sent to the site.

Resource managers handle specific infrastructure types. The VPC Manager orchestrates virtual private cloud creation and management. The Subnet Manager handles network subnet operations within VPCs. The Instance Manager provisions and manages compute instances. The Machine Manager interfaces with physical hardware. The SSH Key Group Manager distributes SSH public keys for access control. The InfiniBand Partition Manager configures high-speed network partitions. The Network Security Group Manager implements firewall rules and security policies. The Tenant Manager handles multi-tenant isolation. The Operating System Manager maintains OS image associations. The Instance Type Manager defines compute resource templates. The Expected Machine Manager tracks anticipated hardware inventory. The SKU Manager maintains product and service definitions.

Each manager follows a common pattern. During initialization, it registers Prometheus metrics for monitoring. During the start phase, it registers its workflows and activities with the Temporal worker. During operation, it responds to workflow execution requests, implementing the business logic for its domain. Publishers allow the agent to initiate workflows in the cloud namespace. Subscribers enable the agent to respond to cloud-initiated workflows.

### HTTP Service Layer

The agent exposes an HTTP service for operational visibility and health checking. The service runs on a configurable port separate from the metrics endpoint. It provides routes for health status, returning the current health state of all managers. It exposes bootstrap state, showing credential download status and configuration. It offers Carbide state, displaying gRPC connection health and error details. Additional routes provide access to each manager's state information for debugging.

The HTTP service uses standard HTTP handlers without additional frameworks, keeping dependencies minimal. Responses use plain text or JSON formats for easy consumption by monitoring tools and operators.

### Metrics and Observability

The agent registers numerous Prometheus metrics for comprehensive monitoring. Version metrics expose the build version and build date as labels with a constant value, allowing dashboards to display version information. Health status metrics provide gauge values for overall agent health, Carbide gRPC health, bootstrap state, and Temporal connection status.

Counter metrics track bootstrap credential downloads attempted and succeeded, gRPC requests succeeded and failed, Temporal workflow connections attempted and succeeded, and various workflow-specific operations. Histogram metrics measure workflow execution duration and activity execution duration, providing percentile information for performance analysis.

All metrics use the "elektra_site_agent" namespace prefix for consistency. Metrics follow Prometheus naming conventions with descriptive help text. The metrics endpoint serves Prometheus text format on the configured metrics port, separate from the main HTTP service to isolate monitoring traffic.

---

## Bootstrap Process

### Initial Site Registration

The bootstrap process is the critical first operation that connects a newly deployed site agent to the cloud infrastructure. This process must complete successfully before the agent can perform any other operations.

The process begins when the agent pod starts in Kubernetes. The Bootstrap Manager reads from a Kubernetes secret mounted at the bootstrap secret path. This secret contains four essential pieces of information: the site UUID uniquely identifying this site, a one-time password provided by cloud administrators, the credentials URL pointing to the Cloud Site Manager, and the CA certificate for validating the credentials server.

The manager validates that all four values are present and non-empty. If any value is missing, the agent refuses to start, logging a fatal error. This fail-fast approach ensures the agent never operates in an invalid state.

### Credential Download

With valid bootstrap configuration, the manager initiates credential download. It constructs an HTTPS client configured with the provided CA certificate in its trust store. For localhost connections during testing, it may skip server verification, but production deployments always validate server certificates.

The manager creates a site credentials request containing the site UUID and one-time password. It marshals this request to JSON and sends it via POST to the credentials URL. The request includes proper headers for content negotiation and distributed tracing integration.

The Cloud Site Manager receives the request, validates the one-time password as described in the Site Manager documentation, and returns a credentials response. This response contains three PEM-encoded values: the private key for client authentication, the client certificate signed by the cloud CA, and the CA certificate for validating cloud services.

The manager implements exponential backoff retry logic for resilience against temporary network issues. Initial attempts use short delays, gradually increasing to prevent overwhelming the credentials server. Maximum retry duration and attempt counts are configurable. Once credentials are successfully downloaded, the manager parses and validates them before storage.

### Credential Storage

Storage behavior depends on the runtime environment. When running in Kubernetes, the manager uses the Kubernetes API to update a secret containing Temporal certificates. It retrieves the existing secret, updates the data map with new certificate values, and writes the secret back. This atomic operation ensures other pods sharing the secret see consistent state.

When running in Docker or development mode, the manager writes certificates to the file system at configured paths. It ensures parent directories exist and sets appropriate file permissions. The OTP value is also stored to prevent redundant downloads for the same credential set.

The manager performs additional validation on downloaded certificates. It parses the CA certificate PEM block, extracts the certificate structure, and reads the expiration date. This expiration is published as a Prometheus metric, enabling monitoring systems to alert before certificate expiry.

### Certificate Rotation

The Bootstrap Manager watches its configuration secret for changes using the file system notify library. When it detects modifications to secret files, it automatically triggers a new credential download cycle. This enables certificate rotation without manual intervention or pod restarts.

The rotation process compares the current OTP value with the stored value from the previous download. If they match, rotation is skipped to avoid unnecessary downloads. If they differ, the full download and storage process executes with the new OTP. This supports several rotation scenarios: scheduled rotations with new OTPs provided by administrators, emergency rotations in response to suspected compromise, and automatic rotations triggered by certificate expiry monitoring.

Success and failure events are tracked as Prometheus metrics. The credentials download attempted counter increments for each attempt. The credentials download succeeded counter increments only for successful completions. Operators can alert on high failure rates or extended periods without successful downloads.

---

## Workflow Orchestration

### Temporal Integration

The Site Agent uses Temporal as its workflow orchestration engine, providing reliable execution guarantees for long-running, complex operations. Temporal's architecture ensures workflows continue execution even through failures, restarts, and network partitions.

The agent creates two separate Temporal client connections. The publisher client connects to the cloud namespace, allowing the agent to initiate workflows that execute in the cloud control plane. This is used for status updates, event notifications, and requests for cloud-side operations. The subscriber client connects to the site-specific namespace, where cloud services send workflows for execution at the site. Each site has a unique namespace identified by its site UUID.

Both clients use identical configuration except for the namespace. Connection options specify the Temporal server hostname and port. When TLS is enabled, connection options include the client certificate key pair and CA certificate for server validation. Data converter configuration ensures proper serialization of Protocol Buffer messages and JSON payloads. Client and worker interceptors enable OpenTelemetry tracing integration. Logger configuration routes Temporal's internal logs to the agent's logging system.

### Workflow Worker

The Temporal worker is the component that actually executes workflows and activities. It polls the subscribe queue in the site namespace, receiving workflow tasks from the Temporal service. When a task arrives, the worker examines the workflow type, locates the registered workflow implementation, and begins execution.

Workflow execution follows Temporal's event sourcing model. Each step in the workflow execution generates events recorded in Temporal's database. If the worker crashes or the pod restarts, a new worker can resume execution from the recorded event history. This provides workflow durability without requiring the agent to manage checkpointing logic.

The worker implements strict workflow and activity separation. Workflows describe the orchestration logic: what steps to execute, in what order, with what retry policies. Activities implement actual operations: making gRPC calls, querying databases, or performing computations. This separation enables Temporal to manage activity retries, timeouts, and compensation logic independently.

Worker configuration specifies several important parameters. The workflow panic policy determines behavior when workflows panic. The agent uses FailWorkflow, which marks the workflow as failed rather than retrying infinitely. Worker interceptors enable tracing and monitoring integration. Concurrent workflow and activity execution limits prevent resource exhaustion.

### Workflow Registration

Each resource manager registers its workflows and activities with the worker during initialization. Registration creates a mapping from workflow type names to implementation functions. When the worker receives a workflow task, it looks up the implementation using this mapping.

The VPC Manager registers workflows for creating VPCs, updating VPC configuration, and deleting VPCs. Each workflow is implemented as a Go function accepting a context and workflow-specific input. The Subnet Manager registers subnet creation, update, and deletion workflows. The Instance Manager registers instance provisioning, modification, and termination workflows.

Activity registration follows the same pattern. Each activity is a function that can be invoked by workflows. Activities for VPC creation might include reserving IP addresses from IPAM, calling the Carbide API to create the network, and updating the database with result status. Activities are reusable across workflows, promoting code reuse and consistency.

### Workflow Execution Flow

When a cloud service needs to perform an operation at a site, it starts a workflow execution in the site's Temporal namespace. The workflow appears in the subscribe queue. The site's worker polls the queue, discovers the new task, and begins execution.

Workflow execution progresses through activities. The workflow code schedules activities, specifying timeout and retry policies. The worker picks up activity tasks from the queue and executes the activity functions. Activities return results to the workflow, which uses those results to determine next steps. This continues until the workflow reaches a terminal state: completed successfully, failed, or cancelled.

Throughout execution, the workflow can emit signals and queries. Signals allow external systems to send events to running workflows. Queries enable external systems to inspect workflow state. The agent uses these mechanisms to provide real-time visibility into long-running operations.

Upon completion, the workflow returns a result value. The agent publishes this result to the cloud namespace using the publisher client. Cloud services can wait synchronously for workflow completion or poll for results asynchronously. Either way, they receive the final state and any output values.

### Error Handling and Retries

Temporal provides sophisticated error handling that the agent leverages extensively. Activity retry policies specify how many times to retry failed activities, with what delay between attempts, and what maximum delay to allow. Exponential backoff prevents overwhelming failing services. Different error types can trigger different retry behaviors.

Workflows can implement compensation logic for partial failures. If a workflow successfully creates a VPC but fails to create a subnet, the compensation logic can delete the VPC to maintain consistent state. The agent implements such patterns where appropriate to prevent resource leaks.

Workflow and activity timeouts ensure operations don't run indefinitely. Start-to-close timeouts limit total execution time. Schedule-to-start timeouts detect when workers are overwhelmed. Schedule-to-close timeouts limit total latency including queue time. The agent configures these timeouts based on expected operation duration and SLAs.

---

## Carbide gRPC Communication

### gRPC Client Architecture

The Carbide Manager maintains a gRPC client for communication with local infrastructure services. This client provides the interface between high-level workflow operations and low-level infrastructure management.

The client supports three security modes. Insecure gRPC mode uses no encryption or authentication, suitable only for development and testing. Server TLS mode encrypts traffic and verifies the server's identity using certificates, providing confidentiality and server authentication. Mutual TLS mode adds client certificate authentication, ensuring both parties verify each other's identity.

Client configuration comes from environment variables and configuration files. The Carbide address specifies the hostname and port of the infrastructure service. The security option selects one of the three modes. Certificate paths point to the CA certificate, client certificate, and client private key files.

### Client Initialization and Lifecycle

The Carbide Manager initializes its gRPC client during the start phase. It reads certificates from the configured paths, parses them to ensure validity, and calculates MD5 hashes of the certificate contents. These hashes enable detection of certificate rotation.

Client creation involves several steps. First, configure TLS if required, loading certificates and keys. Second, create gRPC dial options specifying keepalive parameters, connection timeout, and interceptors. Third, establish the connection to the Carbide address. Fourth, create service stubs for calling remote methods. Fifth, store the client in the manager's state for access by activities.

The client implements automatic reconnection for transient failures. If the connection is lost, gRPC's internal retry logic attempts reconnection with exponential backoff. The agent monitors connection state through gRPC's connectivity state API, updating health metrics accordingly.

### Certificate Rotation

The Carbide Manager implements automatic certificate rotation without requiring agent restart. A background goroutine periodically checks certificate file MD5 hashes. If hashes differ from the initial values, certificate rotation has occurred.

Upon detecting rotation, the manager creates a new gRPC client with the updated certificates. It uses an atomic swap operation to replace the old client with the new one. Ongoing operations using the old client complete normally. New operations use the new client. Once all operations complete, the old client's connection is gracefully closed.

This approach ensures zero-downtime certificate rotation. Operations in flight are not interrupted. New operations immediately use fresh certificates. The system remains secure even during the rotation window.

### Error Handling and Health Tracking

The Carbide Manager tracks gRPC health through success and failure counters. Each gRPC call increments either the success counter or failure counter based on the result. Health status is derived from recent error patterns and error codes.

Certain error codes indicate unhealthy states. Unavailable errors suggest the service is down or unreachable. Unauthenticated errors indicate certificate or authentication problems. These codes cause the health status to transition to unhealthy, triggering alerts and preventing new operations.

Other error codes represent application-level errors rather than connectivity problems. InvalidArgument indicates bad request parameters. NotFound means a requested resource doesn't exist. These don't affect health status since the connection is working properly.

Error details are logged for debugging. The most recent error message is stored in manager state, accessible via the HTTP service. Prometheus metrics provide time-series data for error rates and patterns. This multi-level visibility enables rapid problem diagnosis.

---

## Resource Management

### Resource Manager Pattern

All resource managers follow a common architectural pattern, providing consistency across the codebase and simplifying maintenance. Each manager implements five key operations: initialization, publisher registration, subscriber registration, cron job scheduling, and HTTP state exposure.

During initialization, the manager registers Prometheus metrics specific to its resource type. It may initialize local caches, create client connections, or perform other one-time setup. Initialization errors are typically fatal, preventing the agent from starting in an invalid state.

Publisher registration creates the ability to initiate workflows in the cloud namespace. The manager registers workflow implementations and specifies their type names. When the manager wants to trigger a cloud workflow, it uses the Temporal publisher client to start execution.

Subscriber registration enables the manager to respond to cloud-initiated workflows. The manager registers workflow and activity implementations with the Temporal worker. When workflows arrive in the subscribe queue, the worker routes them to the appropriate manager based on type.

Cron job scheduling sets up periodic operations. Some managers need to periodically poll infrastructure state, reconcile configuration, or publish inventory updates. The cron functionality schedules these operations at configured intervals.

HTTP state exposure provides operational visibility. Each manager implements a GetState method returning human-readable status information. This includes operation counts, error details, and configuration values. The main HTTP service aggregates this information for operator access.

### VPC Management

The VPC Manager handles virtual private cloud operations, orchestrating network isolation and connectivity. VPC creation workflows allocate IP address ranges from IPAM, call the Carbide API to create the network fabric, configure routing tables, and update the site database with VPC metadata. Activities perform each operation with individual retry policies and timeouts.

VPC update workflows modify existing network configurations. They might adjust IP ranges, change routing rules, or update tags and metadata. Updates are applied atomically where possible to prevent inconsistent states. Where atomicity isn't possible, the workflow implements rollback logic.

VPC deletion workflows follow a carefully ordered sequence. First, enumerate all resources within the VPC to ensure it's empty. Second, delete associated routing tables and network ACLs. Third, delete the VPC itself via Carbide API. Fourth, release IP allocations back to IPAM. Fifth, mark the VPC as deleted in the database. This ordering prevents orphaned resources and ensures clean removal.

### Instance Management

The Instance Manager provisions and manages compute instances across the site infrastructure. Instance creation involves selecting appropriate physical hardware based on instance type requirements, allocating IP addresses from the appropriate subnet, configuring user data for initialization scripts, attaching storage volumes, applying SSH keys for access, and configuring network security groups.

Each step is implemented as a Temporal activity. If any step fails, retry logic attempts recovery. If retries are exhausted, the workflow enters a failed state. The cloud service can examine the failure reason and decide whether to retry the entire workflow or abandon the operation.

Instance lifecycle operations include starting, stopping, rebooting, and terminating instances. Each operation translates to a sequence of Carbide API calls. Start operations power on instances and wait for them to become accessible. Stop operations gracefully shut down the operating system before powering off. Reboot operations combine stop and start with appropriate delays. Terminate operations permanently remove instances, freeing all associated resources.

### SSH Key Distribution

The SSH Key Group Manager distributes public SSH keys to instances for access control. It maintains mappings between key groups and instance sets. When a key group is created or updated, the manager identifies all instances that should receive those keys. It generates workflow tasks for each instance, ensuring keys are deployed consistently.

Key distribution workflows connect to instances using the Carbide API, inject public keys into authorized keys files, set appropriate file permissions, and verify successful injection. Verification ensures operators can actually use the keys to access instances.

Key rotation follows a similar pattern. When keys are removed from a group, workflows remove them from all associated instances. When new keys are added, workflows distribute them. This centralized key management prevents security issues from manual key distribution.

### Network Security

The Network Security Group Manager implements firewall rules and security policies. Security groups define allowed ingress and egress traffic using rules specifying protocol, port range, and source or destination IP ranges. When security groups are created or modified, the manager translates rules into infrastructure-specific firewall configurations.

Application of security groups to instances happens through workflows. The workflow identifies the instance, retrieves current security group associations, computes the delta between desired and actual state, and applies necessary changes via Carbide API calls. This idempotent approach allows repeated application without side effects.

Security group updates are particularly sensitive since they control access to instances. The manager implements careful ordering to prevent accidentally blocking critical traffic. When adding more restrictive rules, it validates that management traffic remains allowed. When removing rules, it checks that required traffic patterns still have allow rules.

---

## Health and Monitoring

### Multi-Level Health Checks

The Site Agent implements health checking at multiple levels, providing granular visibility into system state. Component-level health tracks individual manager health. The Bootstrap Manager health reflects whether credentials have been successfully downloaded. The Carbide Manager health indicates gRPC connection status. The Workflow Orchestrator health shows Temporal connectivity.

Each component reports health as an enumeration: Healthy indicates normal operation. Unhealthy indicates a problem requiring attention. Not Known indicates health status cannot be determined, often during initialization. This tri-state model provides more nuance than simple binary health.

Agent-level health aggregates component health. The agent is healthy only if all critical components are healthy. Critical components include Bootstrap, Carbide, and Workflow Orchestrator. Non-critical components can be unhealthy without affecting overall agent health, though they still generate alerts.

Health status is exposed through multiple channels. A Prometheus gauge publishes numeric health values for each component and overall agent health. The HTTP health endpoint returns health status in JSON format with details for each component. Kubernetes liveness and readiness probes query the health endpoint to determine pod health.

### Prometheus Metrics

The metrics exposition provides comprehensive operational visibility. Version information metrics expose build version and build date, allowing operators to verify deployed versions. Component health metrics track Bootstrap Manager status, Carbide Manager status, and Workflow Orchestrator status.

Operation counters track the volume of operations over time. Bootstrap credential downloads attempted and succeeded show credential acquisition patterns. Carbide gRPC requests succeeded and failed indicate infrastructure interaction rates. Temporal workflow executions started and completed reveal orchestration activity.

Resource-specific metrics track operations for each resource type. VPC creation, update, and deletion counters show network provisioning activity. Instance creation and termination counters indicate compute scaling patterns. SSH key distribution operations demonstrate access control changes.

Histogram metrics capture duration distributions. Workflow execution duration shows how long workflows take to complete, with percentiles revealing outliers. Activity execution duration indicates which activities are slow or experiencing issues. HTTP request duration tracks API responsiveness.

### Distributed Tracing

OpenTelemetry integration provides end-to-end request tracing across the distributed system. When a cloud service initiates a workflow, it creates a trace context containing trace ID and span ID. This context propagates to the Site Agent through Temporal's header mechanism.

The agent extracts trace context from incoming workflows and creates child spans for local operations. Each activity execution creates a span, recording start time, end time, attributes, and events. Attributes capture operation-specific metadata like resource IDs and types. Events record significant occurrences like errors or state transitions.

When the agent calls the Carbide API via gRPC, it propagates trace context in gRPC metadata. The Carbide service creates its own child spans, continuing the trace. When all operations complete, the full trace shows the request path from cloud initiation through agent processing to infrastructure execution.

Trace collection can use various backends. The agent supports Lightstep, Jaeger, Zipkin, and other OpenTelemetry-compatible systems. Configuration via environment variables selects the backend and provides access credentials. Sampling rates control what percentage of traces are collected, balancing observability with overhead.

### Structured Logging

The agent uses zerolog for structured logging, enabling efficient log analysis and correlation. All log messages include timestamp, level, and message text. Additional fields provide context specific to the operation.

Log levels follow standard conventions. Debug level captures detailed information useful for troubleshooting. Info level records normal operational events like workflow starts and completions. Warn level indicates unexpected but non-critical conditions. Error level reports failures requiring investigation. Fatal level indicates unrecoverable errors that cause agent shutdown.

Contextual fields enrich log entries with operational metadata. Manager name identifies which component generated the log. Workflow ID and run ID enable correlation with Temporal execution. Resource IDs associate logs with specific infrastructure elements. Error details provide stack traces and root cause information.

Log output goes to standard error, allowing container orchestrators to collect and route logs. JSON formatting enables automated parsing and indexing. Log aggregation systems like Elasticsearch, Splunk, or CloudWatch can ingest agent logs alongside logs from other system components, enabling cross-component analysis.

---

## Configuration Management

### Configuration Sources

The Site Agent loads configuration from multiple sources with a defined precedence order. Environment variables provide the primary configuration mechanism, suitable for containerized deployments. Kubernetes ConfigMaps and Secrets supply credentials and sensitive values. Command-line flags enable overrides for development and testing. Default values provide sensible fallbacks for optional settings.

Configuration loading follows a validation-heavy approach. The agent reads all configuration sources, merges values according to precedence, validates required values are present, validates value formats and ranges, and fails fast if any validation fails. This prevents the agent from starting in an invalid state that would cause runtime failures.

### Key Configuration Parameters

Database configuration specifies connection parameters for the Postgres database. Server address and port identify the database host. Username and password provide authentication credentials. Database name specifies which database to use. Connection pool settings control maximum connections and idle timeouts.

Temporal configuration defines connectivity to the workflow service. Host and port specify the Temporal server address. Publish and subscribe namespace names route workflows appropriately. Publish and subscribe queue names determine which task queues to use. Certificate paths enable TLS authentication. Server name enables server identity validation.

Carbide configuration controls gRPC communication. Address specifies the Carbide service location. Security option selects insecure, server TLS, or mutual TLS mode. Certificate paths provide CA, client certificate, and client key. Skip server authentication flag disables server validation for testing.

Bootstrap configuration enables initial site registration. Bootstrap secret path points to the Kubernetes secret containing UUID and OTP. Temporal secret path points to where downloaded certificates should be stored. Disable bootstrap flag prevents bootstrap on pods that don't need it.

Operational configuration controls runtime behavior. Debug mode enables verbose logging for troubleshooting. Development mode enables features useful for local development. Metrics port specifies where Prometheus metrics are exposed. Watcher interval controls how often periodic tasks run. Pod name and namespace enable leader election and coordination.

### Pod Identity and Leader Election

The Site Agent runs as a Kubernetes StatefulSet with multiple pods for high availability. Only one pod, the master pod, performs certain singleton operations like bootstrap and scheduling. The configuration system determines which pod is the master.

Pod naming in StatefulSets follows a pattern: the StatefulSet name followed by a hyphen and an ordinal index. The agent parses the pod name environment variable, extracts the ordinal index, and checks if it equals zero. Pod zero is designated the master pod. Other pods are replicas providing redundancy and failover.

The master pod runs the Bootstrap Manager, downloading credentials during startup. Replica pods skip bootstrap, instead reading credentials from the shared secret written by the master. The master pod registers scheduled tasks like inventory polling. Replica pods skip scheduling, avoiding duplicate work.

This approach provides both high availability and correctness. If the master pod crashes, Kubernetes restarts it automatically. If the entire StatefulSet is deleted and recreated, the master pod identity is preserved. Manual scaling operations add replica pods without changing the master.

### Configuration Validation

Validation occurs at multiple stages to catch configuration errors early. Format validation ensures values match expected formats: UUIDs are valid UUIDs, ports are valid port numbers, URLs are well-formed, and required fields are non-empty.

Consistency validation checks relationships between values. If TLS is enabled, certificate paths must be specified. If bootstrap is enabled, secret paths must be configured. If scheduled tasks are enabled, watcher interval must be positive.

Runtime validation confirms configuration actually works. Database connections are tested during startup. Temporal connectivity is verified before starting the worker. Certificate files are read and parsed to ensure they're valid PEM. gRPC connections are established to verify network reachability.

Failed validation produces detailed error messages identifying which configuration parameter is problematic and why. This accelerates troubleshooting compared to generic error messages. The agent logs all configuration values at startup for audit trail and debugging.

---

## Deployment Architecture

### Kubernetes StatefulSet

The Site Agent deploys as a Kubernetes StatefulSet rather than a Deployment due to requirements for stable network identities and ordered pod creation. StatefulSets provide numbered pod names, ensuring the master pod has a predictable identity. They provide stable storage through persistent volume claims. They enable ordered, graceful scaling where pods are created and deleted in sequence.

The StatefulSet spec defines replica count, typically three for production deployments. Pod template specifications include container images, resource requests and limits, environment variables, volume mounts for secrets and config, and affinity rules for distribution.

Service configuration provides network access to agent pods. A headless service enables direct pod-to-pod communication. Individual pod DNS names allow targeting specific replicas. Load-balanced service endpoints distribute traffic across healthy pods.

### Persistent Storage

While the agent primarily uses external Postgres for persistence, local storage serves specific purposes. Certificate storage keeps downloaded certificates accessible across container restarts. Log buffering provides temporary storage for log aggregation systems. Temporary files support download operations and file processing.

Persistent volume claims request storage from Kubernetes. Storage class selection determines the underlying storage technology. Volume size specifies capacity requirements, typically modest for agent workloads. Access mode is read-write-once since volumes are pod-specific.

### Secrets Management

Sensitive data like credentials and certificates are stored in Kubernetes Secrets rather than ConfigMaps or environment variables directly. Bootstrap secrets contain site UUID, one-time password, credentials URL, and CA certificate. These values are injected at deployment time by infrastructure automation.

Temporal certificate secrets contain downloaded credentials. The master pod's Bootstrap Manager writes these values. All pods mount this secret to access certificates for Temporal connectivity. Secret updates trigger automatic remounting in running pods.

Database credentials are stored in a separate secret, following the principle of least privilege. Secrets use appropriate Kubernetes RBAC to restrict access. Service accounts have minimum necessary permissions. Pods can only access secrets in their namespace.

### Container Security

The agent container runs with security hardening to minimize attack surface. The container runs as a non-root user with a high UID. The file system is read-only except for explicitly mounted writable volumes. All capabilities are dropped; no privileged Linux capabilities are granted. Security context settings prevent privilege escalation.

Image security practices include regular base image updates to patch vulnerabilities. Vulnerability scanning in the CI/CD pipeline catches known issues. Image signatures verify image authenticity. Private registry usage prevents tampering.

### High Availability Considerations

Production deployments run three agent replicas for resilience. Pod anti-affinity spreads replicas across different nodes, preventing single node failure from losing all replicas. Multiple availability zones further increase resilience for multi-zone clusters.

The StatefulSet update strategy uses rolling updates to prevent downtime. Pods update one at a time, waiting for the new pod to become ready before updating the next. This ensures at least two pods are always available during updates. Readiness probes prevent routing traffic to pods that aren't fully initialized.

Pod disruption budgets limit voluntary disruptions. The budget ensures at least two pods remain available during cluster maintenance operations. This prevents all pods from being evicted simultaneously during node draining or cluster upgrades.

---

## Integration with Cloud Services

### Temporal Bidirectional Communication

The Site Agent maintains bidirectional communication with cloud services through Temporal. Cloud-to-site communication uses workflows published to the site namespace. Cloud services start workflows specifying the site UUID as the namespace. The site's agent, subscribed to that namespace, receives and executes the workflows. Results are returned through Temporal's response mechanism or via publishing to the cloud namespace.

Site-to-cloud communication uses workflows published to the cloud namespace. The agent starts workflows in the cloud namespace for operations requiring cloud-side execution. Examples include registering site completion after bootstrap, reporting inventory changes detected at the site, and requesting cloud-side resource allocation like IP addresses.

This bidirectional pattern enables flexible orchestration. Cloud services can push work to sites. Sites can pull resources from cloud services. Complex workflows can span both cloud and site, with each side executing steps in its domain. Temporal handles coordination, ensuring exactly-once execution semantics.

### Resource Synchronization

The agent periodically synchronizes local infrastructure state with cloud inventory. Scheduled workflows query local infrastructure for current resource state. They compare local state with cloud database state, identifying discrepancies. Created resources missing from cloud inventory are added. Deleted resources remaining in cloud inventory are marked as deleted. Modified resources with stale cloud data trigger update workflows.

Synchronization handles several edge cases. Resources created during network partitions eventually synchronize when connectivity restores. Resources deleted at the site but not yet reflected in cloud state are reconciled. Drift between desired configuration and actual state triggers remediation workflows.

The synchronization interval balances freshness with overhead. Too frequent synchronization creates excessive load. Too infrequent synchronization delays detection of issues. Default intervals are typically ten to thirty minutes, configurable based on operational requirements.

### Event Publishing

The agent publishes events to the cloud for significant occurrences. Resource lifecycle events announce creation, modification, and deletion of infrastructure resources. Health state changes notify when components transition between healthy and unhealthy states. Error events report failures requiring cloud-side intervention.

Event publishing uses Temporal signals for immediate delivery or workflow initiation for guaranteed processing. Critical events use workflows to ensure delivery even if the cloud side is temporarily unavailable. Non-critical events use signals for lower overhead.

Event schemas use Protocol Buffers for type safety and evolution support. Schema versioning enables rolling updates without breaking compatibility. Cloud services subscribe to events of interest, implementing event-driven architecture patterns.

---

## Operational Considerations

### Day One Operations

Initial deployment requires several preparatory steps. Infrastructure preparation includes provisioning Kubernetes cluster at the site, configuring persistent storage for agent pods, and establishing network connectivity to cloud services. Credential preparation involves obtaining site UUID from cloud administrators, receiving one-time password for bootstrap, and collecting Temporal server connection details. Configuration preparation includes creating ConfigMaps with non-sensitive settings, creating Secrets with credentials and certificates, and validating configuration values before deployment.

Deployment steps follow a specific sequence. First, create the namespace for agent components. Second, apply RBAC roles and service accounts. Third, create ConfigMaps and Secrets. Fourth, deploy the StatefulSet. Fifth, verify pods start successfully. Sixth, check logs for bootstrap completion. Seventh, validate Temporal connectivity. Eighth, run test workflows to confirm operation.

Common day-one issues include certificate trust problems preventing Temporal connectivity, network policies blocking required traffic, insufficient RBAC permissions for Kubernetes API access, and bootstrap secret formatting errors. Structured logging and detailed error messages help diagnose these quickly.

### Ongoing Operations

Routine operational tasks include monitoring metrics for degradation, reviewing logs for error patterns, validating certificate expiration dates, updating container images for security patches, and scaling StatefulSets based on load.

Certificate rotation follows a periodic schedule before expiration. Cloud administrators roll the OTP in the Cloud Site Manager. They update the bootstrap secret with the new OTP. The agent detects the change and automatically downloads fresh certificates. Temporal connectivity transitions smoothly to new certificates without downtime.

Version upgrades use StatefulSet rolling update strategies. New container images are validated in staging environments. Production updates proceed one pod at a time. Readiness checks prevent traffic to incompletely initialized pods. Rollback procedures restore previous versions if issues arise.

### Troubleshooting

The agent provides multiple troubleshooting interfaces. HTTP endpoints expose current state of all managers. Logs contain detailed operation history with correlation IDs. Metrics reveal patterns and trends over time. Distributed traces show request flow across components.

Common problems have standard diagnostic procedures. Temporal connectivity issues involve checking certificate validity, verifying DNS resolution of Temporal server, testing network connectivity to Temporal port, and examining Temporal logs for authentication errors.

Carbide gRPC problems involve validating certificate permissions and paths, testing gRPC server health independently, examining gRPC status codes for error categories, and checking for network policies blocking gRPC traffic.

Workflow failures require examining Temporal workflow history for failed activities, reviewing activity error messages and stack traces, checking resource availability on infrastructure side, and verifying workflow input parameters are valid.

### Disaster Recovery

Site failures requiring agent redeployment follow recovery procedures. If persistent volumes survive, StatefulSet recreation recovers state automatically. If volumes are lost, bootstrap repeats with new credentials. Cloud state remains authoritative during recovery. Workflows resubmit for any operations that were in flight during failure.

The agent's stateless design, where Temporal and Postgres hold critical state, simplifies recovery. Container images are immutable and versioned. Secrets and ConfigMaps can be reconstructed from documentation. Recovery primarily involves redeploying workloads and waiting for synchronization to complete.

For complete site destruction and rebuild scenarios, cloud administrators provision a new site UUID. New bootstrap credentials are generated. The new agent deploys with the new identity. Cloud-side automation reconciles the new site identity with physical infrastructure location and organizational ownership.

---

## Conclusion

The Elektra Site Agent serves as the critical orchestration layer at edge sites in NVIDIA's Carbide distributed cloud infrastructure. Its design emphasizes several key principles that enable reliable, scalable edge operations.

Reliability comes from Temporal workflow orchestration providing exactly-once execution semantics, automatic retries with exponential backoff for transient failures, persistent workflow state surviving agent restarts, and graceful degradation when individual components fail.

Security is implemented through mutual TLS for all external communications, certificate-based authentication preventing unauthorized access, regular certificate rotation without service interruption, and secure credential storage in Kubernetes secrets.

Observability provides operational visibility through comprehensive Prometheus metrics for all operations, distributed tracing correlating requests across services, structured logging enabling automated analysis, and multi-level health checks with component granularity.

Extensibility supports evolving requirements through plugin-style resource managers enabling new resource types, standardized patterns for adding workflows and activities, configuration-driven feature enablement, and schema evolution support via Protocol Buffers.

As the Carbide cloud infrastructure scales to thousands of sites managing diverse infrastructure types, the Elektra Site Agent will continue serving as the reliable, secure, observable bridge between centralized cloud control and distributed edge execution. Its architecture supports this growth while maintaining the simplicity and operational clarity required for successful distributed systems.

