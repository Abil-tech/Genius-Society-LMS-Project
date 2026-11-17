import { useEffect, useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import {
  Plus,
  Search,
  RefreshCw,
  FileText,
  ExternalLink,
  Trash2,
  Edit,
  Eye,
  Calendar,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import Pagination from '../../components/common/Pagination'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  getMaterials,
  deleteMaterial,
  getTeacherClasses,
  getTeacherSubjects,
} from '../../services/teacherService'
import type { Material, TeacherClass, TeacherSubject } from '../../types/teacher'

export default function TeacherMaterialsPage() {
  const toast = useToast()
  const [materials, setMaterials] = useState<Material[]>([])
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
  const [deleteTarget, setDeleteTarget] = useState<Material | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [matRes, clsRes, subRes] = await Promise.allSettled([
        getMaterials({ search, classId: selectedClass, subjectId: selectedSubject }),
        getTeacherClasses(),
        getTeacherSubjects(),
      ])

      if (matRes.status === 'fulfilled' && matRes.value.success && matRes.value.data) {
        setMaterials(matRes.value.data)
      } else if (matRes.status === 'fulfilled' && !matRes.value.success) {
        setError(matRes.value.error?.message || 'Gagal mengambil data materi.')
      }

      if (clsRes.status === 'fulfilled' && clsRes.value.success && clsRes.value.data) {
        setClasses(clsRes.value.data)
      }
      if (subRes.status === 'fulfilled' && subRes.value.success && subRes.value.data) {
        setSubjects(subRes.value.data)
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat materi.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [selectedClass, selectedSubject])

  const filteredMaterials = useMemo(() => {
    let list = materials
    if (search.trim()) {
      const q = search.toLowerCase()
      list = list.filter(
        (m) =>
          m.title.toLowerCase().includes(q) ||
          (m.description && m.description.toLowerCase().includes(q))
      )
    }
    return list
  }, [materials, search])

  const totalPages = Math.ceil(filteredMaterials.length / pageSize)
  const paginatedMaterials = useMemo(() => {
    const start = (page - 1) * pageSize
    return filteredMaterials.slice(start, start + pageSize)
  }, [filteredMaterials, page, pageSize])

  const handleDeleteConfirm = async () => {
    if (!deleteTarget) return
    setIsDeleting(true)
    try {
      const res = await deleteMaterial(deleteTarget.id)
      if (res.success) {
        toast.success(
          'Materi Berhasil Dihapus',
          res.data?.archived ? 'Materi telah diarsipkan.' : 'Materi telah dihapus.'
        )
        setDeleteTarget(null)
        fetchData()
      } else {
        toast.error('Gagal Menghapus Materi', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kesalahan sistem.'
      toast.error('Gagal Menghapus', msg)
    } finally {
      setIsDeleting(false)
    }
  }

  const columns: Column<Material>[] = [
    {
      key: 'title',
      header: 'Judul Materi',
      render: (m) => (
        <div className="max-w-xs sm:max-w-sm">
          <Link
            to={`/teacher/materials/${m.id}`}
            className="font-bold text-[#30251D] hover:text-[#FF9800] transition-colors line-clamp-1"
          >
            {m.title}
          </Link>
          {m.description && (
            <p className="mt-0.5 text-xs text-[#687385] line-clamp-1">{m.description}</p>
          )}
          <div className="mt-1 flex items-center gap-2 text-[11px] text-[#687385]">
            {m.file && (
              <span className="inline-flex items-center gap-1 text-emerald-700 bg-emerald-50 px-1.5 py-0.5 rounded font-semibold">
                <FileText size={11} />
                <span>File Lampiran</span>
              </span>
            )}
            {m.link && (
              <span className="inline-flex items-center gap-1 text-blue-700 bg-blue-50 px-1.5 py-0.5 rounded font-semibold">
                <ExternalLink size={11} />
                <span>Tautan Web</span>
              </span>
            )}
          </div>
        </div>
      ),
    },
    {
      key: 'subject',
      header: 'Mata Pelajaran',
      render: (m) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {m.subjectName || 'Mata Pelajaran'}
        </span>
      ),
    },
    {
      key: 'classes',
      header: 'Kelas Dituju',
      render: (m) => (
        <div className="flex flex-wrap gap-1">
          {m.classNames && m.classNames.length > 0 ? (
            m.classNames.map((c, i) => (
              <span
                key={i}
                className="rounded-md bg-[#FFF8F5] border border-[#EBD4C5] px-2 py-0.5 text-[11px] font-bold text-[#30251D]"
              >
                {c}
              </span>
            ))
          ) : (
            <span className="text-gray-400 text-xs">Semua Kelas</span>
          )}
        </div>
      ),
    },
    {
      key: 'publishAt',
      header: 'Tanggal Terbit',
      render: (m) => (
        <span className="inline-flex items-center gap-1 text-xs text-[#687385]">
          <Calendar size={12} className="text-[#FF9800]" />
          <span>
            {new Date(m.publishAt).toLocaleDateString('id-ID', {
              day: 'numeric',
              month: 'short',
              year: 'numeric',
            })}
          </span>
        </span>
      ),
    },
    {
      key: 'actions',
      header: 'Aksi',
      align: 'right',
      render: (m) => (
        <div className="flex items-center justify-end gap-1.5">
          <Link
            to={`/teacher/materials/${m.id}`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Lihat Detail"
          >
            <Eye size={14} />
          </Link>
          <Link
            to={`/teacher/materials/${m.id}/edit`}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
            title="Edit Materi"
          >
            <Edit size={14} />
          </Link>
          <button
            type="button"
            onClick={() => setDeleteTarget(m)}
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-rose-200 text-rose-600 hover:bg-rose-50 transition-colors"
            title="Hapus Materi"
          >
            <Trash2 size={14} />
          </button>
        </div>
      ),
    },
  ]

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Materi Pembelajaran
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Bagikan bahan ajar, modul digital, file presentasi, dan referensi bacaan ke siswa.
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
              to="/teacher/materials/new"
              className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs sm:text-sm font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
            >
              <Plus size={16} />
              <span>Buat Materi</span>
            </Link>
          </div>
        </div>

        {/* Filter bar */}
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
              placeholder="Cari materi berdasarkan judul atau deskripsi..."
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
          <ErrorState title="Gagal Memuat Materi" message={error} onRetry={fetchData} />
        ) : (
          <div className="space-y-4">
            <DataTable
              columns={columns}
              data={paginatedMaterials}
              loading={loading}
              emptyMessage={search ? 'Materi tidak ditemukan' : 'Belum ada materi pembelajaran'}
              emptySubMessage={
                search
                  ? `Tidak ada hasil pencarian untuk "${search}"`
                  : 'Klik tombol "Buat Materi" untuk mulai membagikan materi ke siswa.'
              }
            />

            <Pagination
              currentPage={page}
              totalPages={totalPages}
              onPageChange={setPage}
              totalItems={filteredMaterials.length}
              pageSize={pageSize}
            />
          </div>
        )}

        {/* Delete Confirmation Dialog */}
        <ConfirmDialog
          isOpen={Boolean(deleteTarget)}
          onClose={() => setDeleteTarget(null)}
          onConfirm={handleDeleteConfirm}
          title="Hapus Materi Pembelajaran"
          message={`Apakah Anda yakin ingin menghapus materi "${deleteTarget?.title}"?`}
          confirmText="Hapus Materi"
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
