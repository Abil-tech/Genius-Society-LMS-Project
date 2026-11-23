import { api } from '../api'
import type { StudentDashboard } from '../../types/student'

export const getStudentDashboard = async (): Promise<StudentDashboard> => {
  const res = await api.get<{ data: StudentDashboard }>('/student/dashboard')
  return res.data.data
}
