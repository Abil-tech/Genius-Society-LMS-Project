import { useState, useEffect, type FormEvent } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import {
  BookOpen,
  Search,
  Plus,
  Download,
  Printer,
  ChevronLeft,
  ChevronRight,
  Filter,
  X,
  RefreshCw,
  Eye,
  Edit,
  Trash2,
  AlertCircle,
  Users,
  CheckCircle2,
  Bookmark
} from 'lucide-react'

interface Subject {
  id: string
  code: string
  name: string
  category: string
  gradeLevel: string
  kkm: number
  description: string
  status: string
  teacherCount: number
  createdAt: string
}

interface SubjectStats {
  totalSubjects: number
  activeSubjects: number
  generalSubjects: number
  specialtySubjects: number
}

export default function SubjectsManagementPage() {
  const [subjects, setSubjects] = useState<Subject[]>([])
  const [stats, setStats] = useState<SubjectStats>({
    totalSubjects: 0,
    activeSubjects: 0,
    generalSubjects: 0,
    specialtySubjects: 0,
  })
  const [loading, setLoading] = useState(true)

  // Filters & Pagination
  const [search, setSearch] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [gradeFilter, setGradeFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [limit] = useState(15)
  const [total, setTotal] = useState(0)

  // Modals
  const [createModal, setCreateModal] = useState(false)
  const [detailModal, setDetailModal] = useState<Subject | null>(null)
  const [editModal, setEditModal] = useState<Subject | null>(null)
  const [deleteModal, setDeleteModal] = useState<Subject | null>(null)

  // Form states
  const [createForm, setCreateForm] = useState({
    code: '',
    name: '',
    category: 'Umum',
    gradeLevel: 'Semua Tingkat',
    kkm: 75,
    description: '',
    status: 'active',
  })

  const [editForm, setEditForm] = useState({
    code: '',
    name: '',
    category: 'Umum',
    gradeLevel: 'Semua Tingkat',
    kkm: 75,
    description: '',
    status: 'active',
  })

  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // Fetch Stats
  const fetchStats = async () => {
    try {
      const res = await api.get<ApiResponse<SubjectStats>>('/subjects/stats')
      if (res.data.success && res.data.data) {
        setStats(res.data.data)
      }
    } catch (err) {
      console.error('Failed to load subject stats:', err)
    }
  }

  // Fetch Subjects List
  const fetchSubjects = async () => {
    setLoading(true)
    try {
      const params: Record<string, string | number> = {
        page,
        limit,
      }
      if (search.trim()) params.search = search.trim()
      if (categoryFilter) params.category = categoryFilter
      if (gradeFilter) params.gradeLevel = gradeFilter
      if (statusFilter) params.status = statusFilter

      const res = await api.get<ApiResponse<Subject[]>>('/subjects', { params })
      if (res.data.success && res.data.data) {
        setSubjects(res.data.data)
        if (res.data.meta) {
          setTotal((res.data.meta as { total: number }).total || 0)
        }
      }
    } catch (err) {
      console.error('Failed to load subjects:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStats()
  }, [])

  useEffect(() => {
    fetchSubjects()
  }, [page, categoryFilter, gradeFilter, statusFilter])

  const handleSearchSubmit = (e: FormEvent) => {
    e.preventDefault()
    setPage(1)
    fetchSubjects()
  }

  // Create Subject
  const handleCreateSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setFormError(null)
    setSubmitting(true)
    try {
      const res = await api.post('/subjects', createForm)
      if (res.data.success) {
        setCreateModal(false)
        setCreateForm({
          code: '',
          name: '',
          category: 'Umum',
          gradeLevel: 'Semua Tingkat',
          kkm: 75,
          description: '',
          status: 'active',
        })
        fetchSubjects()
        fetchStats()
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal menambahkan mata pelajaran')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Edit Subject
  const handleEditSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!editModal) return
    setFormError(null)
    setSubmitting(true)
    try {
      await api.put(`/subjects/${editModal.id}`, editForm)
      setEditModal(null)
      fetchSubjects()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal memperbarui data mata pelajaran')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Delete Subject
  const handleDeleteSubject = async () => {
    if (!deleteModal) return
    setSubmitting(true)
    try {
      await api.delete(`/subjects/${deleteModal.id}`)
      setDeleteModal(null)
      fetchSubjects()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert(`Gagal menghapus mata pelajaran: ${err.message}`)
      } else {
        alert('Gagal menghapus mata pelajaran')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Open Edit Modal
  const openEdit = (sub: Subject) => {
    setEditModal(sub)
    setEditForm({
      code: sub.code,
      name: sub.name,
      category: sub.category || 'Umum',
      gradeLevel: sub.gradeLevel || 'Semua Tingkat',
      kkm: sub.kkm || 75,
      description: sub.description || '',
      status: sub.status || 'active',
    })
    setFormError(null)
  }

  // Export CSV
  const handleExportCSV = () => {
    if (subjects.length === 0) return
    const headers = ['Kode Mapel', 'Nama Mata Pelajaran', 'Kategori', 'Tingkat', 'KKM', 'Jumlah Pengampu', 'Status']
    const rows = subjects.map((s) => [
      `"${s.code}"`,
      `"${s.name}"`,
      `"${s.category}"`,
      `"${s.gradeLevel}"`,
      `"${s.kkm}"`,
      `"${s.teacherCount}"`,
      `"${s.status}"`,
    ])
    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map((e) => e.join(','))].join('\n')
    const encodedUri = encodeURI(csvContent)
    const link = document.createElement('a')
    link.setAttribute('href', encodedUri)
    link.setAttribute('download', `data-mapel-${new Date().toISOString().slice(0, 10)}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  // Print
  const handlePrint = () => {
    window.print()
  }

  return (
    <AdminLayout>
      <div className="space-y-6">
        {/* Header Title & Actions */}
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-xl sm:text-2xl font-black text-[#30251D] tracking-tight">
              Manajemen Mata Pelajaran
            </h1>
            <p className="mt-1 text-xs text-gray-500 font-medium">
              Kelola kurikulum mata pelajaran, alokasi tingkat kelas, KKM, dan guru pengampu
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <button
              type="button"
              onClick={handleExportCSV}
              className="inline-flex items-center gap-2 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-gray-700 shadow-sm hover:bg-[#FFF8F5] transition-colors"
            >
              <Download size={14} className="text-gray-500" />
              <span>Export CSV</span>
            </button>
            <button
              type="button"
              onClick={handlePrint}
              className="inline-flex items-center gap-2 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-gray-700 shadow-sm hover:bg-[#FFF8F5] transition-colors"
            >
              <Printer size={14} className="text-gray-500" />
              <span>Cetak</span>
            </button>
            <button
              type="button"
              onClick={() => {
                setCreateModal(true)
                setFormError(null)
              }}
              className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white shadow-sm hover:bg-[#E68600] transition-colors"
            >
              <Plus size={16} />
              <span>Tambah Mapel</span>
            </button>
          </div>
        </div>

        {/* 4 KPI Cards */}
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Total Mapel</span>
              <div className="rounded-xl bg-orange-50 p-2 text-[#FF9800]">
                <BookOpen size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-[#30251D]">{stats.totalSubjects}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Terdaftar di kurikulum</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Mapel Aktif</span>
              <div className="rounded-xl bg-emerald-50 p-2 text-emerald-600">
                <CheckCircle2 size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-[#30251D]">{stats.activeSubjects}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Diajarkan semester ini</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Mapel Umum</span>
              <div className="rounded-xl bg-blue-50 p-2 text-blue-600">
                <Bookmark size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-blue-600">{stats.generalSubjects}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Muatan inti sekolah</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Peminatan / Mulok</span>
              <div className="rounded-xl bg-amber-50 p-2 text-amber-600">
                <Filter size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-amber-600">{stats.specialtySubjects}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Jurusan & kejuruan</p>
          </div>
        </div>

        {/* Filter & Search Toolbar */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm space-y-3">
          <div className="flex flex-col md:flex-row gap-3">
            {/* Search Input */}
            <form onSubmit={handleSearchSubmit} className="relative flex-1">
              <Search size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                type="text"
                placeholder="Cari kode atau nama mata pelajaran..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-10 pr-4 py-2 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:bg-white focus:outline-none transition-colors"
              />
            </form>

            {/* Filter Kategori */}
            <div className="w-full md:w-48">
              <div className="relative">
                <Bookmark size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <select
                  value={categoryFilter}
                  onChange={(e) => {
                    setCategoryFilter(e.target.value)
                    setPage(1)
                  }}
                  className="w-full appearance-none rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-8 pr-8 py-2 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:bg-white focus:outline-none cursor-pointer"
                >
                  <option value="">Semua Kategori</option>
                  <option value="Umum">Umum (Wajib)</option>
                  <option value="Peminatan IPA">Peminatan IPA</option>
                  <option value="Peminatan IPS">Peminatan IPS</option>
                  <option value="Muatan Lokal">Muatan Lokal</option>
                  
                </select>
              </div>
            </div>

            {/* Filter Tingkat */}
            <div className="w-full md:w-40">
              <div className="relative">
                <Filter size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <select
                  value={gradeFilter}
                  onChange={(e) => {
                    setGradeFilter(e.target.value)
                    setPage(1)
                  }}
                  className="w-full appearance-none rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-8 pr-8 py-2 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:bg-white focus:outline-none cursor-pointer"
                >
                  <option value="">Semua Tingkat</option>
                  <option value="Semua Tingkat">Semua Tingkat</option>
                  <option value="Kelas X">Kelas X</option>
                  <option value="Kelas XI">Kelas XI</option>
                  <option value="Kelas XII">Kelas XII</option>
                </select>
              </div>
            </div>

            {/* Filter Status */}
            <div className="w-full md:w-36">
              <div className="relative">
                <Filter size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <select
                  value={statusFilter}
                  onChange={(e) => {
                    setStatusFilter(e.target.value)
                    setPage(1)
                  }}
                  className="w-full appearance-none rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-8 pr-8 py-2 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:bg-white focus:outline-none cursor-pointer"
                >
                  <option value="">Semua Status</option>
                  <option value="active">Aktif</option>
                  <option value="inactive">Nonaktif</option>
                </select>
              </div>
            </div>

            {/* Reset */}
            {(search || categoryFilter || gradeFilter || statusFilter) && (
              <button
                type="button"
                onClick={() => {
                  setSearch('')
                  setCategoryFilter('')
                  setGradeFilter('')
                  setStatusFilter('')
                  setPage(1)
                }}
                className="inline-flex items-center gap-1.5 rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-xs font-bold text-gray-600 hover:bg-gray-100"
              >
                <RefreshCw size={13} />
                <span>Reset</span>
              </button>
            )}
          </div>
        </div>

        {/* Tabel Data Mata Pelajaran */}
        <div className="overflow-hidden rounded-2xl border border-[#EBD4C5] bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600">
                <tr>
                  <th className="py-3 px-4 font-bold">Kode & Nama Mata Pelajaran</th>
                  <th className="py-3 px-4 font-bold">Kategori</th>
                  <th className="py-3 px-4 font-bold">Tingkat Kelas</th>
                  <th className="py-3 px-4 font-bold">KKM</th>
                  <th className="py-3 px-4 font-bold">Guru Pengampu</th>
                  <th className="py-3 px-4 font-bold">Status</th>
                  <th className="py-3 px-4 font-bold text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#EBD4C5]/60 text-[#30251D]">
                {loading ? (
                  <tr>
                    <td colSpan={7} className="py-12 text-center text-gray-400">
                      <RefreshCw size={24} className="mx-auto mb-2 animate-spin text-[#FF9800]" />
                      <span>Memuat data mata pelajaran...</span>
                    </td>
                  </tr>
                ) : subjects.length > 0 ? (
                  subjects.map((sub) => {
                    return (
                      <tr key={sub.id} className="hover:bg-[#FFF8F5]/50 transition-colors">
                        {/* Kode & Nama */}
                        <td className="py-3.5 px-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-orange-100 font-mono font-bold text-[#FF9800] text-xs">
                              {sub.code.slice(0, 3).toUpperCase()}
                            </div>
                            <div className="min-w-0">
                              <p
                                className="truncate font-bold text-[#30251D] hover:underline cursor-pointer"
                                onClick={() => setDetailModal(sub)}
                              >
                                {sub.name}
                              </p>
                              <p className="font-mono text-[11px] text-gray-400">Kode: {sub.code}</p>
                            </div>
                          </div>
                        </td>

                        {/* Kategori */}
                        <td className="py-3.5 px-4">
                          <span
                            className={`inline-flex items-center gap-1 rounded-lg px-2.5 py-1 text-[11px] font-bold ${
                              sub.category === 'Umum'
                                ? 'bg-blue-50 text-blue-700 border border-blue-200'
                                : sub.category.includes('IPA')
                                ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                                : sub.category.includes('IPS')
                                ? 'bg-amber-50 text-amber-700 border border-amber-200'
                                : 'bg-purple-50 text-purple-700 border border-purple-200'
                            }`}
                          >
                            <Bookmark size={11} />
                            {sub.category}
                          </span>
                        </td>

                        {/* Tingkat */}
                        <td className="py-3.5 px-4">
                          <span className="font-semibold text-gray-700">{sub.gradeLevel || 'Semua'}</span>
                        </td>

                        {/* KKM */}
                        <td className="py-3.5 px-4">
                          <span className="inline-flex items-center justify-center rounded-lg bg-orange-50 px-2.5 py-1 font-mono font-bold text-[#FF9800] border border-orange-200">
                            {sub.kkm}
                          </span>
                        </td>

                        {/* Guru Pengampu */}
                        <td className="py-3.5 px-4">
                          <div className="flex items-center gap-1.5">
                            <Users size={14} className="text-gray-400" />
                            <span className="font-bold text-[#30251D]">{sub.teacherCount} Guru</span>
                          </div>
                        </td>

                        {/* Status */}
                        <td className="py-3.5 px-4">
                          <span
                            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-[10px] font-bold ${
                              sub.status === 'active'
                                ? 'bg-emerald-100 text-emerald-800'
                                : 'bg-gray-100 text-gray-600'
                            }`}
                          >
                            <span className="h-1.5 w-1.5 rounded-full bg-current" />
                            {sub.status === 'active' ? 'Aktif' : 'Nonaktif'}
                          </span>
                        </td>

                        {/* Aksi */}
                        <td className="py-3.5 px-4 text-right">
                          <div className="flex items-center justify-end gap-1.5">
                            <button
                              type="button"
                              onClick={() => setDetailModal(sub)}
                              title="Lihat Detail Mapel"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
                            >
                              <Eye size={15} />
                            </button>
                            <button
                              type="button"
                              onClick={() => openEdit(sub)}
                              title="Edit Mapel"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
                            >
                              <Edit size={15} />
                            </button>
                            <button
                              type="button"
                              onClick={() => setDeleteModal(sub)}
                              title="Hapus Mapel"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-red-50 hover:text-red-600 transition-colors"
                            >
                              <Trash2 size={15} />
                            </button>
                          </div>
                        </td>
                      </tr>
                    )
                  })
                ) : (
                  <tr>
                    <td colSpan={7} className="py-12 text-center">
                      <BookOpen size={32} className="mx-auto mb-2 text-gray-300" />
                      <p className="font-bold text-gray-600">Tidak ada data mata pelajaran ditemukan</p>
                      <p className="text-xs text-gray-400">Coba ubah filter atau tambahkan mata pelajaran baru</p>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="flex flex-col sm:flex-row items-center justify-between border-t border-[#EBD4C5] px-4 py-3 text-xs text-gray-500 gap-2">
            <div>
              Menampilkan {subjects.length} dari total {total} mata pelajaran
            </div>
            <div className="flex items-center gap-2">
              <button
                type="button"
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                className="inline-flex items-center gap-1 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 font-bold text-gray-600 shadow-sm hover:bg-[#FFF8F5] disabled:opacity-40"
              >
                <ChevronLeft size={14} />
                <span>Sebelumnya</span>
              </button>
              <span className="font-bold text-[#30251D] px-2">Halaman {page}</span>
              <button
                type="button"
                disabled={subjects.length < limit}
                onClick={() => setPage((p) => p + 1)}
                className="inline-flex items-center gap-1 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 font-bold text-gray-600 shadow-sm hover:bg-[#FFF8F5] disabled:opacity-40"
              >
                <span>Berikutnya</span>
                <ChevronRight size={14} />
              </button>
            </div>
          </div>
        </div>

        {/* ================= MODAL TAMBAH MAPEL ================= */}
        {createModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="my-8 w-full max-w-lg rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#FF9800] text-white">
                    <Plus size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Tambah Mata Pelajaran</h3>
                    <p className="text-xs text-gray-500">Daftarkan kurikulum mata pelajaran baru</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => setCreateModal(false)}
                  className="rounded-lg p-1 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              {formError && (
                <div className="mt-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                  <AlertCircle size={16} />
                  <span>{formError}</span>
                </div>
              )}

              <form onSubmit={handleCreateSubmit} className="mt-4 space-y-3.5 text-xs">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kode Mata Pelajaran *</label>
                    <input
                      type="text"
                      required
                      placeholder="Contoh: MTK-WJB"
                      value={createForm.code}
                      onChange={(e) => setCreateForm({ ...createForm, code: e.target.value.toUpperCase() })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono uppercase focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama Mata Pelajaran *</label>
                    <input
                      type="text"
                      required
                      placeholder="Contoh: Matematika Wajib"
                      value={createForm.name}
                      onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kategori</label>
                    <select
                      value={createForm.category}
                      onChange={(e) => setCreateForm({ ...createForm, category: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    >
                      <option value="Umum">Umum (Wajib)</option>
                      <option value="Peminatan IPA">Peminatan IPA</option>
                      <option value="Peminatan IPS">Peminatan IPS</option>
                      <option value="Muatan Lokal">Muatan Lokal</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tingkat Kelas</label>
                    <select
                      value={createForm.gradeLevel}
                      onChange={(e) => setCreateForm({ ...createForm, gradeLevel: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    >
                      <option value="Semua Tingkat">Semua Tingkat</option>
                      <option value="Kelas X">Kelas X</option>
                      <option value="Kelas XI">Kelas XI</option>
                      <option value="Kelas XII">Kelas XII</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nilai KKM (1-100)</label>
                    <input
                      type="number"
                      min={0}
                      max={100}
                      value={createForm.kkm}
                      onChange={(e) => setCreateForm({ ...createForm, kkm: parseInt(e.target.value) || 0 })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                </div>

                <div>
                  <label className="block mb-1 font-bold text-gray-700">Deskripsi / Kompetensi Inti</label>
                  <textarea
                    rows={3}
                    placeholder="Contoh: Fokus pada pemahaman konsep aljabar, kalkulus dasar, dan geometri analitik."
                    value={createForm.description}
                    onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                  />
                </div>

                <div>
                  <label className="block mb-1 font-bold text-gray-700">Status Mata Pelajaran</label>
                  <select
                    value={createForm.status}
                    onChange={(e) => setCreateForm({ ...createForm, status: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                  >
                    <option value="active">Aktif</option>
                    <option value="inactive">Nonaktif</option>
                  </select>
                </div>

                <div className="flex justify-end gap-2.5 pt-4 border-t border-[#EBD4C5]/70">
                  <button
                    type="button"
                    onClick={() => setCreateModal(false)}
                    className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 font-bold text-gray-600 hover:bg-[#FFF8F5]"
                  >
                    Batal
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="rounded-xl bg-[#FF9800] px-5 py-2 font-bold text-white shadow-sm hover:bg-[#E68600] disabled:opacity-70 transition-colors"
                  >
                    {submitting ? 'Menyimpan...' : 'Simpan Mapel'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ================= MODAL EDIT MAPEL ================= */}
        {editModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="my-8 w-full max-w-lg rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-100 text-[#FF9800]">
                    <Edit size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Edit Mata Pelajaran</h3>
                    <p className="text-xs text-gray-500">Perbarui rincian kurikulum mapel</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => setEditModal(null)}
                  className="rounded-lg p-1 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              {formError && (
                <div className="mt-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                  <AlertCircle size={16} />
                  <span>{formError}</span>
                </div>
              )}

              <form onSubmit={handleEditSubmit} className="mt-4 space-y-3.5 text-xs">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kode Mata Pelajaran *</label>
                    <input
                      type="text"
                      required
                      value={editForm.code}
                      onChange={(e) => setEditForm({ ...editForm, code: e.target.value.toUpperCase() })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono uppercase focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama Mata Pelajaran *</label>
                    <input
                      type="text"
                      required
                      value={editForm.name}
                      onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kategori</label>
                    <select
                      value={editForm.category}
                      onChange={(e) => setEditForm({ ...editForm, category: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    >
                      <option value="Umum">Umum (Wajib)</option>
                      <option value="Peminatan IPA">Peminatan IPA</option>
                      <option value="Peminatan IPS">Peminatan IPS</option>
                      <option value="Muatan Lokal">Muatan Lokal</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tingkat Kelas</label>
                    <select
                      value={editForm.gradeLevel}
                      onChange={(e) => setEditForm({ ...editForm, gradeLevel: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    >
                      <option value="Semua Tingkat">Semua Tingkat</option>
                      <option value="Kelas X">Kelas X</option>
                      <option value="Kelas XI">Kelas XI</option>
                      <option value="Kelas XII">Kelas XII</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nilai KKM</label>
                    <input
                      type="number"
                      min={0}
                      max={100}
                      value={editForm.kkm}
                      onChange={(e) => setEditForm({ ...editForm, kkm: parseInt(e.target.value) || 0 })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                </div>

                <div>
                  <label className="block mb-1 font-bold text-gray-700">Deskripsi / Silabus</label>
                  <textarea
                    rows={3}
                    value={editForm.description}
                    onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                  />
                </div>

                <div>
                  <label className="block mb-1 font-bold text-gray-700">Status</label>
                  <select
                    value={editForm.status}
                    onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                  >
                    <option value="active">Aktif</option>
                    <option value="inactive">Nonaktif</option>
                  </select>
                </div>

                <div className="flex justify-end gap-2.5 pt-4 border-t border-[#EBD4C5]/70">
                  <button
                    type="button"
                    onClick={() => setEditModal(null)}
                    className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 font-bold text-gray-600 hover:bg-[#FFF8F5]"
                  >
                    Batal
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="rounded-xl bg-[#FF9800] px-5 py-2 font-bold text-white shadow-sm hover:bg-[#E68600] disabled:opacity-70"
                  >
                    {submitting ? 'Menyimpan...' : 'Perbarui Mapel'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ================= MODAL DETAIL MAPEL ================= */}
        {detailModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="my-8 w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-100 text-[#FF9800]">
                    <Eye size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Rincian Mata Pelajaran</h3>
                    <p className="text-xs text-gray-500">{detailModal.code}</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => setDetailModal(null)}
                  className="rounded-lg p-1 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              <div className="mt-4 space-y-3.5 text-xs">
                <div className="rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/50 p-4">
                  <span className="font-mono text-[11px] font-bold text-[#FF9800]">{detailModal.code}</span>
                  <h4 className="text-base font-black text-[#30251D] mt-0.5">{detailModal.name}</h4>
                  <p className="text-xs text-gray-500 mt-1">{detailModal.description || 'Tidak ada deskripsi silabus'}</p>
                </div>

                <div className="rounded-xl border border-[#EBD4C5]/60 p-4 space-y-2">
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Kategori:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.category}</span>
                  </div>
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Tingkat Kelas:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.gradeLevel}</span>
                  </div>
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Kriteria Ketuntasan (KKM):</span>
                    <span className="font-bold text-[#FF9800]">{detailModal.kkm}</span>
                  </div>
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Jumlah Guru Pengampu:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.teacherCount} Guru</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-500">Status Kurikulum:</span>
                    <span className="font-bold text-emerald-600">{detailModal.status === 'active' ? 'Aktif' : 'Nonaktif'}</span>
                  </div>
                </div>

                <div className="flex justify-end pt-2">
                  <button
                    type="button"
                    onClick={() => setDetailModal(null)}
                    className="rounded-xl bg-[#30251D] px-5 py-2 font-bold text-white hover:bg-black transition-colors"
                  >
                    Tutup
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* ================= MODAL KONFIRMASI HAPUS MAPEL ================= */}
        {deleteModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-red-100 text-red-600">
                  <Trash2 size={22} />
                </div>
                <div>
                  <h3 className="text-base font-black text-[#30251D]">Hapus Mata Pelajaran</h3>
                  <p className="text-xs text-gray-500">Tindakan ini tidak dapat dibatalkan</p>
                </div>
              </div>

              <div className="mt-4 rounded-xl border border-red-100 bg-red-50/60 p-3.5 text-xs text-red-800 space-y-1.5">
                <p>
                  Apakah Anda yakin ingin menghapus mata pelajaran <strong className="font-bold text-red-900">{deleteModal.name}</strong> ({deleteModal.code})?
                </p>
                {deleteModal.teacherCount > 0 && (
                  <p className="text-[11px] text-red-700 font-semibold">
                    ⚠️ Mata pelajaran ini masih diampu oleh {deleteModal.teacherCount} guru. Anda harus melepaskan penugasan guru terlebih dahulu.
                  </p>
                )}
              </div>

              <div className="mt-6 flex justify-end gap-2.5">
                <button
                  type="button"
                  disabled={submitting}
                  onClick={() => setDeleteModal(null)}
                  className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 text-xs font-bold text-gray-600 hover:bg-[#FFF8F5]"
                >
                  Batal
                </button>
                <button
                  type="button"
                  disabled={submitting}
                  onClick={handleDeleteSubject}
                  className="inline-flex items-center gap-1.5 rounded-xl bg-red-600 px-5 py-2 text-xs font-bold text-white shadow-sm hover:bg-red-700 disabled:opacity-70 transition-colors"
                >
                  <Trash2 size={14} />
                  <span>{submitting ? 'Menghapus...' : 'Ya, Hapus Mapel'}</span>
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
