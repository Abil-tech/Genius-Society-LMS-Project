import { useEffect, useState, useMemo } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  ArrowLeft,
  Search,
  RefreshCw,
  Download,
  Award,
  Clock,
  FileText,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import Modal from '../../components/common/Modal'
import StatusBadge from '../../components/common/StatusBadge'
import ErrorState from '../../components/common/ErrorState'
import { useToast } from '../../components/common/Toast'
import {
  getProjectById,
  getProjectSubmissions,
  gradeProjectSubmission,
  getDownloadPresignedUrl,
} from '../../services/teacherService'
import type { Project, ProjectSubmission } from '../../types/teacher'

export default function TeacherProjectSubmissionsPage() {
  const { id } = useParams<{ id: string }>()
  const toast = useToast()

  const [project, setProject] = useState<Project | null>(null)
  const [submissions, setSubmissions] = useState<ProjectSubmission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')

  // Modal Grading State
  const [activeSubmission, setActiveSubmission] = useState<ProjectSubmission | null>(null)
  const [score, setScore] = useState<string>('')
  const [feedback, setFeedback] = useState<string>('')
  const [isGrading, setIsGrading] = useState(false)
  const [downloadingKey, setDownloadingKey] = useState<string | null>(null)

  const fetchData = async () => {
    if (!id) return
    setLoading(true)
    setError(null)
    try {
      const [pRes, subRes] = await Promise.allSettled([
        getProjectById(id),
        getProjectSubmissions(id),
      ])

      if (pRes.status === 'fulfilled' && pRes.value.success && pRes.value.data) {
        setProject(pRes.value.data)
      }
      if (subRes.status === 'fulfilled' && subRes.value.success && subRes.value.data) {
        setSubmissions(subRes.value.data)
      } else if (subRes.status === 'fulfilled' && !subRes.value.success) {
        setError(subRes.value.error?.message || 'Gagal memuat pengumpulan projek.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat data.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [id])

  const openGradingModal = (sub: ProjectSubmission) => {
    setActiveSubmission(sub)
    setScore(sub.score !== undefined && sub.score !== null ? String(sub.score) : '')
    setFeedback(sub.feedback || '')
  }

  const handleGradeSubmit = async () => {
    if (!id || !activeSubmission) return
    const numScore = parseFloat(score)
    if (isNaN(numScore) || numScore < 0 || numScore > 100) {
      toast.warning('Nilai Tidak Valid', 'Masukkan nilai angka antara 0 hingga 100.')
      return
    }

    setIsGrading(true)
    try {
      const res = await gradeProjectSubmission(id, activeSubmission.id, {
        score: numScore,
        feedback: feedback.trim(),
      })

      if (res.success) {
        toast.success('Nilai Disimpan', `Nilai ${numScore} berhasil diberikan ke ${activeSubmission.studentName}.`)
        setActiveSubmission(null)
        fetchData()
      } else {
        toast.error('Gagal Menyimpan Nilai', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal mengirim nilai ke server.'
      toast.error('Gagal Menyimpan', msg)
    } finally {
      setIsGrading(false)
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
        toast.error('Gagal Mengunduh', res.error?.message || 'URL berkas tidak valid.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal menghasilkan URL unduhan.'
      toast.error('Gagal Mengunduh', msg)
    } finally {
      setDownloadingKey(null)
    }
  }

  const filteredSubmissions = useMemo(() => {
    if (!search.trim()) return submissions
    const q = search.toLowerCase()
    return submissions.filter(
      (s) =>
        s.studentName.toLowerCase().includes(q) ||
        (s.studentNis && s.studentNis.toLowerCase().includes(q))
    )
  }, [submissions, search])

  const columns: Column<ProjectSubmission>[] = [
    {
      key: 'student',
      header: 'Nama Siswa',
      render: (s) => (
        <div>
          <p className="font-bold text-[#30251D]">{s.studentName}</p>
          {s.studentNis && <p className="text-[11px] text-[#687385]">NIS: {s.studentNis}</p>}
        </div>
      ),
    },
    {
      key: 'version',
      header: 'Versi',
      align: 'center',
      render: (s) => (
        <span className="rounded-md bg-[#FFF8F5] border border-[#EBD4C5] px-2 py-0.5 text-xs font-bold text-[#30251D]">
          v{s.version}
        </span>
      ),
    },
    {
      key: 'submittedAt',
      header: 'Waktu Kumpul',
      render: (s) => (
        <span className="inline-flex items-center gap-1 text-xs text-[#687385]">
          <Clock size={12} className="text-[#FF9800]" />
          <span>
            {new Date(s.submittedAt).toLocaleString('id-ID', {
              dateStyle: 'medium',
              timeStyle: 'short',
            })}
          </span>
        </span>
      ),
    },
    {
      key: 'files',
      header: 'Berkas Karya',
      render: (s) => (
        <div className="flex flex-col gap-1">
          {s.files && s.files.length > 0 ? (
            s.files.map((f, idx) => (
              <button
                key={idx}
                type="button"
                onClick={() => handleDownload(f.key, f.name)}
                disabled={downloadingKey === f.key}
                className="inline-flex items-center gap-1.5 text-xs font-semibold text-[#FF9800] hover:underline"
              >
                <Download size={12} />
                <span className="truncate max-w-[140px]">{f.name}</span>
              </button>
            ))
          ) : (
            <span className="text-gray-400 text-xs">Tanpa berkas</span>
          )}
        </div>
      ),
    },
    {
      key: 'status',
      header: 'Status & Nilai',
      render: (s) => {
        const isGraded = s.score !== undefined && s.score !== null
        return isGraded ? (
          <div>
            <StatusBadge status="graded" label={`Nilai: ${s.score}`} />
            {s.feedback && (
              <p className="mt-1 text-[11px] text-[#687385] line-clamp-1 italic">
                "{s.feedback}"
              </p>
            )}
          </div>
        ) : (
          <StatusBadge status="pending" label="Belum Dinilai" />
        )
      },
    },
    {
      key: 'action',
      header: 'Aksi',
      align: 'right',
      render: (s) => (
        <button
          type="button"
          onClick={() => openGradingModal(s)}
          className="inline-flex items-center gap-1.5 rounded-xl bg-[#FFF8F5] border border-[#EBD4C5] px-3 py-1.5 text-xs font-bold text-[#FF9800] hover:bg-[#FEEADC] transition-colors"
        >
          <Award size={14} />
          <span>{s.score !== undefined && s.score !== null ? 'Ubah Nilai' : 'Beri Nilai'}</span>
        </button>
      ),
    },
  ]

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Back Link */}
        <div className="flex items-center justify-between">
          <Link
            to={project ? `/teacher/projects/${project.id}` : '/teacher/projects'}
            className="inline-flex items-center gap-2 text-xs font-bold text-[#FF9800] hover:underline"
          >
            <ArrowLeft size={16} />
            <span>Kembali ke Detail Projek</span>
          </Link>

          <button
            type="button"
            onClick={fetchData}
            disabled={loading}
            className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3 py-1.5 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
          >
            <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
            <span>Segarkan</span>
          </button>
        </div>

        {/* Header */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-white p-5 shadow-xs">
          <div>
            <span className="rounded-full bg-[#FFF8F5] border border-[#EBD4C5] px-2.5 py-0.5 text-[10px] font-extrabold uppercase text-[#FF9800]">
              Pengumpulan Projek Siswa
            </span>
            <h1 className="mt-1 text-xl sm:text-2xl font-extrabold text-[#30251D]">
              {project?.title || 'Pengumpulan Projek'}
            </h1>
            <p className="text-xs text-[#687385]">
              {project?.subjectName ? `${project.subjectName} • ` : ''}
              Total {submissions.length} siswa telah mengumpulkan projek ini
            </p>
          </div>
        </div>

        {/* Search */}
        <div className="relative max-w-md">
          <Search
            size={16}
            className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
          />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari siswa yang mengumpulkan..."
            className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
          />
        </div>

        {/* Content */}
        {error ? (
          <ErrorState title="Gagal Memuat Pengumpulan" message={error} onRetry={fetchData} />
        ) : (
          <DataTable
            columns={columns}
            data={filteredSubmissions}
            loading={loading}
            emptyMessage={search ? 'Siswa tidak ditemukan' : 'Belum ada siswa yang mengumpulkan projek ini'}
            emptySubMessage={
              search ? `Tidak ada hasil pencarian untuk "${search}"` : undefined
            }
          />
        )}

        {/* Modal Penilaian Projek */}
        <Modal
          isOpen={Boolean(activeSubmission)}
          onClose={() => setActiveSubmission(null)}
          title={`Penilaian Projek: ${activeSubmission?.studentName}`}
          description={`Pengumpulan versi ke-${activeSubmission?.version} • Dikirim pada ${
            activeSubmission
              ? new Date(activeSubmission.submittedAt).toLocaleString('id-ID', {
                  dateStyle: 'medium',
                  timeStyle: 'short',
                })
              : ''
          }`}
          footer={
            <>
              <button
                type="button"
                onClick={() => setActiveSubmission(null)}
                className="w-full sm:w-auto rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 text-xs font-semibold text-[#30251D] hover:bg-[#FFF8F5]"
              >
                Batal
              </button>
              <button
                type="button"
                onClick={handleGradeSubmit}
                disabled={isGrading}
                className="w-full sm:w-auto inline-flex items-center justify-center gap-1.5 rounded-xl bg-[#FF9800] px-5 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] disabled:opacity-50"
              >
                {isGrading ? (
                  <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                ) : (
                  <Award size={15} />
                )}
                <span>Simpan Nilai</span>
              </button>
            </>
          }
        >
          <div className="space-y-4">
            {/* Berkas Projek */}
            <div>
              <p className="text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Berkas Hasil Karya Siswa
              </p>
              {activeSubmission?.files && activeSubmission.files.length > 0 ? (
                <div className="space-y-2">
                  {activeSubmission.files.map((file, i) => (
                    <div
                      key={i}
                      className="flex items-center justify-between rounded-xl border border-[#EBD4C5] bg-[#FFF8F5]/60 p-3"
                    >
                      <div className="flex items-center gap-2.5 min-w-0 pr-2">
                        <FileText size={18} className="text-[#FF9800] shrink-0" />
                        <span className="truncate text-xs font-bold text-[#30251D]">
                          {file.name}
                        </span>
                      </div>
                      <button
                        type="button"
                        onClick={() => handleDownload(file.key, file.name)}
                        disabled={downloadingKey === file.key}
                        className="inline-flex items-center gap-1.5 rounded-lg bg-[#FF9800] px-3 py-1 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] shrink-0"
                      >
                        <Download size={12} />
                        <span>Unduh</span>
                      </button>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-xs text-gray-500 italic p-3 bg-stone-50 rounded-xl">
                  Siswa tidak menyertakan berkas karya.
                </p>
              )}
            </div>

            {/* Input Nilai */}
            <div>
              <label htmlFor="score" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Nilai Projek (Skala 0 - 100) <span className="text-rose-500">*</span>
              </label>
              <input
                id="score"
                type="number"
                min="0"
                max="100"
                step="0.5"
                value={score}
                onChange={(e) => setScore(e.target.value)}
                placeholder="Contoh: 90"
                className="w-full sm:w-48 rounded-xl border border-[#EBD4C5] px-3.5 py-2.5 text-base font-extrabold text-[#30251D] focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
              />
            </div>

            {/* Feedback / Catatan Guru */}
            <div>
              <label htmlFor="feedback" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Umpan Balik (Feedback) & Rubrik Projek
              </label>
              <textarea
                id="feedback"
                rows={3}
                value={feedback}
                onChange={(e) => setFeedback(e.target.value)}
                placeholder="Berikan apresiasi dan catatan konstruktif pada hasil karya siswa..."
                className="w-full rounded-xl border border-[#EBD4C5] px-3.5 py-2.5 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
              />
            </div>
          </div>
        </Modal>
      </div>
    </TeacherLayout>
  )
}
