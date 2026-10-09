import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useParams } from '@tanstack/react-router'
import FormField from '../components/FormField'
import {
  useAddTestCase,
  useAdminQuiz,
  useDeleteTestCase,
  useGenerateTestCases,
} from '../hooks/useQuizzes'
import { getErrorMessage } from '../services/api'

const generateSchema = z.object({
  count: z.coerce.number().int().min(1, 'Minimal 1').max(10, 'Maksimal 10'),
  problem_statement: z.string().max(10000, 'Maksimal 10000 karakter'),
})

const testCaseSchema = z.object({
  stdin: z.string().max(10000, 'Maksimal 10000 karakter'),
  expected_output: z.string().min(1, 'Output yang diharapkan wajib diisi').max(10000, 'Maksimal 10000 karakter'),
  is_hidden: z.boolean(),
})

function GenerateForm({ quizId }) {
  const generate = useGenerateTestCases(quizId)
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({ resolver: zodResolver(generateSchema), defaultValues: { count: 5, problem_statement: '' } })

  const onSubmit = async (values) => {
    try {
      await generate.mutateAsync({
        count: values.count,
        problem_statement: values.problem_statement,
      })
    } catch {
      // Ditampilkan dari state mutation.
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-3 p-4" noValidate>
      <h2 className="text-lg font-bold">Generate testcase dengan AI</h2>
      <FormField label="Jumlah testcase (1–10)" error={errors.count?.message}>
        <input type="number" className="input" {...register('count')} />
      </FormField>
      <FormField label="Problem statement khusus (opsional, default: deskripsi soal)" error={errors.problem_statement?.message}>
        <textarea rows={3} className="input" {...register('problem_statement')} />
      </FormField>
      {generate.isError ? <p role="alert" className="text-sm text-fail">{getErrorMessage(generate.error)}</p> : null}
      {generate.isSuccess ? <p className="text-sm text-pass">{generate.data.length} testcase ditambahkan.</p> : null}
      <button type="submit" className="btn-primary" disabled={generate.isPending}>
        {generate.isPending ? 'Menghasilkan…' : 'Generate'}
      </button>
    </form>
  )
}

function ManualForm({ quizId }) {
  const addTestCase = useAddTestCase(quizId)
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(testCaseSchema),
    defaultValues: { stdin: '', expected_output: '', is_hidden: true },
  })

  const onSubmit = async (values) => {
    try {
      await addTestCase.mutateAsync(values)
      reset()
    } catch {
      // Ditampilkan dari state mutation.
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-3 p-4" noValidate>
      <h2 className="text-lg font-bold">Tambah testcase manual</h2>
      <FormField label="Input (stdin)" error={errors.stdin?.message}>
        <textarea rows={3} className="input font-mono" {...register('stdin')} />
      </FormField>
      <FormField label="Output yang diharapkan" error={errors.expected_output?.message}>
        <textarea rows={3} className="input font-mono" {...register('expected_output')} />
      </FormField>
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" {...register('is_hidden')} />
        Sembunyikan dari siswa
      </label>
      {addTestCase.isError ? <p role="alert" className="text-sm text-fail">{getErrorMessage(addTestCase.error)}</p> : null}
      <button type="submit" className="btn-primary" disabled={addTestCase.isPending}>
        {addTestCase.isPending ? 'Menyimpan…' : 'Tambah'}
      </button>
    </form>
  )
}

export default function AdminQuizPage() {
  const { quizId } = useParams({ strict: false })
  const { data: quiz, isLoading, isError, error } = useAdminQuiz(quizId)
  const deleteTestCase = useDeleteTestCase(quizId)

  if (isLoading) return <p className="text-ink-soft">Memuat…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  const testCases = quiz.test_cases ?? []

  return (
    <div className="space-y-6">
      <Link to="/admin" className="text-sm text-signal">
        ← Dashboard admin
      </Link>
      <div>
        <h1 className="text-3xl font-bold">{quiz.title}</h1>
        <p className="mt-2 max-w-3xl whitespace-pre-wrap text-sm text-ink-soft">{quiz.description}</p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <GenerateForm quizId={quizId} />
        <ManualForm quizId={quizId} />
      </div>

      <section>
        <h2 className="mb-3 text-xl font-bold">Testcase ({testCases.length})</h2>
        {deleteTestCase.isError ? (
          <p role="alert" className="mb-2 text-sm text-fail">{getErrorMessage(deleteTestCase.error)}</p>
        ) : null}
        {testCases.length === 0 ? (
          <p className="text-ink-soft">Belum ada testcase. Soal ini belum bisa dikerjakan siswa.</p>
        ) : (
          <ul className="space-y-3">
            {testCases.map((tc, i) => (
              <li key={tc.id} className="panel p-3 text-xs">
                <div className="mb-2 flex items-center justify-between">
                  <span className="font-semibold">
                    #{i + 1} · {tc.source} · {tc.is_hidden ? 'tersembunyi' : 'contoh'}
                  </span>
                  <button
                    type="button"
                    className="text-fail"
                    onClick={() => {
                      if (window.confirm('Hapus testcase ini?')) deleteTestCase.mutate(tc.id)
                    }}
                  >
                    Hapus
                  </button>
                </div>
                <pre className="overflow-auto">stdin:{'\n'}{tc.stdin || '(kosong)'}</pre>
                <pre className="mt-2 overflow-auto">expected:{'\n'}{tc.expected_output}</pre>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}
