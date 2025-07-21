import { ApiError } from '../../services/catalogApi'

interface ErrorDisplayProps {
  error: Error | ApiError | null
  title?: string
  onRetry?: () => void
  className?: string
}

export default function ErrorDisplay({ 
  error, 
  title = 'Something went wrong',
  onRetry,
  className = ''
}: ErrorDisplayProps) {
  if (!error) return null

  const isApiError = error instanceof ApiError
  const errorMessage = error.message
  const statusCode = isApiError ? error.status : undefined

  return (
    <div className={`bg-red-50 border border-red-200 rounded-lg p-6 ${className}`}>
      <div className="flex">
        <div className="flex-shrink-0">
          <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
          </svg>
        </div>
        <div className="ml-3 flex-1">
          <h3 className="text-lg font-medium text-red-800">{title}</h3>
          <div className="mt-2 text-sm text-red-700">
            <p>{errorMessage}</p>
            {statusCode && (
              <p className="mt-1 text-xs text-red-600">Status: {statusCode}</p>
            )}
            {isApiError && statusCode === 500 && (
              <p className="mt-2 text-xs text-red-500">
                Make sure the catalog service is running: <code>./tilt-lab up</code>
              </p>
            )}
          </div>
          {onRetry && (
            <div className="mt-4">
              <button
                onClick={onRetry}
                className="bg-red-100 text-red-800 px-4 py-2 rounded-md text-sm font-medium hover:bg-red-200 transition-colors"
              >
                Try Again
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
} 