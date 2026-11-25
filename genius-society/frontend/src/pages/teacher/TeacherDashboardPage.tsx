import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import {
  DoorOpen,
  Users,
  AlertCircle,
  FileSpreadsheet,
  CheckSquare,
  FolderGit2,
  Calendar,
  Clock,
  RefreshCw,
  ShieldCheck,
  ChevronRight,
  BookOpen,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import StatCard from '../../components/common/StatCard'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import StatusBadge from '../../components/common/StatusBadge'
import { getTeacherDashboard } from '../../services/teacherService'
import type { TeacherDashboardData } from '../../types/teacher'
import { useAuth } from '../../contexts/AuthContext'

export default function TeacherDashboardPage() {
  const { user } = useAuth()
  const [data, setData] = useState<TeacherDashboardData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchDashboard = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getTeacherDashboard()
      if (res.success && res.data) {
        setData(res.data)
      } else {
        // Fallback structure if response body format differs
        setError(res.error?.message || 'Gagal memuat ringkasan dashboard guru.')
      }
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'Terjadi kesalahan saat menghubungi server.'
      setError(message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchDashboard()
  }, [fetchDashboard])

  // Is the teacher a Walas?
  const isWalas = Boolean(
    data?.isWalas ||
      user?.isWalas ||
      user?.is_walas
  )

  const kpi = data?.kpi || {
    classesCount: 0,
    studentsCount: 0,
    pendingGradesCount: 0,
    activeAssignmentsCount: 0,
    activeAssessmentsCount: 0,
    activeProjectsCount: 0,
  }

  return (
    <TeacherLayout isWalas={isWalas}>
      <div className="space-y-6 sm:space-y-8">
        {/* Header Greeting & Action Bar */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Selamat Datang, {user?.name || 'Bapak/Ibu Guru'}!
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Kelola materi pembelajaran, evaluasi tugas & ujian siswa, serta pantau jadwal mengajar.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchDashboard}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5] disabled:opacity-50 transition-colors"
              title="Perbarui data"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Segarkan</span>
            </button>
          </div>
        </div>

        {/* Error Notification Banner if any */}
        {error && (
          <ErrorState
            title="Gagal Memuat Data Dashboard"
            message={error}
            onRetry={fetchDashboard}
            compact
          />
        )}

        {/* Quick Action Bar */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-4 sm:p-5 shadow-xs">
          <p className="text-xs font-bold uppercase tracking-wider text-[#687385] mb-3">
            Aksi Cepat
          </p>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-2.5 sm:gap-3">
            <Link
              to="/teacher/materials/new"
              className="flex items-center gap-2.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] hover:border-[#FF9800] hover:bg-[#FEEADC]/60 transition-all group"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800] group-hover:bg-[#FF9800] group-hover:text-white transition-colors">
                <BookOpen size={14} />
              </div>
              <span className="truncate">Buat Materi</span>
            </Link>

            <Link
              to="/teacher/assignments/new"
              className="flex items-center gap-2.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] hover:border-[#FF9800] hover:bg-[#FEEADC]/60 transition-all group"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800] group-hover:bg-[#FF9800] group-hover:text-white transition-colors">
                <FileSpreadsheet size={14} />
              </div>
              <span className="truncate">Buat Tugas</span>
            </Link>

            <Link
              to="/teacher/assessments/new"
              className="flex items-center gap-2.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] hover:border-[#FF9800] hover:bg-[#FEEADC]/60 transition-all group"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800] group-hover:bg-[#FF9800] group-hover:text-white transition-colors">
                <CheckSquare size={14} />
              </div>
              <span className="truncate">Buat Assessment</span>
            </Link>

            <Link
              to="/teacher/projects/new"
              className="flex items-center gap-2.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] hover:border-[#FF9800] hover:bg-[#FEEADC]/60 transition-all group"
            >
              <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800] group-hover:bg-[#FF9800] group-hover:text-white transition-colors">
                <FolderGit2 size={14} />
              </div>
              <span className="truncate">Buat Projek</span>
            </Link>
          </div>
        </div>

        {/* 6 KPI Cards */}
        <div className="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-3 sm:gap-4">
          <StatCard
            title="Kelas Diajar"
            value={kpi.classesCount}
            icon={DoorOpen}
            linkTo="/teacher/classes"
            loading={loading}
          />
          <StatCard
            title="Total Siswa"
            value={kpi.studentsCount}
            icon={Users}
            linkTo="/teacher/classes"
            loading={loading}
          />
          <StatCard
            title="Belum Dinilai"
            value={kpi.pendingGradesCount}
            icon={AlertCircle}
            highlight={kpi.pendingGradesCount > 0}
            badgeText={kpi.pendingGradesCount > 0 ? 'Perlu Dinilai' : undefined}
            linkTo="/teacher/grades"
            loading={loading}
          />
          <StatCard
            title="Tugas Aktif"
            value={kpi.activeAssignmentsCount}
            icon={FileSpreadsheet}
            linkTo="/teacher/assignments"
            loading={loading}
          />
          <StatCard
            title="Assessment Aktif"
            value={kpi.activeAssessmentsCount}
            icon={CheckSquare}
            linkTo="/teacher/assessments"
            loading={loading}
          />
          <StatCard
            title="Projek Aktif"
            value={kpi.activeProjectsCount}
            icon={FolderGit2}
            linkTo="/teacher/projects"
            loading={loading}
          />
        </div>

        {/* Kartu Khusus Wali Kelas (HANYA jika isWalas) */}
        {isWalas && (
          <div className="overflow-hidden rounded-2xl border border-[#FF9800]/40 bg-gradient-to-r from-[#FFF8F5] via-white to-[#FEEADC] p-5 shadow-xs">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div className="flex items-start gap-3.5">
                <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#FF9800] text-white shadow-xs">
                  <ShieldCheck size={24} />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <span className="rounded-full bg-[#FF9800]/10 px-2 py-0.5 text-[10px] font-extrabold uppercase tracking-wider text-[#FF9800]">
                      Tugas Tambahan
                    </span>
                    <h3 className="text-sm sm:text-base font-bold text-[#30251D]">
                      Wali Kelas {data?.walasClass?.name || 'Perwalian'}
                    </h3>
                  </div>
                  <p className="mt-1 text-xs text-[#687385]">
                    Anda bertindak sebagai Wali Kelas untuk {data?.walasClass?.studentCount ?? 0} siswa.
                    Pantau kehadiran, perkembangan karakter, dan rekap nilai seluruh mata pelajaran.
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-2 sm:shrink-0">
                <Link
                  to="/teacher/walas/students"
                  className="rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] hover:bg-[#FFF8F5] shadow-xs transition-colors"
                >
                  Daftar Siswa
                </Link>
                <Link
                  to="/teacher/walas/grades"
                  className="inline-flex items-center gap-1 rounded-xl bg-[#FF9800] px-3.5 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors"
                >
                  <span>Rekap Nilai</span>
                  <ChevronRight size={14} />
                </Link>
              </div>
            </div>
          </div>
        )}

        {/* 2-Column Section: Jadwal Hari Ini & Pengumpulan Terbaru */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Jadwal Hari Ini */}
          <div className="flex flex-col rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-xs">
            <div className="flex items-center justify-between border-b border-[#EBD4C5]/60 pb-3">
              <div className="flex items-center gap-2">
                <Calendar size={18} className="text-[#FF9800]" />
                <h2 className="text-sm sm:text-base font-bold text-[#30251D]">
                  Jadwal Mengajar Hari Ini
                </h2>
              </div>
              <Link
                to="/teacher/schedule"
                className="text-xs font-bold text-[#FF9800] hover:underline"
              >
                Lihat Semua Jadwal &rarr;
              </Link>
            </div>

            <div className="mt-4 flex-1">
              {loading ? (
                <div className="space-y-3">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="h-14 animate-pulse rounded-xl bg-gray-100" />
                  ))}
                </div>
              ) : !data?.todaySchedule || data.todaySchedule.length === 0 ? (
                <EmptyState
                  icon={Calendar}
                  title="Tidak Ada Jadwal Hari Ini"
                  description="Tidak ada agenda sesi mengajar untuk hari ini. Waktu luang dapat digunakan untuk persiapan materi atau evaluasi."
                />
              ) : (
                <div className="space-y-2.5">
                  {data.todaySchedule.map((item) => (
                    <div
                      key={item.id}
                      className="flex items-center justify-between rounded-xl border border-[#EBD4C5]/70 bg-[#FFF8F5]/60 p-3.5 transition-colors hover:bg-[#FFF8F5]"
                    >
                      <div className="flex items-center gap-3">
                        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800] font-bold text-xs">
                          <Clock size={16} />
                        </div>
                        <div>
                          <p className="text-xs font-bold text-[#30251D]">{item.subject}</p>
                          <p className="text-[11px] text-[#687385]">
                            {item.className} {item.room ? `• Ruang ${item.room}` : ''}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <span className="rounded-lg bg-white border border-[#EBD4C5] px-2.5 py-1 text-xs font-bold text-[#30251D]">
                          {item.startTime} - {item.endTime}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Pengumpulan Terbaru */}
          <div className="flex flex-col rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-xs">
            <div className="flex items-center justify-between border-b border-[#EBD4C5]/60 pb-3">
              <div className="flex items-center gap-2">
                <FileSpreadsheet size={18} className="text-[#FF9800]" />
                <h2 className="text-sm sm:text-base font-bold text-[#30251D]">
                  Pengumpulan Terbaru
                </h2>
              </div>
              <Link
                to="/teacher/grades"
                className="text-xs font-bold text-[#FF9800] hover:underline"
              >
                Lihat Penilaian &rarr;
              </Link>
            </div>

            <div className="mt-4 flex-1">
              {loading ? (
                <div className="space-y-3">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="h-14 animate-pulse rounded-xl bg-gray-100" />
                  ))}
                </div>
              ) : !data?.recentSubmissions || data.recentSubmissions.length === 0 ? (
                <EmptyState
                  icon={FileSpreadsheet}
                  title="Belum Ada Pengumpulan"
                  description="Belum ada siswa yang mengumpulkan tugas atau projek terbaru."
                />
              ) : (
                <div className="space-y-2.5">
                  {data.recentSubmissions.map((sub) => (
                    <div
                      key={sub.id}
                      className="flex items-center justify-between rounded-xl border border-[#EBD4C5]/70 bg-white p-3.5 transition-colors hover:bg-[#FFF8F5]"
                    >
                      <div className="min-w-0 pr-3">
                        <div className="flex items-center gap-2">
                          <span
                            className={`rounded px-1.5 py-0.5 text-[10px] font-extrabold uppercase ${
                              sub.type === 'assignment'
                                ? 'bg-amber-100 text-amber-800'
                                : 'bg-purple-100 text-purple-800'
                            }`}
                          >
                            {sub.type === 'assignment' ? 'Tugas' : 'Projek'}
                          </span>
                          <p className="truncate text-xs font-bold text-[#30251D]">
                            {sub.studentName}
                          </p>
                        </div>
                        <p className="mt-1 truncate text-xs text-[#687385]">{sub.title}</p>
                        <p className="mt-0.5 text-[10px] text-gray-400">
                          {new Date(sub.submittedAt).toLocaleString('id-ID', {
                            dateStyle: 'medium',
                            timeStyle: 'short',
                          })}
                        </p>
                      </div>

                      <div className="shrink-0">
                        {sub.graded ? (
                          <StatusBadge status="graded" label={`Nilai: ${sub.score ?? '-'}`} />
                        ) : (
                          <StatusBadge status="pending" label="Belum Dinilai" />
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </TeacherLayout>
  )
}
