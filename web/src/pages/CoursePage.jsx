import { Link, useParams } from '@tanstack/react-router'
import { useCourse, useStudentEnrollments } from '../hooks/useCourses'
import { getErrorMessage } from '../services/api'

export default function CoursePage() {
  const { courseId } = useParams({ strict: false })
  const { data: course, isLoading, isError, error } = useCourse(courseId)
  const { data: enrollments, isLoading: enrollLoading } = useStudentEnrollments()

  if (isLoading || enrollLoading) return <p className="p-8 text-ink-soft">Memuat data kursus...</p>
  if (isError) return <p role="alert" className="p-8 text-fail">{getErrorMessage(error)}</p>
  if (!course) return <p className="p-8 text-ink-soft">Kursus tidak ditemukan.</p>

  // Cek status enrollment
  const enrollment = enrollments?.find(e => e.course_id === Number(courseId))
  const isEnrolled = enrollment?.status === 'APPROVED'

  return (
    <div>
      <Link to="/courses" className="text-sm text-signal">
        ← Semua kursus
      </Link>
      <h1 className="mt-2 text-3xl font-bold">{course.title}</h1>
      <p className="mt-1 max-w-2xl text-ink-soft">{course.description}</p>

      {/* Enrollment Status */}
      {!isEnrolled ? (
        <div className="mt-6 panel p-4">
          <p className="text-sm">
            {enrollment 
              ? `Status Pendaftaran: ${enrollment.status}` 
              : "Anda belum terdaftar di kursus ini. Silakan hubungi admin untuk pendaftaran."}
          </p>
        </div>
      ) : (
        <div className="mt-6">
          <h2 className="text-xl font-bold mb-4">Daftar Soal</h2>
          {(course.quizzes ?? []).length === 0 ? (
            <p className="text-ink-soft">Belum ada soal untuk kursus ini.</p>
          ) : (
            <ul className="grid gap-2">
              {course.quizzes.map((quiz) => (
                <li key={quiz.id}>
                  <Link
                    to="/quizzes/$quizId"
                    params={{ quizId: String(quiz.id) }}
                    className="panel block p-3 hover:border-ink"
                  >
                    {quiz.title}
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
}
