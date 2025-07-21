interface FooterProps {
  extraInfo?: string
}

export default function Footer({ extraInfo }: FooterProps) {
  return (
    <footer className="bg-white border-t border-gray-200 mt-16">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="text-center text-sm text-gray-500">
          <p>LocalMart - Phase 3.0.0 Frontend • React + React Query + TypeScript</p>
          <div className="mt-1 flex justify-center items-center space-x-4">
            <span>
              API: <code className="bg-gray-100 px-2 py-1 rounded">catalog.kubelab.lan:8081</code>
            </span>
            {extraInfo && (
              <>
                <span>•</span>
                <span>{extraInfo}</span>
              </>
            )}
          </div>
        </div>
      </div>
    </footer>
  )
} 