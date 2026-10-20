import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ShieldCheck, Mail, Lock, AlertCircle, Loader2, ArrowLeft } from 'lucide-react'
import { useAuth } from '../contexts/AuthContext'

export default function AdminLoginPage() {
  const [loginId, setLoginId] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setSubmitting(true)

    try {
      const user = await login(loginId, password)
      
      // Validasi khusus bahwa hanya role admin/super admin yang diproses di portal ini
      if (user.role !== 'super_admin' && user.role !== 'admin') {
        setError('Akses ditolak: Portal ini khusus Administrator dan Super Admin.')
        return
      }

      if (user.role === 'super_admin') {
        navigate('/super-admin')
      } else {
        navigate('/admin')
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('Kredensial salah atau terjadi kesalahan pada sistem.')
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#F1F1F3] px-4 font-sans text-brand-navy">
      <div className="w-full max-w-[420px] rounded-2xl border border-gray-200/90 bg-white p-8 shadow-card sm:p-10">
        <div className="mb-6">
          <Link
            to="/login"
            className="mb-4 inline-flex items-center gap-1.5 text-xs text-muted hover:text-brand-navy transition-colors"
          >
            <ArrowLeft size={14} />
            <span>Kembali ke Login Pengguna</span>
          </Link>

          <div className="flex items-center gap-2.5">
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-peach text-brand-orange">
              <ShieldCheck size={20} />
            </span>
            <div>
              <span className="font-mono text-[10px] font-bold tracking-wider text-brand-orange uppercase">
                PORTAL ADMINISTRASI
              </span>
              <h1 className="text-xl font-extrabold text-brand-navy">
                Genius Society
              </h1>
            </div>
          </div>
          <p className="mt-3 text-xs leading-relaxed text-muted">
            Masuk dengan akun terverifikasi Administrator Sekolah atau Super Admin.
          </p>
        </div>

        <div className="mb-6 border-b border-line" />

        {error && (
          <div className="mb-5 flex items-start gap-2 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-600">
            <AlertCircle size={15} className="mt-0.5 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1.5 block text-xs font-semibold text-gray-700">
              Email atau Username
            </label>
            <div className="relative flex items-center">
              <span className="pointer-events-none absolute left-3.5 text-gray-400">
                <Mail size={15} />
              </span>
              <input
                type="text"
                required
                value={loginId}
                onChange={(e) => setLoginId(e.target.value)}
                placeholder="admin@genius.ac.id"
                className="w-full rounded-lg border border-line bg-gray-50/50 py-2.5 pl-10 pr-3 text-xs text-gray-800 placeholder-gray-400 transition-colors focus:border-brand-orange focus:bg-white focus:outline-none focus:ring-1 focus:ring-brand-orange"
              />
            </div>
          </div>

          <div>
            <label className="mb-1.5 block text-xs font-semibold text-gray-700">
              Password
            </label>
            <div className="relative flex items-center">
              <span className="pointer-events-none absolute left-3.5 text-gray-400">
                <Lock size={15} />
              </span>
              <input
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••••••"
                className="w-full rounded-lg border border-line bg-gray-50/50 py-2.5 pl-10 pr-3 text-xs text-gray-800 placeholder-gray-400 transition-colors focus:border-brand-orange focus:bg-white focus:outline-none focus:ring-1 focus:ring-brand-orange"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="mt-2 flex w-full items-center justify-center gap-2 rounded-lg bg-brand-orange py-2.5 text-xs font-bold text-white shadow-card transition-all hover:bg-brand-orange-dark active:scale-[0.99] disabled:opacity-70"
          >
            {submitting ? (
              <>
                <Loader2 size={14} className="animate-spin" />
                <span>Memverifikasi Akses...</span>
              </>
            ) : (
              <span>Masuk Administrator</span>
            )}
          </button>
        </form>

        <div className="mt-8 border-t border-line/60 pt-4 text-center">
          <p className="text-[11px] text-muted">
            Dilindungi enkripsi sistem dan validasi RBAC server-side.
          </p>
        </div>
      </div>
    </div>
  )
}
