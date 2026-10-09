import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Navigate } from '@tanstack/react-router'
import { api, errorMessage } from '../services/api'
import { useAuth } from '../store/auth'

const schema = z.object({
  name: z.string().optional(),
  email: z.string().email('Enter a valid email address'),
  password: z.string().min(8, 'Use at least 8 characters'),
})

export default function Login() {
  const user = useAuth((s) => s.user)
  const setSession = useAuth((s) => s.setSession)
  const [register, setRegister] = useState(false)
  const [error, setError] = useState('')
  const { register: field, handleSubmit, formState: { errors, isSubmitting } } = useForm({ resolver: zodResolver(schema) })

  if (user) return <Navigate to="/" />

  const onSubmit = async (values) => {
    setError('')
    try {
      if (register) {
        if (!values.name || values.name.length < 2) return setError('Enter your name')
        await api.post('/auth/register', values)
      }
      const r = await api.post('/auth/login', { email: values.email, password: values.password })
      setSession(r.data.data.user, r.data.data.access_token)
    } catch (e) {
      setError(errorMessage(e))
    }
  }

  return (
    <main className="auth">
      <form onSubmit={handleSubmit(onSubmit)} className="panel stack">
        <h1>{register ? 'Create your account' : 'Sign in'}</h1>
        {register && <label>Name<input {...field('name')} autoComplete="name" /></label>}
        <label>Email<input {...field('email')} type="email" autoComplete="email" /></label>
        {errors.email && <p className="err">{errors.email.message}</p>}
        <label>Password<input {...field('password')} type="password" autoComplete={register ? 'new-password' : 'current-password'} /></label>
        {errors.password && <p className="err">{errors.password.message}</p>}
        {error && <p className="err">{error}</p>}
        <button disabled={isSubmitting}>{register ? 'Create account' : 'Sign in'}</button>
        <button type="button" className="ghost" onClick={() => setRegister(!register)}>
          {register ? 'I already have an account' : 'I need an account'}
        </button>
      </form>
    </main>
  )
}
