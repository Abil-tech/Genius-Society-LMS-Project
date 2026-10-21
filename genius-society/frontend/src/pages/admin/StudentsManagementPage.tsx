import { useState, useEffect, type FormEvent } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import {
  GraduationCap,
  Search,
  Plus,
  UserCheck,
  UserX,
  Download,
  Printer,
  ChevronLeft,
  ChevronRight,
  Filter,
  Layers,
  X,
  CheckSquare,
  Square,
  Lock,
  UserPlus,
  RefreshCw,
  Eye,
  Edit,
  Trash2,
  AlertCircle
} from 'lucide-react'

interface Student {
  id: string
  userId?: string
  name: string
  username?: string
  email?: string
  nis: string
  nisn?: string
  gender?: string
  birthPlace?: string
  birthDate?: string
  address?: string
  parentName?: string
  parentPhone?: string
  enrollmentYear?: number
  classId?: string
  className?: string
  status: string
  hasAccount: boolean
  isActive: boolean
  createdAt: string
}

interface StudentStats {
  totalStudents: number
  activeStudents: number
  studentsWithoutClass: number
  studentsWithoutAccount: number
}

interface ClassOption {
  id: string
  name: string
}

export default function StudentsManagementPage() {
  const [students, setStudents] = useState<Student[]>([])
  const [stats, setStats] = useState<StudentStats>({
    totalStudents: 0,
    activeStudents: 0,
    studentsWithoutClass: 0,
    studentsWithoutAccount: 0,
  })
  const [loading, setLoading] = useState(true)
  const [classes, setClasses] = useState<ClassOption[]>([])

  // Search, Filter, Pagination
  const [search, setSearch] = useState('')
  const [classFilter, setClassFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [limit] = useState(15)
  const [total, setTotal] = useState(0)

  // Selection & Bulk
  const [selectedIds, setSelectedIds] = useState<string[]>([])
  const [bulkClassModal, setBulkClassModal] = useState(false)
  const [selectedBulkClass, setSelectedBulkClass] = useState('')

  // Modals
  const [createModal, setCreateModal] = useState(false)
  const [detailModal, setDetailModal] = useState<Student | null>(null)
  const [editModal, setEditModal] = useState<Student | null>(null)
  const [accountModal, setAccountModal] = useState<Student | null>(null)
  const [deleteModal, setDeleteModal] = useState<Student | null>(null)

  // Form states
  const [createForm, setCreateForm] = useState({
    name: '',
    nis: '',
    nisn: '',
    gender: 'L',
    birthPlace: '',
    birthDate: '',
    address: '',
    parentName: '',
    parentPhone: '',
    enrollmentYear: new Date().getFullYear(),
    classId: '',
    status: 'active',
    createAccount: true,
    username: '',
    email: '',
    password: 'Password123!',
  })

  const [editForm, setEditForm] = useState({
    name: '',
    nis: '',
    nisn: '',
    gender: 'L',
    birthPlace: '',
    birthDate: '',
    address: '',
    parentName: '',
    parentPhone: '',
    enrollmentYear: new Date().getFullYear(),
    classId: '',
    status: 'active',
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
      const res = await api.get<ApiResponse<StudentStats>>('/students/stats')
      if (res.data.success && res.data.data) {
        setStats(res.data.data)
      }
    } catch {
      // ignore
    }
  }

  // Fetch Classes for dropdown
  const fetchClasses = async () => {
    try {
      const res = await api.get<ApiResponse<{ classDistribution: Array<{ classId: string; className: string }> }>>('/admin/dashboard')
      if (res.data.success && res.data.data?.classDistribution) {
        setClasses(
          res.data.data.classDistribution.map((c) => ({
            id: c.classId,
            name: c.className,
          }))
        )
      }
    } catch {
      // fallback
      setClasses([
        { id: 'c1', name: 'Kelas X-A' },
        { id: 'c2', name: 'Kelas X-B' },
        { id: 'c3', name: 'Kelas XI-IPA 1' },
      ])
    }
  }

  // Fetch Student List
  const fetchStudents = async () => {
    setLoading(true)
    try {
      const params: Record<string, string | number> = {
        page,
        limit,
      }
      if (search) params.search = search
      if (classFilter) params.classId = classFilter
      if (statusFilter) params.status = statusFilter

      const res = await api.get<ApiResponse<Student[]>>('/students', { params })
      if (res.data.success && res.data.data) {
        setStudents(res.data.data)
        const meta = (res.data as unknown as { meta?: { total?: number } }).meta
        if (meta?.total !== undefined) {
          setTotal(meta.total)
        } else {
          setTotal(res.data.data.length)
        }
      }
    } catch (err: unknown) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchStats()
    fetchClasses()
  }, [])

  useEffect(() => {
    fetchStudents()
  }, [page, classFilter, statusFilter])

  const handleSearchSubmit = (e: FormEvent) => {
    e.preventDefault()
    setPage(1)
    fetchStudents()
  }

  // Checkbox Selection
  const toggleSelectAll = () => {
    if (selectedIds.length === students.length) {
      setSelectedIds([])
    } else {
      setSelectedIds(students.map((s) => s.id))
    }
  }

  const toggleSelectOne = (id: string) => {
    if (selectedIds.includes(id)) {
      setSelectedIds(selectedIds.filter((item) => item !== id))
    } else {
      setSelectedIds([...selectedIds, id])
    }
  }

  // Create Student Submit
  const handleCreateStudent = async (e: FormEvent) => {
    e.preventDefault()
    setFormError(null)
    setSubmitting(true)
    try {
      const res = await api.post('/students', createForm)
      if (res.data.success) {
        setCreateModal(false)
        fetchStudents()
        fetchStats()
        // Reset form
        setCreateForm({
          name: '',
          nis: '',
          nisn: '',
          gender: 'L',
          birthPlace: '',
          birthDate: '',
          address: '',
          parentName: '',
          parentPhone: '',
          enrollmentYear: new Date().getFullYear(),
          classId: '',
          status: 'active',
          createAccount: true,
          username: '',
          email: '',
          password: 'Password123!',
        })
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal menyimpan data siswa')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Edit Student Submit
  const handleEditStudent = async (e: FormEvent) => {
    e.preventDefault()
    if (!editModal) return
    setFormError(null)
    setSubmitting(true)
    try {
      await api.put(`/students/${editModal.id}`, editForm)
      setEditModal(null)
      fetchStudents()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal memperbarui data siswa')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Bulk Status Change
  const handleBulkStatusChange = async (status: string) => {
    if (selectedIds.length === 0) return
    if (!confirm(`Ubah status ${selectedIds.length} siswa menjadi "${status}"?`)) return
    try {
      await api.post('/students/bulk/status', {
        studentIds: selectedIds,
        status,
      })
      setSelectedIds([])
      fetchStudents()
      fetchStats()
    } catch {
      alert('Gagal mengubah status massal')
    }
  }

  // Bulk Assign Class
  const handleBulkAssignClass = async () => {
    if (selectedIds.length === 0 || !selectedBulkClass) return
    try {
      await api.post('/students/bulk/class', {
        studentIds: selectedIds,
        classId: selectedBulkClass,
      })
      setBulkClassModal(false)
      setSelectedIds([])
      fetchStudents()
      fetchStats()
    } catch {
      alert('Gagal mengalokasikan kelas')
    }
  }

  // Create Student Account
  const handleCreateAccountSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!accountModal) return
    setFormError(null)
    setSubmitting(true)
    try {
      await api.post(`/students/${accountModal.id}/account`, accountForm)
      setAccountModal(null)
      fetchStudents()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        setFormError(err.message)
      } else {
        setFormError('Gagal membuat akun murid')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Toggle Account Active / Inactive
  const handleToggleAccount = async (studentId: string) => {
    try {
      await api.patch(`/students/${studentId}/account/toggle`)
      fetchStudents()
      fetchStats()
    } catch {
      alert('Gagal mengubah status akun murid')
    }
  }

  // Delete Student
  const handleDeleteStudent = async () => {
    if (!deleteModal) return
    setSubmitting(true)
    try {
      await api.delete(`/students/${deleteModal.id}`)
      setDeleteModal(null)
      setSelectedIds((prev) => prev.filter((id) => id !== deleteModal.id))
      fetchStudents()
      fetchStats()
    } catch (err: unknown) {
      if (err instanceof Error) {
        alert(`Gagal menghapus siswa: ${err.message}`)
      } else {
        alert('Gagal menghapus data siswa')
      }
    } finally {
      setSubmitting(false)
    }
  }

  // Open Edit Modal
  const openEdit = (s: Student) => {
    setEditModal(s)
    setEditForm({
      name: s.name,
      nis: s.nis,
      nisn: s.nisn || '',
      gender: s.gender || 'L',
      birthPlace: s.birthPlace || '',
      birthDate: s.birthDate || '',
      address: s.address || '',
      parentName: s.parentName || '',
      parentPhone: s.parentPhone || '',
      enrollmentYear: s.enrollmentYear || new Date().getFullYear(),
      classId: s.classId || 'none',
      status: s.status,
    })
  }

  // Print function
  const handlePrint = () => {
    window.print()
  }

  // Export CSV
  const handleExportCSV = () => {
    if (students.length === 0) return
    const headers = ['Nama', 'NIS', 'NISN', 'Kelas', 'Status', 'Akun', 'Wali', 'Kontak']
    const rows = students.map((s) => [
      `"${s.name}"`,
      `"${s.nis}"`,
      `"${s.nisn || '-'}"`,
      `"${s.className}"`,
      `"${s.status}"`,
      `"${s.hasAccount ? s.username : 'Belum Ada'}"`,
      `"${s.parentName || '-'}"`,
      `"${s.parentPhone || '-'}"`,
    ])
    const csvContent = 'data:text/csv;charset=utf-8,' + [headers.join(','), ...rows.map((e) => e.join(','))].join('\n')
    const encodedUri = encodeURI(csvContent)
    const link = document.createElement('a')
    link.setAttribute('href', encodedUri)
    link.setAttribute('download', `Daftar_Siswa_${new Date().toISOString().slice(0, 10)}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  return (
    <AdminLayout>
      <div className="mx-auto max-w-7xl space-y-6">
        {/* ========================================================================= */}
        {/* 1. HEADER & RINGKASAN STATISTIK */}
        {/* ========================================================================= */}
        <div className="flex flex-col gap-4 rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-6 shadow-sm sm:flex-row sm:items-center sm:justify-between">
          <div>
            <div className="flex items-center gap-2">
              <span className="flex h-2.5 w-2.5 rounded-full bg-[#FF9800]" />
              <span className="text-[11px] font-bold uppercase tracking-wider text-[#FF9800]">
                Administrasi Akademik
              </span>
            </div>
            <h1 className="mt-1 text-2xl font-black text-[#30251D] sm:text-3xl tracking-tight">
              Manajemen Siswa
            </h1>
            <p className="mt-1 text-xs text-gray-500">
              Kelola data induk siswa, penempatan rombongan belajar, akun login, dan arsip data.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-2.5">
            <button
              type="button"
              onClick={handlePrint}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] px-3.5 py-2 text-xs font-bold text-[#30251D] hover:bg-[#EBD4C5]/40 transition-colors"
            >
              <Printer size={15} />
              <span>Cetak</span>
            </button>
            <button
              type="button"
              onClick={handleExportCSV}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] px-3.5 py-2 text-xs font-bold text-[#30251D] hover:bg-[#EBD4C5]/40 transition-colors"
            >
              <Download size={15} />
              <span>Export CSV</span>
            </button>
            <button
              type="button"
              onClick={() => setCreateModal(true)}
              className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white shadow-sm hover:bg-[#E68600] transition-colors"
            >
              <Plus size={16} />
              <span>Tambah Siswa Baru</span>
            </button>
          </div>
        </div>

        {/* 4 STAT CARDS */}
        <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-[11px] font-bold uppercase text-gray-400">Total Siswa</span>
              <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#FEEADC] text-[#FF9800]">
                <GraduationCap size={16} />
              </span>
            </div>
            <div className="mt-2 text-2xl font-black text-[#30251D]">{stats.totalStudents}</div>
            <p className="text-[10px] text-gray-500">Terdaftar di database</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-[11px] font-bold uppercase text-gray-400">Siswa Aktif</span>
              <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-100 text-emerald-700">
                <UserCheck size={16} />
              </span>
            </div>
            <div className="mt-2 text-2xl font-black text-[#30251D]">{stats.activeStudents}</div>
            <p className="text-[10px] text-gray-500">Status pembelajaran aktif</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-[11px] font-bold uppercase text-gray-400">Belum Punya Kelas</span>
              <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-red-100 text-red-700">
                <Layers size={16} />
              </span>
            </div>
            <div className="mt-2 text-2xl font-black text-red-600">{stats.studentsWithoutClass}</div>
            <p className="text-[10px] text-gray-500">Perlu alokasi rombel</p>
          </div>

          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-4 shadow-sm">
            <div className="flex items-center justify-between">
              <span className="text-[11px] font-bold uppercase text-gray-400">Belum Punya Akun</span>
              <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-amber-100 text-amber-800">
                <UserX size={16} />
              </span>
            </div>
            <div className="mt-2 text-2xl font-black text-amber-600">{stats.studentsWithoutAccount}</div>
            <p className="text-[10px] text-gray-500">Belum dapat login sistem</p>
          </div>
        </div>

        {/* ========================================================================= */}
        {/* 2. PENCARIAN & FILTER TOOLBAR */}
        {/* ========================================================================= */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-sm space-y-4">
          <form onSubmit={handleSearchSubmit} className="flex flex-col sm:flex-row gap-3">
            <div className="relative flex-1">
              <Search size={16} className="pointer-events-none absolute left-3.5 top-3 text-gray-400" />
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Cari siswa berdasarkan nama, NIS, atau NISN..."
                className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 py-2.5 pl-10 pr-4 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:bg-white focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
              />
            </div>
            <button
              type="submit"
              className="rounded-xl bg-[#FF9800] px-5 py-2.5 text-xs font-bold text-white shadow-sm hover:bg-[#E68600] transition-colors"
            >
              Cari Siswa
            </button>
          </form>

          {/* Filter Bar */}
          <div className="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-[#EBD4C5]/60">
            <div className="flex flex-wrap items-center gap-2.5">
              <div className="flex items-center gap-1.5 text-xs font-semibold text-gray-500">
                <Filter size={14} className="text-[#FF9800]" />
                <span>Filter:</span>
              </div>

              {/* Filter Kelas */}
              <select
                value={classFilter}
                onChange={(e) => {
                  setClassFilter(e.target.value)
                  setPage(1)
                }}
                className="rounded-lg border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs text-[#30251D] focus:border-[#FF9800] focus:outline-none"
              >
                <option value="">Semua Rombel / Kelas</option>
                <option value="none">Khusus Belum Ada Kelas</option>
                {classes.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>

              {/* Filter Status */}
              <select
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value)
                  setPage(1)
                }}
                className="rounded-lg border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs text-[#30251D] focus:border-[#FF9800] focus:outline-none"
              >
                <option value="">Semua Status</option>
                <option value="active">Aktif</option>
                <option value="inactive">Nonaktif</option>
                <option value="graduated">Lulus</option>
                <option value="transferred">Pindah</option>
                <option value="archived">Diarsipkan</option>
              </select>
            </div>

            {/* Aksi Massal (Muncul jika ada siswa yang dipilih) */}
            {selectedIds.length > 0 && (
              <div className="flex items-center gap-2 bg-[#FEEADC] px-3 py-1.5 rounded-xl border border-[#FF9800]/40">
                <span className="text-xs font-bold text-[#FF9800]">
                  {selectedIds.length} Terpilih:
                </span>
                <button
                  type="button"
                  onClick={() => setBulkClassModal(true)}
                  className="rounded-lg bg-white px-2.5 py-1 text-[11px] font-bold text-[#30251D] shadow-sm hover:bg-gray-50"
                >
                  Pindah Kelas
                </button>
                <button
                  type="button"
                  onClick={() => handleBulkStatusChange('graduated')}
                  className="rounded-lg bg-white px-2.5 py-1 text-[11px] font-bold text-emerald-700 shadow-sm hover:bg-emerald-50"
                >
                  Proses Kelulusan
                </button>
                <button
                  type="button"
                  onClick={() => handleBulkStatusChange('archived')}
                  className="rounded-lg bg-white px-2.5 py-1 text-[11px] font-bold text-gray-700 shadow-sm hover:bg-gray-50"
                >
                  Arsipkan
                </button>
              </div>
            )}
          </div>
        </div>

        {/* ========================================================================= */}
        {/* 3. TABEL DAFTAR SISWA */}
        {/* ========================================================================= */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] shadow-sm overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-500 uppercase tracking-wider text-[10px]">
                <tr>
                  <th className="py-3.5 px-4 w-10">
                    <button type="button" onClick={toggleSelectAll} className="text-[#FF9800]">
                      {selectedIds.length === students.length && students.length > 0 ? (
                        <CheckSquare size={16} />
                      ) : (
                        <Square size={16} />
                      )}
                    </button>
                  </th>
                  <th className="py-3.5 px-4 font-bold">Identitas Siswa</th>
                  <th className="py-3.5 px-4 font-bold">NIS / NISN</th>
                  <th className="py-3.5 px-4 font-bold">Kelas</th>
                  <th className="py-3.5 px-4 font-bold">Status Siswa</th>
                  <th className="py-3.5 px-4 font-bold">Akun Murid</th>
                  <th className="py-3.5 px-4 font-bold text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#EBD4C5]/60">
                {loading ? (
                  <tr>
                    <td colSpan={7} className="py-10 text-center text-xs text-gray-400">
                      <div className="flex justify-center items-center gap-2">
                        <RefreshCw size={16} className="animate-spin text-[#FF9800]" />
                        <span>Memuat data siswa...</span>
                      </div>
                    </td>
                  </tr>
                ) : students.length > 0 ? (
                  students.map((s) => {
                    const isSelected = selectedIds.includes(s.id)
                    return (
                      <tr
                        key={s.id}
                        className={`transition-colors ${
                          isSelected ? 'bg-[#FEEADC]/30' : 'hover:bg-[#FFF8F5]/50'
                        }`}
                      >
                        <td className="py-3.5 px-4">
                          <button
                            type="button"
                            onClick={() => toggleSelectOne(s.id)}
                            className="text-[#FF9800]"
                          >
                            {isSelected ? <CheckSquare size={16} /> : <Square size={16} />}
                          </button>
                        </td>

                        {/* Foto & Identitas */}
                        <td className="py-3.5 px-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-[#FEEADC] font-bold text-[#FF9800]">
                              {s.name.charAt(0)}
                            </div>
                            <div>
                              <p className="font-bold text-[#30251D] hover:text-[#FF9800] cursor-pointer" onClick={() => setDetailModal(s)}>
                                {s.name}
                              </p>
                              <span className="text-[10px] text-gray-400">
                                {s.gender === 'L' ? 'Laki-laki' : s.gender === 'P' ? 'Perempuan' : '-'} • Masuk: {s.enrollmentYear || '-'}
                              </span>
                            </div>
                          </div>
                        </td>

                        {/* NIS / NISN */}
                        <td className="py-3.5 px-4 font-mono">
                          <span className="block font-bold text-[#30251D]">{s.nis}</span>
                          <span className="text-[10px] text-gray-400">NISN: {s.nisn || '-'}</span>
                        </td>

                        {/* Kelas */}
                        <td className="py-3.5 px-4">
                          <span className={`inline-block rounded-lg px-2.5 py-1 text-[11px] font-bold ${
                            s.className === 'Belum Ada Kelas'
                              ? 'bg-red-50 text-red-600 border border-red-200'
                              : 'bg-[#FFF8F5] text-[#30251D] border border-[#EBD4C5]'
                          }`}>
                            {s.className}
                          </span>
                        </td>

                        {/* Status Siswa */}
                        <td className="py-3.5 px-4">
                          <span className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-[10px] font-bold capitalize ${
                            s.status === 'active'
                              ? 'bg-emerald-100 text-emerald-700'
                              : s.status === 'graduated'
                              ? 'bg-blue-100 text-blue-700'
                              : s.status === 'transferred'
                              ? 'bg-amber-100 text-amber-800'
                              : 'bg-gray-100 text-gray-600'
                          }`}>
                            <span className="h-1.5 w-1.5 rounded-full bg-current" />
                            {s.status}
                          </span>
                        </td>

                        {/* Akun Murid */}
                        <td className="py-3.5 px-4">
                          {s.hasAccount ? (
                            <div className="flex items-center gap-2">
                              <div>
                                <span className="block font-mono text-[11px] font-bold text-[#30251D]">
                                  @{s.username}
                                </span>
                                <span className={`text-[10px] ${s.isActive ? 'text-emerald-600 font-semibold' : 'text-red-500'}`}>
                                  {s.isActive ? 'Akun Aktif' : 'Akun Nonaktif'}
                                </span>
                              </div>
                              <button
                                type="button"
                                onClick={() => handleToggleAccount(s.id)}
                                title={s.isActive ? 'Nonaktifkan Akun' : 'Aktifkan Akun'}
                                className="rounded p-1 text-gray-400 hover:text-[#FF9800]"
                              >
                                <Lock size={13} />
                              </button>
                            </div>
                          ) : (
                            <button
                              type="button"
                              onClick={() => {
                                setAccountModal(s)
                                setAccountForm({
                                  username: s.nis,
                                  email: `${s.nis}@siswa.genius.ac.id`,
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
                              onClick={() => setDetailModal(s)}
                              title="Lihat Profil Siswa"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
                            >
                              <Eye size={15} />
                            </button>
                            <button
                              type="button"
                              onClick={() => openEdit(s)}
                              title="Edit Data Siswa"
                              className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
                            >
                              <Edit size={15} />
                            </button>
                            <button
                              type="button"
                              onClick={() => setDeleteModal(s)}
                              title="Hapus Data Siswa"
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
                    <td colSpan={7} className="py-10 text-center text-xs text-gray-400">
                      Tidak ada data siswa yang cocok dengan filter atau pencarian.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination Footer */}
          <div className="flex items-center justify-between border-t border-[#EBD4C5]/60 bg-[#FFF8F5]/40 px-5 py-3 text-xs">
            <span className="text-gray-500 text-[11px]">
              Menampilkan {students.length} dari total {total} siswa
            </span>
            <div className="flex items-center gap-2">
              <button
                type="button"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
                className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] bg-white text-gray-600 hover:bg-[#FFF8F5] disabled:opacity-40"
              >
                <ChevronLeft size={14} />
              </button>
              <span className="font-bold text-[#30251D]">Halaman {page}</span>
              <button
                type="button"
                disabled={students.length < limit}
                onClick={() => setPage(page + 1)}
                className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] bg-white text-gray-600 hover:bg-[#FFF8F5] disabled:opacity-40"
              >
                <ChevronRight size={14} />
              </button>
            </div>
          </div>
        </div>

        {/* ========================================================================= */}
        {/* MODAL: TAMBAH SISWA BARU */}
        {/* ========================================================================= */}
        {createModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm">
            <div className="w-full max-w-2xl max-h-[90vh] overflow-y-auto rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl">
              <div className="flex items-center justify-between pb-4 border-b border-[#EBD4C5]/70">
                <div>
                  <h3 className="text-base font-black text-[#30251D]">Tambah Siswa Baru</h3>
                  <p className="text-xs text-gray-400">Masukkan data identitas pokok siswa</p>
                </div>
                <button
                  type="button"
                  onClick={() => setCreateModal(false)}
                  className="rounded-lg p-1.5 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              {formError && (
                <div className="mt-4 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                  <AlertCircle size={15} />
                  <span>{formError}</span>
                </div>
              )}

              <form onSubmit={handleCreateStudent} className="mt-4 space-y-4 text-xs">
                <div className="grid sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama Lengkap Siswa *</label>
                    <input
                      type="text"
                      required
                      value={createForm.name}
                      onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                      placeholder="e.g. Raditya Pratama"
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">NIS (Nomor Induk Siswa) *</label>
                    <input
                      type="text"
                      required
                      value={createForm.nis}
                      onChange={(e) => setCreateForm({ ...createForm, nis: e.target.value })}
                      placeholder="e.g. 2401001"
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid sm:grid-cols-3 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">NISN</label>
                    <input
                      type="text"
                      value={createForm.nisn}
                      onChange={(e) => setCreateForm({ ...createForm, nisn: e.target.value })}
                      placeholder="e.g. 0081234567"
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Jenis Kelamin</label>
                    <select
                      value={createForm.gender}
                      onChange={(e) => setCreateForm({ ...createForm, gender: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    >
                      <option value="L">Laki-laki</option>
                      <option value="P">Perempuan</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Rombel / Kelas</label>
                    <select
                      value={createForm.classId}
                      onChange={(e) => setCreateForm({ ...createForm, classId: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    >
                      <option value="">-- Belum Ada Kelas --</option>
                      {classes.map((c) => (
                        <option key={c.id} value={c.id}>
                          {c.name}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="grid sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tempat Lahir</label>
                    <input
                      type="text"
                      value={createForm.birthPlace}
                      onChange={(e) => setCreateForm({ ...createForm, birthPlace: e.target.value })}
                      placeholder="e.g. Jakarta"
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tanggal Lahir</label>
                    <input
                      type="date"
                      value={createForm.birthDate}
                      onChange={(e) => setCreateForm({ ...createForm, birthDate: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                </div>

                <div>
                  <label className="block mb-1 font-bold text-gray-700">Alamat Tempat Tinggal</label>
                  <textarea
                    rows={2}
                    value={createForm.address}
                    onChange={(e) => setCreateForm({ ...createForm, address: e.target.value })}
                    placeholder="Alamat domisili siswa..."
                    className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                  />
                </div>

                <div className="grid sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama Orang Tua / Wali</label>
                    <input
                      type="text"
                      value={createForm.parentName}
                      onChange={(e) => setCreateForm({ ...createForm, parentName: e.target.value })}
                      placeholder="e.g. Hendra Setiawan"
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">No. HP Orang Tua</label>
                    <input
                      type="text"
                      value={createForm.parentPhone}
                      onChange={(e) => setCreateForm({ ...createForm, parentPhone: e.target.value })}
                      placeholder="e.g. 08123456789"
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                </div>

                {/* Section Buat Akun Murid */}
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-4">
                  <label className="flex items-center gap-2 cursor-pointer font-bold text-[#30251D]">
                    <input
                      type="checkbox"
                      checked={createForm.createAccount}
                      onChange={(e) => setCreateForm({ ...createForm, createAccount: e.target.checked })}
                      className="rounded text-[#FF9800] focus:ring-[#FF9800]"
                    />
                    <span>Sekaligus Buat Akun Login Murid</span>
                  </label>

                  {createForm.createAccount && (
                    <div className="grid sm:grid-cols-3 gap-3 mt-3 pt-3 border-t border-[#EBD4C5]/60">
                      <div>
                        <label className="block mb-1 font-semibold text-gray-600">Username *</label>
                        <input
                          type="text"
                          required={createForm.createAccount}
                          value={createForm.username}
                          onChange={(e) => setCreateForm({ ...createForm, username: e.target.value })}
                          placeholder="e.g. raditya24"
                          className="w-full rounded-lg border border-[#EBD4C5] bg-white p-2 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                        />
                      </div>
                      <div>
                        <label className="block mb-1 font-semibold text-gray-600">Email Akun</label>
                        <input
                          type="email"
                          value={createForm.email}
                          onChange={(e) => setCreateForm({ ...createForm, email: e.target.value })}
                          placeholder="e.g. raditya@siswa.genius.ac.id"
                          className="w-full rounded-lg border border-[#EBD4C5] bg-white p-2 text-xs focus:border-[#FF9800] focus:outline-none"
                        />
                      </div>
                      <div>
                        <label className="block mb-1 font-semibold text-gray-600">Password Awal *</label>
                        <input
                          type="text"
                          required={createForm.createAccount}
                          value={createForm.password}
                          onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })}
                          className="w-full rounded-lg border border-[#EBD4C5] bg-white p-2 text-xs font-mono focus:border-[#FF9800] focus:outline-none"
                        />
                      </div>
                    </div>
                  )}
                </div>

                <div className="flex justify-end gap-2.5 pt-3 border-t border-[#EBD4C5]/70">
                  <button
                    type="button"
                    onClick={() => setCreateModal(false)}
                    className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 text-xs font-bold text-gray-600 hover:bg-gray-50"
                  >
                    Batal
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="rounded-xl bg-[#FF9800] px-5 py-2 text-xs font-bold text-white shadow-sm hover:bg-[#E68600] disabled:opacity-70"
                  >
                    {submitting ? 'Menyimpan...' : 'Simpan Siswa'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ========================================================================= */}
        {/* MODAL: DETAIL & PROFIL SISWA */}
        {/* ========================================================================= */}
        {detailModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm">
            <div className="w-full max-w-xl rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl">
              <div className="flex items-center justify-between pb-4 border-b border-[#EBD4C5]/70">
                <div className="flex items-center gap-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#FEEADC] text-xl font-black text-[#FF9800]">
                    {detailModal.name.charAt(0)}
                  </div>
                  <div>
                    <h3 className="text-lg font-black text-[#30251D]">{detailModal.name}</h3>
                    <p className="font-mono text-xs text-gray-400">NIS: {detailModal.nis} • NISN: {detailModal.nisn || '-'}</p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => setDetailModal(null)}
                  className="rounded-lg p-1.5 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              <div className="mt-4 space-y-4 text-xs">
                <div className="grid grid-cols-2 gap-3 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3.5">
                  <div>
                    <span className="text-gray-400">Rombel / Kelas:</span>
                    <p className="font-bold text-[#30251D]">{detailModal.className}</p>
                  </div>
                  <div>
                    <span className="text-gray-400">Status Siswa:</span>
                    <p className="font-bold capitalize text-[#FF9800]">{detailModal.status}</p>
                  </div>
                  <div>
                    <span className="text-gray-400">Jenis Kelamin:</span>
                    <p className="font-semibold text-[#30251D]">{detailModal.gender === 'L' ? 'Laki-laki' : 'Perempuan'}</p>
                  </div>
                  <div>
                    <span className="text-gray-400">Tahun Masuk:</span>
                    <p className="font-semibold text-[#30251D]">{detailModal.enrollmentYear || '-'}</p>
                  </div>
                </div>

                <div className="rounded-xl border border-[#EBD4C5] p-3.5 space-y-2">
                  <h4 className="font-bold uppercase tracking-wider text-[10px] text-gray-400">Data Domisili & Keluarga</h4>
                  <div className="grid grid-cols-2 gap-2">
                    <div>
                      <span className="text-gray-400">Tempat, Tanggal Lahir:</span>
                      <p className="font-semibold text-[#30251D]">{detailModal.birthPlace || '-'}, {detailModal.birthDate || '-'}</p>
                    </div>
                    <div>
                      <span className="text-gray-400">Nama Orang Tua / Wali:</span>
                      <p className="font-semibold text-[#30251D]">{detailModal.parentName || '-'}</p>
                    </div>
                    <div>
                      <span className="text-gray-400">Kontak Orang Tua:</span>
                      <p className="font-mono font-semibold text-[#30251D]">{detailModal.parentPhone || '-'}</p>
                    </div>
                    <div>
                      <span className="text-gray-400">Alamat:</span>
                      <p className="font-semibold text-[#30251D]">{detailModal.address || '-'}</p>
                    </div>
                  </div>
                </div>

                <div className="rounded-xl border border-[#EBD4C5] p-3.5 space-y-2">
                  <h4 className="font-bold uppercase tracking-wider text-[10px] text-gray-400">Status Akun Login Murid</h4>
                  {detailModal.hasAccount ? (
                    <div className="flex items-center justify-between">
                      <div>
                        <span className="font-mono font-bold text-xs text-[#30251D]">@{detailModal.username}</span>
                        <p className="text-gray-400">{detailModal.email}</p>
                      </div>
                      <span className={`rounded-full px-2.5 py-0.5 text-[10px] font-bold ${
                        detailModal.isActive ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
                      }`}>
                        {detailModal.isActive ? 'Akun Aktif' : 'Akun Nonaktif'}
                      </span>
                    </div>
                  ) : (
                    <p className="text-amber-700">Siswa ini belum memiliki akun login di sistem Genius Society.</p>
                  )}
                </div>
              </div>

              <div className="mt-5 flex justify-end gap-2 border-t border-[#EBD4C5]/70 pt-3">
                <button
                  type="button"
                  onClick={() => setDetailModal(null)}
                  className="rounded-xl bg-gray-100 px-4 py-2 text-xs font-bold text-gray-700 hover:bg-gray-200"
                >
                  Tutup
                </button>
              </div>
            </div>
          </div>
        )}

        {/* ========================================================================= */}
        {/* MODAL: EDIT DATA SISWA */}
        {/* ========================================================================= */}
        {editModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm">
            <div className="w-full max-w-xl max-h-[90vh] overflow-y-auto rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl">
              <div className="flex items-center justify-between pb-4 border-b border-[#EBD4C5]/70">
                <div>
                  <h3 className="text-base font-black text-[#30251D]">Edit Data Siswa</h3>
                  <p className="text-xs text-gray-400">{editModal.name} ({editModal.nis})</p>
                </div>
                <button
                  type="button"
                  onClick={() => setEditModal(null)}
                  className="rounded-lg p-1.5 text-gray-400 hover:text-gray-700"
                >
                  <X size={18} />
                </button>
              </div>

              {formError && (
                <div className="mt-4 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                  <AlertCircle size={15} />
                  <span>{formError}</span>
                </div>
              )}

              <form onSubmit={handleEditStudent} className="mt-4 space-y-4 text-xs">
                <div className="grid sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama Siswa</label>
                    <input
                      type="text"
                      required
                      value={editForm.name}
                      onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">NIS</label>
                    <input
                      type="text"
                      required
                      value={editForm.nis}
                      onChange={(e) => setEditForm({ ...editForm, nis: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid sm:grid-cols-3 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Rombel / Kelas</label>
                    <select
                      value={editForm.classId}
                      onChange={(e) => setEditForm({ ...editForm, classId: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    >
                      <option value="none">-- Belum Ada Kelas --</option>
                      {classes.map((c) => (
                        <option key={c.id} value={c.id}>
                          {c.name}
                        </option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Status</label>
                    <select
                      value={editForm.status}
                      onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    >
                      <option value="active">Aktif</option>
                      <option value="inactive">Nonaktif</option>
                      <option value="graduated">Lulus</option>
                      <option value="transferred">Pindah</option>
                      <option value="archived">Diarsipkan</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tahun Masuk</label>
                    <input
                      type="number"
                      value={editForm.enrollmentYear}
                      onChange={(e) => setEditForm({ ...editForm, enrollmentYear: parseInt(e.target.value) || 0 })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                </div>

                <div className="grid sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Orang Tua / Wali</label>
                    <input
                      type="text"
                      value={editForm.parentName}
                      onChange={(e) => setEditForm({ ...editForm, parentName: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">No. HP Orang Tua</label>
                    <input
                      type="text"
                      value={editForm.parentPhone}
                      onChange={(e) => setEditForm({ ...editForm, parentPhone: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 p-2.5 text-xs font-mono focus:border-[#FF9800] focus:bg-white focus:outline-none"
                    />
                  </div>
                </div>

                <div className="flex justify-end gap-2.5 pt-3 border-t border-[#EBD4C5]/70">
                  <button
                    type="button"
                    onClick={() => setEditModal(null)}
                    className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 text-xs font-bold text-gray-600 hover:bg-gray-50"
                  >
                    Batal
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="rounded-xl bg-[#FF9800] px-5 py-2 text-xs font-bold text-white shadow-sm hover:bg-[#E68600] disabled:opacity-70"
                  >
                    {submitting ? 'Memperbarui...' : 'Simpan Perubahan'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ========================================================================= */}
        {/* MODAL: PINDAH KELAS MASSAL */}
        {/* ========================================================================= */}
        {bulkClassModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-xs">
              <h3 className="text-base font-black text-[#30251D]">Pindahkan Kelas Massal</h3>
              <p className="mt-1 text-gray-500">
                Pilih rombel tujuan untuk <strong>{selectedIds.length}</strong> siswa yang telah Anda pilih.
              </p>

              <div className="mt-4">
                <label className="block mb-1 font-bold text-gray-700">Rombel Tujuan</label>
                <select
                  value={selectedBulkClass}
                  onChange={(e) => setSelectedBulkClass(e.target.value)}
                  className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-2.5 text-xs focus:border-[#FF9800] focus:outline-none"
                >
                  <option value="">-- Pilih Kelas Tujuan --</option>
                  <option value="none">Kosongkan Kelas (Tanpa Kelas)</option>
                  {classes.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="mt-6 flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => setBulkClassModal(false)}
                  className="rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 font-bold text-gray-600"
                >
                  Batal
                </button>
                <button
                  type="button"
                  disabled={!selectedBulkClass}
                  onClick={handleBulkAssignClass}
                  className="rounded-xl bg-[#FF9800] px-5 py-2 font-bold text-white shadow-sm hover:bg-[#E68600] disabled:opacity-50"
                >
                  Terapkan Kelas
                </button>
              </div>
            </div>
          </div>
        )}

        {/* ========================================================================= */}
        {/* MODAL: BUAT AKUN MURID */}
        {/* ========================================================================= */}
        {accountModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-xs">
              <div className="flex items-center justify-between pb-3 border-b border-[#EBD4C5]/70">
                <div>
                  <h3 className="text-base font-black text-[#30251D]">Buat Akun Murid</h3>
                  <p className="text-gray-400">Untuk siswa: {accountModal.name}</p>
                </div>
                <button type="button" onClick={() => setAccountModal(null)} className="text-gray-400 hover:text-gray-700">
                  <X size={18} />
                </button>
              </div>

              {formError && (
                <div className="mt-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-2.5 text-xs text-red-600">
                  <AlertCircle size={15} />
                  <span>{formError}</span>
                </div>
              )}

              <form onSubmit={handleCreateAccountSubmit} className="mt-4 space-y-3">
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

        {/* ================= MODAL KONFIRMASI HAPUS SISWA ================= */}
        {deleteModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-red-100 text-red-600">
                  <Trash2 size={22} />
                </div>
                <div>
                  <h3 className="text-base font-black text-[#30251D]">Hapus Data Siswa</h3>
                  <p className="text-xs text-gray-500">Tindakan ini permanen dan tidak dapat dibatalkan</p>
                </div>
              </div>

              <div className="mt-4 rounded-xl border border-red-100 bg-red-50/60 p-3.5 text-xs text-red-800 space-y-1.5">
                <p>
                  Apakah Anda yakin ingin menghapus data siswa <strong className="font-bold text-red-900">{deleteModal.name}</strong> (NIS: {deleteModal.nis})?
                </p>
                {deleteModal.hasAccount && (
                  <p className="text-[11px] text-red-700 font-semibold">
                    ⚠️ Akun login murid terkait (@{deleteModal.username}) juga akan dihapus dari sistem.
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
                  onClick={handleDeleteStudent}
                  className="inline-flex items-center gap-1.5 rounded-xl bg-red-600 px-5 py-2 text-xs font-bold text-white shadow-sm hover:bg-red-700 disabled:opacity-70 transition-colors"
                >
                  <Trash2 size={14} />
                  <span>{submitting ? 'Menghapus...' : 'Ya, Hapus Siswa'}</span>
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
