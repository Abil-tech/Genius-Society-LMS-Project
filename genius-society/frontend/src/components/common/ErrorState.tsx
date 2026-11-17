import { AlertTriangle, RotateCcw } from 'lucide-react'

interface ErrorStateProps {
  title?: string
  message: string
  onRetry?: () => void
  compact?: boolean
}

export default function ErrorState({
  title = 'Terjadi Kesalahan',
  message,
  onRetry,
  compact = false,
}: ErrorStateProps) {
  if (compact) {
    return (
      <div
        role="alert"
        className="flex items-center justify-between gap-3 rounded-xl border border-rose-200 bg-rose-50/70 p-3 text-xs text-rose-900"
      >
        <div className="flex items-center gap-2 min-w-0">
          <AlertTriangle size={16} className="text-rose-600 shrink-0" aria-hidden="true" />
          <span className="truncate">{message}</span>
        </div>
        {onRetry && (
          <button
            type="button"
            onClick={onRetry}
            className="shrink-0 font-bold text-rose-700 underline hover:text-rose-900 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-500 rounded"
          >
            Coba lagi
          </button>
        )}
      </div>
    )
  }

  return (
    <div
      role="alert"
      className="flex flex-col items-center justify-center rounded-2xl border border-rose-200 bg-rose-50/40 p-8 sm:p-12 text-center"
    >
      <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-rose-100 text-rose-600">
        <AlertTriangle size={24} aria-hidden="true" />
      </div>
      <h3 className="mt-4 text-base font-bold text-rose-900">{title}</h3>
      <p className="mt-1.5 max-w-md text-xs sm:text-sm text-rose-700">{message}</p>

      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="mt-5 inline-flex items-center gap-2 rounded-xl bg-white border border-rose-300 px-4 py-2 text-xs font-bold text-rose-800 shadow-sm hover:bg-rose-50 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-500 transition-colors"
        >
          <RotateCcw size={14} />
          <span>Coba Lagi</span>
        </button>
      )}
    </div>
  )
}
