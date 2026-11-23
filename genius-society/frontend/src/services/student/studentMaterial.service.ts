import { api } from '../api'
import type { Material } from '../../types/student'

export const getMaterials = async (subjectId?: string): Promise<Material[]> => {
  const res = await api.get<{ data: Material[] }>('/student/materials', {
    params: subjectId ? { subjectId } : undefined,
  })
  return res.data?.data ?? []
}

export const getMaterialDetail = async (id: string): Promise<Material> => {
  const res = await api.get<{ data: Material }>(`/student/materials/${id}`)
  return res.data.data
}

export const getMaterialDownloadURL = async (
  id: string,
  fileIndex: number
): Promise<{ url: string; expiresIn: number }> => {
  const res = await api.get<{ data: { url: string; expiresIn: number } }>(
    `/student/materials/${id}/files/${fileIndex}/download`
  )
  return res.data.data
}
