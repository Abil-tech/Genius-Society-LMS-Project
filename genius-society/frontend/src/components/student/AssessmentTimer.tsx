import { useState, useEffect } from 'react'
import { Clock } from 'lucide-react'

interface AssessmentTimerProps {
  durationSeconds: number
  startedAt: string
  onExpired: () => void
}

export default function AssessmentTimer({
  durationSeconds,
  startedAt,
  onExpired,
}: AssessmentTimerProps) {
  const [secondsLeft, setSecondsLeft] = useState<number>(() => {
    const startMs = new Date(startedAt).getTime()
    const nowMs = Date.now()
    const elapsedSeconds = Math.max(0, Math.floor((nowMs - startMs) / 1000))
    return Math.max(0, durationSeconds - elapsedSeconds)
  })

  useEffect(() => {
    if (secondsLeft <= 0) {
      onExpired()
      return
    }

    const interval = setInterval(() => {
      setSecondsLeft((prev) => {
        if (prev <= 1) {
          clearInterval(interval)
          onExpired()
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(interval)
  }, [secondsLeft, onExpired])

  const minutes = Math.floor(secondsLeft / 60)
  const seconds = secondsLeft % 60
  const formattedTime = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`

  let colorClasses = 'text-[#30251D] bg-[#FFF8F5] border-[#EBD4C5]'
  if (secondsLeft < 60) {
    colorClasses = 'text-[#DC2626] bg-rose-50 border-rose-200 animate-pulse'
  } else if (secondsLeft < 300) {
    colorClasses = 'text-[#FF9800] bg-orange-50 border-[#FF9800]/40'
  }

  return (
    <div
      className={`flex items-center gap-2 rounded-xl border px-3.5 py-1.5 text-xs font-extrabold shadow-xs transition-colors ${colorClasses}`}
    >
      <Clock size={15} />
      <span>Sisa Waktu: {formattedTime}</span>
    </div>
  )
}
