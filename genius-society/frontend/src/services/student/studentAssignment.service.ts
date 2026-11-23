import { api } from '../api'
import type { Assignment, Submission } from '../../types/student'

export const getAssignments = async (subjectId?: string): Promise<Assignment[]> => {
  const res = await api.get<{ data: Assignment[] }>('/student/assignments', {
    params: subjectId ? { subjectId } : undefined,
  })
  return res.data?.data ?? []
}

export const getAssignmentDetail = async (id: string): Promise<Assignment> => {
  const res = await api.get<{ data: Assignment }>(`/student/assignments/${id}`)
  return res.data.data
}

export const submitAssignment = async (
  id: string,
  formData: FormData
): Promise<Submission> => {
  const res = await api.post<{ data: Submission }>(
    `/student/assignments/${id}/submit`,
    formData,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
    }
  )
  return res.data.data
}

export const resubmitAssignment = async (
  id: string,
  formData: FormData
): Promise<Submission> => {
  const res = await api.put<{ data: Submission }>(
    `/student/assignments/${id}/resubmit`,
    formData,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
    }
  )
  return res.data.data
}

export const getAssignmentFileDownloadURL = async (
  id: string,
  fileIndex: number
): Promise<{ url: string; expiresIn: number }> => {
  const res = await api.get<{ data: { url: string; expiresIn: number } }>(
    `/student/assignments/${id}/files/${fileIndex}/download`
  )
  return res.data.data
}
