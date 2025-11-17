# Cloud Site Manager - Comprehensive Technical Documentation

## Table of Contents

1. Executive Summary
2. Architecture Overview
3. Core Components
4. Site Lifecycle Management
5. Security Architecture
6. API Reference
7. Data Models and Custom Resource Definitions
8. Bootstrap and Handshake Protocol
9. Certificate Management
10. Observability and Monitoring
11. Integration Strategy: Incorporating carbide-rest-api Components
12. Deployment and Configuration

---

## Executive Summary

The Cloud Site Manager, commonly referred to as CSM, is a Kubernetes-native microservice that serves as the control plane orchestrator for NVIDIA's Carbide distributed cloud infrastructure. Its primary purpose is to manage the secure onboarding, authentication, and lifecycle of remote edge sites, which may be data centers, GPU clusters, or distributed computing locations, into the centralized cloud platform.

### Key Responsibilities

The Site Manager handles four critical functions. First, it manages site registration by onboarding new remote sites into the cloud infrastructure. Second, it provides secure bootstrap capabilities through one-time password based credential exchange for initial site authentication. Third, it handles certificate provisioning by issuing and managing mutual TLS certificates for site-to-cloud communication. Finally, it tracks site lifecycle states using Kubernetes Custom Resources to maintain a comprehensive view of all sites in the system.

The service also integrates sites with their respective service providers and FleetCommand organizations, ensuring proper organizational mapping and access control.

### Technology Stack

The Site Manager is built using Go as the primary programming language. For HTTP routing, it uses Gorilla Mux, while command-line interface functionality is provided by the urfave cli framework. State management is handled through Kubernetes Custom Resource Definitions, which provide a declarative, persistent storage mechanism.

Certificate management operations are delegated to a separate credentials manager service backed by HashiCorp Vault. For observability, the system integrates OpenTelemetry for distributed tracing and Sentry for error tracking and reporting.

Security is paramount, with the system implementing mutual TLS for all communications, one-time password based authentication for initial handshakes, and time-limited credentials that expire after a defined period.

---

## Architecture Overview

### High-Level Design

The Cloud Site Manager operates within the cloud control plane as a central orchestration service. At its core, the manager exposes three main API groups: administrative APIs at the v1 site endpoints for managing site lifecycle, a credentials API at v1 sitecreds for bootstrap operations, and system APIs for health checks and metrics.

The Site Manager core logic handles all business operations, delegating storage to a Kubernetes Custom Resource Definition client. This client communicates with the Kubernetes API Server, where Site CRD data is persisted. The CRD storage maintains the site specification, status information, and one-time password state.

Running alongside the Site Manager within the control plane is a credentials manager service backed by Vault. This service manages the PKI certificate authority, handles certificate issuance requests, and enforces time-to-live policies on issued certificates.

Remote edge sites connect to the Site Manager to complete their initial bootstrap process. They present their one-time password to obtain client certificates, then establish ongoing mutual TLS connections for all subsequent communication.

### Component Interaction Flow

The typical flow begins when a cloud administrator creates a new site via the REST API. The Site Manager generates a one-time password and stores the site definition in a Kubernetes Custom Resource Definition. The edge site receives this one-time password through an out-of-band channel, typically secure email or an administrative portal.

Using this one-time password, the edge site requests credentials from the Site Manager. The manager validates the password and, if valid, requests a client certificate from the Vault service. Upon receiving the certificate bundle, which includes the private key, client certificate, and CA certificate, the edge site configures mutual TLS and establishes a secure connection to the cloud.

Finally, the edge site registers itself as operational, allowing the Site Manager to transition it to a fully registered state. At this point, the site is considered part of the cloud infrastructure and can participate in regular operations.

---

## Core Components

### SiteMgr - The Manager Core

The SiteMgr structure is the central orchestrator that encapsulates all site management functionality. It maintains configuration options, a Kubernetes CRD client for state persistence, an HTTP service for the API server, and a dedicated HTTP client for Vault communication. Additionally, it holds a structured logger for operational visibility and a network listener for handling incoming connections.

The initialization process follows a specific sequence. First, it connects to the Kubernetes API using in-cluster configuration, which provides authentication through the pod's service account. The transport layer is wrapped with OpenTelemetry instrumentation for distributed tracing capabilities.

If configured, Sentry is initialized for error reporting, allowing the operations team to track and respond to issues in production. The manager then acquires TLS certificates, either from provided file paths or by requesting them from the Vault service. Finally, it registers all HTTP handlers for the various API endpoints and starts both the OpenTelemetry daemon and the HTTP service.

### HTTP Handlers

The Site Manager implements two categories of handlers: administrative handlers and bootstrap handlers.

#### Admin Handlers

The Create Handler processes POST requests to the v1 site endpoint. It validates the site creation request, generates a twenty-character one-time password with a twenty-four hour expiration, creates a Site Custom Resource Definition in Kubernetes, and sets the initial state to Await Handshake.

The Delete Handler removes sites from the system by deleting their Custom Resource Definition from Kubernetes. It properly handles cases where the site doesn't exist, returning appropriate error codes.

The Get Handler retrieves site details from the CRD storage and returns comprehensive information including the site specification, current status, one-time password, and expiration timestamp.

The Register Handler transitions sites to the Registration Complete state. This is typically called by the site agent after it has successfully established control plane connectivity and is ready for production operations.

The Roll Handler generates a new one-time password for existing sites and resets their state to Await Handshake. This is useful for credential rotation, recovery from failed bootstraps, or when a one-time password expires before use.

#### Credentials Handler

The Credentials Handler is the most security-critical component of the system. It performs a series of validation checks before issuing certificates. First, it verifies that the site exists in the CRD storage. Second, it ensures the site is in the correct state, specifically Await Handshake, preventing one-time password reuse. Third, it compares the provided one-time password with the stored value to detect forgery attempts. Finally, it verifies the timestamp hasn't exceeded the expiration window.

