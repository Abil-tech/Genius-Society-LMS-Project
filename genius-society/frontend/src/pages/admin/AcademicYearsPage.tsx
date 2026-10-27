import { useState, useEffect, type FormEvent } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import { Plus, Edit, Trash2, X, RefreshCw, AlertCircle, Layers, GraduationCap, Clock, FileText, ArrowUpDown } from 'lucide-react'

interface AcademicYear { id: string; name: string; semester: number; startDate: string; endDate: string; isActive: boolean; status: string }
interface PromoCandidate { studentId: string; nis: string; name: string; sourceClassId: string; sourceClassName: string; educationLevel: string; eligibilityStatus: string; recommendedAction: string; reasonCode?: string; reason?: string; suggestedClassId?: string; targetClassId?: string; averageScore: number; missingGradesCount: number }
interface GradCandidate { studentId: string; nis: string; name: string; classId: string; className: string; educationLevel: string; enrollmentYear: number; eligibilityStatus: string; recommendedAction: string; reasonCode?: string; reason?: string; averageScore: number; missingGradesCount: number }
interface SuspendedDecision { id: string; studentId: string; studentName: string; nis: string; sourceClassName: string; sourceAcademicYearId: string; targetAcademicYearId?: string; processType: string; eligibilityStatus: string; decisionStatus: string; reasonCode: string; reason: string; reviewStatus: string; reviewedByName?: string; reviewedAt?: string; decidedAt: string; createdAt: string }
interface Job { id: string; processType: string; sourceAcademicYearId: string; sourceAcademicYearName: string; targetAcademicYearId?: string; targetAcademicYearName?: string; status: string; summary: { totalStudents: number; promoted: number; retained: number; graduated: number; suspended: number; excluded: number; failed: number }; executedByName: string; errorMessage?: string; startedAt: string; completedAt?: string }

type Tab = 'years' | 'promotion' | 'graduation' | 'decisions'

