import axios from 'axios'
import { useAuth } from '../store/auth'

const baseURL = import.meta.env.VITE_API_URL || '/api/v1'

export const api = axios.create({ baseURL, withCredentials: true })
const bare = axios.create({ baseURL, withCredentials: true })

api.interceptors.request.use((config) => {
  const token = useAuth.getState().token
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let refreshing = null

api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config
    if (error.response?.status === 401 && !original._retried) {
      original._retried = true
      refreshing ??= bare
        .post('/auth/refresh')
        .then((r) => {
          const { user, access_token } = r.data.data
          useAuth.getState().setSession(user, access_token)
          return access_token
        })
        .finally(() => {
          refreshing = null
        })
      try {
        const token = await refreshing
        original.headers.Authorization = `Bearer ${token}`
        return api(original)
      } catch {
        useAuth.getState().clear()
      }
    }
    return Promise.reject(error)
  },
)

export const errorMessage = (e) => e.response?.data?.errors?.[0] || e.response?.data?.message || e.message

export const bootstrapSession = async () => {
  try {
    const r = await bare.post('/auth/refresh')
    useAuth.getState().setSession(r.data.data.user, r.data.data.access_token)
  } catch {
    /* not signed in */
  }
}
