import { Link } from 'react-router-dom'

export default function Hero() {
  return (
    <section className="mx-auto max-w-6xl px-4 pt-6 sm:px-6">
      <div className="rounded-card bg-surface/60 px-4 pb-24 pt-14 text-center sm:px-8 sm:pt-16">
        <span className="inline-block rounded-pill border border-brand-orange/50 bg-surface px-3 py-1 text-[11px] font-medium text-brand-orange">
          Sistem Manajemen Pembelajaran Masa Depan
        </span>

        <h1 className="mx-auto mt-5 max-w-xl text-4xl font-extrabold leading-tight tracking-tight text-brand-navy sm:text-5xl">
          Revolusi Pembelajaran Digital
          <span className="block text-brand-orange">GENIUS SOCIETY</span>
        </h1>

        <p className="mx-auto mt-5 max-w-lg text-sm leading-relaxed text-muted">
          Tingkatkan pengalaman belajar dengan platform terintegrasi: akses materi, kelola tugas, dan
          pantau perkembangan akademik secara real-time dalam lingkungan belajar yang aman dan
          interaktif.
        </p>

        <Link
          to="/login"
          className="mt-7 inline-block rounded-pill bg-brand-orange px-6 py-2.5 text-sm font-semibold text-white shadow-card transition-colors hover:bg-brand-orange-dark"
        >
          Mulai Belajar
        </Link>
      </div>
    </section>
  )
}
