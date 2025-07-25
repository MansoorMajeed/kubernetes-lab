# Kubernetes Observability Lab - Claude Code Assistant Guide

## 🎯 Project Overview

**LocalMart E-commerce Platform** - A production-like observability stack for learning cloud-native patterns locally. This is a comprehensive Kubernetes lab demonstrating microservices, observability, and modern development practices.

**Current Status**: Phase 3.0.0 (React frontend complete) → Planning Phase 4.0.0 (Cart service)

## 🚀 Quick Commands

```bash
# Essential lab commands (ALWAYS use these, never raw kubectl/tilt)
./setup-lab.sh                    # Initialize k3d cluster
./tilt-lab up                      # Deploy entire stack
./kubectl-lab get pods -A          # Check all services
./scripts/simulate-traffic.sh      # Generate observability data
./tilt-lab down                    # Stop lab

# Testing & Validation
npm test                          # Frontend tests (if applicable)
go test ./...                     # Go service tests
./scripts/sync-docs.sh validate   # Validate documentation
```

## 🏗️ Architecture Summary

**Current Services (Phase 3.0.0):**
- ✅ **Catalog Service** (Go + PostgreSQL) - Product API with full observability
- ✅ **Frontend** (React + Tailwind CSS) - E-commerce UI at localmart.kubelab.lan:8081
- ✅ **Observability Stack** - Prometheus, Grafana, Loki, Tempo

**Planned Services:**
- 🔮 **Cart Service** (Phase 4.0.0) - Go + Redis + gRPC
- 🔮 **Review Service** (Phase 5.0.0) - Python + MongoDB

## 📁 Key Files & Directories

```
├── services/
│   ├── catalog/           # Reference implementation (Go service)
│   └── frontend/          # React e-commerce UI
├── k8s/
│   ├── observability/     # Prometheus, Grafana, Loki, Tempo configs
│   └── apps/              # Service deployments
├── scripts/
│   ├── simulate-traffic.sh    # Generate test data
│   └── sync-docs.sh          # Documentation workflow
├── docs/brainstorm/README.md  # Complete architecture plan
├── Tiltfile              # Automated deployment config
└── .cursorrules          # Development guidelines
```

## 🛠️ Development Workflow

### Adding New Services
1. Follow `services/catalog/` structure and patterns
2. Include observability from day one (tracing, logging, metrics)
3. Use clean architecture (handlers → models → database)
4. Update Tiltfile and ingress configurations
5. Create comprehensive README with API documentation

### Testing Strategy
```bash
# Run appropriate tests based on service
cd services/catalog && go test ./...
cd services/frontend && npm test
./scripts/simulate-traffic.sh --duration 60  # Integration testing
```

### Code Standards
- **Go Services**: Clean architecture, structured logging, OpenTelemetry instrumentation
- **Frontend**: React + Tailwind CSS, responsive design, error handling
- **Observability**: All services must emit metrics, logs, and traces
- **Documentation**: Each service needs comprehensive README

## 🔍 Troubleshooting

### Common Issues
```bash
# Check cluster status
./kubectl-lab get pods -A
./kubectl-lab describe pod <pod-name>

# Observability access
# Grafana: grafana.kubelab.lan:8081 (admin/password)
# Prometheus: prometheus.kubelab.lan:8081
# Application: localmart.kubelab.lan:8081

# Tilt debugging
# Tilt UI: localhost:10350
```

### Host Configuration Required
Add to `/etc/hosts`:
```
127.0.0.1 grafana.kubelab.lan
127.0.0.1 prometheus.kubelab.lan
127.0.0.1 catalog.kubelab.lan
127.0.0.1 localmart.kubelab.lan
```

## 🎓 Learning Focus

**Technology Diversity:**
- **Languages**: Go (performance), Python (flexibility), JavaScript (frontend)
- **Databases**: PostgreSQL (relational), MongoDB (document), Redis (cache)
- **Communication**: REST (standard), gRPC (performance)

**Observability Patterns:**
- Distributed tracing across services and protocols
- Structured logging with trace correlation
- Metrics collection and dashboards
- Real-world debugging scenarios

## 📋 Phase Progression

**Current: Phase 3.0.0** ✅
- Complete React frontend with product browsing
- Full integration with Catalog Service
- Frontend observability and performance monitoring

**Next: Phase 4.0.0** 🔮
- Cart Service (Go + Redis)
- gRPC communication (Cart ↔ Catalog)
- Mixed protocol tracing (REST + gRPC)

**Future: Phase 5.0.0** 🔮
- Review Service (Python + MongoDB)
- Cross-language observability
- Complete e-commerce functionality

## ⚠️ Critical Safety Rules

- ✅ ALWAYS use `./kubectl-lab` (never raw `kubectl`)
- ✅ ALWAYS use `./tilt-lab` (never raw `tilt`)
- ✅ Test Kubernetes manifests in isolation before adding to Tiltfile
- ✅ All configurations are LOCAL ONLY
- ✅ Run `./scripts/sync-docs.sh` at end of sessions

## 🤖 AI Development Guidelines

**Strengths:**
- Documentation and code reviews
- Boilerplate generation (K8s manifests, Go handlers)
- Architecture analysis and debugging assistance

**Requires Human Verification:**
- Dependency management and versions
- OpenTelemetry configurations
- Service discovery and networking
- Integration between observability components

**Best Practices:**
- Use official scaffolding tools for dependency management
- Prioritize simple, educational implementations
- Add clear comments explaining learning objectives
- Let tooling handle version management

## 📊 Access URLs

| Service | URL | Purpose |
|---------|-----|---------|
| **LocalMart Frontend** | [localmart.kubelab.lan:8081](http://localmart.kubelab.lan:8081) | E-commerce UI |
| **Grafana** | [grafana.kubelab.lan:8081](http://grafana.kubelab.lan:8081) | Observability dashboards |
| **Prometheus** | [prometheus.kubelab.lan:8081](http://prometheus.kubelab.lan:8081) | Metrics queries |
| **Catalog API** | [catalog.kubelab.lan:8081](http://catalog.kubelab.lan:8081) | Product API |
| **Tilt UI** | [localhost:10350](http://localhost:10350) | Deployment management |

---

*Generated for Claude Code Assistant - Last updated: 2025-01-25*