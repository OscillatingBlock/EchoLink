import {useNavigate} from 'react-router-dom'

export default function CreateBot()
{
  const navigate=useNavigate()
  return (
    <form
      className="mx-auto m-5 h-[80%]  space-y-5 rounded-xl border border-slate-200 bg-white p-6 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]"
      onSubmit={(e) =>
      {
        e.preventDefault();
        navigate('/bots')
      }}
    >
      <div>
        <h1 className="text-2xl font-semibold">Create Bot</h1>
        <p className="text-sm text-slate-500 dark:text-slate-400">Define your bot’s basics</p>
      </div>
      <label className="block text-sm">
        Name
        <input className="mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100" placeholder="e.g. Sales Assistant" />
      </label>
      <label className="block text-sm">
        Description
        <textarea rows={4} className="mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100" placeholder="What does this bot do?" />
      </label>
      <div className="grid gap-4 sm:grid-cols-2">
        <label className="block text-sm">
          Voice
          <select className="mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100">
            <option>Alloy</option>
            <option>Amber</option>
          </select>
        </label>
        <label className="block text-sm">
          Language
          <select className="mt-1 w-full rounded-md border border-slate-300 bg-white px-3 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100">
            <option>English (US)</option>
            <option>English (UK)</option>
          </select>
        </label>
      </div>
      <div className="flex items-center justify-end gap-3">
        <button type="button" onClick={() => navigate('/bots')} className="rounded-md border border-slate-300 px-3 py-2 text-sm dark:border-slate-700">Cancel</button>
        <button type="submit" className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-700">Create</button>
      </div>
    </form>
  )
}
