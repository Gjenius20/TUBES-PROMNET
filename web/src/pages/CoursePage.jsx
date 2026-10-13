import { Link, useParams, useNavigate } from '@tanstack/react-router'
import { useAuthStore } from '../store/authStore'
import { useCourse, useEnroll, useStudentEnrollments } from '../hooks/useCourses'
import { getErrorMessage } from '../services/api'

export default function CoursePage() {
  const { courseId } = useParams({ strict: false })
  const navigate = useNavigate()
  const { data: course, isLoading, isError, error } = useCourse(courseId)
  const { user } = useAuthStore()
  const { data: enrollments } = useStudentEnrollments()
  const enrollMutation = useEnroll(courseId)

  if (isLoading) return <p className="text-ink-soft">Memuat…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  // Cek status enrollment
  const enrollment = enrollments?.find(e => e.course_id === Number(courseId))
  const isEnrolled = enrollment?.status === 'APPROVED'
  const isPending = enrollment?.status === 'PENDING'
  const isRejected = enrollment?.status === 'REJECTED'

  const handleEnroll = async () => {
    try {
      await enrollMutation.mutateAsync()
      // Refresh will be handled by query invalidation
    } catch (err) {
      console.error('Enrollment failed:', err)
    }
  }

  const quizzes = course.quizzes ?? []

  return (
    <div>
      <Link to="/courses" className="text-sm text-signal">
        ← Semua kursus
      </Link>
      <h1 className="mt-2 text-3xl font-bold">{course.title}</h1>
      <p className="mt-1 max-w-2xl text-ink-soft">{course.description}</p>

      {/* Enrollment Status */}
      <div className="mt-6 panel p-4">
        {!enrollment ? (
          <button
            onClick={handleEnroll}
            disabled={enrollMutation.isPending}
            className="btn-primary"
          >
            {enrollMutation.isPending ? 'Mendaftar...' : 'Daftar Kursus'}
          </button>
        ) : isPending ? (
          <span className="badge badge-warning">Menunggu Persetujuan Admin</span>
        ) : isRejected ? (
          <span className="badge badge-error">Ditolak oleh Admin</span>
        ) : (
          <Link to={`/courses/${courseId}/quizzes`} className="btn-primary">
            Masuk Kursus
          </Link>
        )}
      </div>

      {isEnrolled && (
        <>
          <h2 className="mb-3 mt-8 text-xl font-bold">Soal</h2>
          {quizzes.length === 0 ? (
            <p className="text-ink-soft">Belum ada soal di kursus ini.</p>
          ) : (
            <ul className="panel divide-y divide-rule">
              {quizzes.map((quiz) => (
                <li key={quiz.id}>
                  <Link
                    to="/quizzes/$quizId"
                    params={{ quizId: String(quiz.id) }}
                    className="flex items-center justify-between px-5 py-3 hover:bg-paper"
                  >
                    <span className="font-medium">{quiz.title}</span>
                    <span className="text-sm text-signal">Kerjakan</span>
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </>
      )}
    </div>
  )
}
