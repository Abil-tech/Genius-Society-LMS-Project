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
  getAssignments,
  deleteAssignment,
  getTeacherClasses,
  getTeacherSubjects,
} from '../../services/teacherService'
import type { Assignment, TeacherClass, TeacherSubject } from '../../types/teacher'

export default function TeacherAssignmentsPage() {
  const toast = useToast()
  const [assignments, setAssignments] = useState<Assignment[]>([])
  const [classes, setClasses] = useState<TeacherClass[]>([])
  const [subjects, setSubjects] = useState<TeacherSubject[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Filters
  const [search, setSearch] = useState('')
  const [selectedClass, setSelectedClass] = useState('')
  const [selectedSubject, setSelectedSubject] = useState('')
  const [page, setPage] = useState(1)
  const pageSize = 10

  // Delete modal state
  const [deleteTarget, setDeleteTarget] = useState<Assignment | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [assRes, clsRes, subRes] = await Promise.allSettled([
        getAssignments({ search, classId: selectedClass, subjectId: selectedSubject }),
        getTeacherClasses(),
        getTeacherSubjects(),
      ])

      if (assRes.status === 'fulfilled' && assRes.value.success && assRes.value.data) {
        setAssignments(assRes.value.data)
      } else if (assRes.status === 'fulfilled' && !assRes.value.success) {
        setError(assRes.value.error?.message || 'Gagal mengambil data tugas.')
      }

      if (clsRes.status === 'fulfilled' && clsRes.value.success && clsRes.value.data) {
        setClasses(clsRes.value.data)
      }
      if (subRes.status === 'fulfilled' && subRes.value.success && subRes.value.data) {
        setSubjects(subRes.value.data)
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat tugas.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [selectedClass, selectedSubject])

  const filteredAssignments = useMemo(() => {
    let list = assignments
    if (search.trim()) {
      const q = search.toLowerCase()
      list = list.filter(
        (a) =>
          a.title.toLowerCase().includes(q) ||
          (a.description && a.description.toLowerCase().includes(q))
      )
    }
    return list
  }, [assignments, search])

  const totalPages = Math.ceil(filteredAssignments.length / pageSize)
  const paginatedAssignments = useMemo(() => {
    const start = (page - 1) * pageSize
    return filteredAssignments.slice(start, start + pageSize)
  }, [filteredAssignments, page, pageSize])

  const handleDeleteConfirm = async () => {
    if (!deleteTarget) return
    setIsDeleting(true)
    try {
      const res = await deleteAssignment(deleteTarget.id)
      if (res.success) {
        toast.success(
          'Tugas Berhasil Diproses',
          res.data?.archived ? 'Tugas telah diarsipkan karena memiliki submission.' : 'Tugas telah dihapus.'
        )
        setDeleteTarget(null)
        fetchData()
      } else {
        toast.error('Gagal Menghapus Tugas', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kesalahan sistem.'
      toast.error('Gagal Menghapus', msg)
    } finally {
      setIsDeleting(false)
    }
  }

  const columns: Column<Assignment>[] = [
    {
      key: 'title',
      header: 'Judul Tugas',
      render: (a) => {
        const isPast = new Date(a.deadline) < new Date()
        return (
          <div className="max-w-xs sm:max-w-sm">
            <Link
              to={`/teacher/assignments/${a.id}`}
              className="font-bold text-[#30251D] hover:text-[#FF9800] transition-colors line-clamp-1"
            >
              {a.title}
            </Link>
            {a.description && (
              <p className="mt-0.5 text-xs text-[#687385] line-clamp-1">{a.description}</p>
            )}
            <div className="mt-1 flex items-center gap-2">
              <StatusBadge
                status={isPast ? 'ended' : 'active'}
                label={isPast ? 'Berakhir' : 'Aktif'}
              />
              {a.allowResubmit && (
                <span className="text-[10px] text-gray-500 font-medium">Boleh Resubmit</span>
              )}
            </div>
          </div>
        )
      },
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
      key: 'deadline',
      header: 'Tenggat Waktu',
      render: (a) => {
        const isPast = new Date(a.deadline) < new Date()
        return (
          <span
            className={`inline-flex items-center gap-1 text-xs font-semibold ${
              isPast ? 'text-rose-700' : 'text-[#30251D]'
            }`}
          >
            <Clock size={12} className={isPast ? 'text-rose-600' : 'text-[#FF9800]'} />
            <span>
              {new Date(a.deadline).toLocaleString('id-ID', {
                dateStyle: 'medium',
                timeStyle: 'short',
              })}
            </span>
          </span>
        )
      },
    },
    {
      key: 'submissions',
      header: 'Pengumpulan',
      render: (a) => (
        <Link
          to={`/teacher/assignments/${a.id}/submissions`}
          className="inline-flex items-center gap-1.5 rounded-lg bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-1 text-xs font-bold text-[#FF9800] hover:bg-[#FEEADC] transition-colors"
        >
          <Users size={12} />
          <span>
            {a.gradedSubmissions ?? 0}/{a.totalSubmissions ?? 0} Dinilai
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
            to={`/teacher/assignments/${a.id}`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Lihat Detail"
          >
            <Eye size={14} />
          </Link>
          <Link
            to={`/teacher/assignments/${a.id}/edit`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Edit Tugas"
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

  const hasSubmissions = (deleteTarget?.totalSubmissions ?? 0) > 0

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Penugasan Siswa
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Buat latihan mandiri, tugas portofolio, atur tenggat waktu, dan beri nilai pengumpulan siswa.
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
              to="/teacher/assignments/new"
              className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs sm:text-sm font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
            >
              <Plus size={16} />
              <span>Buat Tugas</span>
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
              placeholder="Cari tugas berdasarkan judul atau deskripsi..."
              className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
            />
          </div>

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
          <ErrorState title="Gagal Memuat Tugas" message={error} onRetry={fetchData} />
        ) : (
          <div className="space-y-4">
            <DataTable
              columns={columns}
              data={paginatedAssignments}
              loading={loading}
              emptyMessage={search ? 'Tugas tidak ditemukan' : 'Belum ada tugas pembelajaran'}
              emptySubMessage={
                search
                  ? `Tidak ada hasil pencarian untuk "${search}"`
                  : 'Klik tombol "Buat Tugas" untuk menugaskan latihan ke kelas Anda.'
              }
            />

            <Pagination
              currentPage={page}
              totalPages={totalPages}
              onPageChange={setPage}
              totalItems={filteredAssignments.length}
              pageSize={pageSize}
            />
          </div>
        )}

        {/* Delete Confirmation Dialog */}
        <ConfirmDialog
          isOpen={Boolean(deleteTarget)}
          onClose={() => setDeleteTarget(null)}
          onConfirm={handleDeleteConfirm}
          title={hasSubmissions ? 'Arsipkan Tugas' : 'Hapus Tugas'}
          message={`Apakah Anda yakin ingin menghapus atau mengarsipkan tugas "${deleteTarget?.title}"?`}
          hasSubmissions={hasSubmissions}
          archiveNotice="Tugas ini telah memiliki jawaban/pengumpulan dari siswa. Data akan diarsipkan agar seluruh riwayat pengumpulan siswa tetap tersimpan aman."
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
