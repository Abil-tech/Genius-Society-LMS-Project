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
  getProjects,
  deleteProject,
  getTeacherClasses,
  getTeacherSubjects,
} from '../../services/teacherService'
import type { Project, TeacherClass, TeacherSubject } from '../../types/teacher'

export default function TeacherProjectsPage() {
  const toast = useToast()
  const [projects, setProjects] = useState<Project[]>([])
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
  const [deleteTarget, setDeleteTarget] = useState<Project | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [projRes, clsRes, subRes] = await Promise.allSettled([
        getProjects({ search, classId: selectedClass, subjectId: selectedSubject }),
        getTeacherClasses(),
        getTeacherSubjects(),
      ])

      if (projRes.status === 'fulfilled' && projRes.value.success && projRes.value.data) {
        setProjects(projRes.value.data)
      } else if (projRes.status === 'fulfilled' && !projRes.value.success) {
        setError(projRes.value.error?.message || 'Gagal mengambil data projek.')
      }

      if (clsRes.status === 'fulfilled' && clsRes.value.success && clsRes.value.data) {
        setClasses(clsRes.value.data)
      }
      if (subRes.status === 'fulfilled' && subRes.value.success && subRes.value.data) {
        setSubjects(subRes.value.data)
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat projek.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [selectedClass, selectedSubject])

  const filteredProjects = useMemo(() => {
    let list = projects
    if (search.trim()) {
      const q = search.toLowerCase()
      list = list.filter(
        (p) =>
          p.title.toLowerCase().includes(q) ||
          (p.description && p.description.toLowerCase().includes(q))
      )
    }
    return list
  }, [projects, search])

  const totalPages = Math.ceil(filteredProjects.length / pageSize)
  const paginatedProjects = useMemo(() => {
    const start = (page - 1) * pageSize
    return filteredProjects.slice(start, start + pageSize)
  }, [filteredProjects, page, pageSize])

  const handleDeleteConfirm = async () => {
    if (!deleteTarget) return
    setIsDeleting(true)
    try {
      const res = await deleteProject(deleteTarget.id)
      if (res.success) {
        toast.success(
          'Projek Berhasil Diproses',
          res.data?.archived ? 'Projek telah diarsipkan karena memiliki submission.' : 'Projek telah dihapus.'
        )
        setDeleteTarget(null)
        fetchData()
      } else {
        toast.error('Gagal Menghapus Projek', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kesalahan sistem.'
      toast.error('Gagal Menghapus', msg)
    } finally {
      setIsDeleting(false)
    }
  }

  const columns: Column<Project>[] = [
    {
      key: 'title',
      header: 'Judul Projek',
      render: (p) => {
        const isPast = new Date(p.deadline) < new Date()
        return (
          <div className="max-w-xs sm:max-w-sm">
            <Link
              to={`/teacher/projects/${p.id}`}
              className="font-bold text-[#30251D] hover:text-[#FF9800] transition-colors line-clamp-1"
            >
              {p.title}
            </Link>
            {p.description && (
              <p className="mt-0.5 text-xs text-[#687385] line-clamp-1">{p.description}</p>
            )}
            <div className="mt-1 flex items-center gap-2">
              <StatusBadge
                status={isPast ? 'ended' : 'active'}
                label={isPast ? 'Berakhir' : 'Aktif'}
              />
              {p.allowResubmit && (
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
      render: (p) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {p.subjectName || 'Mata Pelajaran'}
        </span>
      ),
    },
    {
      key: 'classes',
      header: 'Kelas',
      render: (p) => (
        <div className="flex flex-wrap gap-1">
          {p.classNames && p.classNames.length > 0 ? (
            p.classNames.map((c, i) => (
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
      render: (p) => {
        const isPast = new Date(p.deadline) < new Date()
        return (
          <span
            className={`inline-flex items-center gap-1 text-xs font-semibold ${
              isPast ? 'text-rose-700' : 'text-[#30251D]'
            }`}
          >
            <Clock size={12} className={isPast ? 'text-rose-600' : 'text-[#FF9800]'} />
            <span>
              {new Date(p.deadline).toLocaleString('id-ID', {
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
      render: (p) => (
        <Link
          to={`/teacher/projects/${p.id}/submissions`}
          className="inline-flex items-center gap-1.5 rounded-lg bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-1 text-xs font-bold text-[#FF9800] hover:bg-[#FEEADC] transition-colors"
        >
          <Users size={12} />
          <span>
            {p.gradedSubmissions ?? 0}/{p.totalSubmissions ?? 0} Dinilai
          </span>
        </Link>
      ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'right',
      render: (p) => (
        <div className="flex items-center justify-end gap-1.5">
          <Link
            to={`/teacher/projects/${p.id}`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Lihat Detail"
          >
            <Eye size={14} />
          </Link>
          <Link
            to={`/teacher/projects/${p.id}/edit`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Edit Projek"
          >
            <Edit size={14} />
          </Link>
          <button
            type="button"
            onClick={() => setDeleteTarget(p)}
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
              Projek Pembelajaran
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Kelola projek kolaboratif atau karya mandiri siswa, tentukan panduan, dan beri evaluasi terperinci.
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
              to="/teacher/projects/new"
              className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs sm:text-sm font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
            >
              <Plus size={16} />
              <span>Buat Projek</span>
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
              placeholder="Cari projek berdasarkan judul atau deskripsi..."
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
          <ErrorState title="Gagal Memuat Projek" message={error} onRetry={fetchData} />
        ) : (
          <div className="space-y-4">
            <DataTable
              columns={columns}
              data={paginatedProjects}
              loading={loading}
              emptyMessage={search ? 'Projek tidak ditemukan' : 'Belum ada projek pembelajaran'}
              emptySubMessage={
                search
                  ? `Tidak ada hasil pencarian untuk "${search}"`
                  : 'Klik tombol "Buat Projek" untuk menugaskan projek ke kelas Anda.'
              }
            />

            <Pagination
              currentPage={page}
              totalPages={totalPages}
              onPageChange={setPage}
              totalItems={filteredProjects.length}
              pageSize={pageSize}
            />
          </div>
        )}

        {/* Delete Confirmation Dialog */}
        <ConfirmDialog
          isOpen={Boolean(deleteTarget)}
          onClose={() => setDeleteTarget(null)}
          onConfirm={handleDeleteConfirm}
          title={hasSubmissions ? 'Arsipkan Projek' : 'Hapus Projek'}
          message={`Apakah Anda yakin ingin menghapus atau mengarsipkan projek "${deleteTarget?.title}"?`}
          hasSubmissions={hasSubmissions}
          archiveNotice="Projek ini telah memiliki hasil karya yang dikumpulkan siswa. Data akan diarsipkan agar seluruh riwayat portofolio siswa tetap aman."
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
