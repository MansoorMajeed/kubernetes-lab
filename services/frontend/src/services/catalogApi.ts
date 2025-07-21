import type { Product, ProductsResponse, ProductResponse } from '../types/api'

// Custom error class for API errors
export class ApiError extends Error {
  public status?: number
  
  constructor(message: string, status?: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

// Base API configuration
const API_BASE_URL = '/api/v1'

// Generic fetch wrapper with error handling
async function apiRequest<T>(endpoint: string): Promise<T> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`)
    
    if (!response.ok) {
      throw new ApiError(`Failed to fetch: ${response.status} ${response.statusText}`, response.status)
    }
    
    return response.json()
  } catch (error) {
    if (error instanceof ApiError) {
      throw error
    }
    
    // Handle network errors, JSON parsing errors, etc.
    throw new ApiError(
      error instanceof Error ? error.message : 'An unexpected error occurred'
    )
  }
}

// Catalog API service
export const catalogApi = {
  // Get all products with pagination
  getProducts: (): Promise<ProductsResponse> => {
    return apiRequest<ProductsResponse>('/products')
  },

  // Get single product by ID
  getProduct: (id: string | number): Promise<ProductResponse> => {
    return apiRequest<ProductResponse>(`/products/${id}`)
  },

  // Future: Add products, update products, etc.
  // createProduct: (product: Omit<Product, 'id' | 'created_at' | 'updated_at'>): Promise<ProductResponse> => {
  //   return apiRequest<ProductResponse>('/products', {
  //     method: 'POST',
  //     headers: { 'Content-Type': 'application/json' },
  //     body: JSON.stringify(product)
  //   })
  // }
} 