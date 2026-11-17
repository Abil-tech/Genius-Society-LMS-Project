import type { ReactNode } from 'react'

export interface Column<T> {
  key: string
  header: string
  render?: (row: T, index: number) => ReactNode
  className?: string
  align?: 'left' | 'center' | 'right'
}

interface DataTableProps<T> {
  columns: Column<T>[]
  data: T[]
  loading?: boolean
  emptyMessage?: string
  emptySubMessage?: string
  emptyAction?: ReactNode
  keyExtractor?: (row: T, index: number) => string | number
  onRowClick?: (row: T) => void
}

export default function DataTable<T>({
  columns,
  data,
  loading = false,
  emptyMessage = 'Belum ada data',
  emptySubMessage,
  emptyAction,
  keyExtractor,
  onRowClick,
}: DataTableProps<T>) {
  return (
    <div className="w-full overflow-hidden rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] shadow-sm">
      <div className="w-full overflow-x-auto">
        <table className="w-full border-collapse text-left text-xs sm:text-sm">
          <thead>
            <tr className="border-b border-[#EBD4C5] bg-[#FFF8F5]/80 text-[#30251D]">
              {columns.map((col) => (
                <th
                  key={col.key}
                  scope="col"
                  className={`px-4 sm:px-6 py-3.5 text-xs font-bold tracking-wider uppercase text-[#687385] ${
                    col.align === 'center'
                      ? 'text-center'
                      : col.align === 'right'
                      ? 'text-right'
                      : 'text-left'
                  } ${col.className || ''}`}
                >
                  {col.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-[#EBD4C5]/60">
            {loading ? (
              // Skeleton rows
              Array.from({ length: 5 }).map((_, rIdx) => (
                <tr key={`skeleton-${rIdx}`} className="animate-pulse">
                  {columns.map((col) => (
                    <td key={col.key} className="px-4 sm:px-6 py-4">
                      <div className="h-4 w-3/4 rounded bg-gray-200" />
                    </td>
                  ))}
                </tr>
              ))
            ) : data.length === 0 ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-6 py-12 text-center text-[#687385]"
                >
                  <p className="font-semibold text-sm text-[#30251D]">{emptyMessage}</p>
                  {emptySubMessage && (
                    <p className="mt-1 text-xs text-[#687385]">{emptySubMessage}</p>
                  )}
                  {emptyAction && <div className="mt-4">{emptyAction}</div>}
                </td>
              </tr>
            ) : (
              data.map((row, index) => {
                const key = keyExtractor
                  ? keyExtractor(row, index)
                  : (row as { id?: string | number }).id ?? index
                const isClickable = Boolean(onRowClick)

                return (
                  <tr
                    key={key}
                    onClick={() => onRowClick?.(row)}
                    className={`transition-colors ${
                      isClickable
                        ? 'cursor-pointer hover:bg-[#FFF8F5]/80'
                        : 'hover:bg-stone-50/50'
                    }`}
                  >
                    {columns.map((col) => (
                      <td
                        key={col.key}
                        className={`px-4 sm:px-6 py-3.5 text-[#30251D] ${
                          col.align === 'center'
                            ? 'text-center'
                            : col.align === 'right'
                            ? 'text-right'
                            : 'text-left'
                        } ${col.className || ''}`}
                      >
                        {col.render
                          ? col.render(row, index)
                          : String((row as Record<string, unknown>)[col.key] ?? '-')}
                      </td>
                    ))}
                  </tr>
                )
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
