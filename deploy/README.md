# 

# **Carbide Installation**

| Revision | Date | Notes |
| :---- | :- | :---- |
| 0.1 |  |  |
|  |  |  |

## Context

Below is a **prescriptive, BYO‑Kubernetes bring‑up guide** that is **not site‑specific**. It encodes the **order of operations** and the **exact versions/images we verified**, plus what you must configure if you already operate some of the common services yourself.

**Principles**

* All unknown values remain explicit placeholders like \<REPLACE\_ME\>.

* If you **already run** one of the core services (PostgreSQL, cert‑manager, Temporal), do the “**If you already have it**” checklist for that service.

* If you **don’t**, deploy the **“Reference version we ran”** (images and versions below) and apply the config under “**If you deploy our reference**”.


## **0\) Validated baseline (what we ran)**

**Kubernetes & node runtime**

* Control plane: **Kubernetes v1.30.4** (server)  
* Nodes: kubelet **v1.26.15**, container runtime **containerd 1.7.1**  
* CNI: **Calico v3.28.1** (node & controllers)  
* OS: Ubuntu **24.04.1 LTS**

**Networking**

* Ingress: **Project Contour v1.25.2** (controller) \+ **Envoy v1.26.4** (daemonset)  
* Load balancer: **MetalLB v0.14.5** (controller & speaker)

**Secret & cert plumbing**

* **External Secrets Operator v0.8.6**  
* **cert‑manager v1.11.1** (controller/webhook/CA‑injector)  
  * **Approver‑policy v0.6.3**  
     (Pods present as cert-manager, cainjector, webhook, and policy controller.  )

**State**

* **PostgreSQL**: Zalando **Postgres Operator v1.10.1** \+ **Spilo‑15 image 3.0‑p1** (Postgres **15**)

**Temporal & search**

* **Temporal server 1.22.6** (frontend/history/matching/worker)  
  * **Admin tools 1.22.4**, **UI 2.16.2**  
* **Elasticsearch 7.17.3** (for Temporal visibility)

**Monitoring & telemetry (reference stack \-  OPTIONAL)**

* **Prometheus Operator v0.68.0**; **Prometheus v2.47.0**; **Alertmanager v0.26.0**  
* **Grafana v10.1.2**; **kube‑state‑metrics v2.10.0**  
* **OpenTelemetry Collector v0.102.1**  
* **Loki v2.8.4**  
* **Node exporter v1.6.1**

**Carbide components (what this guide installs)**

* **cloud‑api**: nvcr.io/nvidian/nvforge-devel/cloud-api:v0.2.72 (two replicas)  
* **cloud‑workflow**: nvcr.io/nvidian/nvforge-devel/cloud-workflow:v0.2.30 (cloud‑worker, site‑worker)  
* **cloud‑cert‑manager (credsmgr)**: nvcr.io/nvidian/nvforge-devel/cloud-cert-manager:v0.1.16  
* **elektra-site-agent**: nvcr.io/nvidian/nvforge-devel/forge-elektra:v2025.06.20-rc1-0


## **1\) Order of operations (high level)**

1. **Cluster & networking ready**  
   * Kubernetes v1.26–1.30 (tested on 1.30.4), containerd 1.7.x, Calico (or conformant CNI)  
   * Ingress controller (Contour/Envoy) \+ LoadBalancer (MetalLB or cloud LB) available  
   * Available DNS recursive resolvers and NTP available

2. **Foundation services (in this order)**  
   * **External Secrets Operator (ESO) \- Optional**  
   * **cert‑manager** (Issuers/ClusterIssuers in place)  
   * **PostgreSQL** (DB/role/extension prerequisites below)  
   * **Temporal** (server up; register namespaces)

3. **Carbide Rest components**  
   * Deploy **cloud‑api**, **cloud‑workflow (cloud‑worker & site‑worker)**, **cloud‑cert‑manager (credsmgr)**  
   * **Seed DB and register Temporal namespaces** (cloud, site, then the **site UUID**)  
   * **Create OTP & bootstrap secrets** for **elektra‑site‑agent**; roll restart it

4. **Monitoring (optional but recommended)**  
   * Prometheus operator, Grafana, Loki, OTel, node exporter

The rest of this document details **what to configure** for each shared service and **how the NVCarbide workloads use them**.  

## **2\) External Secrets Operator (ESO)**

**Reference version we ran:** ghcr.io/external-secrets/external-secrets:v0.8.6.

**What you must provide**

* A SecretStore/ClusterSecretStore pointing at your secret backend (e.g. **Vault**, cloud provider secrets), and if applicable a Postgres secret namespace. Vault is optional; this repo does not require it.  
* ExternalSecret objects similar to these (namespaces vary by component):  
  * forge-roots-eso → target secret forge-roots with keys site-root, forge-root  
  * DB credentials ExternalSecrets per namespace (e.g., clouddb-db-eso → forge.forge-pg-cluster.credentials)  
* Ensure an **image pull secret** (e.g., imagepullsecret) exists in namespaces that pull from nvcr.io.


## **3\) cert‑manager (TLS & trust)**

**Reference versions we ran**

* Controller/Webhook/CAInjector **v1.11.1**  
* Approver‑policy **v0.6.3**  
* ClusterIssuers present: self-issuer, site-issuer, carbide-ca-issuer, carbide-ca-issuer

**If you already have cert‑manager**

* Ensure it is ≥ **v1.11.1** and that:  
  * Your **ClusterIssuer** objects can issue:  
    * cluster internal certs (service DNS SANs) and  
    * any externally‑facing FQDNs you choose.  
  * Approver flows allow your teams to create Certificate resources for the NVCarbide namespaces.

**If you deploy our reference**

* Install cert‑manager **v1.11.1** and **approver‑policy v0.6.3**.  
* Create ClusterIssuers matching your PKI: \<ISSUER\_NAME\>.  
* Typical **SANs** for NVFORGE services include:  
  * Internal service names (e.g., carbide-api.\<ns\>.svc.cluster.local, carbide-api.forge)  
  * Optional external FQDNs (your domains)


## **4\) Vault not required**

**HashiCorp Vault is not required** for this deployment. PKI is provided by:

* **carbide-rest-cert-manager** (native Go PKI) — issues certificates for Temporal and site bootstrap from a CA stored in a Kubernetes secret.
* **carbide-ca-issuer** (cert-manager.io ClusterIssuer) — uses the same CA secret (**ca-signing-secret**) so that cert-manager can issue certificates (e.g. for site-manager, site-agent, Temporal client) without a separate Vault server.

See **7.0** (CA secret) and **7.1** (cloud-cert-manager) for how to provision the CA and deploy the cert-manager.


## **5\) PostgreSQL (DB)**

**Reference we ran**

* **Zalando Postgres Operator v1.10.1**  
* **Spilo‑15 image 3.0‑p1** (Postgres **15**)

**If you already have Postgres**

* Provide a database \<POSTGRES\_DB\> and role \<POSTGRES\_USER\> with password \<POSTGRES\_PASSWORD\>.  
* Enable **TLS** (recommended) or allow secure network policy between DB and the NVCarbide namespaces.  
* Create extensions (the apps expect these):

| CREATE EXTENSION IF NOT EXISTS btree\_gin;CREATE EXTENSION IF NOT EXISTS pg\_trgm; |
| :---- |

This can be done with a call similar to the following:

| psql "postgres://\<POSTGRES\_USER\>:\<POSTGRES\_PASSWORD\>@\<POSTGRES\_HOST\>:\<POSTGRES\_PORT\>/\<POSTGRES\_DB\>?sslmode=\<POSTGRES\_SSLMODE\>" \\     \-c 'CREATE EXTENSION IF NOT EXISTS btree\_gin;' \\     \-c 'CREATE EXTENSION IF NOT EXISTS pg\_trgm;' |
| :---- |

* Make the DSN available to workloads via ESO targets (per‑namespace credentials):  
  * Examples: forge.forge-pg-cluster.credentials, elektra-site-agent.elektra.forge-pg-cluster.credentials (names are examples—use your own).

**If you deploy our reference**

* Deploy the Zalando operator and a Spilo‑15 cluster sized for your SLOs.  
  Expose a ClusterIP service on **5432** and surface credentials through ExternalSecrets to each namespace that needs them.


## **6\) Temporal**

**Reference we ran**

* Temporal server **1.22.6** (frontend/history/matching/worker)  
* UI **2.16.2**, admin tools **1.22.4**  
* Frontend service endpoint (cluster‑internal):  
   **temporal-frontend.temporal.svc:7233**

**Required namespaces**

* Base: **cloud, site**  
* Per‑site: the **site UUID** 

**If you already have Temporal**

* Ensure the **frontend gRPC endpoint** is reachable from NVCarbide workloads and present the proper **mTLS**/CA if you require TLS.  
* Register namespaces:

