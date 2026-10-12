import { useState } from 'react'
import { Link, useParams } from '@tanstack/react-router'
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
          {r.actual_output !== undefined ? <pre className="overflow-auto">Output kamu:{'\n'}{r.actual_output}</pre> : null}
          {r.message ? <pre className="overflow-auto text-fail">{r.message}</pre> : null}
        </div>
      ))}
    </div>
  )
}

export default function WorkspacePage() {
  const { quizId } = useParams({ strict: false })
  const { data: quiz, isLoading, isError, error } = useQuiz(quizId)
  const { data: history } = useSubmissions(quizId)
  const submit = useSubmit(quizId)
  const [code, setCode] = useState(STARTER_CODE)

  if (isLoading) return <p className="text-ink-soft">Memuat soal…</p>
  if (isError) return <p role="alert" className="text-fail">{getErrorMessage(error)}</p>

  return (
    <div>
      <Link to="/courses/$courseId" params={{ courseId: String(quiz.course_id) }} className="text-sm text-signal">
        ← Kembali ke kursus
      </Link>
      <div className="mt-3 grid gap-6 lg:grid-cols-[minmax(0,2fr)_minmax(0,3fr)]">
        <section className="space-y-4">
          <h1 className="text-2xl font-bold">{quiz.title}</h1>
          <p className="whitespace-pre-wrap text-sm leading-relaxed">{quiz.description}</p>
          {quiz.constraints ? (
            <div>
              <h2 className="mb-1 text-sm font-bold">Batasan</h2>
              <p className="whitespace-pre-wrap text-sm text-ink-soft">{quiz.constraints}</p>
            </div>
          ) : null}
          {quiz.sample_test_cases.map((tc, i) => (
            <div key={i} className="panel p-3 text-xs">
              <p className="mb-1 font-semibold">Contoh {i + 1}</p>
              <pre className="overflow-auto">Input:{'\n'}{tc.stdin || '(kosong)'}</pre>
              <pre className="mt-2 overflow-auto">Output:{'\n'}{tc.expected_output}</pre>
            </div>
          ))}
        </section>

        <section className="space-y-4">
          <CodeEditor value={code} onChange={setCode} />
          <button
            type="button"
            className="btn-primary"
            disabled={submit.isPending || code.trim() === ''}
            onClick={() => submit.mutate({ sourceCode: code })}
          >
            {submit.isPending ? 'Menilai…' : 'Kirim jawaban'}
          </button>

          {submit.isError ? (
            <p role="alert" className="text-sm text-fail">{getErrorMessage(submit.error)}</p>
          ) : null}
          {submit.data ? <ResultPanel submission={submit.data} /> : null}

          {history && history.length > 0 ? (
            <div>
              <h2 className="mb-2 text-sm font-bold">Riwayat pengiriman</h2>
              <ul className="panel divide-y divide-rule text-sm">
                {history.map((s) => (
                  <li key={s.id} className="flex items-center justify-between px-3 py-2">
                    <StatusBadge status={s.status} />
                    <span className="text-ink-soft">
                      {s.score} · {new Date(s.created_at).toLocaleString('id-ID')}
                    </span>
                  </li>
                ))}
              </ul>
            </div>
          ) : null}
        </section>
      </div>
    </div>
  )
}
