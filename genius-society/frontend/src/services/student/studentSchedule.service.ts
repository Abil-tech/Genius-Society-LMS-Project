import { api } from '../api'
import type { ScheduleItem } from '../../types/student'

export const getSchedule = async (): Promise<ScheduleItem[]> => {
  const res = await api.get<{ data: ScheduleItem[] }>('/student/schedule')
  return res.data?.data ?? []
}
