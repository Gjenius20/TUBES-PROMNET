import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'

export function useCourses() {
  return useQuery({
    queryKey: ['courses'],
    queryFn: async () => (await api.get('/courses')).data.data.courses,
  })
}

export function useCourse(courseId) {
  return useQuery({
    queryKey: ['courses', courseId],
    queryFn: async () => (await api.get(`/courses/${courseId}`)).data.data.course,
    enabled: Boolean(courseId),
  })
}

export function useCreateCourse() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (values) => (await api.post('/admin/courses', values)).data.data.course,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['courses'] }),
  })
}

export function useDeleteCourse() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id) => api.delete(`/admin/courses/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['courses'] }),
  })
}

export function useEnroll(courseId) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async () => (await api.post(`/courses/${courseId}/enroll`)).data.data.enrollment,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['studentEnrollments'] }),
  })
}

export function useStudentEnrollments() {
  return useQuery({
    queryKey: ['studentEnrollments'],
    queryFn: async () => (await api.get('/student/enrollments')).data.data.enrollments,
  })
}

export function useAdminEnrollments(status) {
  return useQuery({
    queryKey: ['adminEnrollments', status],
    queryFn: async () => (await api.get('/admin/enrollments', { params: { status } })).data.data.enrollments,
  })
}

export function useUpdateEnrollmentStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, status }) => (await api.patch(`/admin/enrollments/${id}`, { status })).data.data.enrollment,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['adminEnrollments'] }),
  })
}

export function useExportGrades(courseId) {
  return useMutation({
    mutationFn: async () => {
      const res = await api.get(`/admin/courses/${courseId}/export-grades`, {
        responseType: 'blob',
      })
      return res.data
    },
  })
}
