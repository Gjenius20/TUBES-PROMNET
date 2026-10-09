import Editor from '@monaco-editor/react'

// Wrapper Monaco Editor dengan dukungan sintaks bahasa C.
export default function CodeEditor({ value, onChange, height = '420px' }) {
  return (
    <div className="overflow-hidden rounded-lg border border-ink">
      <Editor
        height={height}
        language="c"
        theme="vs-dark"
        value={value}
        onChange={(v) => onChange(v ?? '')}
        options={{
          fontSize: 14,
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 4,
        }}
      />
    </div>
  )
}
