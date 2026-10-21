import { useState, type ReactNode } from 'react'
import AdminSidebar from './AdminSidebar'
import AdminTopbar from './AdminTopbar'

interface AdminLayoutProps {
  children: ReactNode
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <div className="min-h-screen bg-[#FFF8F5] text-[#30251D]">
      {/* Fixed Sidebar */}
      <AdminSidebar mobileOpen={mobileOpen} setMobileOpen={setMobileOpen} />

      {/* Main Content Area: offset 64 (w-64) on md/desktop agar tidak tertutup sidebar */}
      <div className="flex flex-col md:pl-64 min-h-screen">
        <AdminTopbar onOpenMobileSidebar={() => setMobileOpen(true)} />

        <main className="flex-1 p-4 sm:p-6 md:p-8">
          {children}
        </main>
      </div>
    </div>
  )
}
