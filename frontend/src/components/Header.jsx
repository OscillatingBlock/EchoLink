import ThemeToggle from './ThemeToggle'

export default function Header()
{
  return (
    <header className="sticky top-0 z-10 border-b border-slate-200 bg-white/80 backdrop-blur dark:border-[#2e2e2e] dark:bg-[#0c0c0c]">
      <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4">
        <div className="text-sm text-slate-500 dark:text-slate-400">Dashboard</div>
        <div className="flex items-center gap-3">
          <input
            placeholder="Search bots…"
            className="w-64 rounded-md border border-slate-200 bg-slate-100/40 px-3 py-1.5 text-sm text-[#0c0c0c] outline-none transition focus:ring-2 focus:ring-indigo-500/40 dark:border-slate-700 dark:bg-[#363636] dark:text-slate-100"
          />
          <ThemeToggle />
          <div className="grid size-8 place-items-center rounded-full bg-indigo-100 font-semibold text-indigo-700 dark:bg-[#363636] dark:text-indigo-300">
            VC
          </div>
        </div>
      </div>
    </header>
  )
}
