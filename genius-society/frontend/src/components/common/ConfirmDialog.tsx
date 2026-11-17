import { AlertTriangle, Archive, Trash2 } from 'lucide-react'
import Modal from './Modal'

interface ConfirmDialogProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'primary'
  isLoading?: boolean
  hasSubmissions?: boolean
  archiveNotice?: string
}

export default function ConfirmDialog({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  confirmText,
  cancelText = 'Batal',
  variant = 'danger',
  isLoading = false,
  hasSubmissions = false,
  archiveNotice,
}: ConfirmDialogProps) {
  const isArchiveMode = hasSubmissions
  const defaultConfirmText = isArchiveMode
    ? 'Arsipkan Data'
    : variant === 'danger'
    ? 'Hapus'
    : 'Lanjutkan'

  const effectiveConfirmText = confirmText || defaultConfirmText

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={title}
      size="sm"
      footer={
        <>
          <button
            type="button"
            onClick={onClose}
            disabled={isLoading}
            className="w-full sm:w-auto rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 text-xs sm:text-sm font-semibold text-[#30251D] hover:bg-[#FFF8F5] disabled:opacity-50 transition-colors"
          >
            {cancelText}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={isLoading}
            className={`w-full sm:w-auto inline-flex items-center justify-center gap-1.5 rounded-xl px-4 py-2 text-xs sm:text-sm font-bold text-white shadow-sm disabled:opacity-50 transition-colors ${
              isArchiveMode
                ? 'bg-amber-600 hover:bg-amber-700'
                : variant === 'danger'
                ? 'bg-rose-600 hover:bg-rose-700'
                : 'bg-[#FF9800] hover:bg-[#E68A00]'
            }`}
          >
            {isLoading ? (
              <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
            ) : isArchiveMode ? (
              <Archive size={15} />
            ) : (
              <Trash2 size={15} />
            )}
            <span>{effectiveConfirmText}</span>
          </button>
        </>
      }
    >
      <div className="space-y-3">
        <div className="flex items-start gap-3">
          <div
            className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl ${
              isArchiveMode
                ? 'bg-amber-100 text-amber-700'
                : variant === 'danger'
                ? 'bg-rose-100 text-rose-700'
                : 'bg-[#FFF8F5] text-[#FF9800]'
            }`}
          >
            <AlertTriangle size={20} />
          </div>
          <div className="text-xs sm:text-sm text-[#30251D]">
            <p>{message}</p>
          </div>
        </div>

        {isArchiveMode && (
          <div className="rounded-xl border border-amber-200 bg-amber-50 p-3 text-xs text-amber-900">
            <p className="font-bold flex items-center gap-1.5">
              <Archive size={14} className="text-amber-700 shrink-0" />
              Perhatian: Data Memiliki Pengumpulan (Submission)
            </p>
            <p className="mt-1 text-amber-800">
              {archiveNotice ||
                'Item ini sudah memiliki data pengumpulan dari siswa. Sesuai kebijakan integritas akademik, data tidak akan dihapus permanen melainkan diarsipkan agar riwayat nilai tetap aman.'}
            </p>
          </div>
        )}
      </div>
    </Modal>
  )
}
