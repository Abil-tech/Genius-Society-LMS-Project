import { useEffect, useState, useMemo } from 'react'
import {
  Clock,
  RefreshCw,
  MapPin,
  Calendar,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getTeacherSchedule } from '../../services/teacherService'
import type { TeacherScheduleItem } from '../../types/teacher'

const DAYS = [
  { id: 1, name: 'Senin' },
  { id: 2, name: 'Selasa' },
  { id: 3, name: 'Rabu' },
  { id: 4, name: 'Kamis' },
  { id: 5, name: 'Jumat' },
  { id: 6, name: 'Sabtu' },
  { id: 7, name: 'Minggu' },
]

export default function TeacherSchedulePage() {
  const [schedules, setSchedules] = useState<TeacherScheduleItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeDay, setActiveDay] = useState<number>(() => {
    // Current day: Sunday=0 -> 7, Monday=1 -> 1
    const d = new Date().getDay()
    return d === 0 ? 7 : d
  })

  const fetchSchedule = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await getTeacherSchedule()
      if (res.success && res.data) {
        setSchedules(res.data)
      } else {
        setError(res.error?.message || 'Gagal memuat jadwal mengajar.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat jadwal.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSchedule()
  }, [])

  const schedulesByDay = useMemo(() => {
    const map: Record<number, TeacherScheduleItem[]> = {
      1: [],
      2: [],
      3: [],
      4: [],
      5: [],
      6: [],
      7: [],
    }
    schedules.forEach((item) => {
      const d = item.day ?? 1
      if (map[d]) {
        map[d].push(item)
      }
    })
    // Sort by startTime
    Object.keys(map).forEach((k) => {
      map[Number(k)].sort((a, b) => a.startTime.localeCompare(b.startTime))
    })
    return map
  }, [schedules])

  const currentDayItems = schedulesByDay[activeDay] || []

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Jadwal Mengajar Mingguan
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Tampilan jadwal mengajar resmi per rombel dan ruangan (Read-only, ditetapkan oleh tim Kurikulum).
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchSchedule}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Segarkan</span>
            </button>
          </div>
        </div>

        {/* Day Selector Tabs (Senin - Minggu) */}
        <div className="flex overflow-x-auto gap-2 pb-1 scrollbar-none">
          {DAYS.map((day) => {
            const count = schedulesByDay[day.id]?.length ?? 0
            const isActive = activeDay === day.id
            return (
              <button
                key={day.id}
                type="button"
                onClick={() => setActiveDay(day.id)}
                className={`flex shrink-0 items-center gap-2 rounded-2xl border px-4 py-2.5 text-xs font-bold transition-all ${
                  isActive
                    ? 'border-[#FF9800] bg-[#FF9800] text-white shadow-xs'
                    : 'border-[#EBD4C5] bg-white text-[#30251D] hover:bg-[#FFF8F5]'
                }`}
              >
                <span>{day.name}</span>
                <span
                  className={`rounded-full px-1.5 py-0.2 text-[10px] font-extrabold ${
                    isActive ? 'bg-white/20 text-white' : 'bg-[#FFF8F5] text-[#FF9800]'
                  }`}
                >
                  {count}
                </span>
              </button>
            )
          })}
        </div>

        {/* Content */}
        {error ? (
          <ErrorState title="Gagal Memuat Jadwal" message={error} onRetry={fetchSchedule} />
        ) : loading ? (
          <div className="space-y-3">
            {[1, 2, 3, 4].map((i) => (
              <div
                key={i}
                className="h-20 rounded-2xl border border-[#EBD4C5] bg-white p-5 animate-pulse"
              />
            ))}
          </div>
        ) : currentDayItems.length === 0 ? (
          <EmptyState
            icon={Calendar}
            title={`Tidak Ada Jadwal di Hari ${DAYS.find((d) => d.id === activeDay)?.name}`}
            description="Tidak ada jam mengajar yang dijadwalkan pada hari ini. Anda dapat mempersiapkan materi atau memeriksa penugasan."
          />
        ) : (
          <div className="space-y-3">
            {currentDayItems.map((item) => (
              <div
                key={item.id}
                className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs hover:border-[#FF9800]/60 transition-colors"
              >
                <div className="flex items-start sm:items-center gap-3.5">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] text-[#FF9800]">
                    <Clock size={22} />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="rounded-md bg-[#FFF8F5] border border-[#EBD4C5] px-2 py-0.5 text-xs font-black text-[#FF9800]">
                        {item.startTime} - {item.endTime}
                      </span>
                      <span className="text-xs font-bold text-[#687385]">
                        Kelas {item.className}
                      </span>
                    </div>

                    <h3 className="mt-1.5 text-base sm:text-lg font-bold text-[#30251D]">
                      {item.subject}
                    </h3>
                  </div>
                </div>

                <div className="flex items-center gap-4 text-xs font-semibold text-[#687385] pl-15 sm:pl-0">
                  {item.room && (
                    <span className="inline-flex items-center gap-1.5 rounded-lg bg-stone-50 border border-stone-200 px-3 py-1.5 text-[#30251D]">
                      <MapPin size={14} className="text-[#FF9800]" />
                      <span>Ruang: {item.room}</span>
                    </span>
                  )}
                  <span className="inline-flex items-center gap-1 text-emerald-700 font-bold">
                    <span>Terjadwal</span>
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </TeacherLayout>
  )
}
