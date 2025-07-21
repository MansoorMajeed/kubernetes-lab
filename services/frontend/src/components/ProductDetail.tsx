import { useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useProduct } from '../hooks/useProducts'
import Layout from './layout/Layout'
import Loading from './ui/Loading'
import ErrorDisplay from './ui/ErrorDisplay'
import { recordPageView, recordProductView } from '../services/metricsApi'

export default function ProductDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { data, isLoading, error, refetch } = useProduct(id!)

  // Record page view and product view when data loads
  useEffect(() => {
    recordPageView('product_detail')
  }, [])

  useEffect(() => {
    if (data?.data) {
      recordProductView(data.data.id, data.data.name)
    }
  }, [data])

  // Loading state
  if (isLoading) {
    return (
      <Layout showBackButton>
        <Loading message="Loading product..." className="h-64" />
      </Layout>
    )
  }

  // Error state
  if (error || !data) {
    return (
      <Layout showBackButton>
        <ErrorDisplay 
          error={error} 
          title="Error loading product"
          onRetry={() => refetch()}
        />
      </Layout>
    )
  }

  const product = data.data

  return (
    <Layout 
      showBackButton 
      title="Product Detail"
      footerInfo={`Product ID: ${product.id}`}
    >
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="p-8">
          {/* Product Header */}
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-4">{product.name}</h1>
            <div className="flex items-center justify-between">
              <span className="text-4xl font-bold text-blue-600">
                ${product.price.toFixed(2)}
              </span>
              <div className="text-right">
                <div className={`inline-flex px-3 py-1 rounded-full text-sm font-medium ${
                  product.stock_quantity > 10 
                    ? 'bg-green-100 text-green-800' 
                    : product.stock_quantity > 0 
                    ? 'bg-yellow-100 text-yellow-800'
                    : 'bg-red-100 text-red-800'
                }`}>
                  {product.stock_quantity > 10 
                    ? `${product.stock_quantity} in stock` 
                    : product.stock_quantity > 0 
                    ? `Only ${product.stock_quantity} left!`
                    : 'Out of stock'
                  }
                </div>
              </div>
            </div>
          </div>

          {/* Product Description */}
          <div className="mb-8">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Description</h2>
            <p className="text-gray-600 text-lg leading-relaxed">{product.description}</p>
          </div>

          {/* Product Details Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
            <div className="bg-gray-50 p-4 rounded-lg">
              <h3 className="font-semibold text-gray-900 mb-2">Product ID</h3>
              <p className="text-gray-600">{product.id}</p>
            </div>
            <div className="bg-gray-50 p-4 rounded-lg">
              <h3 className="font-semibold text-gray-900 mb-2">Stock Quantity</h3>
              <p className="text-gray-600">{product.stock_quantity} units</p>
            </div>
            <div className="bg-gray-50 p-4 rounded-lg">
              <h3 className="font-semibold text-gray-900 mb-2">Added</h3>
              <p className="text-gray-600">{new Date(product.created_at).toLocaleDateString()}</p>
            </div>
            <div className="bg-gray-50 p-4 rounded-lg">
              <h3 className="font-semibold text-gray-900 mb-2">Last Updated</h3>
              <p className="text-gray-600">{new Date(product.updated_at).toLocaleDateString()}</p>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-4">
            <button
              className={`flex-1 py-3 px-6 rounded-md font-medium transition-colors ${
                product.stock_quantity > 0
                  ? 'bg-blue-600 text-white hover:bg-blue-700'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              }`}
              disabled={product.stock_quantity === 0}
            >
              {product.stock_quantity > 0 ? 'Add to Cart (Coming in Phase 4.0)' : 'Out of Stock'}
            </button>
            <button
              onClick={() => navigate('/')}
              className="flex-1 py-3 px-6 bg-gray-100 text-gray-700 rounded-md font-medium hover:bg-gray-200 transition-colors"
            >
              Continue Shopping
            </button>
          </div>
        </div>
      </div>
    </Layout>
  )
} 