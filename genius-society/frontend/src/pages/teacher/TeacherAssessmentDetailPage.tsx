import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  ArrowLeft,
  Edit,
  Trash2,
  Users,
  Award,
  Layers,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import Modal from '../../components/common/Modal'
import StatusBadge from '../../components/common/StatusBadge'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  getAssessmentById,
  deleteAssessment,
  getAssessmentAttempts,
  getAssessmentAttemptDetail,
  gradeEssayAnswer,
} from '../../services/teacherService'
import type {
  Assessment,
  AssessmentAttempt,
} from '../../types/teacher'

export default function TeacherAssessmentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const toast = useToast()

  const [assessment, setAssessment] = useState<Assessment | null>(null)
  const [attempts, setAttempts] = useState<AssessmentAttempt[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<'soal' | 'peserta'>('peserta')

  // Grading Essay Modal
  const [activeAttempt, setActiveAttempt] = useState<AssessmentAttempt | null>(null)
  const [attemptDetail, setAttemptDetail] = useState<AssessmentAttempt | null>(null)
  const [loadingDetail, setLoadingDetail] = useState(false)
  const [essayScores, setEssayScores] = useState<Record<string, { score: number; feedback?: string }>>({})
  const [submittingGrade, setSubmittingGrade] = useState(false)

  // Delete Dialog
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

  const fetchData = async () => {
    if (!id) return
    setLoading(true)
    setError(null)
    try {
      const [assRes, attRes] = await Promise.allSettled([
        getAssessmentById(id),
        getAssessmentAttempts(id),
      ])

      if (assRes.status === 'fulfilled' && assRes.value.success && assRes.value.data) {
        setAssessment(assRes.value.data)
      }
      if (attRes.status === 'fulfilled' && attRes.value.success && attRes.value.data) {
        setAttempts(attRes.value.data)
      } else if (attRes.status === 'fulfilled' && !attRes.value.success) {
        setError(attRes.value.error?.message || 'Gagal memuat peserta assessment.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat data.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [id])

  const handleDelete = async () => {
    if (!id) return
    setIsDeleting(true)
    try {
      const res = await deleteAssessment(id)
      if (res.success) {
        toast.success(
          'Assessment Diproses',
          res.data?.archived ? 'Assessment telah diarsipkan.' : 'Assessment telah dihapus.'
        )
        navigate('/teacher/assessments')
      } else {
        toast.error('Gagal Menghapus', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat menghapus.'
      toast.error('Gagal Menghapus', msg)
    } finally {
      setIsDeleting(false)
      setDeleteOpen(false)
    }
  }

  // Open Essay Grading Modal
  const openEssayGrading = async (attempt: AssessmentAttempt) => {
    setActiveAttempt(attempt)
    setLoadingDetail(true)
    try {
      const res = await getAssessmentAttemptDetail(id!, attempt.id)
      if (res.success && res.data) {
        setAttemptDetail(res.data)
        // Initialize scores
        const initial: Record<string, { score: number; feedback?: string }> = {}
        res.data.answers?.forEach((ans) => {
          if (ans.type === 'essay') {
            initial[ans.questionId] = {
              score: ans.manualScore ?? 0,
              feedback: ans.feedback || '',
            }
          }
        })
        setEssayScores(initial)
      } else {
        toast.error('Gagal Memuat Jawaban', res.error?.message || 'Data tidak ditemukan.')
      }
    } catch (err: unknown) {
      toast.error('Kendala Sistem', 'Gagal memuat jawaban siswa.')
    } finally {
      setLoadingDetail(false)
    }
  }

  const handleSaveEssayGrade = async (questionId: string) => {
    if (!id || !activeAttempt) return
    const entry = essayScores[questionId]
    if (!entry) return

    setSubmittingGrade(true)
    try {
      const res = await gradeEssayAnswer(id, activeAttempt.id, questionId, {
        score: entry.score,
        feedback: entry.feedback,
      })
      if (res.success) {
        toast.success('Nilai Esai Disimpan', `Skor ${entry.score} berhasil diperbarui.`)
        fetchData()
      } else {
        toast.error('Gagal Menyimpan', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch {
      toast.error('Kendala Sistem', 'Gagal menyimpan nilai esai.')
    } finally {
      setSubmittingGrade(false)
    }
  }

  const attemptColumns: Column<AssessmentAttempt>[] = [
    {
      key: 'student',
      header: 'Nama Siswa',
      render: (a) => (
        <div>
          <p className="font-bold text-[#30251D]">{a.studentName}</p>
          {a.studentNis && <p className="text-[11px] text-[#687385]">NIS: {a.studentNis}</p>}
        </div>
      ),
    },
    {
      key: 'status',
      header: 'Status',
      render: (a) => <StatusBadge status={a.status} />,
    },
    {
      key: 'submittedAt',
      header: 'Waktu Pengumpulan',
      render: (a) =>
        a.submittedAt ? (
          <span className="text-xs text-[#687385]">
            {new Date(a.submittedAt).toLocaleString('id-ID', {
              dateStyle: 'medium',
              timeStyle: 'short',
            })}
          </span>
        ) : (
          <span className="text-xs text-amber-700 italic">Sedang Mengerjakan</span>
        ),
    },
    {
      key: 'scores',
      header: 'Rincian Nilai',
      render: (a) => (
        <div className="text-xs">
          <span className="font-bold text-base text-[#30251D]">{a.finalScore ?? '-'}</span>
          <div className="text-[11px] text-[#687385]">
            Auto: {a.autoScore} • Esai: {a.manualScore}
          </div>
        </div>
      ),
    },
    {
      key: 'action',
      header: 'Aksi',
      align: 'right',
      render: (a) => (
        <button
          type="button"
          onClick={() => openEssayGrading(a)}
          className="inline-flex items-center gap-1.5 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] px-3 py-1.5 text-xs font-bold text-[#FF9800] hover:bg-[#FEEADC] transition-colors"
        >
          <Award size={14} />
          <span>Periksa & Nilai Esai</span>
        </button>
      ),
    },
  ]

  const hasAttempts = attempts.length > 0

  return (
    <TeacherLayout>
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Back Link */}
        <div className="flex items-center justify-between">
          <Link
            to="/teacher/assessments"
            className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
          >
            <ArrowLeft size={16} />
            <span>Kembali ke Daftar Assessment</span>
          </Link>

          {assessment && (
            <div className="flex items-center gap-2">
              <Link
                to={`/teacher/assessments/${assessment.id}/edit`}
                className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] hover:bg-[#FFF8F5] shadow-xs transition-colors"
              >
                <Edit size={14} />
                <span>Ubah</span>
              </Link>
              <button
                type="button"
                onClick={() => setDeleteOpen(true)}
                className="inline-flex items-center gap-1.5 rounded-xl border border-rose-200 bg-white px-3.5 py-2 text-xs font-bold text-rose-600 hover:bg-rose-50 shadow-xs transition-colors"
              >
                <Trash2 size={14} />
                <span>Hapus</span>
              </button>
            </div>
          )}
        </div>

        {error ? (
          <ErrorState title="Assessment Tidak Ditemukan" message={error} onRetry={fetchData} />
        ) : loading ? (
          <div className="space-y-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 animate-pulse">
            <div className="h-8 bg-gray-200 rounded-lg w-2/3" />
            <div className="h-20 bg-gray-200 rounded-lg" />
          </div>
        ) : assessment ? (
          <div className="space-y-6">
            {/* Header Card */}
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 sm:p-8 shadow-xs">
              <div className="flex flex-wrap items-center gap-2 mb-3">
                <span className="rounded-full bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-0.5 text-[11px] font-extrabold uppercase text-[#FF9800]">
                  {assessment.subjectName || 'Mata Pelajaran'}
                </span>
                <span className="rounded px-2 py-0.5 text-[10px] font-extrabold uppercase bg-amber-100 text-amber-800">
                  {assessment.category.toUpperCase()}
                </span>
                <StatusBadge status={assessment.status} />
              </div>

              <h1 className="text-xl sm:text-2xl font-extrabold text-[#30251D]">
                {assessment.title}
              </h1>

              {assessment.description && (
                <p className="mt-2 text-xs sm:text-sm text-[#687385]">{assessment.description}</p>
              )}

              {/* Rincian Ringkas */}
              <div className="mt-6 grid grid-cols-2 sm:grid-cols-4 gap-3 pt-5 border-t border-[#EBD4C5]/60 text-center">
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">Durasi</span>
                  <p className="text-base font-extrabold text-[#30251D]">
                    {assessment.durationMinutes} Menit
                  </p>
                </div>
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">Total Soal</span>
                  <p className="text-base font-extrabold text-[#30251D]">
                    {assessment.questions?.length ?? 0} Butir
                  </p>
                </div>
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">Peserta Mengerjakan</span>
                  <p className="text-base font-extrabold text-[#FF9800]">{attempts.length} Siswa</p>
                </div>
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">Kunci & Nilai</span>
                  <p className="text-sm font-bold text-emerald-700 mt-0.5">
                    {assessment.showResult ? 'Diumumkan' : 'Dirahasiakan'}
                  </p>
                </div>
              </div>
            </div>

            {/* Tab Nav: Peserta Ujian vs Naskah Soal */}
            <div className="flex border-b border-[#EBD4C5]">
              <button
                type="button"
                onClick={() => setActiveTab('peserta')}
                className={`flex items-center gap-2 px-5 py-3 text-xs sm:text-sm font-bold border-b-2 transition-all ${
                  activeTab === 'peserta'
                    ? 'border-[#FF9800] text-[#FF9800]'
                    : 'border-transparent text-[#687385] hover:text-[#30251D]'
                }`}
              >
                <Users size={16} />
                <span>Hasil & Peserta ({attempts.length})</span>
              </button>
              <button
                type="button"
                onClick={() => setActiveTab('soal')}
                className={`flex items-center gap-2 px-5 py-3 text-xs sm:text-sm font-bold border-b-2 transition-all ${
                  activeTab === 'soal'
                    ? 'border-[#FF9800] text-[#FF9800]'
                    : 'border-transparent text-[#687385] hover:text-[#30251D]'
                }`}
              >
                <Layers size={16} />
                <span>Naskah Soal ({assessment.questions?.length ?? 0})</span>
              </button>
            </div>

            {/* Tab 1: Peserta */}
            {activeTab === 'peserta' && (
              <div className="space-y-4">
                <DataTable
                  columns={attemptColumns}
                  data={attempts}
                  loading={loading}
                  emptyMessage="Belum ada siswa yang mengerjakan ujian ini"
                  emptySubMessage="Hasil pengerjaan dan penilaian esai siswa akan muncul di tabel ini."
                />
              </div>
            )}

            {/* Tab 2: Naskah Soal */}
            {activeTab === 'soal' && (
              <div className="space-y-4">
                {assessment.questions?.map((q, idx) => (
                  <div
                    key={idx}
                    className="rounded-2xl border border-[#EBD4C5] bg-white p-5 space-y-3"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="flex h-6 w-6 items-center justify-center rounded-lg bg-[#FF9800] text-xs font-bold text-white">
                          {idx + 1}
                        </span>
                        <span
                          className={`rounded px-2 py-0.5 text-[10px] font-bold uppercase ${
                            q.type === 'mcq'
                              ? 'bg-blue-100 text-blue-800'
                              : 'bg-purple-100 text-purple-800'
                          }`}
                        >
                          {q.type === 'mcq' ? 'Pilihan Ganda' : 'Esai'}
                        </span>
                      </div>
                      <span className="text-xs font-bold text-[#687385]">
                        Bobot: {q.weight} Poin
                      </span>
                    </div>

                    <p className="text-xs sm:text-sm font-medium text-[#30251D]">{q.text}</p>

                    {q.type === 'mcq' && q.options && (
                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 pt-2">
                        {q.options.map((opt, oIdx) => {
                          const isCorrect = q.correctOption === oIdx
                          const letter = String.fromCharCode(65 + oIdx)
                          return (
                            <div
                              key={oIdx}
                              className={`flex items-center gap-2.5 rounded-xl border p-2.5 text-xs ${
                                isCorrect
                                  ? 'border-emerald-300 bg-emerald-50 text-emerald-900 font-bold'
                                  : 'border-[#EBD4C5]/80 bg-[#FFF8F5]/40 text-[#30251D]'
                              }`}
                            >
                              <span
                                className={`flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-bold ${
                                  isCorrect
                                    ? 'bg-emerald-600 text-white'
                                    : 'bg-white border border-[#EBD4C5] text-gray-500'
                                }`}
                              >
                                {letter}
                              </span>
                              <span>{opt}</span>
                              {isCorrect && (
                                <span className="ml-auto text-[10px] font-extrabold uppercase text-emerald-700">
                                  Kunci Jawaban
                                </span>
                              )}
                            </div>
                          )
                        })}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        ) : null}

        {/* Modal Penilaian Esai Siswa */}
        <Modal
          isOpen={Boolean(activeAttempt)}
          onClose={() => setActiveAttempt(null)}
          title={`Penilaian Esai Siswa: ${activeAttempt?.studentName}`}
          description="Periksa teks jawaban esai yang dikerjakan siswa dan input skor manual per butir soal."
          size="lg"
        >
          {loadingDetail ? (
            <div className="space-y-4 animate-pulse p-4">
              <div className="h-20 bg-gray-100 rounded-xl" />
              <div className="h-20 bg-gray-100 rounded-xl" />
            </div>
          ) : (
            <div className="space-y-6">
              {attemptDetail?.answers?.filter((a) => a.type === 'essay').length === 0 ? (
                <p className="text-xs text-gray-500 italic p-4 text-center">
                  Tidak ada soal esai pada ujian ini. Semua butir soal telah dinilai otomatis.
                </p>
              ) : (
                attemptDetail?.answers
                  ?.filter((a) => a.type === 'essay')
                  .map((ans, idx) => {
                    const currentEntry = essayScores[ans.questionId] || { score: 0, feedback: '' }
                    return (
                      <div
                        key={ans.questionId}
                        className="rounded-2xl border border-[#EBD4C5] bg-[#FFF8F5]/30 p-4 space-y-3"
                      >
                        <div className="flex items-center justify-between">
                          <span className="text-xs font-bold uppercase tracking-wider text-[#FF9800]">
                            Soal Esai #{idx + 1} (Bobot Maks: {ans.weight ?? 10})
                          </span>
                        </div>

                        {ans.questionText && (
                          <div className="rounded-xl border border-stone-200 bg-white p-3 text-xs text-[#30251D]">
                            <span className="font-bold block mb-1">Pertanyaan:</span>
                            <p>{ans.questionText}</p>
                          </div>
                        )}

                        <div className="rounded-xl border border-[#EBD4C5] bg-white p-3 text-xs">
                          <span className="font-bold text-[#30251D] block mb-1">
                            Jawaban Siswa:
                          </span>
                          <p className="text-[#30251D] whitespace-pre-wrap leading-relaxed">
                            {ans.essayText || <span className="text-gray-400 italic">Tidak menjawab</span>}
                          </p>
                        </div>

                        <div className="flex flex-col sm:flex-row sm:items-center gap-3 pt-2">
                          <div className="flex items-center gap-2">
                            <label className="text-xs font-bold text-[#30251D]">Skor Nilai:</label>
                            <input
                              type="number"
                              min="0"
                              max={ans.weight ?? 100}
                              value={currentEntry.score}
                              onChange={(e) =>
                                setEssayScores((prev) => ({
                                  ...prev,
                                  [ans.questionId]: {
                                    ...prev[ans.questionId],
                                    score: Number(e.target.value),
                                  },
                                }))
                              }
                              className="w-20 rounded-xl border border-[#EBD4C5] bg-white p-2 text-xs font-bold text-center text-[#30251D]"
                            />
                          </div>

                          <div className="flex-1">
                            <input
                              type="text"
                              placeholder="Catatan / Feedback untuk jawaban ini..."
                              value={currentEntry.feedback || ''}
                              onChange={(e) =>
                                setEssayScores((prev) => ({
                                  ...prev,
                                  [ans.questionId]: {
                                    ...prev[ans.questionId],
                                    feedback: e.target.value,
                                  },
                                }))
                              }
                              className="w-full rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs text-[#30251D]"
                            />
                          </div>

                          <button
                            type="button"
                            onClick={() => handleSaveEssayGrade(ans.questionId)}
                            disabled={submittingGrade}
                            className="rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] disabled:opacity-50"
                          >
                            Simpan
                          </button>
                        </div>
                      </div>
                    )
                  })
              )}
            </div>
          )}
        </Modal>

        {/* Delete Dialog */}
        <ConfirmDialog
          isOpen={deleteOpen}
          onClose={() => setDeleteOpen(false)}
          onConfirm={handleDelete}
          title={hasAttempts ? 'Arsipkan Assessment' : 'Hapus Assessment'}
          message={`Apakah Anda yakin ingin menghapus assessment "${assessment?.title}"?`}
          hasSubmissions={hasAttempts}
          archiveNotice="Assessment ini telah diikuti siswa. Data akan diarsipkan agar seluruh riwayat pengerjaan tetap tersimpan aman."
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