| tctl \--ns cloud namespace registertctl \--ns site  namespace registertctl \--ns \<SITE\_UUID\> namespace register (once you know the site UUID) |
| :---- |

**If you deploy our reference**

* Deploy Temporal as above, expose :7233.  
* Register the same namespaces.


## **7\) What each carbide workload expects \- exact images we used and what resources needs to be applied in what order:**

## **7.0 CA secret for cert‑manager (carbide namespace)**

#### Goal

Before deploying **carbide-rest-cert-manager** and using the **carbide-ca-issuer** ClusterIssuer, a root CA must be provided as a Kubernetes secret. The cert-manager deployment and the ClusterIssuer both use this secret.

**One Secret in the same namespace as the deployment (e.g. carbide):**

* **ca-signing-secret**  
  * `tls.crt` — PEM‑encoded root CA certificate  
  * `tls.key` — PEM‑encoded root CA private key  

The cert-manager deployment in this repo expects these keys under the volume mount `/etc/pki/ca`. You can create the CA with your own PKI tooling. Example:

| openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "/CN=Carbide CA/O=NVIDIA" |
| :---- |

Create the secret in the deployment namespace (e.g. **carbide**). Use keys **tls.crt** and **tls.key**:

| kubectl \-n carbide create secret generic ca-signing-secret \--from-file=tls.crt=./ca.crt \--from-file=tls.key=./ca.key |
| :---- |

Once **ca-signing-secret** exists, **carbide-rest-cert-manager** and **carbide-ca-issuer** can issue certificates.

## **7.1 cloud‑cert‑manager “credsmgr” (carbide-rest-cert-manager)**

### **Role**

carbide-rest-cert-manager (the **credsmgr** deployment) is responsible for issuing mTLS certificates for Temporal access and site bootstrap 

* Uses **native Go PKI** (no Vault server).  
* Reads the root CA from Kubernetes secret **ca-signing-secret** (tls.crt / tls.key).  
* Exposes an HTTPS service on port 8000; **carbide-ca-issuer** ClusterIssuer uses the same CA so cert-manager can issue certificates (site-manager, site-agent, Temporal client, etc.) without a separate Vault.

**Current layout (this repo):** deploy/kustomize/base includes cert-manager/deployment.yaml (single container, native PKI), cert-manager/service.yaml, cert-manager/rbac.yaml, and cert-manager-io/cluster-issuer.yaml (carbide-ca-issuer referencing ca-signing-secret). Ensure **ca-signing-secret** exists (see 7.0) before deploying. No Vault server or vault-token is required.

## **7.2 Install Temporal Certificates \- Reference only\!**

