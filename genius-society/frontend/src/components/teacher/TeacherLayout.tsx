import { useState, type ReactNode } from 'react'
import TeacherSidebar from './TeacherSidebar'
import TeacherTopbar from './TeacherTopbar'

interface TeacherLayoutProps {
  children: ReactNode
  isWalas?: boolean
}

export default function TeacherLayout({ children, isWalas }: TeacherLayoutProps) {
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[#FFF8F5] text-[#30251D]">
      {/* Fixed Sidebar for desktop, drawer for mobile */}
      <TeacherSidebar
        mobileOpen={mobileOpen}
        setMobileOpen={setMobileOpen}
        isWalasOverride={isWalas}
      />

      {/* Main Content Area: offset on lg:pl-64 */}
      <div className="flex flex-col lg:pl-64 min-h-screen">
        <TeacherTopbar
          onOpenMobileSidebar={() => setMobileOpen(true)}
          isWalas={isWalas}
        />

        <main className="flex-1 p-4 sm:p-6 lg:p-8 max-w-7xl w-full mx-auto">
          {children}
        </main>
      </div>
    </div>
  )
}
