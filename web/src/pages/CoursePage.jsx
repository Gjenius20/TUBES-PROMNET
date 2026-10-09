import { Link, useParams } from '@tanstack/react-router'
import { useCourse } from '../hooks/useCourses'
import { getErrorMessage } from '../services/api'

export default function CoursePage() {
  const { courseId } = useParams({ strict: false })
  const { data: course, isLoading, isError, error } = useCourse(courseId)

  if (isLoading) return <p className="text-ink-soft">Memuat…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  const quizzes = course.quizzes ?? []

  return (
    <div>
      <Link to="/courses" className="text-sm text-signal">
        ← Semua kursus
      </Link>
      <h1 className="mt-2 text-3xl font-bold">{course.title}</h1>
      <p className="mt-1 max-w-2xl text-ink-soft">{course.description}</p>

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
    </div>
  )
}