When all validations pass, the handler requests a client certificate from Vault with a ninety-day time-to-live. It also fetches the CA certificate needed for mutual TLS validation. The response includes the private key, client certificate, and CA certificate as a complete bundle. The handler then transitions the site state to Handshake Complete, marking the one-time password as used and preventing any reuse attempts.

### Kubernetes CRD Integration

The Site Custom Resource Definition serves as the source of truth for all site state. It uses standard Kubernetes object metadata and consists of two main sections: the specification and the status.

The specification contains immutable site properties including the unique site identifier, a human-readable name, the service provider name, and the FleetCommand organization mapping. These properties define what the site is and its organizational relationships.

The status section tracks mutable operational state. It includes the one-time password information with both the passcode and expiration timestamp, the current bootstrap state tracking the handshake progress, and the control plane status indicating operational health.

Using Custom Resource Definitions provides several advantages. The declarative state model aligns with Kubernetes' reconciliation patterns. High availability and persistence come from etcd-backed storage with automatic replication. RBAC integration provides Kubernetes-native access control. Audit trails are built-in through Kubernetes event logging. Watch support enables real-time state change notifications for monitoring and automation.

### Certificate Manager Integration

The Site Manager delegates all Public Key Infrastructure operations to a separate credentials manager service backed by HashiCorp Vault. This separation of concerns improves security and maintainability.

For certificate acquisition, the manager sends HTTP POST requests to the PKI cloud cert endpoint with the certificate type, application identifier, and desired time-to-live in hours. The response includes both the private key and signed certificate.

To retrieve the CA certificate, the manager sends a simple GET request to the PKI CA PEM endpoint. This endpoint is typically public since the CA certificate must be distributed to all sites for trust establishment.

The system issues two types of certificates. Server certificates for the Site Manager itself have a one-year time-to-live and use the ingress hostname as the common name. Client certificates for edge sites have a ninety-day time-to-live and embed the site UUID as the application identifier.

The TLS setup logic follows a specific pattern. First, it checks if certificates were provided via command-line flags. If not, it enters a retry loop to acquire certificates from Vault, with a ten-second backoff between attempts. Once acquired, it writes the certificate and key to temporary files and configures the HTTP server to use them.

---

## Site Lifecycle Management

### State Machine

Sites progress through a well-defined state machine during their lifecycle. Understanding this flow is crucial for operations and troubleshooting.

The journey begins with site creation initiated through the admin API. This immediately places the site in the Await Handshake state, where it has a valid one-time password with a twenty-four hour expiration window. During this state, the site is waiting for the edge location to retrieve credentials.

When the edge site presents the one-time password in a credentials request, the system validates it and issues certificates. The site then transitions to Handshake Complete state. In this state, the certificate has been issued and the one-time password is marked as used, preventing any reuse.

From Handshake Complete, the site can progress to Registration Complete by calling the registration endpoint. This indicates the site has successfully established control plane connectivity and is fully operational. The one-time password is no longer accessible in this state.

At any point after creation, an administrator can roll the one-time password, which resets the site back to Await Handshake state with a fresh password. This provides a recovery mechanism for various failure scenarios.

### State Definitions and Valid Transitions

The Await Handshake state indicates an initial condition after creation or password roll. The one-time password is valid and unused. From this state, the only valid transition is to Handshake Complete upon a successful credentials request.

The Handshake Complete state means credentials have been issued and the site has a valid certificate. The one-time password has been consumed. This state can transition to Registration Complete when the site confirms operational status.

The Registration Complete state represents a fully operational site where the one-time password is no longer accessible. This state can transition back to Await Handshake if an administrator rolls the password for re-provisioning.

### Operations by State

Different operations are permitted or denied based on the current state. In Await Handshake, credentials requests are allowed and will succeed if the password matches and hasn't expired. Duplicate credentials requests are denied to prevent double-issuance. One-time password rolls are permitted to allow recovery from various issues.

In Handshake Complete state, credentials requests are denied since the password has already been used. Registration is allowed to move the site to operational status. One-time password rolls are still permitted for emergency re-provisioning scenarios.

In Registration Complete state, credentials requests are denied as the site is already provisioned. Duplicate registration attempts are denied as redundant. One-time password rolls are allowed to support re-provisioning if the site needs to be rebuilt or recovered.

---

## Security Architecture

### Multi-Layer Security Model

The Site Manager implements defense in depth through multiple overlapping security mechanisms.

#### One-Time Password Based Initial Authentication

The password generation process uses cryptographically secure random number generation to create a twenty-character passcode. With sixty-two possible characters per position, this provides approximately one hundred and nineteen bits of entropy, making brute force attacks computationally infeasible.

Each password expires after twenty-four hours, limiting the attack window for compromised credentials. The passwords are single-use, with state transitions preventing any reuse even within the expiration window. The generation relies on the crypto rand package to ensure true randomness rather than pseudorandom values.

The validation process follows a strict sequence of checks, failing fast to prevent timing attacks. First, verify the site exists in the CRD storage. Second, confirm the site is in the correct Await Handshake state. Third, compare the provided passcode with the stored value using constant-time comparison. Finally, verify the timestamp hasn't exceeded the expiration window.

This design mitigates several attack vectors. Brute force attacks are impractical due to high entropy and TLS overhead on each attempt. Replay attacks are prevented by the state transition making each password single-use. Time-based attacks are mitigated by the twenty-four hour expiration. Man-in-the-middle attacks are prevented by TLS protecting the password in transit.

#### Mutual TLS Authentication

The server side, meaning the Cloud Site Manager itself, uses a certificate from Vault with the ingress hostname as the common name. This certificate has a one-year validity period with automated rotation capabilities. All communication uses TLS version one point two or higher with strong cipher suites.

