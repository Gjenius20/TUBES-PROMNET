import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'

export function useStudents(search, page) {
  return useQuery({
    queryKey: ['students', search, page],
    queryFn: async () => (await api.get('/admin/students', { params: { search, page, limit: 10 } })).data.data,
  })
}

export function useBulkEnroll(courseId) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (student_ids) => api.post(`/admin/courses/${courseId}/enroll-bulk`, { student_ids }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['adminEnrollments'] }),
  })
}

export function useRemoveStudent() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ courseId, studentId }) => api.delete(`/admin/courses/${courseId}/enrollments/${studentId}`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['adminEnrollments'] })
      // Invalidate course detail queries broadly as we don't know which course specifically here
      qc.invalidateQueries({ queryKey: ['courseDetail'] })
    },
  })
}