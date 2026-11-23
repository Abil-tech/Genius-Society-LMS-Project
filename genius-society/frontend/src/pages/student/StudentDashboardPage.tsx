import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import {
  FileSpreadsheet,
  CheckSquare,
  Clock,
  Award,
  CalendarDays,
  ArrowRight,
} from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import StatCard from '../../components/common/StatCard'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getStudentDashboard } from '../../services/student/studentDashboard.service'
import type { StudentDashboard } from '../../types/student'

export default function StudentDashboardPage() {
  const [data, setData] = useState<StudentDashboard | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getStudentDashboard()
        if (isMounted) setData(res)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat dashboard'
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

  if (error) {
    return (
      <StudentLayout>
        <ErrorState
          title="Terjadi Kesalahan"
          message={error}
          onRetry={() => window.location.reload()}
        />
      </StudentLayout>
    )
  }

  return (
    <StudentLayout>
      <div className="space-y-6">
        {/* Header Greeting */}
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
            Selamat Datang di Portal Belajar
          </h1>
          <p className="mt-1 text-xs text-[#687385]">
            Pantau tugas, ikuti assessment, dan lihat perkembangan nilai akademikmu di sini.
          </p>
        </div>

        {/* Stat Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <StatCard
            title="Tugas Aktif"
            value={loading ? '...' : (data?.activeAssignmentsCount ?? 0)}
            icon={FileSpreadsheet}
            highlight={true}
            description="Tugas yang perlu dikerjakan"
          />
          <StatCard
            title="Assessment Aktif"
            value={loading ? '...' : (data?.activeAssessmentsCount ?? 0)}
            icon={CheckSquare}
            description="Kuis / Ujian yang tersedia"
          />
          <StatCard
            title="Jadwal Hari Ini"
            value={loading ? '...' : (data?.todaySchedule?.length ?? 0)}
            icon={CalendarDays}
            description="Mata pelajaran terjadwal"
          />
        </div>

        {/* Deadlines & Today Schedule */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Deadline Terdekat */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-xs">
            <div className="flex items-center justify-between pb-3 border-b border-[#EBD4C5]/60">
              <div className="flex items-center gap-2">
                <Clock size={16} className="text-[#FF9800]" />
                <h2 className="text-sm font-bold text-[#30251D]">Deadline Terdekat</h2>
              </div>
              <Link
                to="/student/assignments"
                className="inline-flex items-center gap-1 text-xs font-bold text-[#FF9800] hover:underline"
              >
                <span>Lihat Semua</span>
                <ArrowRight size={12} />
              </Link>
            </div>

            <div className="mt-4 space-y-3">
              {loading ? (
                <div className="animate-pulse space-y-2">
                  <div className="h-12 bg-gray-100 rounded-xl" />
                  <div className="h-12 bg-gray-100 rounded-xl" />
                </div>
              ) : !data?.upcomingDeadlines || data.upcomingDeadlines.length === 0 ? (
                <EmptyState
                  title="Tidak ada deadline terdekat"
                  description="Semua tugas dan projekmu telah selesai atau belum ada tugas baru."
                />
              ) : (
                data.upcomingDeadlines.map((item) => (
                  <div
                    key={item.id}
                    className="flex items-center justify-between p-3 rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/60"
                  >
                    <div className="min-w-0 pr-3">
                      <p className="truncate text-xs font-bold text-[#30251D]">
                        {item.title}
                      </p>
                      <p className="text-[11px] text-[#687385] font-semibold mt-0.5">
                        {item.type === 'assignment' ? 'Tugas' : 'Projek'} &bull; Deadline:{' '}
                        {new Date(item.deadline).toLocaleDateString('id-ID', {
                          day: 'numeric',
                          month: 'short',
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </p>
                    </div>
                    <span
                      className={`shrink-0 rounded-lg px-2.5 py-1 text-[10px] font-extrabold ${
                        item.submitted
                          ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                          : 'bg-[#FEEADC] text-[#FF9800] border border-[#FF9800]/30'
                      }`}
                    >
                      {item.submitted ? 'Sudah Submit' : 'Belum Submit'}
                    </span>
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Jadwal Hari Ini */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-xs">
            <div className="flex items-center justify-between pb-3 border-b border-[#EBD4C5]/60">
              <div className="flex items-center gap-2">
                <CalendarDays size={16} className="text-[#FF9800]" />
                <h2 className="text-sm font-bold text-[#30251D]">Jadwal Pelajaran Hari Ini</h2>
              </div>
              <Link
                to="/student/schedule"
                className="inline-flex items-center gap-1 text-xs font-bold text-[#FF9800] hover:underline"
              >
                <span>Jadwal Lengkap</span>
                <ArrowRight size={12} />
              </Link>
            </div>

            <div className="mt-4 space-y-3">
              {loading ? (
                <div className="animate-pulse space-y-2">
                  <div className="h-12 bg-gray-100 rounded-xl" />
                  <div className="h-12 bg-gray-100 rounded-xl" />
                </div>
              ) : !data?.todaySchedule || data.todaySchedule.length === 0 ? (
                <EmptyState
                  title="Tidak ada jadwal hari ini"
                  description="Hari ini tidak ada mata pelajaran terjadwal untuk kelasmu."
                />
              ) : (
                data.todaySchedule.map((sc) => (
                  <div
                    key={sc.id}
                    className="flex items-center justify-between p-3 rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/60"
                  >
                    <div>
                      <p className="text-xs font-bold text-[#30251D]">{sc.subject}</p>
                      <p className="text-[11px] text-[#687385] mt-0.5">
                        Guru: {sc.teacher || '-'} &bull; Ruang: {sc.room || '-'}
                      </p>
                    </div>
                    <span className="shrink-0 rounded-lg bg-white px-2.5 py-1 text-xs font-mono font-bold text-[#FF9800] border border-[#EBD4C5]">
                      {sc.startTime} - {sc.endTime}
                    </span>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>

        {/* Nilai Terbaru */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-xs">
          <div className="flex items-center justify-between pb-3 border-b border-[#EBD4C5]/60">
            <div className="flex items-center gap-2">
              <Award size={16} className="text-[#FF9800]" />
              <h2 className="text-sm font-bold text-[#30251D]">Nilai Terbaru</h2>
            </div>
            <Link
              to="/student/grades"
              className="inline-flex items-center gap-1 text-xs font-bold text-[#FF9800] hover:underline"
            >
              <span>Semua Nilai</span>
              <ArrowRight size={12} />
            </Link>
          </div>

          <div className="mt-4">
            {loading ? (
              <div className="animate-pulse space-y-2">
                <div className="h-12 bg-gray-100 rounded-xl" />
              </div>
            ) : !data?.recentGrades || data.recentGrades.length === 0 ? (
              <EmptyState
                title="Belum ada riwayat nilai"
                description="Nilai tugas atau kuis yang telah dinilai guru akan muncul di sini."
              />
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {data.recentGrades.map((g) => (
                  <div
                    key={g.id}
                    className="flex items-center justify-between p-3.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/40"
                  >
                    <div>
                      <span className="text-[10px] font-bold uppercase tracking-wider text-[#FF9800]">
                        {g.subject}
                      </span>
                      <p className="text-xs font-bold text-[#30251D] mt-0.5 truncate max-w-[170px]">
                        {g.title}
                      </p>
                      <span className="text-[10px] text-[#687385]">
                        {new Date(g.gradedAt).toLocaleDateString('id-ID')}
                      </span>
                    </div>
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#FFFFFF] border border-[#EBD4C5] text-sm font-black text-[#30251D]">
                      {g.score}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </StudentLayout>
  )
}