Branch with These changes \-\> [link](https://gitlab-master.nvidia.com/nvmetal/carbide-external/-/tree/sa-enablement/apps/cloud-certificates?ref_type=heads)

This step pre‑provisions the TLS certs that Temporal and the cloud workloads use for mTLS:

* **Client certs** used by cloud-api and cloud-workflow when calling Temporal.  
* **Server certs** used by Temporal for its cloud, site, and inter‑service endpoints.

All of these are issued by the carbide-ca-issuer ClusterIssuer you created in **7.1 cloud‑cert‑manager** and ultimately chain back to the CA secret from **7.0** (ca-signing-secret).

This section is **internal only**: CNs/SANs are hard‑coded to \*.temporal.nvidia.com. Customers should treat this as a pattern and replace hostnames and namespaces with their own.  

### **7.2.1 Kustomize layout**

From the internal repo location you already referenced for SAs, there is a Kustomize base with:

| client-cloud-api.yaml         \# client cert for cloud-apiclient-cloud-workflow.yaml    \# client cert for cloud-workflowserver-cloud.yaml             \# server cert for "cloud" Temporal endpointserver-interservice.yaml      \# server cert for Temporal inter-service trafficserver-site.yaml              \# server cert for "site" Temporal endpointkustomization.yaml            \# references all of the above |
| :---- |

kustomization.yaml simply lists these five Certificate resources as resources

### **7.2.2 What the certs do**

**Client certificates**

* **client-cloud-api.yaml**  
  * Namespace: cloud-api  
  * Resource: Certificate temporal-client-cloud-certs  
  * Produces Secret: temporal-client-cloud-certs  
  * CN / DNS:  
    * commonName: cloud.client.temporal.nvidia.com  
    * dnsNames: \[cloud.client.temporal.nvidia.com\]  
  * Intended use: mTLS client cert for the **cloud‑api** deployment to call Temporal.  
* **client-cloud-workflow.yaml**  
  * Namespace: cloud-workflow  
  * Resource: Certificate temporal-client-cloud-certs  
  * Produces Secret: temporal-client-cloud-certs  
  * CN / DNS identical to the cloud‑api client cert.  
  * Intended use: mTLS client cert for **cloud-worker** and **site-worker** to call Temporal.

Both client certs:

* issuerRef: name: carbide-ca-issuer, kind: ClusterIssuer, group: cert-manager.io  
* usages: server auth, client auth  
* Duration: 2160h (90 days), renewBefore: 360h (15 days)

Cloud‑side deployments mount these Secrets at:

* /var/secrets/temporal/certs/client-cloud/ in cloud-api  
* /var/secrets/temporal/certs/client-cloud/ in cloud-workflow


**Server certificates (Temporal namespace)**

All three server certs live in **namespace temporal** and are issued by carbide-ca-issuer with usages: \[server auth, client auth\].

* **server-cloud.yaml**  
  * Certificate server-cloud-certs  
  * Secret: server-cloud-certs  
  * CN / DNS:  
    * cloud.server.temporal.nvidia.com  
  * Duration: 2h, renewBefore: 1h  
     (short lifetime for internal/dev testing; adjust if you need something longer‑lived.)  
* **server-interservice.yaml**  
  * Certificate server-interservice-certs  
  * Secret: server-interservice-certs  
  * CN / DNS:  
    * interservice.server.temporal.nvidia.com  
  * Duration: 2160h / renewBefore: 360h  
* **server-site.yaml**  
  * Certificate server-site-certs  
  * Secret: server-site-certs  
  * CN / DNS:  
    * site.server.temporal.nvidia.com  
  * Duration: 2160h / renewBefore: 360h

These Secrets are consumed by your Temporal deployment (frontends, gateways, etc.) according to your internal Temporal Helm values / manifests.

### **7.2.3 Applying the cloud certs**

**Prerequisites:**

* 7.0 root CA secrets (ca-signing-secret) exist in cert-manager.  
* 7.1 cloud-cert-manager / credsmgr and the carbide-ca-issuer ClusterIssuer are healthy.

From the directory containing the kustomization.yaml and the five cert YAMLs above, apply:

| kubectl apply \-k \<path-to-cloud-certs-kustomize\> |
| :---- |

This will create:

* temporal-client-cloud-certs in namespaces:  
  * cloud-api  
  * cloud-workflow  
* server-cloud-certs, server-interservice-certs, server-site-certs in namespace temporal.

You can verify issuance with:

| kubectl \-n cloud-api        get certificate,secret temporal-client-cloud-certskubectl \-n cloud-workflow   get certificate,secret temporal-client-cloud-certskubectl \-n temporal         get certificate,secret server-\*-certs |
| :---- |

## **7.3 Install Temporal (reference only)**

Branch with These changes \-\> [link](https://gitlab-master.nvidia.com/nvmetal/temporal-helm-charts/-/merge_requests/35)  
**Audience note:**  
This section is **for internal SAs only**. Customers are expected to deploy Temporal themselves (self‑hosted or managed/hosted), following Temporal’s public documentation. They do **not** need to copy our values file; they just need to satisfy the integration constraints and the certs from **7.2**.

### **7.3.1 What we ran internally**

We used the internal temporal-helm-charts/temporal chart with a custom values.yaml that:

* Deploys Temporal **server 1.22.6** (temporalio/server:1.22.6).  
* Enables:  
  * admintools (temporalio/admin-tools:1.22.4)  
  * web (temporalio/ui:2.16.2)  
  * Embedded **Elasticsearch 7.17.3** in the same namespace for visibility.  
* Uses **PostgreSQL** for both persistence and visibility:  
  * Main DB: temporal as user temporal\_user  
    * Host: forge-pg-cluster.postgres.svc.cluster.local:5432  
    * Credentials from Secret temporal-user.forge-pg-cluster.credentials (key password)  
  * Visibility DB: temporal\_visibility as user temporal\_visibility\_user  
    * Credentials from Secret temporal-visibility-user.forge-pg-cluster.credentials  
* Enables **TLS everywhere** using the certs created in **7.2**:  
  * Mounts these secrets:  
    * server-interservice-certs  
    * server-cloud-certs  
    * server-site-certs  
  * Configures global.tls:  
    * Internode server/client use server-interservice-certs  
    * Frontend hostOverrides:  
      * interservice.server.temporal.nvidia.com → server-interservice-certs  
      * cloud.server.temporal.nvidia.com       → server-cloud-certs  
      * site.server.temporal.nvidia.com        → server-site-certs  
  * CLI and UI (admintools \+ web) set:  
    * \*\_TLS\_SERVER\_NAME=interservice.server.temporal.nvidia.com  
    * \*\_TLS\_CERT/TLS\_KEY/TLS\_CA pointing at /var/secrets/temporal/certs/server-interservice/\*  
* Exposes the **frontend** on a ClusterIP service on port **7233** in namespace temporal (internal only).

We installed it with a simple Helm command from the chart directory, e.g.:

| helm dependency update ./ helm upgrade \--install temporal . \\  \--namespace temporal \\  \-f ./values.yaml |
| :---- |

After a successful install in our reference cluster, kubectl \-n temporal get pods showed:

| elasticsearch-master-0temporal-admintools-...temporal-frontend-...temporal-history-...temporal-matching-...temporal-web-...temporal-worker-... |
| :---- |

All in Running state.

* In the customer environment, they can choose any hostnames / ports, as long as:

  * cloud-api and cloud-workflow configs point at the correct host:port.

  * The TLS settings and CA chain match the secrets they configure (analogous to the certs in **7.2**).

### **7.3.3 Ensure base Temporal namespaces (cloud and site)**

Once the Temporal cluster is up and all pods are Running, two **Temporal namespaces** must exist for Carbide:

* cloud  
* site

Customers can create these namespaces however they like (e.g. tctl namespace register, the Temporal Web UI, or automation).

For **our SAs**, the temporal-helm-charts repo includes a helper script that does this idempotently:

* **Location (repo root):** ./register-temporal-namespaces.sh  
  **Optional env:** set KUBECTL\_CONTEXT if you need a non‑default kube context.

What the script does at a high level:

1. Waits for the **admintools** pod in namespace temporal to be Ready.  
2. Waits for at least one **worker** pod to be Ready.  
3. For each namespace in cloud and site:  
   * Runs tctl \--ns \<NS\> namespace describe inside the admintools pod.  
   * If it doesn’t exist, runs tctl \--ns \<NS\> namespace register.  
   * Retries up to 10 times with a 120‑second backoff.

Example invocation and (trimmed) output:

| ./register-temporal-namespaces.sh⏳  Waiting for Temporal admintools pod to be Ready...✅  admintools pod temporal-admintools-... is Ready⏳  Waiting for at least one Temporal worker pod to be Ready...✅  worker pod temporal-worker-... is ReadyCreating Temporal namespace 'cloud' (attempt 1/10)...✓  Namespace 'cloud' created.Creating Temporal namespace 'site' (attempt 1/10)...✓  Namespace 'site' created.🎉 Temporal namespaces 'cloud' and 'site' ensured. |
| :---- |

For customers, the equivalent manual commands (run wherever tctl can reach their Temporal frontend) are:

| tctl \--ns cloud namespace registertctl \--ns site  namespace register |
| :---- |

## 

## **7.4 cloud-db (namespace cloud-db)**

Branch with These changes \-\> [link](https://gitlab-master.nvidia.com/nvmetal/cloud-db/-/merge_requests/524)

#### Role

cloud-db is a **migration job** that initializes or upgrades the cloud database schema. It’s **not** a long-running service: it runs once with a given version, applies migrations to Postgres, and exits.

#### Reference behavior we run internally

* **Namespace:** cloud-db  
* **Workload type:** one-off Job (e.g., cloud-db-migration-0-1-45)  
* **Image / version (reference only):**  
   nvcr.io/nvidian/nvforge-devel/cloud-db:v0.1.45  
   In customer environments, this must be replaced with an image they build from the GitHub repo and push to their own registry.  
* **Config source:** a ConfigMap with:  
  * dbHost → hostname/FQDN of the Postgres service  
  * dbPort → Postgres port (typically 5432\)  
  * dbName → name of the application database (e.g., forge or cloud)  
  * dbUser → logical Postgres user (e.g., forge)  
* **Credentials source:** a single DB credential Secret (in our reference: forge.forge-pg-cluster.credentials) with keys:  
  * username  
  * password  
* **Runtime behavior:**  
  * On start, the Job logs:  
    * PGUSER, DBNAME and a masked DB URL: postgres://PGUSER:\*\*\*\*@PGHOST:PGPORT/DBNAME  
  * If PGUSER, PGPASSWORD, or DBNAME are missing, the Job fails immediately.  
  * Otherwise, it executes a migration entrypoint (e.g., /usr/local/bin/migrate.sh).

### **What a solution architect must adjust for a customer environment**

When you hand this to a customer, they need to **review and customize** the following fields:

1. **Postgres host & port**  
   * **dbHost** must be set to **their Postgres service FQDN or hostname**, most likely not our internal forge-pg-cluster.postgres.svc.cluster.local.  
   * **dbPort** must match their DB port (usually still 5432).  
2. **Database name & user**  
   * **dbName** should be the name of the database where cloud schemas live (e.g., cloud, forge, or whatever they use).  
   * **dbUser** should match the Postgres user they’ve provisioned for cloud components.  
3. **Database credential Secret**  
   * The Job expects to read PGUSER and PGPASSWORD from a K8s Secret. In our reference, that Secret is forge.forge-pg-cluster.credentials with keys username and password.  
   * Customers can either:  
     * Reuse that Secret name (and ESO mapping) and just make sure the values are correct for their DB, **or**  
     * Change the Job’s environment to point at **their own Secret name** (e.g., cloud-db-credentials) while preserving the key names (username / password).  
4. **Image and registry**  
   * The **image name and tag** must be updated to point at the **customer’s registry**, not nvcr.io.  
   * For example:  
     * Build from the repo: docker build \-t \<REGISTRY\>/\<PROJECT\>/cloud-db:\<TAG\> .  
     * Push it: docker push \<REGISTRY\>/\<PROJECT\>/cloud-db:\<TAG\>  
     * Update the deployment to use image: \<REGISTRY\>/\<PROJECT\>/cloud-db:\<TAG\>.  
5. **Image pull details**  
   * If the registry is private, the imagePullSecrets entry must match the **customer’s Docker registry secret name**.  
   * If they use a public registry or anonymous access, this can be dropped.  
6. **Migration job naming & versioning**  
   * The Job name (cloud-db-migration-0-1-17 in our example) encodes a migration version.  
   * For new releases, **architects should bump both**:  
     * The Job name (e.g., cloud-db-migration-0-1-18) and  
     * The image tag (e.g., v0.1.18),  
        so operators can see which migration ran.

### **How to apply this component**

First if the following extensions are not setup:

| \-c 'CREATE EXTENSION IF NOT EXISTS btree\_gin;' \\\-c 'CREATE EXTENSION IF NOT EXISTS pg\_trgm;' |
| :---- |

Run [create-postgres-extensions.sh](http://create-postgres-extensions.sh) while pointing to the cluster:

| \~/go/src/cloud-db/deploy🌴git:(sa-enablement-distroless) ⌚ 21:40:53 » ./create-postgres-extensions.sh📦 Using database name: forge⏳  Waiting for StatefulSet postgres/forge-pg-cluster replicas to be Ready...✅  Postgres StatefulSet is Ready (3/3)🔑  Running extension SQL inside pod forge-pg-cluster-0Defaulted container "postgres" out of: postgres, postgres-exporterCREATE EXTENSIONCREATE EXTENSION✅  Postgres extensions ensured. |
| :---- |

Once the customer-specific fields are updated in the cloud-api Kustomize base (image, registry auth, config, DB secret name, service type), the architect applies it with:

| kubectl apply \-k deploy/kustomize/base  |
| :---- |

## **7.5 cloud‑workflow (namespace cloud-workflow)**

Branch with These changes \-\> [link](https://gitlab-master.nvidia.com/nvmetal/cloud-workflow/-/tree/sa-enablement/deploy/kustomize?ref_type=heads)

### **Role**

cloud-workflow provides the **cloud-side Temporal workers**. It has two deployments:

* cloud-worker → processes workflows in the **cloud** Temporal namespace & queue.  
* site-worker → processes workflows in the **site** Temporal namespace & queue.

Both read a shared config file and differ only in TEMPORAL\_NAMESPACE / TEMPORAL\_QUEUE.

### **7.5.1 Dependencies**

Before cloud-workflow is applied, the following must already be in place:

1. **PostgreSQL \+ schema**  
   * Postgres endpoint reachable (same DB that cloud-db migrates).  
   * DB created (e.g., forge or cloud).  
   * cloud-db migration Job has succeeded.  
   * Credential Secret (in our reference: forge.forge-pg-cluster.credentials) with username / password.  
2. **Temporal**  
   * Frontend gRPC endpoint reachable from this namespace (reference: temporal-frontend-headless.temporal.svc.cluster.local:7233).  
   * Temporal namespaces cloud and site are registered.  
3. **Temporal client TLS Secret**  
   * Secret temporal-client-cloud-certs containing tls.crt, tls.key, ca.crt.  
   * Mounted at /var/secrets/temporal/certs/client-cloud/.  
4. **Image registry**  
   * A registry controlled by the customer where they will **build and push** the cloud-workflow image.

### **7.5.2 Behavior in our reference environment**

* **Namespace:** cloud-workflow  
* **Workload type:**  
  * Deployment cloud-worker (3 replicas).  
  * Deployment site-worker (3 replicas).  
* **Container port:** 8899 (HTTP health/ready; metrics served on 9360).  
* **Liveness / Readiness probes:**  
  * GET /healthz on 8899  
  * GET /readyz on 8899  
* **Config file:** /etc/cloud-workflow/config.yaml  
  * CONFIG\_FILE\_PATH env is set to this path.  
* **Volumes (both deployments):**  
  * forge-pg-cluster-auth → Secret forge.forge-pg-cluster.credentials → /var/secrets/db/auth/  
  * cloud-workflow-config → ConfigMap cloud-workflow-config → /etc/cloud-workflow/  
  * temporal-client-cloud-certs → Secret temporal-client-cloud-certs → /var/secrets/temporal/certs/client-cloud/  
* **Temporal env per deployment:**  
  * cloud-worker: TEMPORAL\_NAMESPACE=cloud, TEMPORAL\_QUEUE=cloud  
  * site-worker:  TEMPORAL\_NAMESPACE=site,  TEMPORAL\_QUEUE=site

#### Internal image we validated with (for reference only):

| nvcr.io/nvidian/nvforge-devel/cloud-workflow:v0.2.33 |
| :---- |

This NVCR image is **only** what we ran internally. Customers **must build their own image from source and push it to their registry**, then update the deployment to use that registry URL.

#### 7.5.3 Kustomize layout (sa-enablement branch)

In the cloud-workflow repo’s sa-enablement branch, we expect:

| deploy/kustomize/base/cloud-workflow/  secrets/    temporal-client-cloud-certs.yaml  \# placeholder TLS secret, empty values  configmap.yaml                      \# cloud-workflow-config (config.yaml)  deployment.yaml                     \# cloud-worker \+ site-worker  service.yaml                        \# Services for each worker deployment  imagepullsecret.yaml  namespace.yaml  kustomization.yaml |
| :---- |

#### 7.5.4 Config (what SAs must change)

**File:** deploy/kustomize/base/cloud-workflow/configmap.yaml

**Key:** data key config.yaml inside the cloud-workflow-config ConfigMap.

Fields that are **customer-specific**:

* log.level → e.g., debug, info, warn.  
* db.\*:  
  * db.host → customer’s Postgres host/FQDN.  
  * db.port → customer’s port (usually 5432).  
  * db.name → name of customer’s DB.  
  * db.user → DB user name.  
  * db.passwordPath → path under /var/secrets/db/auth/ that the app expects.

* temporal.\*:  
  * temporal.host → Temporal frontend DNS name (e.g., temporal-frontend-headless.temporal.svc.cluster.local or the customer’s host).  
  * temporal.port → typically 7233\.  
  * temporal.serverName → server name used for TLS validation (if applicable).  
  * temporal.namespace → should match the worker’s TEMPORAL\_NAMESPACE env.  
  * temporal.queue → should match the worker’s TEMPORAL\_QUEUE env.  
  * temporal.tls.certPath, keyPath, caPath → must align with where temporal-client-cloud-certs is mounted.  
* siteManager.svcEndpoint → full URL for the site API (e.g., https://sitemgr.cloud-site-manager:8100/v1/site).

#### 7.5.5 Secrets (placeholders vs real values)

**File:** deploy/kustomize/base/cloud-workflow/secrets/temporal-client-cloud-certs.yaml

* Contains empty tls.crt, tls.key, ca.crt. This is a placeholder and can be applied as is\!

**DB credentials**

* forge.forge-pg-cluster.credentials must exist in the cluster (ESO or manual).  
* If the customer uses a different secret, they must update the volumes\[\].secret.secretName in deployment.yaml accordingly.

#### 7.5.6 Image and registry (critical change vs internal NVCR)

In the repo’s deployment.yaml, the containers currently reference:

| image: nvcr.io/nvidian/nvforge-devel/cloud-workflow:v0.2.33 |
| :---- |

For customer deployments:

1. **Build from source** (in the cloud-workflow repo):

| docker build \-t \<CUSTOM\_REGISTRY\>/\<PROJECT\>/cloud-workflow:\<TAG\> .docker push \<CUSTOM\_REGISTRY\>/\<PROJECT\>/cloud-workflow:\<TAG\> |
| :---- |

2. **Update the image reference** in the manifests to use **their registry URL**, not NVCR:  
   * Either via **Kustomize images** in deploy/kustomize/base/cloud-workflow/kustomization.yaml:

| images:  \- name: cloud-workflow          \# logical name in deployment.yaml    newName: \<CUSTOM\_REGISTRY\>/\<PROJECT\>/cloud-workflow    newTag: \<TAG\> |
| :---- |

   * Or by editing deployment.yaml directly if you are not using the images: transform.  
3. Ensure imagePullSecrets points at the customer’s registry secret when the registry is private.

#### 7.5.7 Services

The repo defines **internal Services** for each worker:

* cloud-worker Service:  
  * Port 8899 → HTTP/health  
  * Port 9360 → metrics  
* site-worker Service:  
  * Same ports as above.

These are intended for **internal use** (probes/metrics scraping), not direct external traffic.


#### 7.5.8 How to apply cloud-workflow

After:

* cloud-db has completed its migrations,  
* DB credentials exist in the expected Secret,  
* Temporal is reachable and cloud / site namespaces exist,  
* temporal-client-cloud-certs and any OTEL secrets are populated appropriately,

the solution architect can deploy cloud-workflow with:

| kubectl apply \-k deploy/kustomize |
| :---- |

This will create:

* Namespace cloud-workflow,  
* ConfigMap cloud-workflow-config,  
* Secrets in cloud-workflow/secrets (as provided/overridden),  
* Deployments cloud-worker and site-worker,  
* Services for cloud-worker and site-worker.

## **7.6 cloud‑api (namespace cloud-api)**

Branch with These changes \-\> [link](https://gitlab-master.nvidia.com/nvmetal/cloud-api/-/merge_requests/1181)

#### Role

cloud-api is the **front-end API** for cloud-side operations. It exposes HTTP endpoints and (optionally) metrics; it reads configuration from a YAML file, connects to the cloud DB, and integrates with Temporal(and grpc backend through temporal),  site-manager, and optional telemetry backends.

#### Reference behavior we run internally

* **Namespace:** cloud-api  
* **Workload type:** Deployment (1 replica per node)  
* **Image / version (reference only):**  
   nvcr.io/nvidian/nvforge-devel/cloud-api:v0.2.76  
* **Port / probes:**  
  * Container port **8388** (named api).  
  * Liveness probe: GET /healthz on 8388\.  
  * Readiness probe: GET /readyz on 8388\.  
* **Service:** LoadBalancer with:  
  * Port **80** → targetPort **8388** (http)  
  * Port **9360** → targetPort **9360** (metrics)  
* **Config file:**  
  * CONFIG\_FILE\_PATH=/etc/cloud-api/config.yaml  
  * Volume cloud-api-config mounts that config at /etc/cloud-api/config.yaml.  
* **Secrets / volumes expected by the deployment:**  
  * forge.forge-pg-cluster.credentials → mounted at /var/secrets/db/auth/  
     (DB username/password)  
  * temporal-client-cloud-certs → mounted at /var/secrets/temporal/certs/client-cloud/  
     (Temporal client TLS: tls.crt, tls.key, ca.crt)

#### Kustomize layout (sa-enablement branch)

Under the cloud-api repo’s sa-enablement [branch](https://gitlab-master.nvidia.com/nvmetal/cloud-api/-/merge_requests/1181):

| deploy/kustomize/base/  secrets/    ssa-client-secret.yaml              \# placeholder    temporal-client-cloud-certs.yaml    \# placeholder  configmap.yaml                        \# cloud-api-config (config.yaml)  Deployment.yaml   Imagepullsecret.yaml                  \# replace the secret value  service.yaml  namespace.yaml  kustomization.yaml |
| :---- |

### **Config & integration details**

#### 1\. Application config (config.yaml)

This is the **main file SAs need to customize** for each customer.

* **File:** deploy/kustomize/base/cloud-api/configmap.yaml  
* **Key:** cloud-api-config.data\[\`config.yaml\`\]  
* **Important fields inside config.yaml:**  
  * log.level → log level (e.g., debug, info).  
  * kas.legacyJwksUrl, kas.ssaJwksUrl, kas.starfleetJwksUrl → JWKS endpoints for various IdPs / SSA integration (customer URLs).  
  * db.host, db.port, db.name, db.user, db.passwordPath → DB connection:  
    * host and port must match **customer’s Postgres service** (often the same as cloud-db-config values).  
    * name and user must correspond to their DB and role.  
    * passwordPath should match the path under /var/secrets/db/auth/ where the password is mounted.  
  * temporal.host, temporal.port, temporal.serverName, temporal.namespace, temporal.queue:  
    * host must be the Temporal frontend or headless service for Temporal (e.g., temporal-frontend-headless.temporal.svc.cluster.local).  
    * namespace/queue must match the **Temporal namespace / task queue** used by cloud workers (e.g., **cloud**).  
    * tls.certPath, keyPath, caPath must match how temporal-client-cloud-certs is mounted.  
  * siteManager.enabled, siteManager.svcEndpoint → URL for site-manager’s v1/site endpoint.

#### Customer actions (config.yaml):

* Update all **URLs**, **hostnames**, and **ports** to match the customer’s auth, Temporal, DB, site-manager, and ZincSearch endpoints.  
* Update db.\* to align with their Postgres configuration.  
* Confirm TLS paths line up with the volume mount paths in deployment.yaml.

#### 2\. Deployment & container image

* **File:** deploy/kustomize/base/cloud-api/**deployment.yaml**  
* **Key elements:**  
  * spec.replicas: default **2**; tune per customer.  
  * Container:  
    * Name: cloud-api  
    * Image (reference): nvcr.io/nvidian/nvforge-devel/cloud-api:v0.2.76  
    * Env: CONFIG\_FILE\_PATH=/etc/cloud-api/config.yaml  
    * Ports: 8388 named api  
  * Volume mounts:  
    * /var/secrets/db/auth/ → forge.forge-pg-cluster.credentials  
    * /var/secrets/ssa/ → ssa-client-secret  
    * /etc/cloud-api/ → cloud-api-config  
    * /var/secrets/temporal/certs/client-cloud/ → temporal-client-cloud-certs

**Customer actions (deployment):**

* Replace the image with their own build:  
  * Build: docker build \-t \<REGISTRY\>/\<PROJECT\>/cloud-api:\<TAG\> .  
  * Push:  docker push \<REGISTRY\>/\<PROJECT\>/cloud-api:\<TAG\>  
* Update the Kustomize images mapping:  
  * **File:** deploy/kustomize/base/cloud-api/kustomization.yaml  
     images: → set newName: \<REGISTRY\>/\<PROJECT\>/cloud-api, newTag: \<TAG\>.  
* Decide if they want imagePullSecrets:  
  * For private registries, ensure:

| imagePullSecrets:  \- name: \<CUSTOM\_IMAGEPULLSECRET\> |
| :---- |

#### 3\. Database credentials Secret

* The deployment mounts a secret called forge.forge-pg-cluster.credentials at /var/secrets/db/auth/.  
* That secret is expected to have **username** and **password** keys.

**Customer actions (DB Secret):**

* Either:  
  * Keep the name forge.forge-pg-cluster.credentials and create/populate it (ideally via ESO), **or**  
  * Change volumes\[\].secret.secretName \+ any env.valueFrom.secretKeyRef.name in deployment.yaml to their Secret name (e.g., cloud-db-credentials).

**File:** deploy/kustomize/base/cloud-api/deployment.yaml

#### 4\. SSA client secret (ssa-client-secret)

There are two related files:

* Placeholder Secret with empty value (original pattern)  
   **File:** deploy/kustomize/base/cloud-api/secrets/ssa-client-secret.yaml  
  * client-secret: "" (base64 of empty) — **safe for non-SSA environments**.  
* Secret with placeholder value c3NhLXJlcGxhY2U= (base64 for "ssa-replace")  
   **File:** deploy/kustomize/base/cloud-api/ssa-placeholder-secret.yaml  
  * This is **only a placeholder**; customers must override it with:  
    * A real SSA client secret (if using SSA), or  
    * An empty value in non-SSA deployments.

**Customer actions (SSA Secret):**

* If SSA is in scope:  
  * Set client-secret to the **base64-encoded** SSA client secret.  
* If SSA is **not** used:  
  * Keep the placeholder

#### 5\. Temporal client TLS certs (temporal-client-cloud-certs)

* Placeholder TLS Secret with keys tls.crt, tls.key, ca.crt set to "" (empty).  
   **File:** deploy/kustomize/base/cloud-api/secrets/temporal-client-cloud-certs.yaml  
* Mounted at /var/secrets/temporal/certs/client-cloud/.

**Customer actions (Temporal TLS):**

* If Temporal requires mTLS (our reference assumes yes):  
  * Fill in **base64-encoded** tls.crt, tls.key, and ca.crt for the client cert, key, and CA.  
* If they run Temporal without TLS:  
  * They can leave the cert values blank but must ensure config.yaml TLS paths and Temporal deployment tolerate non-TLS (or remove TLS config from config.yaml).

#### 6\. Service exposure

* **File:** deploy/kustomize/base/cloud-api/service.yaml  
* **Behavior in reference:**  
  * Type LoadBalancer.  
  * Port **80/TCP** → targetPort: 8388 (HTTP API).  
  * Port **9360/TCP** → targetPort: 9360 (metrics).  
* Selector:  
  * app: cloud-api (matches labels in deployment.yaml).

**Customer actions (Service):**

* Decide **how** cloud-api should be exposed:  
  * **Internal only**: set type to ClusterIP and front with Ingress/HTTPProxy.  
  * **External**: keep or adjust LoadBalancer and connect via their LB/Ingress.  
* Adjust ports if they need different external port numbers, but keep targetPort aligned with the container’s 8388 for API and 9360 for metrics.


**How to apply cloud-api**

Once all customer-specific knobs are updated:

* configmap.yaml (app config),  
* deployment.yaml (image, registry, secrets, volumes),  
* service.yaml (type/ports),  
* secrets in secrets/ and ssa-placeholder-secret.yaml (as needed),  
* and otel.env (if telemetry is used),

the architect can deploy cloud-api by running:

| kubectl apply \-k deploy/kustomize/base/cloud-api |
| :---- |

https://docs.google.com/document/d/174MUDmsOAtZiE7qrYQ7D-cJFQR40B83xeBHJFT\_ivkM/edit?tab=t.0

## **7.7 cloud‑site‑manager (namespace cloud-site-manager)**

#### Role

cloud-site-manager (a.k.a. *sitemgr*) is “site registry” service. It:

* Owns the Site CRD and stores per‑site metadata as Site objects.  
* Talks to **credsmgr** to mint per‑site credentials / CA material.  
* Exposes an HTTPS API on port **8100** that the bootstrap flow calls to:  
  * create a new site (/v1/site), and  
  * later allow the site agent to fetch its One TIme Password(OTP) credentials.


### **7.7.1 Dependencies**

Before you deploy cloud-site-manager, you should already have:

* **cert-manager \+ cloud-cert-manager (credsmgr)**:  
  * carbide-ca-issuer ClusterIssuer configured and working.  
  * credsmgr Service reachable at https://credsmgr.cert-manager:8000 (or your equivalent).  
* **CA secret**: ca-signing-secret in the deployment namespace (see 7.0).  
* **Kubernetes CRD support** (standard; no extra setup).

There is **no direct Postgres or Temporal dependency** in this manifest; those get exercised later by other components through the site bootstrap workflows.

### **7.7.2 Behavior in our reference environment**

**Namespace**: cloud-site-manager

**Workloads**:

* **CRD**: sites.forge.nvidia.io  
  * Kind: Site  
  * Scope: Namespaced  
  * spec and status are both “free‑form” (x-kubernetes-preserve-unknown-fields: true), so the CRD doesn’t constrain your schema.  
* **Deployment**: csm  
  * Replicas: **3**  
  * Container: sitemgr  
  * Image (reference only):  
     nvcr.io/nvidian/nvforge-devel/cloud-site-manager:v0.1.16  
  * Arguments:  
    * \--listen-port=8100  
    * \--creds-manager-url=https://credsmgr.cert-manager:8000  
    * \--ingress-host=sitemgr.cloud-site-manager.svc.cluster.local  
  * Environment:  
    * SITE\_MANAGER\_NS from metadata.namespace  
    * envFrom: secretRef: otel-lightstep (optional telemetry; can be stubbed or removed)  
  * Probes:  
    * Liveness / Readiness / Startup: HTTPS GET /healthz on port 8100  
  * ServiceAccount: site-manager  
* **RBAC** (namespace‑scoped):  
  * Role site-manager grants full access to sites and sites/status in API group forge.nvidia.io.  
  * RoleBinding site-manager binds the Role to the site-manager ServiceAccount in cloud-site-manager.  
* **Service**: sitemgr  
  * Type: ClusterIP  
  * Port 8100 → targetPort 8100  
  * Selector: app.kubernetes.io/name: csm  
  * This is what your bootstrap flow curls:  
     https://sitemgr.cloud-site-manager:8100/v1/site


### **7.7.3 Kustomize layout (sa‑enablement branch)**

In the cloud-site-manager repo (sa‑enablement branch), the base Kustomize tree for this component looks like:

* deploy/kustomize/  
  * namespace.yaml              – defines cloud-site-manager namespace  
  * crd-site.yaml               – Site CRD (sites.forge.nvidia.io)  
  * deployment.yaml             – csm Deployment (sitemgr container)  
  * serviceaccount.yaml         – site-manager ServiceAccount  
  * serviceaccount-role.yaml    – Role for CRD access  
  * serviceaccount-rolebinding.yaml – binds Role to site-manager SA  
  * sitemgr-service.yaml        – ClusterIP Service on port 8100  
  * kustomization.yaml          – includes the above and overrides the image

kustomization.yaml sets the **reference** image to:

| images:  \- name: sitemgr    newName: nvcr.io/nvidian/nvforge-devel/cloud-site-manager    newTag: v0.1.16 |
| :---- |

### 

### **7.7.4 Customer knobs (what SAs must adjust)**

1. **Image & registry**

Customers must build cloud-site-manager from source and push to their own registry, then update the image mapping in kustomization.yaml.

Example:

| \# Build and pushdocker build \-t \<REGISTRY\>/\<PROJECT\>/cloud-site-manager:\<TAG\> .docker push \<REGISTRY\>/\<PROJECT\>/cloud-site-manager:\<TAG\> |
| :---- |

Then in the Kustomize base:

| images:  \- name: sitemgr    newName: \<REGISTRY\>/\<PROJECT\>/cloud-site-manager    newTag: \<TAG\> |
| :---- |

If their registry is private, ensure an appropriate imagePullSecrets entry (and corresponding Secret) is present in the Deployment; otherwise it can be removed.

2. **credsmgr URL**

The Deployment passes:

| \--creds-manager-url=https://credsmgr.cert-manager:8000 |
| :---- |

If the customer:

* Renames the credsmgr Service,  
* Runs it in a different namespace,

they **must update this flag** to point at their actual URL. Whatever CA service this URL references must:

* Trust the same root CA that site agents will trust, and  
* Be able to issue the per‑site certs/OTP material your bootstrap flow expects.


3. **Ingress / internal host**

The \--ingress-host flag is currently:

| \--ingress-host=sitemgr.cloud-site-manager.svc.cluster.local |
| :---- |

This is the hostname sitemgr uses when constructing URLs and, in some flows, may be placed into certificates or responses.

Customer options:

* Keep it as the internal DNS name if the service will only ever be called in‑cluster.  
* Change it to an internal or external FQDN (e.g. sitemgr.\<customer-domain\>), but then:  
  * Make sure their DNS / ingress routes that hostname to the sitemgr Service.  
  * Ensure any TLS certs used for sitemgr include this hostname in their SANs.

### **7.7.5 How to apply cloud-site-manager**

After:

* cloud-cert-manager / credsmgr is running and reachable at the URL you specified,  
* The CA secret (ca-signing-secret) exists in the deployment namespace (e.g. carbide),  
* The image reference has been updated to the customer’s registry,  
* Any telemetry Secret decisions are made,

the SA can deploy cloud-site-manager with:

| $ kubectl apply \-k deploy/kustomize namespace/cloud-site-manager created customresourcedefinition.apiextensions.k8s.io/sites.forge.nvidia.io created serviceaccount/site-manager created role.rbac.authorization.k8s.io/site-manager created rolebinding.rbac.authorization.k8s.io/site-manager created service/sitemgr created deployment.apps/csm created  |
| :---- |

This creates:

* Namespace cloud-site-manager,  
* The Site CRD in the cluster,  
* ServiceAccount \+ Role \+ RoleBinding for site-manager,  
* The csm Deployment (sitemgr) and sitemgr Service.

## **7.13 SA Enablement \- Carbide bring up**

Requirements:

* All common services running  
* CA secret (ca-signing-secret) and cert-manager ready

### **Terraform**

We need to add proper policies to the cluster using terraform (e.g. for PKI or network policies). This work may be abstracted in your environment’s automation (e.g. carbide-external repository).

* Please pull the latest main in carbide-external repository.  
* Make sure you have access to the cluster  
  * Run ssh tunnel command in the background o another terminal  
    ssh \-D 8888 \-N renojump 

* Switch back to repo and make sure you kubectx is pointed to the correct cluster  
  export KUBECONFIG=/PATH/TO/SITE/kubeconfig 

* Terraform might complain about kubeconfig certificate since it looks at   
  \~/.kube/config   
  Make sure to to backup your default config. Then copy over  
  cp /PATH/TO/SITE/kubeconfig  \~/.kube/config 

* Now you can run the terraform script  
  ./manifests/setup\_terraform.sh dev6 

### **Image Pull Secret**

* Since we are still using internal NVIDIA images we need NVCR auth token to pull these images. Export in an env variable  
  export NVCR\_AUTH\_TOKEN=XXX   
* Then create the secret:

| kubectl create secret generic imagepullsecret \\      \--type=kubernetes.io/dockerconfigjson \\      \-n default \\      \--from-literal=.dockerconfigjson="{\\"auths\\":{\\"nvcr.io\\":{\\"auth\\":\\"$NVCR\_AUTH\_TOKEN\\"}}}" |
| :---- |


* 

There is already a branch in the carbide repository called sa-enablement\_dev8.  
* Talks to **Carbide** over mTLS gRPC to orchestrate hardware operations.  
* Reads/writes to a **site‑local PostgreSQL** database for persistent state.  
* Uses **cert‑manager (carbide-ca-issuer)** to obtain a client cert and trust the CA.  
* Is bootstrapped **per site** using an OTP \+ CA bundle \+ credentials delivered by the cloud side (cloud‑site‑manager).

Workload shape:

* StatefulSet elektra-site-agent  
  * replicas: 3  
  * serviceAccount: site-agent  
* Services:  
  * elektra-headless (headless, metrics on 2112\)  
  * elektra (ClusterIP, metrics on 2112\)  
* Metrics: port **2112** (scraped by Prometheus via ServiceMonitor if present).


### 

### 

### **7.14.2 What must exist before you deploy Elektra**

Before applying the Elektra Kustomize overlays, your cluster must already have:

#### PostgreSQL

* A database reachable from elektra-site-agent namespace with:  
  * DB\_HOST, DB\_PORT, DB\_NAME encoded in ConfigMap elektra-database-config.  
* A Secret referenced as **elektra-site-agent.elektra.forge-pg-cluster.credentials** (or your renamed equivalent) with keys:  
  * username  
  * password  
     Typically from your secret backend (e.g. ESO) or created manually.

#### Temporal

* A reachable frontend endpoint, e.g.:  
  * temporal-frontend-headless.temporal.svc.cluster.local:7233 (in‑cluster), or  
  * site.server.temporal.\<your-domain\>:7233 (if you front Temporal with a gateway).  
* Namespaces:  
  * cloud  
  * site  
  * A **per‑site namespace** that will equal cluster\_id / SITE\_UUID (created later by create-site-in-db.sh).

If enable\_tls=true in Elektra, you also need a way to obtain an mTLS client certificate and CA bundle for the site (wired into the temporal-cert Secret later).

#### cert‑manager (carbide-ca-issuer)

* A ClusterIssuer for site‑agent mTLS, e.g. **carbide-ca-issuer**.  
* A CertificateRequestPolicy that allows the Elektra namespace to request a client cert with a SPIFFE URI like:  
  * spiffe://forge.local/elektra-site-agent/sa/\<SA\_NAME\>

The grpc‑client overlay already defines:

* CertificateRequestPolicy site-agent-approver-policy  
* ClusterRole/ClusterRoleBinding giving cert-manager permission to **use** that policy  
* A Certificate grpc-client-cert issuing Secret elektra-site-agent-grpc-client-cert via carbide-ca-issuer.

You just need to ensure your **ClusterIssuer name / SPIFFE root** match your PKI.


### 

### 

### **7.14.3 Configuration knobs (ConfigMaps → env)**

The **base** config lives in base/files/config.properties and is turned into ConfigMap elektra-config-map. That, plus elektra-database-config, feeds almost all Elektra environment variables via the StatefulSet.

Key knobs SAs/customers must understand:

#### Temporal & TLS

* temporal\_server  
   Hostname of the site Temporal gateway Elektra connects to, e.g.:  
  * site.server.temporal.nvidia.com (external name), or  
  * temporal-frontend-headless.temporal.svc.cluster.local (in‑cluster).  
* temporal\_cert\_path  
   Filesystem path where the Temporal mTLS secret is mounted, typically:  
  * /var/secrets/temporal/certs  
* Must match the temporal-auth volume mount in the StatefulSet (which projects the temporal-cert Secret).  
* enable\_tls (true/false)  
   Whether Elektra uses mTLS for Temporal.  
* temporal\_inventory\_schedule  
   Cron‑style string for periodic inventory workflows (e.g. @every 3m).  
* temporal\_subscribe\_queue, temporal\_publish\_queue  
   Task queues used for inbound/outbound workflows (cloud‑workflow must match these).

#### Site identity & Temporal namespace

* cluster\_id  
   The **site UUID**. Used as:  
  * CLUSTER\_ID env  
  * TEMPORAL\_SUBSCRIBE\_NAMESPACE (via cluster\_id)  
* It must match:  
  * The SITE\_UUID used in site.sql and in DB rows,  
  * The Temporal namespace created for the site, and  
  * The Site CR created by cloud‑site‑manager.

In the overlay, you will explicitly set cluster\_id=\<SITE\_UUID\> and temporal\_subscribe\_namespace=\<SITE\_UUID\>; see 7.4.5/7.4.6.

#### Database (elektra-database-config)

elektra-database-config is generated by the base kustomization and contains:

* DB\_HOST – host for the site DB (default: forge-pg-cluster.postgres.svc.cluster.local)  
* DB\_PORT – port (default: 5432\)  
* DB\_NAME – DB name (default: elektra)  
* SECRET\_REF – name of the DB credentials Secret in your environment (defaults to elektra-site-agent.elektra.forge-pg-cluster.credentials.postgresql.acid.zalan.do)

The StatefulSet expects a Secret named **elektra-site-agent.elektra.forge-pg-cluster.credentials** (without the .postgresql.acid… suffix) for DB\_USER and DB\_PASSWORD. Adjust the Secret name and references consistently if you rename them.


#### Operational knobs

* dev\_mode, enable\_debug, clientset\_log  
* metrics\_port (default 2112\)  
* site\_workflow\_version, cloud\_workflow\_version – semantic tags that must match the versions of your site/cloud workflows in Temporal.  
* upgrade\_cron\_schedule, upgrade\_frequency\_week\_nums, upgrade\_batches – govern when/how upgrade workflows run.

If you integrate with Lightstep/OTel, set these and ensure an appropriate otel-lightstep-access Secret exists. Otherwise, leave them unset and/or disable the OTEL integration paths in Elektra.


### **7.14.4 Certificates & gRPC client identity**

#### gRPC client cert (SPIFFE)

The grpc‑client overlay defines:

* CertificateRequestPolicy site-agent-approver-policy  
  * issuerRef → ClusterIssuer carbide-ca-issuer  
  * Scoped to namespace: elektra-site-agent  
  * Allows DNS SANs:  
    * \*.svc, \*.cluster.local, \*.svc.cluster.local  
  * Allows URIs with prefix:  
    * spiffe://forge.local/elektra-site-agent/sa/\*  
* Certificate grpc-client-cert  
  * secretName: elektra-site-agent-grpc-client-cert  
  * dnsNames: elektra.elektra-site-agent.svc.cluster.local  
  * uris: spiffe://forge.local/elektra-site-agent/sa/elektra-site-agent

The StatefulSet mounts this Secret as:

* /etc/carbide (volume spiffe).

#### Customer knobs

If your PKI differs, update:

* selector.issuerRef.name in the CertificateRequestPolicy,  
* spec.issuerRef.name, and SAN/URI values in grpc-client-cert,  
* The secretName used by the spiffe volume in the StatefulSet.

Ensure cert‑manager (and carbide-ca-issuer) allow the **site-agent** namespace to request that cert. Elektra won’t establish mTLS correctly with Carbide/other services if this cert is missing and enable\_tls=true.

#### Dynamic secrets used by Elektra

1. **Database credentials Secret**  
   * Name: elektra-site-agent.elektra.\<your-db-cluster\>.credentials (or equivalent).  
   * Keys: username, password.  
   * Source: your secret backend (e.g. ESO) or manual Secret.  
   * Used by envs DB\_USER and DB\_PASSWORD.

2. **Site registration Secret bootstrap-info** – **runtime only**  
    Not in static manifests. Expected keys:  
   * site-uuid – must equal cluster\_id.  
   * otp – OTP returned by cloud‑site‑manager.  
   * creds-url – URL Elektra calls to fetch long‑lived site credentials.  
   * cacert – CA chain to trust that URL.  
3. Mounted at /etc/sitereg (one file per key). This is created by setup-site-bootstrap.sh (see 7.4.5).  
4. **Temporal certificate Secret temporal-cert** – **runtime**  
    Base manifests include a placeholder Secret with AutoGenerated values to satisfy tooling. Elektra expects keys:  
   * otp  
   * cacertificate  
   * certificate  
   * key  
5. Mounted via projected volume temporal-auth at /var/secrets/temporal/certs, which must match temporal\_cert\_path.

    In a real deployment, your bootstrap flow (or operator) replaces the placeholder with real mTLS material. setup-site-bootstrap.sh creates an empty scaffold (values are blank) so you can wire in real certs later.


### **7.14.5 Site bootstrap workflow (helper scripts)**

To avoid hand‑editing UUIDs, SQL, and secrets, the Elektra repo includes three helper scripts in elektra-site-agent/deploy/. SAs should run them **in order**.

All scripts default to the current kube context; you can override via KUBECTL\_CONTEXT=\<context\>.

#### 1\) gen-site-sql.sh – generate site.sql \+ SITE\_UUID

From elektra-site-agent/deploy:

| ./gen-site-sql.sh |
| :---- |

What it does:

* Generates a **SITE\_UUID** (or uses SITE\_UUID env if set).  
* Uses the directory name as **site name** (or SITE\_NAME env).  
* Writes a site.sql file (next to the script) containing two inserts:  
  * one into infrastructure\_provider,  
  * one into site with id \= SITE\_UUID.

Typical output:

| 🔖  Site name : deploy🔑  SITE\_UUID : 12bd0e80-3382-40d6-82d0-0b20dc14be11📌  First-state log: SITE\_UUID=12bd0e80-3382-40d6-82d0-0b20dc14be11✅  Generated .../deploy/site.sql |
| :---- |

#### 2\) create-site-in-db.sh – insert site rows \+ Temporal namespace

Then:

| ./create-site-in-db.sh\# or ./create-site-in-db.sh /path/to/site.sql |
| :---- |

What it does:

* Extracts SITE\_UUID from site.sql.  
* Reads the cloud DB name from ConfigMap cloud-db-config in namespace cloud-db.  
* Locates the master Postgres pod in namespace postgres.  
* If no site row for this SITE\_UUID exists, pipes site.sql into psql to insert it.  
* Uses the Temporal admintools pod to ensure a namespace named SITE\_UUID exists  
   (tctl \--ns \<SITE\_UUID\> namespace register if needed).

Typical output:

| Using SITE\_UUID=12bd0e80-3382-40d6-82d0-0b20dc14be11⏳  Checking if site rows already exist...Inserting site data into Postgres...INSERT 0 1INSERT 0 1⏳  Ensuring Temporal namespace 12bd0e80-3382-40d6-82d0-0b20dc14be11 exists...Namespace 12bd0e80-3382-40d6-82d0-0b20dc14be11 successfully registered.✓  Temporal namespace 12bd0e80-3382-40d6-82d0-0b20dc14be11 registered.🎉  Site DB rows & Temporal namespace ensured for SITE\_UUID=12bd0e80-3382-40d6-82d0-0b20dc14be11. |
| :---- |

#### 3\) setup-site-bootstrap.sh – create site CR \+ secrets

Finally:

| ./setup-site-bootstrap.sh\# Optional: SITE\_UUID=\<existing\> ./setup-site-bootstrap.sh |
| :---- |

What it does:

1. **Determines SITE\_UUID**  
   * From SITE\_UUID env or by parsing site.sql.  
2. **Ensures nettools-pod is running**  
   * kubectl apply \-f nettools-pod.yaml, waits for phase Running.  
3. **Creates/ensures the Site CR via cloud‑site‑manager**  
   * Calls https://sitemgr.cloud-site-manager:8100/v1/site from nettools-pod with JSON payload containing siteuuid=\<SITE\_UUID\>.  
4. **Fetches the CA cert from credsmgr**  
   * curl https://credsmgr.cert-manager:8000/v1/pki/ca/pem → /tmp/cacert.pem, then copies it locally.  
5. **Reads the OTP from the Site CR**  
   * kubectl get site site-\<SITE\_UUID\> \-n cloud-site-manager \-o jsonpath='{.status.otp.passcode}'  
6. **Ensures namespace elektra-site-agent exists**.  
7. **Creates/updates two secrets in elektra-site-agent:**  
   * bootstrap-info:  
     * literals site-uuid, otp, creds-url, and cacert file  
   * temporal-cert:  
     * blank literal values for cacertificate, certificate, key, otp (scaffold only).

Typical output:

| Using SITE\_UUID=12bd0e80-3382-40d6-82d0-0b20dc14be11Ensuring nettools-pod exists...nettools-pod is Running.Creating site via site-manager...sites.forge.nvidia.io "site-12bd0e80-3382-40d6-82d0-0b20dc14be11" already existsFetching CA certificate from credsmgr...Reading OTP from Site CR...OTP=4waJOq-GEETjPTS134HCCw0DIiU=Ensuring namespace elektra-site-agent exists...namespace/elektra-site-agent createdCreating / updating bootstrap-info secret in elektra-site-agent...secret/bootstrap-info createdCreating / updating temporal-cert secret in elektra-site-agent...secret/temporal-cert created🎉  Site-agent bootstrap secrets created. No restart performed. |
| :---- |

At this point:

* DB rows and Temporal namespace exist.  
* Site CR, OTP, CA, and bootstrap-info / temporal-cert secrets exist.  
* You are ready to wire the UUID into the overlay and deploy Elektra.


### 

### 

### **7.14.6 Overlay configuration (per‑site UUID) & deployment**

The overlay lives under elektra-site-agent/deploy/overlay/ and imports the base:

| resources:  \- ../baseimages:  \- name: nvcr.io/nvidian/nvforge-devel/forge-elektra    newTag: v2025.10.10-rc1-0configMapGenerator:  \- name: elektra-config-map    behavior: merge    env: config.properties |
| :---- |

The **only fields you must always adjust per site** are in overlay/config.properties:

| cluster\_id=12bd0e80-3382-40d6-82d0-0b20dc14be11temporal\_host=temporal-frontend-headless.temporal.svc.cluster.localtemporal\_port=7233temporal\_publish\_namespace=sitetemporal\_publish\_queue=sitetemporal\_subscribe\_namespace=12bd0e80-3382-40d6-82d0-0b20dc14be11temporal\_subscribe\_queue=sitecarbide\_grpc=2 |
| :---- |

**SA actions:**

1. Set cluster\_id to the **exact SITE\_UUID** printed by gen-site-sql.sh.  
2. Set temporal\_subscribe\_namespace to the **same SITE\_UUID**.  
3. Adjust temporal\_host (and temporal\_port if needed) if your Temporal frontend Service name differs.  
4. Optionally tweak carbide\_grpc / other knobs according to policy.

This ensures:

* Elektra’s CLUSTER\_ID and TEMPORAL\_SUBSCRIBE\_NAMESPACE match the DB/site/Temporal namespace.  
* Elektra connects to the correct Temporal host for this site.

#### Deploy Elektra

With:

* DB Secret \+ elektra-database-config in place,  
* Temporal endpoint \+ namespaces (cloud, site, and \<SITE\_UUID\>) ready,  
* cert‑manager and carbide-ca-issuer configured,  
* bootstrap-info & temporal-cert created via setup-site-bootstrap.sh,  
* overlay/config.properties updated as above,

deploy via:

| kubectl apply \-k ./overlay |
| :---- |

Expected creation:

| namespace/elektra-site-agent configuredserviceaccount/site-agent createdrole.rbac.authorization.k8s.io/site-agent createdclusterrole.rbac.authorization.k8s.io/cert-manager-policy:elektra-site-agent createdrolebinding.rbac.authorization.k8s.io/site-agent createdclusterrolebinding.rbac.authorization.k8s.io/cert-manager-policy:elektra-site-agent createdconfigmap/elektra-config-map createdconfigmap/elektra-database-config createdservice/elektra createdservice/elektra-headless createdstatefulset.apps/elektra-site-agent createdcertificate.cert-manager.io/grpc-client-cert createdcertificaterequestpolicy.policy.cert-manager.io/site-agent-approver-policy created |
| :---- |

Verify pods:

| kubectl \-n elektra-site-agent get pods |
| :---- |

You should eventually see something like:

| elektra-site-agent-0   1/1   Running   ...elektra-site-agent-1   1/1   Running   ...elektra-site-agent-2   1/1   Running   ... |
| :---- |

At this point, the site agent is fully wired:

* Uses **mTLS** & SPIFFE to talk to Carbide (and optionally other services).  
* Connects to Temporal with:  
  * publish namespace/queue: site / site  
  * subscribe namespace: \<SITE\_UUID\> / queue site  
* Reads its per‑site identity and credentials from the DB, Temporal, bootstrap-info, and temporal-cert secrets.

### 

## 

## **8\) Networking checklist (Contour/Envoy \+ LB)**

* **Contour/Envoy** (reference: Contour **1.25.2**, Envoy **1.26.4**)  
  * Create Ingress/HTTPProxy entries for any external UIs/APIs you expose.  
  * Use **TLS** via your cert‑manager Issuer (secretName: \<YOUR\_CERT\_SECRET\>).

* **LoadBalancer**  
  * If on‑prem: configure **MetalLB v0.14.5** with your own IPAddressPool CIDRs and **(optionally) BGP** peerings.  
  * If on cloud: use the provider LB and map Service type LoadBalancer on required ports (e.g., cloud-api on 80/9360, carbide-api on 443).

## **9\) Host Ingestion**

After carbide is up and running you can begin ingesting managed hosts. Before you begin that make sure that:

* You have the **forge-admin-cli** command available   
  You can compile it from sources or you can use pre compiled binary. Another choice is to use a containerized version.  
* You can access the carbide site using admin cli.   
  The api service is running at IP address CARBIDE\_API\_EXTERNAL. It is recommended that you add t  
  Run the command 

  **admin-cli \-c [https://api-](https://api-)\<ENVIRONMENT\_NAME\>.\<SITE\_DOMAIN\_NAME\> mi show**  
* DHCP requests from all managed host IPMI networks have been forwarded to service running in Carbide whose IP address is CARBIDE\_DHCP\_EXTERNAL. THis must be completed by the site's networking team.  
* You have the following information for all hosts that need to be ingested:  
  * Mac address of the host BMC  
  * Chassis serial number   
  * Host BMC user name (typically this is factory default username)  
  * Host BMC password (typically this is factory default password)

Upon first login Carbide requires to change factory default credentials to a new set of site wide defaults. You need to set these new credentials as well. These are the new credentials that need to be set:

1. Host BMC Credential  
2. DPU BMC Credential  
3. Host UEFI password  
4. DPU UEFI password

### **Update site default credentials**

Run these commands to set the default credentials. \<**API-URL\>** is typically the following:  
**https://api-\<ENVIRONMENT\_NAME\>.\<SITE\_DOMAIN\_NAME\>**

### **Update DPU UEFI password**

| https://api-\<ENVIRONMENT\_NAME\>.\<SITE\_DOMAIN\_NAME\> |
| :---- |

### **Update Host UEFI password**

Generate a new password using this command:

| admin-cli \-c \<api-url\> host generate-host-uefi-password |
| :---- |

Run this command to update host uefi password:

| admin-cli \-c \<api-url\> credential add-uefi \--kind=host \--password='x' |
| :---- |

	

### **Host and DPU BMC Password**

| admin-cli \-c \<api-url\> credential add-bmc \--kind=site-wide-root \--password='x' |
| :---- |

### **Approve all machine for ingestion**

You need to configure a rule to automatically approve all machines. Failing to do this will cause these machines to raise health alerts. Run this command to complete this step:

| admin-cli \-c \<api-url\> mb site trusted-machine approve \\\* persist \--pcr-registers="0,3,5,6" |
| :---- |

### 

### **Add Expected Machines Table**

First prepare expected machine JSON file as follows:

| {  "expected\_machines": \[    {      "bmc\_mac\_address": "C4:5A:B1:C8:38:0D",      "bmc\_username": "root",      "bmc\_password": "default-password1",      "chassis\_serial\_number": "SERIAL-1"    },    {      "bmc\_mac\_address": "C4:5A:FF:FF:FF:FF",      "bmc\_username": "root",      "bmc\_password": "default-password2",      "chassis\_serial\_number": "SERIAL-2"    }  \]} |
| :---- |

Only servers listed in this table will be ingested. So you have to include all servers in this file. When the file is ready upload it to the site by this command:

| admin-cli \-c \<api-url\> credential em replace-all \--filename expected\_machines.json |
| :---- |

