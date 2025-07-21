# LocalMart Frontend Service

**Phase 3.0.0**: React + TypeScript + Tailwind CSS e-commerce frontend

## 🎯 Overview

The LocalMart frontend provides a modern, responsive e-commerce interface built with React and TypeScript. It demonstrates production-ready frontend patterns including API integration, routing, responsive design, and modern CSS frameworks.

## ✨ Features

### Core Functionality
- **Product Browsing**: Responsive grid layout with product cards
- **Product Details**: Individual product pages with navigation
- **API Integration**: Consumes Catalog Service REST API
- **Client-Side Routing**: Single Page Application with React Router
- **Responsive Design**: Mobile-first approach with Tailwind CSS

### Technical Features
- **TypeScript**: Type-safe development
- **Tailwind CSS v4**: Latest utility-first CSS framework
- **Vite**: Fast development server and build tool
- **React Router**: Declarative routing for SPA
- **Error Handling**: Graceful API error states
- **Loading States**: User feedback during API calls
- **Frontend Observability**: Real User Monitoring with Core Web Vitals
- **React Query**: Caching and API state management

## 🏗️ Architecture

### Component Structure
```
src/
├── components/
│   ├── ProductList.tsx     # Product grid with API integration
│   ├── ProductDetail.tsx   # Individual product view
│   ├── layout/            # Header, Footer, Layout components
│   └── ui/                # Reusable UI components + ErrorBoundary
├── hooks/
│   └── useProducts.ts     # React Query hooks with metrics
├── services/
│   ├── catalogApi.ts      # API service layer
│   └── metricsApi.ts      # Frontend metrics collection
├── types/
│   └── api.ts            # TypeScript type definitions
├── App.tsx               # Router setup with Query Provider
├── main.tsx             # React app entry point
└── index.css            # Tailwind configuration
```

### Technology Stack
- **React 18**: Modern React with hooks
- **TypeScript**: Static typing for better DX
- **Tailwind CSS v4**: Utility-first styling
- **React Router v6**: Client-side routing
- **React Query**: Data fetching and caching
- **Vite**: Build tool and dev server

### React Query Integration
```typescript
// Modern data fetching with caching and metrics
export function useProducts() {
  return useQuery({
    queryKey: ['products'],
    queryFn: async () => {
      const startTime = performance.now()
      const result = await catalogApi.getProducts()
      const duration = performance.now() - startTime
      
      // Automatic performance tracking
      recordAPICall('/api/v1/products', duration, true)
      return result
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}
```

## 🚀 Development

### Local Development (Vite)
```bash
cd services/frontend
npm install
npm run dev
# Access at http://localhost:5173
```

### Production Build
```bash
npm run build
# Creates optimized build in dist/
```

### Code Quality
```bash
npm run lint    # ESLint checking
npm run preview # Preview production build
```


### Nginx Configuration
- **Static File Serving**: Optimized for React SPA
- **API Proxy**: Routes `/api/*` to catalog service
- **SPA Routing**: Fallback to index.html for client-side routes
- **Performance**: Gzip compression and caching headers

## ☸️ Kubernetes Deployment

### Resources
```bash
k8s/apps/frontend/
├── namespace.yaml     # frontend namespace
├── deployment.yaml    # React app deployment
├── service.yaml      # ClusterIP service
└── ingress.yaml      # External access
```

### Deployment Features
- **Replicas**: 2 pods for availability
- **Health Checks**: Readiness and liveness probes
- **Resource Limits**: CPU and memory constraints
- **Rolling Updates**: Zero-downtime deployments

### Access
- **External**: http://localmart.kubelab.lan:8081
- **Internal**: `frontend-service.frontend.svc.cluster.local`

## 🎨 UI/UX Features

### Responsive Design
```css
/* Mobile-first responsive grid */
.grid {
  @apply grid grid-cols-1 gap-6;
  @apply sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4;
}
```

### Design System
- **Color Palette**: Modern neutral tones
- **Typography**: Clear hierarchy with Tailwind
- **Spacing**: Consistent spacing scale
- **Components**: Reusable card and button patterns

