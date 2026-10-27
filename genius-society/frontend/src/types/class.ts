export interface Class {
  id: string;
  code: string;
  name: string;
  educationLevel: string;
  department?: string;
  academicYearId: string;
  homeroomTeacherId?: string;
  capacity: number;
  status: string;
  activeStudentCount?: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateClassReq {
  code: string;
  name: string;
  educationLevel: string;
  department?: string;
  academicYearId: string;
  homeroomTeacherId?: string;
  capacity: number;
  status: string;
}

export interface UpdateClassReq {
  code?: string;
  name?: string;
  educationLevel?: string;
  department?: string;
  academicYearId?: string;
  homeroomTeacherId?: string;
  capacity?: number;
  status?: string;
}
