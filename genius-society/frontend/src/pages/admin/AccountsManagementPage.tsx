import { useState, useEffect, type FormEvent } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import { Plus, Edit, Trash2, X, RefreshCw, AlertCircle, Search } from 'lucide-react'

interface Account { id: string; name: string; email: string; username: string; role: string; isActive: boolean }

export default function AccountsManagementPage() {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const limit = 20
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('')
  const [createModal, setCreateModal] = useState(false)
  const [editModal, setEditModal] = useState<Account | null>(null)
  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const emptyForm = { name: '', email: '', username: '', password: '', role: 'teacher', isActive: true }
  const [form, setForm] = useState({ ...emptyForm })

  const fetchAccounts = async () => {
    setLoading(true)
    try {
      const params: any = { page, limit }
      if (roleFilter) params.role = roleFilter
      const res = await api.get<ApiResponse<Account[]>>('/users', { params })
      if (res.data.success && res.data.data) {
        setAccounts(res.data.data)
        const meta: any = res.data.meta
        if (meta) setTotal(meta.total || res.data.data.length)
      }
    } catch (e) { console.error(e) } finally { setLoading(false) }
  }
  useEffect(() => { fetchAccounts() }, [page, roleFilter])

  const errMsg = (err: unknown, fb: string) => (err as any)?.response?.data?.error?.message || (err instanceof Error ? err.message : fb)

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault(); setFormError(null); setSubmitting(true)
    try {
      const res = await api.post('/users', { name: form.name, email: form.email, username: form.username, password: form.password, role: form.role })
      if (res.data.success) { setCreateModal(false); setForm({ ...emptyForm }); fetchAccounts() }
    } catch (err) { setFormError(errMsg(err, 'Gagal membuat akun')) } finally { setSubmitting(false) }
  }
  const handleEdit = async (e: FormEvent) => {
    e.preventDefault(); if (!editModal) return; setFormError(null); setSubmitting(true)
    try {
      const payload: any = { name: form.name, email: form.email, username: form.username, isActive: form.isActive }
      if (form.password) payload.password = form.password
      const res = await api.put(`/users/${editModal.id}`, payload)
      if (res.data.success) { setEditModal(null); fetchAccounts() }
    } catch (err) { setFormError(errMsg(err, 'Gagal memperbarui')) } finally { setSubmitting(false) }
  }
  const handleDelete = async (a: Account) => {
    if (!confirm(`Hapus akun ${a.name}?`)) return
    try { await api.delete(`/users/${a.id}`); fetchAccounts() } catch { alert('Gagal') }
  }
  const openEdit = (a: Account) => { setEditModal(a); setForm({ name: a.name, email: a.email, username: a.username, password: '', role: a.role, isActive: a.isActive }) }

  const filtered = search ? accounts.filter((a) => (a.role !== 'super_admin' && a.role !== 'admin') && (a.name.toLowerCase().includes(search.toLowerCase()) || a.email.toLowerCase().includes(search.toLowerCase()) || a.username.toLowerCase().includes(search.toLowerCase()))) : accounts.filter((a) => a.role !== 'super_admin' && a.role !== 'admin')

  return (
    <AdminLayout>
      <div className="mx-auto max-w-7xl space-y-6">
        <div className="flex flex-col gap-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-sm sm:flex-row sm:items-center sm:justify-between">
          <div><h1 className="text-2xl font-black text-[#30251D]">Manajemen Akun</h1><p className="mt-1 text-xs text-gray-500">Kelola akun pengguna sekolah.</p></div>
          <button onClick={() => setCreateModal(true)} className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white hover:bg-[#E68600]"><Plus size={16} /><span>Tambah Akun</span></button>
        </div>
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-sm flex flex-col md:flex-row gap-3">
          <div className="relative flex-1"><Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" /><input value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Cari nama/email/username..." className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 pl-10 pr-4 py-2 text-xs" /></div>
          <select value={roleFilter} onChange={(e) => { setRoleFilter(e.target.value); setPage(1) }} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
            <option value="">Semua Role</option><option value="teacher">Guru</option><option value="student">Murid</option><option value="curriculum">Kurikulum</option><option value="principal">Kepala Sekolah</option>
          </select>
        </div>
        <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-sm overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 uppercase"><tr><th className="py-3 px-4 font-bold">Nama</th><th className="py-3 px-4 font-bold">Email</th><th className="py-3 px-4 font-bold">Username</th><th className="py-3 px-4 font-bold">Role</th><th className="py-3 px-4 font-bold">Status</th><th className="py-3 px-4 font-bold text-right">Aksi</th></tr></thead>
            <tbody className="divide-y divide-[#EBD4C5]/60 text-[#30251D]">
              {loading ? <tr><td colSpan={6} className="py-12 text-center text-gray-400"><RefreshCw size={24} className="mx-auto mb-2 animate-spin text-[#FF9800]" />Memuat...</td></tr>
              : filtered.length > 0 ? filtered.map((a) => (
                <tr key={a.id} className="hover:bg-[#FFF8F5]/50">
                  <td className="py-3.5 px-4 font-bold">{a.name}</td>
                  <td className="py-3.5 px-4">{a.email}</td>
                  <td className="py-3.5 px-4">{a.username}</td>
                  <td className="py-3.5 px-4">{a.role}</td>
                  <td className="py-3.5 px-4"><span className={`inline-flex rounded-full px-2.5 py-0.5 text-[10px] font-bold ${a.isActive ? 'bg-emerald-100 text-emerald-700' : 'bg-gray-200 text-gray-600'}`}>{a.isActive ? 'Aktif' : 'Nonaktif'}</span></td>
                  <td className="py-3.5 px-4 text-right"><div className="flex items-center justify-end gap-1.5"><button onClick={() => openEdit(a)} className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800]"><Edit size={15} /></button><button onClick={() => handleDelete(a)} className="rounded-lg p-1.5 text-gray-500 hover:bg-red-50 hover:text-red-600"><Trash2 size={15} /></button></div></td>
                </tr>
              )) : <tr><td colSpan={6} className="py-12 text-center text-gray-500">Tidak ada akun.</td></tr>}
            </tbody>
          </table>
        </div>
        <div className="flex items-center justify-between text-xs text-gray-500 px-4">
          <div>Total {total}</div>
          <div className="flex gap-2"><button disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))} className="px-3 py-1.5 rounded border disabled:opacity-40">Prev</button><span className="px-2 font-bold">Hal {page}</span><button disabled={accounts.length < limit} onClick={() => setPage((p) => p + 1)} className="px-3 py-1.5 rounded border disabled:opacity-40">Next</button></div>
        </div>
      </div>

      {createModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
          <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
            <div className="flex items-center justify-between border-b pb-3 mb-4"><h3 className="text-lg font-bold">Tambah Akun</h3><button onClick={() => setCreateModal(false)}><X size={20} /></button></div>
            {formError && <div className="mb-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600"><AlertCircle size={16}/><span>{formError}</span></div>}
            <form onSubmit={handleCreate} className="space-y-4 text-xs">
              <div><label className="block mb-1 font-bold text-gray-700">Nama *</label><input required value={form.name} onChange={(e)=>setForm({...form,name:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Email *</label><input type="email" required value={form.email} onChange={(e)=>setForm({...form,email:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Username *</label><input required value={form.username} onChange={(e)=>setForm({...form,username:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Password *</label><input type="password" required value={form.password} onChange={(e)=>setForm({...form,password:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Role *</label><select required value={form.role} onChange={(e)=>setForm({...form,role:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2"><option value="admin">Admin</option><option value="teacher">Guru</option><option value="student">Murid</option><option value="curriculum">Kurikulum</option><option value="principal">Kepala Sekolah</option></select></div>
              <div className="flex justify-end gap-2 mt-4"><button type="button" onClick={()=>setCreateModal(false)} className="px-4 py-2 text-xs bg-gray-200 rounded">Batal</button><button type="submit" disabled={submitting} className="px-4 py-2 text-xs bg-[#FF9800] text-white rounded hover:bg-[#E68600]">Simpan</button></div>
            </form>
          </div>
        </div>
      )}

      {editModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
          <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
            <div className="flex items-center justify-between border-b pb-3 mb-4"><h3 className="text-lg font-bold">Edit Akun</h3><button onClick={() => setEditModal(null)}><X size={20} /></button></div>
            {formError && <div className="mb-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600"><AlertCircle size={16}/><span>{formError}</span></div>}
            <form onSubmit={handleEdit} className="space-y-4 text-xs">
              <div><label className="block mb-1 font-bold text-gray-700">Nama</label><input value={form.name} onChange={(e)=>setForm({...form,name:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Email</label><input type="email" value={form.email} onChange={(e)=>setForm({...form,email:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Username</label><input value={form.username} onChange={(e)=>setForm({...form,username:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Password (kosongkan jika tidak diubah)</label><input type="password" value={form.password} onChange={(e)=>setForm({...form,password:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
              <div><label className="block mb-1 font-bold text-gray-700">Status</label><select value={form.isActive ? 'true' : 'false'} onChange={(e)=>setForm({...form,isActive:e.target.value==='true'})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2"><option value="true">Aktif</option><option value="false">Nonaktif</option></select></div>
              <div className="flex justify-end gap-2 mt-4"><button type="button" onClick={()=>setEditModal(null)} className="px-4 py-2 text-xs bg-gray-200 rounded">Batal</button><button type="submit" disabled={submitting} className="px-4 py-2 text-xs bg-[#FF9800] text-white rounded hover:bg-[#E68600]">Simpan</button></div>
            </form>
          </div>
        </div>
      )}
    </AdminLayout>
  )
}