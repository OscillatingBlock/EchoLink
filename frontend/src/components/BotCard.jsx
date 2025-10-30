export default function BotCard({bot, onClick})
{
  return (
    <button
      onClick={onClick}
      className="group w-full rounded-lg border border-slate-200 bg-white p-4 text-left shadow-sm transition hover:shadow-md bg-gradient-to-r  from-[#010101] to-[#22094f]"
    >
      <div className="flex items-start justify-between">
        <div>
          <div className="text-base font-semibold text-slate-900 group-hover:text-indigo-600 dark:text-slate-100 dark:group-hover:text-indigo-300">
            {bot.name}
          </div>
          <div className="mt-1 text-sm text-slate-500 dark:text-slate-400">{bot.description}</div>
        </div>
        <span className="ml-3 rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700">
          {bot.status}
        </span>
      </div>
      <div className="mt-3 flex items-center gap-4 text-xs text-slate-500 dark:text-slate-400">
        <div>Calls: {bot.calls}</div>
        <div>Success: {bot.successRate}%</div>
      </div>
    </button>
  )
}