The client side, meaning edge sites, receives certificates issued after one-time password validation. These have a ninety-day validity period to encourage regular rotation. The site UUID is embedded as an application identifier for tracking and authorization purposes.

The trust model centers on a single certificate authority. All certificates are signed by this central CA, and the CA certificate is distributed during the bootstrap process. This enables bidirectional authentication on every connection, with both parties verifying each other's identity.

#### Kubernetes Role-Based Access Control

The CRD client uses an in-cluster service account with minimal permissions following the principle of least privilege. The service account can only perform operations on Site resources within its designated namespace. It has get, list, watch, create, update, and delete permissions on sites, plus additional get, update, and patch permissions on the sites status subresource.

This ensures the Site Manager can only manage Site Custom Resource Definitions and cannot access other cluster resources or namespaces, limiting the blast radius of any potential compromise.

#### Transport Security

The system enforces HTTPS-only communication with no HTTP fallback option. If TLS certificate acquisition fails, the service refuses to start rather than operating in a degraded mode. This fail-closed approach ensures security properties are never compromised.

Communication with the Vault service uses TLS, though certificate verification is skipped since it's a local service within the trusted cluster network. All requests include timeout protection set to twenty seconds, preventing resource exhaustion from hung connections. Retry logic with backoff prevents thundering herd problems.

#### Observability Security

The Sentry integration uses a configurable DSN provided via environment variable, allowing different configurations per environment. Breadcrumbs are collected for debugging context but exclude personally identifiable information. Error contexts are sanitized before transmission to prevent credential leakage.

OpenTelemetry distributed tracing includes span naming for component identification but excludes sensitive data from span attributes. Request correlation is maintained without exposing security-critical information in trace metadata.

---

## API Reference

### Admin Endpoints

All administrative endpoints require authentication, typically enforced by an API gateway or ingress controller in front of the Site Manager.

#### Create Site

The create site endpoint accepts POST requests to the v1 site path. The request must include a JSON body with the site UUID, optional human-readable name, service provider name, and FleetCommand organization.

A successful response returns HTTP status two hundred and one Created with no body. Administrators should check logs to retrieve the generated one-time password.

Error conditions include four hundred Bad Request for invalid JSON or when a site with that UUID already exists, and five hundred Internal Server Error for CRD creation failures.

The operation creates a Site Custom Resource Definition in Kubernetes, generates and stores a one-time password in the status section, and sets the bootstrap state to Await Handshake.

#### Get Site

The get site endpoint accepts GET requests to the v1 site path followed by the site UUID. No request body is required.

A successful response returns HTTP status two hundred OK with a JSON body containing the site UUID, name, provider, FleetCommand organization, bootstrap state, control plane status, one-time password, and OTP expiry timestamp.

Error conditions include four hundred Bad Request for missing UUID, four hundred and four Not Found when the site doesn't exist, and five hundred Internal Server Error for CRD retrieval failures.

#### Delete Site

The delete site endpoint accepts DELETE requests to the v1 site path followed by the site UUID. No request body is required.

A successful response returns HTTP status two hundred OK with no body. The site is removed from the system.

Error conditions include four hundred Bad Request for missing UUID, four hundred and four Not Found when the site doesn't exist, and five hundred Internal Server Error for CRD deletion failures.

Important note: this operation does NOT revoke any previously issued certificates. If the site had obtained credentials, those certificates remain valid until their natural expiration.

#### Roll OTP

The roll OTP endpoint accepts POST requests to the v1 site roll path followed by the site UUID. No request body is required.

A successful response returns HTTP status two hundred OK with no body. Administrators should check logs for the newly generated one-time password.

Error conditions include four hundred Bad Request for missing UUID, four hundred and four Not Found when the site doesn't exist, and five hundred Internal Server Error for OTP generation or CRD update failures.

The operation generates a new one-time password with a fresh twenty-four hour expiration, resets the state to Await Handshake, and invalidates the previous one-time password.

Common use cases include recovering from failed initial bootstrap attempts, handling expired one-time passwords before use, responding to suspected password compromise, and site re-provisioning scenarios.

#### Register Site

The register site endpoint accepts POST requests to the v1 site register path followed by the site UUID. No request body is required.

A successful response returns HTTP status two hundred OK with no body.

Error conditions include four hundred Bad Request for missing UUID, four hundred and four Not Found when the site doesn't exist, and five hundred Internal Server Error for state update failures.

The operation transitions the site state to Registration Complete, indicating the site is fully operational. This endpoint is typically called by the site agent after establishing control plane connectivity.

### Bootstrap Endpoint

#### Request Credentials

The request credentials endpoint accepts POST requests to the v1 sitecreds path. The request must include a JSON body with the site UUID and one-time password.

A successful response returns HTTP status two hundred OK with a JSON body containing three PEM-encoded strings: the private key, the client certificate, and the CA certificate.

Error conditions include four hundred Bad Request for invalid JSON, and five hundred Internal Server Error with various specific messages. The message "site UUID not found" indicates the site doesn't exist in the system. "OTP already used" means the site is not in Await Handshake state. "Bad OTP" indicates passcode mismatch. "OTP expired" means the timestamp has exceeded the twenty-four hour window. "Error getting creds, check logs" indicates Vault communication failure. "Error getting CA, check logs" indicates CA retrieval failure.

Upon success, several side effects occur. A client certificate is issued from Vault with a ninety-day time-to-live. The site state transitions to Handshake Complete. The one-time password is marked as used and cannot be reused.

Security considerations are critical for this endpoint. The password is single-use, becoming invalid after a successful call. Edge sites must securely store the private key with appropriate file permissions and encryption at rest. The certificate expires in ninety days, requiring planning for rotation procedures.

### System Endpoints

#### Health Check

The health check endpoint accepts GET requests to the health path. A successful response returns HTTP status two hundred OK with a JSON body containing a status field set to "healthy".

