import { useState, useEffect } from 'react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getGrades } from '../../services/student/studentGrade.service'
import type { Grade } from '../../types/student'

export default function StudentGradesPage() {
  const [grades, setGrades] = useState<Grade[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getGrades()
        if (isMounted) setGrades(Array.isArray(res) ? res : [])
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat rekap nilai'
          setError(msg)
          setGrades([])
        }
      } finally {
        if (isMounted) setLoading(false)
      }
    }
    fetchData()
    return () => {
      isMounted = false
    }
  }, [])

  const getScoreColor = (score?: number) => {
    if (score === undefined || score === null) return 'text-gray-400'
    if (score >= 75) return 'text-emerald-700 font-extrabold'
    if (score >= 60) return 'text-amber-600 font-extrabold'
    return 'text-rose-600 font-extrabold'
  }

  return (
    <StudentLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
            Buku Nilai & Transkrip Akademik
          </h1>
          <p className="mt-1 text-xs text-[#687385]">
            Rekap perolehan nilai mata pelajaran berdasarkan bobot: (Tugas + Assessment Harian) 60% + UTS 20% + UAS 20%.
          </p>
        </div>

        {error ? (
          <ErrorState
            title="Gagal Memuat Nilai"
            message={error}
            onRetry={() => window.location.reload()}
          />
        ) : loading ? (
          <div className="space-y-3">
            <div className="h-44 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
          </div>
        ) : !grades || grades.length === 0 ? (
          <EmptyState
            title="Belum ada data nilai"
            description="Nilai akademikmu akan muncul di sini setelah gurumu memberikan penilaian tugas atau ujian."
          />
        ) : (
          <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-xs overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="bg-[#FFF8F5] border-b border-[#EBD4C5] text-[#30251D] font-extrabold uppercase tracking-wider text-[11px]">
                  <tr>
                    <th className="py-3.5 px-4">Mata Pelajaran</th>
                    <th className="py-3.5 px-3 text-center">Tugas</th>
                    <th className="py-3.5 px-3 text-center">Assessment</th>
                    <th className="py-3.5 px-3 text-center">Projek</th>
                    <th className="py-3.5 px-3 text-center">UTS (20%)</th>
                    <th className="py-3.5 px-3 text-center">UAS (20%)</th>
                    <th className="py-3.5 px-4 text-center bg-[#FEEADC]/40">Nilai Akhir</th>
                    <th className="py-3.5 px-3 text-center">Predikat</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[#EBD4C5]/60">
                  {grades.map((g) => (
                    <tr key={g.subjectId} className="hover:bg-[#FFF8F5]/40 transition-colors">
                      <td className="py-3.5 px-4">
                        <span className="font-bold text-[#30251D] block">
                          {g.subjectName}
                        </span>
                        <span className="text-[10px] text-[#687385] block">
                          Kode: {g.subjectCode}
                        </span>
                      </td>
                      <td className="py-3.5 px-3 text-center font-mono">
                        {g.assignmentAverage !== undefined ? g.assignmentAverage.toFixed(1) : '-'}
                      </td>
                      <td className="py-3.5 px-3 text-center font-mono">
                        {g.assessmentAverage !== undefined ? g.assessmentAverage.toFixed(1) : '-'}
                      </td>
                      <td className="py-3.5 px-3 text-center font-mono">
                        {g.projectAverage !== undefined ? g.projectAverage.toFixed(1) : '-'}
                      </td>
                      <td className="py-3.5 px-3 text-center font-mono">
                        {g.utsScore !== undefined ? g.utsScore.toFixed(1) : '-'}
                      </td>
                      <td className="py-3.5 px-3 text-center font-mono">
                        {g.uasScore !== undefined ? g.uasScore.toFixed(1) : '-'}
                      </td>
                      <td className={`py-3.5 px-4 text-center text-sm font-mono bg-[#FEEADC]/20 ${getScoreColor(g.finalScore)}`}>
                        {g.finalScore !== undefined ? g.finalScore.toFixed(1) : '-'}
                      </td>
                      <td className="py-3.5 px-3 text-center">
                        <span className="inline-block rounded-md bg-[#FFF8F5] border border-[#EBD4C5] px-2 py-0.5 text-xs font-black text-[#30251D]">
                          {g.gradeLetter}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </StudentLayout>
  )
}
