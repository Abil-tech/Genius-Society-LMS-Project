import { createContext, useContext, useState, type ReactNode } from 'react'
import { CheckCircle2, AlertTriangle, AlertCircle, Info, X } from 'lucide-react'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface ToastMessage {
  id: string
  type: ToastType
  title: string
  message?: string
  duration?: number
}

interface ToastContextType {
  toasts: ToastMessage[]
  addToast: (toast: Omit<ToastMessage, 'id'>) => void
  removeToast: (id: string) => void
  success: (title: string, message?: string) => void
  error: (title: string, message?: string) => void
  warning: (title: string, message?: string) => void
  info: (title: string, message?: string) => void
}

const ToastContext = createContext<ToastContextType | undefined>(undefined)

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastMessage[]>([])

  const removeToast = (id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }

  const addToast = (toast: Omit<ToastMessage, 'id'>) => {
    const id = Math.random().toString(36).substring(2, 9)
    const newToast: ToastMessage = { ...toast, id }
    setToasts((prev) => [...prev, newToast])

    const dur = toast.duration ?? 4500
    if (dur > 0) {
      setTimeout(() => {
        removeToast(id)
      }, dur)
    }
  }

  const success = (title: string, message?: string) => addToast({ type: 'success', title, message })
  const error = (title: string, message?: string) => addToast({ type: 'error', title, message })
  const warning = (title: string, message?: string) => addToast({ type: 'warning', title, message })
  const info = (title: string, message?: string) => addToast({ type: 'info', title, message })

  const getIcon = (type: ToastType) => {
    switch (type) {
      case 'success':
        return <CheckCircle2 className="text-emerald-600 shrink-0" size={18} />
      case 'error':
        return <AlertCircle className="text-rose-600 shrink-0" size={18} />
      case 'warning':
        return <AlertTriangle className="text-amber-600 shrink-0" size={18} />
      default:
        return <Info className="text-sky-600 shrink-0" size={18} />
    }
  }

  const getBg = (type: ToastType) => {
    switch (type) {
      case 'success':
        return 'border-emerald-200 bg-white'
      case 'error':
        return 'border-rose-200 bg-white'
      case 'warning':
        return 'border-amber-200 bg-white'
      default:
        return 'border-sky-200 bg-white'
    }
  }

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast, success, error, warning, info }}>
      {children}

      {/* Toast container floating top-right */}
      <div
        aria-live="polite"
        className="fixed top-5 right-5 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none"
      >
        {toasts.map((toast) => (
          <div
            key={toast.id}
            role="status"
            className={`pointer-events-auto flex items-start gap-3 rounded-2xl border p-4 shadow-lg transition-all animate-in slide-in-from-top-2 ${getBg(
              toast.type
            )}`}
          >
            {getIcon(toast.type)}
            <div className="flex-1 min-w-0 pr-1">
              <p className="text-xs sm:text-sm font-bold text-[#30251D]">{toast.title}</p>
              {toast.message && (
                <p className="mt-0.5 text-xs text-[#687385] leading-relaxed">{toast.message}</p>
              )}
            </div>
            <button
              type="button"
              onClick={() => removeToast(toast.id)}
              aria-label="Tutup notifikasi"
              className="text-gray-400 hover:text-gray-600 p-0.5 rounded"
            >
              <X size={15} />
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  )
}

export function useToast() {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider')
  }
  return context
}
