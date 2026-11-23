import { api } from '../api'
import type { ClassRankingResponse } from '../../types/student'

export const getRanking = async (): Promise<ClassRankingResponse> => {
  const res = await api.get<{ data: ClassRankingResponse }>('/student/ranking')
  return res.data.data
}
