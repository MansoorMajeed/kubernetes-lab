# Release Strategy for Kubernetes Lab Phases

## Overview

This document outlines the strategy for creating and managing phased releases in the Kubernetes Lab project. The goal is to provide a progressive learning path where each phase builds upon the previous one.

## Git Tag Strategy

### Tag Naming Convention
- **Format**: `v<major>.<minor>.<patch>-<description>`
- **Examples**:
  - `v1.0.0-monitoring-foundation`
  - `v1.1.0-loki-integration`
  - `v2.0.0-basic-service`
  - `v2.1.0-prometheus-metrics`

### Semantic Versioning Guidelines
- **Major version** (v1, v2, v3): Represents a new phase with significant new concepts
- **Minor version** (x.1, x.2, x.3): Incremental improvements within a phase
- **Patch version** (x.x.1): Bug fixes and documentation updates

## Release Process

### Option 1: Direct Tagging (Will follow this while solo)
```bash
# 1. Ensure all changes are committed
git add .
git commit -m "feat: implement prometheus metrics for catalog service"

# 2. Create the tag
git tag -a v2.1.0-prometheus-metrics -m "Phase 2.1.0: Add Prometheus metrics to catalog service"

# 3. Push to remote
git push origin v2.1.0-prometheus-metrics

# 4. Create GitHub release (manual step)
```

### Option 2: Pull Request Workflow (Team Collaboration)
```bash
# 1. Create feature branch
git checkout -b phase-2.1.0-prometheus-metrics

# 2. Make changes and commit
git add .
git commit -m "feat: add prometheus metrics to catalog service"

# 3. Push branch and create PR
git push origin phase-2.1.0-prometheus-metrics

# 4. After PR is merged, create tag on main
git checkout main
git pull origin main
git tag -a v2.1.0-prometheus-metrics -m "Phase 2.1.0: Add Prometheus metrics to catalog service"
git push origin v2.1.0-prometheus-metrics
```

## GitHub Releases

### Creating Releases
1. Go to GitHub repository → Releases → "Create a new release"
2. Select the tag you created
3. Add release notes with:
   - **What's new**: Key features and changes
   - **Learning objectives**: What users will learn
   - **Prerequisites**: Required previous phases
   - **Duration**: Estimated completion time
   - **Breaking changes**: If any (rare in learning environment)

### Release Notes Template
```markdown
## Phase X.Y.Z: [Phase Name]

### 🎯 Learning Objectives
- Learn about [concept 1]
- Implement [technology/pattern]
- Understand [principle]

### 📋 Prerequisites
- Complete Phase X.Y.Z-1
- Basic understanding of [concept]

### ⏱️ Estimated Duration
X-Y hours

### 🚀 What's New
- Added [feature/service]
- Implemented [observability pattern]
- Enhanced [existing functionality]

### 🔧 Getting Started
1. Checkout this phase: `git checkout vX.Y.Z-phase-name`
2. Start the lab: `./start-lab.sh`
3. Follow the guide: `cat phases/phase-X/README.md`

### 📖 Documentation
- Phase Guide: `phases/phase-X/README.md`
- Exercises: `phases/phase-X/EXERCISES.md`
- Troubleshooting: `phases/phase-X/TROUBLESHOOTING.md`

### 🐛 Known Issues
- [Issue 1 and workaround]
- [Issue 2 and workaround]

### 🔄 Next Phase
After completing this phase, move to Phase X.Y.Z+1: [Next Phase Name]
```


This will:
1. Analyze your git history
2. Find appropriate commits for each phase
3. Create tags at historical points
4. Set up the foundation for future phases

### Manual Retrospective Process
If you need more control:

```bash
# 1. Find the commit before Go service was added
git log --oneline --grep="service\|catalog" | tail -n 1

# 2. Create tags at appropriate commits
git tag -a v1.0.0-monitoring-foundation <commit-hash> -m "Phase 1.0.0: Basic monitoring stack"
git tag -a v1.1.0-loki-integration <commit-hash> -m "Phase 1.1.0: Add Loki for log aggregation"

# 3. Tag current state as v2.0.0
git tag -a v2.0.0-basic-service HEAD -m "Phase 2.0.0: Deploy catalog service"

# 4. Push all tags
git push origin --tags
```

## Phase Testing

### Before Creating a Tag
1. **Test the phase independently**:
   ```bash
   ./start-lab.sh
   # Test all functionality
   ```

2. **Verify documentation**:
   ```bash
   cat phases/phase-X/README.md
   # Ensure instructions are clear and complete
   ```

3. **Check for completeness**:
   - All promised features work
   - Documentation is up-to-date
   - Examples and exercises are functional

### After Creating a Tag
1. **Test the checkout process**:
   ```bash
   ./scripts/manage-phases.sh checkout vX.Y.Z-phase-name
   ./start-lab.sh
   ```

2. **Verify the learning path**:
   - Can a learner complete the phase?
   - Are prerequisites clear?
   - Do exercises work as expected?

## Best Practices

### Development Workflow
1. **Always work on main** for this learning project (unless collaborating)
2. **Document as you go** - update phase docs with each change
3. **Test frequently** - ensure each phase works independently
4. **Version incrementally** - don't skip version numbers

### Tag Management
1. **Use descriptive tags** - tag names should be self-explanatory
2. **Create tags promptly** - don't let features accumulate without tagging
3. **Follow the progression** - ensure phases build logically
4. **Test phase transitions** - verify learners can move between phases

### Documentation
1. **Update phase docs** with each release
2. **Include troubleshooting** for common issues
3. **Provide clear examples** and exercises
4. **Link to relevant resources** and documentation

## Troubleshooting

### Common Issues

**Tags not showing up**:
```bash
git push origin --tags
```

**Wrong commit tagged**:
```bash
git tag -d v1.0.0-wrong-tag
git push origin :refs/tags/v1.0.0-wrong-tag
```

**Phase won't start**:
```bash
./start-lab.sh --reset
```

## Tools and Scripts

### Phase Management Script
Use `./scripts/manage-phases.sh` for common operations:
- `list-phases`: Show all available phases
- `checkout <phase>`: Switch to a specific phase
- `create-tag <tag> <message>`: Create a new phase tag
- `current-phase`: Show current phase
- `create-retrospective-tags`: Set up historical tags

### GitHub Integration
- Use GitHub releases for user-friendly documentation
- Tag releases with clear version numbers
- Include comprehensive release notes
- Link to phase documentation

---

## Summary

The phased approach allows learners to:
1. Start with foundational concepts
2. Build complexity incrementally
3. Understand each technology before moving to the next
4. Have clear checkpoints and recovery points
5. Focus on specific learning objectives

This strategy transforms a complex system into a manageable learning journey. 