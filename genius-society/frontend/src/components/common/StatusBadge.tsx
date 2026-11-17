import type { LucideIcon } from 'lucide-react'
import { CheckCircle2, Clock, AlertCircle, FileText, XCircle, HelpCircle } from 'lucide-react'

export type BadgeVariant =
  | 'active'
  | 'ended'
  | 'draft'
  | 'published'
  | 'closed'
  | 'graded'
  | 'pending'
  | 'in_progress'
  | 'submitted'
  | 'success'
  | 'warning'
  | 'danger'
  | 'info'

interface StatusBadgeProps {
  status: BadgeVariant | string
  label?: string
  size?: 'sm' | 'md'
  showIcon?: boolean
}

const statusConfig: Record<
  string,
  { label: string; icon: LucideIcon; bg: string; text: string; border: string }
> = {
  active: {
    label: 'Aktif',
    icon: CheckCircle2,
    bg: 'bg-emerald-50',
    text: 'text-emerald-800',
    border: 'border-emerald-300',
  },
  published: {
    label: 'Dipublikasikan',
    icon: CheckCircle2,
    bg: 'bg-emerald-50',
    text: 'text-emerald-800',
    border: 'border-emerald-300',
  },
  success: {
    label: 'Berhasil',
    icon: CheckCircle2,
    bg: 'bg-emerald-50',
    text: 'text-emerald-800',
    border: 'border-emerald-300',
  },
  graded: {
    label: 'Sudah Dinilai',
    icon: CheckCircle2,
    bg: 'bg-blue-50',
    text: 'text-blue-800',
    border: 'border-blue-300',
  },
  pending: {
    label: 'Belum Dinilai',
    icon: Clock,
    bg: 'bg-amber-50',
    text: 'text-amber-800',
    border: 'border-amber-300',
  },
  in_progress: {
    label: 'Sedang Dikerjakan',
    icon: Clock,
    bg: 'bg-amber-50',
    text: 'text-amber-800',
    border: 'border-amber-300',
  },
  submitted: {
    label: 'Telah Dikumpulkan',
    icon: FileText,
    bg: 'bg-indigo-50',
    text: 'text-indigo-800',
    border: 'border-indigo-300',
  },
  draft: {
    label: 'Draf',
    icon: FileText,
    bg: 'bg-stone-100',
    text: 'text-stone-800',
    border: 'border-stone-300',
  },
  ended: {
    label: 'Berakhir',
    icon: Clock,
    bg: 'bg-gray-100',
    text: 'text-gray-700',
    border: 'border-gray-300',
  },
  closed: {
    label: 'Ditutup',
    icon: XCircle,
    bg: 'bg-gray-100',
    text: 'text-gray-700',
    border: 'border-gray-300',
  },
  warning: {
    label: 'Perhatian',
    icon: AlertCircle,
    bg: 'bg-amber-50',
    text: 'text-amber-800',
    border: 'border-amber-300',
  },
  danger: {
    label: 'Batal / Gagal',
    icon: XCircle,
    bg: 'bg-rose-50',
    text: 'text-rose-800',
    border: 'border-rose-300',
  },
  info: {
    label: 'Informasi',
    icon: HelpCircle,
    bg: 'bg-sky-50',
    text: 'text-sky-800',
    border: 'border-sky-300',
  },
}

export default function StatusBadge({
  status,
  label,
  size = 'sm',
  showIcon = true,
}: StatusBadgeProps) {
  const normKey = status.toLowerCase()
  const conf = statusConfig[normKey] || {
    label: label || status,
    icon: HelpCircle,
    bg: 'bg-gray-100',
    text: 'text-gray-800',
    border: 'border-gray-300',
  }

  const Icon = conf.icon
  const displayLabel = label || conf.label

  return (
    <span
      className={`inline-flex items-center gap-1.5 font-semibold rounded-full border ${conf.bg} ${conf.text} ${conf.border} ${
        size === 'sm' ? 'px-2.5 py-0.5 text-xs' : 'px-3 py-1 text-sm'
      }`}
    >
      {showIcon && <Icon size={size === 'sm' ? 12 : 14} className="shrink-0" aria-hidden="true" />}
      <span>{displayLabel}</span>
    </span>
  )
}
