function CourseAdminItem({ course, deleteCourse, deleteQuiz }) {
  const { data, fetchProgress, exportExcel, isLoading } = useCourseProgress(course.id)
  const exportGrades = useExportGrades(course.id)
  const [showProgress, setShowProgress] = useState(false)

  const handleToggleProgress = () => {
    if (showProgress) {
      setShowProgress(false)
    } else {
      setShowProgress(true)
      fetchProgress()
    }
  }

  return (
    <div key={course.id} className="panel p-4">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h3 className="text-lg font-bold">{course.title}</h3>
          <p className="text-sm text-ink-soft">{course.description}</p>
        </div>
        <div className="flex gap-2">
          <button onClick={handleToggleProgress} className="btn-secondary text-xs" disabled={isLoading}>
            {showProgress ? 'Sembunyikan Progress' : 'Tampilkan Progress'}
          </button>
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

      {showProgress && (
        <div className="mt-4 text-xs">
          {isLoading ? (
            <p className="text-center text-ink-soft py-4">Memuat progress...</p>
          ) : data ? (
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
          ) : null}
        </div>
      )}
