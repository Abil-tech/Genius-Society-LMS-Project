import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import type { LucideIcon } from 'lucide-react'
import { Inbox } from 'lucide-react'

interface EmptyStateProps {
  icon?: LucideIcon
  title: string
  description?: string
  actionLabel?: string
  actionOnClick?: () => void
  actionHref?: string
  actionIcon?: LucideIcon
  extra?: ReactNode
}

export default function EmptyState({
  icon: Icon = Inbox,
  title,
  description,
  actionLabel,
  actionOnClick,
  actionHref,
  actionIcon: ActionIcon,
  extra,
}: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-[#EBD4C5] bg-[#FFFFFF] p-8 sm:p-12 text-center shadow-sm">
      <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-[#FFF8F5] border border-[#EBD4C5] text-[#FF9800]">
        <Icon size={28} aria-hidden="true" />
      </div>
      <h3 className="mt-4 text-base sm:text-lg font-bold text-[#30251D]">{title}</h3>
      {description && (
        <p className="mt-1.5 max-w-md text-xs sm:text-sm text-[#687385]">{description}</p>
      )}

      {actionLabel && (
        <div className="mt-6">
          {actionHref ? (
            <Link
              to={actionHref}
              className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2.5 text-xs sm:text-sm font-bold text-white shadow-sm hover:bg-[#E68A00] focus:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9800] focus-visible:ring-offset-2 transition-colors"
            >
              {ActionIcon && <ActionIcon size={16} />}
              <span>{actionLabel}</span>
            </Link>
          ) : (
            <button
              type="button"
              onClick={actionOnClick}
              className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2.5 text-xs sm:text-sm font-bold text-white shadow-sm hover:bg-[#E68A00] focus:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9800] focus-visible:ring-offset-2 transition-colors"
            >
              {ActionIcon && <ActionIcon size={16} />}
              <span>{actionLabel}</span>
            </button>
          )}
        </div>
      )}

      {extra && <div className="mt-4">{extra}</div>}
    </div>
  )
}
