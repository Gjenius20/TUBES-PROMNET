import { create } from 'zustand'

// Access token hanya disimpan di memori (bukan localStorage). Sesi dipulihkan
// saat reload lewat refresh token di cookie HttpOnly (lihat bootstrapAuth).
export const useAuthStore = create((set) => ({
  accessToken: null,
  user: null,
  initialized: false,
  setSession: ({ accessToken, user }) => set({ accessToken, user, initialized: true }),
  clearSession: () => set({ accessToken: null, user: null, initialized: true }),
}))
