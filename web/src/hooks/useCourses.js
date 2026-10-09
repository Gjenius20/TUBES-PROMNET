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