#### Version Info

The version info endpoint accepts GET requests to the version path. A successful response returns HTTP status two hundred OK with a JSON body containing version information.

#### Metrics

The metrics endpoint accepts GET requests to the metrics path. The response is in Prometheus text format, including HTTP request counters by method, path, and status code.

---

## Data Models and Custom Resource Definitions

### Kubernetes Site CRD

The Site Custom Resource follows standard Kubernetes patterns with API version, kind, metadata, spec, and status sections.

The metadata includes the object name following the pattern "site-" followed by the UUID, and the namespace where the Site Manager operates, typically "csm".

The specification section is immutable after creation and contains the site UUID, human-readable site name, service provider name, and FleetCommand organization. These fields define the site's identity and organizational relationships.

The status section is mutable and tracks operational state. It contains one-time password information with both the passcode and RFC three three three nine formatted timestamp, the bootstrap state indicating progress through the handshake sequence, and control plane status showing operational health.

The object naming convention uses "site-" as a prefix followed by the full UUID, making objects easily identifiable and preventing naming conflicts.

Access patterns use the CRD client with standard Kubernetes operations. Creates use the Sites interface with a Create call. Gets use the Sites interface with a Get call. Updates use the Sites interface with an Update call. Status updates use the Sites interface with an UpdateStatus call. Deletes use the Sites interface with a Delete call.

The status subresource separation is important for several reasons. It separates desired state in the spec from observed state in the status. This allows different RBAC policies for spec versus status updates. It prevents accidental spec modification during status updates. Controllers can update status without triggering spec validation.

### REST API Models

The Site Create Request contains fields for site UUID, name, provider, and FleetCommand organization. All fields are optional in the JSON schema but site UUID is required for the operation to succeed.

The Site Get Response returns comprehensive site information including UUID, name, provider, FleetCommand organization, bootstrap state, control plane status, one-time password, and OTP expiry. This provides administrators with a complete view of site status.

The Site Credentials Request requires only the site UUID and one-time password, keeping the authentication surface minimal.

The Site Credentials Response contains the three components needed for mutual TLS: the private key in PEM format, the client certificate in PEM format, and the CA certificate in PEM format for server validation.

---

## Bootstrap and Handshake Protocol

### Detailed Sequence

The bootstrap process involves coordination between multiple components across cloud and edge locations.

#### Phase One: Site Provisioning

The cloud administrator initiates the process by sending a POST request to the v1 site endpoint with the site details. The API gateway routes this request to the Site Manager service.

The Site Manager authenticates to Kubernetes using its in-cluster service account, acquiring a bearer token for subsequent API requests. It then acquires its own TLS certificate from Vault, requesting a server certificate with a one-year time-to-live using the ingress hostname as the common name. This enables HTTPS for the API server.

Next, the Site Manager creates a Site Custom Resource Definition in Kubernetes. The object name follows the pattern "site-" plus the UUID, and the spec is populated from the request while the status is initially empty.

The one-time password generation follows, creating a twenty-character random string with a twenty-four hour expiration timestamp. This information is written to the Site status section, and the bootstrap state is set to Await Handshake.

The administrator retrieves the one-time password by issuing a GET request for the site details. The response includes the full site information including the password and expiry time. The administrator then communicates this password to the edge site operator through a secure out-of-band channel such as encrypted email, a secure web portal, SSH session, or in some cases physical delivery.

#### Phase Two: Site Bootstrap

The edge site operator receives the site UUID and one-time password through the secure channel. Using these credentials, the site agent sends a POST request to the v1 sitecreds endpoint with the UUID and password in the JSON body.

The Site Manager performs critical security validations in sequence. First, it verifies the site exists by querying the CRD storage. Second, it ensures the site is in Await Handshake state, preventing password reuse. Third, it compares the provided password with the stored value using constant-time comparison to prevent timing attacks. Finally, it verifies the timestamp hasn't exceeded the expiration window.

If all validations pass, the manager requests a client certificate from Vault. The request specifies certificate type "client", the site UUID as the application identifier, and a time-to-live of two thousand one hundred and sixty hours, which equals ninety days. Vault responds with a private key and signed certificate in PEM format.

The manager also retrieves the CA certificate through a separate GET request to the PKI CA PEM endpoint. Edge sites need this CA certificate to validate the cloud server certificate during mutual TLS handshakes.

With both certificates acquired, the manager transitions the site status to Handshake Complete, marking the one-time password as used. It then returns the complete credential bundle to the edge site.

#### Phase Three: Site Registration

The edge site stores all three certificate components securely, typically in protected files with restricted permissions and optionally encrypted at rest.

The site then configures its HTTP client or service mesh for mutual TLS. This involves loading the certificate keypair, creating a CA certificate pool from the provided CA cert, and configuring the TLS settings with both the client certificate and root CA pool.

Using this mutual TLS configuration, the edge site sends a POST request to the v1 site register endpoint followed by its UUID. This confirms the site has successfully established connectivity and is ready for operational traffic.

The Site Manager responds by updating the site status to Registration Complete, indicating the site is fully operational. At this point, the one-time password is no longer accessible or needed for future operations.

---

## Certificate Management

### Certificate Lifecycle

The system manages two distinct certificate types with different lifecycles and purposes.

#### Server Certificate

The server certificate authenticates the Cloud Site Manager to incoming clients. Its common name matches the ingress hostname, such as "sitemanager.forge.nvidia.com". The certificate has a validity period of three hundred and sixty-five days, or one year. Certificates are stored at the filesystem paths /tmp/mgr.cert and /tmp/mgr.key. Administrators can override automatic acquisition by providing custom certificate paths via the tls-cert-path and tls-key-path command-line flags.

