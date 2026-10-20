import { useAuth } from '../contexts/AuthContext'
import { LogOut, UserCircle } from 'lucide-react'
import { useNavigate } from 'react-router-dom'

interface RoleDashboardProps {
  roleTitle: string
  badgeColor?: string
}

export default function RoleDashboardPlaceholder({ roleTitle, badgeColor = 'bg-brand-peach text-brand-orange' }: RoleDashboardProps) {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  return (
    <div className="min-h-screen bg-[#FFF8F5] text-[#30251D]">
      {/* Top Header */}
      <header className="sticky top-0 z-20 flex h-16 items-center justify-between border-b border-[#EBD4C5] bg-white px-6 shadow-sm">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-[#FF9800] font-bold text-white shadow-sm">
            GS
          </div>
          <div>
            <h1 className="text-sm font-bold text-[#30251D]">Genius Society LMS</h1>
            <span className={`inline-block rounded-full px-2 py-0.5 text-[10px] font-semibold ${badgeColor}`}>
              {roleTitle}
            </span>
          </div>
        </div>

        <button
          onClick={handleLogout}
          className="flex items-center gap-2 rounded-lg border border-[#EBD4C5] px-3.5 py-1.5 text-xs font-semibold text-gray-700 transition-colors hover:bg-red-50 hover:text-red-600"
        >
          <LogOut size={14} />
          <span>Keluar</span>
        </button>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-4xl p-6 md:p-10">
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-8 shadow-sm">
          <div className="flex items-center gap-4">
            <span className="flex h-14 w-14 items-center justify-center rounded-2xl bg-orange-50 text-[#FF9800]">
              <UserCircle size={36} />
            </span>
            <div>
              <p className="text-xs font-medium uppercase tracking-wider text-gray-400">
                Dashboard {roleTitle}
              </p>
              <h2 className="text-2xl font-black text-[#30251D]">
                Selamat datang {user?.name}
              </h2>
            </div>
          </div>

          <div className="mt-6 rounded-xl border border-dashed border-[#EBD4C5] bg-[#FFF8F5] p-5">
            <h3 className="text-xs font-bold uppercase tracking-wider text-gray-500">
              Informasi Akun
            </h3>
            <div className="mt-3 grid gap-3 sm:grid-cols-3 text-xs">
              <div>
                <span className="text-gray-400">Username:</span>
                <p className="font-mono font-semibold text-[#30251D]">{user?.username}</p>
              </div>
              <div>
                <span className="text-gray-400">Email:</span>
                <p className="font-semibold text-[#30251D]">{user?.email}</p>
              </div>
              <div>
                <span className="text-gray-400">Role Terdaftar:</span>
                <p className="font-semibold capitalize text-[#FF9800]">{user?.role?.replace('_', ' ')}</p>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
