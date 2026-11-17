import { useEffect, useState, type FormEvent } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import {
  ArrowLeft,
  Save,
  Plus,
  Trash2,
  CheckSquare,
  AlertCircle,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  createAssessment,
  updateAssessment,
  getAssessmentById,
  getTeacherClasses,
  getTeacherSubjects,
} from '../../services/teacherService'
import type {
  TeacherClass,
  TeacherSubject,
  AssessmentCategory,
  AssessmentStatus,
  AssessmentQuestion,
  CreateAssessmentPayload,
} from '../../types/teacher'

export default function TeacherAssessmentFormPage() {
  const { id } = useParams<{ id: string }>()
  const isEdit = Boolean(id)
  const navigate = useNavigate()
  const toast = useToast()

  // Basic Info States
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [durationMinutes, setDurationMinutes] = useState(60)
  const [subjectId, setSubjectId] = useState('')
  const [selectedClassIds, setSelectedClassIds] = useState<string[]>([])
  const [category, setCategory] = useState<AssessmentCategory>('daily')
  const [status, setStatus] = useState<AssessmentStatus>('draft')
  const [showResult, setShowResult] = useState(true)
  const [startAt, setStartAt] = useState('')
  const [endAt, setEndAt] = useState('')

  // Questions Builder
  const [questions, setQuestions] = useState<AssessmentQuestion[]>([
    {
      type: 'mcq',
      text: '',
      options: ['', '', '', ''],
      correctOption: 0,
      weight: 10,
      order: 1,
    },
  ])

  // Options
  const [classes, setClasses] = useState<TeacherClass[]>([])
  const [subjects, setSubjects] = useState<TeacherSubject[]>([])
  const [pageLoading, setPageLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState<string | null>(null)
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    // Default schedule
    const now = new Date()
    const later = new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000)
    setStartAt(now.toISOString().substring(0, 16))
    setEndAt(later.toISOString().substring(0, 16))

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
          const assRes = await getAssessmentById(id)
          if (assRes.success && assRes.data) {
            const a = assRes.data
            setTitle(a.title)
            setDescription(a.description || '')
            setDurationMinutes(a.durationMinutes || 60)
            setSubjectId(a.subjectId)
            setSelectedClassIds(a.classIds || [])
            setCategory(a.category)
            setStatus(a.status)
            setShowResult(a.showResult)
            if (a.startAt) setStartAt(new Date(a.startAt).toISOString().substring(0, 16))
            if (a.endAt) setEndAt(new Date(a.endAt).toISOString().substring(0, 16))
            if (a.questions && a.questions.length > 0) {
              setQuestions(a.questions)
            }
          } else {
            setFormError('Assessment tidak ditemukan.')
          }
        }
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : 'Gagal memuat form assessment.'
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

  // Question Management
  const addQuestion = (type: 'mcq' | 'essay') => {
    setQuestions((prev) => [
      ...prev,
      {
        type,
        text: '',
        options: type === 'mcq' ? ['', '', '', ''] : undefined,
        correctOption: type === 'mcq' ? 0 : undefined,
        weight: 10,
        order: prev.length + 1,
      },
    ])
  }

  const removeQuestion = (index: number) => {
    if (questions.length <= 1) {
      toast.warning('Minimal Soal', 'Assessment harus memiliki minimal 1 butir soal.')
      return
    }
    setQuestions((prev) => prev.filter((_, i) => i !== index))
  }

  const updateQuestionText = (index: number, text: string) => {
    setQuestions((prev) => {
      const copy = [...prev]
      copy[index] = { ...copy[index], text }
      return copy
    })
  }

  const updateQuestionWeight = (index: number, weight: number) => {
    setQuestions((prev) => {
      const copy = [...prev]
      copy[index] = { ...copy[index], weight: Math.max(1, weight) }
      return copy
    })
  }

  const updateMCQOption = (qIndex: number, optIndex: number, val: string) => {
    setQuestions((prev) => {
      const copy = [...prev]
      const currentOpts = copy[qIndex].options ? [...copy[qIndex].options!] : []
      currentOpts[optIndex] = val
      copy[qIndex] = { ...copy[qIndex], options: currentOpts }
      return copy
    })
  }

  const setCorrectOption = (qIndex: number, optIndex: number) => {
    setQuestions((prev) => {
      const copy = [...prev]
      copy[qIndex] = { ...copy[qIndex], correctOption: optIndex }
      return copy
    })
  }

  const addOptionToMCQ = (qIndex: number) => {
    setQuestions((prev) => {
      const copy = [...prev]
      const currentOpts = copy[qIndex].options || []
      if (currentOpts.length >= 5) {
        toast.warning('Batas Opsi', 'Maksimal 5 opsi jawaban per soal.')
        return prev
      }
      copy[qIndex] = { ...copy[qIndex], options: [...currentOpts, ''] }
      return copy
    })
  }

  const removeOptionFromMCQ = (qIndex: number, optIndex: number) => {
    setQuestions((prev) => {
      const copy = [...prev]
      const currentOpts = copy[qIndex].options || []
      if (currentOpts.length <= 2) {
        toast.warning('Batas Opsi', 'Minimal 2 opsi pilihan ganda.')
        return prev
      }
      const newOpts = currentOpts.filter((_, i) => i !== optIndex)
      let correct = copy[qIndex].correctOption ?? 0
      if (correct >= newOpts.length) correct = newOpts.length - 1
      copy[qIndex] = { ...copy[qIndex], options: newOpts, correctOption: correct }
      return copy
    })
  }

  // Calculate total weight real-time
  const totalWeight = questions.reduce((sum, q) => sum + (Number(q.weight) || 0), 0)

  const validate = () => {
    const errors: Record<string, string> = {}
    if (!title.trim()) errors.title = 'Judul ujian wajib diisi'
    if (!subjectId) errors.subjectId = 'Mata pelajaran wajib dipilih'
    if (selectedClassIds.length === 0) errors.classes = 'Pilih minimal satu kelas tujuan'
    if (durationMinutes < 5) errors.duration = 'Durasi pengerjaan minimal 5 menit'

    // Check questions
    for (let i = 0; i < questions.length; i++) {
      const q = questions[i]
      if (!q.text.trim()) {
        errors.questions = `Pertanyaan pada butir nomor ${i + 1} belum diisi`
        break
      }
      if (q.type === 'mcq') {
        if (!q.options || q.options.some((opt) => !opt.trim())) {
          errors.questions = `Semua opsi jawaban pada butir nomor ${i + 1} harus diisi`
          break
        }
      }
    }

    setValidationErrors(errors)
    return Object.keys(errors).length === 0
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!validate()) return

    setSubmitting(true)
    setFormError(null)

    const payload: CreateAssessmentPayload = {
      title,
      description,
      durationMinutes: Number(durationMinutes),
      subjectId,
      classIds: selectedClassIds,
      category,
      status,
      showResult,
      startAt: startAt ? new Date(startAt).toISOString() : undefined,
      endAt: endAt ? new Date(endAt).toISOString() : undefined,
      questions: questions.map((q, idx) => ({ ...q, order: idx + 1 })),
    }

    try {
      if (isEdit && id) {
        const res = await updateAssessment(id, payload)
        if (res.success) {
          toast.success('Assessment Diperbarui', 'Perubahan assessment dan butir soal berhasil disimpan.')
          navigate(`/teacher/assessments/${id}`)
        } else {
          setFormError(res.error?.message || 'Gagal menyimpan perubahan.')
        }
      } else {
        const res = await createAssessment(payload)
        if (res.success) {
          toast.success('Assessment Dibuat', 'Assessment baru siap diujikan ke siswa.')
          navigate('/teacher/assessments')
        } else {
          setFormError(res.error?.message || 'Gagal membuat assessment baru.')
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
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Back Link */}
        <Link
          to={isEdit && id ? `/teacher/assessments/${id}` : '/teacher/assessments'}
          className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
        >
          <ArrowLeft size={16} />
          <span>Kembali</span>
        </Link>

        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-[#FF9800]">
              <CheckSquare size={20} />
            </div>
            <div>
              <h1 className="text-lg sm:text-xl font-extrabold text-[#30251D]">
                {isEdit ? 'Ubah Assessment & Soal' : 'Buat Assessment & Soal Baru'}
              </h1>
              <p className="text-xs text-[#687385]">
                Susun butir pertanyaan pilihan ganda dan esai dengan penilaian otomatis serta manual.
              </p>
            </div>
          </div>

          {/* Real-time Weight Pill */}
          <div className="rounded-xl border border-[#FF9800]/40 bg-[#FFF8F5] px-4 py-2 text-center sm:text-right">
            <span className="block text-[10px] font-bold uppercase text-[#687385]">
              Total Bobot Nilai
            </span>
            <span
              className={`text-lg font-black ${
                totalWeight === 100 ? 'text-emerald-700' : 'text-[#FF9800]'
              }`}
            >
              {totalWeight} Poin
            </span>
          </div>
        </div>

        {formError && (
          <ErrorState title="Periksa Kembali Form" message={formError} compact />
        )}

        {validationErrors.questions && (
          <div className="rounded-xl border border-rose-200 bg-rose-50 p-3 text-xs text-rose-800 flex items-center gap-2">
            <AlertCircle size={15} className="shrink-0" />
            <span>{validationErrors.questions}</span>
          </div>
        )}

        {pageLoading ? (
          <div className="space-y-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 animate-pulse">
            <div className="h-10 bg-gray-200 rounded-xl" />
            <div className="h-28 bg-gray-200 rounded-xl" />
            <div className="h-10 bg-gray-200 rounded-xl" />
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Section 1: Informasi Dasar */}
            <div className="space-y-5 rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
              <h2 className="text-sm font-extrabold uppercase tracking-wider text-[#30251D] border-b border-[#EBD4C5]/60 pb-3">
                1. Informasi Dasar Assessment
              </h2>

              {/* Judul */}
              <div>
                <label htmlFor="title" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                  Judul Assessment <span className="text-rose-500">*</span>
                </label>
                <input
                  id="title"
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Contoh: Penilaian Harian Bab 2 - Dinamika Gerak"
                  className={`w-full rounded-xl border px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:outline-none focus:ring-1 ${
                    validationErrors.title
                      ? 'border-rose-400 focus:border-rose-500 focus:ring-rose-500'
                      : 'border-[#EBD4C5] focus:border-[#FF9800] focus:ring-[#FF9800]'
                  }`}
                />
              </div>

              {/* Kategori, Status & Durasi */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div>
                  <label htmlFor="category" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                    Kategori Penilaian
                  </label>
                  <select
                    id="category"
                    value={category}
                    onChange={(e) => setCategory(e.target.value as AssessmentCategory)}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2.5 text-xs sm:text-sm font-medium text-[#30251D] focus:border-[#FF9800] focus:outline-none"
                  >
                    <option value="daily">Harian (Daily - Bobot 60%)</option>
                    <option value="uts">UTS (Tengah Semester - 20%)</option>
                    <option value="uas">UAS (Akhir Semester - 20%)</option>
                  </select>
                </div>

                <div>
                  <label htmlFor="status" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                    Status Publikasi
                  </label>
                  <select
                    id="status"
                    value={status}
                    onChange={(e) => setStatus(e.target.value as AssessmentStatus)}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2.5 text-xs sm:text-sm font-medium text-[#30251D] focus:border-[#FF9800] focus:outline-none"
                  >
                    <option value="draft">Draf (Disimpan Sementara)</option>
                    <option value="published">Dipublikasikan (Siap Dikerjakan)</option>
                    <option value="closed">Ditutup</option>
                  </select>
                </div>

                <div>
                  <label htmlFor="durationMinutes" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                    Durasi (Menit) <span className="text-rose-500">*</span>
                  </label>
                  <input
                    id="durationMinutes"
                    type="number"
                    min="5"
                    max="360"
                    value={durationMinutes}
                    onChange={(e) => setDurationMinutes(Number(e.target.value))}
                    className="w-full rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2.5 text-xs sm:text-sm font-bold text-[#30251D] focus:border-[#FF9800] focus:outline-none"
                  />
                </div>
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
                  className="w-full rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] focus:border-[#FF9800] focus:outline-none"
                >
                  <option value="">-- Pilih Mata Pelajaran --</option>
                  {subjects.map((sub) => (
                    <option key={sub.id} value={sub.id}>
                      {sub.name} ({sub.code})
                    </option>
                  ))}
                </select>
              </div>

              {/* Kelas Dituju */}
              <div>
                <label className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-2">
                  Kelas Sasaran <span className="text-rose-500">*</span>
                </label>
                {availableClasses.length === 0 ? (
                  <p className="text-xs text-gray-500 italic p-3 bg-[#FFF8F5] rounded-xl border border-[#EBD4C5]">
                    Pilih mata pelajaran untuk melihat kelas yang diajar.
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
              </div>

              {/* Jadwal Akses & Opsi Hasil */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-2">
                <div>
                  <label htmlFor="startAt" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                    Waktu Mulai Bisa Diakses
                  </label>
                  <input
                    id="startAt"
                    type="datetime-local"
                    value={startAt}
                    onChange={(e) => setStartAt(e.target.value)}
                    className="w-full rounded-xl border border-[#EBD4C5] py-2.5 px-3.5 text-xs sm:text-sm text-[#30251D] bg-white focus:border-[#FF9800] focus:outline-none"
                  />
                </div>

                <div>
                  <label htmlFor="endAt" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                    Waktu Berakhir (Selesai Ujian)
                  </label>
                  <input
                    id="endAt"
                    type="datetime-local"
                    value={endAt}
                    onChange={(e) => setEndAt(e.target.value)}
                    className="w-full rounded-xl border border-[#EBD4C5] py-2.5 px-3.5 text-xs sm:text-sm text-[#30251D] bg-white focus:border-[#FF9800] focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label className="flex items-center gap-3 p-3 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={showResult}
                    onChange={(e) => setShowResult(e.target.checked)}
                    className="h-4 w-4 rounded text-[#FF9800] focus:ring-[#FF9800]"
                  />
                  <span className="text-xs font-bold text-[#30251D]">
                    Tampilkan hasil nilai ke siswa segera setelah selesai dikerjakan / dinilai
                  </span>
                </label>
              </div>
            </div>

            {/* Section 2: Butir Soal (MCQ & Essay Builder) */}
            <div className="space-y-5 rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-[#EBD4C5]/60 pb-3">
                <div>
                  <h2 className="text-sm font-extrabold uppercase tracking-wider text-[#30251D]">
                    2. Butir Soal ({questions.length} Butir)
                  </h2>
                  <p className="text-xs text-[#687385]">
                    Pilihan ganda dinilai otomatis. Soal esai akan dinilai oleh guru secara manual.
                  </p>
                </div>

                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    onClick={() => addQuestion('mcq')}
                    className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] px-3 py-1.5 text-xs font-bold text-[#30251D] hover:bg-[#FEEADC] hover:text-[#FF9800] transition-colors"
                  >
                    <Plus size={14} />
                    <span>+ Pilihan Ganda</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => addQuestion('essay')}
                    className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] px-3 py-1.5 text-xs font-bold text-[#30251D] hover:bg-[#FEEADC] hover:text-[#FF9800] transition-colors"
                  >
                    <Plus size={14} />
                    <span>+ Esai</span>
                  </button>
                </div>
              </div>

              {/* Questions List */}
              <div className="space-y-5">
                {questions.map((q, qIdx) => (
                  <div
                    key={qIdx}
                    className="rounded-2xl border border-[#EBD4C5] bg-[#FFF8F5]/30 p-5 space-y-4"
                  >
                    <div className="flex items-center justify-between gap-3">
                      <div className="flex items-center gap-2">
                        <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-[#FF9800] text-xs font-bold text-white shadow-xs">
                          {qIdx + 1}
                        </span>
                        <span
                          className={`rounded px-2 py-0.5 text-[11px] font-bold uppercase ${
                            q.type === 'mcq'
                              ? 'bg-blue-100 text-blue-800'
                              : 'bg-purple-100 text-purple-800'
                          }`}
                        >
                          {q.type === 'mcq' ? 'Pilihan Ganda' : 'Esai'}
                        </span>
                      </div>

                      <div className="flex items-center gap-3">
                        <div className="flex items-center gap-1.5">
                          <label htmlFor={`weight-${qIdx}`} className="text-xs font-bold text-[#687385]">Bobot:</label>
                          <input
                            id={`weight-${qIdx}`}
                            type="number"
                            min="1"
                            max="100"
                            value={q.weight}
                            onChange={(e) => updateQuestionWeight(qIdx, Number(e.target.value))}
                            className="w-16 rounded-lg border border-[#EBD4C5] bg-white px-2 py-1 text-xs font-bold text-center text-[#30251D]"
                          />
                        </div>

                        <button
                          type="button"
                          onClick={() => removeQuestion(qIdx)}
                          className="flex h-8 w-8 items-center justify-center rounded-lg border border-rose-200 bg-white text-rose-600 hover:bg-rose-50 transition-colors"
                          title="Hapus Soal"
                        >
                          <Trash2 size={14} />
                        </button>
                      </div>
                    </div>

                    {/* Teks Pertanyaan */}
                    <div>
                      <label htmlFor={`q-text-${qIdx}`} className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1">
                        Teks Pertanyaan
                      </label>
                      <textarea
                        id={`q-text-${qIdx}`}
                        rows={3}
                        value={q.text}
                        onChange={(e) => updateQuestionText(qIdx, e.target.value)}
                        placeholder="Tuliskan pertanyaan soal..."
                        className="w-full rounded-xl border border-[#EBD4C5] bg-white p-3 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
                      />
                    </div>

                    {/* MCQ Options */}
                    {q.type === 'mcq' && (
                      <div className="space-y-2 pt-2">
                        <div className="flex items-center justify-between">
                          <span className="text-xs font-bold text-[#30251D]">
                            Opsi Pilihan (Pilih lingkaran untuk menetapkan kunci jawaban):
                          </span>
                          {(q.options?.length ?? 0) < 5 && (
                            <button
                              type="button"
                              onClick={() => addOptionToMCQ(qIdx)}
                              className="text-xs font-bold text-[#FF9800] hover:underline"
                            >
                              + Tambah Opsi
                            </button>
                          )}
                        </div>

                        {q.options?.map((opt, optIdx) => {
                          const letter = String.fromCharCode(65 + optIdx)
                          const isCorrect = q.correctOption === optIdx
                          return (
                            <div key={optIdx} className="flex items-center gap-2">
                              <button
                                type="button"
                                onClick={() => setCorrectOption(qIdx, optIdx)}
                                title={isCorrect ? 'Kunci Jawaban' : 'Jadikan Kunci'}
                                className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-extrabold border transition-colors ${
                                  isCorrect
                                    ? 'bg-emerald-600 text-white border-emerald-600 shadow-xs'
                                    : 'bg-white border-[#EBD4C5] text-gray-500 hover:border-[#FF9800]'
                                }`}
                              >
                                {letter}
                              </button>

                              <input
                                type="text"
                                value={opt}
                                onChange={(e) => updateMCQOption(qIdx, optIdx, e.target.value)}
                                placeholder={`Teks opsi ${letter}...`}
                                className={`flex-1 rounded-xl border px-3 py-1.5 text-xs sm:text-sm text-[#30251D] bg-white focus:outline-none focus:ring-1 ${
                                  isCorrect
                                    ? 'border-emerald-300 focus:border-emerald-500 focus:ring-emerald-500 bg-emerald-50/20'
                                    : 'border-[#EBD4C5] focus:border-[#FF9800] focus:ring-[#FF9800]'
                                }`}
                              />

                              {(q.options?.length ?? 0) > 2 && (
                                <button
                                  type="button"
                                  onClick={() => removeOptionFromMCQ(qIdx, optIdx)}
                                  className="text-gray-400 hover:text-rose-600 p-1 rounded"
                                  title="Hapus opsi"
                                >
                                  <Trash2 size={14} />
                                </button>
                              )}
                            </div>
                          )
                        })}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {/* Submit Actions */}
            <div className="flex items-center justify-end gap-3 pt-4 border-t border-[#EBD4C5]">
              <Link
                to="/teacher/assessments"
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
                <span>{isEdit ? 'Simpan Perubahan' : 'Terbitkan Assessment'}</span>
              </button>
            </div>
          </form>
        )}
      </div>
    </TeacherLayout>
  )
}
