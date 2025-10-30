import {NavLink} from 'react-router-dom'
import logo from "/echolink.png"
const navLinkClass=({isActive}) =>
  `flex items-center gap-3 rounded-md px-3 py-2 transition-colors ${isActive
    ? 'bg-indigo-500/10 text-indigo-600 dark:text-indigo-300'
    :'text-slate-600 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white'
  }`

export default function Sidebar()
{
  return (
    <aside className="h-full w-64 shrink-0 border-r border-slate-200 bg-white dark:border-[#1d1d1d] dark:bg-[#0c0c0c]">
      <div className="border-b border-slate-200 p-4 dark:border-[#1e1e1e]">
        <div className="h-10 font-semibold text-indigo-600  dark:text-indigo-400"><img className='h-full' src={logo} alt="" /></div>
      </div>
      <nav className="space-y-1 p-3">
        <NavLink className={navLinkClass} to="/dashboard">
          <span>Home</span>
        </NavLink>
        <NavLink className={navLinkClass} to="/bots">
          <span>Bots</span>
        </NavLink>
        <NavLink className={navLinkClass} to="/create">
          <span>Create Bot</span>
        </NavLink>
      </nav>
    </aside>
  )
}
