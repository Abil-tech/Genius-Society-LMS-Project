import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

interface ProtectedRouteProps {
  allowedRoles?: string[]
}

export default function ProtectedRoute({ allowedRoles }: ProtectedRouteProps) {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-[#F1F1F3]">
        <div className="flex items-center gap-3 rounded-lg border border-line bg-surface px-6 py-4 shadow-sm">
          <div className="h-5 w-5 animate-spin rounded-full border-2 border-brand-orange border-t-transparent" />
          <span className="text-sm font-medium text-brand-navy">Memuat sesi pengguna...</span>
        </div>
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (allowedRoles && !allowedRoles.includes(user.role)) {
    // Redirect ke dashboard role-nya sendiri jika mencoba akses role lain
    const roleRoutes: Record<string, string> = {
      super_admin: '/super-admin',
      admin: '/admin',
      teacher: '/teacher',
      student: '/student',
      curriculum: '/curriculum',
      principal: '/principal',
    }
    const target = roleRoutes[user.role] || '/'
    return <Navigate to={target} replace />
  }

  return <Outlet />
}
