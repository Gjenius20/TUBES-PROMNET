import { Link } from '@tanstack/react-router'
import { useCourses } from '../hooks/useCourses'
import { getErrorMessage } from '../services/api'

export default function CoursesPage() {
  const { data, isLoading, isError, error } = useCourses()

  if (isLoading) return <p className="text-ink-soft">Memuat kursus…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  return (
    <div>
      <h1 className="mb-6 text-3xl font-bold">Kursus</h1>
      {data.length === 0 ? (
        <p className="text-ink-soft">Belum ada kursus.</p>
      ) : (
        <ul className="grid gap-4 sm:grid-cols-2">
          {data.map((course) => (
            <li key={course.id}>
              <Link
                to="/courses/$courseId"
                params={{ courseId: String(course.id) }}
                className="panel block p-5 hover:border-ink"
              >
                <h2 className="text-xl font-bold">{course.title}</h2>
                <p className="mt-1 line-clamp-2 text-sm text-ink-soft">{course.description}</p>
                <p className="mt-3 text-xs text-ink-soft">{(course.quizzes ?? []).length} soal</p>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
