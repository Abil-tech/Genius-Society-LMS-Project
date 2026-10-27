import { useState, useEffect, type FormEvent } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import { Search, Plus, Edit, Trash2, X, RefreshCw, ArrowUpDown, AlertCircle } from 'lucide-react'

interface ClassItem {
  id: string
  code: string
  name: string
  educationLevel: string
  department?: string
  academicYearId: string
  homeroomTeacherId?: string
  capacity: number
  status: string
  activeStudentCount?: number
}

interface AcademicYearOption {
  id: string
  name: string
  semester?: number
  isActive?: boolean
}

interface TeacherOption {
  id: string
  name: string
}

export default function ClassesManagementPage() {
  const [classes, setClasses] = useState<ClassItem[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const limit = 15

  // Filters & search
  const [search, setSearch] = useState('')
  const [educationLevelFilter, setEducationLevelFilter] = useState('')
  const [departmentFilter, setDepartmentFilter] = useState('')
  const [academicYearFilter, setAcademicYearFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')

  // Options for selects
  const [academicYears, setAcademicYears] = useState<AcademicYearOption[]>([])
  const [academicYearsLoading, setAcademicYearsLoading] = useState(false)
  const [academicYearsError, setAcademicYearsError] = useState<string | null>(null)
  const [teachers, setTeachers] = useState<TeacherOption[]>([])

  // Modal state
  const [createModal, setCreateModal] = useState(false)
  const [editModal, setEditModal] = useState<ClassItem | null>(null)
  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // Form state for create / edit
  const emptyForm = {
    code: '',
    name: '',
    educationLevel: '',
    department: '',
    academicYearId: '',
    homeroomTeacherId: '',
    capacity: 0,
    status: 'active',
  }
  const [form, setForm] = useState({ ...emptyForm })

  // ---------------------------------------------------------------------------
  // Data fetching helpers
  // ---------------------------------------------------------------------------
  const fetchClasses = async () => {
    setLoading(true)
    try {
      const params: Record<string, string | number> = { page, limit }
      if (search) params.search = search
      if (educationLevelFilter) params.educationLevel = educationLevelFilter
      if (departmentFilter) params.department = departmentFilter
      if (academicYearFilter) params.academicYearId = academicYearFilter
      if (statusFilter) params.status = statusFilter
      const res = await api.get<ApiResponse<ClassItem[]>>('/classes', { params })
      if (res.data.success && res.data.data) {
        setClasses(res.data.data)
        if (res.data.meta) setTotal((res.data.meta as any).total || 0)
        else setTotal(res.data.data.length)
      }
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const fetchAcademicYears = async () => {
    setAcademicYearsLoading(true)
    setAcademicYearsError(null)
    try {
      const res = await api.get<ApiResponse<AcademicYearOption[]>>('/academic-years')
      if (res.data.success && Array.isArray(res.data.data)) {
        setAcademicYears(res.data.data)
      } else {
        setAcademicYears([])
      }
    } catch (err: unknown) {
      console.error('Error fetching academic years:', err)
      setAcademicYearsError('Gagal memuat tahun ajaran')
      setAcademicYears([])
    } finally {
      setAcademicYearsLoading(false)
    }
  }

  const fetchTeachers = async () => {
    try {
      const res = await api.get<ApiResponse<TeacherOption[]>>('/teachers')
      if (res.data.success && res.data.data) setTeachers(res.data.data)
    } catch {
      setTeachers([])
    }
  }

  useEffect(() => {
    fetchAcademicYears()
    fetchTeachers()
  }, [])

  useEffect(() => {
    fetchClasses()
  }, [page, educationLevelFilter, departmentFilter, academicYearFilter, statusFilter])

  const handleSearchSubmit = (e: FormEvent) => {
    e.preventDefault()
    setPage(1)
    fetchClasses()
  }

  // ---------------------------------------------------------------------------
  // Create / Edit handlers
  // ---------------------------------------------------------------------------
  const handleCreateSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setFormError(null)
    setSubmitting(true)
    try {
      const payload: any = {
        code: form.code,
        name: form.name,
        educationLevel: form.educationLevel,
        capacity: form.capacity,
        status: form.status,
      }
      if (form.department) payload.department = form.department
      if (form.academicYearId) payload.academicYearId = form.academicYearId
      if (form.homeroomTeacherId) payload.homeroomTeacherId = form.homeroomTeacherId
      const res = await api.post('/classes', payload)
      if (res.data.success) {
        setCreateModal(false)
        setForm({ ...emptyForm })
        fetchClasses()
      }
    } catch (err: unknown) {
      if (err instanceof Error) setFormError(err.message)
      else setFormError('Gagal membuat kelas')
    } finally {
      setSubmitting(false)
    }
  }

  const openEdit = (cls: ClassItem) => {
    setEditModal(cls)
    setForm({
      code: cls.code,
      name: cls.name,
      educationLevel: cls.educationLevel,
      department: cls.department || '',
      academicYearId: cls.academicYearId,
      homeroomTeacherId: cls.homeroomTeacherId || '',
      capacity: cls.capacity,
      status: cls.status,
    })
  }

  const handleEditSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!editModal) return
    setFormError(null)
    setSubmitting(true)
    try {
      const payload: any = {}
      if (form.name) payload.name = form.name
      if (form.educationLevel) payload.educationLevel = form.educationLevel
      if (form.department) payload.department = form.department
      if (form.academicYearId) payload.academicYearId = form.academicYearId
      if (form.homeroomTeacherId) payload.homeroomTeacherId = form.homeroomTeacherId
      if (form.capacity) payload.capacity = form.capacity
      if (form.status) payload.status = form.status

      const res = await api.put(`/classes/${editModal.id}`, payload)
      if (res.data.success) {
        setEditModal(null)
        fetchClasses()
      }
    } catch (err: unknown) {
      if (err instanceof Error) setFormError(err.message)
      else setFormError('Gagal memperbarui kelas')
    } finally {
      setSubmitting(false)
    }
  }

  // ---------------------------------------------------------------------------
  // Archive handler
  // ---------------------------------------------------------------------------
  const handleArchive = async (cls: ClassItem) => {
    if (!confirm(`Apakah Anda yakin ingin mengarsipkan kelas ${cls.name}?`)) return
    try {
      await api.delete(`/classes/${cls.id}`)
      fetchClasses()
    } catch {
      alert('Gagal mengarsipkan kelas')
    }
  }

  // ---------------------------------------------------------------------------
  // Render helpers
  // ---------------------------------------------------------------------------
  const getTeacherName = (id?: string) => {
    const t = teachers.find((t) => t.id === id)
    return t ? t.name : '-'
  }

  const getAcademicYearName = (id: string) => {
    const ay = academicYears.find((a) => a.id === id)
    return ay ? ay.name : id.slice(0, 6)
  }

  return (
    <AdminLayout>
      <div className="mx-auto max-w-7xl space-y-6">
        {/* HEADER */}
        <div className="flex flex-col gap-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-sm sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-black text-[#30251D]">Manajemen Kelas</h1>
            <p className="mt-1 text-xs text-gray-500">Kelola data kelas, wali kelas, dan kapasitas.</p>
          </div>
          <button type="button" onClick={() => setCreateModal(true)} className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white hover:bg-[#E68600]">
            <Plus size={16} />
            <span>Tambah Kelas</span>
          </button>
        </div>

        {/* FILTER & SEARCH */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm flex flex-col md:flex-row gap-3">
          <form onSubmit={handleSearchSubmit} className="relative flex-1">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Cari kode atau nama kelas..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-10 pr-4 py-2 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:bg-white focus:outline-none"
            />
          </form>
          <select value={educationLevelFilter} onChange={(e) => { setEducationLevelFilter(e.target.value); setPage(1) }} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
            <option value="">Semua Tingkat</option>
            <option value="SD">SD</option>
            <option value="SMP">SMP</option>
            <option value="SMA">SMA</option>
            <option value="SMK">SMK</option>
          </select>
          <input type="text" placeholder="Departemen" value={departmentFilter} onChange={(e) => { setDepartmentFilter(e.target.value); setPage(1) }} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs" />
          <select
            value={academicYearFilter}
            onChange={(e) => { setAcademicYearFilter(e.target.value); setPage(1) }}
            className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs"
            disabled={academicYearsLoading}
          >
            <option value="">{academicYearsLoading ? 'Memuat tahun ajaran...' : 'Semua Tahun Ajaran'}</option>
            {academicYears.map((ay) => (
              <option key={ay.id} value={ay.id}>
                {ay.name}{ay.semester ? ` (Sem ${ay.semester})` : ''}
              </option>
            ))}
          </select>
          <select value={statusFilter} onChange={(e) => { setStatusFilter(e.target.value); setPage(1) }} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
            <option value="">Semua Status</option>
            <option value="active">Aktif</option>
            <option value="inactive">Tidak Aktif</option>
            <option value="archived">Diarsipkan</option>
          </select>
          {(search || educationLevelFilter || departmentFilter || academicYearFilter || statusFilter) && (
            <button type="button" onClick={() => { setSearch(''); setEducationLevelFilter(''); setDepartmentFilter(''); setAcademicYearFilter(''); setStatusFilter(''); setPage(1) }} className="inline-flex items-center gap-1.5 rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-xs font-bold text-gray-600 hover:bg-gray-100">
              <RefreshCw size={13} />
              <span>Reset</span>
            </button>
          )}
        </div>

        {/* TABLE */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-sm overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 uppercase">
              <tr>
                <th className="py-3 px-4 font-bold">Kode</th>
                <th className="py-3 px-4 font-bold">Nama</th>
                <th className="py-3 px-4 font-bold">Tingkat</th>
                <th className="py-3 px-4 font-bold">Departemen</th>
                <th className="py-3 px-4 font-bold">Tahun Ajaran</th>
                <th className="py-3 px-4 font-bold">Wali Kelas</th>
                <th className="py-3 px-4 font-bold">Kapasitas</th>
                <th className="py-3 px-4 font-bold">Status</th>
                <th className="py-3 px-4 font-bold text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[#EBD4C5]/60 text-[#30251D]">
              {loading ? (
                <tr>
                  <td colSpan={9} className="py-12 text-center text-gray-400">
                    <RefreshCw size={24} className="mx-auto mb-2 animate-spin text-[#FF9800]" />
                    Memuat data kelas...
                  </td>
                </tr>
              ) : classes.length > 0 ? (
                classes.map((cls) => (
                  <tr key={cls.id} className="hover:bg-[#FFF8F5]/50 transition-colors">
                    <td className="py-3.5 px-4 font-mono">{cls.code}</td>
                    <td className="py-3.5 px-4 font-bold">{cls.name}</td>
                    <td className="py-3.5 px-4">{cls.educationLevel}</td>
                    <td className="py-3.5 px-4">{cls.department || '-'}</td>
                    <td className="py-3.5 px-4">{getAcademicYearName(cls.academicYearId)}</td>
                    <td className="py-3.5 px-4">{getTeacherName(cls.homeroomTeacherId)}</td>
                    <td className="py-3.5 px-4 text-center">{cls.capacity}</td>
                    <td className="py-3.5 px-4">
                      <span className={`inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-[10px] font-bold ${cls.status === 'active' ? 'bg-emerald-100 text-emerald-700' : cls.status === 'inactive' ? 'bg-gray-100 text-gray-600' : 'bg-gray-200 text-gray-600'}`}> {cls.status} </span>
                    </td>
                    <td className="py-3.5 px-4 text-right">
                      <div className="flex items-center justify-end gap-1.5">
                        <button type="button" onClick={() => openEdit(cls)} title="Edit" className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800]"><Edit size={15} /></button>
                        <button type="button" onClick={() => handleArchive(cls)} title="Arsipkan" className="rounded-lg p-1.5 text-gray-500 hover:bg-red-50 hover:text-red-600"><Trash2 size={15} /></button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={9} className="py-12 text-center text-gray-500">Tidak ada kelas ditemukan.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div className="flex flex-col sm:flex-row items-center justify-between border-t border-[#EBD4C5] px-4 py-3 text-xs text-gray-500 gap-2">
          <div>Menampilkan {classes.length} dari total {total} kelas</div>
          <div className="flex items-center gap-2">
            <button type="button" disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))} className="inline-flex items-center gap-1 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 font-bold text-gray-600 hover:bg-[#FFF8F5] disabled:opacity-40">
              <ArrowUpDown size={14} />
              <span>Sebelumnya</span>
            </button>
            <span className="font-bold px-2">Halaman {page}</span>
            <button type="button" disabled={classes.length < limit} onClick={() => setPage((p) => p + 1)} className="inline-flex items-center gap-1 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 font-bold text-gray-600 hover:bg-[#FFF8F5] disabled:opacity-40">
              <span>Berikutnya</span>
            </button>
          </div>
        </div>

        {/* CREATE MODAL */}
        {createModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="w-full max-w-2xl rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b pb-3 mb-4">
                <h3 className="text-lg font-bold">Tambah Kelas Baru</h3>
                <button type="button" onClick={() => setCreateModal(false)} className="text-gray-400 hover:text-gray-600"><X size={20} /></button>
              </div>
              {formError && (
                <div className="mb-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                  <AlertCircle size={16} />
                  <span>{formError}</span>
                </div>
              )}
              <form onSubmit={handleCreateSubmit} className="space-y-4 text-xs">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kode *</label>
                    <input type="text" required value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama *</label>
                    <input type="text" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tingkat Pendidikan *</label>
                    <select required value={form.educationLevel} onChange={(e) => setForm({ ...form, educationLevel: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2">
                      <option value="">Pilih Tingkat</option>
                      <option value="SD">SD</option>
                      <option value="SMP">SMP</option>
                      <option value="SMA">SMA</option>
                      <option value="SMK">SMK</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Departemen</label>
                    <input type="text" value={form.department} onChange={(e) => setForm({ ...form, department: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tahun Ajaran *</label>
                    <select
                      required
                      value={form.academicYearId}
                      onChange={(e) => setForm({ ...form, academicYearId: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2"
                      disabled={academicYearsLoading}
                    >
                      <option value="">{academicYearsLoading ? 'Memuat tahun ajaran...' : 'Pilih Tahun Ajaran'}</option>
                      {academicYears.map((ay) => (
                        <option key={ay.id} value={ay.id}>
                          {ay.name}{ay.semester ? ` (Sem ${ay.semester})` : ''}
                        </option>
                      ))}
                    </select>
                    {academicYearsError && (
                      <p className="mt-1 text-[11px] text-red-500">{academicYearsError}</p>
                    )}
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Wali Kelas</label>
                    <select value={form.homeroomTeacherId} onChange={(e) => setForm({ ...form, homeroomTeacherId: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2">
                      <option value="">Tidak ada</option>
                      {teachers.map((t) => (
                        <option key={t.id} value={t.id}>{t.name}</option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kapasitas *</label>
                    <input type="number" min={1} required value={form.capacity} onChange={(e) => setForm({ ...form, capacity: Number(e.target.value) })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Status *</label>
                    <select required value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2">
                      <option value="active">Aktif</option>
                      <option value="inactive">Tidak Aktif</option>
                    </select>
                  </div>
                </div>
                <div className="flex justify-end gap-2 mt-4">
                  <button type="button" onClick={() => setCreateModal(false)} className="px-4 py-2 text-xs bg-gray-200 rounded">Batal</button>
                  <button type="submit" disabled={submitting} className="px-4 py-2 text-xs bg-[#FF9800] text-white rounded hover:bg-[#E68600]">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* EDIT MODAL */}
        {editModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="w-full max-w-2xl rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b pb-3 mb-4">
                <h3 className="text-lg font-bold">Edit Kelas</h3>
                <button type="button" onClick={() => setEditModal(null)} className="text-gray-400 hover:text-gray-600"><X size={20} /></button>
              </div>
              {formError && (
                <div className="mb-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600">
                  <AlertCircle size={16} />
                  <span>{formError}</span>
                </div>
              )}
              <form onSubmit={handleEditSubmit} className="space-y-4 text-xs">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kode</label>
                    <input type="text" disabled value={form.code} className="w-full rounded-xl border border-[#EBD4C5] bg-gray-100 p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Nama *</label>
                    <input type="text" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tingkat Pendidikan *</label>
                    <select required value={form.educationLevel} onChange={(e) => setForm({ ...form, educationLevel: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2">
                      <option value="">Pilih Tingkat</option>
                      <option value="SD">SD</option>
                      <option value="SMP">SMP</option>
                      <option value="SMA">SMA</option>
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Departemen</label>
                    <input type="text" value={form.department} onChange={(e) => setForm({ ...form, department: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Tahun Ajaran *</label>
                    <select
                      required
                      value={form.academicYearId}
                      onChange={(e) => setForm({ ...form, academicYearId: e.target.value })}
                      className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2"
                      disabled={academicYearsLoading}
                    >
                      <option value="">{academicYearsLoading ? 'Memuat tahun ajaran...' : 'Pilih Tahun Ajaran'}</option>
                      {academicYears.map((ay) => (
                        <option key={ay.id} value={ay.id}>
                          {ay.name}{ay.semester ? ` (Sem ${ay.semester})` : ''}
                        </option>
                      ))}
                    </select>
                    {academicYearsError && (
                      <p className="mt-1 text-[11px] text-red-500">{academicYearsError}</p>
                    )}
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Wali Kelas</label>
                    <select value={form.homeroomTeacherId} onChange={(e) => setForm({ ...form, homeroomTeacherId: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2">
                      <option value="">Tidak ada</option>
                      {teachers.map((t) => (
                        <option key={t.id} value={t.id}>{t.name}</option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Kapasitas *</label>
                    <input type="number" min={1} required value={form.capacity} onChange={(e) => setForm({ ...form, capacity: Number(e.target.value) })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" />
                  </div>
                  <div>
                    <label className="block mb-1 font-bold text-gray-700">Status *</label>
                    <select required value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2">
                      <option value="active">Aktif</option>
                      <option value="inactive">Tidak Aktif</option>
                      <option value="archived">Diarsipkan</option>
                    </select>
                  </div>
                </div>
                <div className="flex justify-end gap-2 mt-4">
                  <button type="button" onClick={() => setEditModal(null)} className="px-4 py-2 text-xs bg-gray-200 rounded">Batal</button>
                  <button type="submit" disabled={submitting} className="px-4 py-2 text-xs bg-[#FF9800] text-white rounded hover:bg-[#E68600]">Simpan</button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}
