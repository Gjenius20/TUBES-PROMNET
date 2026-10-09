import { useEffect } from 'react'
import { createRootRoute, createRoute, createRouter, Link, Navigate, Outlet, useNavigate } from '@tanstack/react-router'
import { api, bootstrapSession } from './services/api'
import { useAuth } from './store/auth'
import Login from './pages/Login'
import Workspace from './pages/Workspace'
import Admin from './pages/Admin'

function Layout() {
  const { user, ready, clear } = useAuth()
  const navigate = useNavigate()
  useEffect(() => {
    bootstrapSession().finally(() => useAuth.setState({ ready: true }))
  }, [])
  if (!ready) return <p className="muted pad">Loading…</p>

  const logout = async () => {
    await api.post('/auth/logout').catch(() => {})
    clear()
    navigate({ to: '/login' })
  }
  return (
    <>
      {user && (
        <header className="bar">
          <strong>C Autograder</strong>
          <nav>
            <Link to="/">Workspace</Link>
            {user.role === 'ADMIN' && <Link to="/admin">Admin</Link>}
          </nav>
          <span className="spacer" />
          <span className="muted">{user.name}</span>
          <button className="ghost" onClick={logout}>Sign out</button>
        </header>
      )}
      <Outlet />
    </>
  )
}

const Guard = ({ admin, children }) => {
  const user = useAuth((s) => s.user)
  if (!user) return <Navigate to="/login" />
  if (admin && user.role !== 'ADMIN') return <Navigate to="/" />
  return children
}

const root = createRootRoute({ component: Layout })
const routeTree = root.addChildren([
  createRoute({ getParentRoute: () => root, path: '/login', component: Login }),
  createRoute({ getParentRoute: () => root, path: '/', component: () => <Guard><Workspace /></Guard> }),
  createRoute({ getParentRoute: () => root, path: '/admin', component: () => <Guard admin><Admin /></Guard> }),
])

export const router = createRouter({ routeTree })
