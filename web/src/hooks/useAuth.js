import { useMutation } from '@tanstack/react-query'
import { api } from '../services/api'
import { useAuthStore } from '../store/authStore'
import { queryClient } from '../queryClient'

function saveSession(res) {
  const { access_token: accessToken, user } = res.data.data
  useAuthStore.getState().setSession({ accessToken, user })
  return user
}

export function useLogin() {
  return useMutation({
    mutationFn: async (values) => saveSession(await api.post('/auth/login', values)),
  })
}

export function useRegister() {
  return useMutation({
    mutationFn: async (values) => saveSession(await api.post('/auth/register', values)),
  })
}

export function useLogout() {
  return useMutation({
    mutationFn: async () => {
      try {
        await api.post('/auth/logout')
      } catch {
        // Tetap logout di sisi klien walau request gagal.
      }
      useAuthStore.getState().clearSession()
      queryClient.clear()
    },
  })
}