The acquisition flow follows a specific pattern. First, the system checks whether user-provided certificates exist at the configured paths. If certificates are provided, it uses them immediately without contacting Vault. If not, it enters an acquisition loop where it requests a certificate from Vault with a retry mechanism. Failed requests result in a ten-second sleep before retrying, continuing indefinitely until successful.

Once acquired, the certificates are written to disk at the configured paths with appropriate file permissions. The HTTP server is then configured to use these certificates for all incoming connections.

Rotation strategies vary based on deployment model. For manual rotation, administrators delete the pod and allow Kubernetes to create a replacement, which automatically acquires a fresh certificate. For automated rotation, a Kubernetes CronJob can trigger rolling restarts before certificate expiry.

#### Client Certificate

Client certificates authenticate edge sites to cloud services. The common name is simply "client" for all sites. The subject alternative name, stored in the application field, contains the site UUID for identification and authorization purposes. These certificates have a validity period of ninety days, or two thousand one hundred and sixty hours, encouraging regular rotation.

All client certificates are signed by the cloud PKI certificate authority, establishing a chain of trust. The acquisition process requires a POST request to the credentials manager service with the certificate name set to "client", the application field set to the site UUID, and the time-to-live set to two thousand one hundred and sixty hours.

The response includes both the private key and the signed certificate in PEM-encoded format. Edge sites must store these securely and protect the private key with appropriate file permissions.

Rotation strategies focus on automation since manual intervention at each edge site would be operationally expensive. Site agents should monitor certificate expiry dates and automatically request renewals before expiration. The ninety-day time-to-live encourages regular rotation as a best practice. For re-provisioning scenarios, administrators can roll the one-time password and have the site request new certificates through the normal bootstrap flow.

### Vault Integration Details

The credentials manager exposes two primary PKI endpoints. The issue certificate endpoint accepts POST requests at v1 pki cloud-cert. It requires mutual TLS authentication using the Cloud Site Manager's server certificate. Vault may impose rate limiting to prevent abuse.

The request body specifies the certificate name, which is either "client" for edge sites or a hostname for server certificates. The application field contains the site UUID for client certificates or is empty for server certificates. The time-to-live field specifies validity in hours.

The response contains the private key and certificate, both in PEM-encoded format.

The second endpoint retrieves the CA certificate. It accepts GET requests at v1 pki ca pem. This is typically a public endpoint requiring no authentication since the CA certificate must be distributed widely. The response is a PEM-encoded CA certificate that clients use to validate server certificates.

Error handling in the client code follows a defensive pattern. Network errors from connection failures or timeouts are returned immediately to the caller. HTTP status codes outside the two hundred range trigger error responses with the status and body content included. A special case handles HTTP status two hundred and four, No Content, by returning nil without error. All other successful responses have their body content returned to the caller.

Retry logic varies by context. During initial certificate acquisition at startup, the system retries infinitely with ten-second backoff until successful, since the service cannot operate without certificates. During runtime certificate requests for edge sites, the system fails fast and returns errors to the caller, allowing the caller to decide on retry strategies.

---

## Observability and Monitoring

### Logging

The logging framework uses Logrus for structured logging with configurable log levels and rich context fields.

Log levels follow standard severity classifications. Info level captures normal operations such as site creation, one-time password rolls, and successful registration. Warn level indicates non-fatal issues like one-time password duration overrides or configuration warnings. Error level signals errors requiring investigation, including CRD operation failures, Vault communication errors, and certificate acquisition problems. Fatal level indicates the service cannot start due to critical failures like Kubernetes connection failure or TLS setup errors.

Sample log messages illustrate the level of detail captured. An info message might state "Site created" followed by the full site details including UUID, name, provider, and FleetCommand organization. Another info message would report "expire at" followed by the UTC timestamp. When processing credentials requests, an info message states "Getting creds for site" followed by the UUID. Upon successful registration, an info message reports the site object name followed by "registration complete".

### Distributed Tracing with OpenTelemetry

The system instruments three key integration points for distributed tracing. All HTTP server requests are traced, providing end-to-end visibility into request processing. Kubernetes API calls through the CRD client are instrumented to track state storage operations. Certificate requests to the Vault client are traced to monitor PKI operations.

Span hierarchies provide context for complex operations. For example, a site creation request generates a parent span named "csm-site-create". This has child spans for "site-crd-client: Create Site", "site-crd-client: Get Site", and "site-crd-client: UpdateStatus".

The implementation wraps the Kubernetes client transport with OpenTelemetry instrumentation. A custom span name formatter ensures consistent naming as "site-crd-client" for all Kubernetes operations. Similarly, the Vault HTTP client is wrapped with OpenTelemetry, using "cert-manager-client" as the span name.

HTTP handlers are wrapped at registration time. Each handler is wrapped with an OpenTelemetry handler using an operation-specific name like "csm-site-create". If Sentry is configured, an additional Sentry wrapper is applied on top for error tracking.

Trace context propagation uses standard W3C Trace Context headers, allowing traces to span multiple services. Context is propagated to both the Kubernetes API and Vault service, enabling end-to-end request correlation across the entire infrastructure.

### Error Tracking with Sentry

When configured with a DSN, Sentry initialization occurs at service startup. The configuration enables debug mode for verbose logging and attaches stack traces to all reported errors.

The wrapper implementation creates a Sentry hub clone for each request, providing isolation between concurrent requests. The hub scope captures the HTTP request context, including method, path, and headers. A deferred panic recovery handler catches any panics, reports them to Sentry with full context, and then re-panics to allow proper cleanup.

Error reporting is automatic for all panics, capturing the panic value, stack trace, and request context. This allows the operations team to diagnose and respond to production issues quickly.

### Metrics with Prometheus

The metrics endpoint exposes several key measurements. The cloud_site_manager_requests_total counter tracks total HTTP requests, labeled by method, path, and status code. The cloud_site_manager_request_duration_seconds histogram measures request latency distribution. The cloud_site_manager_request_size_bytes histogram tracks incoming request sizes. The cloud_site_manager_response_size_bytes histogram monitors outgoing response sizes.

