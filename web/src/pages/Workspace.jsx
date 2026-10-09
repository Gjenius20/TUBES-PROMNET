import { useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import Editor from '@monaco-editor/react'
import { api, errorMessage } from '../services/api'

const TEMPLATE = `#include <stdio.h>

int main(void) {
    return 0;
}
`

export default function Workspace() {
  const [activeId, setActiveId] = useState(null)
  const [code, setCode] = useState(TEMPLATE)

  const quizzes = useQuery({ queryKey: ['quizzes'], queryFn: () => api.get('/quizzes').then((r) => r.data.data) })
  const active = quizzes.data?.find((q) => q.id === activeId)

  const submit = useMutation({
    mutationFn: () => api.post('/submissions', { quiz_id: activeId, source_code: code }).then((r) => r.data.data),
  })

  return (
    <div className="split">
      <aside className="side">
        <h2>Problems</h2>
        {quizzes.isLoading && <p className="muted">Loading problems…</p>}
        {quizzes.data?.length === 0 && <p className="muted">No problems yet. An admin needs to add one.</p>}
        {quizzes.data?.map((q) => (
          <button key={q.id} className={q.id === activeId ? 'item on' : 'item'} onClick={() => { setActiveId(q.id); setCode(TEMPLATE); submit.reset() }}>
            {q.title}
          </button>
        ))}
      </aside>
      <section className="work">
        {!active ? (
          <p className="muted pad">Choose a problem to start coding.</p>
        ) : (
          <>
            <div className="statement">
              <h2>{active.title}</h2>
              <p>{active.statement}</p>
            </div>
            <div className="editor">
              <Editor height="100%" language="c" value={code} onChange={(v) => setCode(v ?? '')} options={{ minimap: { enabled: false }, fontSize: 14 }} />
            </div>
            <div className="result">
              <button onClick={() => submit.mutate()} disabled={submit.isPending}>{submit.isPending ? 'Running tests…' : 'Run and submit'}</button>
              {submit.isError && <span className="err">{errorMessage(submit.error)}</span>}
              {submit.data && (
                <span className={submit.data.status === 'PASSED' ? 'ok' : 'err'}>
                  {submit.data.status.replace('_', ' ')} · score {Math.round(submit.data.score)}
                </span>
              )}
            </div>
          </>
        )}
      </section>
    </div>
  )
}
