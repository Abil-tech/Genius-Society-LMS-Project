import { api, type ApiResponse } from './api'
import type {
  TeacherDashboardData,
  TeacherClass,
  TeacherStudent,
  TeacherSubject,
  Material,
  Assignment,
  AssignmentSubmission,
  Assessment,
  AssessmentAttempt,
  Project,
  ProjectSubmission,
  StudentGradeRow,
  GradeItem,
  TeacherScheduleItem,
  WalasStudentDetail,
  WalasGradeRekapRow,
  CreateMaterialPayload,
  CreateAssignmentPayload,
  CreateAssessmentPayload,
  CreateProjectPayload,
} from '../types/teacher'

export interface FilterParams {
  search?: string
  classId?: string
  subjectId?: string
  page?: number
  limit?: number
}

// 1. Dashboard
export const getTeacherDashboard = async () => {
  const res = await api.get<ApiResponse<TeacherDashboardData>>('/teacher/dashboard')
  return res.data
}

// 2. Classes & Students
export const getTeacherClasses = async () => {
  const res = await api.get<ApiResponse<TeacherClass[]>>('/teacher/classes')
  return res.data
}

export const getClassDetail = async (id: string) => {
  const res = await api.get<ApiResponse<TeacherClass>>(`/teacher/classes/${id}`)
  return res.data
}

export const getClassStudents = async (classId: string) => {
  const res = await api.get<ApiResponse<TeacherStudent[]>>(`/teacher/classes/${classId}/students`)
  return res.data
}

// 3. Subjects
export const getTeacherSubjects = async () => {
  const res = await api.get<ApiResponse<TeacherSubject[]>>('/teacher/subjects')
  return res.data
}

// 4. Materials
export const getMaterials = async (params?: FilterParams) => {
  const res = await api.get<ApiResponse<Material[]>>('/teacher/materials', { params })
  return res.data
}

export const getMaterialById = async (id: string) => {
  const res = await api.get<ApiResponse<Material>>(`/teacher/materials/${id}`)
  return res.data
}

export const createMaterial = async (data: FormData | CreateMaterialPayload) => {
  const isMultipart = data instanceof FormData
  const res = await api.post<ApiResponse<Material>>('/teacher/materials', data, {
    headers: isMultipart ? { 'Content-Type': 'multipart/form-data' } : undefined,
  })
  return res.data
}

export const updateMaterial = async (id: string, data: FormData | Partial<CreateMaterialPayload>) => {
  const isMultipart = data instanceof FormData
  const res = await api.put<ApiResponse<Material>>(`/teacher/materials/${id}`, data, {
    headers: isMultipart ? { 'Content-Type': 'multipart/form-data' } : undefined,
  })
  return res.data
}

export const deleteMaterial = async (id: string) => {
  const res = await api.delete<ApiResponse<{ message: string; archived?: boolean }>>(`/teacher/materials/${id}`)
  return res.data
}

// 5. Assignments
export const getAssignments = async (params?: FilterParams) => {
  const res = await api.get<ApiResponse<Assignment[]>>('/teacher/assignments', { params })
  return res.data
}

export const getAssignmentById = async (id: string) => {
  const res = await api.get<ApiResponse<Assignment>>(`/teacher/assignments/${id}`)
  return res.data
}

export const createAssignment = async (data: FormData | CreateAssignmentPayload) => {
  const isMultipart = data instanceof FormData
  const res = await api.post<ApiResponse<Assignment>>('/teacher/assignments', data, {
    headers: isMultipart ? { 'Content-Type': 'multipart/form-data' } : undefined,
  })
  return res.data
}

export const updateAssignment = async (id: string, data: FormData | Partial<CreateAssignmentPayload>) => {
  const isMultipart = data instanceof FormData
  const res = await api.put<ApiResponse<Assignment>>(`/teacher/assignments/${id}`, data, {
    headers: isMultipart ? { 'Content-Type': 'multipart/form-data' } : undefined,
  })
  return res.data
}

export const deleteAssignment = async (id: string) => {
  const res = await api.delete<ApiResponse<{ message: string; archived?: boolean }>>(`/teacher/assignments/${id}`)
  return res.data
}

export const getAssignmentSubmissions = async (assignmentId: string) => {
  const res = await api.get<ApiResponse<AssignmentSubmission[]>>(`/teacher/assignments/${assignmentId}/submissions`)
  return res.data
}

export const gradeAssignmentSubmission = async (
  assignmentId: string,
  submissionId: string,
  data: { score: number; feedback: string }
) => {
  const res = await api.post<ApiResponse<AssignmentSubmission>>(
    `/teacher/assignments/${assignmentId}/submissions/${submissionId}/grade`,
    data
  )
  return res.data
}

