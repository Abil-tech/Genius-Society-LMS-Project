import { useEffect, useState, type FormEvent } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { ArrowLeft, Save, Upload, FolderGit2, AlertCircle, Clock } from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  createProject,
  updateProject,
  getProjectById,
  getTeacherClasses,
  getTeacherSubjects,
} from '../../services/teacherService'
import type { TeacherClass, TeacherSubject } from '../../types/teacher'

export default function TeacherProjectFormPage() {
  const { id } = useParams<{ id: string }>()
  const isEdit = Boolean(id)
  const navigate = useNavigate()
  const toast = useToast()

  // Form states
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [instructions, setInstructions] = useState('')
  const [subjectId, setSubjectId] = useState('')
  const [selectedClassIds, setSelectedClassIds] = useState<string[]>([])
  const [deadline, setDeadline] = useState('')
  const [allowResubmit, setAllowResubmit] = useState(true)
  const [files, setFiles] = useState<File[]>([])

  // Options
  const [classes, setClasses] = useState<TeacherClass[]>([])
  const [subjects, setSubjects] = useState<TeacherSubject[]>([])
  const [pageLoading, setPageLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState<string | null>(null)
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    // Default deadline: 14 days from now
    const d = new Date()
    d.setDate(d.getDate() + 14)
    d.setHours(23, 59, 0, 0)
    setDeadline(d.toISOString().substring(0, 16))

    const init = async () => {
      setPageLoading(true)
      try {
        const [clsRes, subRes] = await Promise.all([
          getTeacherClasses(),
          getTeacherSubjects(),
        ])

        if (clsRes.success && clsRes.data) setClasses(clsRes.data)
        if (subRes.success && subRes.data) {
          setSubjects(subRes.data)
          if (!isEdit && subRes.data.length > 0) {
            setSubjectId(subRes.data[0].id)
          }
        }

        if (isEdit && id) {
          const pRes = await getProjectById(id)
          if (pRes.success && pRes.data) {
            const p = pRes.data
            setTitle(p.title)
            setDescription(p.description || '')
            setInstructions(p.instructions || '')
            setSubjectId(p.subjectId)
            setSelectedClassIds(p.classIds || [])
            setAllowResubmit(p.allowResubmit)
            if (p.deadline) {
              setDeadline(new Date(p.deadline).toISOString().substring(0, 16))
            }
          } else {
            setFormError('Projek tidak ditemukan.')
          }
        }
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : 'Gagal memuat form projek.'
        setFormError(msg)
      } finally {
        setPageLoading(false)
      }
    }

    init()
  }, [id, isEdit])

  const availableClasses = classes.filter((cls) => {
    if (!subjectId) return true
    const currentSub = subjects.find((s) => s.id === subjectId)
    if (currentSub?.classIds && currentSub.classIds.length > 0) {
      return currentSub.classIds.includes(cls.id)
    }
    return true
  })

  const handleClassToggle = (classId: string) => {
    setSelectedClassIds((prev) =>
      prev.includes(classId) ? prev.filter((i) => i !== classId) : [...prev, classId]
    )
  }

  const handleSelectAllClasses = () => {
    if (selectedClassIds.length === availableClasses.length) {
      setSelectedClassIds([])
    } else {
      setSelectedClassIds(availableClasses.map((c) => c.id))
    }
  }

  const validate = () => {
    const errors: Record<string, string> = {}
    if (!title.trim()) errors.title = 'Judul projek wajib diisi'
    if (!subjectId) errors.subjectId = 'Mata pelajaran wajib dipilih'
    if (selectedClassIds.length === 0) errors.classes = 'Pilih minimal satu kelas sasaran'
    if (!deadline) errors.deadline = 'Tenggat waktu pengumpulan projek wajib diisi'
    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    setSubmitting(true)
    setFormError(null)

    try {
      const formData = new FormData()
      formData.append('title', title)
      formData.append('description', description)
      formData.append('instructions', instructions)
      formData.append('subjectId', subjectId)
      selectedClassIds.forEach((cId) => formData.append('classIds', cId))
      formData.append('deadline', new Date(deadline).toISOString())
      formData.append('allowResubmit', String(allowResubmit))
      files.forEach((f) => formData.append('attachments', f))

      if (isEdit && id) {
        const res = await updateProject(id, formData)
        if (res.success) {
          toast.success('Projek Diperbarui', 'Perubahan rincian projek berhasil disimpan.')
          navigate(`/teacher/projects/${id}`)
        } else {
          setFormError(res.error?.message || 'Gagal menyimpan perubahan.')
        }
      } else {
        const res = await createProject(formData)
        if (res.success) {
          toast.success('Projek Dibuat', 'Projek baru telah ditugaskan ke kelas sasaran.')
          navigate('/teacher/projects')
        } else {
          setFormError(res.error?.message || 'Gagal membuat projek baru.')
        }
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat menyimpan data.'
      setFormError(msg)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <TeacherLayout>
      <div className="max-w-3xl mx-auto space-y-6">
        {/* Back Link */}
        <Link
          to={isEdit && id ? `/teacher/projects/${id}` : '/teacher/projects'}
          className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
        >
          <ArrowLeft size={16} />
          <span>Kembali</span>
        </Link>

        {/* Header */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-[#FF9800]">
              <FolderGit2 size={20} />
            </div>
            <div>
              <h1 className="text-lg sm:text-xl font-extrabold text-[#30251D]">
                {isEdit ? 'Ubah Rincian Projek' : 'Buat Projek Pembelajaran Baru'}
              </h1>
              <p className="text-xs text-[#687385]">
                Tentukan panduan projek, rubrik penilaian, tenggat waktu, dan dokumen pendukung.
              </p>
            </div>
          </div>
        </div>

        {formError && (
          <ErrorState title="Periksa Kembali Form" message={formError} compact />
        )}

        {pageLoading ? (
          <div className="space-y-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 animate-pulse">
            <div className="h-10 bg-gray-200 rounded-xl" />
            <div className="h-28 bg-gray-200 rounded-xl" />
            <div className="h-10 bg-gray-200 rounded-xl" />
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-6 rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
            {/* Judul Projek */}
            <div>
              <label htmlFor="title" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Judul Projek <span className="text-rose-500">*</span>
              </label>
              <input
                id="title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Contoh: Projek Sains - Pembuatan Miniatur Ekosistem Terarium"
                className={`w-full rounded-xl border px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:outline-none focus:ring-1 ${
                  validationErrors.title
                    ? 'border-rose-400 focus:border-rose-500 focus:ring-rose-500'
                    : 'border-[#EBD4C5] focus:border-[#FF9800] focus:ring-[#FF9800]'
                }`}
              />
              {validationErrors.title && (
                <p className="mt-1 text-xs text-rose-600 flex items-center gap-1">
                  <AlertCircle size={12} />
                  <span>{validationErrors.title}</span>
                </p>
              )}
            </div>

            {/* Mata Pelajaran */}
            <div>
              <label htmlFor="subjectId" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Mata Pelajaran <span className="text-rose-500">*</span>
              </label>
              <select
                id="subjectId"
                value={subjectId}
                onChange={(e) => {
                  setSubjectId(e.target.value)
                  setSelectedClassIds([])
                }}
                className={`w-full rounded-xl border px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] bg-white focus:outline-none focus:ring-1 ${
                  validationErrors.subjectId
                    ? 'border-rose-400 focus:border-rose-500 focus:ring-rose-500'
                    : 'border-[#EBD4C5] focus:border-[#FF9800] focus:ring-[#FF9800]'
                }`}
              >
                <option value="">-- Pilih Mata Pelajaran --</option>
                {subjects.map((sub) => (
                  <option key={sub.id} value={sub.id}>
                    {sub.name} ({sub.code})
                  </option>
                ))}
              </select>
              {validationErrors.subjectId && (
                <p className="mt-1 text-xs text-rose-600 flex items-center gap-1">
                  <AlertCircle size={12} />
                  <span>{validationErrors.subjectId}</span>
                </p>
              )}
            </div>

            {/* Deskripsi */}
            <div>
              <label htmlFor="description" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Ringkasan Projek
              </label>
              <textarea
                id="description"
                rows={2}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Ringkasan singkat mengenai tujuan pembelajaran projek..."
                className="w-full rounded-xl border border-[#EBD4C5] px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
              />
            </div>

            {/* Instruksi & Panduan Pelaksanaan */}
            <div>
              <label htmlFor="instructions" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Instruksi & Rubrik Penilaian
              </label>
              <textarea
                id="instructions"
                rows={4}
                value={instructions}
                onChange={(e) => setInstructions(e.target.value)}
                placeholder="Jelaskan tahapan pengerjaan, format laporan, kriteria penilaian (rubrik), dll..."
                className="w-full rounded-xl border border-[#EBD4C5] px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
              />
            </div>

            {/* Kelas Sasaran */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <label className="block text-xs font-bold uppercase tracking-wider text-[#30251D]">
                  Kelas Sasaran <span className="text-rose-500">*</span>
                </label>
                {availableClasses.length > 0 && (
                  <button
                    type="button"
                    onClick={handleSelectAllClasses}
                    className="text-xs font-bold text-[#FF9800] hover:underline"
                  >
                    {selectedClassIds.length === availableClasses.length
                      ? 'Batalkan Semua'
                      : 'Pilih Semua'}
                  </button>
                )}
              </div>

              {availableClasses.length === 0 ? (
                <p className="text-xs text-gray-500 italic p-3 bg-[#FFF8F5] rounded-xl border border-[#EBD4C5]">
                  Pilih mata pelajaran terlebih dahulu untuk melihat daftar kelas yang diajar.
                </p>
              ) : (
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                  {availableClasses.map((cls) => {
                    const checked = selectedClassIds.includes(cls.id)
                    return (
                      <button
                        key={cls.id}
                        type="button"
                        onClick={() => handleClassToggle(cls.id)}
                        className={`flex items-center gap-2 rounded-xl border p-2.5 text-xs font-semibold text-left transition-colors ${
                          checked
                            ? 'border-[#FF9800] bg-[#FEEADC] text-[#FF9800]'
                            : 'border-[#EBD4C5] bg-white text-[#30251D] hover:bg-[#FFF8F5]'
                        }`}
                      >
                        <input
                          type="checkbox"
                          checked={checked}
                          onChange={() => {}}
                          className="rounded text-[#FF9800] focus:ring-[#FF9800]"
                        />
                        <span className="truncate">{cls.name}</span>
                      </button>
                    )
                  })}
                </div>
              )}
              {validationErrors.classes && (
                <p className="mt-1 text-xs text-rose-600 flex items-center gap-1">
                  <AlertCircle size={12} />
                  <span>{validationErrors.classes}</span>
                </p>
              )}
            </div>

            {/* Tenggat Waktu & Opsi Resubmit */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label htmlFor="deadline" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                  Tenggat Waktu Projek <span className="text-rose-500">*</span>
                </label>
                <div className="relative">
                  <Clock size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    id="deadline"
                    type="datetime-local"
                    value={deadline}
                    onChange={(e) => setDeadline(e.target.value)}
                    className="w-full rounded-xl border border-[#EBD4C5] py-2.5 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] bg-white focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
                  />
                </div>
                {validationErrors.deadline && (
                  <p className="mt-1 text-xs text-rose-600 flex items-center gap-1">
                    <AlertCircle size={12} />
                    <span>{validationErrors.deadline}</span>
                  </p>
                )}
              </div>

              <div>
                <span className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                  Kebijakan Revisi
                </span>
                <label className="flex items-center gap-3 p-3 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={allowResubmit}
                    onChange={(e) => setAllowResubmit(e.target.checked)}
                    className="h-4 w-4 rounded text-[#FF9800] focus:ring-[#FF9800]"
                  />
                  <span className="text-xs font-bold text-[#30251D]">
                    Izinkan revisi berkas projek sebelum batas tenggat berakhir
                  </span>
                </label>
              </div>
            </div>

            {/* Lampiran Panduan Projek */}
            <div>
              <label className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Lampiran Panduan / Lembar Kerja Projek (Cloudflare R2)
              </label>
              <div className="rounded-xl border border-dashed border-[#EBD4C5] p-5 text-center bg-[#FFF8F5]/60">
                <Upload size={24} className="mx-auto text-[#FF9800]" />
                <p className="mt-2 text-xs font-semibold text-[#30251D]">
                  {files.length > 0
                    ? `${files.length} file dipilih: ${files.map((f) => f.name).join(', ')}`
                    : 'Pilih berkas panduan projek dari komputer'}
                </p>
                <p className="text-[11px] text-[#687385] mt-0.5">
                  Mendukung PDF, Word, PowerPoint, ZIP (Maks 100MB)
                </p>
                <label className="mt-3 inline-block cursor-pointer rounded-xl bg-white border border-[#EBD4C5] px-4 py-2 text-xs font-bold text-[#30251D] hover:bg-[#FFF8F5] shadow-xs">
                  <span>Pilih File</span>
                  <input
                    type="file"
                    multiple
                    className="hidden"
                    onChange={(e) => {
                      if (e.target.files) {
                        setFiles(Array.from(e.target.files))
                      }
                    }}
                  />
                </label>
              </div>
            </div>

            {/* Submit Actions */}
            <div className="flex items-center justify-end gap-3 pt-4 border-t border-[#EBD4C5]">
              <Link
                to="/teacher/projects"
                className="rounded-xl border border-[#EBD4C5] bg-white px-5 py-2.5 text-xs sm:text-sm font-semibold text-[#30251D] hover:bg-[#FFF8F5] transition-colors"
              >
                Batal
              </Link>
              <button
                type="submit"
                disabled={submitting}
                className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-6 py-2.5 text-xs sm:text-sm font-bold text-white shadow-xs hover:bg-[#E68A00] disabled:opacity-50 transition-colors"
              >
                {submitting ? (
                  <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                ) : (
                  <Save size={16} />
                )}
                <span>{isEdit ? 'Simpan Perubahan' : 'Terbitkan Projek'}</span>
              </button>
            </div>
          </form>
        )}
      </div>
    </TeacherLayout>
  )
}
