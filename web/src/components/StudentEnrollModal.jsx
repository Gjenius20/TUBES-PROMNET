import { useState, useEffect, useRef } from 'react'
import { useStudents, useBulkEnroll } from '../hooks/useUsers'

export default function StudentEnrollModal({ courseId, isOpen, onClose }) {
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [selectedIds, setSelectedIds] = useState(new Set());

  const modalRef = useRef(null);
  const searchInputRef = useRef(null);
  const previousActiveElement = useRef(null);

  const { data, isLoading, isFetching, isError, error } = useStudents(search, page);
  const bulkEnroll = useBulkEnroll(courseId);

  // Focus Management & Scroll Lock
  useEffect(() => {
    if (isOpen) {
      previousActiveElement.current = document.activeElement;
      document.body.style.overflow = 'hidden';
      setTimeout(() => searchInputRef.current?.focus(), 100);
    } else {
      document.body.style.overflow = '';
      previousActiveElement.current?.focus();
    }
    return () => { document.body.style.overflow = ''; };
  }, [isOpen]);

  if (!isOpen) return null;

  const toggleSelect = (id) => {
    const next = new Set(selectedIds);
    next.has(id) ? next.delete(id) : next.add(id);
    setSelectedIds(next);
  };

  const toggleSelectAll = () => {
    if (selectedIds.size === data?.users?.length && data?.users?.length > 0) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(data?.users?.map(u => u.id) || []));
    }
  };

  const handleEnroll = async () => {
    await bulkEnroll.mutateAsync(Array.from(selectedIds));
    setSelectedIds(new Set());
    onClose();
  };

  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
      onClick={onClose}
    >
      <div
        className="bg-slate-900 p-6 rounded-lg w-full max-w-2xl border border-slate-600 shadow-xl"
        onClick={e => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        onKeyDown={(e) => e.key === 'Escape' && onClose()}
      >
        <h2 id="modal-title" className="text-xl font-bold text-slate-100 mb-4">Tambahkan Siswa</h2>

        <label htmlFor="student-search" className="sr-only">Cari Siswa</label>
        <input
          id="student-search"
          ref={searchInputRef}
          type="search"
          placeholder="Cari NIM, Nama, atau Email..."
          aria-label="Cari siswa"
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="w-full px-4 py-2 bg-slate-800 border border-slate-600 rounded-lg text-slate-100 placeholder-slate-400 mb-4 focus:ring-2 focus:ring-cyan-400"
        />

        {isError ? (
          <p className="text-fail p-4">Gagal memuat data: {error.message}</p>
        ) : (
          <div className="overflow-y-auto max-h-[50vh] border border-slate-600 rounded-lg">
            <table className="w-full text-sm" role="grid">
              <thead className="bg-slate-800 sticky top-0">
                <tr>
                  <th scope="col" className="px-4 py-2 w-12">
                    <input type="checkbox" aria-label="Pilih semua" onChange={toggleSelectAll} checked={data?.users?.length > 0 && selectedIds.size === data.users.length} />
                  </th>
                  <th scope="col" className="px-4 py-2 text-left">NIM</th>
                  <th scope="col" className="px-4 py-2 text-left">Nama</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-700">
                {isLoading ? <tr><td colSpan={3} className="text-center p-4">Memuat...</td></tr> :
                 data?.users.map(u => (
                  <tr key={u.id} className="hover:bg-slate-800">
                    <td className="px-4 py-2"><input type="checkbox" aria-label={`Pilih ${u.name}`} checked={selectedIds.has(u.id)} onChange={() => toggleSelect(u.id)} /></td>
                    <td className="px-4 py-2 text-slate-300">{u.nim}</td>
                    <td className="px-4 py-2 text-slate-100">{u.name}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="flex justify-end gap-2 mt-4 pt-4 border-t border-slate-700">
          <button onClick={onClose} className="px-4 py-2 bg-slate-800 hover:bg-slate-700 rounded-lg text-sm">Batal</button>
          <button onClick={handleEnroll} disabled={selectedIds.size === 0 || bulkEnroll.isPending} className="px-4 py-2 bg-cyan-600 hover:bg-cyan-700 rounded-lg text-white text-sm">
            {bulkEnroll.isPending ? 'Memproses...' : `Enroll ${selectedIds.size} Siswa`}
          </button>
        </div>
      </div>
    </div>
  )
}