import Sidebar from '../components/Sidebar'
import Header from '../components/Header'
import { Outlet } from 'react-router-dom'

export default function DashboardLayout()
{
    return (
        <div className="h-full grid grid-cols-[16rem_1fr]">
            <Sidebar />
            <div className="flex min-w-0 flex-col">
                <Header />
                <main className="mx-auto w-full background max-w-7xl flex-1 p-4">
                    <Outlet />
                </main>
            </div>
        </div>
    )
}
