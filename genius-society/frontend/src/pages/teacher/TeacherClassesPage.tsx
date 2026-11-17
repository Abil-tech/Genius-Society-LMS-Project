import { useEffect, useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { DoorOpen, Users, Search, BookOpen, ChevronRight, RefreshCw } from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getTeacherClasses } from '../../services/teacherService'
import type { TeacherClass } from '../../types/teacher'

export default function TeacherClassesPage() {
  const [classes, setClasses] = useState<TeacherClass[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchClasses = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getTeacherClasses()
      if (res.success && res.data) {
        setClasses(res.data)
      } else {
        setError(res.error?.message || 'Gagal mengambil daftar kelas.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala koneksi.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchClasses()
  }, [])

  const filteredClasses = useMemo(() => {
    if (!search.trim()) return classes
    const q = search.toLowerCase()
    return classes.filter(
      (c) =>
        c.name.toLowerCase().includes(q) ||
        c.code.toLowerCase().includes(q) ||
        c.educationLevel.toLowerCase().includes(q)
    )
  }, [classes, search])

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Kelas yang Diajar
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Daftar rombongan belajar tempat Anda ditugaskan mengajar materi dan mata pelajaran.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchClasses}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5] transition-colors"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Segarkan</span>
            </button>
          </div>
        </div>

        {/* Search Bar */}
        <div className="relative max-w-md">
          <Search
            size={16}
            className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
          />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari kelas berdasarkan nama atau tingkat..."
            className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
          />
        </div>

        {/* Content State */}
        {error ? (
          <ErrorState title="Gagal Memuat Kelas" message={error} onRetry={fetchClasses} />
        ) : loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div
                key={i}
                className="h-44 rounded-2xl border border-[#EBD4C5] bg-white p-5 animate-pulse"
              />
            ))}
          </div>
        ) : filteredClasses.length === 0 ? (
          <EmptyState
            icon={DoorOpen}
            title={search ? 'Kelas Tidak Ditemukan' : 'Belum Ada Penugasan Kelas'}
            description={
              search
                ? `Tidak ditemukan kelas yang cocok dengan kata kunci "${search}".`
                : 'Anda belum ditugaskan mengajar di kelas manapun oleh tim Kurikulum.'
            }
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredClasses.map((cls) => (
              <div
                key={cls.id}
                className="flex flex-col justify-between rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs transition-all hover:border-[#FF9800]/60 hover:shadow-md group"
              >
                <div>
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-[#FF9800] group-hover:bg-[#FF9800] group-hover:text-white transition-colors">
                      <DoorOpen size={22} />
                    </div>
                    <span className="rounded-full bg-stone-100 border border-stone-200 px-2.5 py-0.5 text-[11px] font-bold text-[#30251D]">
                      Tingkat {cls.educationLevel}
                    </span>
                  </div>

                  <h3 className="mt-3 text-base sm:text-lg font-bold text-[#30251D] group-hover:text-[#FF9800] transition-colors">
                    {cls.name}
                  </h3>
                  <p className="text-xs text-[#687385]">Kode Rombel: {cls.code}</p>

                  <div className="mt-4 flex items-center gap-4 text-xs text-[#687385] border-t border-[#EBD4C5]/60 pt-3">
                    <span className="inline-flex items-center gap-1.5 font-medium">
                      <Users size={14} className="text-[#FF9800]" />
                      <span>{cls.studentCount ?? 0} Siswa</span>
                    </span>
                    {cls.subjects && cls.subjects.length > 0 && (
                      <span className="inline-flex items-center gap-1.5 font-medium">
                        <BookOpen size={14} className="text-[#FF9800]" />
                        <span>{cls.subjects.map((s) => s.name).join(', ')}</span>
                      </span>
                    )}
                  </div>
                </div>

                <div className="mt-5 pt-3 border-t border-[#EBD4C5]/40">
                  <Link
                    to={`/teacher/classes/${cls.id}`}
                    className="inline-flex w-full items-center justify-between rounded-xl bg-[#FFF8F5] px-3.5 py-2 text-xs font-bold text-[#30251D] hover:bg-[#FEEADC] hover:text-[#FF9800] transition-colors"
                  >
                    <span>Daftar Siswa & Detail</span>
                    <ChevronRight size={14} />
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </TeacherLayout>
  )
}
