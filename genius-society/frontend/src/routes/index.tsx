import { Navigate, Route, Routes } from 'react-router-dom'
import LandingPage from '../pages/LandingPage'
import LoginPage from '../pages/LoginPage'
import AdminLoginPage from '../pages/AdminLoginPage'
import ProtectedRoute from '../components/ProtectedRoute'
import SuperAdminDashboard from '../pages/SuperAdminDashboard'
import AdminDashboard from '../pages/admin/AdminDashboardPage'
import StudentsManagementPage from '../pages/admin/StudentsManagementPage'
import TeachersManagementPage from '../pages/admin/TeachersManagementPage'
import SubjectsManagementPage from '../pages/admin/SubjectsManagementPage'
import CurriculumDashboard from '../pages/CurriculumDashboard'
import PrincipalDashboard from '../pages/PrincipalDashboard'
import ClassesManagementPage from '../pages/admin/ClassesManagementPage'
import AccountsManagementPage from '../pages/admin/AccountsManagementPage'
import AcademicYearsPage from '../pages/admin/AcademicYearsPage'

// Teacher Pages
import TeacherDashboardPage from '../pages/teacher/TeacherDashboardPage'
import TeacherClassesPage from '../pages/teacher/TeacherClassesPage'
import TeacherClassDetailPage from '../pages/teacher/TeacherClassDetailPage'
import TeacherMaterialsPage from '../pages/teacher/TeacherMaterialsPage'
import TeacherMaterialDetailPage from '../pages/teacher/TeacherMaterialDetailPage'
import TeacherMaterialFormPage from '../pages/teacher/TeacherMaterialFormPage'
import TeacherAssignmentsPage from '../pages/teacher/TeacherAssignmentsPage'
import TeacherAssignmentDetailPage from '../pages/teacher/TeacherAssignmentDetailPage'
import TeacherAssignmentFormPage from '../pages/teacher/TeacherAssignmentFormPage'
import TeacherAssignmentSubmissionsPage from '../pages/teacher/TeacherAssignmentSubmissionsPage'
import TeacherAssessmentsPage from '../pages/teacher/TeacherAssessmentsPage'
import TeacherAssessmentDetailPage from '../pages/teacher/TeacherAssessmentDetailPage'
import TeacherAssessmentFormPage from '../pages/teacher/TeacherAssessmentFormPage'
import TeacherProjectsPage from '../pages/teacher/TeacherProjectsPage'
import TeacherProjectDetailPage from '../pages/teacher/TeacherProjectDetailPage'
import TeacherProjectFormPage from '../pages/teacher/TeacherProjectFormPage'
import TeacherProjectSubmissionsPage from '../pages/teacher/TeacherProjectSubmissionsPage'
import TeacherGradesPage from '../pages/teacher/TeacherGradesPage'
import TeacherSchedulePage from '../pages/teacher/TeacherSchedulePage'
import TeacherWalasStudentsPage from '../pages/teacher/TeacherWalasStudentsPage'
import TeacherWalasGradesPage from '../pages/teacher/TeacherWalasGradesPage'

// Student Pages
import StudentDashboardPage from '../pages/student/StudentDashboardPage'
import StudentMaterialsPage from '../pages/student/StudentMaterialsPage'
import StudentMaterialDetailPage from '../pages/student/StudentMaterialDetailPage'
import StudentAssignmentsPage from '../pages/student/StudentAssignmentsPage'
import StudentAssignmentDetailPage from '../pages/student/StudentAssignmentDetailPage'
import StudentAssessmentsPage from '../pages/student/StudentAssessmentsPage'
import StudentAssessmentDetailPage from '../pages/student/StudentAssessmentDetailPage'
import StudentAssessmentWorkPage from '../pages/student/StudentAssessmentWorkPage'
import StudentProjectsPage from '../pages/student/StudentProjectsPage'
import StudentProjectDetailPage from '../pages/student/StudentProjectDetailPage'
import StudentGradesPage from '../pages/student/StudentGradesPage'
import StudentSchedulePage from '../pages/student/StudentSchedulePage'
import StudentRankingPage from '../pages/student/StudentRankingPage'

