import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from '@tanstack/react-router'
import FormField from '../components/FormField'
import { useCourses, useCreateCourse, useDeleteCourse } from '../hooks/useCourses'
import { useCreateQuiz, useDeleteQuiz } from '../hooks/useQuizzes'
import { getErrorMessage } from '../services/api'

const courseSchema = z.object({
  title: z.string().min(2, 'Minimal 2 karakter').max(150, 'Maksimal 150 karakter'),
  description: z.string().max(5000, 'Maksimal 5000 karakter'),
})

const quizSchema = z.object({
  course_id: z.coerce.number().int().positive('Pilih kursus'),
  title: z.string().min(2, 'Minimal 2 karakter').max(150, 'Maksimal 150 karakter'),
  description: z.string().min(1, 'Deskripsi soal wajib diisi').max(10000, 'Maksimal 10000 karakter'),
  constraints: z.string().max(5000, 'Maksimal 5000 karakter'),
})

function CourseForm() {
  const createCourse = useCreateCourse()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({ resolver: zodResolver(courseSchema), defaultValues: { title: '', description: '' } })

  const onSubmit = async (values) => {
    try {
      await createCourse.mutateAsync(values)
      reset()
    } catch {
      // Ditampilkan dari state mutation.
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-3 p-4" noValidate>
      <h2 className="text-lg font-bold">Kursus baru</h2>
      <FormField label="Judul" error={errors.title?.message}>
        <input className="input" {...register('title')} />
      </FormField>
      <FormField label="Deskripsi" error={errors.description?.message}>
        <textarea rows={3} className="input" {...register('description')} />
      </FormField>
      {createCourse.isError ? <p role="alert" className="text-sm text-fail">{getErrorMessage(createCourse.error)}</p> : null}
      <button type="submit" className="btn-primary" disabled={createCourse.isPending}>
        {createCourse.isPending ? 'Menyimpan…' : 'Buat kursus'}
      </button>
    </form>
  )
}

function QuizForm({ courses }) {
  const createQuiz = useCreateQuiz()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(quizSchema),
    defaultValues: { course_id: '', title: '', description: '', constraints: '' },
  })

  const onSubmit = async (values) => {
    try {
      await createQuiz.mutateAsync(values)
      reset()
    } catch {
      // Ditampilkan dari state mutation.
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-3 p-4" noValidate>
      <h2 className="text-lg font-bold">Soal baru</h2>
      <FormField label="Kursus" error={errors.course_id?.message}>
        <select className="input" {...register('course_id')}>
          <option value="">Pilih kursus…</option>
          {courses.map((c) => (
            <option key={c.id} value={c.id}>
              {c.title}
            </option>
          ))}
        </select>
      </FormField>
      <FormField label="Judul soal" error={errors.title?.message}>
        <input className="input" {...register('title')} />
      </FormField>
      <FormField label="Deskripsi soal (problem statement)" error={errors.description?.message}>
        <textarea rows={5} className="input" {...register('description')} />
      </FormField>
      <FormField label="Batasan (opsional)" error={errors.constraints?.message}>
        <textarea rows={2} className="input" {...register('constraints')} />
      </FormField>
      {createQuiz.isError ? <p role="alert" className="text-sm text-fail">{getErrorMessage(createQuiz.error)}</p> : null}
      <button type="submit" className="btn-primary" disabled={createQuiz.isPending}>
        {createQuiz.isPending ? 'Menyimpan…' : 'Buat soal'}
      </button>
    </form>
  )
}

export default function AdminPage() {
  const { data: courses, isLoading, isError, error } = useCourses()
  const deleteCourse = useDeleteCourse()
  const deleteQuiz = useDeleteQuiz()

  if (isLoading) return <p className="text-ink-soft">Memuat…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  return (
    <div className="space-y-8">
      <h1 className="text-3xl font-bold">Dashboard Admin</h1>

      <div className="grid gap-6 lg:grid-cols-2">
        <CourseForm />
        <QuizForm courses={courses} />
      </div>

      <section>
        <h2 className="mb-3 text-xl font-bold">Kursus & soal</h2>
        {deleteCourse.isError || deleteQuiz.isError ? (
          <p role="alert" className="mb-2 text-sm text-fail">
            {getErrorMessage(deleteCourse.error || deleteQuiz.error)}
          </p>
        ) : null}
        <div className="space-y-4">
          {courses.map((course) => (
            <div key={course.id} className="panel p-4">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <h3 className="text-lg font-bold">{course.title}</h3>
                  <p className="text-sm text-ink-soft">{course.description}</p>
                </div>
                <button
                  type="button"
                  className="btn-danger"
                  onClick={() => {
                    if (window.confirm(`Hapus kursus "${course.title}" beserta semua soalnya?`)) {
                      deleteCourse.mutate(course.id)
                    }
                  }}
                >
                  Hapus
                </button>
              </div>
              <ul className="mt-3 divide-y divide-rule text-sm">
                {(course.quizzes ?? []).map((quiz) => (
                  <li key={quiz.id} className="flex items-center justify-between py-2">
                    <span>{quiz.title}</span>
                    <span className="flex items-center gap-3">
                      <Link
                        to="/admin/quizzes/$quizId"
                        params={{ quizId: String(quiz.id) }}
                        className="font-semibold text-signal"
                      >
                        Testcase
                      </Link>
                      <button
                        type="button"
                        className="text-fail"
                        onClick={() => {
                          if (window.confirm(`Hapus soal "${quiz.title}"?`)) deleteQuiz.mutate(quiz.id)
                        }}
                      >
                        Hapus
                      </button>
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
