import { useEffect, useState } from 'react'
import AdminLayout from '../../components/admin/AdminLayout'
import { api, type ApiResponse } from '../../services/api'
import { 
  GraduationCap, 
  Users, 
  DoorOpen, 
  UserCheck, 
  AlertTriangle, 
  CheckCircle2, 
  PlusCircle, 
  Calendar,
  Layers,
  ArrowUpRight,
  RefreshCw,
  Clock
} from 'lucide-react'

interface AdminDashboardData {
  academicYear: {
    name: string
    semester: number
    isActive: boolean
  }
  currentDate: string
  kpi: {
    totalStudents: number
    totalTeachers: number
    totalClasses: number
    totalAccounts: number
  }
  classDistribution: Array<{
    classId: string
    className: string
    count: number
  }>
  dataCompleteness: {
    studentsWithoutClass: number
    classesWithoutWalas: number
    teachersWithoutMapel: number
    inactiveAccounts: number
  }
  recentUsers: Array<{
    id: string
    name: string
    username: string
    email: string
    role: string
    createdAt: string
    isActive: boolean
  }>
}

export default function AdminDashboardPage() {
  const [data, setData] = useState<AdminDashboardData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchDashboardData = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await api.get<ApiResponse<AdminDashboardData>>('/admin/dashboard')
      if (res.data.success && res.data.data) {
        setData(res.data.data)
      } else {
        setError(res.data.error?.message || 'Gagal memuat data dashboard')
      }
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message)
      } else {
        setError('Gagal terhubung ke server atau memuat data dashboard.')
      }
      setData(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const handleQuickAction = (actionName: string) => {
    alert(`Aksi Cepat: "${actionName}" akan membuka modal form pembuatan data. Fitur modal detail pada modul berikutnya.`)
  }

  return (
    <AdminLayout>
      <div className="mx-auto max-w-7xl space-y-6">
        {/* ========================================================================= */}
        {/* 1. HEADING: "Dashboard Admin", Tahun Ajaran, Semester Aktif, Tanggal Hari Ini */}
        {/* ========================================================================= */}
        <div className="flex flex-col gap-3 rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-6 shadow-sm sm:flex-row sm:items-center sm:justify-between">
          <div>
            <div className="flex items-center gap-2">
              <span className="flex h-2.5 w-2.5 rounded-full bg-[#FF9800] animate-pulse" />
              <span className="text-[11px] font-bold uppercase tracking-wider text-[#FF9800]">
                Ringkasan Eksekutif
              </span>
            </div>
            <h1 className="mt-1 text-2xl font-black text-[#30251D] sm:text-3xl tracking-tight">
              Dashboard Admin
            </h1>
            <p className="mt-1 text-xs text-gray-500">
              Pantau status administrasi akademik, distribusi kelas, dan kelengkapan data sekolah.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            {/* Badge Tahun Ajaran & Semester */}
            <div className="flex items-center gap-2 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] px-3.5 py-2">
              <Calendar size={15} className="text-[#FF9800]" />
              <div className="text-left">
                <span className="block text-[10px] uppercase font-bold text-gray-400">
                  Tahun Ajaran Aktif
                </span>
                <span className="block text-xs font-bold text-[#30251D]">
                  {data?.academicYear.name || '2024/2025'} • Semester {data?.academicYear.semester || 1}
                </span>
              </div>
            </div>

            {/* Badge Tanggal Hari Ini */}
            <div className="flex items-center gap-2 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] px-3.5 py-2">
              <Clock size={15} className="text-[#FF9800]" />
              <div className="text-left">
                <span className="block text-[10px] uppercase font-bold text-gray-400">
                  Tanggal Hari Ini
                </span>
                <span className="block text-xs font-bold text-[#30251D]">
                  {data?.currentDate || '20 Oktober 2026'}
                </span>
              </div>
            </div>

            <button
              onClick={fetchDashboardData}
              title="Refresh Data"
              className="flex h-10 w-10 items-center justify-center rounded-xl border border-[#EBD4C5] bg-white text-gray-600 hover:bg-[#FFF8F5] hover:text-[#FF9800] transition-colors"
            >
              <RefreshCw size={15} className={loading ? 'animate-spin' : ''} />
            </button>
          </div>
        </div>

        {error && (
          <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-xs text-red-600">
            {error}
          </div>
        )}

        {/* ========================================================================= */}
        {/* 2. KPI CARDS (4): Total Siswa, Total Guru, Total Kelas, Total Akun */}
        {/* ========================================================================= */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {/* Card 1: Total Siswa */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-sm transition-all hover:border-[#FF9800]/50 hover:shadow-md">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold uppercase tracking-wider text-gray-500">
                Total Siswa
              </span>
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#FEEADC] text-[#FF9800]">
                <GraduationCap size={20} />
              </span>
            </div>
            <div className="mt-4">
              <span className="text-3xl font-black text-[#30251D]">
                {data?.kpi.totalStudents ?? 0}
              </span>
              <span className="ml-2 text-xs font-medium text-gray-400">Siswa Terdaftar</span>
            </div>
            <p className="mt-2 text-[11px] text-gray-500">
              Status aktif di seluruh angkatan
            </p>
          </div>

          {/* Card 2: Total Guru */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-sm transition-all hover:border-[#FF9800]/50 hover:shadow-md">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold uppercase tracking-wider text-gray-500">
                Total Guru
              </span>
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-orange-50 text-[#FF9800]">
                <Users size={20} />
              </span>
            </div>
            <div className="mt-4">
              <span className="text-3xl font-black text-[#30251D]">
                {data?.kpi.totalTeachers ?? 0}
              </span>
              <span className="ml-2 text-xs font-medium text-gray-400">Tenaga Pendidik</span>
            </div>
            <p className="mt-2 text-[11px] text-gray-500">
              Guru tetap & pembimbing wali kelas
            </p>
          </div>

          {/* Card 3: Total Kelas */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-sm transition-all hover:border-[#FF9800]/50 hover:shadow-md">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold uppercase tracking-wider text-gray-500">
                Total Kelas
              </span>
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#FEEADC] text-[#FF9800]">
                <DoorOpen size={20} />
              </span>
            </div>
            <div className="mt-4">
              <span className="text-3xl font-black text-[#30251D]">
                {data?.kpi.totalClasses ?? 0}
              </span>
              <span className="ml-2 text-xs font-medium text-gray-400">Rombel Aktif</span>
            </div>
            <p className="mt-2 text-[11px] text-gray-500">
              Tingkat X, XI, dan XII
            </p>
          </div>

          {/* Card 4: Total Akun */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-sm transition-all hover:border-[#FF9800]/50 hover:shadow-md">
            <div className="flex items-center justify-between">
              <span className="text-xs font-bold uppercase tracking-wider text-gray-500">
                Total Akun
              </span>
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-orange-50 text-[#FF9800]">
                <UserCheck size={20} />
              </span>
            </div>
            <div className="mt-4">
              <span className="text-3xl font-black text-[#30251D]">
                {data?.kpi.totalAccounts ?? 0}
              </span>
              <span className="ml-2 text-xs font-medium text-gray-400">Pengguna Sistem</span>
            </div>
            <p className="mt-2 text-[11px] text-gray-500">
              Termasuk admin, guru, murid & staf
            </p>
          </div>
        </div>

        {/* ========================================================================= */}
        {/* 3. AKSI CEPAT: Tambah Siswa, Tambah Guru, Tambah Kelas, Tambah Akun */}
        {/* ========================================================================= */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-5 shadow-sm">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h2 className="text-sm font-bold text-[#30251D]">Aksi Cepat Administrasi</h2>
              <p className="text-[11px] text-gray-500">Pintasan praktis untuk input dan pembuatan entitas baru</p>
            </div>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
            <button
              type="button"
              onClick={() => handleQuickAction('Tambah Siswa')}
              className="flex items-center justify-center gap-2 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] transition-all hover:bg-[#FF9800] hover:text-white hover:border-[#FF9800] hover:shadow-sm"
            >
              <PlusCircle size={15} />
              <span>Tambah Siswa</span>
            </button>
            <button
              type="button"
              onClick={() => handleQuickAction('Tambah Guru')}
              className="flex items-center justify-center gap-2 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] transition-all hover:bg-[#FF9800] hover:text-white hover:border-[#FF9800] hover:shadow-sm"
            >
              <PlusCircle size={15} />
              <span>Tambah Guru</span>
            </button>
            <button
              type="button"
              onClick={() => handleQuickAction('Tambah Kelas')}
              className="flex items-center justify-center gap-2 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] transition-all hover:bg-[#FF9800] hover:text-white hover:border-[#FF9800] hover:shadow-sm"
            >
              <PlusCircle size={15} />
              <span>Tambah Kelas</span>
            </button>
            <button
              type="button"
              onClick={() => handleQuickAction('Tambah Akun')}
              className="flex items-center justify-center gap-2 rounded-xl border border-[#EBD4C5] bg-[#FFF8F5] p-3 text-xs font-bold text-[#30251D] transition-all hover:bg-[#FF9800] hover:text-white hover:border-[#FF9800] hover:shadow-sm"
            >
              <PlusCircle size={15} />
              <span>Tambah Akun</span>
            </button>
          </div>
        </div>

        {/* ========================================================================= */}
        {/* 4. KONTEN UTAMA & PANEL PENDUKUNG (GRID 2 KOLOM) */}
        {/* ========================================================================= */}
        <div className="grid gap-6 lg:grid-cols-3">
          {/* Konten Utama (2 Kolom): Distribusi Siswa per Kelas */}
          <div className="lg:col-span-2 rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-6 shadow-sm">
            <div className="flex items-center justify-between mb-5">
              <div>
                <div className="flex items-center gap-2">
                  <Layers size={16} className="text-[#FF9800]" />
                  <h2 className="text-base font-bold text-[#30251D]">Distribusi Siswa per Kelas</h2>
                </div>
                <p className="mt-1 text-xs text-gray-500">
                  Kapasitas dan perbandingan populasi murid per rombongan belajar
                </p>
              </div>
            </div>

            <div className="space-y-4">
              {data?.classDistribution && data.classDistribution.length > 0 ? (
                data.classDistribution.map((item) => {
                  const maxCap = 36 // kapasitas ideal standar kelas
                  const percentage = Math.min(Math.round((item.count / maxCap) * 100), 100)
                  return (
                    <div key={item.classId} className="rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5]/60 p-4">
                      <div className="flex items-center justify-between text-xs font-bold text-[#30251D]">
                        <span className="flex items-center gap-2">
                          <span className="h-2 w-2 rounded-full bg-[#FF9800]" />
                          {item.className}
                        </span>
                        <span>
                          {item.count} <span className="text-gray-400 font-normal">/ {maxCap} Siswa</span>
                        </span>
                      </div>
                      {/* Bar Progress Indicator */}
                      <div className="mt-2.5 h-2 w-full overflow-hidden rounded-full bg-[#EBD4C5]/50">
                        <div
                          className="h-full rounded-full bg-[#FF9800] transition-all duration-500"
                          style={{ width: `${Math.max(percentage, 4)}%` }}
                        />
                      </div>
                      <div className="mt-1.5 flex justify-between text-[10px] text-gray-400">
                        <span>Terisi {percentage}%</span>
                        <span>Sisa Kuota: {Math.max(maxCap - item.count, 0)} kursi</span>
                      </div>
                    </div>
                  )
                })
              ) : (
                <div className="py-8 text-center text-xs text-gray-400">
                  Belum ada data distribusi kelas.
                </div>
              )}
            </div>
          </div>

          {/* Panel Pendukung (1 Kolom): Kelengkapan Data (Actionable Items) */}
          <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-6 shadow-sm flex flex-col justify-between">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <AlertTriangle size={16} className="text-[#FF9800]" />
                <h2 className="text-base font-bold text-[#30251D]">Kelengkapan Data</h2>
              </div>
              <p className="text-xs text-gray-500 mb-5">
                Perlu ditindaklanjuti untuk menjaga integritas data akademik
              </p>

              <div className="space-y-3">
                {/* 1. Siswa belum punya kelas */}
                <div className="flex items-center justify-between rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5] p-3.5">
                  <div>
                    <p className="text-xs font-bold text-[#30251D]">Siswa Belum Punya Kelas</p>
                    <p className="text-[10px] text-gray-500">Perlu penempatan rombel</p>
                  </div>
                  <span className={`rounded-full px-2.5 py-1 text-xs font-bold ${
                    (data?.dataCompleteness.studentsWithoutClass ?? 0) > 0
                      ? 'bg-red-100 text-red-700'
                      : 'bg-emerald-100 text-emerald-700'
                  }`}>
                    {data?.dataCompleteness.studentsWithoutClass ?? 0}
                  </span>
                </div>

                {/* 2. Kelas tanpa wali kelas */}
                <div className="flex items-center justify-between rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5] p-3.5">
                  <div>
                    <p className="text-xs font-bold text-[#30251D]">Kelas Tanpa Wali Kelas</p>
                    <p className="text-[10px] text-gray-500">Perlu penetapan Walas</p>
                  </div>
                  <span className={`rounded-full px-2.5 py-1 text-xs font-bold ${
                    (data?.dataCompleteness.classesWithoutWalas ?? 0) > 0
                      ? 'bg-amber-100 text-amber-800'
                      : 'bg-emerald-100 text-emerald-700'
                  }`}>
                    {data?.dataCompleteness.classesWithoutWalas ?? 0}
                  </span>
                </div>

                {/* 3. Guru tanpa mata pelajaran */}
                <div className="flex items-center justify-between rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5] p-3.5">
                  <div>
                    <p className="text-xs font-bold text-[#30251D]">Guru Tanpa Mata Pelajaran</p>
                    <p className="text-[10px] text-gray-500">Perlu penugasan mapel</p>
                  </div>
                  <span className={`rounded-full px-2.5 py-1 text-xs font-bold ${
                    (data?.dataCompleteness.teachersWithoutMapel ?? 0) > 0
                      ? 'bg-amber-100 text-amber-800'
                      : 'bg-emerald-100 text-emerald-700'
                  }`}>
                    {data?.dataCompleteness.teachersWithoutMapel ?? 0}
                  </span>
                </div>

                {/* 4. Akun nonaktif */}
                <div className="flex items-center justify-between rounded-xl border border-[#EBD4C5]/60 bg-[#FFF8F5] p-3.5">
                  <div>
                    <p className="text-xs font-bold text-[#30251D]">Akun Nonaktif</p>
                    <p className="text-[10px] text-gray-500">Perlu aktivasi atau review</p>
                  </div>
                  <span className={`rounded-full px-2.5 py-1 text-xs font-bold ${
                    (data?.dataCompleteness.inactiveAccounts ?? 0) > 0
                      ? 'bg-red-100 text-red-700'
                      : 'bg-emerald-100 text-emerald-700'
                  }`}>
                    {data?.dataCompleteness.inactiveAccounts ?? 0}
                  </span>
                </div>
              </div>
            </div>

            <div className="mt-5 rounded-xl border border-[#EBD4C5] bg-[#FEEADC]/40 p-3 flex items-center gap-2 text-xs text-[#30251D]">
              <CheckCircle2 size={16} className="text-[#FF9800] shrink-0" />
              <span>Semua indikator kelengkapan divalidasi langsung dari database.</span>
            </div>
          </div>
        </div>

        {/* ========================================================================= */}
        {/* 5. TABEL: Akun atau Siswa yang Baru Ditambahkan */}
        {/* ========================================================================= */}
        <div className="rounded-2xl border border-[#EBD4C5] bg-[#FFFFFF] p-6 shadow-sm">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-5">
            <div>
              <h2 className="text-base font-bold text-[#30251D]">Akun & Siswa Terbaru</h2>
              <p className="text-xs text-gray-500">Daftar pengguna dan siswa yang baru ditambahkan ke sistem</p>
            </div>
            <button
              type="button"
              onClick={() => handleQuickAction('Lihat Semua Pengguna')}
              className="inline-flex items-center gap-1.5 text-xs font-bold text-[#FF9800] hover:underline"
            >
              <span>Kelola Seluruh Akun</span>
              <ArrowUpRight size={14} />
            </button>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="border-b border-[#EBD4C5] bg-[#FFF8F5] text-gray-500 uppercase tracking-wider text-[10px]">
                <tr>
                  <th className="py-3 px-4 font-bold">Nama Pengguna</th>
                  <th className="py-3 px-4 font-bold">Username</th>
                  <th className="py-3 px-4 font-bold">Email</th>
                  <th className="py-3 px-4 font-bold">Role</th>
                  <th className="py-3 px-4 font-bold">Waktu Dibuat</th>
                  <th className="py-3 px-4 font-bold">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#EBD4C5]/60">
                {data?.recentUsers && data.recentUsers.length > 0 ? (
                  data.recentUsers.map((u) => (
                    <tr key={u.id} className="hover:bg-[#FFF8F5]/50 transition-colors">
                      <td className="py-3.5 px-4 font-bold text-[#30251D]">
                        {u.name}
                      </td>
                      <td className="py-3.5 px-4 font-mono text-gray-600">
                        @{u.username}
                      </td>
                      <td className="py-3.5 px-4 text-gray-500">
                        {u.email}
                      </td>
                      <td className="py-3.5 px-4">
                        <span className="inline-block rounded-md bg-[#FEEADC] px-2 py-0.5 text-[10px] font-bold capitalize text-[#FF9800]">
                          {u.role.replace('_', ' ')}
                        </span>
                      </td>
                      <td className="py-3.5 px-4 text-gray-500">
                        {u.createdAt}
                      </td>
                      <td className="py-3.5 px-4">
                        <span className={`inline-flex items-center gap-1 text-[11px] font-semibold ${
                          u.isActive ? 'text-emerald-600' : 'text-red-500'
                        }`}>
                          <span className={`h-1.5 w-1.5 rounded-full ${
                            u.isActive ? 'bg-emerald-500' : 'bg-red-500'
                          }`} />
                          {u.isActive ? 'Aktif' : 'Nonaktif'}
                        </span>
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} className="py-6 text-center text-gray-400">
                      Tidak ada pengguna yang baru ditambahkan.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </AdminLayout>
  )
}
