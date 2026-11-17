import { useEffect, useState, useMemo } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, Users, Search, RefreshCw, Mail, Phone } from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import ErrorState from '../../components/common/ErrorState'
import { getClassStudents, getClassDetail } from '../../services/teacherService'
import type { TeacherStudent, TeacherClass } from '../../types/teacher'

export default function TeacherClassDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [classInfo, setClassInfo] = useState<TeacherClass | null>(null)
  const [students, setStudents] = useState<TeacherStudent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchData = async () => {
    if (!id) return
    setLoading(true)
    setError(null)
    try {
      const [classRes, studentsRes] = await Promise.allSettled([
        getClassDetail(id),
        getClassStudents(id),
      ])

      if (classRes.status === 'fulfilled' && classRes.value.success && classRes.value.data) {
        setClassInfo(classRes.value.data)
      }
      if (studentsRes.status === 'fulfilled' && studentsRes.value.success && studentsRes.value.data) {
        setStudents(studentsRes.value.data)
      } else if (studentsRes.status === 'fulfilled' && !studentsRes.value.success) {
        setError(studentsRes.value.error?.message || 'Gagal mengambil data siswa.')
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
  }, [id])

  const filteredStudents = useMemo(() => {
    if (!search.trim()) return students
    const q = search.toLowerCase()
    return students.filter(
      (s) =>
        s.name.toLowerCase().includes(q) ||
        (s.nis && s.nis.toLowerCase().includes(q)) ||
        (s.email && s.email.toLowerCase().includes(q))
    )
  }, [students, search])

  const columns: Column<TeacherStudent>[] = [
    {
      key: 'no',
      header: 'No',
      className: 'w-12',
      render: (_, idx) => <span className="font-medium text-gray-500">{idx + 1}</span>,
    },
    {
      key: 'nis',
      header: 'NIS / NISN',
      render: (s) => (
        <div>
          <span className="font-bold text-[#30251D]">{s.nis || '-'}</span>
          {s.nisn && <span className="block text-[11px] text-[#687385]">NISN: {s.nisn}</span>}
        </div>
      ),
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
            {s.email && (
              <span className="flex items-center gap-1 text-[11px] text-[#687385]">
                <Mail size={11} />
                <span>{s.email}</span>
              </span>
            )}
          </div>
        </div>
      ),
    },
    {
      key: 'gender',
      header: 'L / P',
      align: 'center',
      render: (s) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {s.gender === 'L' || s.gender === 'M'
            ? 'Laki-laki'
            : s.gender === 'P' || s.gender === 'F'
            ? 'Perempuan'
            : s.gender || '-'}
        </span>
      ),
    },
    {
      key: 'phone',
      header: 'Kontak',
      render: (s) =>
        s.phone ? (
          <span className="flex items-center gap-1 text-xs text-[#687385]">
            <Phone size={12} className="text-[#FF9800]" />
            <span>{s.phone}</span>
          </span>
        ) : (
          <span className="text-gray-400 text-xs">-</span>
        ),
    },
    {
      key: 'status',
      header: 'Status',
      align: 'right',
      render: (s) => (
        <span className="inline-flex rounded-full bg-emerald-50 border border-emerald-200 px-2 py-0.5 text-[11px] font-bold text-emerald-800">
          {s.status || 'Aktif'}
        </span>
      ),
    },
  ]

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Breadcrumb / Back button */}
        <div className="flex items-center justify-between">
          <Link
            to="/teacher/classes"
            className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
          >
            <ArrowLeft size={16} />
            <span>Kembali ke Daftar Kelas</span>
          </Link>

          <button
            type="button"
            onClick={fetchData}
            disabled={loading}
            className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
          >
            <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
            <span>Segarkan</span>
          </button>
        </div>

        {/* Header */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div>
              <span className="rounded-full bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-0.5 text-[10px] font-extrabold uppercase tracking-wider text-[#FF9800]">
                Rombongan Belajar
              </span>
              <h1 className="mt-1 text-xl sm:text-2xl font-extrabold text-[#30251D]">
                {classInfo?.name || 'Detail Kelas'}
              </h1>
              <p className="text-xs text-[#687385]">
                {classInfo?.code ? `Kode: ${classInfo.code} • ` : ''}Tingkat {classInfo?.educationLevel || '-'}
              </p>
            </div>

            <div className="flex items-center gap-3">
              <div className="flex items-center gap-2 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] px-4 py-2">
                <Users size={18} className="text-[#FF9800]" />
                <div>
                  <span className="block text-[10px] font-semibold text-[#687385] uppercase">
                    Total Siswa
                  </span>
                  <span className="text-base font-extrabold text-[#30251D]">
                    {students.length} Siswa
                  </span>
                </div>
              </div>
            </div>
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
            placeholder="Cari nama atau NIS siswa..."
            className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
          />
        </div>

        {/* Content */}
        {error ? (
          <ErrorState title="Gagal Memuat Daftar Siswa" message={error} onRetry={fetchData} />
        ) : (
          <DataTable
            columns={columns}
            data={filteredStudents}
            loading={loading}
            emptyMessage={search ? 'Siswa tidak ditemukan' : 'Belum ada data siswa di kelas ini'}
            emptySubMessage={
              search ? `Tidak ada hasil untuk kata kunci "${search}"` : undefined
            }
          />
        )}
      </div>
    </TeacherLayout>
  )
}
