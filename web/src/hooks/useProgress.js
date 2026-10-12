import { useState } from 'react'
import { api } from '../services/api'

export function useCourseProgress(courseId) {
  const [data, setData] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState(null)

  const fetchProgress = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const res = await api.get(`/admin/courses/${courseId}/progress`)
      setData(res.data.data.progress)
    } catch (err) {
      setError(err)
    } finally {
      setIsLoading(false)
    }
  }

  const exportExcel = async () => {
    try {
      const res = await api.get(`/admin/courses/${courseId}/export-excel`, {
        responseType: 'blob',
      })
      const url = window.URL.createObjectURL(new Blob([res.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `course-report-${courseId}.xlsx`)
      document.body.appendChild(link)
      link.click()
      link.remove()
    } catch (err) {
      setError(err)
    }
  }

  return { data, isLoading, error, fetchProgress, exportExcel }
}
