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

  return { data, isLoading, error, fetchProgress }
}
