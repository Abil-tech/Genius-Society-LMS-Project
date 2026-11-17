import { useEffect, useState, useMemo } from 'react'
import {
  Search,
  RefreshCw,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import ErrorState from '../../components/common/ErrorState'
import { getWalasGrades } from '../../services/teacherService'
import type { WalasGradeRekapRow } from '../../types/teacher'

export default function TeacherWalasGradesPage() {
  const [rekap, setRekap] = useState<WalasGradeRekapRow[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getWalasGrades()
      if (res.success && res.data) {
        setRekap(res.data)
      } else {
        setError(res.error?.message || 'Gagal memuat rekap nilai kelas perwalian.')
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
  }, [])

  // Extract distinct subjects from all rows
  const allSubjectNames = useMemo(() => {
    const names = new Set<string>()
    rekap.forEach((row) => {
      row.subjectScores?.forEach((s) => names.add(s.subjectName))
    })
    return Array.from(names)
  }, [rekap])

  const filteredRekap = useMemo(() => {
    if (!search.trim()) return rekap
    const q = search.toLowerCase()
    return rekap.filter(
      (r) =>
        r.studentName.toLowerCase().includes(q) ||
        (r.studentNis && r.studentNis.toLowerCase().includes(q))
    )
  }, [rekap, search])

  return (
    <TeacherLayout isWalas={true}>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <div className="flex items-center gap-2">
              <span className="rounded-full bg-[#FF9800]/10 px-2.5 py-0.5 text-[10px] font-extrabold uppercase text-[#FF9800]">
                Tugas Wali Kelas
              </span>
            </div>
            <h1 className="mt-1 text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Rekapitulasi Nilai Seluruh Mata Pelajaran
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Transkrip nilai seluruh mata pelajaran kelas perwalian untuk persiapan rapor dan peringkat siswa (Read-only).
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchData}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Segarkan</span>
            </button>
          </div>
        </div>

        {/* Search */}
        <div className="relative max-w-md">
          <Search
            size={16}
            className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
          />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari siswa dalam rekap rapor..."
            className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
          />
        </div>

        {/* Table Content */}
        {error ? (
          <ErrorState title="Gagal Memuat Rekap Rapor" message={error} onRetry={fetchData} />
        ) : (
          <div className="w-full overflow-hidden rounded-2xl border border-[#EBD4C5] bg-white shadow-xs">
            <div className="w-full overflow-x-auto">
              <table className="w-full border-collapse text-left text-xs sm:text-sm">
                <thead>
                  <tr className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-[#30251D]">
                    <th className="px-4 py-3.5 text-xs font-bold uppercase text-[#687385] text-center w-12">
                      Rank
                    </th>
                    <th className="px-4 py-3.5 text-xs font-bold uppercase text-[#687385]">
                      NIS
                    </th>
                    <th className="px-4 py-3.5 text-xs font-bold uppercase text-[#687385]">
                      Nama Siswa
                    </th>
                    {allSubjectNames.map((sName) => (
                      <th
                        key={sName}
                        className="px-3 py-3.5 text-xs font-bold uppercase text-[#687385] text-center whitespace-nowrap"
                      >
                        {sName}
                      </th>
                    ))}
                    <th className="px-4 py-3.5 text-xs font-bold uppercase text-[#FF9800] text-right">
                      Rata-Rata
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[#EBD4C5]/60">
                  {loading ? (
                    Array.from({ length: 5 }).map((_, rIdx) => (
                      <tr key={`sk-${rIdx}`} className="animate-pulse">
                        <td colSpan={4 + allSubjectNames.length} className="px-4 py-4">
                          <div className="h-4 bg-gray-200 rounded" />
                        </td>
                      </tr>
                    ))
                  ) : filteredRekap.length === 0 ? (
                    <tr>
                      <td
                        colSpan={4 + allSubjectNames.length}
                        className="px-6 py-12 text-center text-[#687385]"
                      >
                        <p className="font-semibold text-sm text-[#30251D]">
                          Belum ada data rekap nilai kelas perwalian
                        </p>
                      </td>
                    </tr>
                  ) : (
                    filteredRekap.map((row, idx) => {
                      const rank = row.rank ?? idx + 1
                      return (
                        <tr key={row.studentId} className="hover:bg-[#FFF8F5]/60 transition-colors">
                          <td className="px-4 py-3.5 text-center">
                            <span
                              className={`inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-black ${
                                rank === 1
                                  ? 'bg-amber-100 text-amber-900 border border-amber-300'
                                  : rank === 2
                                  ? 'bg-stone-200 text-stone-800'
                                  : rank === 3
                                  ? 'bg-orange-100 text-orange-800'
                                  : 'text-gray-500'
                              }`}
                            >
                              {rank}
                            </span>
                          </td>
                          <td className="px-4 py-3.5 font-bold text-[#30251D]">
                            {row.studentNis}
                          </td>
                          <td className="px-4 py-3.5 font-bold text-[#30251D] whitespace-nowrap">
                            {row.studentName}
                          </td>
                          {allSubjectNames.map((sName) => {
                            const sub = row.subjectScores?.find((s) => s.subjectName === sName)
                            const score = sub ? sub.finalScore : null
                            return (
                              <td
                                key={sName}
                                className="px-3 py-3.5 text-center font-semibold text-xs text-[#30251D]"
                              >
                                {score !== null && score !== undefined ? score : '-'}
                              </td>
                            )
                          })}
                          <td className="px-4 py-3.5 text-right font-black text-sm text-[#FF9800]">
                            {row.overallAverage !== undefined && row.overallAverage !== null
                              ? row.overallAverage.toFixed(1)
                              : '-'}
                          </td>
                        </tr>
                      )
                    })
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </TeacherLayout>
  )
}
