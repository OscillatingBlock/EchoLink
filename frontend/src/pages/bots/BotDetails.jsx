import {useParams, Link} from 'react-router-dom'

export default function BotDetails()
{
  const {id}=useParams()
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Bot: {id}</h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">Configuration and performance</p>
        </div>
        <Link to="/bots" className="text-sm text-indigo-600 hover:underline dark:text-indigo-400">Back to Bots</Link>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <div className="space-y-4 md:col-span-2">
          <section className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]">
            <h2 className="mb-2 text-lg font-semibold">Voice Settings</h2>
            <div className="grid gap-3 sm:grid-cols-2">
              <label className="text-sm">Voice
                <select className="mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100">
                  <option>Alloy</option>
                  <option>Amber</option>
                </select>
              </label>
              <label className="text-sm">Language
                <select className="mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100">
                  <option>English (US)</option>
                  <option>English (UK)</option>
                </select>
              </label>
            </div>
          </section>

          <section className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]">
            <h2 className="mb-2 text-lg font-semibold">Scripts</h2>
            <textarea rows={6} className="w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100" placeholder="Intro script..." />
          </section>
        </div>

        <aside className="space-y-4">
          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]">
            <div className="text-sm text-slate-500 dark:text-slate-400">Status</div>
            <div className="mt-1 inline-flex items-center gap-2 rounded-full bg-green-100 px-2.5 py-1 text-xs font-medium text-green-700">Active</div>
          </div>
          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]">
            <div className="text-sm text-slate-500 dark:text-slate-400">Today</div>
            <div className="mt-1 text-2xl font-semibold">42 calls</div>
          </div>
        </aside>
      </div>
    </div>
  )
}
