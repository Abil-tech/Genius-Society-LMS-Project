import { api } from '../api'
import type {
  Assessment,
  AssessmentStartResponse,
  AnswerSubmission,
  AssessmentResultResponse,
} from '../../types/student'

export const getAssessments = async (subjectId?: string): Promise<Assessment[]> => {
  const res = await api.get<{ data: Assessment[] }>('/student/assessments', {
    params: subjectId ? { subjectId } : undefined,
  })
  return res.data?.data ?? []
}

export const getAssessmentDetail = async (id: string): Promise<Assessment> => {
  const res = await api.get<{ data: Assessment }>(`/student/assessments/${id}`)
  return res.data.data
}

export const startAssessment = async (id: string): Promise<AssessmentStartResponse> => {
  const res = await api.post<{ data: AssessmentStartResponse }>(
    `/student/assessments/${id}/start`
  )
  return res.data.data
}

export const submitAssessment = async (
  id: string,
  answers: AnswerSubmission[]
): Promise<AssessmentResultResponse> => {
  const res = await api.post<{ data: AssessmentResultResponse }>(
    `/student/assessments/${id}/submit`,
    { answers }
  )
  return res.data.data
}

export const getAssessmentResult = async (
  id: string
): Promise<AssessmentResultResponse> => {
  const res = await api.get<{ data: AssessmentResultResponse }>(
    `/student/assessments/${id}/result`
  )
  return res.data.data
}
