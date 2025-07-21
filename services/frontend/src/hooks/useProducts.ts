import { useQuery } from '@tanstack/react-query'
import { catalogApi } from '../services/catalogApi'
import { recordReactQueryMetric, recordAPICall } from '../services/metricsApi'

// Query keys for React Query cache management
export const productKeys = {
  all: ['products'] as const,
  lists: () => [...productKeys.all, 'list'] as const,
  list: (filters: Record<string, unknown>) => [...productKeys.lists(), { filters }] as const,
  details: () => [...productKeys.all, 'detail'] as const,
  detail: (id: string | number) => [...productKeys.details(), id] as const,
}

// Hook to fetch all products
export function useProducts() {
  return useQuery({
    queryKey: productKeys.lists(),
    queryFn: async () => {
      const startTime = performance.now()
      
      try {
        const result = await catalogApi.getProducts()
        const duration = performance.now() - startTime
        
        // Record metrics
        recordAPICall('/api/v1/products', duration, true)
        recordReactQueryMetric('products_list', duration, false)
        
        return result
      } catch (error) {
        const duration = performance.now() - startTime
        recordAPICall('/api/v1/products', duration, false)
        throw error
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
  })
}

// Hook to fetch a single product by ID
export function useProduct(id: string | number) {
  return useQuery({
    queryKey: productKeys.detail(id),
    queryFn: async () => {
      const startTime = performance.now()
      
      try {
        const result = await catalogApi.getProduct(id)
        const duration = performance.now() - startTime
        
        // Record metrics
        recordAPICall(`/api/v1/products/${id}`, duration, true)
        recordReactQueryMetric(`product_detail_${id}`, duration, false)
        
        return result
      } catch (error) {
        const duration = performance.now() - startTime
        recordAPICall(`/api/v1/products/${id}`, duration, false)
        throw error
      }
    },
    enabled: !!id, // Only run query if ID is provided
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
  })
}

// Future: Hook for product mutations (create, update, delete)
// export function useCreateProduct() {
//   const queryClient = useQueryClient()
//   
//   return useMutation({
//     mutationFn: catalogApi.createProduct,
//     onSuccess: () => {
//       // Invalidate and refetch products list
//       queryClient.invalidateQueries({ queryKey: productKeys.lists() })
//     },
//   })
// } 