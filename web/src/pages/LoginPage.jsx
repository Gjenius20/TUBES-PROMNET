import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useNavigate } from '@tanstack/react-router'
import FormField from '../components/FormField'
import { useLogin } from '../hooks/useAuth'
import { getErrorMessage } from '../services/api'

const schema = z.object({
  email: z.string().email('Email tidak valid'),
  password: z.string().min(1, 'Password wajib diisi'),
})

export default function LoginPage() {
  const navigate = useNavigate()
  const login = useLogin()
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({ resolver: zodResolver(schema) })

  const onSubmit = async (values) => {
    try {
      await login.mutateAsync(values)
      navigate({ to: '/courses' })
    } catch {
      // Pesan error ditampilkan dari state mutation.
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-sm">
        <h1 className="mb-1 text-3xl font-bold">Masuk ke C-Grader</h1>
        <p className="mb-6 text-sm text-ink-soft">Kerjakan soal pemrograman C dan lihat nilainya langsung.</p>
        <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-4 p-5" noValidate>
          <FormField label="Email" error={errors.email?.message}>
            <input type="email" autoComplete="email" className="input" {...register('email')} />
          </FormField>
          <FormField label="Password" error={errors.password?.message}>
            <input type="password" autoComplete="current-password" className="input" {...register('password')} />
          </FormField>
          {login.isError ? (
            <p role="alert" className="text-sm text-fail">
              {getErrorMessage(login.error)}
            </p>
          ) : null}
          <button type="submit" className="btn-primary w-full" disabled={login.isPending}>
            {login.isPending ? 'Memproses…' : 'Masuk'}
          </button>
        </form>
        <p className="mt-4 text-center text-sm text-ink-soft">
          Belum punya akun?{' '}
          <Link to="/register" className="font-semibold text-signal">
            Daftar
          </Link>
        </p>
      </div>
    </div>
  )
}
