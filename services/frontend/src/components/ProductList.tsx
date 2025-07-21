import { useEffect } from 'react'
import { useProducts } from '../hooks/useProducts'
import Layout from './layout/Layout'
import Loading from './ui/Loading'
import ErrorDisplay from './ui/ErrorDisplay'
import ProductCard from './ui/ProductCard'
import { recordPageView } from '../services/metricsApi'

export default function ProductList() {
  const { data, isLoading, error, refetch } = useProducts()

  // Record page view on component mount
  useEffect(() => {
    recordPageView('product_list')
  }, [])

  return (
    <Layout title="Products">
      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Products</h2>
        <p className="text-gray-600">Browse our catalog powered by the Catalog Service</p>
      </div>

      {/* Loading State */}
      {isLoading && (
        <Loading message="Loading products..." className="h-64" />
      )}

      {/* Error State */}
      {error && (
        <ErrorDisplay 
          error={error} 
          title="Error loading products"
          onRetry={() => refetch()}
        />
      )}

      {/* Products Grid */}
      {!isLoading && !error && data && (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {data.data.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}

      {/* Empty State */}
      {!isLoading && !error && data && data.data.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 mb-4">
            <svg className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No products found</h3>
          <p className="text-gray-600">The catalog service returned no products.</p>
          <button
            onClick={() => refetch()}
            className="mt-4 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
          >
            Refresh
          </button>
        </div>
      )}
    </Layout>
  )
} 