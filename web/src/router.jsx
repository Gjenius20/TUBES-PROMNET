import { createRootRoute, createRoute, createRouter, Outlet, redirect } from '@tanstack/react-router'
import { bootstrapAuth } from './services/api'
import { useAuthStore } from './store/authStore'
import AppLayout from './components/AppLayout'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import CoursesPage from './pages/CoursesPage'
import CoursePage from './pages/CoursePage'
import WorkspacePage from './pages/WorkspacePage'
import AdminPage from './pages/AdminPage'
import AdminQuizPage from './pages/AdminQuizPage'

const rootRoute = createRootRoute({
  component: Outlet,
  // Pulihkan sesi (lewat refresh cookie) sebelum route mana pun diputuskan.
  beforeLoad: async () => {
    await bootstrapAuth()
  },
  notFoundComponent: () => <p className="p-8 text-ink-soft">Halaman tidak ditemukan.</p>,
})

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: () => {
    const { user } = useAuthStore.getState()
    throw redirect({ to: user ? '/courses' : '/login' })
  },
})

// Halaman tamu: pengguna yang sudah login diarahkan ke /courses.
function redirectIfAuthenticated() {
  if (useAuthStore.getState().user) throw redirect({ to: '/courses' })
}

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  beforeLoad: redirectIfAuthenticated,
  component: LoginPage,
})

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  beforeLoad: redirectIfAuthenticated,
  component: RegisterPage,
})

// Layout route tanpa path untuk semua halaman yang wajib login.
const authedRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'authed',
  beforeLoad: () => {
    if (!useAuthStore.getState().user) throw redirect({ to: '/login' })
  },
  component: AppLayout,
})

const coursesRoute = createRoute({ getParentRoute: () => authedRoute, path: '/courses', component: CoursesPage })
const courseRoute = createRoute({ getParentRoute: () => authedRoute, path: '/courses/$courseId', component: CoursePage })
const quizRoute = createRoute({ getParentRoute: () => authedRoute, path: '/quizzes/$quizId', component: WorkspacePage })

function requireAdmin() {
  if (useAuthStore.getState().user?.role !== 'ADMIN') throw redirect({ to: '/courses' })
}

const adminRoute = createRoute({
  getParentRoute: () => authedRoute,
  path: '/admin',
  beforeLoad: requireAdmin,
  component: AdminPage,
})
const adminQuizRoute = createRoute({
  getParentRoute: () => authedRoute,
  path: '/admin/quizzes/$quizId',
  beforeLoad: requireAdmin,
  component: AdminQuizPage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  registerRoute,
  authedRoute.addChildren([coursesRoute, courseRoute, quizRoute, adminRoute, adminQuizRoute]),
])

export const router = createRouter({ routeTree })
