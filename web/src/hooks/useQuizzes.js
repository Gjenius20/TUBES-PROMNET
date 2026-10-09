import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'

export function useQuiz(quizId) {
  return useQuery({
    queryKey: ['quiz', quizId],
    queryFn: async () => (await api.get(`/quizzes/${quizId}`)).data.data.quiz,
    enabled: Boolean(quizId),
  })
}

export function useAdminQuiz(quizId) {
  return useQuery({
    queryKey: ['admin-quiz', quizId],
    queryFn: async () => (await api.get(`/admin/quizzes/${quizId}`)).data.data.quiz,
    enabled: Boolean(quizId),
  })
}

export function useSubmissions(quizId) {
  return useQuery({
    queryKey: ['submissions', quizId],
    queryFn: async () =>
      (await api.get('/submissions', { params: { quiz_id: quizId } })).data.data.submissions,
    enabled: Boolean(quizId),
  })
}

export function useSubmit(quizId) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({ sourceCode }) =>
      (await api.post('/submissions', { quiz_id: Number(quizId), source_code: sourceCode })).data
        .data.submission,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['submissions', quizId] }),
  })
}

export function useCreateQuiz() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (values) => (await api.post('/admin/quizzes', values)).data.data.quiz,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['courses'] }),
  })
}

export function useDeleteQuiz() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id) => api.delete(`/admin/quizzes/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['courses'] }),
  })
}

export function useAddTestCase(quizId) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (values) =>
      (await api.post(`/admin/quizzes/${quizId}/testcases`, values)).data.data.test_case,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-quiz', quizId] }),
  })
}

export function useDeleteTestCase(quizId) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id) => api.delete(`/admin/testcases/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-quiz', quizId] }),
  })
}

export function useGenerateTestCases(quizId) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (values) =>
      (await api.post('/admin/quizzes/generate-testcases', { quiz_id: Number(quizId), ...values }))
        .data.data.test_cases,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['admin-quiz', quizId] }),
  })
}
