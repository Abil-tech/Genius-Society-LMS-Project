import { NavLink, useNavigate } from 'react-router-dom'
import {
  LayoutDashboard,
  BookOpen,
  FileSpreadsheet,
  CheckSquare,
  FolderGit2,
  Award,
  CalendarDays,
  Trophy,
  LogOut,
  X,
  GraduationCap,
} from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

interface StudentSidebarProps {
  mobileOpen: boolean
  setMobileOpen: (open: boolean) => void
}

const navItems = [
  { to: '/student', label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: '/student/materials', label: 'Materi Belajar', icon: BookOpen },
  { to: '/student/assignments', label: 'Tugas', icon: FileSpreadsheet },
  { to: '/student/assessments', label: 'Assessment', icon: CheckSquare },
  { to: '/student/projects', label: 'Projek', icon: FolderGit2 },
  { to: '/student/grades', label: 'Nilai Saya', icon: Award },
  { to: '/student/schedule', label: 'Jadwal Pelajaran', icon: CalendarDays },
  { to: '/student/ranking', label: 'Ranking Kelas', icon: Trophy },
]

export default function StudentSidebar({ mobileOpen, setMobileOpen }: StudentSidebarProps) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    try {
      await logout()
    } finally {
      navigate('/login')
    }
  }

  const sidebarContent = (
    <div className="flex h-full flex-col justify-between bg-[#FFFFFF] border-r border-[#EBD4C5]">
      {/* Brand & Header */}
      <div className="overflow-y-auto">
        <div className="flex h-16 items-center justify-between px-6 border-b border-[#EBD4C5]/60">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-[#FF9800] font-black text-white shadow-sm">
              GS
            </div>
            <div>
              <span className="block text-xs font-black tracking-wider text-[#30251D]">
                GENIUS SOCIETY
              </span>
              <span className="block text-[10px] font-bold tracking-widest uppercase text-[#FF9800]">
                PORTAL SISWA
              </span>
            </div>
          </div>
          {/* Close button on mobile */}
          <button
            type="button"
            onClick={() => setMobileOpen(false)}
            className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] lg:hidden"
            aria-label="Tutup navigasi"
          >
            <X size={18} />
          </button>
        </div>

        {/* Navigation list */}
        <nav className="space-y-1.5 p-4">
          <span className="px-3 text-[10px] font-bold uppercase tracking-wider text-[#687385]">
            Menu Pembelajaran
          </span>
          <div className="mt-2 space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon
              return (
                <NavLink
                  key={item.to}
                  to={item.to}
                  end={item.end}
                  onClick={() => setMobileOpen(false)}
                  className={({ isActive }) =>
                    `flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-xs font-bold transition-all ${
                      isActive
                        ? 'bg-[#FEEADC] text-[#FF9800] shadow-xs'
                        : 'text-[#687385] hover:bg-[#FFF8F5] hover:text-[#30251D]'
                    }`
                  }
                >
                  <Icon size={16} />
                  <span>{item.label}</span>
                </NavLink>
              )
            })}
          </div>
        </nav>
      </div>

      {/* Footer Profile & Logout */}
      <div className="p-4 border-t border-[#EBD4C5]/60">
        <div className="flex items-center gap-3 rounded-xl bg-[#FFF8F5] p-3 border border-[#EBD4C5]">
          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-[#FF9800] text-xs font-bold text-white shadow-xs">
            <GraduationCap size={16} />
          </div>
          <div className="min-w-0 flex-1">
            <p className="truncate text-xs font-bold text-[#30251D]">
              {user?.name || 'Siswa'}
            </p>
            <p className="truncate text-[10px] text-[#687385] font-semibold">
              {user?.username || 'Murid'}
            </p>
          </div>
          <button
            type="button"
            onClick={handleLogout}
            className="rounded-lg p-1.5 text-gray-400 hover:bg-white hover:text-rose-600 transition-colors"
            title="Keluar"
            aria-label="Keluar dari akun"
          >
            <LogOut size={16} />
          </button>
        </div>
      </div>
    </div>
  )

  return (
    <>
      {/* Desktop Fixed Sidebar */}
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-64 lg:block">
        {sidebarContent}
      </aside>

      {/* Mobile Drawer with Backdrop */}
      {mobileOpen && (
        <div className="fixed inset-0 z-40 lg:hidden">
          <div
            className="fixed inset-0 bg-black/40 backdrop-blur-xs transition-opacity"
            onClick={() => setMobileOpen(false)}
          />
          <div className="fixed inset-y-0 left-0 w-64 shadow-xl">
            {sidebarContent}
          </div>
        </div>
      )}
    </>
  )
}
