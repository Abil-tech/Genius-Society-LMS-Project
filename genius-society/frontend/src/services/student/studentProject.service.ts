import { api } from '../api'
import type { Project, Submission } from '../../types/student'

export const getProjects = async (subjectId?: string): Promise<Project[]> => {
  const res = await api.get<{ data: Project[] }>('/student/projects', {
    params: subjectId ? { subjectId } : undefined,
  })
  return res.data?.data ?? []
}

export const getProjectDetail = async (id: string): Promise<Project> => {
  const res = await api.get<{ data: Project }>(`/student/projects/${id}`)
  return res.data.data
}

export const submitProject = async (
  id: string,
  formData: FormData
): Promise<Submission> => {
  const res = await api.post<{ data: Submission }>(
    `/student/projects/${id}/submit`,
    formData,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
    }
  )
  return res.data.data
}

export const resubmitProject = async (
  id: string,
  formData: FormData
): Promise<Submission> => {
  const res = await api.put<{ data: Submission }>(
    `/student/projects/${id}/resubmit`,
    formData,
    {
      headers: { 'Content-Type': 'multipart/form-data' },
    }
  )
  return res.data.data
}

export const getProjectFileDownloadURL = async (
  id: string,
  fileIndex: number
): Promise<{ url: string; expiresIn: number }> => {
  const res = await api.get<{ data: { url: string; expiresIn: number } }>(
    `/student/projects/${id}/files/${fileIndex}/download`
  )
  return res.data.data
}
