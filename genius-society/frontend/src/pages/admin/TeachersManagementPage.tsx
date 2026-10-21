import { useState, useEffect, type FormEvent } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import {
  Users,
  Search,
  Plus,
  UserCheck,
  Download,
  Printer,
  ChevronLeft,
  ChevronRight,
  Filter,
  X,
  Lock,
  UserPlus,
  RefreshCw,
  Eye,
  Edit,
  Trash2,
  AlertCircle,
  BookOpen,
  Award
} from 'lucide-react'

interface Teacher {
  id: string
  userId?: string
  name: string
  username?: string
  email?: string
  nip: string
  gender?: string
  phone?: string
  address?: string
  subject?: string
  status: string
  isWalas: boolean
  walasClasses?: string[]
  hasAccount: boolean
  isActive: boolean
  createdAt: string
}

interface TeacherStats {
  totalTeachers: number
  activeTeachers: number
  walasCount: number
  teachersWithoutMapel: number
}

interface ClassOption {
  id: string
  name: string
}

export default function TeachersManagementPage() {
  const [teachers, setTeachers] = useState<Teacher[]>([])
  const [stats, setStats] = useState<TeacherStats>({
    totalTeachers: 0,
    activeTeachers: 0,
    walasCount: 0,
    teachersWithoutMapel: 0,
  })
  const [loading, setLoading] = useState(true)
  const [classes, setClasses] = useState<ClassOption[]>([])

  // Filters & Pagination
  const [search, setSearch] = useState('')
  const [subjectFilter, setSubjectFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [walasFilter, setWalasFilter] = useState('')
  const [page, setPage] = useState(1)
  const [limit] = useState(15)
  const [total, setTotal] = useState(0)

  // Modals
  const [createModal, setCreateModal] = useState(false)
  const [detailModal, setDetailModal] = useState<Teacher | null>(null)
  const [editModal, setEditModal] = useState<Teacher | null>(null)
  const [accountModal, setAccountModal] = useState<Teacher | null>(null)
  const [deleteModal, setDeleteModal] = useState<Teacher | null>(null)

  // Form states
  const [createForm, setCreateForm] = useState({
    name: '',
    nip: '',
    gender: 'L',
    phone: '',
    address: '',
    subject: '',
    status: 'active',
    isWalas: false,
    walasClassIds: [] as string[],
    createAccount: true,
    username: '',
    email: '',
    password: 'Password123!',
  })

  const [editForm, setEditForm] = useState({
    name: '',
    nip: '',
    gender: 'L',
    phone: '',
    address: '',
    subject: '',
    status: 'active',
    isWalas: false,
    walasClassIds: [] as string[],
  })

  const [accountForm, setAccountForm] = useState({
    username: '',
    email: '',
    password: '',
  })

  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // Fetch Stats
  const fetchStats = async () => {
    try {
      const res = await api.get<ApiResponse<TeacherStats>>('/teachers/stats')
      if (res.data.success && res.data.data) {
        setStats(res.data.data)
      }
    } catch (err) {
      console.error('Failed to load teacher stats:', err)
    }
  }

  // Fetch Classes for Walas dropdown
  const fetchClasses = async () => {
    try {
      const res = await api.get('/admin/dashboard')
      if (res.data?.data?.classDistribution) {
        setClasses(
          res.data.data.classDistribution.map((c: { classId: string; className: string }) => ({
            id: c.classId,
            name: c.className,
          }))
        )
      }
    } catch {
      setClasses([
        { id: '1', name: 'X IPA 1' },
        { id: '2', name: 'X IPA 2' },
        { id: '3', name: 'XI IPA 1' },
        { id: '4', name: 'XII IPA 1' },
      ])
    }
  }

  // Fetch Teachers List
  const fetchTeachers = async () => {
    setLoading(true)
    try {
      const params: Record<string, string | number | boolean> = {
        page,
        limit,
      }
      if (search.trim()) params.search = search.trim()
      if (subjectFilter) params.subject = subjectFilter
      if (statusFilter) params.status = statusFilter
      if (walasFilter !== '') params.isWalas = walasFilter === 'true'

      const res = await api.get<ApiResponse<Teacher[]>>('/teachers', { params })
      if (res.data.success && res.data.data) {
        setTeachers(res.data.data)
        if (res.data.meta) {
          setTotal((res.data.meta as { total: number }).total || 0)
        }
      }
    } catch (err) {
      console.error('Failed to load teachers:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStats()
    fetchClasses()
  }, [])

  useEffect(() => {
    fetchTeachers()
  }, [page, subjectFilter, statusFilter, walasFilter])

  const handleSearchSubmit = (e: FormEvent) => {
    e.preventDefault()
    setPage(1)
    fetchTeachers()
  }

  // Create Teacher
  const handleCreateSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setFormError(null)
    setSubmitting(true)
    try {
      const res = await api.post('/teachers', createForm)
      if (res.data.success) {
        setCreateModal(false)
        setCreateForm({
          name: '',
          nip: '',
          gender: 'L',
          phone: '',
          address: '',
          subject: '',
          status: 'active',
          isWalas: false,
          walasClassIds: [],
          createAccount: true,
          username: '',
          email: '',
          password: 'Password123!',
        })
        fetchTeachers()
        fetchStats()
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal menambahkan data guru')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Edit Teacher
  const handleEditSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!editModal) return
    setFormError(null)
    setSubmitting(true)
    try {
      await api.put(`/teachers/${editModal.id}`, editForm)
      setEditModal(null)
      fetchTeachers()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal memperbarui data guru')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Create Account
  const handleCreateAccountSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!accountModal) return
    setFormError(null)
    setSubmitting(true)
    try {
      await api.post(`/teachers/${accountModal.id}/account`, accountForm)
      setAccountModal(null)
      fetchTeachers()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal membuat akun guru')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Toggle Account Active / Inactive
  const handleToggleAccount = async (teacherId: string) => {
    try {
      await api.patch(`/teachers/${teacherId}/account/toggle`)
      fetchTeachers()
      fetchStats()
    } catch {
      alert('Gagal mengubah status akun guru')
    }
  }

  // Delete Teacher
  const handleDeleteTeacher = async () => {
    if (!deleteModal) return
    setSubmitting(true)
    try {
      await api.delete(`/teachers/${deleteModal.id}`)
      setDeleteModal(null)
      fetchTeachers()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert(`Gagal menghapus data guru: ${err.message}`)
      } else {
        alert('Gagal menghapus data guru')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Open Edit Modal
  const openEdit = (t: Teacher) => {
    setEditModal(t)
    setEditForm({
      name: t.name,
      nip: t.nip,
      gender: t.gender || 'L',
      phone: t.phone || '',
      address: t.address || '',
      subject: t.subject || '',
      status: t.status || 'active',
      isWalas: t.isWalas,
      walasClassIds: [],
    })
    setFormError(null)
  }

  // Export CSV
  const handleExportCSV = () => {
    if (teachers.length === 0) return
    const headers = ['NIP', 'Nama Guru', 'Mata Pelajaran', 'Wali Kelas', 'Status', 'Username', 'Email']
    const rows = teachers.map((t) => [
      `"${t.nip}"`,
      `"${t.name}"`,
      `"${t.subject || '-'}"`,
      `"${t.isWalas ? (t.walasClasses?.join(', ') || 'Ya') : 'Bukan'}"`,
      `"${t.status}"`,
      `"${t.username || '-'}"`,
      `"${t.email || '-'}"`,
    ])
    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map((e) => e.join(','))].join('\n')
    const encodedUri = encodeURI(csvContent)
    const link = document.createElement('a')
    link.setAttribute('href', encodedUri)
    link.setAttribute('download', `data-guru-${new Date().toISOString().slice(0, 10)}.csv`)
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
              Manajemen Guru
            </h1>
            <p className="mt-1 text-xs text-gray-500 font-medium">
              Kelola direktori pendidik, penugasan mata pelajaran, wali kelas, dan akun login
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
              <span>Tambah Guru</span>
            </button>
          </div>
        </div>

        {/* KPI Cards */}
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Total Guru</span>
              <div className="rounded-xl bg-orange-50 p-2 text-[#FF9800]">
                <Users size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-[#30251D]">{stats.totalTeachers}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Terdaftar di sistem</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Guru Aktif</span>
              <div className="rounded-xl bg-emerald-50 p-2 text-emerald-600">
                <UserCheck size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-[#30251D]">{stats.activeTeachers}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Status mengajar aktif</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Wali Kelas</span>
              <div className="rounded-xl bg-blue-50 p-2 text-blue-600">
                <Award size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-blue-600">{stats.walasCount}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Memegang amanah kelas</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold text-gray-500">Belum Ada Mapel</span>
              <div className="rounded-xl bg-amber-50 p-2 text-amber-600">
                <BookOpen size={16} />
              </div>
            </div>
            <div className="mt-2 text-2xl font-black text-amber-600">{stats.teachersWithoutMapel}</div>
            <p className="mt-0.5 text-[10px] text-gray-400 font-medium">Perlu penugasan kurikulum</p>
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
                placeholder="Cari nama atau NIP guru..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-10 pr-4 py-2 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:bg-white focus:outline-none transition-colors"
              />
            </form>

            {/* Filter Mata Pelajaran */}
            <div className="w-full md:w-48">
              <div className="relative">
                <BookOpen size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <select
                  value={subjectFilter}
                  onChange={(e) => {
                    setSubjectFilter(e.target.value)
                    setPage(1)
                  }}
                  className="w-full appearance-none rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-8 pr-8 py-2 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:bg-white focus:outline-none cursor-pointer"
                >
                  <option value="">Semua Mapel</option>
                  <option value="Matematika">Matematika</option>
                  <option value="Bahasa Indonesia">Bahasa Indonesia</option>
                  <option value="Bahasa Inggris">Bahasa Inggris</option>
                  <option value="Fisika">Fisika</option>
                  <option value="Biologi">Biologi</option>
                  <option value="Kimia">Kimia</option>
                  <option value="Sejarah">Sejarah</option>
                  <option value="Pendidikan Agama">Pendidikan Agama</option>
                </select>
              </div>
            </div>

            {/* Filter Wali Kelas */}
            <div className="w-full md:w-40">
              <div className="relative">
                <Award size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <select
                  value={walasFilter}
                  onChange={(e) => {
                    setWalasFilter(e.target.value)
                    setPage(1)
                  }}
                  className="w-full appearance-none rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-8 pr-8 py-2 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:bg-white focus:outline-none cursor-pointer"
                >
                  <option value="">Semua Guru</option>
                  <option value="true">Wali Kelas</option>
                  <option value="false">Bukan Walas</option>
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
                  <option value="leave">Cuti</option>
                </select>
              </div>
            </div>

            {/* Reset */}
            {(search || subjectFilter || statusFilter || walasFilter !== '') && (
              <button
                type="button"
                onClick={() => {
                  setSearch('')
                  setSubjectFilter('')
                  setStatusFilter('')
                  setWalasFilter('')
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

        {/* Tabel Data Guru */}
        <div className="overflow-hidden rounded-2xl border border-[#EBD4C5] bg-white shadow-sm">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600">
                <tr>
                  <th className="py-3 px-4 font-bold">NIP & Identitas Guru</th>
                  <th className="py-3 px-4 font-bold">Mata Pelajaran</th>
                  <th className="py-3 px-4 font-bold">Wali Kelas</th>
                  <th className="py-3 px-4 font-bold">Status</th>
                  <th className="py-3 px-4 font-bold">Akun Login</th>
                  <th className="py-3 px-4 font-bold text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#EBD4C5]/60 text-[#30251D]">
                {loading ? (
                  <tr>
                    <td colSpan={6} className="py-12 text-center text-gray-400">
                      <RefreshCw size={24} className="mx-auto mb-2 animate-spin text-[#FF9800]" />
                      <span>Memuat data guru...</span>
                    </td>
                  </tr>
                ) : teachers.length > 0 ? (
                  teachers.map((t) => {
                    return (
                      <tr key={t.id} className="hover:bg-[#FFF8F5]/50 transition-colors">
                        {/* NIP & Nama */}
                        <td className="py-3.5 px-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-orange-100 font-bold text-[#FF9800] text-xs">
                              {t.name
                                .split(' ')
                                .map((n) => n[0])
                                .slice(0, 2)
                                .join('')
                                .toUpperCase()}
                            </div>
                            <div className="min-w-0">
                              <p className="truncate font-bold text-[#30251D] hover:underline cursor-pointer" onClick={() => setDetailModal(t)}>
                                {t.name}
                              </p>
                              <p className="font-mono text-[11px] text-gray-400">
                                NIP: {t.nip} {t.gender && `• (${t.gender})`}
                              </p>
                            </div>
                          </div>
                        </td>

                        {/* Mata Pelajaran */}
                        <td className="py-3.5 px-4">
                          {t.subject ? (
                            <span className="inline-flex items-center gap-1 rounded-lg border border-[#EBD4C5] bg-[#FFF8F5] px-2.5 py-1 text-[11px] font-bold text-[#30251D]">
                              <BookOpen size={11} className="text-[#FF9800]" />
                              {t.subject}
                            </span>
                          ) : (
                            <span className="text-[11px] font-medium text-amber-600 bg-amber-50 px-2 py-0.5 rounded-md border border-amber-200">
                              Belum ditentukan
                            </span>
                          )}
                        </td>

                        {/* Wali Kelas */}
                        <td className="py-3.5 px-4">
                          {t.isWalas ? (
                            <div className="flex items-center gap-1.5">
                              <span className="inline-flex items-center gap-1 rounded-full bg-blue-100 px-2.5 py-0.5 text-[10px] font-bold text-blue-800">
                                <Award size={10} />
                                Walas
                              </span>
                              {t.walasClasses && t.walasClasses.length > 0 && (
                                <span className="font-semibold text-[11px] text-[#30251D]">
                                  ({t.walasClasses.join(', ')})
                                </span>
                              )}
                            </div>
                          ) : (
                            <span className="text-gray-400 text-[11px]">-</span>
                          )}
                        </td>

                        {/* Status */}
                        <td className="py-3.5 px-4">
                          <span
                            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-[10px] font-bold ${
                              t.status === 'active'
                                ? 'bg-emerald-100 text-emerald-800'
                                : t.status === 'leave'
                                ? 'bg-amber-100 text-amber-800'
                                : 'bg-gray-100 text-gray-600'
                            }`}
                          >
                            <span className="h-1.5 w-1.5 rounded-full bg-current" />
                            {t.status}
                          </span>
                        </td>

                        {/* Akun Guru */}
                        <td className="py-3.5 px-4">
                          {t.hasAccount ? (
                            <div className="flex items-center gap-2">
                              <div>
                                <span className="block font-mono text-[11px] font-bold text-[#30251D]">
                                  @{t.username}
                                </span>
                                <span className={`text-[10px] ${t.isActive ? 'text-emerald-600 font-semibold' : 'text-red-500'}`}>
                                  {t.isActive ? 'Akun Aktif' : 'Akun Nonaktif'}
                                </span>
                              </div>
                              <button
                                type="button"
                                onClick={() => handleToggleAccount(t.id)}
                                title={t.isActive ? 'Nonaktifkan Akun' : 'Aktifkan Akun'}
                                className="rounded p-1 text-gray-400 hover:text-[#FF9800]"
                              >
                                <Lock size={13} />
                              </button>
                            </div>
                          ) : (
                            <button
                              type="button"
                              onClick={() => {
                                setAccountModal(t)
                                setAccountForm({
                                  username: t.nip,
                                  email: `${t.nip}@guru.genius.ac.id`,
                                  password: 'Password123!',
                                })
                              }}
                              className="inline-flex items-center gap-1 rounded-lg border border-dashed border-[#FF9800] bg-orange-50/50 px-2 py-1 text-[11px] font-bold text-[#FF9800] hover:bg-[#FF9800] hover:text-white transition-colors"
                            >
                              <UserPlus size={12} />
                              <span>Buat Akun</span>
                            </button>
                          )}
                        </td>

                        {/* Aksi */}
                        <td className="py-3.5 px-4 text-right">
                          <div className="flex items-center justify-end gap-1.5">
                            <button
                              type="button"
                              onClick={() => setDetailModal(t)}
                              title="Lihat Profil Guru"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
                            >
                              <Eye size={15} />
                            </button>
                            <button
                              type="button"
                              onClick={() => openEdit(t)}
                              title="Edit Data Guru"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
                            >
                              <Edit size={15} />
                            </button>
                            <button
                              type="button"
                              onClick={() => setDeleteModal(t)}
                              title="Hapus Data Guru"
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
                    <td colSpan={6} className="py-12 text-center">
                      <Users size={32} className="mx-auto mb-2 text-gray-300" />
                      <p className="font-bold text-gray-600">Tidak ada data guru ditemukan</p>
                      <p className="text-xs text-gray-400">Coba ubah kata kunci pencarian atau filter Anda</p>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="flex flex-col sm:flex-row items-center justify-between border-t border-[#EBD4C5] px-4 py-3 text-xs text-gray-500 gap-2">
            <div>
              Menampilkan {teachers.length} dari total {total} guru
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
                disabled={teachers.length < limit}
                onClick={() => setPage((p) => p + 1)}
                className="inline-flex items-center gap-1 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 font-bold text-gray-600 shadow-sm hover:bg-[#FFF8F5] disabled:opacity-40"
              >
                <span>Berikutnya</span>
                <ChevronRight size={14} />
              </button>
            </div>
          </div>
        </div>

        {/* ================= MODAL TAMBAH GURU ================= */}
        {createModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="my-8 w-full max-w-2xl rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#FF9800] text-white">
                    <Plus size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Tambah Data Guru Baru</h3>
                    <p className="text-xs text-gray-500">Lengkapi formulir pendidik dan akun login</p>
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

              <form onSubmit={handleCreateSubmit} className="mt-4 space-y-4 text-xs">
                {/* Bagian 1: Identitas Guru */}
                <div className="rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/30 p-3.5 space-y-3">
                  <h4 className="font-bold text-[#FF9800] uppercase tracking-wider text-[11px]">
                    1. Identitas Guru
                  </h4>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div>
                      <label className="block mb-1 font-bold text-gray-700">Nama Lengkap & Gelar *</label>
                      <input
                        type="text"
                        required
                        placeholder="Contoh: Dra. Nurhaliza, M.Pd"
                        value={createForm.name}
                        onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                      />
                    </div>
                    <div>
                      <label className="block mb-1 font-bold text-gray-700">NIP (Nomor Induk Pegawai) *</label>
                      <input
                        type="text"
                        required
                        placeholder="Contoh: 198503122010012003"
                        value={createForm.nip}
                        onChange={(e) => setCreateForm({ ...createForm, nip: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                    <div>
                      <label className="block mb-1 font-bold text-gray-700">Jenis Kelamin</label>
                      <select
                        value={createForm.gender}
                        onChange={(e) => setCreateForm({ ...createForm, gender: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                      >
                        <option value="L">Laki-laki (L)</option>
                        <option value="P">Perempuan (P)</option>
                      </select>
                    </div>

                    <div>
                      <label className="block mb-1 font-bold text-gray-700">Mata Pelajaran Diampu</label>
                      <input
                        type="text"
                        placeholder="Contoh: Matematika"
                        value={createForm.subject}
                        onChange={(e) => setCreateForm({ ...createForm, subject: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                      />
                    </div>

                    <div>
                      <label className="block mb-1 font-bold text-gray-700">Status Kepegawaian</label>
                      <select
                        value={createForm.status}
                        onChange={(e) => setCreateForm({ ...createForm, status: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                      >
                        <option value="active">Aktif</option>
                        <option value="inactive">Nonaktif</option>
                        <option value="leave">Cuti</option>
                      </select>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div>
                      <label className="block mb-1 font-bold text-gray-700">Nomor Telepon / WhatsApp</label>
                      <input
                        type="text"
                        placeholder="Contoh: 081234567890"
                        value={createForm.phone}
                        onChange={(e) => setCreateForm({ ...createForm, phone: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                      />
                    </div>

                    <div>
                      <label className="block mb-1 font-bold text-gray-700">Alamat Tempat Tinggal</label>
                      <input
                        type="text"
                        placeholder="Contoh: Jl. Merdeka No. 45"
                        value={createForm.address}
                        onChange={(e) => setCreateForm({ ...createForm, address: e.target.value })}
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                      />
                    </div>
                  </div>

                  {/* Walas Checkbox */}
                  <div className="pt-2 border-t border-[#EBD4C5]/40">
                    <label className="flex items-center gap-2 cursor-pointer font-bold text-[#30251D]">
                      <input
                        type="checkbox"
                        checked={createForm.isWalas}
                        onChange={(e) => setCreateForm({ ...createForm, isWalas: e.target.checked })}
                        className="rounded border-[#EBD4C5] text-[#FF9800] focus:ring-[#FF9800]"
                      />
                      <span>Tugaskan sebagai Wali Kelas (Walas)</span>
                    </label>

                    {createForm.isWalas && (
                      <div className="mt-2 pl-6">
                        <label className="block mb-1 font-bold text-gray-700">Pilih Kelas Binaan</label>
                        <select
                          onChange={(e) => {
                            if (e.target.value && !createForm.walasClassIds.includes(e.target.value)) {
                              setCreateForm({
                                ...createForm,
                                walasClassIds: [...createForm.walasClassIds, e.target.value],
                              })
                            }
                          }}
                          className="w-full sm:w-64 rounded-xl border border-[#EBD4C5] bg-white p-2 text-xs focus:border-[#FF9800] focus:outline-none"
                        >
                          <option value="">+ Tambah Kelas Binaan</option>
                          {classes.map((c) => (
                            <option key={c.id} value={c.id}>
                              {c.name}
                            </option>
                          ))}
                        </select>
                      </div>
                    )}
                  </div>
                </div>

                {/* Bagian 2: Akun Login */}
                <div className="rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/30 p-3.5 space-y-3">
                  <div className="flex items-center justify-between">
                    <h4 className="font-bold text-[#FF9800] uppercase tracking-wider text-[11px]">
                      2. Akun Login Sistem
                    </h4>
                    <label className="flex items-center gap-1.5 cursor-pointer text-gray-700 font-bold">
                      <input
                        type="checkbox"
                        checked={createForm.createAccount}
                        onChange={(e) => setCreateForm({ ...createForm, createAccount: e.target.checked })}
                        className="rounded border-[#EBD4C5] text-[#FF9800] focus:ring-[#FF9800]"
                      />
                      <span>Buat Akun Sekaligus</span>
                    </label>
                  </div>

                  {createForm.createAccount && (
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-1">
                      <div>
                        <label className="block mb-1 font-bold text-gray-700">Username *</label>
                        <input
                          type="text"
                          required={createForm.createAccount}
                          placeholder={createForm.nip || 'guru_nip'}
                          value={createForm.username}
                          onChange={(e) => setCreateForm({ ...createForm, username: e.target.value })}
                          className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                        />
                      </div>
                      <div>
                        <label className="block mb-1 font-bold text-gray-700">Email Akun *</label>
                        <input
                          type="email"
                          required={createForm.createAccount}
                          placeholder="guru@genius.ac.id"
                          value={createForm.email}
                          onChange={(e) => setCreateForm({ ...createForm, email: e.target.value })}
                          className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                        />
                      </div>
                      <div>
                        <label className="block mb-1 font-bold text-gray-700">Password Awal *</label>
                        <input
                          type="text"
                          required={createForm.createAccount}
                          value={createForm.password}
                          onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })}
                          className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                        />
                      </div>
                    </div>
                  )}
                </div>

                <div className="flex justify-end gap-2.5 pt-3 border-t border-[#EBD4C5]/70">
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
                    {submitting ? 'Menyimpan...' : 'Simpan Data Guru'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ================= MODAL EDIT GURU ================= */}
        {editModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="my-8 w-full max-w-xl rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-100 text-[#FF9800]">
                    <Edit size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Edit Data Guru</h3>
                    <p className="text-xs text-gray-500">Perbarui informasi pengajar</p>
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

              <form onSubmit={handleEditSubmit} className="mt-4 space-y-3 text-xs">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama Lengkap & Gelar *</label>
                    <input
                      type="text"
                      required
                      value={editForm.name}
                      onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">NIP *</label>
                    <input
                      type="text"
                      required
                      value={editForm.nip}
                      onChange={(e) => setEditForm({ ...editForm, nip: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Jenis Kelamin</label>
                    <select
                      value={editForm.gender}
                      onChange={(e) => setEditForm({ ...editForm, gender: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    >
                      <option value="L">Laki-laki (L)</option>
                      <option value="P">Perempuan (P)</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Mata Pelajaran</label>
                    <input
                      type="text"
                      value={editForm.subject}
                      onChange={(e) => setEditForm({ ...editForm, subject: e.target.value })}
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
                      <option value="leave">Cuti</option>
                    </select>
                  </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Telepon / WhatsApp</label>
                    <input
                      type="text"
                      value={editForm.phone}
                      onChange={(e) => setEditForm({ ...editForm, phone: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Alamat</label>
                    <input
                      type="text"
                      value={editForm.address}
                      onChange={(e) => setEditForm({ ...editForm, address: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-white p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                    />
                  </div>
                </div>

                <div className="pt-2 border-t border-[#EBD4C5]/50">
                  <label className="flex items-center gap-2 cursor-pointer font-bold text-[#30251D]">
                    <input
                      type="checkbox"
                      checked={editForm.isWalas}
                      onChange={(e) => setEditForm({ ...editForm, isWalas: e.target.checked })}
                      className="rounded border-[#EBD4C5] text-[#FF9800] focus:ring-[#FF9800]"
                    />
                    <span>Tugaskan sebagai Wali Kelas (Walas)</span>
                  </label>
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
                    {submitting ? 'Menyimpan...' : 'Perbarui Guru'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ================= MODAL DETAIL PROFIL GURU ================= */}
        {detailModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="my-8 w-full max-w-lg rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-100 text-[#FF9800]">
                    <Eye size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Profil Lengkap Guru</h3>
                    <p className="text-xs text-gray-500">Detail identitas, tugas mengajar, dan akun</p>
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

              <div className="mt-4 space-y-4 text-xs">
                {/* Header Profile */}
                <div className="flex items-center gap-4 rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/50 p-4">
                  <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-[#FF9800] text-lg font-black text-white shadow-sm">
                    {detailModal.name
                      .split(' ')
                      .map((n) => n[0])
                      .slice(0, 2)
                      .join('')
                      .toUpperCase()}
                  </div>
                  <div>
                    <h4 className="text-sm font-black text-[#30251D]">{detailModal.name}</h4>
                    <p className="font-mono text-gray-500">NIP: {detailModal.nip}</p>
                    <div className="mt-1 flex items-center gap-2">
                      <span className="rounded-full bg-emerald-100 px-2 py-0.5 text-[10px] font-bold text-emerald-800">
                        {detailModal.status}
                      </span>
                      {detailModal.isWalas && (
                        <span className="rounded-full bg-blue-100 px-2 py-0.5 text-[10px] font-bold text-blue-800">
                          Wali Kelas
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                {/* Detail Data */}
                <div className="rounded-xl border border-[#EBD4C5]/60 p-4 space-y-2.5">
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Mata Pelajaran:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.subject || 'Belum Ditentukan'}</span>
                  </div>
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Jenis Kelamin:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.gender === 'P' ? 'Perempuan' : 'Laki-laki'}</span>
                  </div>
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">No. Telepon / WA:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.phone || '-'}</span>
                  </div>
                  <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                    <span className="text-gray-500">Alamat:</span>
                    <span className="font-bold text-[#30251D]">{detailModal.address || '-'}</span>
                  </div>
                  {detailModal.isWalas && (
                    <div className="flex justify-between border-b border-[#EBD4C5]/40 pb-2">
                      <span className="text-gray-500">Kelas Binaan:</span>
                      <span className="font-bold text-blue-700">{detailModal.walasClasses?.join(', ') || 'Semua Kelas Binaan'}</span>
                    </div>
                  )}
                  <div className="flex justify-between">
                    <span className="text-gray-500">Akun Login:</span>
                    <span className="font-bold text-[#30251D]">
                      {detailModal.hasAccount ? `@${detailModal.username} (${detailModal.email})` : 'Belum Ada Akun'}
                    </span>
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

        {/* ================= MODAL BUAT AKUN GURU ================= */}
        {accountModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b border-[#EBD4C5]/70 pb-3">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-100 text-[#FF9800]">
                    <UserPlus size={18} />
                  </div>
                  <div>
                    <h3 className="text-base font-black text-[#30251D]">Buat Akun Login Guru</h3>
                    <p className="text-xs text-gray-500">Untuk {accountModal.name}</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => setAccountModal(null)}
                  className="rounded-lg p-1 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              {formError && (
                <div className="mt-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-2.5 text-xs text-red-600">
                  <AlertCircle size={15} />
                  <span>{formError}</span>
                </div>
              )}

              <form onSubmit={handleCreateAccountSubmit} className="mt-4 space-y-3 text-xs">
                <div>
                  <label className="block mb-1 font-bold text-gray-700">Username *</label>
                  <input
                    type="text"
                    required
                    value={accountForm.username}
                    onChange={(e) => setAccountForm({ ...accountForm, username: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block mb-1 font-bold text-gray-700">Email Akun *</label>
                  <input
                    type="email"
                    required
                    value={accountForm.email}
                    onChange={(e) => setAccountForm({ ...accountForm, email: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                  />
                </div>
                <div>
                  <label className="block mb-1 font-bold text-gray-700">Password Awal *</label>
                  <input
                    type="text"
                    required
                    value={accountForm.password}
                    onChange={(e) => setAccountForm({ ...accountForm, password: e.target.value })}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                  />
                </div>

                <div className="flex justify-end gap-2 pt-3 border-t border-[#EBD4C5]/70">
                  <button
                    type="button"
                    onClick={() => setAccountModal(null)}
                    className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 font-bold text-gray-600"
                  >
                    Batal
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="rounded-xl bg-[#FF9800] px-5 py-2 font-bold text-white shadow-sm hover:bg-[#E68600] disabled:opacity-70"
                  >
                    {submitting ? 'Memproses...' : 'Buat Akun'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ================= MODAL KONFIRMASI HAPUS GURU ================= */}
        {deleteModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-red-100 text-red-600">
                  <Trash2 size={22} />
                </div>
                <div>
                  <h3 className="text-base font-black text-[#30251D]">Hapus Data Guru</h3>
                  <p className="text-xs text-gray-500">Tindakan ini permanen dan tidak dapat dibatalkan</p>
                </div>
              </div>

              <div className="mt-4 rounded-xl border border-red-100 bg-red-50/60 p-3.5 text-xs text-red-800 space-y-1.5">
                <p>
                  Apakah Anda yakin ingin menghapus data guru <strong className="font-bold text-red-900">{deleteModal.name}</strong> (NIP: {deleteModal.nip})?
                </p>
                {deleteModal.hasAccount && (
                  <p className="text-[11px] text-red-700 font-semibold">
                    ⚠️ Akun login guru terkait (@{deleteModal.username}) juga akan dihapus dari sistem.
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
                  onClick={handleDeleteTeacher}
                  className="inline-flex items-center gap-1.5 rounded-xl bg-red-600 px-5 py-2 text-xs font-bold text-white shadow-sm hover:bg-red-700 disabled:opacity-70 transition-colors"
                >
                  <Trash2 size={14} />
                  <span>{submitting ? 'Menghapus...' : 'Ya, Hapus Guru'}</span>
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
