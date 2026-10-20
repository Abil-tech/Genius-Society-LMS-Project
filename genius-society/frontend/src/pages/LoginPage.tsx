import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { User, Lock, AlertCircle, Loader2 } from 'lucide-react'
import { useAuth } from '../contexts/AuthContext'

export default function LoginPage() {
  const [username, setUsername] = useState('')
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
      const user = await login(username, password)
      // Arahkan ke dashboard sesuai role
      switch (user.role) {
        case 'super_admin':
          navigate('/super-admin')
          break
        case 'admin':
          navigate('/admin')
          break
        case 'teacher':
          navigate('/teacher')
          break
        case 'student':
          navigate('/student')
          break
        case 'curriculum':
          navigate('/curriculum')
          break
        case 'principal':
          navigate('/principal')
          break
        default:
          navigate('/')
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('Terjadi kesalahan saat masuk. Silakan periksa kembali akun Anda.')
      }
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-[#F1F1F3] px-4 font-sans text-brand-navy">
      {/* Top Left Watermark / Code seperti di login-page.png */}
      <div className="absolute left-6 top-6 select-none font-mono text-[11px] font-semibold tracking-wider text-[#A37B3E]">
        GENIUS SOCIETY
        <span className="block text-[9px] text-gray-400 font-normal">+ 00.00.01</span>
      </div>

      {/* Bottom Corner Coordinates seperti di referensi */}
      <div className="pointer-events-none absolute bottom-4 left-6 select-none font-mono text-[9px] text-gray-400">
        + 99.00.01
      </div>
      <div className="pointer-events-none absolute bottom-4 right-6 select-none font-mono text-[9px] text-gray-400">
        + 49.99.01
      </div>

      {/* Card Login */}
      <div className="relative w-full max-w-[440px] rounded-2xl border border-gray-200/80 bg-white p-8 shadow-sm sm:p-10">
        {/* Subtle top-right cut/fold line accent */}
        <div className="absolute right-0 top-0 h-6 w-6 border-b border-l border-gray-100 bg-gray-50/50 rounded-bl" />

        <div className="mb-6">
          <span className="block font-mono text-[10px] font-medium uppercase tracking-widest text-gray-400">
            LOGIN PAGE
          </span>
          <h1 className="mt-1 text-2xl font-black tracking-tight text-[#111827]">
            GENIUS SOCIETY
          </h1>
          <p className="mt-1 text-xs text-gray-500">
            Akun dan Password tanyakan kepada walas
          </p>
        </div>

        <div className="mb-6 border-b border-gray-200/80" />

        {error && (
          <div className="mb-5 flex items-center gap-2 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-600">
            <AlertCircle size={15} className="shrink-0" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1.5 block text-xs font-semibold text-gray-700">
              Username
            </label>
            <div className="relative flex items-center">
              <span className="pointer-events-none absolute left-3.5 text-gray-400">
                <User size={15} />
              </span>
              <input
                type="text"
                required
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="ENG-000000"
                className="w-full rounded-md border border-orange-200/70 bg-[#FFFBF7] py-2.5 pl-10 pr-3 font-mono text-xs text-gray-800 placeholder-gray-400 transition-colors focus:border-brand-orange focus:bg-white focus:outline-none focus:ring-1 focus:ring-brand-orange"
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
                className="w-full rounded-md border border-orange-200/70 bg-[#FFFBF7] py-2.5 pl-10 pr-3 font-mono text-xs text-gray-800 placeholder-gray-400 transition-colors focus:border-brand-orange focus:bg-white focus:outline-none focus:ring-1 focus:ring-brand-orange"
              />
            </div>
          </div>

          <div className="flex justify-end pt-1">
            <button
              type="button"
              onClick={() => alert('Silakan hubungi Wali Kelas atau Admin sekolah Anda untuk reset akun.')}
              className="text-[11px] font-semibold text-[#B87A2C] transition-colors hover:text-brand-orange"
            >
              Lupa Password
            </button>
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="mt-2 flex w-full items-center justify-center gap-2 rounded-md bg-[#FF9500] py-3 text-xs font-bold uppercase tracking-wider text-white shadow-sm transition-all hover:bg-[#E68600] active:scale-[0.99] disabled:opacity-70"
          >
            {submitting ? (
              <>
                <Loader2 size={14} className="animate-spin" />
                <span>Memproses...</span>
              </>
            ) : (
              <span>LOGIN &rarr;</span>
            )}
          </button>
        </form>

        <div className="mt-8 border-t border-gray-100 pt-4 text-center">
          <Link
            to="/admin/login"
            className="text-[11px] font-medium text-gray-400 hover:text-brand-navy transition-colors"
          >
            Masuk sebagai Administrator / Super Admin &rarr;
          </Link>
        </div>
      </div>
    </div>
  )
}
