import { Menu, Bell } from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

interface AdminTopbarProps {
  onOpenMobileSidebar: () => void
}

export default function AdminTopbar({ onOpenMobileSidebar }: AdminTopbarProps) {
  const { user } = useAuth()

  return (
    <header className="sticky top-0 z-20 flex h-16 w-full items-center justify-between border-b border-[#EBD4C5] bg-[#FFFFFF]/90 px-4 sm:px-8 backdrop-blur">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onOpenMobileSidebar}
          className="rounded-lg p-2 text-gray-600 hover:bg-[#FFF8F5] md:hidden"
          aria-label="Buka menu navigasi"
        >
          <Menu size={20} />
        </button>
        <span className="text-xs font-semibold uppercase tracking-wider text-gray-400 hidden sm:inline-block">
          Sistem Administrasi Akademik Sekolah
        </span>
      </div>

      <div className="flex items-center gap-4">
        {/* Notifikasi Bell */}
        <button
          type="button"
          className="relative flex h-9 w-9 items-center justify-center rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 transition-colors hover:text-[#FF9800]"
          title="Notifikasi"
        >
          <Bell size={16} />
          <span className="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-[#FF9800]" />
        </button>

        {/* User Badge */}
        <div className="flex items-center gap-2.5 pl-2 border-l border-[#EBD4C5]/70">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#FF9800] text-xs font-bold text-white shadow-sm">
            {user?.name?.charAt(0) || 'A'}
          </div>
          <div className="hidden sm:block text-left">
            <span className="block text-xs font-bold text-[#30251D]">
              {user?.name || 'Administrator'}
            </span>
            <span className="block text-[10px] font-semibold text-[#FF9800]">
              Admin Sekolah
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
