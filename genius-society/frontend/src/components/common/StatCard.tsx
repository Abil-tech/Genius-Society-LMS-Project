import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import type { LucideIcon } from 'lucide-react'

interface StatCardProps {
  title: string
  value: number | string
  icon: LucideIcon
  description?: string
  trend?: string
  highlight?: boolean
  linkTo?: string
  loading?: boolean
  badgeText?: string
  customFooter?: ReactNode
}

export default function StatCard({
  title,
  value,
  icon: Icon,
  description,
  trend,
  highlight = false,
  linkTo,
  loading = false,
  badgeText,
  customFooter,
}: StatCardProps) {
  const content = (
    <div
      className={`relative overflow-hidden rounded-2xl border p-5 transition-all duration-200 ${
        highlight
          ? 'border-[#FF9800] bg-gradient-to-br from-[#FFF8F5] via-[#FFFFFF] to-[#FEEADC] shadow-md ring-2 ring-[#FF9800]/30'
          : 'border-[#EBD4C5] bg-[#FFFFFF] shadow-sm hover:border-[#FF9800]/50 hover:shadow-md'
      }`}
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <p className="text-xs font-semibold uppercase tracking-wider text-[#687385] truncate">
            {title}
          </p>
          {loading ? (
            <div className="mt-2 h-8 w-20 animate-pulse rounded-lg bg-gray-200" />
          ) : (
            <div className="mt-1 flex items-baseline gap-2">
              <span
                className={`text-2xl sm:text-3xl font-extrabold tracking-tight ${
                  highlight ? 'text-[#FF9800]' : 'text-[#30251D]'
                }`}
              >
                {value}
              </span>
              {badgeText && (
                <span className="rounded-full bg-[#FF9800]/10 px-2 py-0.5 text-[11px] font-bold text-[#FF9800]">
                  {badgeText}
                </span>
              )}
            </div>
          )}
          {description && (
            <p className="mt-1 text-xs text-[#687385] line-clamp-1">{description}</p>
          )}
          {trend && (
            <p className="mt-1 text-[11px] font-medium text-emerald-600">{trend}</p>
          )}
        </div>

        <div
          className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-xl transition-colors ${
            highlight
              ? 'bg-[#FF9800] text-white shadow-sm'
              : 'bg-[#FFF8F5] text-[#FF9800] border border-[#EBD4C5]'
          }`}
          aria-hidden="true"
        >
          <Icon size={22} />
        </div>
      </div>

      {customFooter && <div className="mt-3 pt-3 border-t border-[#EBD4C5]/60">{customFooter}</div>}

      {linkTo && (
        <span className="mt-3 inline-flex items-center text-xs font-bold text-[#FF9800] group-hover:underline">
          Lihat detail &rarr;
        </span>
      )}
    </div>
  )

  if (linkTo) {
    return (
      <Link to={linkTo} className="group block focus:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9800] rounded-2xl">
        {content}
      </Link>
    )
  }

  return content
}
