import { useEffect, useState, useMemo } from 'react'
import {
  Search,
  RefreshCw,
  Phone,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import ErrorState from '../../components/common/ErrorState'
import { getWalasStudents } from '../../services/teacherService'
import type { WalasStudentDetail } from '../../types/teacher'

export default function TeacherWalasStudentsPage() {
  const [students, setStudents] = useState<WalasStudentDetail[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getWalasStudents()
      if (res.success && res.data) {
        setStudents(res.data)
      } else {
        setError(res.error?.message || 'Gagal memuat daftar siswa perwalian.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat data.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const filteredStudents = useMemo(() => {
    if (!search.trim()) return students
    const q = search.toLowerCase()
    return students.filter(
      (s) =>
        s.name.toLowerCase().includes(q) ||
        (s.nis && s.nis.toLowerCase().includes(q)) ||
        (s.parentName && s.parentName.toLowerCase().includes(q))
    )
  }, [students, search])

  const columns: Column<WalasStudentDetail>[] = [
    {
      key: 'rank',
      header: 'Rank',
      className: 'w-12',
      align: 'center',
      render: (s, idx) => {
        const rank = s.rank ?? idx + 1
        return (
          <span
            className={`flex h-6 w-6 items-center justify-center rounded-full text-xs font-black ${
              rank === 1
                ? 'bg-amber-100 text-amber-900 border border-amber-300'
                : rank === 2
                ? 'bg-stone-200 text-stone-800'
                : rank === 3
                ? 'bg-orange-100 text-orange-800'
                : 'text-gray-500'
            }`}
          >
            {rank}
          </span>
        )
      },
    },
    {
      key: 'nis',
      header: 'NIS',
      render: (s) => <span className="font-bold text-[#30251D]">{s.nis}</span>,
    },
    {
      key: 'name',
      header: 'Nama Siswa',
      render: (s) => (
        <div className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-[#FF9800]/10 text-xs font-bold text-[#FF9800]">
            {s.name.charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="font-bold text-[#30251D]">{s.name}</p>
            {s.email && <p className="text-[11px] text-[#687385]">{s.email}</p>}
          </div>
        </div>
      ),
    },
    {
      key: 'gender',
      header: 'L/P',
      align: 'center',
      render: (s) => (
        <span className="text-xs font-semibold text-[#30251D]">
          {s.gender === 'L' || s.gender === 'M' ? 'L' : 'P'}
        </span>
      ),
    },
    {
      key: 'parent',
      header: 'Orang Tua / Wali',
      render: (s) => (
        <div className="text-xs">
          <p className="font-semibold text-[#30251D]">{s.parentName || '-'}</p>
          {s.parentPhone && (
            <span className="flex items-center gap-1 text-[11px] text-[#687385]">
              <Phone size={10} className="text-[#FF9800]" />
              <span>{s.parentPhone}</span>
            </span>
          )}
        </div>
      ),
    },
    {
      key: 'average',
      header: 'Rata-rata Nilai',
      align: 'right',
      render: (s) => (
        <div className="text-right">
          <span className="text-sm font-black text-[#FF9800]">
            {s.averageScore !== undefined && s.averageScore !== null
              ? s.averageScore.toFixed(1)
              : '-'}
          </span>
        </div>
      ),
    },
  ]

  return (
    <TeacherLayout isWalas={true}>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div className="flex items-center gap-2">
              <span className="rounded-full bg-[#FF9800]/10 px-2.5 py-0.5 text-[10px] font-extrabold uppercase text-[#FF9800]">
                Tugas Wali Kelas
              </span>
            </div>
            <h1 className="mt-1 text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Daftar Siswa Perwalian
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Biodata lengkap, kontak wali, dan ringkasan akademik siswa di kelas binaan Anda (Read-only).
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchData}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Segarkan</span>
            </button>
          </div>
        </div>

        {/* Search */}
        <div className="relative max-w-md">
          <Search
            size={16}
            className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
          />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari siswa atau orang tua..."
            className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
          />
        </div>

        {/* Content */}
        {error ? (
          <ErrorState title="Gagal Memuat Siswa Perwalian" message={error} onRetry={fetchData} />
        ) : (
          <DataTable
            columns={columns}
            data={filteredStudents}
            loading={loading}
            emptyMessage={search ? 'Siswa tidak ditemukan' : 'Belum ada data siswa perwalian'}
            emptySubMessage={
              search ? `Tidak ada hasil pencarian untuk "${search}"` : undefined
            }
          />
        )}
      </div>
    </TeacherLayout>
  )
}
