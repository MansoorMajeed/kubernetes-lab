import { useNavigate } from 'react-router-dom'

interface HeaderProps {
  showBackButton?: boolean
  title?: string
}

export default function Header({ showBackButton = false, title }: HeaderProps) {
  const navigate = useNavigate()

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <h1 className="text-2xl font-bold text-gray-900">🏪 LocalMart</h1>
            {title && (
              <span className="ml-4 text-lg text-gray-600">• {title}</span>
            )}
          </div>
          
          <div className="flex items-center space-x-4">
            {showBackButton && (
              <button
                onClick={() => navigate('/')}
                className="text-blue-600 hover:text-blue-700 font-medium flex items-center transition-colors"
              >
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                </svg>
                Back to Products
              </button>
            )}
            
            <div className="text-sm text-gray-600">
              Phase 3.0.0 - React + React Query
            </div>
          </div>
        </div>
      </div>
    </header>
  )
} 