import axios from 'axios'
import { useAuthStore } from '../store/authStore'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export const api = axios.create({ baseURL: API_URL, withCredentials: true })

// Instance tanpa interceptor, khusus untuk refresh (mencegah loop tak berujung).
const bare = axios.create({ baseURL: API_URL, withCredentials: true })

let refreshPromise = null

// Single-flight: banyak request 401 bersamaan hanya memicu satu refresh.
export function refreshSession() {
  if (!refreshPromise) {
    refreshPromise = bare
      .post('/auth/refresh')
      .then((res) => {
        const { access_token: accessToken, user } = res.data.data
        useAuthStore.getState().setSession({ accessToken, user })
        return accessToken
      })
      .finally(() => {
        refreshPromise = null
      })
  }
  return refreshPromise
}

let bootstrapPromise = null

// Memulihkan sesi sekali saat aplikasi dibuka.
export function bootstrapAuth() {
  if (useAuthStore.getState().initialized) return Promise.resolve()
  if (!bootstrapPromise) {
    bootstrapPromise = refreshSession()
      .catch(() => useAuthStore.getState().clearSession())
      .finally(() => {
        bootstrapPromise = null
      })
  }
  return bootstrapPromise
}

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

const AUTH_PATHS = ['/auth/login', '/auth/register', '/auth/refresh']

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config
    const status = error.response?.status
    const isAuthCall = AUTH_PATHS.some((p) => original?.url?.startsWith(p))

    if (status === 401 && original && !original._retry && !isAuthCall) {
      original._retry = true
      try {
        const token = await refreshSession()
        original.headers.Authorization = `Bearer ${token}`
        return api(original)
      } catch (refreshError) {
        useAuthStore.getState().clearSession()
        window.location.assign('/login')
        return Promise.reject(refreshError)
      }
    }
    return Promise.reject(error)
  },
)

export function getErrorMessage(error, fallback = 'Terjadi kesalahan. Coba lagi.') {
  const data = error?.response?.data
  if (data?.errors?.length) return `${data.message}: ${data.errors.join(', ')}`
  return data?.message || fallback
}
