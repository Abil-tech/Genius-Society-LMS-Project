import { useState, useEffect, type FormEvent } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  ArrowLeft,
  Clock,
  Download,
  Upload,
  CheckCircle2,
  FileText,
} from 'lucide-react'
import StudentLayout from '../../components/student/StudentLayout'
import ErrorState from '../../components/common/ErrorState'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import { useToast } from '../../components/common/Toast'
import {
  getProjectDetail,
  submitProject,
  resubmitProject,
  getProjectFileDownloadURL,
} from '../../services/student/studentProject.service'
import type { Project } from '../../types/student'

export default function StudentProjectDetailPage() {
  const { id } = useParams<{ id: string }>()
  const toast = useToast()

  const [project, setProject] = useState<Project | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [file, setFile] = useState<File | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [downloadingIdx, setDownloadingIdx] = useState<number | null>(null)

  const fetchData = async () => {
    if (!id) return
    try {
      setLoading(true)
      setError(null)
      const res = await getProjectDetail(id)
      setProject(res)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal memuat projek'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [id])

  const handleDownload = async (index: number) => {
    if (!id) return
    try {
      setDownloadingIdx(index)
      const { url } = await getProjectFileDownloadURL(id, index)
      window.open(url, '_blank')
    } catch {
      toast.error('Gagal Mengunduh', 'Berkas panduan projek tidak dapat diunduh.')
    } finally {
      setDownloadingIdx(null)
    }
  }

  const handleConfirmSubmit = async () => {
    if (!id || !file) return
    try {
      setSubmitting(true)
      const formData = new FormData()
      formData.append('file', file)

      const isResubmit = Boolean(project?.submission?.submitted)
      if (isResubmit) {
        await resubmitProject(id, formData)
        toast.success('Berhasil Mengirim Ulang', 'Revisi projek telah tersimpan.')
      } else {
        await submitProject(id, formData)
        toast.success('Projek Berhasil Dikumpulkan', 'Laporan projekmu telah diterima oleh guru.')
      }

      setFile(null)
      await fetchData()
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal mengumpulkan projek'
      toast.error('Pengumpulan Gagal', msg)
    } finally {
      setSubmitting(false)
      setConfirmOpen(false)
    }
  }

  if (error) {
    return (
      <StudentLayout>
        <ErrorState title="Projek Tidak Ditemukan" message={error} onRetry={fetchData} />
      </StudentLayout>
    )
  }

  if (loading || !project) {
    return (
      <StudentLayout>
        <div className="animate-pulse space-y-4 max-w-4xl">
          <div className="h-6 w-32 bg-gray-200 rounded-lg" />
          <div className="h-44 bg-white border border-[#EBD4C5] rounded-2xl" />
        </div>
      </StudentLayout>
    )
  }

  const isSubmitted = project.submission?.submitted
  const canSubmit = !isSubmitted || project.allowResubmit

  return (
    <StudentLayout>
      <div className="space-y-6 max-w-4xl">
        <Link
          to="/student/projects"
          className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
        >
          <ArrowLeft size={16} />
          <span>Kembali ke Daftar Projek</span>
        </Link>

        {/* Project Info */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
          <div className="flex flex-wrap items-center gap-2 text-xs">
            <span className="rounded-lg bg-[#FEEADC] px-2.5 py-1 font-extrabold text-[#FF9800]">
              {project.subjectName}
            </span>
            <span className="flex items-center gap-1 text-[#687385]">
              <Clock size={13} className="text-[#FF9800]" />
              <span>
                Deadline:{' '}
                {new Date(project.deadline).toLocaleDateString('id-ID', {
                  day: 'numeric',
                  month: 'long',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </span>
            </span>
            <span className="text-[#687385]">&bull; Guru: {project.teacherName || '-'}</span>
          </div>

          <h1 className="mt-3 text-lg sm:text-2xl font-black text-[#30251D]">
            {project.title}
          </h1>

          {project.description && (
            <div className="mt-3 text-xs sm:text-sm text-[#687385] leading-relaxed">
              {project.description}
            </div>
          )}

          {/* Instructions & Rubric */}
          {project.instructions && (
            <div className="mt-5 border-t border-[#EBD4C5]/60 pt-4">
              <span className="text-xs font-bold uppercase tracking-wider text-[#30251D]">
                Panduan Pelaksanaan & Rubrik Penilaian:
              </span>
              <div className="mt-2 text-xs sm:text-sm text-[#30251D] whitespace-pre-line leading-relaxed p-4 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/60">
                {project.instructions}
              </div>
            </div>
          )}

          {/* Question Attachment */}
          {project.attachments && project.attachments.length > 0 && (
            <div className="mt-5 border-t border-[#EBD4C5]/60 pt-4">
              <span className="text-xs font-bold uppercase tracking-wider text-[#687385]">
                Lampiran Berkas Panduan:
              </span>
              <div className="mt-2 space-y-2">
                {project.attachments.map((att, idx) => (
                  <div
                    key={att.key || idx}
                    className="flex items-center justify-between p-3 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]"
                  >
                    <div className="flex items-center gap-2 min-w-0 pr-2">
                      <FileText size={15} className="text-[#FF9800] shrink-0" />
                      <span className="text-xs font-semibold text-[#30251D] truncate">
                        {att.name}
                      </span>
                    </div>
                    <button
                      type="button"
                      onClick={() => handleDownload(idx)}
                      disabled={downloadingIdx === idx}
                      className="inline-flex items-center gap-1 text-xs font-bold text-[#FF9800] hover:underline shrink-0"
                    >
                      <Download size={13} />
                      <span>{downloadingIdx === idx ? 'Mengunduh...' : 'Unduh'}</span>
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Submission Status Section */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
          <h2 className="text-sm font-bold text-[#30251D] pb-3 border-b border-[#EBD4C5]/60 flex items-center justify-between">
            <span>Status Pengumpulan Projek</span>
            {isSubmitted && (
              <span className="flex items-center gap-1 text-xs font-bold text-emerald-600">
                <CheckCircle2 size={15} />
                <span>Sudah Dikumpulkan (Versi {project.submission?.version})</span>
              </span>
            )}
          </h2>

          {isSubmitted ? (
            <div className="mt-4 space-y-4">
              {project.submission?.score !== undefined && (
                <div className="flex items-center justify-between p-4 rounded-xl border border-emerald-200 bg-emerald-50/60">
                  <div>
                    <span className="text-xs font-bold text-emerald-800">
                      Nilai Projek dari Guru:
                    </span>
                    {project.submission?.feedback && (
                      <p className="mt-1 text-xs text-emerald-700 italic">
                        "{project.submission.feedback}"
                      </p>
                    )}
                  </div>
                  <span className="text-2xl font-black text-emerald-700">
                    {project.submission.score}
                  </span>
                </div>
              )}

              <div className="text-xs text-[#687385] space-y-1">
                <p>
                  Waktu Submit:{' '}
                  {project.submission?.submittedAt
                    ? new Date(project.submission.submittedAt).toLocaleString('id-ID')
                    : '-'}
                  {project.submission?.isLate && (
                    <span className="ml-2 font-bold text-rose-600">(Terlambat)</span>
                  )}
                </p>
                {project.submission?.files && project.submission.files.length > 0 && (
                  <p>
                    Berkas Dikirim:{' '}
                    {project.submission.files.map((f) => f.name).join(', ')}
                  </p>
                )}
              </div>
            </div>
          ) : (
            <p className="mt-4 text-xs text-[#687385]">
              Kamu belum mengumpulkan berkas hasil projek ini.
            </p>
          )}

          {/* Form Upload */}
          {canSubmit && (
            <form
              onSubmit={(e: FormEvent) => {
                e.preventDefault()
                setConfirmOpen(true)
              }}
              className="mt-6 pt-5 border-t border-[#EBD4C5]/60 space-y-4"
            >
              <h3 className="text-xs font-bold uppercase tracking-wider text-[#30251D]">
                {isSubmitted ? 'Unggah Pengumpulan Ulang Projek' : 'Kirim Laporan Projek'}
              </h3>

              <div>
                <label className="block text-xs font-bold text-[#30251D] mb-1">
                  Pilih Berkas Laporan Projek (PDF, ZIP, PPTX, DOCX maks 50MB)
                </label>
                <div className="flex items-center gap-3">
                  <label className="inline-flex items-center gap-2 rounded-xl bg-white border border-[#EBD4C5] px-4 py-2 text-xs font-bold text-[#30251D] hover:bg-[#FFF8F5] cursor-pointer shadow-xs">
                    <Upload size={14} className="text-[#FF9800]" />
                    <span>Pilih Berkas</span>
                    <input
                      type="file"
                      className="hidden"
                      onChange={(e) => {
                        if (e.target.files && e.target.files[0]) {
                          setFile(e.target.files[0])
                        }
                      }}
                    />
                  </label>
                  <span className="text-xs text-[#687385] truncate max-w-xs">
                    {file ? file.name : 'Belum ada berkas dipilih'}
                  </span>
                </div>
              </div>

              <button
                type="submit"
                disabled={submitting || !file}
                className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-5 py-2.5 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors disabled:opacity-50"
              >
                <Upload size={14} />
                <span>{isSubmitted ? 'Kirim Ulang Projek' : 'Kumpulkan Projek'}</span>
              </button>
            </form>
          )}
        </div>

        {/* Confirm Modal */}
        <ConfirmDialog
          isOpen={confirmOpen}
          title={isSubmitted ? 'Kirim Ulang Projek?' : 'Kumpulkan Projek Sekarang?'}
          message="Pastikan berkas hasil projek yang kamu unggah sudah lengkap dan siap dinilai oleh gurumu."
          confirmText={submitting ? 'Mengirim...' : 'Ya, Kumpulkan'}
          cancelText="Batal"
          variant="primary"
          onConfirm={handleConfirmSubmit}
          onClose={() => setConfirmOpen(false)}
        />
      </div>
    </StudentLayout>
  )
}
