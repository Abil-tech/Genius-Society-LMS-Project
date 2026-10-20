import { BarChart3, ClipboardList, FileCheck2, MessageSquare, MonitorPlay, Trophy } from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

interface Feature {
  icon: LucideIcon
  title: string
  description: string
}

const features: Feature[] = [
  {
    icon: MonitorPlay,
    title: 'Interactive Learning',
    description: 'Modul interaktif dengan integrasi multimedia untuk meningkatkan pemahaman siswa.',
  },
  {
    icon: ClipboardList,
    title: 'Assignment Management',
    description: 'Pengelolaan tugas terstruktur dengan tenggat waktu dan notifikasi otomatis.',
  },
  {
    icon: FileCheck2,
    title: 'Online Assessment',
    description: 'Sistem ujian online terintegrasi dengan beragam tipe soal dan penilaian otomatis.',
  },
  {
    icon: BarChart3,
    title: 'Learning Analytics',
    description: 'Pantau perkembangan akademik siswa dengan laporan komprehensif dan visual.',
  },
  {
    icon: MessageSquare,
    title: 'Discussion Forum',
    description: 'Ruang kolaborasi siswa dan guru untuk tanya jawab dan diskusi materi.',
  },
  {
    icon: Trophy,
    title: 'Digital Certificate',
    description: 'Sertifikat digital otentik bagi siswa yang menyelesaikan modul kompetensi.',
  },
]

export default function Features() {
  return (
    <section id="fitur" className="mx-auto max-w-6xl px-4 pb-24 pt-20 sm:px-6">
      <div className="text-center">
        <h2 className="text-xl font-bold text-brand-navy sm:text-2xl">Keunggulan Platform LMS Kami</h2>
        <p className="mx-auto mt-2 max-w-md text-sm text-muted">
          Dirancang khusus untuk mendukung kegiatan belajar mengajar yang efektif, interaktif, dan
          terukur.
        </p>
      </div>

      <ul className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {features.map(({ icon: Icon, title, description }) => (
          <li
            key={title}
            className="rounded-card border border-line/70 bg-surface p-6 shadow-card transition-shadow hover:shadow-md"
          >
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-peach text-brand-orange">
              <Icon size={16} aria-hidden="true" />
            </span>
            <h3 className="mt-5 text-sm font-bold text-brand-navy">{title}</h3>
            <p className="mt-2 text-xs leading-relaxed text-muted">{description}</p>
          </li>
        ))}
      </ul>
    </section>
  )
}
