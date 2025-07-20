# End-of-Session Documentation Workflow

## 🎯 Purpose
Keep documentation synchronized with code changes in a token-efficient way. This workflow ensures AI assistants can quickly find accurate, up-to-date information without wasting tokens on stale documentation.

## 📋 When to Run This
- **End of development sessions** (before git commit/push)
- **Before creating git tags** for phase releases
- **After significant architectural changes**
- **When adding new services or major features**

## 🚀 Quick Workflow

### 1. Analyze Changes
```bash
# Quick analysis of recent changes
./scripts/sync-docs.sh analyze

# Or check since last tag
./scripts/sync-docs.sh analyze --since-tag

# Or check specific range
./scripts/sync-docs.sh analyze HEAD~10
```

**What it does:**
- Analyzes git diff to see what files changed
- Maps changes to specific documentation sections  
- Suggests exactly which docs need updates (no guessing!)

### 2. Update Suggested Documentation
The script will output targeted suggestions like:
```
🔧 Service code changed → Check services/catalog/README.md#api-reference
📊 Observability code changed → Update services/catalog/README.md#observability-features-deep-dive
☸️  K8s manifests changed → Check README.md#architecture diagrams
```

**Token-efficient approach:** Only read and update the specific sections mentioned.

### 3. Validate Documentation Structure
```bash
# Quick validation (doesn't read everything, just checks structure)
./scripts/sync-docs.sh validate
```

**What it checks:**
- Required sections exist in key documentation files
- .cursorrules pointers are valid
- Basic functional test (if services are running)

### 4. Full Analysis (Optional)
```bash
# Run both analyze and validate in one command
./scripts/sync-docs.sh full
```

## 🎯 Token-Efficient Principles

### Smart Change Detection
Instead of reading all docs to understand what's current:
- ✅ **Git diff analysis** tells us exactly what changed
- ✅ **Mapping rules** tell us which docs are affected
- ✅ **Targeted updates** focus only on specific sections

### Predictable Documentation Ownership
```
File Pattern                    → Documentation Owner
services/*/internal/*.go        → services/*/README.md#api-reference
services/*/internal/tracing/    → services/*/README.md#observability
k8s/observability/*.yaml       → README.md#architecture  
k8s/apps/*/                    → README.md#detailed-setup
scripts/*.sh                   → scripts/README.md
Tiltfile                       → README.md#launch-the-lab
```

### Quick Validation (Not Full Reads)
```bash
# Test functionality, don't read everything
curl -s catalog.kubelab.lan:8081/health     # API works?
./kubectl-lab get pods -A | grep Running    # Services deployed?
grep -q "Expected Section" README.md        # Structure intact?
```

## 📚 Documentation Update Guidelines

### When Service Code Changes
- **API changes** → Update `services/*/README.md#testing-examples` with new curl examples
- **Observability changes** → Update `services/*/README.md#observability-features-deep-dive`
- **New patterns** → Consider updating `.cursorrules` decision tree

### When Infrastructure Changes  
- **K8s manifests** → Update `README.md#architecture` diagrams if needed
- **Helm values** → Update `README.md#detailed-setup` instructions
- **Ingress changes** → Update `README.md#host-configuration`

### When Adding New Services
- **Create service README** following `services/catalog/README.md` pattern
- **Update main README** architecture diagrams
- **Update .cursorrules** documentation roadmap table
- **Update traffic simulation** if applicable

## 🤖 AI Assistant Integration

### For AI: Use This Workflow
```bash
# At end of session, run analysis
./scripts/sync-docs.sh analyze --since-tag

# Follow the suggestions (read only specific sections mentioned)
# Update only the targeted documentation sections
# Validate the results

./scripts/sync-docs.sh validate
```

### Smart Response Pattern
When user asks "Is documentation up to date?":
1. **Run the script** to see what changed
2. **Read only suggested sections** (not everything)
3. **Validate specific integration points** (don't test everything)
4. **Report findings** with specific actions needed

## 🎯 Example Session End

```bash
# 1. What changed?
./scripts/sync-docs.sh analyze --since-tag

# Output:
# 🔧 Service code changed → Check services/catalog/README.md#api-reference
# 📊 Observability code changed → Update services/catalog/README.md#observability-features-deep-dive

# 2. Update only those specific sections (token-efficient!)

# 3. Validate
./scripts/sync-docs.sh validate

# 4. Commit with confidence that docs are synchronized
```

## 💡 Benefits

✅ **Token Efficient** - Only read/update what actually changed  
✅ **Consistent** - Predictable documentation ownership  
✅ **Fast Context** - AI knows exactly where to look  
✅ **Prevents Drift** - Systematic but lightweight process  
✅ **Scalable** - Works as project grows  

This workflow turns documentation maintenance from "read everything to understand everything" into "analyze changes to know exactly what to check." 