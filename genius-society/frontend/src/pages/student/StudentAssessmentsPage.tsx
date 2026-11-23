import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Search, Clock, ArrowRight } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getAssessments } from '../../services/student/studentAssessment.service'
import type { Assessment } from '../../types/student'

export default function StudentAssessmentsPage() {
  const [assessments, setAssessments] = useState<Assessment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getAssessments()
        if (isMounted) setAssessments(res)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat assessment'
          setError(msg)
        }
      } finally {
        if (isMounted) setLoading(false)
      }
    }
    fetchData()
    return () => {
      isMounted = false
    }
  }, [])

  const filtered = assessments.filter((a) =>
    a.title.toLowerCase().includes(search.toLowerCase())
  )

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'in_progress':
        return (
          <span className="rounded-lg bg-amber-50 px-2.5 py-1 text-[10px] font-extrabold text-amber-700 border border-amber-200 animate-pulse">
            Sedang Dikerjakan
          </span>
        )
      case 'submitted':
        return (
          <span className="rounded-lg bg-emerald-50 px-2.5 py-1 text-[10px] font-extrabold text-emerald-700 border border-emerald-200">
            Sudah Submit
          </span>
        )
      case 'completed':
        return (
          <span className="rounded-lg bg-blue-50 px-2.5 py-1 text-[10px] font-extrabold text-blue-700 border border-blue-200">
            Selesai
          </span>
        )
      case 'closed':
        return (
          <span className="rounded-lg bg-gray-100 px-2.5 py-1 text-[10px] font-extrabold text-gray-600 border border-gray-200">
            Ditutup
          </span>
        )
      default:
        return (
          <span className="rounded-lg bg-[#FEEADC] px-2.5 py-1 text-[10px] font-extrabold text-[#FF9800] border border-[#FF9800]/30">
            Tersedia
          </span>
        )
    }
  }

  return (
    <StudentLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
            Assessment & Ujian Online
          </h1>
          <p className="mt-1 text-xs text-[#687385]">
            Ikuti kuis dan evaluasi berkala untuk mengukur pemahaman materi.
          </p>
        </div>

        {/* Search */}
        <div className="relative">
          <Search size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
          <input
            type="text"
            placeholder="Cari assessment..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2.5 pl-10 pr-4 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
          />
        </div>

        {/* List Content */}
        {error ? (
          <ErrorState
            title="Gagal Memuat Assessment"
            message={error}
            onRetry={() => window.location.reload()}
          />
        ) : loading ? (
          <div className="space-y-3">
            <div className="h-24 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
          </div>
        ) : filtered.length === 0 ? (
          <EmptyState
            title="Tidak ada assessment aktif"
            description="Saat ini belum ada assessment atau ujian yang dijadwalkan untuk kelasmu."
          />
        ) : (
          <div className="space-y-3">
            {filtered.map((a) => (
              <div
                key={a.id}
                className="flex flex-col sm:flex-row sm:items-center justify-between p-4 sm:p-5 rounded-2xl border border-[#EBD4C5] bg-white shadow-xs hover:border-[#FF9800] transition-colors gap-3"
              >
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="rounded-md bg-[#FEEADC] px-2 py-0.5 text-[10px] font-extrabold text-[#FF9800]">
                      {a.subjectName}
                    </span>
                    {getStatusBadge(a.status)}
                  </div>
                  <h3 className="text-sm sm:text-base font-extrabold text-[#30251D] truncate">
                    {a.title}
                  </h3>
                  <div className="mt-1 flex items-center gap-3 text-[11px] text-[#687385]">
                    <span className="flex items-center gap-1">
                      <Clock size={13} className="text-[#FF9800]" />
                      <span>Durasi: {a.durationMinutes} Menit</span>
                    </span>
                    <span>&bull; Kategori: {a.category?.toUpperCase()}</span>
                  </div>
                </div>

                <div className="flex items-center justify-between sm:justify-end gap-3 pt-2 sm:pt-0 border-t sm:border-0 border-[#EBD4C5]/60">
                  {a.score !== undefined && (
                    <div className="flex items-center gap-1.5 px-3 py-1 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-xs font-black text-[#30251D]">
                      <span>Skor:</span>
                      <span className="text-[#FF9800]">{a.score}</span>
                    </div>
                  )}

                  <Link
                    to={`/student/assessments/${a.id}`}
                    className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
                  >
                    <span>
                      {a.status === 'in_progress'
                        ? 'Lanjutkan'
                        : a.status === 'completed' || a.status === 'submitted'
                        ? 'Lihat Hasil'
                        : 'Mulai'}
                    </span>
                    <ArrowRight size={13} />
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </StudentLayout>
  )
}