### User Experience
- **Loading States**: Skeleton loading for better UX
- **Error Handling**: User-friendly error messages
- **Navigation**: Intuitive product browsing
- **Performance**: Optimized images and assets

## 📊 Frontend Observability


> Note: this is highly experimental. I could not find much reference to this pattern of using
> prometheus to store frontend metrics. While it sounds odd, I think it is cool

### Real User Monitoring (RUM)
The frontend automatically collects performance and business metrics using modern browser APIs:

**Performance Metrics:**
- **Core Web Vitals**: LCP, FID, CLS scores
- **Page Load Times**: Complete loading duration
- **API Performance**: Request duration from user perspective
- **React Query Metrics**: Cache hit rates and query performance

**Business Metrics:**
- **Page Views**: Track user navigation patterns
- **Product Views**: Individual product engagement
- **Error Tracking**: JavaScript exceptions and React errors
- **User Sessions**: Journey tracking across the application

### Integration with Prometheus
```typescript
// Metrics automatically sent to backend every 30 seconds
POST /api/v1/frontend-metrics
{
  "performance_metrics": [...],
  "business_events": [...],
  "error_events": [...]
}

// Viewable in existing Grafana dashboards alongside backend metrics
```

### Key Components
- **MetricsService**: Automatic browser performance monitoring
- **React Query Hooks**: API call instrumentation  
- **ErrorBoundary**: React error catching and reporting
- **Business Event Tracking**: User interaction monitoring

## 🔗 API Integration

### Catalog Service Integration
```typescript
interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  image_url: string;
}

// API endpoints consumed
GET /api/products           # Product listing
GET /api/products/:id       # Product details
GET /api/health            # Service health
```

### Error Handling
```typescript
const [error, setError] = useState<string | null>(null);

// Graceful error handling
.catch(err => {
  console.error('API Error:', err);
  setError('Failed to load products. Please try again.');
});
```

## 🚦 Testing the Frontend

### Manual Testing
```bash
# 1. Start the lab
./start-lab.sh
./tilt-lab up

# 2. Visit the frontend
open http://localmart.kubelab.lan:8081

# 3. Test functionality
# - Product grid loads
# - Click product -> details page
# - Back navigation works
# - Responsive on mobile
```

### API Testing
```bash
# Test API proxy through frontend
curl http://localmart.kubelab.lan:8081/api/products
curl http://localmart.kubelab.lan:8081/api/health
```

## 🎯 Learning Objectives

### Frontend Development
- Modern React patterns with hooks
- TypeScript integration and benefits
- Component composition and reusability
- Client-side routing in SPAs

### CSS & Styling
- Tailwind CSS v4 features and setup
- Responsive design principles
- Mobile-first development
- CSS Grid and Flexbox patterns

### API Integration
- Fetch API usage and best practices
- Error handling and loading states
- CORS considerations in development
- REST API consumption patterns

### DevOps & Deployment
- Multi-stage Docker builds
- Nginx configuration for SPAs
- Kubernetes deployment patterns
- Production optimization techniques

## 🔄 Phase Progression

### Phase 3.0.0 ✅ (Current)
- React frontend with product browsing
- Responsive design with Tailwind CSS
- API integration with catalog service
- Docker and Kubernetes deployment

### Phase 3.1.0 🔮 (Next)
- Search functionality
- Enhanced error handling
- Loading skeletons
- Performance optimizations

### Phase 3.2.0 🔮 (Future)
- Frontend observability (tracing)
- Performance monitoring
- Analytics integration
- Advanced UX patterns

## 🛠️ Development Tips

### Tailwind CSS v4
```css
/* Key differences in v4 */
@import "tailwindcss";           /* New import syntax */
/* Content detection via comments */
/* lg:grid-cols-4 xl:grid-cols-4 */
```

### React Router Setup
```typescript
// Clean routing pattern
<Router>
  <Routes>
    <Route path="/" element={<ProductList />} />
    <Route path="/product/:id" element={<ProductDetail />} />
  </Routes>
</Router>
```

### Vite Configuration
```typescript
// Proxy for development
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
});
```

