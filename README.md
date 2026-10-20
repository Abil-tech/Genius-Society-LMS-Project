<div align="center">

<img src="https://readme-typing-svg.demolab.com?font=Fira+Code&weight=700&size=34&pause=1000&color=29B6F6&center=true&vCenter=true&width=520&height=60&lines=Genius+Society" alt="Genius Society" />

**Platform E-Learning sekolah: materi, tugas, assessment, projek, nilai, jadwal, dan ranking dalam satu sistem.**

<br>

<a href="https://github.com/USERNAME/genius-society"><img src="https://img.shields.io/badge/GITHUB-181717?style=for-the-badge&logo=github&logoColor=white" alt="GitHub" /></a>
<img src="https://komarev.com/ghpvc/?username=USERNAME-genius-society&label=VIEWS&color=0e75b6&style=for-the-badge" alt="Views" />

</div>

<br>

<hr>

<div align="center">

## Tech Stack & Tools

<br>

<img src="https://img.shields.io/badge/VITE-646CFF?style=for-the-badge&logo=vite&logoColor=white" alt="Vite" />
<img src="https://img.shields.io/badge/REACT-20232A?style=for-the-badge&logo=react&logoColor=61DAFB" alt="React" />
<img src="https://img.shields.io/badge/TYPESCRIPT-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript" />
<img src="https://img.shields.io/badge/TAILWIND_CSS-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white" alt="Tailwind CSS" />

<img src="https://img.shields.io/badge/REACT_ROUTER-CA4245?style=for-the-badge&logo=reactrouter&logoColor=white" alt="React Router" />
<img src="https://img.shields.io/badge/AXIOS-5A29E4?style=for-the-badge&logo=axios&logoColor=white" alt="Axios" />
<img src="https://img.shields.io/badge/LUCIDE-F56565?style=for-the-badge&logo=lucide&logoColor=white" alt="Lucide" />

<img src="https://img.shields.io/badge/GO-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
<img src="https://img.shields.io/badge/GIN-008ECF?style=for-the-badge&logo=gin&logoColor=white" alt="Gin" />
<img src="https://img.shields.io/badge/MONGODB-47A248?style=for-the-badge&logo=mongodb&logoColor=white" alt="MongoDB" />
<img src="https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT" />

<img src="https://img.shields.io/badge/BACKBLAZE_B2-E21E29?style=for-the-badge&logo=backblaze&logoColor=white" alt="Backblaze B2" />
<img src="https://img.shields.io/badge/VERCEL-000000?style=for-the-badge&logo=vercel&logoColor=white" alt="Vercel" />
<img src="https://img.shields.io/badge/DOCKER-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
<img src="https://img.shields.io/badge/GIT-F05032?style=for-the-badge&logo=git&logoColor=white" alt="Git" />

</div>

<br>

<hr>

<div align="center">

## Fitur Utama

</div>

| Modul | Deskripsi |
|---|---|
| **Materi** | File (Backblaze B2) dan link, dipublikasikan ke beberapa kelas sekaligus |
| **Tugas** | Deadline, file pendukung, submit, resubmit opsional, nilai dan feedback |
| **Assessment** | Pilihan ganda (auto-grading) dan essay (manual), timer, attempt, bobot soal |
| **Projek** | Instruksi, deadline, submission file, nilai dan feedback |
| **Sistem Nilai** | Tugas + Assessment Harian 60%, UTS 20%, UAS 20% (bobot dapat dikonfigurasi) |
| **Jadwal** | Guru, kelas, ruangan, dan wali kelas diatur oleh Kurikulum |
| **Notifikasi** | Tugas/assessment/materi baru, deadline, nilai keluar, pengumuman |
| **Ranking** | Peringkat berdasarkan nilai akhir, mengikuti permission |

<br>

<hr>

<div align="center">

## Role & Akses

</div>

| Role | Cakupan |
|---|---|
| **Super Admin** | Administrasi sistem dan konfigurasi. Tanpa akses akademik |
| **Admin** | Data siswa, guru, kelas, mata pelajaran, tahun ajaran, akun |
| **Guru** | Materi, tugas, assessment, projek, penilaian. Walas adalah atribut tambahan Guru |
| **Murid** | Belajar, mengumpulkan tugas/projek, assessment, nilai, ranking |
| **Kurikulum** | Jadwal, penugasan guru, kelas, wali kelas |
| **Kepala Sekolah** | Read-only penuh, dienforce di backend |

<br>

<hr>

<div align="center">

## Keamanan

</div>

- Password di-hash dengan bcrypt
- JWT disimpan di HTTP-only cookie, tidak di localStorage; `Secure` pada production
- Authorization (role + permission) divalidasi di middleware backend
- Validasi input, CORS, rate limiting pada endpoint sensitif
- Validasi MIME type dan ukuran file, object key aman
- Seluruh credential berasal dari environment variable

<br>

<hr>

<div align="center">

## Getting Started

</div>

```bash
git clone https://github.com/USERNAME/genius-society.git
cd genius-society

# Backend
cd backend
cp .env.example .env
go run ./cmd/...

# Frontend
cd ../frontend
cp .env.example .env
npm install
npm run dev
```

### Environment Variables

| Area | Variable |
|---|---|
| Frontend | `VITE_API_URL` |
| Backend | `APP_ENV` `PORT` `MONGODB_URI` `MONGODB_DATABASE` `JWT_SECRET` `COOKIE_DOMAIN` `CORS_ORIGIN` |
| Backblaze B2 | `B2_ACCOUNT_ID` `B2_APPLICATION_KEY_ID` `B2_APPLICATION_KEY` `B2_BUCKET_NAME` `B2_S3_ENDPOINT` `B2_PUBLIC_URL` |

### Seed Data (development only)

1 Super Admin, 1 Admin, 2 Guru, 1 Kurikulum, 1 Kepala Sekolah, 10 Murid, 2 Kelas, 3 Mata Pelajaran.
Kredensial development didokumentasikan terpisah dan tidak boleh dipakai di production.

<br>

<hr>

<div align="center">

## Struktur Project

</div>

```
genius-society/
├── frontend/   # Vite + React + TypeScript
│   └── src/{components,layouts,pages,routes,services,hooks,contexts,types,utils}
└── backend/    # Go + Gin
    ├── cmd/
    ├── internal/{config,handler,service,repository,middleware,model,dto,storage}
    ├── routes/
    └── migrations/
```

<br>

<hr>

<div align="center">

**Genius Society** · Built with React, Go, and MongoDB

</div>
