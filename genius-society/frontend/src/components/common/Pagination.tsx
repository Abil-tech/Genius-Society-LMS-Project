import { ChevronLeft, ChevronRight } from 'lucide-react'

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  totalItems?: number
  pageSize?: number
}

export default function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  totalItems,
  pageSize,
}: PaginationProps) {
  if (totalPages <= 1 && !totalItems) return null

  // Calculate pages window
  const getPages = () => {
    const pages: (number | 'ellipsis')[] = []
    if (totalPages <= 5) {
      for (let i = 1; i <= totalPages; i++) pages.push(i)
    } else {
      pages.push(1)
      if (currentPage > 3) pages.push('ellipsis')
      const start = Math.max(2, currentPage - 1)
      const end = Math.min(totalPages - 1, currentPage + 1)
      for (let i = start; i <= end; i++) pages.push(i)
      if (currentPage < totalPages - 2) pages.push('ellipsis')
      pages.push(totalPages)
    }
    return pages
  }

  const startItem = (currentPage - 1) * (pageSize || 10) + 1
  const endItem = Math.min(currentPage * (pageSize || 10), totalItems || 0)

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-3 pt-4 text-xs text-[#687385]">
      {totalItems !== undefined && (
        <p className="order-2 sm:order-1 text-center sm:text-left">
          Menampilkan <span className="font-bold text-[#30251D]">{totalItems > 0 ? startItem : 0}</span> -{' '}
          <span className="font-bold text-[#30251D]">{endItem}</span> dari{' '}
          <span className="font-bold text-[#30251D]">{totalItems}</span> data
        </p>
      )}

      {totalPages > 1 && (
        <nav
          aria-label="Paginasi data"
          className="order-1 sm:order-2 flex items-center gap-1"
        >
          <button
            type="button"
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1}
            aria-label="Halaman sebelumnya"
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] bg-white text-[#30251D] disabled:opacity-40 disabled:cursor-not-allowed hover:bg-[#FFF8F5] transition-colors"
          >
            <ChevronLeft size={16} />
          </button>

          {getPages().map((page, idx) =>
            page === 'ellipsis' ? (
              <span key={`el-${idx}`} className="px-2 text-gray-400">
                ...
              </span>
            ) : (
              <button
                key={page}
                type="button"
                onClick={() => onPageChange(page)}
                aria-current={currentPage === page ? 'page' : undefined}
                className={`flex h-8 min-w-[32px] px-2 items-center justify-center rounded-lg text-xs font-bold transition-colors ${
                  currentPage === page
                    ? 'bg-[#FF9800] text-white shadow-sm'
                    : 'border border-[#EBD4C5] bg-white text-[#30251D] hover:bg-[#FFF8F5]'
                }`}
              >
                {page}
              </button>
            )
          )}

          <button
            type="button"
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            aria-label="Halaman berikutnya"
            className="flex h-8 w-8 items-center justify-center rounded-lg border border-[#EBD4C5] bg-white text-[#30251D] disabled:opacity-40 disabled:cursor-not-allowed hover:bg-[#FFF8F5] transition-colors"
          >
            <ChevronRight size={16} />
          </button>
        </nav>
      )}
    </div>
  )
}