These metrics enable monitoring dashboards, alerting rules, and capacity planning based on actual usage patterns.

### Health and Readiness Probes

The health endpoint at the /health path provides a simple liveness check. Kubernetes can probe this endpoint to determine if the pod is alive and should remain running.

The liveness probe configuration typically uses an HTTPS GET request to the health endpoint on port eight thousand one hundred. It allows thirty seconds for initial startup, then checks every ten seconds. The timeout is set to five seconds with a failure threshold of three consecutive failures before restart.

Readiness is implicitly determined by successful startup. The service only begins accepting traffic after successful TLS setup and all initialization completes. For more robust readiness checking, organizations could add explicit Kubernetes CRD connectivity verification to ensure the service can reach the API server before marking itself ready.

---

## Integration Strategy: Incorporating carbide-rest-api Components

This section explores how to enhance the Site Manager by selectively integrating components from the carbide-rest-api service. The strategy explicitly excludes multi-site management functionality, which the Site Manager already handles through its core CRD operations.

### Why Integrate carbide-rest-api Components

The carbide-rest-api provides several mature, production-ready capabilities that could significantly enhance the Site Manager's functionality and operational characteristics.

First, it offers sophisticated JWT-based authentication supporting multiple identity providers including NGC, Keycloak, SSA, and custom token formats. This would enable proper access control for administrative operations.

Second, it provides comprehensive audit logging with request and response capture, including intelligent obfuscation of sensitive fields. This addresses compliance requirements for many regulatory frameworks.

Third, it includes Postgres-backed persistence using a data access object pattern for structured data. This would enable user management and audit trail storage beyond what CRDs provide.

Fourth, it implements advanced middleware for CORS support, security headers, body dumping for debugging, and structured logging with rich context. These are production-hardened implementations following best practices.

Fifth, it manages user accounts, roles, and organization mappings, providing a foundation for more sophisticated authorization models.

Finally, it integrates Temporal workflow orchestration for complex, long-running operations that require reliable execution guarantees.

### Integration Approach: Phased Adoption

The integration follows a phased approach, starting with low-risk, high-value components and progressing to more complex integrations.

#### Phase One: Middleware and Security Hardening

This phase integrates three key middleware components: CORS middleware from the middleware package, security headers middleware for protection against common web vulnerabilities, and structured logger middleware for improved log formatting with zerolog.

The benefits include enhanced security posture through multiple defense mechanisms, cross-origin API access enabling web-based frontends, and improved log formatting with structured fields and consistent formats.

The implementation requires replacing Gorilla Mux with Echo framework for consistency with the cloud-api codebase. The security middleware stack includes recovery middleware to catch panics, CORS middleware with appropriate allowed origins, security headers middleware setting HSTS and other protections, and structured logger middleware for request logging.

The Site Manager routes are refactored to use Echo handler patterns. All existing functionality is preserved while gaining the benefits of the enhanced middleware stack.

The migration effort is estimated at two to three days. This includes refactoring handlers from Gorilla Mux to Echo patterns, testing CORS configuration with actual frontend code, and validating that log format changes don't break existing monitoring or alerting.

#### Phase Two: Authentication and User Management

This phase brings in JWT authentication middleware from the carbide-rest-auth package, JWT origin configuration supporting multiple identity providers, user data access objects from the carbide-rest-db package, and Postgres database session management.

The current system has a significant gap: the Site Manager has no authentication on admin endpoints, relying entirely on network-level access controls. This phase addresses that gap.

The benefits are substantial. Secure admin API access uses industry-standard JWT bearer tokens. Support for NGC, Keycloak, SSA, and custom identity providers enables integration with existing identity infrastructure. User tracking in audit logs provides accountability for all operations. Role-based access control foundations enable future authorization enhancements.

The architecture changes shift from a direct flow of admin to Site Manager to CRD, to a flow where the JWT authentication middleware sits between the admin and Site Manager, with the middleware querying a User DAO backed by Postgres for user information.

The database schema requires a new users table with fields for unique identifier, auxiliary ID from the identity provider, first and last names, email address, starfleet ID for internal systems, organization data in JSON format, timestamp of last workflow update, and standard created and updated timestamps. Indexes on auxiliary ID and starfleet ID ensure performant lookups.

The configuration system expands to include database connection parameters with host, port, database name, username, and password. JWT origin configuration specifies trusted issuers and their validation parameters. All sensitive values are sourced from environment variables or Kubernetes secrets.

Handler updates demonstrate the integration pattern. For example, the create site handler extracts the authenticated user from the request context, which was populated by the authentication middleware. It uses this user information for audit logging and potentially for authorization decisions. Site CRD annotations record which user created the site, providing a permanent audit trail.

The migration effort is estimated at five to seven days. This includes setting up Postgres database instances and running schema migrations, integrating the JWT middleware into the handler stack, updating all handlers to the Echo framework if not done in phase one, configuring JWT issuers for NGC, Keycloak, and other providers based on organizational requirements, and comprehensive integration testing with real tokens from each identity provider.

#### Phase Three: Audit Logging

This phase integrates audit middleware components from the middleware package and audit entry data access objects from the database package.

The benefits address compliance and operational requirements. A complete audit trail captures all admin actions with full context. Compliance frameworks like SOC2 and ISO 27001 often require such audit trails. Forensic investigation capabilities enable security teams to reconstruct attack timelines. Request and response body capture, with intelligent obfuscation, provides debugging capabilities without exposing sensitive data.

The database schema adds an audit_entries table tracking endpoint path, query parameters in JSON format, HTTP method, status code, status message for errors, client IP address, organization name from request context, user ID foreign key reference, timestamp of the request, request duration, API version, and request body in JSON format with obfuscation applied. Indexes on user ID, timestamp, and endpoint path enable efficient querying of audit data.

