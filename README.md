# Kubernetes Observability Lab

Note: This is a work in progress

**Production-like observability stack in 5 minutes, locally, completely free** ⚡

I made this as a way to learn Kubernetes + Microservices + observability in a modern production environment. Main motivation was to have
everything run locally so that you don't need spend money on cloud accounts (especially if you are a student or looking for your first job).
I decided to go with Tilt + k3d to setup a local kubernetes environment

Additionally, I wanted to have this local infrastructure mimic a real world setup, so I am building an e-commerce platform locally called "LocalMart".
For now it has only one backned service (catalog service), but the idea is to have more microservices and show the interaction between them, have logs, metrics and traces etc.

I think this is especially useful for those who have not managed any production environments to get an idea of how things might look.

I think this would be helpful for new DevOps engineers, SREs and developers. It might be helpful for some experienced folks as well.

Making all of these manually is going to take a ton of time, but luckily, I was able to use generative AI (Cursor + Claude + Gemin 2.5 pro) for brainstorming, writing docs, writing boilerplate code etc and this sped things up significantly (check the section on AI usage at the bottom to read about my  learnings)

I have organized this repo into different "phases", take a look at the [Release page](https://github.com/MansoorMajeed/kubernetes-lab/releases) to get an idea. You can switch any of the tags and focus only on that phase, if you would like

If you have any questions, you can open an issue in this repo. 

## ⚡ Quick Start (Returning Users)

**Requirements:** Docker, k3d, kubectl, Tilt, Helm

```bash
# 1. Start the lab
./start-lab.sh

# 2. Deploy everything  
./tilt-lab up
```

**Access:** Add to `/etc/hosts` then visit [localmart.kubelab.lan:8081](http://localmart.kubelab.lan:8081) or [grafana.kubelab.lan:8081](http://grafana.kubelab.lan:8081) (admin/password)
Check the `Host configuration` section below for all the hosts to add to the hosts file

## 🎯 Current Status: Phase 3.0.0 Complete ✅

**✅ Phase 1-2**: Observability stack + Go catalog service  
**✅ Phase 3.0.0**: React frontend with product browsing and details  
**🔮 Phase 4.0.0**: Shopping cart service (Go + Redis + gRPC)  
**🔮 Phase 5.0.0**: Review service (Python + MongoDB)  

**Try it now:** [LocalMart E-commerce →](http://localmart.kubelab.lan:8081)

---

## 🎯 What You'll Build

A **complete observability stack** using production patterns, optimized for local learning:

> **Production Patterns, Local Scale**: This setup mirrors real-world observability architecture and configurations, but scaled down to run efficiently on your laptop. Perfect for learning without the complexity or cost of cloud infrastructure.

✅ **Distributed Tracing** - See request flows across microservices  
✅ **Structured Logging** - JSON logs with trace correlation  
✅ **Metrics & Dashboards** - Prometheus + Grafana monitoring  
✅ **Real Microservice** - Go catalog API (more services to come) with full instrumentation  
✅ **Traffic Simulation** - Generate realistic observability data  

**Learning Focus:** Hands-on experience with production patterns, not toy examples.

---

## 🏗️ Architecture

### Lab Infrastructure
```mermaid
graph TB
    Dev[👩‍💻 Developer] --> Tilt[🔄 Tilt<br/>Auto-deploy]
    Tilt --> k3d[☸️ k3d Cluster<br/>Local Kubernetes]
    
    k3d --> Observability[📊 Observability Stack]
    k3d --> Services[🛍️ Application Services]
    
    Observability --> Prometheus[📈 Prometheus<br/>Metrics]
    Observability --> Grafana[📊 Grafana<br/>Dashboards] 
    Observability --> Loki[📝 Loki<br/>Logs]
    Observability --> Tempo[🔍 Tempo<br/>Traces]
    Observability --> Alloy[🔗 Alloy<br/>Collector]
    
    Services --> Catalog[🛒 Catalog Service<br/>Go + PostgreSQL]
    
    style Dev fill:#e1f5fe
    style Observability fill:#f3e5f5
    style Services fill:#e8f5e8
```

### Service Architecture (Current ✅ + Future 🔮)
```mermaid
graph TB
    Client[🌐 Client] --> Frontend[⚛️ LocalMart Frontend<br/>✅ React + Tailwind CSS]
    Frontend --> LB[⚖️ Load Balancer<br/>✅ k3d ingress]
    
    %% Current Services (Solid lines)
    LB --> CatalogAPI[🛒 Catalog Service<br/>✅ Go + REST API<br/>Current v2.2.0]
    CatalogAPI --> PostgreSQL[(🗄️ PostgreSQL<br/>✅ Products)]
    
    %% Future Services (Dotted lines)  
    LB -.-> CartAPI[🛍️ Cart Service<br/>🔮 Go + REST API<br/>Phase 4.0.0]
    LB -.-> ReviewAPI[⭐ Review Service<br/>🔮 Python + REST API<br/>Phase 5.0.0]
    
    %% Future Service Communication (Dotted lines)
    CartAPI -.->|"🚀 gRPC<br/>Fast validation"| CatalogGRPC[🛒 Catalog gRPC<br/>🔮 Product validation]
    ReviewAPI -.->|"🌐 REST<br/>Rich product data"| CatalogAPI
    CatalogGRPC -.-> PostgreSQL
    
    %% Future Databases (Dotted lines)
    CartAPI -.-> Redis[(🔴 Redis<br/>🔮 Cart Sessions)]
    ReviewAPI -.-> MongoDB[(🍃 MongoDB<br/>🔮 Reviews + Ratings)]
    
    %% Observability Stack (Current)
    subgraph "📊 Observability Stack ✅"
        Prometheus[📈 Prometheus]
        Grafana[📊 Grafana] 
        Loki[📝 Loki]
        Tempo[🔍 Tempo]
    end
    
    %% Current Observability Connections (Solid lines)
    Frontend --> Prometheus
    CatalogAPI --> Prometheus
    CatalogAPI --> Loki
    CatalogAPI --> Tempo
    
    %% Future Observability Connections (Dotted lines)
    CartAPI -.-> Prometheus
    ReviewAPI -.-> Prometheus
    CartAPI -.-> Loki  
    ReviewAPI -.-> Loki
    CartAPI -.-> Tempo
    ReviewAPI -.-> Tempo
    
    %% Styling - Current (bright colors)
    style Frontend fill:#e3f2fd,stroke:#1976d2,stroke-width:3px
    style CatalogAPI fill:#c8e6c9,stroke:#388e3c,stroke-width:3px
    style PostgreSQL fill:#e1f5fe,stroke:#0277bd,stroke-width:3px
    style LB fill:#f3e5f5,stroke:#7b1fa2,stroke-width:3px
    
    %% Styling - Future (muted colors)
    style CartAPI fill:#fff3e0,stroke:#f57c00,stroke-width:2px,stroke-dasharray: 5 5
    style ReviewAPI fill:#fce4ec,stroke:#c2185b,stroke-width:2px,stroke-dasharray: 5 5
    style CatalogGRPC fill:#e8f5e8,stroke:#4caf50,stroke-width:2px,stroke-dasharray: 5 5
    style Redis fill:#ffebee,stroke:#f44336,stroke-width:2px,stroke-dasharray: 5 5
    style MongoDB fill:#e0f2f1,stroke:#00695c,stroke-width:2px,stroke-dasharray: 5 5
```

**Legend:**
- **✅ Solid lines & bright colors**: Currently implemented and working
- **🔮 Dotted lines & muted colors**: Planned for future phases

### **What's Coming Next** 🚀

| **Phase** | **Component** | **Technology** | **Learning Focus** |
|-----------|-------------|----------------|-------------------|
| **3.0.0** | 🌐 **Frontend** | React/Vue | Frontend-backend integration, full-stack tracing |
| **4.0.0** | 🛍️ **Cart Service + UI** | Go + Redis + gRPC | High-speed caching, gRPC performance |
| **5.0.0** | ⭐ **Review Service + UI** | Python + MongoDB | NoSQL patterns, multi-language stack |


### **Why This Progression** 🎯
- **🎨 Immediate Value**: Working UI showcases existing services from Phase 3.0
- **🔄 Iterative Development**: Each phase enhances the same UI with new backend capabilities
- **📚 Technology Diversity**: React/Vue → Go + Redis + gRPC → Python + MongoDB
- **🏭 Real-World Flow**: Frontend first, then progressive backend enhancement (like real startups!)
- **📊 Complete Learning**: Full-stack tracing, protocol comparison, database diversity

**🎉 Ready to explore the future?** [**Dive into the architectural planning →**](./docs/brainstorm/)

---

## 🚀 Detailed Setup

### Prerequisites

| Tool | Purpose | Install |
|------|---------|---------|
| [Docker](https://docs.docker.com/engine/install/) | Container runtime | Required |
| [k3d](https://k3d.io/stable/#installation) | Local Kubernetes | Required |
| [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) | Kubernetes CLI | Required |
| [Tilt](https://docs.tilt.dev/) | Auto-deployment | Required |
| [Helm](https://github.com/helm/helm/releases) | Package manager | Required |

### Host Configuration

Add these entries to your `/etc/hosts` file:

```bash
# Observability Stack
127.0.0.1 grafana.kubelab.lan
127.0.0.1 prometheus.kubelab.lan
127.0.0.1 tempo.kubelab.lan

# Application Services  
127.0.0.1 catalog.kubelab.lan
127.0.0.1 localmart.kubelab.lan

# Demo Apps
127.0.0.1 nginx-hello.kubelab.lan
```

**Platform-specific:**
- **macOS/Linux:** `sudo nano /etc/hosts`
- **Windows:** Edit `C:\Windows\System32\drivers\etc\hosts` as Administrator

### Launch the Lab

```bash
# Start k3d cluster and basic setup
./start-lab.sh

# Deploy the complete observability stack
./tilt-lab up
```

**What happens:**
1. Creates local Kubernetes cluster with k3d
2. Configures ingress for local domain access
3. Deploys Prometheus, Grafana, Loki, Tempo via Helm
4. Builds and deploys the catalog microservice
5. Sets up traffic simulation

---

## 🔍 Explore the Lab

### 📊 Observability Stack

| Component | Access | Purpose |
|-----------|--------|---------|
| **Grafana** | [grafana.kubelab.lan:8081](http://grafana.kubelab.lan:8081) | Dashboards, explore traces/logs/metrics |
| **Prometheus** | [prometheus.kubelab.lan:8081](http://prometheus.kubelab.lan:8081) | Metrics collection and queries |
| **Tempo** | [tempo.kubelab.lan:8081](http://tempo.kubelab.lan:8081) | Distributed tracing backend |
| **Tilt UI** | [localhost:10350](http://localhost:10350/) | Deployment management |

**Credentials:** admin / password

### 🛍️ Application Services

| Service | Access | Documentation |
|---------|--------|---------------|
| **LocalMart Frontend** | [localmart.kubelab.lan:8081](http://localmart.kubelab.lan:8081) | Complete e-commerce UI (Phase 3.0.0) |
| **Catalog API** | [catalog.kubelab.lan:8081](http://catalog.kubelab.lan:8081) | [Full API docs →](./services/catalog/) |
| **Demo App** | [nginx-hello.kubelab.lan:8081](http://nginx-hello.kubelab.lan:8081) | Simple test service |

### 🛠️ Helper Scripts

```bash
# Use lab-specific kubectl (ensures correct context)
./kubectl-lab get pods -A

# Generate realistic traffic for observability data
./scripts/simulate-traffic.sh

# Stop/start the lab
./tilt-lab down
./tilt-lab up
```

---

## 📚 Learning Resources

### 🎓 Progressive Learning Path

**Phase System:** Use git tags to explore different complexity levels:

```bash
# See all available phases
git tag | grep -E "^v[0-9]"

# Example: Start with monitoring basics
git checkout v1.0.0-monitoring-foundation

# Progress to distributed tracing  
git checkout v2.3.0-distributed-tracing
```

[**📖 View all phases and releases →**](https://github.com/mansoormajeed/kubernetes-lab/releases)

### 🏗️ Service Deep Dives

- **[Catalog Service →](./services/catalog/)** - Complete microservice with observability
- **[Traffic Simulation →](./scripts/)** - Generate realistic e-commerce patterns

### 🔧 Customization

- **Helm Values:** Modify `k8s/observability/*-values.yaml` 
- **Service Config:** Edit `k8s/apps/*/` deployments
- **Add Services:** Follow catalog service pattern

---

## 🚀 Traffic Simulation

Generate realistic observability data:

```bash
# Quick test with realistic e-commerce patterns
./scripts/simulate-traffic.sh --duration 300 --interval 2

# Heavy load for observability testing
./scripts/simulate-traffic.sh --duration 600 --interval 1 --verbose

# Seed database with sample products only
./scripts/simulate-traffic.sh --seed-only
```

**Generates:** Rich traces, correlated logs, realistic metrics perfect for learning observability patterns.

---

## 🤖 Built with AI Assistance

This project extensively used **AI pair programming** with Cursor/Claude. Here are key learnings for others using LLMs for infrastructure work:

### ✅ What Works Incredibly Well

1. **📖 Documentation & Code Review** - LLMs excel at reading existing code and writing comprehensive docs
2. **🏗️ Boilerplate Generation** - Perfect for repetitive Kubernetes manifests, Go handlers, etc.
3. **🔍 Architecture Analysis** - Great at suggesting improvements and identifying patterns
4. **🐛 Debugging Assistance** - Excellent at analyzing error logs and suggesting fixes

### ⚠️ Important Caveats

1. **📦 Don't Let LLMs Manage Dependencies**
   - Always check official docs for latest versions
   - Use proper tools (`go get`, `helm repo add`) instead of LLM suggestions
   - LLM knowledge of versions is often outdated

2. **🧪 Always Verify Generated Configs**
   - Test Kubernetes manifests in isolation
   - Validate Helm values against chart documentation  
   - Check that environment variables actually exist

3. **🔗 Integration Points Need Human Review**
   - Service discovery, networking, RBAC policies
   - OpenTelemetry configuration and endpoints
   - Database connection strings and credentials

### 💡 Pro Tips for AI-Assisted Infrastructure

- **Ask for comprehensive code reviews** of your entire setup
- **Use AI to explain complex error messages** (especially Kubernetes events)
- **Generate test scenarios** for validating your observability setup
- **Document architecture decisions** with AI help, but verify technical accuracy

**Result:** This approach enabled rapid development of production-like patterns while maintaining code quality and learning value.

---

## 🤝 Contributing

- 🐛 **Found a bug?** Open an issue with details
- 💡 **Have an idea?** Suggest new learning scenarios  
- 🔧 **Want to add a service?** Follow the catalog service pattern
- 📚 **Improve docs?** All improvements welcome

**Learning Focus:** Contributions should maintain educational value and production-like patterns.

---

**🎉 Ready to explore production observability?** Start with `./start-lab.sh` and dive into the Grafana dashboards!

