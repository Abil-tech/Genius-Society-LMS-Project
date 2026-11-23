import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { ArrowLeft, Download, ExternalLink, FileText, Calendar } from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import EmptyState from '../../components/common/EmptyState'
import ErrorState from '../../components/common/ErrorState'
import {
  getMaterialDetail,
  getMaterialDownloadURL,
} from '../../services/student/studentMaterial.service'
import type { Material } from '../../types/student'

export default function StudentMaterialDetailPage() {
  const { id } = useParams<{ id: string }>()
  const [material, setMaterial] = useState<Material | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [downloadingIdx, setDownloadingIdx] = useState<number | null>(null)

  useEffect(() => {
    let isMounted = true
    const fetchData = async () => {
      if (!id) return
      try {
        setLoading(true)
        setError(null)
        const res = await getMaterialDetail(id)
        if (isMounted) setMaterial(res)
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
  }, [id])

  const handleDownload = async (index: number) => {
    if (!id) return
    try {
      setDownloadingIdx(index)
      const { url } = await getMaterialDownloadURL(id, index)
      window.open(url, '_blank')
    } catch {
      alert('Gagal mengunduh berkas. Coba lagi beberapa saat lagi.')
    } finally {
      setDownloadingIdx(null)
    }
  }

  if (error) {
    return (
      <StudentLayout>
        <ErrorState
          title="Materi Tidak Ditemukan"
          message={error}
          onRetry={() => window.location.reload()}
        />
      </StudentLayout>
    )
  }

  if (loading || !material) {
    return (
      <StudentLayout>
        <div className="animate-pulse space-y-4">
          <div className="h-6 w-32 bg-gray-200 rounded-lg" />
          <div className="h-40 bg-white border border-[#EBD4C5] rounded-2xl" />
        </div>
      </StudentLayout>
    )
  }

  return (
    <StudentLayout>
      <div className="space-y-6 max-w-4xl">
        <Link
          to="/student/materials"
          className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
        >
          <ArrowLeft size={16} />
          <span>Kembali ke Daftar Materi</span>
        </Link>

        {/* Header Material Info */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
          <div className="flex flex-wrap items-center gap-2 text-xs">
            <span className="rounded-lg bg-[#FEEADC] px-2.5 py-1 font-extrabold text-[#FF9800]">
              {material.subjectName}
            </span>
            <span className="flex items-center gap-1 text-[#687385]">
              <Calendar size={13} />
              <span>
                {new Date(material.createdAt).toLocaleDateString('id-ID', {
                  day: 'numeric',
                  month: 'long',
                  year: 'numeric',
                })}
              </span>
            </span>
            <span className="text-[#687385]">&bull; Guru: {material.teacherName || '-'}</span>
          </div>

          <h1 className="mt-3 text-lg sm:text-2xl font-black text-[#30251D]">
            {material.title}
          </h1>

          <div className="mt-4 prose prose-sm text-[#30251D] border-t border-[#EBD4C5]/60 pt-4 whitespace-pre-line leading-relaxed">
            {material.description || 'Tidak ada deskripsi tambahan untuk materi ini.'}
          </div>

          {material.link && (
            <div className="mt-4 p-3.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] flex items-center justify-between">
              <span className="text-xs font-semibold text-[#30251D] truncate max-w-md">
                Tautan: {material.link}
              </span>
              <a
                href={material.link}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-1 text-xs font-bold text-[#FF9800] hover:underline shrink-0"
              >
                <span>Buka Tautan</span>
                <ExternalLink size={14} />
              </a>
            </div>
          )}
        </div>

        {/* Attachments Section */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
          <h2 className="text-sm font-bold text-[#30251D] flex items-center gap-2 pb-3 border-b border-[#EBD4C5]/60">
            <FileText size={16} className="text-[#FF9800]" />
            <span>Berkas Lampiran Pembelajaran</span>
          </h2>

          <div className="mt-4 space-y-3">
            {!material.attachments || material.attachments.length === 0 ? (
              <EmptyState
                title="Tidak ada berkas lampiran"
                description="Materi ini tidak menyertakan berkas dokumen untuk diunduh."
              />
            ) : (
              material.attachments.map((file, idx) => (
                <div
                  key={file.key || idx}
                  className="flex items-center justify-between p-3.5 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/60"
                >
                  <div className="flex items-center gap-3 min-w-0 pr-3">
                    <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800]">
                      <FileText size={16} />
                    </div>
                    <div className="min-w-0">
                      <p className="truncate text-xs font-bold text-[#30251D]">
                        {file.name}
                      </p>
                      <p className="text-[11px] text-[#687385] mt-0.5">
                        {file.size ? `${(file.size / (1024 * 1024)).toFixed(2)} MB` : '-'}
                      </p>
                    </div>
                  </div>

                  <button
                    type="button"
                    onClick={() => handleDownload(idx)}
                    disabled={downloadingIdx === idx}
                    className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-3.5 py-1.5 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors disabled:opacity-60 shrink-0"
                  >
                    <Download size={13} />
                    <span>{downloadingIdx === idx ? 'Menyiapkan...' : 'Unduh'}</span>
                  </button>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </StudentLayout>
  )
}