The implementation applies audit middleware to the versioned route group. The audit body middleware captures request and response bodies, applying obfuscation rules. The audit log middleware captures metadata like timing and status codes. These run before the authentication middleware, ensuring even failed authentication attempts are logged.

The obfuscation configuration extends the existing rules to include the one-time password field, preventing password leakage in audit logs even if the entire request body is captured.

Sample audit entries show the comprehensive information captured. Each entry records the full request path, HTTP method, success or error status, client IP for network forensics, authenticated user reference, precise timestamp for timeline reconstruction, request duration for performance analysis, API version for compatibility tracking, and request body with sensitive fields replaced by asterisks.

The migration effort is estimated at two to three days. This includes creating the audit_entries table with appropriate indexes, enabling the middleware in the route configuration, configuring obfuscation fields based on security review, and testing audit log generation for all endpoint types to ensure complete coverage.

#### Phase Four: Advanced Features

This optional phase considers more complex integrations requiring significant effort but providing advanced capabilities.

Temporal workflow integration could orchestrate complex site onboarding workflows involving multiple steps such as DNS record creation, firewall rule configuration, and site agent deployment. The effort is estimated at ten to fifteen days due to the complexity of workflow design and testing.

Keycloak integration for authorization would enable fine-grained role-based access control. For example, provider admins could manage only sites belonging to their provider organization. The effort is estimated at five to seven days including role mapping and policy enforcement implementation.

Pagination and filtering capabilities would support large-scale deployments with hundreds or thousands of sites. The pagination component from the api package provides cursor-based pagination and filter query parsing. The effort is estimated at three to five days for implementation and testing.

Advanced error handling using standardized error response structures with error codes would improve API consistency and client error handling. The effort is estimated at two to three days for implementation across all handlers.

### Integration Roadmap

The phased rollout follows a week-by-week schedule.

Week one focuses on middleware and security improvements. This includes implementing CORS, security headers, and structured logging, migrating to the Echo framework, and requires three to five days of effort.

Weeks two and three address authentication and user management. This involves Postgres database setup, JWT middleware integration, user DAO and authentication flow implementation, and requires seven to ten days of effort.

Weeks three and four implement audit logging. This includes creating audit tables and DAOs, implementing request and response capture, and configuring obfuscation rules. The effort is three to five days.

Month two and beyond, optionally, tackles advanced features including Temporal workflow orchestration, RBAC with Keycloak, and pagination and filtering. The combined effort is twenty to thirty days.

### Benefits Summary

Phase one middleware integration provides strong security improvements and moderate auditability improvements with minimal complexity. It's a low-risk, high-value starting point.

Phase two authentication integration provides excellent security improvements through proper access control, strong auditability through user tracking, good scalability through standard authentication protocols, but moderate complexity requiring database integration.

Phase three audit logging provides moderate direct security improvements but excellent auditability for compliance and investigations, with moderate scalability impact and low complexity.

Phase four advanced features provide strong security through fine-grained RBAC, strong auditability through comprehensive workflow tracking, excellent scalability supporting large deployments, but high complexity requiring significant development and testing effort.

### Key Considerations

Several important factors must be considered during integration.

Backward compatibility requires maintaining the existing v1 sitecreds endpoint without authentication, since edge sites may not have identity tokens during bootstrap. A migration path must be provided for existing integrations. If breaking changes are necessary, versioned API endpoints like v2 should be introduced.

Database dependencies add Postgres as a new infrastructure requirement. This requires database migration tooling, which could leverage the existing carbide-rest-db migrations command. High availability and backup procedures must be established for the database.

Configuration management should move from command-line flags to configuration files in YAML or JSON format for complex settings. Environment variables should be supported for Kubernetes deployments following twelve-factor app principles. ConfigMaps can store JWT issuer configuration, allowing updates without pod restarts.

Testing strategy must include unit tests for all new handlers with comprehensive coverage, integration tests with real Postgres and Kubernetes backends, and end-to-end tests with actual JWT token generation from supported identity providers.

Documentation updates are essential. API documentation should be provided in Swagger or OpenAPI format. Deployment guides need updates for new database requirements. A security configuration guide should explain JWT issuer setup and RBAC policies.

### Example: Complete Authenticated Handler Pattern

A fully integrated handler demonstrates all the patterns discussed. The create site handler extracts the authenticated user from the request context, which was populated by the authentication middleware. If no user is found, it returns an unauthorized error.

Request parsing uses the Echo binding mechanism for automatic JSON unmarshalling. Validation checks ensure required fields are present and properly formatted, such as verifying the site UUID is a valid UUID format.

One-time password generation remains unchanged, using the existing cryptographically secure random generation. The expiry timestamp is calculated and formatted.

Site CRD creation now includes rich annotations recording the user who created the site, their full name for human readability, and the creation timestamp. Labels on the CRD enable efficient querying by provider or organization.

Error handling provides specific, actionable error messages. A conflict status indicates the UUID is already in use. Internal server errors include enough context for debugging without exposing sensitive internals.

Structured logging records comprehensive context including site identification, user identification with both UUID and name, bootstrap state, and one-time password expiry. This creates a complete audit trail even before the database audit log is written.

The response returns a created status with useful information including the site UUID for confirmation, a success message, the one-time password so the administrator can provision the edge site, and the password expiry timestamp for planning purposes. HTTPS protects this response in transit, and the audit log will obfuscate the password in storage.

This example demonstrates how all the integrated components work together to create a secure, auditable, user-friendly administrative interface while maintaining the core site management functionality.

---

## Deployment and Configuration

### Kubernetes Deployment

The Site Manager deploys to Kubernetes using standard resource types.

The namespace definition creates an isolated environment named "csm" for all site manager resources.

