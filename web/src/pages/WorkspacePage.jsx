import { useState } from 'react'
import { useParams } from '@tanstack/react-router'
import CodeEditor from '../components/CodeEditor'
import StatusBadge from '../components/StatusBadge'
import VerdictStrip from '../components/VerdictStrip'
import { useQuiz, useSubmissions, useSubmit } from '../hooks/useQuizzes'
import { getErrorMessage } from '../services/api'

const STARTER_CODE = `#include <stdio.h>

int main(void) {
    return 0;
}
`

function ResultPanel({ submission }) {
  const results = submission.results ?? []
  const visibleFailures = results.filter((r) => !r.hidden && !r.passed)

  return (
    <div className="panel space-y-3 p-4" aria-live="polite">
      <div className="flex flex-wrap items-center gap-3">
        <StatusBadge status={submission.status} />
        <span className="text-sm text-ink-soft">
          Nilai <strong className="text-ink">{submission.score}</strong> · lulus {submission.passed_count}/
          {submission.total_count} testcase
        </span>
      </div>

      {results.length > 0 ? <VerdictStrip results={results} /> : null}

      {submission.status === 'COMPILE_ERROR' ? (
        <pre className="max-h-60 overflow-auto rounded-md bg-ink p-3 text-xs text-paper">
          {submission.error_message}
        </pre>
      ) : null}

      {visibleFailures.map((r) => (
        <div key={r.index} className="rounded-md border border-rule p-3 text-xs">
          <p className="mb-1 font-semibold">Testcase {r.index}: {r.status}</p>
          {r.stdin ? <pre className="overflow-auto">Input:{'\n'}{r.stdin}</pre> : null}
          <pre className="overflow-auto">Diharapkan:{'\n'}{r.expected_output}</pre>
          {r.actual_output !== undefined ? <pre className="overflow-auto">Output Anda:{'\n'}{r.actual_output}</pre> : null}
        </div>
      ))}
    </div>
  )
}

function EmptySubmission() {
  return (
    <div className="panel p-6 text-center text-ink-soft">
      <p>Belum ada submission untuk soal ini.</p>
    </div>
  )
}

function EditorPanel({ quiz, submitMutation }) {
  const [code, setCode] = useState(STARTER_CODE)

  const handleSubmit = async (e) => {
    e.preventDefault()
    await submitMutation.mutateAsync({ sourceCode: code })
  }

  return (
    <div className="space-y-4">
      <form onSubmit={handleSubmit} className="space-y-3">
        <div>
          <label htmlFor="editor" className="block text-sm font-medium mb-1">
            Editor C
          </label>
          <CodeEditor value={code} onChange={setCode} />
        </div>

        <button
          type="submit"
          className="btn-primary w-full"
          disabled={submitMutation.isPending}
        >
          {submitMutation.isPending ? 'Mengirim...' : 'Kirim & Jalankan'}
        </button>
      </form>

      {submitMutation.isError && (
        <p role="alert" className="text-fail text-sm">
          {getErrorMessage(submitMutation.error)}
        </p>
      )}

      {submitMutation.isSuccess && (
        <p role="status" className="text-success text-sm">
          Submission berhasil dikirim!
        </p>
      )}
    </div>
  )
}

export default function WorkspacePage() {
  const { quizId } = useParams({ strict: false })
  const { data: quiz, isLoading, isError, error } = useQuiz(quizId)
  const { data: submissions } = useSubmissions(quizId)
  const submitMutation = useSubmit(quizId)

  if (isLoading) return <p className="p-8 text-ink-soft">Memuat soal...</p>
  if (isError) return <p role="alert" className="p-8 text-fail">{getErrorMessage(error)}</p>
  if (!quiz) return <p className="p-8 text-ink-soft">Soal tidak ditemukan.</p>

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Detail Soal */}
      <div className="panel p-4 space-y-4">
        <h2 className="text-2xl font-bold">{quiz.title}</h2>
        <div className="prose text-sm text-ink-soft">{quiz.description}</div>

        {quiz.constraints && (
          <div className="bg-ink/5 p-3 rounded text-sm">
            <strong>Batasan:</strong> {quiz.constraints}
          </div>
        )}

        {quiz.sampleTestCases?.length > 0 && (
          <div>
            <h3 className="font-semibold mb-2">Contoh Testcase</h3>
            <div className="space-y-2 text-xs font-mono">
              {quiz.sampleTestCases.map((tc, i) => (
                <div key={i} className="bg-ink/5 p-2 rounded">
                  <strong>Input:</strong>
                  <pre className="whitespace-pre-wrap">{tc.stdin}</pre>
                  <strong>Output:</strong>
                  <pre className="whitespace-pre-wrap">{tc.expectedOutput}</pre>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Editor & Hasil */}
      <div className="space-y-4">
        <EditorPanel quiz={quiz} submitMutation={submitMutation} />

        {/* Riwayat Submission */}
        <div className="panel p-4">
          <h3 className="font-semibold mb-3">Riwayat Submission</h3>
          {submissions?.length > 0 ? (
            <div className="space-y-2 max-h-96 overflow-auto">
              {submissions.map((sub) => (
                <div key={sub.id} className="border border-rule rounded p-2">
                  <div className="flex items-center gap-2">
                    <StatusBadge status={sub.status} />
                    <span className="text-sm">
                      Skor: <strong>{sub.score}</strong> · {sub.passed_count}/{sub.total_count} testcase
                    </span>
                    <span className="text-xs text-ink-soft ml-auto">
                      {new Date(sub.created_at).toLocaleString()}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <EmptySubmission />
          )}
        </div>
      </div>
    </div>
  )
}