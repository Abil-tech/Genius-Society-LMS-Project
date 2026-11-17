import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  ArrowLeft,
  Calendar,
  Download,
  ExternalLink,
  Edit,
  Trash2,
  FileText,
  DoorOpen,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  getMaterialById,
  deleteMaterial,
  getDownloadPresignedUrl,
} from '../../services/teacherService'
import type { Material } from '../../types/teacher'

export default function TeacherMaterialDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const toast = useToast()

  const [material, setMaterial] = useState<Material | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [downloading, setDownloading] = useState(false)

  const fetchDetail = async () => {
    if (!id) return
    setLoading(true)
    setError(null)
    try {
      const res = await getMaterialById(id)
      if (res.success && res.data) {
        setMaterial(res.data)
      } else {
        setError(res.error?.message || 'Materi tidak ditemukan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal memuat materi.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchDetail()
  }, [id])

  const handleDelete = async () => {
    if (!id) return
    setIsDeleting(true)
    try {
      const res = await deleteMaterial(id)
      if (res.success) {
        toast.success('Materi Dihapus', 'Materi berhasil dihapus dari sistem.')
        navigate('/teacher/materials')
      } else {
        toast.error('Gagal Menghapus', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat menghapus.'
      toast.error('Gagal Menghapus', msg)
    } finally {
      setIsDeleting(false)
      setDeleteOpen(false)
    }
  }

  const handleDownload = async (key: string, fileName: string) => {
    setDownloading(true)
    try {
      const res = await getDownloadPresignedUrl(key)
      if (res.success && res.data?.url) {
        // Trigger download
        const a = document.createElement('a')
        a.href = res.data.url
        a.download = fileName
        a.target = '_blank'
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
      } else {
        toast.error('Gagal Mengunduh', res.error?.message || 'URL file tidak valid.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal menghasilkan URL unduhan.'
      toast.error('Gagal Mengunduh', msg)
    } finally {
      setDownloading(false)
    }
  }

  return (
    <TeacherLayout>
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Back Link */}
        <div className="flex items-center justify-between">
          <Link
            to="/teacher/materials"
            className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
          >
            <ArrowLeft size={16} />
            <span>Kembali ke Daftar Materi</span>
          </Link>

          {material && (
            <div className="flex items-center gap-2">
              <Link
                to={`/teacher/materials/${material.id}/edit`}
                className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] hover:bg-[#FFF8F5] shadow-xs transition-colors"
              >
                <Edit size={14} />
                <span>Ubah</span>
              </Link>
              <button
                type="button"
                onClick={() => setDeleteOpen(true)}
                className="inline-flex items-center gap-1.5 rounded-xl border border-rose-200 bg-white px-3.5 py-2 text-xs font-bold text-rose-600 hover:bg-rose-50 shadow-xs transition-colors"
              >
                <Trash2 size={14} />
                <span>Hapus</span>
              </button>
            </div>
          )}
        </div>

        {error ? (
          <ErrorState title="Materi Tidak Ditemukan" message={error} onRetry={fetchDetail} />
        ) : loading ? (
          <div className="space-y-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 animate-pulse">
            <div className="h-8 bg-gray-200 rounded-lg w-2/3" />
            <div className="h-20 bg-gray-200 rounded-lg" />
          </div>
        ) : material ? (
          <div className="space-y-6">
            {/* Card Content */}
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 sm:p-8 shadow-xs">
              <div className="flex flex-wrap items-center gap-2 mb-3">
                <span className="rounded-full bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-0.5 text-[11px] font-extrabold uppercase text-[#FF9800]">
                  {material.subjectName || 'Mata Pelajaran'}
                </span>
                <span className="inline-flex items-center gap-1 text-xs text-[#687385]">
                  <Calendar size={13} className="text-[#FF9800]" />
                  <span>
                    Terbit:{' '}
                    {new Date(material.publishAt).toLocaleDateString('id-ID', {
                      day: 'numeric',
                      month: 'long',
                      year: 'numeric',
                    })}
                  </span>
                </span>
              </div>

              <h1 className="text-xl sm:text-2xl font-extrabold text-[#30251D]">
                {material.title}
              </h1>

              {material.description && (
                <div className="mt-4 prose prose-stone max-w-none text-xs sm:text-sm text-[#30251D] leading-relaxed whitespace-pre-wrap">
                  {material.description}
                </div>
              )}

              {/* Kelas Dituju */}
              <div className="mt-6 pt-5 border-t border-[#EBD4C5]/60">
                <p className="text-xs font-bold uppercase tracking-wider text-[#687385] mb-2 flex items-center gap-1.5">
                  <DoorOpen size={14} className="text-[#FF9800]" />
                  <span>Diberikan Kepada Kelas:</span>
                </p>
                <div className="flex flex-wrap gap-1.5">
                  {material.classNames && material.classNames.length > 0 ? (
                    material.classNames.map((c, i) => (
                      <span
                        key={i}
                        className="rounded-lg bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-1 text-xs font-bold text-[#30251D]"
                      >
                        {c}
                      </span>
                    ))
                  ) : (
                    <span className="text-xs text-gray-400">Seluruh Kelas</span>
                  )}
                </div>
              </div>

              {/* Tautan Web */}
              {material.link && (
                <div className="mt-5 pt-4 border-t border-[#EBD4C5]/60">
                  <p className="text-xs font-bold uppercase tracking-wider text-[#687385] mb-2">
                    Tautan Web Terkait
                  </p>
                  <a
                    href={material.link}
                    target="_blank"
                    rel="noreferrer"
                    className="inline-flex items-center gap-2 rounded-xl border border-blue-200 bg-blue-50/60 px-4 py-2.5 text-xs sm:text-sm font-bold text-blue-800 hover:bg-blue-100/70 transition-colors"
                  >
                    <ExternalLink size={15} />
                    <span className="truncate max-w-md">{material.link}</span>
                  </a>
                </div>
              )}

              {/* File Lampiran R2 */}
              {material.file && (
                <div className="mt-5 pt-4 border-t border-[#EBD4C5]/60">
                  <p className="text-xs font-bold uppercase tracking-wider text-[#687385] mb-2">
                    Berkas Lampiran
                  </p>
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-4">
                    <div className="flex items-center gap-3">
                      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800]">
                        <FileText size={20} />
                      </div>
                      <div className="min-w-0">
                        <p className="font-bold text-xs sm:text-sm text-[#30251D] truncate">
                          {material.file.name}
                        </p>
                        <p className="text-[11px] text-[#687385]">
                          {(material.file.size / (1024 * 1024)).toFixed(2)} MB • {material.file.mime}
                        </p>
                      </div>
                    </div>

                    <button
                      type="button"
                      onClick={() => handleDownload(material.file!.key, material.file!.name)}
                      disabled={downloading}
                      className="inline-flex items-center justify-center gap-2 rounded-xl bg-[#FF9800] px-4 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] disabled:opacity-50 transition-colors"
                    >
                      <Download size={14} />
                      <span>{downloading ? 'Mengunduh...' : 'Unduh Berkas'}</span>
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        ) : null}

        {/* Confirm Delete */}
        <ConfirmDialog
          isOpen={deleteOpen}
          onClose={() => setDeleteOpen(false)}
          onConfirm={handleDelete}
          title="Hapus Materi"
          message={`Apakah Anda yakin ingin menghapus materi "${material?.title}"?`}
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
