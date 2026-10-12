import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useNavigate } from '@tanstack/react-router'
import FormField from '../components/FormField'
import { useRegister } from '../hooks/useAuth'
import { getErrorMessage } from '../services/api'

const schema = z.object({
  name: z.string().min(2, 'Minimal 2 karakter').max(100, 'Maksimal 100 karakter'),
  email: z.string().email('Email tidak valid'),
  password: z.string().min(8, 'Minimal 8 karakter').max(72, 'Maksimal 72 karakter'),
})

export default function RegisterPage() {
  const navigate = useNavigate()
  const registerUser = useRegister()
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({ resolver: zodResolver(schema) })

  const onSubmit = async (values) => {
    try {
      await registerUser.mutateAsync(values)
      navigate({ to: '/courses' })
    } catch {
      // Pesan error ditampilkan dari state mutation.
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-sm">
        <h1 className="mb-6 text-3xl font-bold">Buat akun</h1>
        <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-4 p-5" noValidate>
          <FormField label="Nama" error={errors.name?.message}>
            <input type="text" autoComplete="name" className="input" {...register('name')} />
          </FormField>
          <FormField label="Email" error={errors.email?.message}>
            <input type="email" autoComplete="email" className="input" {...register('email')} />
          </FormField>
          <FormField label="Password" error={errors.password?.message}>
            <input type="password" autoComplete="new-password" className="input" {...register('password')} />
          </FormField>
          {registerUser.isError ? (
            <p role="alert" className="text-sm text-fail">
              {getErrorMessage(registerUser.error)}
            </p>
          ) : null}
          <button type="submit" className="btn-primary w-full" disabled={registerUser.isPending}>
            {registerUser.isPending ? 'Memproses…' : 'Daftar'}
          </button>
        </form>
        <p className="mt-4 text-center text-sm text-ink-soft">
          Sudah punya akun?{' '}
          <Link to="/login" className="font-semibold text-signal">
            Masuk
          </Link>
        </p>
      </div>
    </div>
  )
}
