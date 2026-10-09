// Satu sel per testcase: hijau = lulus, merah = gagal. Elemen visual utama hasil penilaian.
export default function VerdictStrip({ results }) {
  return (
    <ol className="flex flex-wrap gap-1.5" aria-label="Hasil per testcase">
      {results.map((r) => (
        <li
          key={r.index}
          title={`Testcase ${r.index}${r.hidden ? ' (tersembunyi)' : ''}: ${r.status}`}
          className={`flex h-9 min-w-9 items-center justify-center rounded-md px-2 text-xs font-semibold text-white ${
            r.passed ? 'bg-pass' : 'bg-fail'
          } ${r.hidden ? 'opacity-80' : ''}`}
        >
          {r.index}
        </li>
      ))}
    </ol>
  )
}
