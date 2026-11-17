import { Menu, Bell } from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

interface TeacherTopbarProps {
  onOpenMobileSidebar: () => void
  isWalas?: boolean
}

export default function TeacherTopbar({
  onOpenMobileSidebar,
  isWalas,
}: TeacherTopbarProps) {
  const { user } = useAuth()

  const displayName = user?.name || 'Bapak/Ibu Guru'
  const initial = displayName.charAt(0).toUpperCase()

  return (
    <header className="sticky top-0 z-20 flex h-16 w-full items-center justify-between border-b border-[#EBD4C5] bg-[#FFFFFF]/90 px-4 sm:px-8 backdrop-blur">
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onOpenMobileSidebar}
          className="rounded-lg p-2 text-[#30251D] hover:bg-[#FFF8F5] lg:hidden focus:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9800]"
          aria-label="Buka menu navigasi samping"
        >
          <Menu size={20} />
        </button>
        <div className="hidden sm:block">
          <span className="text-xs font-semibold uppercase tracking-wider text-[#687385]">
            Portal Pembelajaran Guru
          </span>
        </div>
      </div>

      <div className="flex items-center gap-4">
        {/* Notifikasi Bell */}
        <button
          type="button"
          className="relative flex h-9 w-9 items-center justify-center rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] text-[#30251D] transition-colors hover:text-[#FF9800] focus:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9800]"
          title="Notifikasi"
          aria-label="Lihat notifikasi"
        >
          <Bell size={16} />
          <span className="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-[#FF9800]" />
        </button>

        {/* User Badge */}
        <div className="flex items-center gap-2.5 pl-2 border-l border-[#EBD4C5]/70">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#FF9800] text-xs font-bold text-white shadow-xs">
            {initial}
          </div>
          <div className="hidden sm:block text-left">
            <span className="block text-xs font-bold text-[#30251D] truncate max-w-[150px]">
              {displayName}
            </span>
            <span className="block text-[10px] font-semibold text-[#FF9800]">
              {isWalas ? 'Guru & Wali Kelas' : 'Guru Pengajar'}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
