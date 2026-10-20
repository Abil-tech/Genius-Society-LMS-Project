import { GraduationCap } from 'lucide-react'
import { Link } from 'react-router-dom'

export default function Navbar() {
  return (
    <header className="sticky top-0 z-30 border-b border-line/60 bg-surface/90 backdrop-blur">
      <nav
        aria-label="Navigasi utama"
        className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4 sm:px-6"
      >
        <Link to="/" className="flex items-center gap-2 text-sm font-bold text-brand-navy">
          <GraduationCap size={18} strokeWidth={2} aria-hidden="true" />
          Genius Society
        </Link>
        <Link
          to="/login"
          className="rounded-pill bg-brand-orange px-5 py-2 text-xs font-semibold text-white shadow-card transition-colors hover:bg-brand-orange-dark"
        >
          Masuk Belajar
        </Link>
      </nav>
    </header>
  )
}
