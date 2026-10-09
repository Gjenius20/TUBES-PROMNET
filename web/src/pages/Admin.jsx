import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api, errorMessage } from '../services/api'

const schema = z.object({
  title: z.string().min(3, 'Give the problem a title'),
  statement: z.string().min(20, 'Describe the problem in at least 20 characters'),
})

export default function Admin() {
  const qc = useQueryClient()
  const refresh = () => qc.invalidateQueries({ queryKey: ['quizzes'] })
  const quizzes = useQuery({ queryKey: ['quizzes'], queryFn: () => api.get('/quizzes').then((r) => r.data.data) })
  const { register, handleSubmit, reset, formState: { errors } } = useForm({ resolver: zodResolver(schema) })

  const create = useMutation({
    mutationFn: (v) => api.post('/admin/quizzes', v),
    onSuccess: () => { reset(); refresh() },
  })
  const generate = useMutation({ mutationFn: (id) => api.post('/admin/quizzes/generate-testcases', { quiz_id: id, count: 5 }) })
  const remove = useMutation({ mutationFn: (id) => api.delete(`/admin/quizzes/${id}`), onSuccess: refresh })

  return (
    <main className="admin">
      <form className="panel stack" onSubmit={handleSubmit((v) => create.mutate(v))}>
        <h2>New problem</h2>
        <label>Title<input {...register('title')} /></label>
        {errors.title && <p className="err">{errors.title.message}</p>}
        <label>Problem statement<textarea rows={6} {...register('statement')} /></label>
        {errors.statement && <p className="err">{errors.statement.message}</p>}
        {create.isError && <p className="err">{errorMessage(create.error)}</p>}
        <button disabled={create.isPending}>Create problem</button>
      </form>

      <section className="panel stack">
        <h2>Problems</h2>
        {quizzes.data?.length === 0 && <p className="muted">Nothing here yet. Create your first problem.</p>}
        {quizzes.data?.map((q) => (
          <div key={q.id} className="row">
            <span>{q.title}</span>
            <span className="spacer" />
            <button className="ghost" disabled={generate.isPending} onClick={() => generate.mutate(q.id)}>Generate test cases</button>
            <button className="ghost danger" onClick={() => window.confirm(`Delete "${q.title}" with its test cases and submissions?`) && remove.mutate(q.id)}>Delete</button>
          </div>
        ))}
        {generate.isSuccess && <p className="ok">Added {generate.data.data.data.length} test cases.</p>}
        {generate.isError && <p className="err">{errorMessage(generate.error)}</p>}
        {remove.isError && <p className="err">{errorMessage(remove.error)}</p>}
      </section>
    </main>
  )
}
