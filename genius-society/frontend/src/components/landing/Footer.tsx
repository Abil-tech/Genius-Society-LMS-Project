import { GraduationCap } from 'lucide-react'

const links = ['Privacy Policy', 'Terms of Service', 'FAQ', 'Support']

export default function Footer() {
  return (
    <footer className="mx-4 rounded-t-card bg-brand-orange text-white sm:mx-6">
      <div className="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="flex items-center gap-2 text-sm font-bold">
            <GraduationCap size={18} aria-hidden="true" />
            Genius Society
          </p>
          <p className="mt-2 max-w-xs text-xs text-white/90">
            Platform pembelajaran digital untuk generasi profesional masa depan.
          </p>
        </div>

        <div className="sm:text-right">
          <ul className="flex flex-wrap gap-x-5 gap-y-1 text-xs sm:justify-end">
            {links.map((label) => (
              <li key={label}>
                <a href="#" className="underline-offset-2 hover:underline">
                  {label}
                </a>
              </li>
            ))}
          </ul>
          <p className="mt-3 text-[11px] text-white/90">
            © {new Date().getFullYear()} Genius Society. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}
