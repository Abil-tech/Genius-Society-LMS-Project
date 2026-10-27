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
import TeacherDashboard from '../pages/TeacherDashboard'
import StudentDashboard from '../pages/StudentDashboard'
import CurriculumDashboard from '../pages/CurriculumDashboard'
import PrincipalDashboard from '../pages/PrincipalDashboard'
import ClassesManagementPage from '../pages/admin/ClassesManagementPage'
import AccountsManagementPage from '../pages/admin/AccountsManagementPage'
import AcademicYearsPage from '../pages/admin/AcademicYearsPage'

export default function AppRoutes() {
  return (
    <Routes>
      {/* Public routes */}
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/admin/login" element={<AdminLoginPage />} />

      {/* Protected routes by role */}
      <Route element={<ProtectedRoute allowedRoles={['super_admin']} />}>
        <Route path="/super-admin/*" element={<SuperAdminDashboard />} />
      </Route>

      <Route element={<ProtectedRoute allowedRoles={['admin', 'super_admin']} />}>
        <Route path="/admin" element={<AdminDashboard />} />
        <Route path="/admin/students" element={<StudentsManagementPage />} />
        <Route path="/admin/teachers" element={<TeachersManagementPage />} />
        <Route path="/admin/subjects" element={<SubjectsManagementPage />} />
        <Route path="/admin/classes" element={<ClassesManagementPage />} />
<Route path="/admin/accounts" element={<AccountsManagementPage />} />
<Route path="/admin/academic-year" element={<AcademicYearsPage />} />
      </Route>

      <Route element={<ProtectedRoute allowedRoles={['teacher']} />}>
        <Route path="/teacher/*" element={<TeacherDashboard />} />
      </Route>

      <Route element={<ProtectedRoute allowedRoles={['student']} />}>
        <Route path="/student/*" element={<StudentDashboard />} />
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
