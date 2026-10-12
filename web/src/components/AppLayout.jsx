import { Link, Outlet, useNavigate } from '@tanstack/react-router'
import { useAuthStore } from '../store/authStore'
import { useLogout } from '../hooks/useAuth'

export default function AppLayout() {
  const user = useAuthStore((s) => s.user)
  const navigate = useNavigate()
  const logout = useLogout()

  const handleLogout = async () => {
    await logout.mutateAsync()
    navigate({ to: '/login' })
  }

  return (
    <div className="min-h-screen">
      <header className="border-b border-rule bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
          <nav className="flex items-center gap-6">
            <Link to="/courses" className="font-display text-xl font-bold text-ink">
              C-Grader
            </Link>
            <Link to="/courses" className="text-sm text-ink-soft hover:text-ink">
              Kursus
            </Link>
            {user?.role === 'ADMIN' ? (
              <Link to="/admin" className="text-sm text-ink-soft hover:text-ink">
                Admin
              </Link>
            ) : null}
          </nav>
          <div className="flex items-center gap-3 text-sm">
            <span className="text-ink-soft">{user?.name}</span>
            <button type="button" className="btn-ghost" onClick={handleLogout}>
              Keluar
            </button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-6xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
