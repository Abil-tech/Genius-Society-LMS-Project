import { NavLink, useNavigate } from 'react-router-dom'
import {
  LayoutDashboard,
  DoorOpen,
  BookOpen,
  FileSpreadsheet,
  CheckSquare,
  FolderGit2,
  Award,
  CalendarDays,
  Users2,
  GraduationCap,
  LogOut,
  X,
  ShieldCheck,
} from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

interface TeacherSidebarProps {
  mobileOpen: boolean
  setMobileOpen: (open: boolean) => void
  isWalasOverride?: boolean
}

interface NavItem {
  to: string
  label: string
  icon: typeof LayoutDashboard
  end?: boolean
}

const mainNavItems: NavItem[] = [
  { to: '/teacher', label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: '/teacher/classes', label: 'Kelas Diajar', icon: DoorOpen },
  { to: '/teacher/materials', label: 'Materi Belajar', icon: BookOpen },
  { to: '/teacher/assignments', label: 'Tugas Siswa', icon: FileSpreadsheet },
  { to: '/teacher/assessments', label: 'Assessment', icon: CheckSquare },
  { to: '/teacher/projects', label: 'Projek Siswa', icon: FolderGit2 },
  { to: '/teacher/grades', label: 'Nilai & Rekap', icon: Award },
  { to: '/teacher/schedule', label: 'Jadwal Mengajar', icon: CalendarDays },
]

const walasNavItems: NavItem[] = [
  { to: '/teacher/walas/students', label: 'Siswa Perwalian', icon: Users2 },
  { to: '/teacher/walas/grades', label: 'Rekap Nilai Kelas', icon: GraduationCap },
]

export default function TeacherSidebar({
  mobileOpen,
  setMobileOpen,
  isWalasOverride,
}: TeacherSidebarProps) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  // Guru isWalas can come from user object directly — no need for override prop
  const isWalas =
    isWalasOverride !== undefined
      ? isWalasOverride
      : Boolean(user?.isWalas ?? user?.is_walas)

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
              <span className="block text-xs font-black tracking-tight text-[#30251D]">
                GENIUS SOCIETY
              </span>
              <span className="block text-[10px] font-semibold text-[#FF9800]">
                {isWalas ? 'GURU & WALI KELAS' : 'PORTAL GURU'}
              </span>
            </div>
          </div>

          {/* Close button on mobile */}
          <button
            type="button"
            onClick={() => setMobileOpen(false)}
            aria-label="Tutup menu navigasi"
            className="lg:hidden text-gray-400 hover:text-gray-700 p-1 rounded-lg"
          >
            <X size={18} />
          </button>
        </div>

        {/* Navigation Section */}
        <div className="mt-4 px-3 space-y-6">
          {/* Menu Akademik */}
          <div>
            <p className="px-3 text-[10px] font-bold uppercase tracking-wider text-[#687385]">
              Pembelajaran
            </p>
            <nav aria-label="Menu Pembelajaran" className="mt-2 space-y-1">
              {mainNavItems.map((item) => {
                const Icon = item.icon
                return (
                  <NavLink
                    key={item.to}
                    to={item.to}
                    end={item.end}
                    onClick={() => setMobileOpen(false)}
                    className={({ isActive }) =>
                      `flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-xs font-semibold transition-all ${
                        isActive
                          ? 'bg-[#FEEADC] text-[#FF9800] font-bold shadow-xs'
                          : 'text-[#30251D]/80 hover:bg-[#FFF8F5] hover:text-[#30251D]'
                      }`
                    }
                  >
                    <Icon size={16} className="shrink-0" />
                    <span>{item.label}</span>
                  </NavLink>
                )
              })}
            </nav>
          </div>

          {/* Section Wali Kelas (Hanya jika isWalas) */}
          {isWalas && (
            <div className="pt-2 border-t border-[#EBD4C5]/40">
              <div className="flex items-center gap-1.5 px-3">
                <ShieldCheck size={12} className="text-[#FF9800]" />
                <p className="text-[10px] font-bold uppercase tracking-wider text-[#FF9800]">
                  Wali Kelas
                </p>
              </div>
              <nav aria-label="Menu Wali Kelas" className="mt-2 space-y-1">
                {walasNavItems.map((item) => {
                  const Icon = item.icon
                  return (
                    <NavLink
                      key={item.to}
                      to={item.to}
                      onClick={() => setMobileOpen(false)}
                      className={({ isActive }) =>
                        `flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-xs font-semibold transition-all ${
                          isActive
                            ? 'bg-[#FEEADC] text-[#FF9800] font-bold shadow-xs'
                            : 'text-[#30251D]/80 hover:bg-[#FFF8F5] hover:text-[#30251D]'
                        }`
                      }
                    >
                      <Icon size={16} className="shrink-0" />
                      <span>{item.label}</span>
                    </NavLink>
                  )
                })}
              </nav>
            </div>
          )}
        </div>
      </div>

      {/* User profile & logout footer */}
      <div className="p-3 border-t border-[#EBD4C5]/60 bg-[#FFFFFF]">
        <div className="flex items-center justify-between rounded-xl bg-[#FFF8F5] p-3 border border-[#EBD4C5]/50">
          <div className="min-w-0 pr-2">
            <p className="truncate text-xs font-bold text-[#30251D]">
              {user?.name || 'Tenaga Pengajar'}
            </p>
            <p className="truncate text-[10px] text-[#687385]">
              {user?.email || 'guru@sekolah.id'}
            </p>
          </div>
          <button
            type="button"
            onClick={handleLogout}
            title="Keluar dari akun"
            aria-label="Logout"
            className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-rose-50 hover:text-rose-600 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-500 transition-colors"
          >
            <LogOut size={16} />
          </button>
        </div>
      </div>
    </div>
  )

  return (
    <>
      {/* Desktop Fixed Sidebar (w-64) */}
      <aside className="fixed inset-y-0 left-0 z-30 hidden lg:block w-64">
        {sidebarContent}
      </aside>

      {/* Mobile Drawer Backdrop */}
      {mobileOpen && (
        <div
          role="presentation"
          onClick={() => setMobileOpen(false)}
          className="fixed inset-0 z-40 bg-black/40 backdrop-blur-xs lg:hidden transition-opacity"
        />
      )}

      {/* Mobile Drawer (w-72) */}
      <aside
        className={`fixed inset-y-0 left-0 z-50 w-72 transform transition-transform duration-300 ease-in-out lg:hidden ${
          mobileOpen ? 'translate-x-0' : '-translate-x-full'
        }`}
      >
        {sidebarContent}
      </aside>
    </>
  )
}
