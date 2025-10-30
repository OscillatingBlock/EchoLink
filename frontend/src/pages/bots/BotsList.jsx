import { useNavigate } from 'react-router-dom'
import BotCard from '../../components/BotCard'

const MOCK_BOTS = [
  { id: 'sales-1', name: 'Sales Assistant', description: 'Outbound sales calls to warm leads', status: 'Active', calls: 245, successRate: 90 },
  { id: 'support-1', name: 'Support Bot', description: 'Level 1 support triage', status: 'Active', calls: 123, successRate: 94 },
  { id: 'survey-1', name: 'Survey Bot', description: 'NPS and post-call surveys', status: 'Paused', calls: 58, successRate: 88 },
]

export default function BotsList() {
  const navigate = useNavigate()
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Bots</h1>
        <a href="/create" className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-medium text-white transition-colors hover:bg-indigo-700">New Bot</a>
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {MOCK_BOTS.map((b) => (
          <BotCard key={b.id} bot={b} onClick={() => navigate(`/bots/${b.id}`)} />
        ))}
      </div>
    </div>
  )
}
