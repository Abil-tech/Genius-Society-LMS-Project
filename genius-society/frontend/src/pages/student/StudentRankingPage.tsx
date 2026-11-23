import { useState, useEffect } from 'react'
import { Trophy, Medal, Award } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getRanking } from '../../services/student/studentRanking.service'
import type { ClassRankingResponse } from '../../types/student'

export default function StudentRankingPage() {
  const [data, setData] = useState<ClassRankingResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getRanking()
        if (isMounted) setData(res)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat ranking kelas'
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
  }, [])

  const getRankBadge = (rank: number) => {
    if (rank === 1) {
      return (
        <span className="flex h-7 w-7 items-center justify-center rounded-xl bg-amber-100 text-amber-700 font-black shadow-xs">
          <Trophy size={16} />
        </span>
      )
    }
    if (rank === 2) {
      return (
        <span className="flex h-7 w-7 items-center justify-center rounded-xl bg-slate-200 text-slate-700 font-black shadow-xs">
          <Medal size={16} />
        </span>
      )
    }
    if (rank === 3) {
      return (
        <span className="flex h-7 w-7 items-center justify-center rounded-xl bg-orange-100 text-orange-800 font-black shadow-xs">
          <Award size={16} />
        </span>
      )
    }
    return (
      <span className="flex h-7 w-7 items-center justify-center rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-xs font-bold text-[#30251D]">
        {rank}
      </span>
    )
  }

  return (
    <StudentLayout>
      <div className="space-y-6 max-w-4xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
              Peringkat & Ranking Kelas
            </h1>
            <p className="mt-1 text-xs text-[#687385]">
              Klasemen peringkat siswa di kelas {data?.className || ''} berdasarkan akumulasi nilai akhir.
            </p>
          </div>

          {data?.myRank ? (
            <div className="flex items-center gap-3 rounded-2xl border border-[#FF9800]/40 bg-[#FEEADC] px-4 py-2.5 shadow-xs shrink-0">
              <Trophy size={20} className="text-[#FF9800]" />
              <div>
                <span className="text-[10px] font-bold uppercase tracking-wider text-[#FF9800] block">
                  Peringkat Kamu
                </span>
                <span className="text-base font-black text-[#30251D] block">
                  Juara #{data.myRank} dari {data.totalStudents} Siswa
                </span>
              </div>
            </div>
          ) : null}
        </div>

        {error ? (
          <ErrorState
            title="Gagal Memuat Ranking"
            message={error}
            onRetry={() => window.location.reload()}
          />
        ) : loading ? (
          <div className="space-y-3">
            <div className="h-44 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
          </div>
        ) : !data?.rankings || data.rankings.length === 0 ? (
          <EmptyState
            title="Belum ada data ranking"
            description="Ranking kelas akan terbentuk secara otomatis setelah nilai akhir dihitung."
          />
        ) : (
          <div className="rounded-2xl border border-[#EBD4C5] bg-white shadow-xs overflow-hidden">
            <table className="w-full text-left text-xs">
              <thead className="bg-[#FFF8F5] border-b border-[#EBD4C5] text-[#30251D] font-extrabold uppercase tracking-wider text-[11px]">
                <tr>
                  <th className="py-3.5 px-4 w-20 text-center">Posisi</th>
                  <th className="py-3.5 px-4">Nama Siswa</th>
                  <th className="py-3.5 px-4 text-center">NIS</th>
                  <th className="py-3.5 px-4 text-right">Nilai Akhir</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#EBD4C5]/60">
                {data.rankings.map((r) => (
                  <tr
                    key={r.studentId}
                    className={`transition-colors ${
                      r.isSelf
                        ? 'bg-[#FEEADC] font-bold text-[#30251D]'
                        : 'hover:bg-[#FFF8F5]/40 text-[#30251D]'
                    }`}
                  >
                    <td className="py-3.5 px-4 text-center">
                      <div className="flex justify-center">{getRankBadge(r.rank)}</div>
                    </td>
                    <td className="py-3.5 px-4">
                      <span className="font-bold text-xs sm:text-sm block">
                        {r.name} {r.isSelf && '(Kamu)'}
                      </span>
                    </td>
                    <td className="py-3.5 px-4 text-center font-mono text-[11px] text-[#687385]">
                      {r.nis}
                    </td>
                    <td className="py-3.5 px-4 text-right font-mono font-black text-sm text-[#FF9800]">
                      {r.finalScore.toFixed(1)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </StudentLayout>
  )
}
