import { useEffect, useState, useMemo } from 'react'
import {
  Search,
  RefreshCw,
  Calculator,
  Save,
} from 'lucide-react'
import TeacherLayout from '../../components/teacher/TeacherLayout'
import DataTable, { type Column } from '../../components/common/DataTable'
import ErrorState from '../../components/common/ErrorState'
import Modal from '../../components/common/Modal'
import { useToast } from '../../components/common/Toast'
import {
  getTeacherGrades,
  getTeacherClasses,
  getTeacherSubjects,
  updateStudentGrade,
} from '../../services/teacherService'
import type {
  StudentGradeRow,
  TeacherClass,
  TeacherSubject,
} from '../../types/teacher'

export default function TeacherGradesPage() {
  const toast = useToast()
  const [grades, setGrades] = useState<StudentGradeRow[]>([])
  const [classes, setClasses] = useState<TeacherClass[]>([])
  const [subjects, setSubjects] = useState<TeacherSubject[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Filters
  const [selectedClass, setSelectedClass] = useState('')
  const [selectedSubject, setSelectedSubject] = useState('')
  const [search, setSearch] = useState('')

  // Edit Single Grade Modal
  const [editingStudent, setEditingStudent] = useState<StudentGradeRow | null>(null)
  const [editScoreType, setEditScoreType] = useState<'uts' | 'uas'>('uts')
  const [editScoreValue, setEditScoreValue] = useState('')
  const [savingGrade, setSavingGrade] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    setError(null)
    try {
      const [clsRes, subRes] = await Promise.all([
        getTeacherClasses(),
        getTeacherSubjects(),
      ])

      if (clsRes.success && clsRes.data) {
        setClasses(clsRes.data)
        if (!selectedClass && clsRes.data.length > 0) {
          setSelectedClass(clsRes.data[0].id)
        }
      }

      if (subRes.success && subRes.data) {
        setSubjects(subRes.data)
        if (!selectedSubject && subRes.data.length > 0) {
          setSelectedSubject(subRes.data[0].id)
        }
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal memuat filter kelas.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const fetchGrades = async () => {
    if (!selectedClass || !selectedSubject) return
    setLoading(true)
    setError(null)
    try {
      const res = await getTeacherGrades({
        classId: selectedClass,
        subjectId: selectedSubject,
      })
      if (res.success && res.data) {
        setGrades(res.data)
      } else {
        setError(res.error?.message || 'Gagal memuat rekap nilai siswa.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Terjadi kendala saat memuat data nilai.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (selectedClass && selectedSubject) {
      fetchGrades()
    }
  }, [selectedClass, selectedSubject])

  const filteredGrades = useMemo(() => {
    if (!search.trim()) return grades
    const q = search.toLowerCase()
    return grades.filter(
      (g) =>
        g.studentName.toLowerCase().includes(q) ||
        (g.studentNis && g.studentNis.toLowerCase().includes(q))
    )
  }, [grades, search])

  // Formula 60 / 20 / 20 Calculation
  const calculateFinal = (daily = 0, uts = 0, uas = 0) => {
    return Math.round((daily * 0.6 + uts * 0.2 + uas * 0.2) * 10) / 10
  }

  const openEditModal = (row: StudentGradeRow, type: 'uts' | 'uas') => {
    setEditingStudent(row)
    setEditScoreType(type)
    setEditScoreValue(
      type === 'uts'
        ? String(row.utsScore ?? '')
        : String(row.uasScore ?? '')
    )
  }

  const handleSaveGrade = async () => {
    if (!editingStudent || !selectedSubject || !selectedClass) return
    const num = parseFloat(editScoreValue)
    if (isNaN(num) || num < 0 || num > 100) {
      toast.warning('Nilai Tidak Valid', 'Masukkan angka antara 0 hingga 100.')
      return
    }

    setSavingGrade(true)
    try {
      const res = await updateStudentGrade({
        studentId: editingStudent.studentId,
        subjectId: selectedSubject,
        classId: selectedClass,
        type: editScoreType,
        score: num,
      })

      if (res.success) {
        toast.success(
          'Nilai Diperbarui',
          `Nilai ${editScoreType.toUpperCase()} siswa ${editingStudent.studentName} disimpan.`
        )
        setEditingStudent(null)
        fetchGrades()
      } else {
        toast.error('Gagal Menyimpan', res.error?.message || 'Terjadi kesalahan.')
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Gagal menyimpan nilai.'
      toast.error('Gagal Menyimpan', msg)
    } finally {
      setSavingGrade(false)
    }
  }

  const columns: Column<StudentGradeRow>[] = [
    {
      key: 'no',
      header: 'No',
      className: 'w-10',
      render: (_, i) => <span className="font-medium text-gray-500">{i + 1}</span>,
    },
    {
      key: 'student',
      header: 'Siswa',
      render: (g) => (
        <div>
          <p className="font-bold text-[#30251D]">{g.studentName}</p>
          <p className="text-[11px] text-[#687385]">NIS: {g.studentNis || '-'}</p>
        </div>
      ),
    },
    {
      key: 'assignmentScore',
      header: 'Tugas',
      align: 'center',
      render: (g) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {g.assignmentScore !== undefined && g.assignmentScore !== null ? g.assignmentScore : '-'}
        </span>
      ),
    },
    {
      key: 'assessmentScore',
      header: 'Kuis/Harian',
      align: 'center',
      render: (g) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {g.assessmentScore !== undefined && g.assessmentScore !== null ? g.assessmentScore : '-'}
        </span>
      ),
    },
    {
      key: 'projectScore',
      header: 'Projek',
      align: 'center',
      render: (g) => (
        <span className="font-semibold text-xs text-[#30251D]">
          {g.projectScore !== undefined && g.projectScore !== null ? g.projectScore : '-'}
        </span>
      ),
    },
    {
      key: 'dailyAvg',
      header: 'Harian (60%)',
      align: 'center',
      render: (g) => {
        const dAvg = g.dailyAvg ?? 0
        return (
          <span className="rounded-md bg-[#FFF8F5] border border-[#EBD4C5] px-2 py-0.5 text-xs font-bold text-[#FF9800]">
            {dAvg > 0 ? dAvg : '-'}
          </span>
        )
      },
    },
    {
      key: 'utsScore',
      header: 'UTS (20%)',
      align: 'center',
      render: (g) => (
        <button
          type="button"
          onClick={() => openEditModal(g, 'uts')}
          className="rounded-md border border-[#EBD4C5] hover:border-[#FF9800] bg-white px-2.5 py-0.5 text-xs font-bold text-[#30251D] hover:text-[#FF9800] transition-colors"
          title="Klik untuk ubah nilai UTS"
        >
          {g.utsScore !== undefined && g.utsScore !== null ? g.utsScore : 'Input'}
        </button>
      ),
    },
    {
      key: 'uasScore',
      header: 'UAS (20%)',
      align: 'center',
      render: (g) => (
        <button
          type="button"
          onClick={() => openEditModal(g, 'uas')}
          className="rounded-md border border-[#EBD4C5] hover:border-[#FF9800] bg-white px-2.5 py-0.5 text-xs font-bold text-[#30251D] hover:text-[#FF9800] transition-colors"
          title="Klik untuk ubah nilai UAS"
        >
          {g.uasScore !== undefined && g.uasScore !== null ? g.uasScore : 'Input'}
        </button>
      ),
    },
    {
      key: 'finalScore',
      header: 'Nilai Akhir',
      align: 'right',
      render: (g) => {
        const final =
          g.finalScore ??
          calculateFinal(g.dailyAvg, g.utsScore, g.uasScore)
        const isPassed = final >= 75

        return (
          <div className="text-right">
            <span
              className={`text-sm font-black ${
                isPassed ? 'text-emerald-700' : 'text-rose-600'
              }`}
            >
              {final > 0 ? final : '-'}
            </span>
            <span className="block text-[10px] text-[#687385]">
              {isPassed ? 'Tuntas' : 'Perlu Remidi'}
            </span>
          </div>
        )
      },
    },
  ]

  return (
    <TeacherLayout>
      <div className="space-y-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-xl sm:text-2xl font-extrabold tracking-tight text-[#30251D]">
              Penilaian & Rekap Nilai Siswa
            </h1>
            <p className="mt-1 text-xs sm:text-sm text-[#687385]">
              Rekapitulasi nilai otomatis berdasarkan formula bobot resmi: Harian (60%), UTS (20%), dan UAS (20%).
            </p>
          </div>

          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={fetchGrades}
              disabled={loading}
              className="inline-flex items-center gap-1.5 rounded-xl border border-[#EBD4C5] bg-white px-3.5 py-2 text-xs font-bold text-[#30251D] shadow-xs hover:bg-[#FFF8F5]"
            >
              <RefreshCw size={14} className={loading ? 'animate-spin text-[#FF9800]' : ''} />
              <span>Hitung & Segarkan</span>
            </button>
          </div>
        </div>

        {/* Formula Explainer Card */}
        <div className="rounded-2xl border border-[#FF9800]/30 bg-gradient-to-r from-[#FFF8F5] via-white to-[#FEEADC]/40 p-4 shadow-xs">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-[#FF9800] text-white">
                <Calculator size={18} />
              </div>
              <div className="text-xs text-[#30251D]">
                <span className="font-bold block">Formula Penilaian Akademik Genius Society:</span>
                <span className="text-[#687385]">
                  Nilai Akhir = (Rata-rata Tugas + Assessment + Projek) × 60% + (UTS × 20%) + (UAS × 20%)
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Filter Controls: Kelas & Mata Pelajaran */}
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="relative flex-1">
            <Search
              size={16}
              className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400"
            />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Cari siswa dalam rekap..."
              className="w-full rounded-xl border border-[#EBD4C5] bg-white py-2 pl-9 pr-4 text-xs sm:text-sm text-[#30251D] placeholder-gray-400 focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
            />
          </div>

          <div className="flex items-center gap-2">
            <select
              value={selectedClass}
              onChange={(e) => setSelectedClass(e.target.value)}
              aria-label="Pilih kelas"
              className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs sm:text-sm font-bold text-[#30251D] focus:border-[#FF9800] focus:outline-none"
            >
              {classes.map((c) => (
                <option key={c.id} value={c.id}>
                  Kelas: {c.name}
                </option>
              ))}
            </select>

            <select
              value={selectedSubject}
              onChange={(e) => setSelectedSubject(e.target.value)}
              aria-label="Pilih mata pelajaran"
              className="rounded-xl border border-[#EBD4C5] bg-white px-3 py-2 text-xs sm:text-sm font-bold text-[#30251D] focus:border-[#FF9800] focus:outline-none"
            >
              {subjects.map((s) => (
                <option key={s.id} value={s.id}>
                  Mapel: {s.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Content Table */}
        {error ? (
          <ErrorState title="Gagal Memuat Rekap Nilai" message={error} onRetry={fetchGrades} />
        ) : (
          <DataTable
            columns={columns}
            data={filteredGrades}
            loading={loading}
            emptyMessage={
              !selectedClass || !selectedSubject
                ? 'Pilih kelas dan mata pelajaran untuk melihat nilai'
                : 'Belum ada data nilai di kelas dan mata pelajaran ini'
            }
            emptySubMessage="Nilai akan terisi otomatis seiring tugas, assessment, dan projek diselesaikan dan dinilai."
          />
        )}

        {/* Modal Input Nilai UTS / UAS */}
        <Modal
          isOpen={Boolean(editingStudent)}
          onClose={() => setEditingStudent(null)}
          title={`Input Nilai ${editScoreType.toUpperCase()}: ${editingStudent?.studentName}`}
          description={`Kelas: ${classes.find((c) => c.id === selectedClass)?.name || ''} • Mapel: ${
            subjects.find((s) => s.id === selectedSubject)?.name || ''
          }`}
          size="sm"
          footer={
            <>
              <button
                type="button"
                onClick={() => setEditingStudent(null)}
                className="w-full sm:w-auto rounded-xl border border-[#EBD4C5] bg-white px-4 py-2 text-xs font-semibold text-[#30251D] hover:bg-[#FFF8F5]"
              >
                Batal
              </button>
              <button
                type="button"
                onClick={handleSaveGrade}
                disabled={savingGrade}
                className="w-full sm:w-auto inline-flex items-center justify-center gap-1.5 rounded-xl bg-[#FF9800] px-5 py-2 text-xs font-bold text-white shadow-xs hover:bg-[#E68A00] disabled:opacity-50"
              >
                {savingGrade ? (
                  <span className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                ) : (
                  <Save size={15} />
                )}
                <span>Simpan Nilai</span>
              </button>
            </>
          }
        >
          <div className="space-y-4">
            <div>
              <label htmlFor="editScore" className="block text-xs font-bold uppercase tracking-wider text-[#30251D] mb-1.5">
                Nilai {editScoreType.toUpperCase()} (Skala 0 - 100)
              </label>
              <input
                id="editScore"
                type="number"
                min="0"
                max="100"
                step="0.5"
                value={editScoreValue}
                onChange={(e) => setEditScoreValue(e.target.value)}
                placeholder="Contoh: 88"
                className="w-full rounded-xl border border-[#EBD4C5] px-3.5 py-2.5 text-base font-extrabold text-[#30251D] focus:border-[#FF9800] focus:outline-none focus:ring-1 focus:ring-[#FF9800]"
              />
              <p className="mt-1 text-xs text-[#687385]">
                Nilai ini akan langsung memperbarui bobot 20% pada perhitungan nilai akhir semester siswa.
              </p>
            </div>
          </div>
        </Modal>
      </div>
    </TeacherLayout>
  )
}
