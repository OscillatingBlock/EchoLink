import {useEffect} from 'react'
import {useDispatch, useSelector} from 'react-redux'
import {fetchMyNumber} from '../../features/auth/authSlice'

export default function Home()
{
  const dispatch=useDispatch()
  const {access_token, phone_number, bots_count}=useSelector(s => s.auth)

  useEffect(() =>
  {
    if (access_token) dispatch(fetchMyNumber())
  }, [access_token, dispatch])

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Overview</h1>
        <p className="text-sm text-slate-500 dark:text-slate-400">Key metrics for your voice bots</p>
      </div>
      {access_token&&(
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm  bg-gradient-to-r  from-[#010101] to-[#22094f]">
            <div className="text-sm text-slate-500 dark:text-slate-400">My Number</div>
            <div className="mt-1 text-2xl font-semibold">{phone_number||'—'}</div>
          </div>
          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]">
            <div className="text-sm text-slate-500 dark:text-slate-400">Bots</div>
            <div className="mt-1 text-2xl font-semibold">{typeof bots_count==='number'? bots_count:'—'}</div>
          </div>
        </div>
      )}

      <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[
          {label: 'Total Bots', value: 8},
          {label: 'Active Calls', value: 12},
          {label: 'Success Rate', value: '92%'},
          {label: 'Avg Call Time', value: '3m 42s'},
        ].map((card) => (
          <div
            key={card.label}
            className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:shadow-md bg-gradient-to-r  from-[#010101] to-[#22094f]"
          >
            <div className="text-sm text-slate-500 dark:text-slate-400">{card.label}</div>
            <div className="mt-1 text-2xl font-semibold">{card.value}</div>
          </div>
        ))}
      </section>

      <section className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm bg-gradient-to-r  from-[#010101] to-[#22094f]">
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-lg font-semibold">Recent Activity</h2>
          <a href="/bots" className="text-sm text-indigo-600 hover:underline dark:text-indigo-400">View all</a>
        </div>
        <ul className="divide-y divide-slate-200 text-sm dark:divide-slate-800">
          {[
            'Deployed Sales Assistant bot',
            'Updated Voice Parameters for Support Bot',
            'Created new Outreach bot',
          ].map((item, i) => (
            <li key={i} className="py-2 text-slate-700 dark:text-slate-300">{item}</li>
          ))}
        </ul>
      </section>
    </div>
  )
}
