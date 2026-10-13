import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from '@tanstack/react-router'
import FormField from '../components/FormField'
import StudentEnrollModal from '../components/StudentEnrollModal'
import { useCourseProgress } from '../hooks/useProgress'
import { useCourses, useCreateCourse, useDeleteCourse, useExportGrades } from '../hooks/useCourses'
import { useCreateQuiz, useDeleteQuiz } from '../hooks/useQuizzes'
import { useBulkEnroll, useRemoveStudent, useStudents } from '../hooks/useUsers'
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
      <button type="submit" className="btn-primary" disabled={createCourse.isPending}>
        {createCourse.isPending ? 'Menyimpan...' : 'Simpan Kursus'}
      </button>
      {createCourse.isError && <p className="text-fail text-sm">{getErrorMessage(createCourse.error)}</p>}
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
  } = useForm({ resolver: zodResolver(quizSchema), defaultValues: { course_id: '', title: '', description: '', constraints: '' } })

  const onSubmit = async (values) => {
    try {
      await createQuiz.mutateAsync(values)
      reset()
    } catch {
      // handled by mutation
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="panel space-y-3 p-4" noValidate>
      <h2 className="text-lg font-bold">Soal baru</h2>
      <FormField label="Kursus" error={errors.course_id?.message}>
        <select className="input" {...register('course_id')}>
          <option value="">Pilih kursus</option>
          {courses.map((c) => (
            <option key={c.id} value={c.id}>
              {c.title}
            </option>
          ))}
        </select>
      </FormField>
      <FormField label="Judul" error={errors.title?.message}>
        <input className="input" {...register('title')} />
      </FormField>
      <FormField label="Deskripsi Soal" error={errors.description?.message}>
        <textarea rows={4} className="input" {...register('description')} />
      </FormField>
      <FormField label="Batasan (opsional)" error={errors.constraints?.message}>
        <textarea rows={2} className="input" {...register('constraints')} />
      </FormField>
      <button type="submit" className="btn-primary" disabled={createQuiz.isPending}>
        {createQuiz.isPending ? 'Menyimpan...' : 'Simpan Soal'}
      </button>
      {createQuiz.isError && <p className="text-fail text-sm">{getErrorMessage(createQuiz.error)}</p>}
    </form>
  )
}

function CourseAdminItem({ course, deleteCourse, deleteQuiz, onManageStudents }) {
  const { data: progressData, isLoading: progressLoading, fetchProgress } = useCourseProgress(course.id)
  const exportGrades = useExportGrades(course.id)
  const [showProgress, setShowProgress] = useState(false)

  const handleToggleProgress = () => {
    if (!showProgress) {
      fetchProgress()
    }
    setShowProgress(!showProgress)
  }

  return (
    <div className="panel p-4 space-y-3">
      <div className="flex justify-between items-center">
        <h3 className="font-semibold">{course.title}</h3>
        <div className="flex gap-2">
          <button
            type="button"
            className="btn-primary text-sm"
            onClick={() => onManageStudents(course.id)}
          >
            + Siswa
          </button>
          <button
            type="button"
            className="btn-secondary text-sm"
            onClick={handleToggleProgress}
            disabled={progressLoading}
          >
            {showProgress ? 'Sembunyikan' : 'Tampilkan'} Progress
          </button>
          <button
            type="button"
            className="btn-secondary text-sm"
            onClick={() => exportGrades.mutate()}
            disabled={exportGrades.isPending}
          >
            {exportGrades.isPending ? 'Mengunduh...' : 'Unduh Nilai'}
          </button>
          <button
            type="button"
            className="text-signal"
            onClick={() => {
              if (window.confirm(`Hapus kursus "${course.title}"?`)) deleteCourse.mutate(course.id)
            }}
          >
            Hapus
          </button>
        </div>
      </div>

      <p className="text-sm text-ink-soft">{course.description}</p>

      {showProgress && progressData && (
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead>
              <tr className="border-b border-rule text-left">
                <th className="p-2">Nama</th>
                <th className="p-2">NIM</th>
                <th className="p-2">Rata-rata</th>
                <th className="p-2">Lulus / Total Soal</th>
              </tr>
            </thead>
            <tbody>
              {progressData.map((p) => (
                <tr key={p.user_id || p.nim} className="border-b border-rule">
                  <td className="p-2">{p.name}</td>
                  <td className="p-2">{p.nim}</td>
                  <td className="p-2">{p.average_score?.toFixed(2) ?? '0'}%</td>
                  <td className="p-2">{p.passed_count}/{p.total_count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <h4 className="font-medium">Daftar Soal</h4>
      {(course.quizzes ?? []).length === 0 ? (
        <p className="text-sm text-ink-soft">Belum ada soal.</p>
      ) : (
        <ul className="space-y-2">
          {(course.quizzes ?? []).map((quiz) => (
            <li key={quiz.id} className="flex items-center justify-between py-2 border-t border-rule">
              <span className="text-sm">{quiz.title}</span>
              <div className="flex items-center gap-3">
                <Link
                  to="/admin/quizzes/$quizId"
                  params={{ quizId: String(quiz.id) }}
                  className="font-semibold text-signal text-sm"
                >
                  Testcase
                </Link>
                <button
                  type="button"
                  className="text-fail text-sm"
                  onClick={() => {
                    if (window.confirm(`Hapus soal "${quiz.title}"?`)) deleteQuiz.mutate(quiz.id)
                  }}
                >
                  Hapus
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default function AdminPage() {
  const { data: courses, isLoading, isError, error } = useCourses()
  const deleteCourse = useDeleteCourse()
  const deleteQuiz = useDeleteQuiz()
  const [selectedCourseId, setSelectedCourseId] = useState(null)
  const [enrollModalOpen, setEnrollModalOpen] = useState(false)

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
        <h2 className="mb-4 text-xl font-bold">Daftar Kursus</h2>
        {courses.length === 0 ? (
          <p className="text-ink-soft">Belum ada kursus.</p>
        ) : (
          <div className="space-y-4">
            {courses.map((course) => (
              <CourseAdminItem
                key={course.id}
                course={course}
                deleteCourse={deleteCourse}
                deleteQuiz={deleteQuiz}
                onManageStudents={(id) => {
                  setSelectedCourseId(id)
                  setEnrollModalOpen(true)
                }}
              />
            ))}
          </div>
        )}
      </section>

      <StudentEnrollModal
        courseId={selectedCourseId}
        isOpen={enrollModalOpen}
        onClose={() => setEnrollModalOpen(false)}
      />
    </div>
  )
}