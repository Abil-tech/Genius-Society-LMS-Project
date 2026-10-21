import { NavLink, useNavigate } from 'react-router-dom'
import { 
  LayoutDashboard, 
  Users, 
  GraduationCap, 
  DoorOpen, 
  BookOpen, 
  CalendarDays, 
  UserCog, 
  LogOut,
  X
} from 'lucide-react'
import { useAuth } from '../../contexts/AuthContext'

interface AdminSidebarProps {
  mobileOpen: boolean
  setMobileOpen: (open: boolean) => void
}

const navItems = [
  { to: '/admin', label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: '/admin/students', label: 'Data Siswa', icon: GraduationCap },
  { to: '/admin/teachers', label: 'Data Guru', icon: Users },
  { to: '/admin/classes', label: 'Data Kelas', icon: DoorOpen },
  { to: '/admin/subjects', label: 'Mata Pelajaran', icon: BookOpen },
  { to: '/admin/academic-year', label: 'Tahun Ajaran', icon: CalendarDays },
  { to: '/admin/accounts', label: 'Manajemen Akun', icon: UserCog },
]

export default function AdminSidebar({ mobileOpen, setMobileOpen }: AdminSidebarProps) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const sidebarContent = (
    <div className="flex h-full flex-col justify-between bg-[#FFFFFF] border-r border-[#EBD4C5]">
      {/* Brand & School Header */}
      <div>
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
                ADMINISTRATOR
              </span>
            </div>
          </div>

          {/* Close button on mobile */}
          <button 
            type="button"
            onClick={() => setMobileOpen(false)}
            className="md:hidden text-gray-400 hover:text-gray-700 p-1"
          >
            <X size={18} />
          </button>
        </div>

        {/* Navigation list */}
        <nav aria-label="Menu Admin" className="mt-4 px-3 space-y-1">
          {navItems.map((item) => {
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
                      ? 'bg-[#FEEADC] text-[#FF9800] font-bold shadow-sm'
                      : 'text-gray-600 hover:bg-[#FFF8F5] hover:text-[#30251D]'
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

      {/* User profile & logout footer */}
      <div className="p-3 border-t border-[#EBD4C5]/60">
        <div className="flex items-center justify-between rounded-xl bg-[#FFF8F5] p-3 border border-[#EBD4C5]/50">
          <div className="min-w-0 pr-2">
            <p className="truncate text-xs font-bold text-[#30251D]">
              {user?.name || 'Administrator'}
            </p>
            <p className="truncate text-[10px] text-gray-500 font-mono">
              @{user?.username || 'admin'}
            </p>
          </div>
          <button
            type="button"
            onClick={handleLogout}
            title="Keluar dari sistem"
            className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border border-[#EBD4C5] bg-white text-gray-600 transition-colors hover:bg-red-50 hover:text-red-600 hover:border-red-200"
          >
            <LogOut size={14} />
          </button>
        </div>
      </div>
    </div>
  )

  return (
    <>
      {/* Desktop Fixed Sidebar (tetap di posisinya saat halaman di-scroll) */}
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-64 md:block">
        {sidebarContent}
      </aside>

      {/* Mobile Drawer */}
      {mobileOpen && (
        <div className="fixed inset-0 z-50 md:hidden">
          {/* Backdrop */}
          <div 
            className="fixed inset-0 bg-black/40 backdrop-blur-sm transition-opacity"
            onClick={() => setMobileOpen(false)}
          />
          {/* Drawer content */}
          <div className="fixed inset-y-0 left-0 w-72 max-w-[85vw] shadow-2xl">
            {sidebarContent}
          </div>
        </div>
      )}
    </>
  )
}
