import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import ErrorState from '../../components/common/ErrorState'
import {
  getAssessmentDetail,
  getAssessmentResult,
} from '../../services/student/studentAssessment.service'
import type { Assessment, AssessmentResultResponse } from '../../types/student'

export default function StudentAssessmentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  const [assessment, setAssessment] = useState<Assessment | null>(null)
  const [result, setResult] = useState<AssessmentResultResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      if (!id) return
      try {
        setLoading(true)
        setError(null)
        const res = await getAssessmentDetail(id)
        if (!isMounted) return
        setAssessment(res)

        if (res.attemptStatus === 'submitted' || res.attemptStatus === 'graded') {
          try {
            const resResult = await getAssessmentResult(id)
            if (isMounted) setResult(resResult)
          } catch {
            // result not open yet
          }
        }
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat assessment'
          setError(msg)
        }
      } finally {
        if (isMounted) setLoading(false)
      }
    }
    fetchData()
    return () => {
      isMounted = false
    }
  }, [id])

  if (error) {
    return (
      <StudentLayout>
        <ErrorState
          title="Assessment Tidak Ditemukan"
          message={error}
          onRetry={() => window.location.reload()}
        />
      </StudentLayout>
    )
  }

  if (loading || !assessment) {
    return (
      <StudentLayout>
        <div className="animate-pulse space-y-4 max-w-3xl">
          <div className="h-6 w-32 bg-gray-200 rounded-lg" />
          <div className="h-44 bg-white border border-[#EBD4C5] rounded-2xl" />
        </div>
      </StudentLayout>
    )
  }

  const isCompleted =
    assessment.attemptStatus === 'submitted' || assessment.attemptStatus === 'graded'
  const inProgress = assessment.attemptStatus === 'in_progress'

  return (
    <StudentLayout>
      <div className="space-y-6 max-w-3xl">
        <Link
          to="/student/assessments"
          className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
        >
          <ArrowLeft size={16} />
          <span>Kembali ke Daftar Assessment</span>
        </Link>

        {/* Card Info Assessment */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
          <div className="flex flex-wrap items-center gap-2 text-xs">
            <span className="rounded-lg bg-[#FEEADC] px-2.5 py-1 font-extrabold text-[#FF9800]">
              {assessment.subjectName}
            </span>
            <span className="text-[#687385]">&bull; Guru: {assessment.teacherName || '-'}</span>
          </div>

          <h1 className="mt-3 text-lg sm:text-2xl font-black text-[#30251D]">
            {assessment.title}
          </h1>

          <div className="mt-3 text-xs sm:text-sm text-[#687385] leading-relaxed">
            {assessment.description || 'Tidak ada instruksi khusus untuk assessment ini.'}
          </div>

          {/* Details specs */}
          <div className="mt-6 grid grid-cols-2 sm:grid-cols-3 gap-3 border-t border-[#EBD4C5]/60 pt-5">
            <div className="p-3 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5]/60">
              <span className="text-[10px] font-bold text-[#687385] uppercase tracking-wider block">
                Durasi Waktu
              </span>
              <span className="text-sm font-extrabold text-[#30251D] mt-0.5 block">
                {assessment.durationMinutes} Menit
              </span>
            </div>
            <div className="p-3 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5]/60">
              <span className="text-[10px] font-bold text-[#687385] uppercase tracking-wider block">
                Jumlah Soal
              </span>
              <span className="text-sm font-extrabold text-[#30251D] mt-0.5 block">
                {assessment.questionCount || 0} Butir Soal
              </span>
            </div>
            <div className="p-3 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5]/60 col-span-2 sm:col-span-1">
              <span className="text-[10px] font-bold text-[#687385] uppercase tracking-wider block">
                Status Pengerjaan
              </span>
              <span className="text-sm font-extrabold text-[#FF9800] mt-0.5 block capitalize">
                {assessment.attemptStatus === 'not_started'
                  ? 'Belum Dimulai'
                  : assessment.attemptStatus === 'in_progress'
                  ? 'Sedang Berjalan'
                  : 'Selesai'}
              </span>
            </div>
          </div>

          {/* Action Button */}
          <div className="mt-6 pt-5 border-t border-[#EBD4C5]/60 flex items-center justify-between">
            {isCompleted ? (
              <div className="flex items-center gap-3">
                <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-700 font-black text-xl">
                  {assessment.totalScore ?? result?.totalScore ?? '-'}
                </div>
                <div>
                  <span className="text-xs font-bold text-[#30251D] block">
                    Assessment Selesai
                  </span>
                  <span className="text-[11px] text-[#687385] block">
                    Skor akhir telah tercatat di buku nilai akademikmu.
                  </span>
                </div>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => navigate(`/student/assessments/${id}/work`)}
                className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-6 py-2.5 text-xs sm:text-sm font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
              >
                <span>{inProgress ? 'Lanjutkan Pengerjaan' : 'Mulai Assessment Sekarang'}</span>
                <ArrowRight size={15} />
              </button>
            )}
          </div>
        </div>

        {/* Answers Review if Available */}
        {result?.answers && result.answers.length > 0 && (
          <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs space-y-4">
            <h2 className="text-sm font-bold text-[#30251D] pb-3 border-b border-[#EBD4C5]/60">
              Pembahasan Jawaban
            </h2>
            <div className="space-y-4">
              {result.answers.map((ans, idx) => (
                <div
                  key={ans.questionId || idx}
                  className="p-4 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40 space-y-2"
                >
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-extrabold text-[#30251D]">
                      Soal No. {idx + 1} ({ans.type.toUpperCase()})
                    </span>
                    <span
                      className={`text-xs font-extrabold ${
                        ans.isCorrect ? 'text-emerald-700' : 'text-rose-600'
                      }`}
                    >
                      Poin: {ans.score} / {ans.maxWeight}
                    </span>
                  </div>
                  <p className="text-xs text-[#30251D] font-medium">{ans.text}</p>
                  <p className="text-[11px] text-[#687385]">
                    Jawabanmu: <span className="font-bold text-[#30251D]">{ans.studentAnswer || '(Kosong)'}</span>
                  </p>
                  {ans.feedback && (
                    <p className="text-[11px] text-blue-700 italic">
                      Umpan balik guru: {ans.feedback}
                    </p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </StudentLayout>
  )
}