// 6. Assessments
export const getAssessments = async (params?: FilterParams) => {
  const res = await api.get<ApiResponse<Assessment[]>>('/teacher/assessments', { params })
  return res.data
}

export const getAssessmentById = async (id: string) => {
  const res = await api.get<ApiResponse<Assessment>>(`/teacher/assessments/${id}`)
  return res.data
}

export const createAssessment = async (data: CreateAssessmentPayload) => {
  const res = await api.post<ApiResponse<Assessment>>('/teacher/assessments', data)
  return res.data
}

export const updateAssessment = async (id: string, data: Partial<CreateAssessmentPayload>) => {
  const res = await api.put<ApiResponse<Assessment>>(`/teacher/assessments/${id}`, data)
  return res.data
}

export const deleteAssessment = async (id: string) => {
  const res = await api.delete<ApiResponse<{ message: string; archived?: boolean }>>(`/teacher/assessments/${id}`)
  return res.data
}

export const getAssessmentAttempts = async (assessmentId: string) => {
  const res = await api.get<ApiResponse<AssessmentAttempt[]>>(`/teacher/assessments/${assessmentId}/attempts`)
  return res.data
}

export const getAssessmentAttemptDetail = async (assessmentId: string, attemptId: string) => {
  const res = await api.get<ApiResponse<AssessmentAttempt>>(`/teacher/assessments/${assessmentId}/attempts/${attemptId}`)
  return res.data
}

export const gradeEssayAnswer = async (
  assessmentId: string,
  attemptId: string,
  questionId: string,
  data: { score: number; feedback?: string }
) => {
  const res = await api.post<ApiResponse<AssessmentAttempt>>(
    `/teacher/assessments/${assessmentId}/attempts/${attemptId}/questions/${questionId}/grade`,
    data
  )
  return res.data
}

// 7. Projects
export const getProjects = async (params?: FilterParams) => {
  const res = await api.get<ApiResponse<Project[]>>('/teacher/projects', { params })
  return res.data
}

export const getProjectById = async (id: string) => {
  const res = await api.get<ApiResponse<Project>>(`/teacher/projects/${id}`)
  return res.data
}

export const createProject = async (data: FormData | CreateProjectPayload) => {
  const isMultipart = data instanceof FormData
  const res = await api.post<ApiResponse<Project>>('/teacher/projects', data, {
    headers: isMultipart ? { 'Content-Type': 'multipart/form-data' } : undefined,
  })
  return res.data
}

export const updateProject = async (id: string, data: FormData | Partial<CreateProjectPayload>) => {
  const isMultipart = data instanceof FormData
  const res = await api.put<ApiResponse<Project>>(`/teacher/projects/${id}`, data, {
    headers: isMultipart ? { 'Content-Type': 'multipart/form-data' } : undefined,
  })
  return res.data
}

export const deleteProject = async (id: string) => {
  const res = await api.delete<ApiResponse<{ message: string; archived?: boolean }>>(`/teacher/projects/${id}`)
  return res.data
}

export const getProjectSubmissions = async (projectId: string) => {
  const res = await api.get<ApiResponse<ProjectSubmission[]>>(`/teacher/projects/${projectId}/submissions`)
  return res.data
}

export const gradeProjectSubmission = async (
  projectId: string,
  submissionId: string,
  data: { score: number; feedback: string }
) => {
  const res = await api.post<ApiResponse<ProjectSubmission>>(
    `/teacher/projects/${projectId}/submissions/${submissionId}/grade`,
    data
  )
  return res.data
}

// 8. Grades
export const getTeacherGrades = async (params?: { classId?: string; subjectId?: string }) => {
  const res = await api.get<ApiResponse<StudentGradeRow[]>>('/teacher/grades', { params })
  return res.data
}

export const updateStudentGrade = async (data: {
  studentId: string
  subjectId: string
  classId: string
  type: string
  score: number
}) => {
  const res = await api.put<ApiResponse<GradeItem>>('/teacher/grades', data)
  return res.data
}

// 9. Schedule
export const getTeacherSchedule = async () => {
  const res = await api.get<ApiResponse<TeacherScheduleItem[]>>('/teacher/schedule')
  return res.data
}

// 10. Walas
export const getWalasStudents = async () => {
  const res = await api.get<ApiResponse<WalasStudentDetail[]>>('/teacher/walas/students')
  return res.data
}

export const getWalasGrades = async () => {
  const res = await api.get<ApiResponse<WalasGradeRekapRow[]>>('/teacher/walas/grades')
  return res.data
}

// 11. Presigned file download URL
export const getDownloadPresignedUrl = async (key: string) => {
  const res = await api.get<ApiResponse<{ url: string }>>('/teacher/files/download', {
    params: { key },
  })
  return res.data
}
