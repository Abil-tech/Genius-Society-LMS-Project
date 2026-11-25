import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { ChevronLeft, ChevronRight, Send } from 'lucide-react'
import AssessmentTimer from '../../components/student/AssessmentTimer'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  startAssessment,
  submitAssessment,
} from '../../services/student/studentAssessment.service'
import type { Question, AnswerSubmission } from '../../types/student'

export default function StudentAssessmentWorkPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const toast = useToast()

  const [title, setTitle] = useState('')
  const [durationSeconds, setDurationSeconds] = useState(0)
  const [startedAt, setStartedAt] = useState('')
  const [questions, setQuestions] = useState<Question[]>([])
  const [currentIndex, setCurrentIndex] = useState(0)
  const [answers, setAnswers] = useState<Record<string, string>>({})

  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [confirmOpen, setConfirmOpen] = useState(false)

  useEffect(() => {
    let isMounted = true
    const init = async () => {
      if (!id) return
      try {
        setLoading(true)
        setError(null)
        const res = await startAssessment(id)
        if (!isMounted) return
        setTitle(res.title)
        setDurationSeconds(res.durationMinutes * 60)
        setStartedAt(res.startedAt)
        const qs = res.questions ?? []
        setQuestions(qs)

        const initialAns: Record<string, string> = {}
        qs.forEach((q) => {
          if (q.answer) initialAns[q.id] = q.answer
        })
        setAnswers(initialAns)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memulai ujian'
          setError(msg)
        }
      } finally {
        if (isMounted) setLoading(false)
      }
    }
    init()
    return () => {
      isMounted = false
    }
  }, [id])

  const handleSelectOption = (questionId: string, option: string) => {
    setAnswers((prev) => ({
      ...prev,
      [questionId]: option,
    }))
  }

  const handleTextAnswer = (questionId: string, text: string) => {
    setAnswers((prev) => ({
      ...prev,
      [questionId]: text,
    }))
  }

  const handleSubmitFinal = async () => {
    if (!id) return
    try {
      setSubmitting(true)
      const payload: AnswerSubmission[] = Object.entries(answers).map(
        ([questionId, answer]) => ({
          questionId,
          answer,
        })
      )
      await submitAssessment(id, payload)
      toast.success('Assessment Selesai', 'Semua jawabanmu berhasil dikirim.')
      navigate(`/student/assessments/${id}`)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal mengumpulkan assessment'
      toast.error('Gagal Submit', msg)
    } finally {
      setSubmitting(false)
      setConfirmOpen(false)
    }
  }

  if (error) {
    return (
      <div className="min-h-screen bg-[#FFF8F5] p-6 flex items-center justify-center">
        <div className="max-w-md w-full">
          <ErrorState
            title="Tidak Dapat Mengikuti Ujian"
            message={error}
            onRetry={() => navigate('/student/assessments')}
          />
        </div>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-[#FFF8F5] flex items-center justify-center">
        <div className="text-center space-y-3">
          <div className="h-10 w-10 mx-auto rounded-full border-3 border-[#FF9800] border-t-transparent animate-spin" />
          <p className="text-xs font-bold text-[#30251D]">Menyiapkan Lembar Soal...</p>
        </div>
      </div>
    )
  }

  const currentQ = questions[currentIndex]
  const answeredCount = Object.keys(answers).filter((k) => Boolean(answers[k])).length

  return (
    <div className="min-h-screen bg-[#FFF8F5] text-[#30251D] flex flex-col">
      {/* Topbar Focus Examination */}
      <header className="sticky top-0 z-20 flex h-16 w-full items-center justify-between border-b border-[#EBD4C5] bg-white px-4 sm:px-8 shadow-xs">
        <div>
          <h1 className="text-xs sm:text-sm font-black text-[#30251D] truncate max-w-xs sm:max-w-md">
            {title}
          </h1>
          <span className="text-[11px] text-[#687385] font-semibold">
            {answeredCount} dari {questions.length} soal dijawab
          </span>
        </div>

        {durationSeconds > 0 && startedAt && (
          <AssessmentTimer
            durationSeconds={durationSeconds}
            startedAt={startedAt}
            onExpired={handleSubmitFinal}
          />
        )}
      </header>

      {/* Main Examination View */}
      <div className="flex-1 max-w-4xl w-full mx-auto p-4 sm:p-6 lg:p-8 flex flex-col justify-between">
        {currentQ ? (
          <div className="space-y-6">
            {/* Question Header & Card */}
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs space-y-4">
              <div className="flex items-center justify-between pb-3 border-b border-[#EBD4C5]/60">
                <span className="text-xs font-black uppercase tracking-wider text-[#FF9800]">
                  Soal No. {currentIndex + 1} dari {questions.length}
                </span>
                <span className="text-[11px] font-bold text-[#687385]">
                  Bobot: {currentQ.weight} Poin
                </span>
              </div>

              <p className="text-sm sm:text-base font-bold text-[#30251D] leading-relaxed">
                {currentQ.text}
              </p>

              {/* Multiple Choice Options */}
              {currentQ.type === 'mcq' && currentQ.options && (
                <div className="mt-4 space-y-2.5">
                  {currentQ.options.map((opt, optIdx) => {
                    const isSelected = answers[currentQ.id] === opt || answers[currentQ.id] === String(optIdx)
                    return (
                      <button
                        key={optIdx}
                        type="button"
                        onClick={() => handleSelectOption(currentQ.id, opt)}
                        className={`w-full flex items-center gap-3 p-3.5 rounded-xl border text-left transition-all text-xs sm:text-sm font-semibold ${
                          isSelected
                            ? 'border-[#FF9800] bg-[#FEEADC] text-[#FF9800]'
                            : 'border-[#EBD4C5] bg-white text-[#30251D] hover:bg-[#FFF8F5]'
                        }`}
                      >
                        <span
                          className={`flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs font-bold ${
                            isSelected
                              ? 'bg-[#FF9800] text-white'
                              : 'bg-[#FFF8F5] border border-[#EBD4C5] text-[#30251D]'
                          }`}
                        >
                          {String.fromCharCode(65 + optIdx)}
                        </span>
                        <span>{opt}</span>
                      </button>
                    )
                  })}
                </div>
              )}

              {/* Essay TextArea */}
              {currentQ.type === 'essay' && (
                <div className="mt-4">
                  <textarea
                    rows={6}
                    value={answers[currentQ.id] || ''}
                    onChange={(e) => handleTextAnswer(currentQ.id, e.target.value)}
                    placeholder="Tuliskan jawaban essay kamu secara lengkap dan jelas..."
                    className="w-full rounded-xl border border-[#EBD4C5] p-3.5 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none"
                  />
                </div>
              )}
            </div>

            {/* Questions Number Palette */}
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-4 shadow-xs">
              <span className="block text-[11px] font-bold uppercase tracking-wider text-[#687385] mb-2.5">
                Navigasi Nomor Soal:
              </span>
              <div className="flex flex-wrap gap-2">
                {questions.map((q, qIdx) => {
                  const isCurrent = currentIndex === qIdx
                  const isAnswered = Boolean(answers[q.id])
                  return (
                    <button
                      key={q.id}
                      type="button"
                      onClick={() => setCurrentIndex(qIdx)}
                      className={`h-9 w-9 rounded-xl text-xs font-bold transition-all ${
                        isCurrent
                          ? 'border-2 border-[#FF9800] bg-white text-[#FF9800]'
                          : isAnswered
                          ? 'bg-[#FF9800] text-white shadow-xs'
                          : 'bg-[#FFF8F5] border border-[#EBD4C5] text-[#30251D] hover:bg-gray-100'
                      }`}
                    >
                      {qIdx + 1}
                    </button>
                  )
                })}
              </div>
            </div>
          </div>
        ) : null}

        {/* Action Controls */}
        <div className="mt-6 pt-4 border-t border-[#EBD4C5] flex items-center justify-between">
          <button
            type="button"
            onClick={() => setCurrentIndex((prev) => Math.max(0, prev - 1))}
            disabled={currentIndex === 0}
            className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-4 py-2.5 text-xs font-bold text-[#30251D] hover:bg-[#FFF8F5] disabled:opacity-40 transition-colors"
          >
            <ChevronLeft size={16} />
            <span>Sebelumnya</span>
          </button>

          {currentIndex < questions.length - 1 ? (
            <button
              type="button"
              onClick={() => setCurrentIndex((prev) => Math.min(questions.length - 1, prev + 1))}
              className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-5 py-2.5 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
            >
              <span>Selanjutnya</span>
              <ChevronRight size={16} />
            </button>
          ) : (
            <button
              type="button"
              onClick={() => setConfirmOpen(true)}
              disabled={submitting}
              className="inline-flex items-center gap-1.5 rounded-xl bg-emerald-600 px-6 py-2.5 text-xs font-bold text-white shadow-xs hover:bg-emerald-700 transition-colors disabled:opacity-50"
            >
              <Send size={15} />
              <span>Selesaikan & Kumpulkan</span>
            </button>
          )}
        </div>
      </div>

      {/* Confirmation Dialog */}
      <ConfirmDialog
        isOpen={confirmOpen}
        title="Kumpulkan Semua Jawaban?"
        message={`Kamu telah menjawab ${answeredCount} dari ${questions.length} butir soal. Apakah kamu yakin ingin menyelesaikan assessment ini?`}
        confirmText={submitting ? 'Mengirim...' : 'Ya, Kumpulkan Sekarang'}
        cancelText="Periksa Lagi"
        variant="primary"
        onConfirm={handleSubmitFinal}
        onClose={() => setConfirmOpen(false)}
      />
    </div>
  )
}
