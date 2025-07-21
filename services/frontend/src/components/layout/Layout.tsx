import type { ReactNode } from 'react'
import Header from './Header'
import Footer from './Footer'

interface LayoutProps {
  children: ReactNode
  title?: string
  showBackButton?: boolean
  footerInfo?: string
}

export default function Layout({ 
  children, 
  title, 
  showBackButton = false,
  footerInfo 
}: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <Header title={title} showBackButton={showBackButton} />
      
      <main className="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 w-full">
        {children}
      </main>
      
      <Footer extraInfo={footerInfo} />
    </div>
  )
} 