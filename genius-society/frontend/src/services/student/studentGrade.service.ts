import { api } from '../api'
import type { Grade, SubjectGradeDetail } from '../../types/student'

export const getGrades = async (): Promise<Grade[]> => {
  const res = await api.get<{ data: Grade[] }>('/student/grades')
  return res.data?.data ?? []
}

export const getSubjectGradeDetail = async (
  subjectId: string
): Promise<SubjectGradeDetail> => {
  const res = await api.get<{ data: SubjectGradeDetail }>(
    `/student/grades/${subjectId}`
  )
  return res.data.data
}
