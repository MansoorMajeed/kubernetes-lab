import { useQuery } from '@tanstack/react-query'
import { catalogApi } from '../services/catalogApi'
import type { Product } from '../types/api'

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
    queryFn: () => catalogApi.getProducts(),
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
  })
}

// Hook to fetch a single product by ID
export function useProduct(id: string | number) {
  return useQuery({
    queryKey: productKeys.detail(id),
    queryFn: () => catalogApi.getProduct(id),
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