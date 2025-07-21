// API Response Types
export interface ApiResponse<T> {
  data: T
}

export interface PaginatedApiResponse<T> {
  data: T[]
  page: number
  limit: number
  count: number
}

// Product Types
export interface Product {
  id: number
  name: string
  description: string
  price: number
  stock_quantity: number
  created_at: string
  updated_at: string
}

// API Endpoint Response Types
export type ProductsResponse = PaginatedApiResponse<Product>
export type ProductResponse = ApiResponse<Product> 