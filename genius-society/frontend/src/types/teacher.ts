// Comprehensive types for Teacher role API responses and requests

export interface FileRef {
  key: string
  name: string
  size: number
  mime: string
  url?: string
}

export interface TeacherDashboardKPI {
  classesCount: number
  studentsCount: number
  pendingGradesCount: number
  activeAssignmentsCount: number
  activeAssessmentsCount: number
  activeProjectsCount: number
}

export interface TeacherScheduleItem {
  id: string
  day: number // 1=Senin, ..., 7=Minggu
  startTime: string // "07:30"
  endTime: string   // "09:00"
  subject: string
  className: string
  room?: string
  classId?: string
  subjectId?: string
}

export interface TeacherSubmissionItem {
  id: string
  studentId: string
  studentName: string
  studentNis?: string
  type: 'assignment' | 'project'
  title: string
  submittedAt: string
  graded: boolean
  score?: number
  assignmentId?: string
  projectId?: string
}

export interface WalasClassSummary {
  id: string
  name: string
  level: string
  studentCount: number
}

export interface TeacherDashboardData {
  kpi: TeacherDashboardKPI
  todaySchedule: TeacherScheduleItem[]
  recentSubmissions: TeacherSubmissionItem[]
  isWalas: boolean
  walasClass?: WalasClassSummary
}

export interface TeacherClass {
  id: string
  code: string
  name: string
  educationLevel: string
  department?: string
  academicYearId?: string
  studentCount?: number
  subjects?: Array<{ id: string; name: string }>
}

export interface TeacherStudent {
  id: string
  name: string
  nis: string
  nisn?: string
  gender?: string
  className?: string
  classId?: string
  email?: string
  phone?: string
  status?: string
}

export interface TeacherSubject {
  id: string
  code: string
  name: string
  classIds?: string[]
}

// Material
export interface Material {
  id: string
  title: string
  description?: string
  subjectId: string
  subjectName?: string
  classIds: string[]
  classNames?: string[]
  file?: FileRef
  link?: string
  publishAt: string
  createdAt: string
  updatedAt: string
}

export interface CreateMaterialPayload {
  title: string
  description?: string
  subjectId: string
  classIds: string[]
  link?: string
  publishAt?: string
  file?: File
}

// Assignment
export interface Assignment {
  id: string
  title: string
  description?: string
  subjectId: string
  subjectName?: string
  classIds: string[]
  classNames?: string[]
  deadline: string
  attachments: FileRef[]
  allowResubmit: boolean
  totalSubmissions?: number
  gradedSubmissions?: number
  createdAt: string
  updatedAt: string
}

export interface CreateAssignmentPayload {
  title: string
  description?: string
  subjectId: string
  classIds: string[]
  deadline: string
  allowResubmit: boolean
  attachments?: File[]
}

export interface AssignmentSubmission {
  id: string
  assignmentId: string
  studentId: string
  studentName: string
  studentNis?: string
  files: FileRef[]
  note?: string
  submittedAt: string
  version: number
  score?: number
  feedback?: string
  gradedAt?: string
}

// Assessment
export type AssessmentCategory = 'daily' | 'uts' | 'uas'
export type AssessmentStatus = 'draft' | 'published' | 'closed'
export type QuestionType = 'mcq' | 'essay'

export interface AssessmentQuestion {
  id?: string
  type: QuestionType
  text: string
  options?: string[]
  correctOption?: number // 0-indexed
  weight: number
  order: number
}

export interface Assessment {
  id: string
  title: string
  description?: string
  durationMinutes: number
  subjectId: string
  subjectName?: string
  classIds: string[]
  classNames?: string[]
  startAt?: string
  endAt?: string
  showResult: boolean
  category: AssessmentCategory
  status: AssessmentStatus
  questions?: AssessmentQuestion[]
  totalQuestions?: number
  totalWeight?: number
  totalAttempts?: number
  gradedAttempts?: number
  createdAt: string
  updatedAt: string
}

export interface CreateAssessmentPayload {
  title: string
  description?: string
  durationMinutes: number
  subjectId: string
  classIds: string[]
  startAt?: string
  endAt?: string
  showResult: boolean
  category: AssessmentCategory
  status: AssessmentStatus
  questions: AssessmentQuestion[]
}

export interface AssessmentAttempt {
  id: string
  assessmentId: string
  studentId: string
  studentName: string
  studentNis?: string
  startedAt: string
  expiresAt: string
  submittedAt?: string
  status: 'in_progress' | 'submitted' | 'graded'
  autoScore: number
  manualScore: number
  finalScore: number
  answers?: AssessmentAnswer[]
}

export interface AssessmentAnswer {
  questionId: string
  type: QuestionType
  questionText?: string
  options?: string[]
  selectedOption?: number
  correctOption?: number
  essayText?: string
  autoScore?: number
  manualScore?: number
  weight?: number
  feedback?: string
}

// Project
export interface Project {
  id: string
  title: string
  description?: string
  instructions: string
  deadline: string
  subjectId: string
  subjectName?: string
  classIds: string[]
  classNames?: string[]
  attachments: FileRef[]
  allowResubmit: boolean
  totalSubmissions?: number
  gradedSubmissions?: number
  createdAt: string
  updatedAt: string
}

export interface CreateProjectPayload {
  title: string
  description?: string
  instructions: string
  deadline: string
  subjectId: string
  classIds: string[]
  allowResubmit: boolean
  attachments?: File[]
}

export interface ProjectSubmission {
  id: string
  projectId: string
  studentId: string
  studentName: string
  studentNis?: string
  files: FileRef[]
  submittedAt: string
  version: number
  score?: number
  feedback?: string
}

// Grades
export interface GradeItem {
  id: string
  studentId: string
  studentName: string
  studentNis?: string
  subjectId: string
  subjectName?: string
  classId: string
  className?: string
  academicYearId?: string
  semester?: number
  type: 'assignment' | 'assessment' | 'project' | 'uts' | 'uas'
  sourceTitle?: string
  score: number
  createdAt: string
}

export interface StudentGradeRow {
  studentId: string
  studentName: string
  studentNis: string
  assignmentScore?: number
  assessmentScore?: number
  projectScore?: number
  dailyAvg?: number // 60%
  utsScore?: number // 20%
  uasScore?: number // 20%
  finalScore?: number
}

// Walas
export interface WalasStudentDetail {
  id: string
  name: string
  nis: string
  gender: string
  email?: string
  phone?: string
  parentName?: string
  parentPhone?: string
  averageScore?: number
  rank?: number
}

export interface WalasSubjectScore {
  subjectId: string
  subjectName: string
  dailyAvg: number
  uts: number
  uas: number
  finalScore: number
}

export interface WalasGradeRekapRow {
  studentId: string
  studentName: string
  studentNis: string
  subjectScores: WalasSubjectScore[]
  overallAverage: number
  rank?: number
}