export default function AppRoutes() {
  return (
    <Routes>
      {/* Public routes */}
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/admin/login" element={<AdminLoginPage />} />

      {/* Protected routes: Super Admin */}
      <Route element={<ProtectedRoute allowedRoles={['super_admin']} />}>
        <Route path="/super-admin/*" element={<SuperAdminDashboard />} />
      </Route>

      {/* Protected routes: Admin & Super Admin */}
      <Route element={<ProtectedRoute allowedRoles={['admin', 'super_admin']} />}>
        <Route path="/admin" element={<AdminDashboard />} />
        <Route path="/admin/students" element={<StudentsManagementPage />} />
        <Route path="/admin/teachers" element={<TeachersManagementPage />} />
        <Route path="/admin/subjects" element={<SubjectsManagementPage />} />
        <Route path="/admin/classes" element={<ClassesManagementPage />} />
        <Route path="/admin/accounts" element={<AccountsManagementPage />} />
        <Route path="/admin/academic-year" element={<AcademicYearsPage />} />
      </Route>

      {/* Protected routes: Teacher (Only 'teacher' role allowed) */}
      <Route element={<ProtectedRoute allowedRoles={['teacher']} />}>
        {/* Dashboard */}
        <Route path="/teacher" element={<TeacherDashboardPage />} />

        {/* Kelas */}
        <Route path="/teacher/classes" element={<TeacherClassesPage />} />
        <Route path="/teacher/classes/:id" element={<TeacherClassDetailPage />} />

        {/* Materi */}
        <Route path="/teacher/materials" element={<TeacherMaterialsPage />} />
        <Route path="/teacher/materials/new" element={<TeacherMaterialFormPage />} />
        <Route path="/teacher/materials/:id" element={<TeacherMaterialDetailPage />} />
        <Route path="/teacher/materials/:id/edit" element={<TeacherMaterialFormPage />} />

        {/* Tugas */}
        <Route path="/teacher/assignments" element={<TeacherAssignmentsPage />} />
        <Route path="/teacher/assignments/new" element={<TeacherAssignmentFormPage />} />
        <Route path="/teacher/assignments/:id" element={<TeacherAssignmentDetailPage />} />
        <Route path="/teacher/assignments/:id/edit" element={<TeacherAssignmentFormPage />} />
        <Route
          path="/teacher/assignments/:id/submissions"
          element={<TeacherAssignmentSubmissionsPage />}
        />

        {/* Assessment */}
        <Route path="/teacher/assessments" element={<TeacherAssessmentsPage />} />
        <Route path="/teacher/assessments/new" element={<TeacherAssessmentFormPage />} />
        <Route path="/teacher/assessments/:id" element={<TeacherAssessmentDetailPage />} />
        <Route path="/teacher/assessments/:id/edit" element={<TeacherAssessmentFormPage />} />

        {/* Projek */}
        <Route path="/teacher/projects" element={<TeacherProjectsPage />} />
        <Route path="/teacher/projects/new" element={<TeacherProjectFormPage />} />
        <Route path="/teacher/projects/:id" element={<TeacherProjectDetailPage />} />
        <Route path="/teacher/projects/:id/edit" element={<TeacherProjectFormPage />} />
        <Route
          path="/teacher/projects/:id/submissions"
          element={<TeacherProjectSubmissionsPage />}
        />

        {/* Penilaian & Rekap */}
        <Route path="/teacher/grades" element={<TeacherGradesPage />} />

        {/* Jadwal Mengajar */}
        <Route path="/teacher/schedule" element={<TeacherSchedulePage />} />

        {/* Wali Kelas */}
        <Route
          path="/teacher/walas"
          element={<Navigate to="/teacher/walas/students" replace />}
        />
        <Route path="/teacher/walas/students" element={<TeacherWalasStudentsPage />} />
        <Route path="/teacher/walas/grades" element={<TeacherWalasGradesPage />} />
      </Route>

      {/* Protected routes: Student */}
      <Route element={<ProtectedRoute allowedRoles={['student']} />}>
        <Route path="/student" element={<StudentDashboardPage />} />
        <Route path="/student/materials" element={<StudentMaterialsPage />} />
        <Route path="/student/materials/:id" element={<StudentMaterialDetailPage />} />
        <Route path="/student/assignments" element={<StudentAssignmentsPage />} />
        <Route path="/student/assignments/:id" element={<StudentAssignmentDetailPage />} />
        <Route path="/student/assessments" element={<StudentAssessmentsPage />} />
        <Route path="/student/assessments/:id" element={<StudentAssessmentDetailPage />} />
        <Route path="/student/assessments/:id/work" element={<StudentAssessmentWorkPage />} />
        <Route path="/student/projects" element={<StudentProjectsPage />} />
        <Route path="/student/projects/:id" element={<StudentProjectDetailPage />} />
        <Route path="/student/grades" element={<StudentGradesPage />} />
        <Route path="/student/schedule" element={<StudentSchedulePage />} />
        <Route path="/student/ranking" element={<StudentRankingPage />} />
      </Route>

      <Route element={<ProtectedRoute allowedRoles={['curriculum']} />}>
        <Route path="/curriculum/*" element={<CurriculumDashboard />} />
      </Route>

      <Route element={<ProtectedRoute allowedRoles={['principal']} />}>
        <Route path="/principal/*" element={<PrincipalDashboard />} />
      </Route>

      {/* Fallback */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
