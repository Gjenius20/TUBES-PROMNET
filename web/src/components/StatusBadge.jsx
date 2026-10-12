const STYLES = {
  PASSED: { label: 'Lulus', cls: 'bg-pass/10 text-pass border-pass' },
  WRONG_ANSWER: { label: 'Jawaban salah', cls: 'bg-fail/10 text-fail border-fail' },
  COMPILE_ERROR: { label: 'Gagal kompilasi', cls: 'bg-warn/10 text-warn border-warn' },
  RUNTIME_ERROR: { label: 'Error saat berjalan', cls: 'bg-warn/10 text-warn border-warn' },
}

export default function StatusBadge({ status }) {
  const s = STYLES[status] ?? { label: status, cls: 'border-rule text-ink-soft' }
  return (
    <span className={`inline-block rounded border px-2 py-0.5 text-xs font-semibold ${s.cls}`}>
      {s.label}
    </span>
  )
}