Service account and RBAC configuration follows the principle of least privilege. A service account named "site-manager" provides the pod identity. A role in the csm namespace grants permissions only on forge.nvidia.com API group resources. Specifically, it allows get, list, watch, create, update, patch, and delete verbs on sites resources, and get, update, and patch verbs on sites status subresources. A role binding connects the service account to the role.

The deployment specification creates a highly available configuration with two replicas. The pod template includes the site-manager service account assignment. The container runs the site manager image from NVIDIA's container registry. It exposes port eight thousand one hundred for HTTPS traffic.

Environment variables configure the namespace through SITE_MANAGER_NS and the Sentry DSN through a secret reference for security. Command arguments specify listen port, credentials manager URL pointing to the in-cluster service, ingress hostname for certificate generation, and OTP duration in hours.

Liveness and readiness probes both query the health endpoint over HTTPS. The liveness probe starts after thirty seconds and checks every ten seconds. The readiness probe starts after ten seconds and checks every five seconds. These ensure Kubernetes only routes traffic to healthy pods.

Resource requests and limits prevent resource contention. Requests guarantee one hundred millicores CPU and one hundred twenty-eight megabytes memory. Limits cap usage at five hundred millicores CPU and five hundred twelve megabytes memory.

The service definition creates a cluster IP service named "site-manager" in the csm namespace. It selects pods with the app label "site-manager" and exposes port eight thousand one hundred.

The ingress configuration uses cert-manager for automatic TLS certificate provisioning with Let's Encrypt. Annotations configure nginx to use HTTPS for backend communication and force SSL redirect for all traffic. The TLS section specifies the hostname and secret name for certificate storage. Rules route all traffic for sitemanager.forge.nvidia.com to the site-manager service on port eight thousand one hundred.

### Configuration Options

Command-line flags and environment variables provide flexible configuration.

The listen-port flag defaults to eight thousand one hundred and specifies the HTTPS port the service listens on.

The ingress-host flag defaults to sitemanager.forge.nvidia.com and specifies the hostname for server certificate generation.

The creds-manager-url flag defaults to https://localhost:8000 and specifies the Vault service URL for certificate operations.

The tls-cert-path and tls-key-path flags have no default, causing automatic acquisition from Vault. When specified, they provide paths to existing server certificates.

The namespace flag maps to the SITE_MANAGER_NS environment variable, defaults to "csm", and specifies the Kubernetes namespace for CRD operations.

The otp-duration flag defaults to twenty-four hours and specifies one-time password validity in hours.

The sentry-dsn flag maps to the SENTRY_DSN environment variable, has no default meaning disabled, and specifies the Sentry error tracking DSN when enabled.

The debug flag defaults to false and enables debug-level logging when set to true.

### High Availability Considerations

Production deployments should implement several high availability patterns.

Multiple replicas provide redundancy, with two or more instances recommended for production. This ensures service availability during rolling updates and pod failures.

A pod disruption budget ensures at least one pod remains available during voluntary disruptions like cluster upgrades or node maintenance. The specification sets minAvailable to one and selects pods with the app label "site-manager".

Pod anti-affinity spreads replicas across different nodes. The required during scheduling rule prevents multiple site manager pods on the same node using the kubernetes.io/hostname topology key. This protects against node-level failures.

Graceful shutdown allows in-flight requests to complete before termination. The CLI code includes a five-second grace period after receiving the shutdown signal, allowing active requests to finish before the process exits.

### Security Hardening

Production deployments should implement additional security controls beyond the base configuration.

Running as non-root prevents privilege escalation attacks. The security context sets runAsNonRoot to true, runAsUser to one thousand, and fsGroup to one thousand for proper file permissions.

A read-only root filesystem prevents runtime modification of the container image. The security context sets readOnlyRootFilesystem to true. A writable tmpfs volume is mounted at /tmp for certificates and temporary files that must be writable.

Dropping all capabilities reduces the attack surface by removing unnecessary Linux capabilities. The security context capabilities drop list includes ALL, removing every optional capability.

Network policies implement defense in depth at the network layer. The policy selects site-manager pods and defines both ingress and egress rules. Ingress allows traffic from the ingress-nginx namespace on port eight thousand one hundred. Egress allows traffic to kube-dns for name resolution, to credsmgr pods for certificate operations, and to the Kubernetes API server for CRD operations. All other traffic is implicitly denied.

---

## Conclusion

The Cloud Site Manager is a critical infrastructure component that bridges cloud control planes with distributed edge sites. Its design prioritizes several key principles.

Security is paramount, implemented through one-time password based bootstrap preventing credential theft, mutual TLS ensuring bidirectional authentication, and time-limited credentials reducing the impact of compromise.

Reliability comes from Kubernetes-native state management using etcd-backed CRDs, high availability deployment patterns with multiple replicas, and graceful degradation when dependent services are unavailable.

Observability provides operational visibility through OpenTelemetry distributed tracing correlating requests across services, Sentry error tracking for proactive issue detection, and Prometheus metrics for capacity planning and alerting.

Simplicity keeps the system maintainable with a focused API surface covering only essential operations, a clear state machine with well-defined transitions, and minimal dependencies reducing operational complexity.

By selectively integrating components from carbide-rest-api, specifically authentication for access control, audit logging for compliance, and advanced middleware for security, the Site Manager can evolve into an enterprise-grade site lifecycle management platform while maintaining its core strengths of security, reliability, and operational simplicity.

The phased integration approach allows organizations to adopt capabilities as needed, starting with low-risk middleware improvements and progressing to comprehensive authentication and audit logging. Each phase delivers incremental value while building toward a complete solution.

As the Carbide cloud infrastructure grows, the Site Manager will continue serving as the secure gateway for edge site onboarding, ensuring every site joins the cloud with proper authentication, authorization, and audit controls in place.
