import { BarChart3, Users } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

interface Stat {
  icon: LucideIcon
  value: string
  label: string
}

// Nilai statis sesuai Figma; nanti dapat diganti data dari API publik.
const stats: Stat[] = [
  { icon: BarChart3, value: '3000+', label: 'Siswa Aktif' },
  { icon: Users, value: '150+', label: 'Guru Aktif' },
]

export default function Stats() {
  return (
    <section aria-label="Statistik platform" className="mx-auto -mt-12 max-w-6xl px-4 sm:px-6">
      <ul className="mx-auto flex max-w-md flex-col gap-3 sm:flex-row sm:justify-center">
        {stats.map(({ icon: Icon, value, label }) => (
          <li
            key={label}
            className="flex flex-1 items-center gap-3 rounded-card border border-line/70 bg-surface px-5 py-4 shadow-card"
          >
            <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-brand-peach text-brand-orange">
              <Icon size={16} aria-hidden="true" />
            </span>
            <div>
              <p className="text-lg font-bold leading-tight text-brand-navy">{value}</p>
              <p className="text-[11px] text-muted">{label}</p>
            </div>
          </li>
        ))}
      </ul>
    </section>
  )
}
