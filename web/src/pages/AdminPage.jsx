import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link } from '@tanstack/react-router'
import FormField from '../components/FormField'
import { useCourseProgress } from '../hooks/useProgress'
import { useCourses, useCreateCourse, useDeleteCourse, useAdminEnrollments, useUpdateEnrollmentStatus, useExportGrades } from '../hooks/useCourses'
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

function CourseAdminItem({ course, deleteCourse, deleteQuiz }) {
  const { data, fetchProgress, exportExcel } = useCourseProgress(course.id)
  const exportGrades = useExportGrades(course.id)

  return (
    <div key={course.id} className="panel p-4">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h3 className="text-lg font-bold">{course.title}</h3>
          <p className="text-sm text-ink-soft">{course.description}</p>
        </div>
        <div className="flex gap-2">
          <button onClick={fetchProgress} className="btn-secondary text-xs">Progress</button>
          <button onClick={exportExcel} className="btn-primary text-xs">Excel Progress</button>
          <button
            onClick={() => exportGrades.mutate()}
            disabled={exportGrades.isPending}
            className="btn-success text-xs"
          >
            {exportGrades.isPending ? 'Mengunduh...' : 'Unduh Nilai (Excel)'}
          </button>
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
      </div>
      
      {data && (
        <div className="mt-4 text-xs">
          <table className="w-full border-collapse">
            <thead>
              <tr className="border-b border-rule">
                <th className="text-left p-1">NIM</th>
                <th className="text-left p-1">Nama</th>
                <th className="text-left p-1">Status</th>
                <th className="text-right p-1">Progress</th>
              </tr>
            </thead>
            <tbody>
              {data.map(p => (
                <tr key={p.id} className="border-b border-rule">
                  <td className="p-1">{p.nim}</td>
                  <td className="p-1">{p.name}</td>
                  <td className="p-1">{p.status}</td>
                  <td className="text-right p-1">{p.progress_pct}%</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

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
  )
}

export default function AdminPage() {
  const { data: courses, isLoading, isError, error } = useCourses()
  const deleteCourse = useDeleteCourse()
  const deleteQuiz = useDeleteQuiz()
  const { data: enrollments, isLoading: enrollLoading } = useAdminEnrollments('PENDING')
  const updateEnrollment = useUpdateEnrollmentStatus()
  const [activeTab, setActiveTab] = useState('courses')
  const [enrollmentStatus, setEnrollmentStatus] = useState('PENDING')

  if (isLoading) return <p className="text-ink-soft">Memuat…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  return (
    <div className="space-y-8">
      <h1 className="text-3xl font-bold">Dashboard Admin</h1>

      {/* Tab Navigation */}
      <div className="border-b border-rule">
        <nav className="flex space-x-8" aria-label="Tabs">
          <button
            onClick={() => setActiveTab('courses')}
            className={`py-3 px-1 border-b-2 font-medium text-sm ${activeTab === 'courses' ? 'border-signal text-signal' : 'border-transparent text-ink-soft hover:text-ink'}`}
          >
            Kelola Kursus & Soal
          </button>
          <button
            onClick={() => setActiveTab('enrollments')}
            className={`py-3 px-1 border-b-2 font-medium text-sm ${activeTab === 'enrollments' ? 'border-signal text-signal' : 'border-transparent text-ink-soft hover:text-ink'}`}
          >
            Kelola Pendaftaran
          </button>
        </nav>
      </div>

      {activeTab === 'courses' && (
        <>
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
                <CourseAdminItem
                  key={course.id}
                  course={course}
                  deleteCourse={deleteCourse}
                  deleteQuiz={deleteQuiz}
                />
              ))}
            </div>
          </section>
        </>
      )}

      {activeTab === 'enrollments' && (
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold">Kelola Pendaftaran Kursus</h2>
            <select
              value={enrollmentStatus}
              onChange={(e) => setEnrollmentStatus(e.target.value)}
              className="input w-auto"
            >
              <option value="PENDING">Menunggu (Pending)</option>
              <option value="APPROVED">Disetujui (Approved)</option>
              <option value="REJECTED">Ditolak (Rejected)</option>
            </select>
          </div>

          {enrollLoading ? (
            <p className="text-ink-soft">Memuat…</p>
          ) : (
            <div className="panel overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-rule">
                    <th className="text-left p-3">Siswa</th>
                    <th className="text-left p-3">NIM</th>
                    <th className="text-left p-3">Email</th>
                    <th className="text-left p-3">Kursus</th>
                    <th className="text-left p-3">Status</th>
                    <th className="text-left p-3">Tanggal Daftar</th>
                    <th className="text-right p-3">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {enrollments.length === 0 ? (
                    <tr>
                      <td colSpan={7} className="text-center p-8 text-ink-soft">Tidak ada pengajuan pendaftaran</td>
                    </tr>
                  ) : (
                    enrollments.map((enrollment) => (
                      <tr key={enrollment.id} className="border-b border-rule">
                        <td className="p-3">{enrollment.user?.name}</td>
                        <td className="p-3">{enrollment.user?.nim}</td>
                        <td className="p-3">{enrollment.user?.email}</td>
                        <td className="p-3">{enrollment.course?.title}</td>
                        <td className="p-3">
                          <span className={`badge ${
                            enrollment.status === 'APPROVED' ? 'badge-success' :
                            enrollment.status === 'PENDING' ? 'badge-warning' :
                            'badge-error'
                          }`}>
                            {enrollment.status}
                          </span>
                        </td>
                        <td className="p-3">{new Date(enrollment.created_at).toLocaleDateString('id-ID')}</td>
                        <td className="p-3 text-right">
                          {enrollment.status === 'PENDING' && (
                            <div className="flex items-center justify-end gap-2">
                              <button
                                onClick={() => updateEnrollment.mutate({ id: enrollment.id, status: 'APPROVED' })}
                                disabled={updateEnrollment.isPending}
                                className="btn-success btn-sm"
                              >
                                Setujui
                              </button>
                              <button
                                onClick={() => updateEnrollment.mutate({ id: enrollment.id, status: 'REJECTED' })}
                                disabled={updateEnrollment.isPending}
                                className="btn-error btn-sm"
                              >
                                Tolak
                              </button>
                            </div>
                          )}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          )}
        </section>
      )}
    </div>
  )
}
