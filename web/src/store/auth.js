import { create } from 'zustand'

export const useAuth = create((set) => ({
  user: null,
  token: null,
  ready: false,
  setSession: (user, token) => set({ user, token, ready: true }),
  clear: () => set({ user: null, token: null, ready: true }),
}))
