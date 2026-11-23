import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Search, FileText, ExternalLink, ArrowRight } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import { getMaterials } from '../../services/student/studentMaterial.service'
import type { Material } from '../../types/student'

export default function StudentMaterialsPage() {
  const [materials, setMaterials] = useState<Material[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [selectedSubject, setSelectedSubject] = useState('')

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        const res = await getMaterials(selectedSubject || undefined)
        if (isMounted) setMaterials(res)
      } catch (err: unknown) {
        if (isMounted) {
          const msg = err instanceof Error ? err.message : 'Gagal memuat materi'
          setError(msg)
        }
      } finally {
        if (isMounted) setLoading(false)
      }
    }
    fetchData()
    return () => {
      isMounted = false
    }
  }, [selectedSubject])

  // Extract unique subjects
  const subjects = Array.from(
    new Map(materials.map((m) => [m.subjectId, m.subjectName])).entries()
  )

  const filteredMaterials = materials.filter((m) =>
    m.title.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <StudentLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-xl sm:text-2xl font-black text-[#30251D]">
            Materi Pembelajaran
          </h1>
          <p className="mt-1 text-xs text-[#687385]">
            Akses materi bacaan, dokumen panduan, dan tautan belajar dari gurumu.
          </p>
        </div>

        {/* Filter & Search Bar */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Cari materi pembelajaran..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2.5 pl-10 pr-4 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
            />
          </div>
          {subjects.length > 0 && (
            <select
              value={selectedSubject}
              onChange={(e) => setSelectedSubject(e.target.value)}
              className="rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2.5 text-xs font-semibold text-[#30251D] focus:border-[#FF9800] focus:outline-none"
            >
              <option value="">Semua Mata Pelajaran</option>
              {subjects.map(([id, name]) => (
                <option key={id} value={id}>
                  {name}
                </option>
              ))}
            </select>
          )}
        </div>

        {/* Material Grid */}
        {error ? (
          <ErrorState
            title="Gagal Memuat Materi"
            message={error}
            onRetry={() => setSelectedSubject('')}
          />
        ) : loading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="h-36 rounded-2xl bg-white border border-[#EBD4C5] p-5 animate-pulse" />
            <div className="h-36 rounded-2xl bg-white border border-[#EBD4C5] p-5 animate-pulse" />
          </div>
        ) : filteredMaterials.length === 0 ? (
          <EmptyState
            title="Belum ada materi"
            description={
              search
                ? `Tidak ada materi yang cocok dengan pencarian "${search}".`
                : 'Belum ada materi pembelajaran yang dipublikasikan untuk kelasmu.'
            }
          />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {filteredMaterials.map((m) => (
              <div
                key={m.id}
                className="flex flex-col justify-between rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs hover:border-[#FF9800] transition-colors"
              >
                <div>
                  <div className="flex items-center justify-between gap-2">
                    <span className="rounded-md bg-[#FEEADC] px-2.5 py-0.5 text-[10px] font-extrabold text-[#FF9800]">
                      {m.subjectName}
                    </span>
                    <span className="text-[11px] text-[#687385]">
                      {new Date(m.createdAt).toLocaleDateString('id-ID', {
                        day: 'numeric',
                        month: 'short',
                        year: 'numeric',
                      })}
                    </span>
                  </div>

                  <h3 className="mt-2 text-sm sm:text-base font-extrabold text-[#30251D] line-clamp-1">
                    {m.title}
                  </h3>
                  <p className="mt-1 text-xs text-[#687385] line-clamp-2">
                    {m.description || 'Tidak ada deskripsi tambahan.'}
                  </p>
                </div>

                <div className="mt-4 pt-3 border-t border-[#EBD4C5]/60 flex items-center justify-between">
                  <div className="flex items-center gap-3 text-xs text-[#687385]">
                    {m.fileCount > 0 && (
                      <span className="flex items-center gap-1">
                        <FileText size={14} className="text-[#FF9800]" />
                        <span>{m.fileCount} Berkas</span>
                      </span>
                    )}
                    {m.link && (
                      <span className="flex items-center gap-1">
                        <ExternalLink size={14} className="text-[#FF9800]" />
                        <span>Tautan</span>
                      </span>
                    )}
                  </div>

                  <Link
                    to={`/student/materials/${m.id}`}
                    className="inline-flex items-center gap-1 text-xs font-bold text-[#FF9800] hover:underline"
                  >
                    <span>Buka Materi</span>
                    <ArrowRight size={14} />
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </StudentLayout>
  )
}
