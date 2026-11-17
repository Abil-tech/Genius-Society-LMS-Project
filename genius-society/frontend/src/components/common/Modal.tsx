import { useEffect, type ReactNode } from 'react'
import { X } from 'lucide-react'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string
  description?: string
  children: ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  footer?: ReactNode
}

export default function Modal({
  isOpen,
  onClose,
  title,
  description,
  children,
  size = 'md',
  footer,
}: ModalProps) {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose()
      }
    }
    if (isOpen) {
      document.body.style.overflow = 'hidden'
      window.addEventListener('keydown', handleKeyDown)
    }
    return () => {
      document.body.style.overflow = ''
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen, onClose])

  if (!isOpen) return null

  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
    full: 'max-w-6xl',
  }[size]

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6 overflow-y-auto"
    >
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 backdrop-blur-xs transition-opacity"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Modal Dialog Card */}
      <div
        className={`relative z-10 w-full ${sizeClasses} rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] shadow-xl transition-all duration-200 my-auto`}
      >
        {/* Header */}
        <div className="flex items-start justify-between border-b border-[#EBD4C5] p-5 sm:p-6">
          <div className="pr-4">
            <h3 id="modal-title" className="text-base sm:text-lg font-bold text-[#30251D]">
              {title}
            </h3>
            {description && (
              <p className="mt-1 text-xs sm:text-sm text-[#687385]">{description}</p>
            )}
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Tutup dialog"
            className="rounded-lg p-1.5 text-gray-400 hover:bg-[#FFF8F5] hover:text-[#30251D] focus:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9800] transition-colors"
          >
            <X size={18} />
          </button>
        </div>

        {/* Content */}
        <div className="p-5 sm:p-6 max-h-[75vh] overflow-y-auto">{children}</div>

        {/* Footer */}
        {footer && (
          <div className="flex flex-col-reverse sm:flex-row items-center justify-end gap-2 border-t border-[#EBD4C5] bg-[#FFF8F5]/60 p-4 sm:px-6">
            {footer}
          </div>
        )}
      </div>
    </div>
  )
}
