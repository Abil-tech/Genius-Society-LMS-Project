import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  ArrowLeft,
  Clock,
  Download,
  Edit,
  Trash2,
  FileText,
  DoorOpen,
  Users,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import ConfirmDialog from '../../components/common/ConfirmDialog'
import ErrorState from '../../components/common/ErrorState'
import StatusBadge from '../../components/common/StatusBadge'
import { useToast } from '../../components/common/Toast'
import {
  getAssignmentById,
  deleteAssignment,
  getDownloadPresignedUrl,
} from '../../services/teacherService'
import type { Assignment } from '../../types/teacher'

export default function TeacherAssignmentDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const toast = useToast()

  const [assignment, setAssignment] = useState<Assignment | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const [downloadingKey, setDownloadingKey] = useState<string | null>(null)

  const fetchDetail = async () => {
    if (!id) return
    setLoading(true)
    setError(null)
    try {
      const res = await getAssignmentById(id)
      if (res.success && res.data) {
        setAssignment(res.data)
      } else {
        setError(res.error?.message || 'Tugas tidak ditemukan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal memuat detail tugas.'
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
      const res = await deleteAssignment(id)
      if (res.success) {
        toast.success(
          'Tugas Berhasil Diproses',
          res.data?.archived ? 'Tugas telah diarsipkan.' : 'Tugas telah dihapus.'
        )
        navigate('/teacher/assignments')
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
    setDownloadingKey(key)
    try {
      const res = await getDownloadPresignedUrl(key)
      if (res.success && res.data?.url) {
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
      setDownloadingKey(null)
    }
  }

  const isPast = assignment ? new Date(assignment.deadline) < new Date() : false
  const hasSubmissions = (assignment?.totalSubmissions ?? 0) > 0

  return (
    <TeacherLayout>
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Back Link */}
        <div className="flex items-center justify-between">
          <Link
            to="/teacher/assignments"
            className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
          >
            <ArrowLeft size={16} />
            <span>Kembali ke Daftar Tugas</span>
          </Link>

          {assignment && (
            <div className="flex items-center gap-2">
              <Link
                to={`/teacher/assignments/${assignment.id}/submissions`}
                className="inline-flex items-center gap-1.5 rounded-xl bg-[#FF9800] px-3.5 py-2 text-xs font-bold text-white hover:bg-[#E68A00] shadow-xs transition-colors"
              >
                <Users size={14} />
                <span>Lihat Pengumpulan ({assignment.totalSubmissions ?? 0})</span>
              </Link>
              <Link
                to={`/teacher/assignments/${assignment.id}/edit`}
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
          <ErrorState title="Tugas Tidak Ditemukan" message={error} onRetry={fetchDetail} />
        ) : loading ? (
          <div className="space-y-4 rounded-2xl border border-[#EBD4C5] bg-white p-6 animate-pulse">
            <div className="h-8 bg-gray-200 rounded-lg w-2/3" />
            <div className="h-20 bg-gray-200 rounded-lg" />
          </div>
        ) : assignment ? (
          <div className="space-y-6">
            <div className="rounded-2xl border border-[#EBD4C5] bg-white p-6 sm:p-8 shadow-xs">
              <div className="flex flex-wrap items-center gap-2 mb-3">
                <span className="rounded-full bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-0.5 text-[11px] font-extrabold uppercase text-[#FF9800]">
                  {assignment.subjectName || 'Mata Pelajaran'}
                </span>
                <StatusBadge
                  status={isPast ? 'ended' : 'active'}
                  label={isPast ? 'Tenggat Berakhir' : 'Aktif'}
                />
                <span className="inline-flex items-center gap-1 text-xs text-[#687385]">
                  <Clock size={13} className="text-[#FF9800]" />
                  <span>
                    Deadline:{' '}
                    {new Date(assignment.deadline).toLocaleString('id-ID', {
                      dateStyle: 'medium',
                      timeStyle: 'short',
                    })}
                  </span>
                </span>
              </div>

              <h1 className="text-xl sm:text-2xl font-extrabold text-[#30251D]">
                {assignment.title}
              </h1>

              {assignment.description && (
                <div className="mt-4 prose prose-stone max-w-none text-xs sm:text-sm text-[#30251D] leading-relaxed whitespace-pre-wrap">
                  {assignment.description}
                </div>
              )}

              {/* Status Ringkasan Pengumpulan */}
              <div className="mt-6 grid grid-cols-2 sm:grid-cols-3 gap-3 pt-5 border-t border-[#EBD4C5]/60">
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3 text-center">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">
                    Total Dikumpulkan
                  </span>
                  <p className="text-lg font-extrabold text-[#30251D]">
                    {assignment.totalSubmissions ?? 0}
                  </p>
                </div>
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3 text-center">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">
                    Telah Dinilai
                  </span>
                  <p className="text-lg font-extrabold text-emerald-700">
                    {assignment.gradedSubmissions ?? 0}
                  </p>
                </div>
                <div className="rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3 text-center col-span-2 sm:col-span-1">
                  <span className="text-[10px] font-bold uppercase text-[#687385]">
                    Aturan Resubmit
                  </span>
                  <p className="text-sm font-bold text-[#30251D] mt-1">
                    {assignment.allowResubmit ? 'Diizinkan' : 'Tidak Diizinkan'}
                  </p>
                </div>
              </div>

              {/* Kelas Dituju */}
              <div className="mt-6 pt-5 border-t border-[#EBD4C5]/60">
                <p className="text-xs font-bold uppercase tracking-wider text-[#687385] mb-2 flex items-center gap-1.5">
                  <DoorOpen size={14} className="text-[#FF9800]" />
                  <span>Diberikan Kepada Kelas:</span>
                </p>
                <div className="flex flex-wrap gap-1.5">
                  {assignment.classNames && assignment.classNames.length > 0 ? (
                    assignment.classNames.map((c, i) => (
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

              {/* Lampiran Soal */}
              {assignment.attachments && assignment.attachments.length > 0 && (
                <div className="mt-6 pt-5 border-t border-[#EBD4C5]/60">
                  <p className="text-xs font-bold uppercase tracking-wider text-[#687385] mb-2">
                    Berkas Lampiran Soal / Panduan
                  </p>
                  <div className="space-y-2">
                    {assignment.attachments.map((att, idx) => (
                      <div
                        key={idx}
                        className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/50 p-3.5"
                      >
                        <div className="flex items-center gap-3">
                          <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-white border border-[#EBD4C5] text-[#FF9800]">
                            <FileText size={18} />
                          </div>
                          <div className="min-w-0">
                            <p className="font-bold text-xs sm:text-sm text-[#30251D] truncate">
                              {att.name}
                            </p>
                            <p className="text-[11px] text-[#687385]">
                              {(att.size / (1024 * 1024)).toFixed(2)} MB • {att.mime}
                            </p>
                          </div>
                        </div>

                        <button
                          type="button"
                          onClick={() => handleDownload(att.key, att.name)}
                          disabled={downloadingKey === att.key}
                          className="inline-flex items-center justify-center gap-1.5 rounded-xl bg-[#FF9800] px-3.5 py-1.5 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] disabled:opacity-50 transition-colors"
                        >
                          <Download size={13} />
                          <span>{downloadingKey === att.key ? 'Mengunduh...' : 'Unduh'}</span>
                        </button>
                      </div>
                    ))}
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
          title={hasSubmissions ? 'Arsipkan Tugas' : 'Hapus Tugas'}
          message={`Apakah Anda yakin ingin menghapus atau mengarsipkan tugas "${assignment?.title}"?`}
          hasSubmissions={hasSubmissions}
          archiveNotice="Tugas ini sudah memiliki pengumpulan dari siswa. Demi integritas akademik, data akan diarsipkan."
          isLoading={isDeleting}
        />
      </div>
    </TeacherLayout>
  )
}