export default function AcademicYearsPage() {
  const [tab, setTab] = useState<Tab>('years')

  // Years
  const [years, setYears] = useState<AcademicYear[]>([])
  const [loading, setLoading] = useState(true)
  const [formError, setFormError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [createModal, setCreateModal] = useState(false)
  const [editModal, setEditModal] = useState<AcademicYear | null>(null)
  const emptyForm = { name: '', semester: 1, startDate: '', endDate: '' }
  const [form, setForm] = useState({ ...emptyForm })

  const fetchYears = async () => {
    setLoading(true)
    try { const res = await api.get<ApiResponse<AcademicYear[]>>('/academic-years'); if (res.data.success && res.data.data) setYears(res.data.data) }
    catch (e) { console.error(e) } finally { setLoading(false) }
  }
  useEffect(() => { fetchYears() }, [])

  const errMsg = (err: unknown, fb: string) => (err as any)?.response?.data?.error?.message || (err instanceof Error ? err.message : fb)

  const handleCreate = async (e: FormEvent) => {
    e.preventDefault(); setFormError(null); setSubmitting(true)
    try {
      const res = await api.post('/academic-years', { name: form.name, semester: form.semester, startDate: new Date(form.startDate).toISOString(), endDate: new Date(form.endDate).toISOString() })
      if (res.data.success) { setCreateModal(false); setForm({ ...emptyForm }); fetchYears() }
    } catch (err) { setFormError(errMsg(err, 'Gagal membuat')) } finally { setSubmitting(false) }
  }
  const handleEdit = async (e: FormEvent) => {
    e.preventDefault(); if (!editModal) return; setFormError(null); setSubmitting(true)
    try {
      const payload: any = {}
      if (form.name) payload.name = form.name
      if (form.semester) payload.semester = form.semester
      if (form.startDate) payload.startDate = new Date(form.startDate).toISOString()
      if (form.endDate) payload.endDate = new Date(form.endDate).toISOString()
      const res = await api.put(`/academic-years/${editModal.id}`, payload)
      if (res.data.success) { setEditModal(null); fetchYears() }
    } catch (err) { setFormError(errMsg(err, 'Gagal memperbarui')) } finally { setSubmitting(false) }
  }
  const handleActivate = async (y: AcademicYear) => { if (!confirm(`Aktifkan ${y.name}?`)) return; try { await api.patch(`/academic-years/${y.id}/activate`); fetchYears() } catch { alert('Gagal') } }
  const handleArchive = async (y: AcademicYear) => { if (!confirm(`Arsipkan ${y.name}?`)) return; try { await api.patch(`/academic-years/${y.id}/archive`); fetchYears() } catch { alert('Gagal') } }
  const handleDelete = async (y: AcademicYear) => { if (!confirm(`Hapus ${y.name}?`)) return; try { await api.delete(`/academic-years/${y.id}`); fetchYears() } catch { alert('Gagal') } }
  const openEdit = (y: AcademicYear) => { setEditModal(y); setForm({ name: y.name, semester: y.semester, startDate: y.startDate.slice(0, 10), endDate: y.endDate.slice(0, 10) }) }

  // Promotion
  const [promoSource, setPromoSource] = useState('')
  const [promoTarget, setPromoTarget] = useState('')
  const [promoCandidates, setPromoCandidates] = useState<PromoCandidate[]>([])
  const [promoPreview, setPromoPreview] = useState<any>(null)
  const [promoJobs, setPromoJobs] = useState<Job[]>([])
  const [promoLoading, setPromoLoading] = useState(false)
  const fetchPromoCandidates = async () => {
    if (!promoSource) return
    setPromoLoading(true)
    try { const res = await api.get<ApiResponse<PromoCandidate[]>>(`/academic-years/${promoSource}/promotion/candidates`); if (res.data.success && res.data.data) setPromoCandidates(res.data.data) }
    catch (e) { console.error(e) } finally { setPromoLoading(false) }
  }
  const runPromoPreview = async () => {
    if (!promoSource || !promoTarget) return
    setPromoLoading(true)
    try {
      const res = await api.post<ApiResponse<any>>(`/academic-years/${promoSource}/promotion/preview`, { sourceAcademicYearId: promoSource, targetAcademicYearId: promoTarget })
      if (res.data.success && res.data.data) setPromoPreview(res.data.data)
    } catch (e) { console.error(e) } finally { setPromoLoading(false) }
  }
  const runPromoExecute = async () => {
    if (!promoSource || !promoTarget || !promoCandidates.length) return
    if (!confirm('Eksekusi kenaikan kelas?')) return
    setPromoLoading(true)
    try {
      const decisions = promoCandidates.map((c) => ({ studentId: c.studentId, action: c.recommendedAction || 'promoted', targetClassId: c.targetClassId || c.suggestedClassId }))
      const res = await api.post<ApiResponse<any>>(`/academic-years/${promoSource}/promotion/execute`, { sourceAcademicYearId: promoSource, targetAcademicYearId: promoTarget, decisions, idempotencyKey: crypto.randomUUID() })
      if (res.data.success) { alert('Selesai'); fetchPromoHistory() }
    } catch (e) { console.error(e); alert('Gagal') } finally { setPromoLoading(false) }
  }
  const fetchPromoHistory = async () => {
    if (!promoSource) return
    try { const res = await api.get<ApiResponse<Job[]>>(`/academic-years/${promoSource}/promotion/history`); if (res.data.success && res.data.data) setPromoJobs(res.data.data) } catch (e) { console.error(e) }
  }
  useEffect(() => { if (tab === 'promotion' && promoSource) { fetchPromoCandidates(); fetchPromoHistory() } }, [tab, promoSource])

  // Graduation
  const [gradYear, setGradYear] = useState('')
  const [gradCandidates, setGradCandidates] = useState<GradCandidate[]>([])
  const [gradPreview, setGradPreview] = useState<any>(null)
  const [gradJobs, setGradJobs] = useState<Job[]>([])
  const [gradLoading, setGradLoading] = useState(false)
  const fetchGradCandidates = async () => {
    if (!gradYear) return
    setGradLoading(true)
    try { const res = await api.get<ApiResponse<GradCandidate[]>>(`/academic-years/${gradYear}/graduation/candidates`); if (res.data.success && res.data.data) setGradCandidates(res.data.data) }
    catch (e) { console.error(e) } finally { setGradLoading(false) }
  }
  const runGradPreview = async () => {
    if (!gradYear) return
    setGradLoading(true)
    try {
      const res = await api.post<ApiResponse<any>>(`/academic-years/${gradYear}/graduation/preview`, { academicYearId: gradYear })
      if (res.data.success && res.data.data) setGradPreview(res.data.data)
    } catch (e) { console.error(e) } finally { setGradLoading(false) }
  }
  const runGradExecute = async () => {
    if (!gradYear || !gradCandidates.length) return
    if (!confirm('Eksekusi kelulusan?')) return
    setGradLoading(true)
    try {
      const decisions = gradCandidates.map((c) => ({ studentId: c.studentId, action: c.recommendedAction || 'graduated' }))
      const res = await api.post<ApiResponse<any>>(`/academic-years/${gradYear}/graduation/execute`, { academicYearId: gradYear, decisions, idempotencyKey: crypto.randomUUID() })
      if (res.data.success) { alert('Selesai'); fetchGradHistory() }
    } catch (e) { console.error(e); alert('Gagal') } finally { setGradLoading(false) }
  }
  const fetchGradHistory = async () => {
    if (!gradYear) return
    try { const res = await api.get<ApiResponse<Job[]>>(`/academic-years/${gradYear}/graduation/history`); if (res.data.success && res.data.data) setGradJobs(res.data.data) } catch (e) { console.error(e) }
  }
  useEffect(() => { if (tab === 'graduation' && gradYear) { fetchGradCandidates(); fetchGradHistory() } }, [tab, gradYear])

  // Decisions
  const [decisions, setDecisions] = useState<SuspendedDecision[]>([])
  const [decLoading, setDecLoading] = useState(false)
  const [selectedDecision, setSelectedDecision] = useState<SuspendedDecision | null>(null)
  const [reviewNotes, setReviewNotes] = useState('')
  const [resolveAction, setResolveAction] = useState('promoted')
  const [resolveTarget, setResolveTarget] = useState('')
  const [resolveNotes, setResolveNotes] = useState('')
  const fetchDecisions = async () => {
    if (!promoSource) return
    setDecLoading(true)
    try { const res = await api.get<ApiResponse<SuspendedDecision[]>>(`/academic-years/${promoSource}/decisions/suspended`); if (res.data.success && res.data.data) setDecisions(res.data.data) }
    catch (e) { console.error(e) } finally { setDecLoading(false) }
  }
  const submitReview = async () => {
    if (!selectedDecision) return
    try { await api.post(`/academic-decisions/${selectedDecision.id}/review`, { reviewNotes }); alert('Reviewed'); fetchDecisions() }
    catch { alert('Gagal') }
  }
  const submitResolve = async () => {
    if (!selectedDecision) return
    try { await api.post(`/academic-decisions/${selectedDecision.id}/resolve`, { action: resolveAction, targetClassId: resolveTarget, notes: resolveNotes }); alert('Resolved'); setSelectedDecision(null); fetchDecisions() }
    catch { alert('Gagal') }
  }
  useEffect(() => { if (tab === 'decisions') fetchDecisions() }, [tab, promoSource])

  const yearOptions = years

  return (
    <AdminLayout>
      <div className="mx-auto max-w-7xl space-y-6">
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-black text-[#30251D]">Tahun Ajaran</h1>
          <p className="mt-1 text-xs text-gray-500">Kelola tahun ajaran, kenaikan kelas, kelulusan, dan penangguhan.</p>
        </div>

        <div className="flex gap-2 border-b border-[#EBD4C5]">
          {[
            { id: 'years', label: 'Tahun Ajaran', icon: Layers },
            { id: 'promotion', label: 'Kenaikan Kelas', icon: ArrowUpDown },
            { id: 'graduation', label: 'Kelulusan', icon: GraduationCap },
            { id: 'decisions', label: 'Penangguhan', icon: Clock },
          ].map((t) => {
            const Icon = t.icon
            return <button key={t.id} onClick={() => setTab(t.id as Tab)} className={`inline-flex items-center gap-2 px-4 py-2 text-xs font-bold border-b-2 ${tab === t.id ? 'border-[#FF9800] text-[#FF9800]' : 'border-transparent text-gray-500'}`}><Icon size={14} /> {t.label}</button>
          })}
        </div>

        {tab === 'years' && (
          <div className="space-y-4">
            <div className="flex justify-end">
              <button onClick={() => setCreateModal(true)} className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white hover:bg-[#E68600]"><Plus size={16} /> Tambah</button>
            </div>
            <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-sm overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 uppercase">
                  <tr><th className="py-3 px-4 font-bold">Nama</th><th className="py-3 px-4 font-bold">Sem</th><th className="py-3 px-4 font-bold">Mulai</th><th className="py-3 px-4 font-bold">Berakhir</th><th className="py-3 px-4 font-bold">Status</th><th className="py-3 px-4 font-bold text-right">Aksi</th></tr>
                </thead>
                <tbody className="divide-y divide-[#EBD4C5]/60 text-[#30251D]">
                  {loading ? <tr><td colSpan={6} className="py-12 text-center text-gray-400"><RefreshCw size={24} className="mx-auto mb-2 animate-spin text-[#FF9800]" />Memuat...</td></tr>
                  : years.length > 0 ? years.map((y) => (
                    <tr key={y.id} className="hover:bg-[#FFF8F5]/50">
                      <td className="py-3.5 px-4 font-bold">{y.name}</td>
                      <td className="py-3.5 px-4">{y.semester}</td>
                      <td className="py-3.5 px-4">{new Date(y.startDate).toLocaleDateString('id-ID')}</td>
                      <td className="py-3.5 px-4">{new Date(y.endDate).toLocaleDateString('id-ID')}</td>
                      <td className="py-3.5 px-4"><span className={`inline-flex rounded-full px-2.5 py-0.5 text-[10px] font-bold ${y.status === 'active' ? 'bg-emerald-100 text-emerald-700' : y.status === 'archived' ? 'bg-gray-200 text-gray-600' : 'bg-gray-100 text-gray-600'}`}>{y.status}</span></td>
                      <td className="py-3.5 px-4 text-right">
                        <div className="flex items-center justify-end gap-1.5">
                          <button onClick={() => openEdit(y)} className="rounded-lg p-1.5 text-gray-500 hover:bg-[#FFF8F5] hover:text-[#FF9800]"><Edit size={15} /></button>
                          {y.status !== 'active' && <button onClick={() => handleActivate(y)} className="rounded-lg p-1.5 text-gray-500 hover:bg-emerald-50 hover:text-emerald-600"><FileText size={15} /></button>}
                          {y.status !== 'active' && <button onClick={() => handleArchive(y)} className="rounded-lg p-1.5 text-gray-500 hover:bg-gray-100 hover:text-gray-600"><Trash2 size={15} /></button>}
                          {y.status !== 'active' && <button onClick={() => handleDelete(y)} className="rounded-lg p-1.5 text-gray-500 hover:bg-red-50 hover:text-red-600"><Trash2 size={15} /></button>}
                        </div>
                      </td>
                    </tr>
                  )) : <tr><td colSpan={6} className="py-12 text-center text-gray-500">Tidak ada tahun ajaran.</td></tr>}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {tab === 'promotion' && (
          <div className="space-y-4">
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 flex flex-wrap gap-3 items-end">
              <div><label className="block text-xs font-bold mb-1">Tahun Sumber</label>
                <select value={promoSource} onChange={(e) => setPromoSource(e.target.value)} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
                  <option value="">Pilih</option>{yearOptions.map((y) => <option key={y.id} value={y.id}>{y.name}</option>)}
                </select></div>
              <div><label className="block text-xs font-bold mb-1">Tahun Tujuan</label>
                <select value={promoTarget} onChange={(e) => setPromoTarget(e.target.value)} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
                  <option value="">Pilih</option>{yearOptions.map((y) => <option key={y.id} value={y.id}>{y.name}</option>)}
                </select></div>
              <button onClick={fetchPromoCandidates} className="rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white">Ambil Kandidat</button>
              <button onClick={runPromoPreview} className="rounded-xl border border-[#EBD4C5] px-4 py-2 text-xs font-bold">Preview</button>
              <button onClick={runPromoExecute} className="rounded-xl bg-emerald-600 px-4 py-2 text-xs font-bold text-white">Eksekusi</button>
            </div>
            {promoPreview && <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 text-xs">Preview: total {promoPreview.totalStudents}, eligible {promoPreview.eligible}, pending {promoPreview.pendingReview}, promoted {promoPreview.promoted}, retained {promoPreview.retained}</div>}
            <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-sm overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 uppercase"><tr><th className="py-3 px-4">Nama</th><th className="py-3 px-4">Kelas</th><th className="py-3 px-4">Eligibility</th><th className="py-3 px-4">Rekomendasi</th><th className="py-3 px-4">Skor</th></tr></thead>
                <tbody className="divide-y divide-[#EBD4C5]/60">
                  {promoLoading ? <tr><td colSpan={5} className="py-8 text-center text-gray-400">Memuat...</td></tr>
                  : promoCandidates.length ? promoCandidates.map((c) => (
                    <tr key={c.studentId}><td className="py-2 px-4 font-bold">{c.name}</td><td className="py-2 px-4">{c.sourceClassName}</td><td className="py-2 px-4">{c.eligibilityStatus}</td><td className="py-2 px-4">{c.recommendedAction}</td><td className="py-2 px-4">{c.averageScore}</td></tr>
                  )) : <tr><td colSpan={5} className="py-8 text-center text-gray-500">Tidak ada kandidat.</td></tr>}
                </tbody>
              </table>
            </div>
            {promoJobs.length > 0 && <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 text-xs"><h3 className="font-bold mb-2">Riwayat Job</h3>{promoJobs.map((j) => <div key={j.id} className="border-b py-1">{j.status} - promoted {j.summary.promoted}/{j.summary.totalStudents}</div>)}</div>}
          </div>
        )}

        {tab === 'graduation' && (
          <div className="space-y-4">
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 flex flex-wrap gap-3 items-end">
              <div><label className="block text-xs font-bold mb-1">Tahun Ajaran</label>
                <select value={gradYear} onChange={(e) => setGradYear(e.target.value)} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
                  <option value="">Pilih</option>{yearOptions.map((y) => <option key={y.id} value={y.id}>{y.name}</option>)}
                </select></div>
              <button onClick={fetchGradCandidates} className="rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white">Ambil Kandidat</button>
              <button onClick={runGradPreview} className="rounded-xl border border-[#EBD4C5] px-4 py-2 text-xs font-bold">Preview</button>
              <button onClick={runGradExecute} className="rounded-xl bg-emerald-600 px-4 py-2 text-xs font-bold text-white">Eksekusi</button>
            </div>
            {gradPreview && <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 text-xs">Preview: total {gradPreview.totalStudents}, eligible {gradPreview.eligible}, graduated {gradPreview.graduated}</div>}
            <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-sm overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 uppercase"><tr><th className="py-3 px-4">Nama</th><th className="py-3 px-4">Kelas</th><th className="py-3 px-4">Eligibility</th><th className="py-3 px-4">Rekomendasi</th></tr></thead>
                <tbody className="divide-y divide-[#EBD4C5]/60">
                  {gradLoading ? <tr><td colSpan={4} className="py-8 text-center text-gray-400">Memuat...</td></tr>
                  : gradCandidates.length ? gradCandidates.map((c) => (
                    <tr key={c.studentId}><td className="py-2 px-4 font-bold">{c.name}</td><td className="py-2 px-4">{c.className}</td><td className="py-2 px-4">{c.eligibilityStatus}</td><td className="py-2 px-4">{c.recommendedAction}</td></tr>
                  )) : <tr><td colSpan={4} className="py-8 text-center text-gray-500">Tidak ada kandidat.</td></tr>}
                </tbody>
              </table>
            </div>
            {gradJobs.length > 0 && <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 text-xs"><h3 className="font-bold mb-2">Riwayat Job</h3>{gradJobs.map((j) => <div key={j.id} className="border-b py-1">{j.status} - graduated {j.summary.graduated}/{j.summary.totalStudents}</div>)}</div>}
          </div>
        )}

        {tab === 'decisions' && (
          <div className="space-y-4">
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 flex flex-wrap gap-3 items-end">
              <div><label className="block text-xs font-bold mb-1">Tahun Sumber</label>
                <select value={promoSource} onChange={(e) => setPromoSource(e.target.value)} className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs">
                  <option value="">Pilih</option>{yearOptions.map((y) => <option key={y.id} value={y.id}>{y.name}</option>)}
                </select></div>
              <button onClick={fetchDecisions} className="rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white">Muat</button>
            </div>
            <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-sm overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-600 uppercase"><tr><th className="py-3 px-4">Siswa</th><th className="py-3 px-4">Proses</th><th className="py-3 px-4">Status</th><th className="py-3 px-4">Review</th><th className="py-3 px-4">Aksi</th></tr></thead>
                <tbody className="divide-y divide-[#EBD4C5]/60">
                  {decLoading ? <tr><td colSpan={5} className="py-8 text-center text-gray-400">Memuat...</td></tr>
                  : decisions.length ? decisions.map((d) => (
                    <tr key={d.id}><td className="py-2 px-4 font-bold">{d.studentName}</td><td className="py-2 px-4">{d.processType}</td><td className="py-2 px-4">{d.decisionStatus}</td><td className="py-2 px-4">{d.reviewStatus}</td><td className="py-2 px-4"><button onClick={() => setSelectedDecision(d)} className="text-[#FF9800] font-bold">Detail</button></td></tr>
                  )) : <tr><td colSpan={5} className="py-8 text-center text-gray-500">Tidak ada penangguhan.</td></tr>}
                </tbody>
              </table>
            </div>
            {selectedDecision && (
              <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                <div className="w-full max-w-lg rounded-2xl bg-white p-6 space-y-3 text-xs">
                  <div className="flex justify-between"><h3 className="font-bold text-lg">{selectedDecision.studentName}</h3><button onClick={() => setSelectedDecision(null)}><X size={18} /></button></div>
                  <p>Alasan: {selectedDecision.reason}</p>
                  <div><label className="font-bold">Catatan Review</label><textarea value={reviewNotes} onChange={(e) => setReviewNotes(e.target.value)} className="w-full border rounded p-2" /></div>
                  <button onClick={submitReview} className="bg-[#FF9800] text-white px-3 py-1.5 rounded">Kirim Review</button>
                  <div><label className="font-bold">Resolusi</label>
                    <select value={resolveAction} onChange={(e) => setResolveAction(e.target.value)} className="border rounded p-1.5"><option value="promoted">promoted</option><option value="retained">retained</option><option value="graduated">graduated</option><option value="excluded">excluded</option></select>
                    <input placeholder="Target Class ID" value={resolveTarget} onChange={(e) => setResolveTarget(e.target.value)} className="border rounded p-1.5 ml-2" />
                    <input placeholder="Notes" value={resolveNotes} onChange={(e) => setResolveNotes(e.target.value)} className="border rounded p-1.5 ml-2" />
                  </div>
                  <button onClick={submitResolve} className="bg-emerald-600 text-white px-3 py-1.5 rounded">Resolve</button>
                </div>
              </div>
            )}
          </div>
        )}

        {/* CREATE/EDIT modals same as before */}
        {createModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b pb-3 mb-4"><h3 className="text-lg font-bold">Tambah Tahun Ajaran</h3><button onClick={() => setCreateModal(false)}><X size={20} /></button></div>
              {formError && <div className="mb-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600"><AlertCircle size={16}/><span>{formError}</span></div>}
              <form onSubmit={handleCreate} className="space-y-4 text-xs">
                <div><label className="block mb-1 font-bold text-gray-700">Nama *</label><input required value={form.name} onChange={(e)=>setForm({...form,name:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
                <div><label className="block mb-1 font-bold text-gray-700">Semester *</label><select required value={form.semester} onChange={(e)=>setForm({...form,semester:Number(e.target.value)})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2"><option value={1}>1</option><option value={2}>2</option></select></div>
                <div><label className="block mb-1 font-bold text-gray-700">Tanggal Mulai *</label><input type="date" required value={form.startDate} onChange={(e)=>setForm({...form,startDate:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
                <div><label className="block mb-1 font-bold text-gray-700">Tanggal Berakhir *</label><input type="date" required value={form.endDate} onChange={(e)=>setForm({...form,endDate:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
                <div className="flex justify-end gap-2 mt-4"><button type="button" onClick={()=>setCreateModal(false)} className="px-4 py-2 text-xs bg-gray-200 rounded">Batal</button><button type="submit" disabled={submitting} className="px-4 py-2 text-xs bg-[#FF9800] text-white rounded hover:bg-[#E68600]">Simpan</button></div>
              </form>
            </div>
          </div>
        )}
        {editModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 overflow-y-auto">
            <div className="w-full max-w-md rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xl text-[#30251D]">
              <div className="flex items-center justify-between border-b pb-3 mb-4"><h3 className="text-lg font-bold">Edit Tahun Ajaran</h3><button onClick={() => setEditModal(null)}><X size={20} /></button></div>
              {formError && <div className="mb-3 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-3 text-xs text-red-600"><AlertCircle size={16}/><span>{formError}</span></div>}
              <form onSubmit={handleEdit} className="space-y-4 text-xs">
                <div><label className="block mb-1 font-bold text-gray-700">Nama</label><input value={form.name} onChange={(e)=>setForm({...form,name:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
                <div><label className="block mb-1 font-bold text-gray-700">Semester</label><select value={form.semester} onChange={(e)=>setForm({...form,semester:Number(e.target.value)})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2"><option value={1}>1</option><option value={2}>2</option></select></div>
                <div><label className="block mb-1 font-bold text-gray-700">Tanggal Mulai</label><input type="date" value={form.startDate} onChange={(e)=>setForm({...form,startDate:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
                <div><label className="block mb-1 font-bold text-gray-700">Tanggal Berakhir</label><input type="date" value={form.endDate} onChange={(e)=>setForm({...form,endDate:e.target.value})} className="w-full rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-2" /></div>
                <div className="flex justify-end gap-2 mt-4"><button type="button" onClick={()=>setEditModal(null)} className="px-4 py-2 text-xs bg-gray-200 rounded">Batal</button><button type="submit" disabled={submitting} className="px-4 py-2 text-xs bg-[#FF9800] text-white rounded hover:bg-[#E68600]">Simpan</button></div>
              </form>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  )
}