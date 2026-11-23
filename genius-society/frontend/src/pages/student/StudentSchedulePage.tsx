import { useState, useEffect } from 'react'
import { Clock, MapPin, User } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getSchedule } from '../../services/student/studentSchedule.service'
import type { ScheduleItem } from '../../types/student'

const days = [
  { id: 1, name: 'Senin' },
  { id: 2, name: 'Selasa' },
  { id: 3, name: 'Rabu' },
  { id: 4, name: 'Kamis' },
  { id: 5, name: 'Jumat' },
  { id: 6, name: 'Sabtu' },
]

export default function StudentSchedulePage() {
  const [schedules, setSchedules] = useState<ScheduleItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedDay, setSelectedDay] = useState<number>(() => {
    const cur = new Date().getDay()
    return cur >= 1 && cur <= 6 ? cur : 1
  })

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getSchedule()
        if (isMounted) setSchedules(res)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat jadwal'
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

  const currentDaySchedules = schedules.filter((s) => s.day === selectedDay)

  return (
    <StudentLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
            Jadwal Pelajaran Kelas
          </h1>
          <p className="mt-1 text-xs text-[#687385]">
            Jadwal mingguan mata pelajaran, guru pengajar, dan ruang kelasmu.
          </p>
        </div>

        {/* Day Tabs */}
        <div className="flex items-center gap-2 overflow-x-auto pb-1">
          {days.map((d) => (
            <button
              key={d.id}
              type="button"
              onClick={() => setSelectedDay(d.id)}
              className={`rounded-xl px-4 py-2 text-xs font-bold transition-all shrink-0 ${
                selectedDay === d.id
                  ? 'bg-[#FF9800] text-white shadow-xs'
                  : 'bg-white border border-[#EBD4C5] text-[#30251D] hover:bg-[#FFF8F5]'
              }`}
            >
              {d.name}
            </button>
          ))}
        </div>

        {/* Schedule List for the Day */}
        {error ? (
          <ErrorState
            title="Gagal Memuat Jadwal"
            message={error}
            onRetry={() => window.location.reload()}
          />
        ) : loading ? (
          <div className="space-y-3">
            <div className="h-20 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
            <div className="h-20 bg-white border border-[#EBD4C5] rounded-2xl animate-pulse" />
          </div>
        ) : currentDaySchedules.length === 0 ? (
          <EmptyState
            title={`Tidak ada jadwal di hari ${days.find((d) => d.id === selectedDay)?.name}`}
            description="Tidak ada jam pelajaran yang dijadwalkan pada hari ini."
          />
        ) : (
          <div className="space-y-3">
            {currentDaySchedules.map((sc) => (
              <div
                key={sc.id}
                className="flex flex-col sm:flex-row sm:items-center justify-between p-4 sm:p-5 rounded-2xl border border-[#EBD4C5] bg-white shadow-xs gap-3"
              >
                <div>
                  <h3 className="text-sm sm:text-base font-extrabold text-[#30251D]">
                    {sc.subject}
                  </h3>
                  <div className="mt-1 flex flex-wrap items-center gap-3 text-xs text-[#687385]">
                    <span className="flex items-center gap-1">
                      <User size={13} className="text-[#FF9800]" />
                      <span>Guru: {sc.teacher || '-'}</span>
                    </span>
                    <span className="flex items-center gap-1">
                      <MapPin size={13} className="text-[#FF9800]" />
                      <span>Ruangan: {sc.room || '-'}</span>
                    </span>
                  </div>
                </div>

                <div className="flex items-center gap-2 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] px-3.5 py-1.5 shrink-0">
                  <Clock size={14} className="text-[#FF9800]" />
                  <span className="text-xs font-mono font-bold text-[#30251D]">
                    {sc.startTime} - {sc.endTime}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </StudentLayout>
  )
}
