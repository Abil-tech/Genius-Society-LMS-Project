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
  getAssignmentDetail,
  submitAssignment,
  resubmitAssignment,
  getAssignmentFileDownloadURL,
} from '../../services/student/studentAssignment.service'
import type { Assignment } from '../../types/student'

export default function StudentAssignmentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const toast = useToast()

  const [assignment, setAssignment] = useState<Assignment | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [file, setFile] = useState<File | null>(null)
  const [note, setNote] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [downloadingIdx, setDownloadingIdx] = useState<number | null>(null)

  const fetchData = async () => {
    if (!id) return
    try {
      setLoading(true)
      setError(null)
      const res = await getAssignmentDetail(id)
      setAssignment(res)
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal memuat tugas'
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
      const { url } = await getAssignmentFileDownloadURL(id, index)
      window.open(url, '_blank')
    } catch {
      toast.error('Gagal Mengunduh', 'Berkas soal tidak dapat diunduh.')
    } finally {
      setDownloadingIdx(null)
    }
  }

  const handleConfirmSubmit = async () => {
    if (!id) return
    try {
      setSubmitting(true)
      const formData = new FormData()
      if (file) formData.append('file', file)
      if (note) formData.append('note', note)

      const isResubmit = Boolean(assignment?.submission?.submitted)
      if (isResubmit) {
        await resubmitAssignment(id, formData)
        toast.success('Berhasil Mengirim Ulang', 'Revisi tugas telah tersimpan.')
      } else {
        await submitAssignment(id, formData)
        toast.success('Tugas Berhasil Dikumpulkan', 'Jawabanmu telah diterima oleh guru.')
      }

      setFile(null)
      setNote('')
      await fetchData()
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal mengumpulkan tugas'
      toast.error('Pengumpulan Gagal', msg)
    } finally {
      setSubmitting(false)
      setConfirmOpen(false)
    }
  }

  if (error) {
    return (
      <StudentLayout>
        <ErrorState title="Tugas Tidak Ditemukan" message={error} onRetry={fetchData} />
      </StudentLayout>
    )
  }

  if (loading || !assignment) {
    return (
      <StudentLayout>
        <div className="animate-pulse space-y-4 max-w-4xl">
          <div className="h-6 w-32 bg-gray-200 rounded-lg" />
          <div className="h-44 bg-white border border-[#EBD4C5] rounded-2xl" />
        </div>
      </StudentLayout>
    )
  }

  const isSubmitted = assignment.submission?.submitted
  const canSubmit = !isSubmitted || assignment.allowResubmit

  return (
    <StudentLayout>
      <div className="space-y-6 max-w-4xl">
        <Link
          to="/student/assignments"
          className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
        >
          <ArrowLeft size={16} />
          <span>Kembali ke Daftar Tugas</span>
        </Link>

        {/* Task Info */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 shadow-xs">
          <div className="flex flex-wrap items-center gap-2 text-xs">
            <span className="rounded-lg bg-[#FEEADC] px-2.5 py-1 font-extrabold text-[#FF9800]">
              {assignment.subjectName}
            </span>
            <span className="flex items-center gap-1 text-[#687385]">
              <Clock size={13} className="text-[#FF9800]" />
              <span>
                Deadline:{' '}
                {new Date(assignment.deadline).toLocaleDateString('id-ID', {
                  day: 'numeric',
                  month: 'long',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </span>
            </span>
            <span className="text-[#687385]">&bull; Guru: {assignment.teacherName || '-'}</span>
          </div>

          <h1 className="mt-3 text-lg sm:text-2xl font-black text-[#30251D]">
            {assignment.title}
          </h1>

          <div className="mt-4 border-t border-[#EBD4C5]/60 pt-4 text-xs sm:text-sm text-[#30251D] whitespace-pre-line leading-relaxed">
            {assignment.description || 'Tidak ada deskripsi rincian tugas.'}
          </div>

          {/* Question Attachment */}
          {assignment.attachments && assignment.attachments.length > 0 && (
            <div className="mt-5 border-t border-[#EBD4C5]/60 pt-4">
              <span className="text-xs font-bold uppercase tracking-wider text-[#687385]">
                Lampiran Lembar Soal / Template:
              </span>
              <div className="mt-2 space-y-2">
                {assignment.attachments.map((att, idx) => (
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
            <span>Status Pengumpulan Jawaban</span>
            {isSubmitted && (
              <span className="flex items-center gap-1 text-xs font-bold text-emerald-600">
                <CheckCircle2 size={15} />
                <span>Sudah Dikumpulkan (Versi {assignment.submission?.version})</span>
              </span>
            )}
          </h2>

          {isSubmitted ? (
            <div className="mt-4 space-y-4">
              {assignment.submission?.score !== undefined && (
                <div className="flex items-center justify-between p-4 rounded-xl border border-emerald-200 bg-emerald-50/60">
                  <div>
                    <span className="text-xs font-bold text-emerald-800">
                      Nilai Diberikan Guru:
                    </span>
                    {assignment.submission?.feedback && (
                      <p className="mt-1 text-xs text-emerald-700 italic">
                        "{assignment.submission.feedback}"
                      </p>
                    )}
                  </div>
                  <span className="text-2xl font-black text-emerald-700">
                    {assignment.submission.score}
                  </span>
                </div>
              )}

              <div className="text-xs text-[#687385] space-y-1">
                <p>
                  Waktu Submit:{' '}
                  {assignment.submission?.submittedAt
                    ? new Date(assignment.submission.submittedAt).toLocaleString('id-ID')
                    : '-'}
                  {assignment.submission?.isLate && (
                    <span className="ml-2 font-bold text-rose-600">(Terlambat)</span>
                  )}
                </p>
                {assignment.submission?.note && (
                  <p>Catatan Siswa: {assignment.submission.note}</p>
                )}
                {assignment.submission?.files && assignment.submission.files.length > 0 && (
                  <p>
                    Berkas Dikirim:{' '}
                    {assignment.submission.files.map((f) => f.name).join(', ')}
                  </p>
                )}
              </div>
            </div>
          ) : (
            <p className="mt-4 text-xs text-[#687385]">
              Kamu belum mengumpulkan tugas ini. Silakan unggah berkas jawaban di bawah.
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
                {isSubmitted ? 'Unggah Pengumpulan Ulang (Resubmit)' : 'Kirim Jawaban Tugas'}
              </h3>

              <div>
                <label className="block text-xs font-bold text-[#30251D] mb-1">
                  Pilih Berkas Jawaban (PDF, DOCX, XLSX, ZIP, Gambar maks 50MB)
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

              <div>
                <label className="block text-xs font-bold text-[#30251D] mb-1">
                  Catatan untuk Guru (Opsional)
                </label>
                <textarea
                  rows={2}
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  placeholder="Tambahkan catatan jika ada penjelasan khusus..."
                  className="w-full rounded-xl border border-[#EBD4C5] px-3.5 py-2 text-xs text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none"
                />
              </div>

              <button
                type="submit"
                disabled={submitting || (!file && !note)}
                className="inline-flex items-center gap-2 rounded-xl bg-[#FF9800] px-5 py-2.5 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] transition-colors disabled:opacity-50"
              >
                <Upload size={14} />
                <span>{isSubmitted ? 'Kirim Ulang Jawaban' : 'Kumpulkan Tugas'}</span>
              </button>
            </form>
          )}
        </div>

        {/* Confirm Modal */}
        <ConfirmDialog
          isOpen={confirmOpen}
          title={isSubmitted ? 'Kirim Ulang Jawaban?' : 'Kumpulkan Tugas Sekarang?'}
          message="Pastikan berkas jawaban yang kamu pilih sudah benar sebelum mengirimkan ke gurumu."
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
