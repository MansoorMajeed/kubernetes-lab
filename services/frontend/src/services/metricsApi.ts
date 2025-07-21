// Frontend Metrics Collection Service

// So this is highly experimental idea cursor came up with. I could not find much information
// on this pattern on the internet. But I think it might work at least for learning


export interface PerformanceMetric {
  name: string
  value: number
  timestamp: number
  labels?: Record<string, string>
}

export interface BusinessEvent {
  event: string
  properties?: Record<string, any>
  timestamp: number
}

export interface ErrorEvent {
  error: string
  stack?: string
  component?: string
  timestamp: number
  url: string
}

// Metrics collection service
class MetricsService {
  private metricsQueue: PerformanceMetric[] = []
  private eventsQueue: BusinessEvent[] = []
  private errorsQueue: ErrorEvent[] = []
  private flushInterval = 30000 // 30 seconds
  private maxQueueSize = 100

  constructor() {
    this.startPerformanceObserver()
    this.startErrorListener()
    this.startFlushTimer()
  }

  // Performance Metrics
  recordPerformanceMetric(name: string, value: number, labels?: Record<string, string>) {
    const metric: PerformanceMetric = {
      name,
      value,
      timestamp: Date.now(),
      labels: {
        ...labels,
        user_agent: navigator.userAgent.slice(0, 50), // Truncate for storage
        url: window.location.pathname,
      }
    }
    
    this.metricsQueue.push(metric)
    this.checkQueueSize()
  }

  // Business Events
  recordEvent(event: string, properties?: Record<string, any>) {
    const businessEvent: BusinessEvent = {
      event,
      properties: {
        ...properties,
        url: window.location.pathname,
        referrer: document.referrer,
      },
      timestamp: Date.now(),
    }
    
    this.eventsQueue.push(businessEvent)
    this.checkQueueSize()
  }

  // Error Tracking
  recordError(error: Error, component?: string) {
    const errorEvent: ErrorEvent = {
      error: error.message,
      stack: error.stack?.slice(0, 500), // Truncate stack trace
      component,
      timestamp: Date.now(),
      url: window.location.pathname,
    }
    
    this.errorsQueue.push(errorEvent)
    this.checkQueueSize()
  }

  // Core Web Vitals Observer
  private startPerformanceObserver() {
    // Check if browser supports Performance Observer
    if ('PerformanceObserver' in window) {
      try {
        // Largest Contentful Paint (LCP)
        const lcpObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries()
          const lastEntry = entries[entries.length - 1] as any
          this.recordPerformanceMetric('frontend_lcp_seconds', lastEntry.startTime / 1000)
        })
        lcpObserver.observe({ entryTypes: ['largest-contentful-paint'] })

        // First Input Delay (FID)
        const fidObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries()
          entries.forEach((entry: any) => {
            this.recordPerformanceMetric('frontend_fid_milliseconds', entry.processingStart - entry.startTime)
          })
        })
        fidObserver.observe({ entryTypes: ['first-input'] })

        // Cumulative Layout Shift (CLS)
        let clsValue = 0
        const clsObserver = new PerformanceObserver((list) => {
          const entries = list.getEntries()
          entries.forEach((entry: any) => {
            if (!entry.hadRecentInput) {
              clsValue += entry.value
            }
          })
          this.recordPerformanceMetric('frontend_cls_score', clsValue)
        })
        clsObserver.observe({ entryTypes: ['layout-shift'] })

        // Page Load Time
        window.addEventListener('load', () => {
          const loadTime = performance.now()
          this.recordPerformanceMetric('frontend_page_load_seconds', loadTime / 1000)
        })

      } catch (error) {
        console.warn('Performance Observer not fully supported:', error)
      }
    }
  }

  // Global Error Listener
  private startErrorListener() {
    window.addEventListener('error', (event) => {
      this.recordError(new Error(event.message), 'global')
    })

    window.addEventListener('unhandledrejection', (event) => {
      this.recordError(new Error(event.reason), 'promise')
    })
  }

  // Auto-flush timer
  private startFlushTimer() {
    setInterval(() => {
      this.flush()
    }, this.flushInterval)
  }

  // Queue size management
  private checkQueueSize() {
    if (this.metricsQueue.length + this.eventsQueue.length + this.errorsQueue.length > this.maxQueueSize) {
      this.flush()
    }
  }

  // Send metrics to backend
  async flush() {
    if (this.metricsQueue.length === 0 && this.eventsQueue.length === 0 && this.errorsQueue.length === 0) {
      return
    }

    const payload = {
      performance_metrics: [...this.metricsQueue],
      business_events: [...this.eventsQueue],
      error_events: [...this.errorsQueue],
      timestamp: Date.now(),
      session_id: this.getSessionId(),
    }

    try {
      // Send to catalog service metrics endpoint
      await fetch('/api/v1/frontend-metrics', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })

      // Clear queues after successful send
      this.metricsQueue = []
      this.eventsQueue = []
      this.errorsQueue = []

    } catch (error) {
      console.warn('Failed to send metrics:', error)
      // Don't clear queues on failure - will retry on next flush
    }
  }

  // Simple session ID generation
  private getSessionId(): string {
    let sessionId = sessionStorage.getItem('frontend_session_id')
    if (!sessionId) {
      sessionId = Math.random().toString(36).substring(2, 15)
      sessionStorage.setItem('frontend_session_id', sessionId)
    }
    return sessionId
  }

  // Manual flush for critical events
  flushImmediate() {
    return this.flush()
  }
}

// Singleton instance
export const metricsService = new MetricsService()

// Convenience functions
export const recordPageView = (page: string) => {
  metricsService.recordEvent('page_view', { page })
}

export const recordProductView = (productId: string | number, productName: string) => {
  metricsService.recordEvent('product_view', { 
    product_id: productId.toString(),
    product_name: productName 
  })
}

export const recordAPICall = (endpoint: string, duration: number, success: boolean) => {
  metricsService.recordPerformanceMetric('frontend_api_call_duration_seconds', duration / 1000, {
    endpoint,
    success: success.toString(),
  })
}

export const recordReactQueryMetric = (queryKey: string, duration: number, fromCache: boolean) => {
  metricsService.recordPerformanceMetric('frontend_react_query_duration_seconds', duration / 1000, {
    query_key: queryKey,
    from_cache: fromCache.toString(),
  })
} 