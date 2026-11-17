import { useEffect, useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import {
  Plus,
  Search,
  RefreshCw,
  Clock,
  Trash2,
  Edit,
  Eye,
  Users,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import Pagination from '../../components/common/Pagination'
import StatusBadge from '../../components/common/StatusBadge'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  getAssessments,
  deleteAssessment,
  getTeacherClasses,
  getTeacherSubjects,
} from '../../services/teacherService'
import type { Assessment, TeacherClass, TeacherSubject } from '../../types/teacher'

export default function TeacherAssessmentsPage() {
  const toast = useToast()
  const [assessments, setAssessments] = useState<Assessment[]>([])
  const [classes, setClasses] = useState<TeacherClass[]>([])
  const [subjects, setSubjects] = useState<TeacherSubject[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Filters
  const [search, setSearch] = useState('')
  const [selectedClass, setSelectedClass] = useState('')
  const [selectedSubject, setSelectedSubject] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 10

  // Delete target
  const [deleteTarget, setDeleteTarget] = useState<Assessment | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [assRes, clsRes, subRes] = await Promise.allSettled([
        getAssessments({ search, classId: selectedClass, subjectId: selectedSubject }),
        getTeacherClasses(),
        getTeacherSubjects(),
      ])

      if (assRes.status === 'fulfilled' && assRes.value.success && assRes.value.data) {
        setAssessments(assRes.value.data)
      } else if (assRes.status === 'fulfilled' && !assRes.value.success) {
        setError(assRes.value.error?.message || 'Gagal mengambil data assessment.')
      }

      if (clsRes.status === 'fulfilled' && clsRes.value.success && clsRes.value.data) {
        setClasses(clsRes.value.data)
      }
      if (subRes.status === 'fulfilled' && subRes.value.success && subRes.value.data) {
        setSubjects(subRes.value.data)
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat assessment.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [selectedClass, selectedSubject])

  const filteredAssessments = useMemo(() => {
    let list = assessments
    if (selectedCategory) {
      list = list.filter((a) => a.category === selectedCategory)
    }
    if (search.trim()) {
      const q = search.toLowerCase()
      list = list.filter(
        (a) =>
          a.title.toLowerCase().includes(q) ||
          (a.description && a.description.toLowerCase().includes(q))
      )
    }
    return list
  }, [assessments, search, selectedCategory])

  const totalPages = Math.ceil(filteredAssessments.length / pageSize)
  const paginatedAssessments = useMemo(() => {
    const start = (page - 1) * pageSize
    return filteredAssessments.slice(start, start + pageSize)
  }, [filteredAssessments, page, pageSize])

  const handleDeleteConfirm = async () => {
    if (!deleteTarget) return
    setIsDeleting(true)
    try {
      const res = await deleteAssessment(deleteTarget.id)
      if (res.success) {
        toast.success(
          'Assessment Diproses',
          res.data?.archived ? 'Assessment telah diarsipkan karena memiliki peserta.' : 'Assessment telah dihapus.'
        )
        setDeleteTarget(null)
        fetchData()
      } else {
        toast.error('Gagal Menghapus', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala sistem.'
      toast.error('Gagal Menghapus', msg)
    } finally {
      setIsDeleting(false)
    }
  }

  const columns: Column<Assessment>[] = [
    {
      key: 'title',
      header: 'Judul Ujian / Assessment',
      render: (a) => (
        <div className="max-w-xs sm:max-w-sm">
          <Link
            to={`/teacher/assessments/${a.id}`}
            className="font-bold text-[#30251D] hover:text-[#FF9800] transition-colors line-clamp-1"
          >
            {a.title}
          </Link>
          <div className="mt-1 flex items-center gap-2">
            <span
              className={`rounded px-1.5 py-0.5 text-[10px] font-extrabold uppercase ${
                a.category === 'uts'
                  ? 'bg-amber-100 text-amber-800'
                  : a.category === 'uas'
                  ? 'bg-rose-100 text-rose-800'
                  : 'bg-emerald-100 text-emerald-800'
              }`}
            >
              {a.category.toUpperCase()}
            </span>
            <StatusBadge status={a.status} />
          </div>
        </div>
      ),
    },
    {
      key: 'subject',
      header: 'Mata Pelajaran',
      render: (a) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {a.subjectName || 'Mata Pelajaran'}
        </span>
      ),
    },
    {
      key: 'duration',
      header: 'Durasi & Soal',
      render: (a) => (
        <div className="text-xs">
          <span className="font-bold text-[#30251D] flex items-center gap-1">
            <Clock size={12} className="text-[#FF9800]" />
            <span>{a.durationMinutes} Menit</span>
          </span>
          <span className="text-[11px] text-[#687385] block mt-0.5">
            {a.totalQuestions ?? a.questions?.length ?? 0} Butir Soal
          </span>
        </div>
      ),
    },
    {
      key: 'classes',
      header: 'Kelas',
      render: (a) => (
        <div className="flex flex-wrap gap-1">
          {a.classNames && a.classNames.length > 0 ? (
            a.classNames.map((c, i) => (
              <span
                key={i}
                className="rounded-md bg-[#FFF8F5] border border-[#EBD4C5] px-2 py-0.5 text-[11px] font-bold text-[#30251D]"
              >
                {c}
              </span>
            ))
          ) : (
            <span className="text-gray-400 text-xs">-</span>
          )}
        </div>
      ),
    },
    {
      key: 'attempts',
      header: 'Peserta',
      render: (a) => (
        <Link
          to={`/teacher/assessments/${a.id}`}
          className="inline-flex items-center gap-1.5 rounded-lg bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-1 text-xs font-bold text-[#FF9800] hover:bg-[#FEEADC] transition-colors"
        >
          <Users size={12} />
          <span>
            {a.gradedAttempts ?? 0}/{a.totalAttempts ?? 0} Dinilai
          </span>
        </Link>
      ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'right',
      render: (a) => (
        <div className="flex items-center justify-end gap-1.5">
          <Link
            to={`/teacher/assessments/${a.id}`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Lihat Detail & Soal"
          >
            <Eye size={14} />
          </Link>
          <Link
            to={`/teacher/assessments/${a.id}/edit`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Edit Assessment"
          >
            <Edit size={14} />
          </Link>
          <button
            type="button"
            onClick={() => setDeleteTarget(a)}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-rose-200 text-rose-600 hover:bg-rose-50 transition-colors"
            title="Hapus / Arsipkan"
          >
            <Trash2 size={14} />
          </button>
        </div>
      ),
    },
  ]

  const hasAttempts = (deleteTarget?.totalAttempts ?? 0) > 0

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Assessment & Ujian
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Kelola kuis harian, UTS, dan UAS. Mendukung soal pilihan ganda otomatis dan esai manual.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchData}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Segarkan</span>
            </button>

            <Link
              to="/teacher/assessments/new"
              className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs sm:text-sm font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
            >
              <Plus size={16} />
              <span>Buat Assessment</span>
            </Link>
          </div>
        </div>

        {/* Filters */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <Search
              size={16}
              className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
            />
            <input
              type="text"
              value={search}
              onChange={(e) => {
                setSearch(e.target.value)
                setPage(1)
              }}
              placeholder="Cari assessment berdasarkan judul..."
              className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
            />
          </div>

          <select
            value={selectedCategory}
            onChange={(e) => {
              setSelectedCategory(e.target.value)
              setPage(1)
            }}
            aria-label="Filter berdasarkan kategori ujian"
            className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs sm:text-sm font-medium text-[#30251D] focus:border-[#FF9800] focus:outline-none"
          >
            <option value="">Semua Kategori</option>
            <option value="daily">Harian (Daily)</option>
            <option value="uts">UTS (Tengah Semester)</option>
            <option value="uas">UAS (Akhir Semester)</option>
          </select>

          <select
            value={selectedClass}
            onChange={(e) => {
              setSelectedClass(e.target.value)
              setPage(1)
            }}
            aria-label="Filter berdasarkan kelas"
            className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs sm:text-sm font-medium text-[#30251D] focus:border-[#FF9800] focus:outline-none"
          >
            <option value="">Semua Kelas</option>
            {classes.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>

          <select
            value={selectedSubject}
            onChange={(e) => {
              setSelectedSubject(e.target.value)
              setPage(1)
            }}
            aria-label="Filter berdasarkan mata pelajaran"
            className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs sm:text-sm font-medium text-[#30251D] focus:border-[#FF9800] focus:outline-none"
          >
            <option value="">Semua Mata Pelajaran</option>
            {subjects.map((s) => (
              <option key={s.id} value={s.id}>
                {s.name}
              </option>
            ))}
          </select>
        </div>

        {/* Content */}
        {error ? (
          <ErrorState title="Gagal Memuat Assessment" message={error} onRetry={fetchData} />
        ) : (
          <div className="space-y-4">
            <DataTable
              columns={columns}
              data={paginatedAssessments}
              loading={loading}
              emptyMessage={search ? 'Assessment tidak ditemukan' : 'Belum ada ujian / assessment'}
              emptySubMessage={
                search
                  ? `Tidak ada hasil pencarian untuk "${search}"`
                  : 'Klik "Buat Assessment" untuk menyusun soal ujian baru.'
              }
            />

            <Pagination
              currentPage={page}
              totalPages={totalPages}
              onPageChange={setPage}
              totalItems={filteredAssessments.length}
              pageSize={pageSize}
            />
          </div>
        )}

        {/* Delete Confirmation Dialog */}
        <ConfirmDialog
          isOpen={Boolean(deleteTarget)}
          onClose={() => setDeleteTarget(null)}
          onConfirm={handleDeleteConfirm}
          title={hasAttempts ? 'Arsipkan Assessment' : 'Hapus Assessment'}
          message={`Apakah Anda yakin ingin menghapus atau mengarsipkan ujian "${deleteTarget?.title}"?`}
          hasSubmissions={hasAttempts}
          archiveNotice="Assessment ini sudah diikuti oleh siswa. Sesuai regulasi penilaian, data akan diarsipkan agar rekap nilai siswa tetap aman."
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
