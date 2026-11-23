export interface FileRef {
  key: string
  name: string
  size: number
  mime: string
}

export interface DeadlineItem {
  id: string
  type: 'assignment' | 'project'
  title: string
  subject: string
  deadline: string
  submitted: boolean
}

export interface RecentGrade {
  id: string
  type: string
  title: string
  subject: string
  score: number
  gradedAt: string
}

export interface ScheduleItem {
  id: string
  day: number
  dayName: string
  startTime: string
  endTime: string
  subject: string
  teacher: string
  room: string
}

export interface StudentDashboard {
  activeAssignmentsCount: number
  activeAssessmentsCount: number
  upcomingDeadlines: DeadlineItem[]
  recentGrades: RecentGrade[]
  todaySchedule: ScheduleItem[]
}

export interface Material {
  id: string
  title: string
  description: string
  subjectId: string
  subjectName: string
  teacherName: string
  fileCount: number
  link?: string
  publishAt?: string
  createdAt: string
  attachments?: FileRef[]
}

export interface Submission {
  id: string
  submitted: boolean
  submittedAt?: string
  files: FileRef[]
  note?: string
  version: number
  score?: number
  feedback?: string
  gradedAt?: string
  isLate: boolean
}

export interface Assignment {
  id: string
  title: string
  description?: string
  subjectId: string
  subjectName: string
  teacherName: string
  deadline: string
  allowResubmit: boolean
  status: 'pending' | 'submitted' | 'graded' | 'overdue'
  attachments?: FileRef[]
  submission?: Submission
  createdAt: string
}

export interface Project {
  id: string
  title: string
  description?: string
  instructions?: string
  subjectId: string
  subjectName: string
  teacherName: string
  deadline: string
  allowResubmit: boolean
  status: 'pending' | 'submitted' | 'graded' | 'overdue'
  attachments?: FileRef[]
  submission?: Submission
  createdAt: string
}

export interface Assessment {
  id: string
  title: string
  description: string
  durationMinutes: number
  subjectId: string
  subjectName: string
  teacherName: string
  category: string
  startAt?: string
  endAt?: string
  status: 'available' | 'in_progress' | 'submitted' | 'completed' | 'closed'
  attemptId?: string
  score?: number
  questionCount?: number
  attemptStatus?: string
  startedAt?: string
  expiresAt?: string
  totalScore?: number
  showResult?: boolean
}

export interface Question {
  id: string
  type: 'mcq' | 'essay'
  text: string
  options?: string[]
  weight: number
  order: number
  answer?: string
}

export interface AssessmentStartResponse {
  attemptId: string
  assessmentId: string
  title: string
  durationMinutes: number
  startedAt: string
  expiresAt: string
  questions: Question[]
}

export interface AnswerSubmission {
  questionId: string
  answer: string
}

export interface AnswerResult {
  questionId: string
  type: 'mcq' | 'essay'
  text: string
  options?: string[]
  studentAnswer: string
  correctOption?: number
  isCorrect?: boolean
  score: number
  maxWeight: number
  feedback?: string
}

export interface AssessmentResultResponse {
  assessmentId: string
  title: string
  subjectName: string
  totalScore: number
  autoScore: number
  manualScore: number
  status: string
  submittedAt?: string
  answers?: AnswerResult[]
}

export interface Grade {
  subjectId: string
  subjectName: string
  subjectCode: string
  teacherName: string
  assignmentAverage?: number
  assessmentAverage?: number
  projectAverage?: number
  dailyAverage?: number
  utsScore?: number
  uasScore?: number
  finalScore?: number
  gradeLetter: string
}

export interface GradeDetailItem {
  id: string
  type: string
  title: string
  score: number
  gradedAt: string
  feedback?: string
}

export interface SubjectGradeDetail {
  subjectId: string
  subjectName: string
  finalScore?: number
  gradeLetter: string
  assignments: GradeDetailItem[]
  assessments: GradeDetailItem[]
  projects: GradeDetailItem[]
  uts?: GradeDetailItem
  uas?: GradeDetailItem
}

export interface RankingItem {
  rank: number
  studentId: string
  name: string
  nis: string
  finalScore: number
  isSelf: boolean
}

export interface ClassRankingResponse {
  className: string
  myRank: number
  totalStudents: number
  rankings: RankingItem[]
}
