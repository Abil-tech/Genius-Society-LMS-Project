import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Search, Clock, ArrowRight } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getAssignments } from '../../services/student/studentAssignment.service'
import type { Assignment } from '../../types/student'

export default function StudentAssignmentsPage() {
  const [assignments, setAssignments] = useState<Assignment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getAssignments()
        if (isMounted) setAssignments(res)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat tugas'
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

  const filtered = assignments.filter((a) => {
    const matchSearch = a.title.toLowerCase().includes(search.toLowerCase())
    if (!matchSearch) return false
    if (statusFilter === 'all') return true
    return a.status === statusFilter
  })

  const getStatusBadge = (status: string, deadline: string) => {
    const isUrgent =
      status === 'pending' &&
      new Date(deadline).getTime() - Date.now() < 24 * 3600 * 1000 &&
      new Date(deadline).getTime() > Date.now()

    if (status === 'graded') {
      return (
        <span className="rounded-lg bg-blue-50 px-2.5 py-1 text-[10px] font-extrabold text-blue-700 border border-blue-200">
          Sudah Dinilai
        </span>
      )
    }
    if (status === 'submitted') {
      return (
        <span className="rounded-lg bg-emerald-50 px-2.5 py-1 text-[10px] font-extrabold text-emerald-700 border border-emerald-200">
          Sudah Dikumpulkan
        </span>
      )
    }
    if (status === 'overdue') {
      return (
        <span className="rounded-lg bg-rose-50 px-2.5 py-1 text-[10px] font-extrabold text-rose-700 border border-rose-200">
          Terlewat
        </span>
      )
    }
    return (
      <span
        className={`rounded-lg px-2.5 py-1 text-[10px] font-extrabold border ${
          isUrgent
            ? 'bg-rose-50 text-rose-600 border-rose-300 animate-pulse'
            : 'bg-[#FEEADC] text-[#FF9800] border-[#FF9800]/30'
        }`}
      >
        {isUrgent ? 'Segera Berakhir' : 'Belum Dikumpulkan'}
      </span>
    )
  }

  return (
    <StudentLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
            Daftar Tugas
          </h1>
          <p className="mt-1 text-xs text-[#687385]">
            Kerjakan dan kumpulkan tugas sebelum batas waktu yang ditentukan.
          </p>
        </div>

        {/* Filter bar */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Cari tugas..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2.5 pl-10 pr-4 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
            />
          </div>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2.5 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:outline-none"
          >
            <option value="all">Semua Status</option>
            <option value="pending">Belum Dikumpulkan</option>
            <option value="submitted">Sudah Dikumpulkan</option>
            <option value="graded">Sudah Dinilai</option>
            <option value="overdue">Terlewat (Overdue)</option>
          </select>
        </div>

        {/* List Content */}
        {error ? (
          <ErrorState
            title="Gagal Memuat Tugas"
            message={error}
            onRetry={() => window.location.reload()}
          />
        ) : loading ? (
          <div className="space-y-3">
            <div className="h-24 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
            <div className="h-24 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
          </div>
        ) : filtered.length === 0 ? (
          <EmptyState
            title="Tidak ada tugas ditemukan"
            description="Saat ini belum ada tugas baru yang perlu dikumpulkan."
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
                    {getStatusBadge(a.status, a.deadline)}
                  </div>
                  <h3 className="text-sm sm:text-base font-extrabold text-[#30251D] truncate">
                    {a.title}
                  </h3>
                  <div className="mt-1 flex items-center gap-2 text-[11px] text-[#687385]">
                    <Clock size={13} className="text-[#FF9800]" />
                    <span>
                      Tenggat:{' '}
                      {new Date(a.deadline).toLocaleDateString('id-ID', {
                        day: 'numeric',
                        month: 'short',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </span>
                  </div>
                </div>

                <div className="flex items-center justify-between sm:justify-end gap-3 pt-2 sm:pt-0 border-t sm:border-0 border-[#EBD4C5]/60">
                  {a.submission?.score !== undefined && (
                    <div className="flex items-center gap-1.5 px-3 py-1 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-xs font-black text-[#30251D]">
                      <span>Nilai:</span>
                      <span className="text-[#FF9800]">{a.submission.score}</span>
                    </div>
                  )}

                  <Link
                    to={`/student/assignments/${a.id}`}
                    className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
                  >
                    <span>{a.status === 'pending' ? 'Kerjakan' : 'Lihat Detail'}</span>
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
