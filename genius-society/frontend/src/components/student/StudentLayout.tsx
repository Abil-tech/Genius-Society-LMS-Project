import { useState, type ReactNode } from 'react'
import StudentSidebar from './StudentSidebar'
import StudentTopbar from './StudentTopbar'

export default function StudentLayout({ children }: { children: ReactNode }) {
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[#FFF8F5] text-[#30251D]">
      <StudentSidebar mobileOpen={mobileOpen} setMobileOpen={setMobileOpen} />
      <div className="flex flex-col lg:pl-64 min-h-screen">
        <StudentTopbar onOpenMobileSidebar={() => setMobileOpen(true)} />
        <main className="flex-1 p-4 sm:p-6 lg:p-8 max-w-7xl w-full mx-auto">
          {children}
        </main>
      </div>
    </div>
  )
}
