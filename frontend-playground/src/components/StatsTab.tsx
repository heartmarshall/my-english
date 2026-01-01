import { useQuery } from '@apollo/client/react'
import { GET_STATS } from '../graphql/queries'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Button } from './ui/button'

export default function StatsTab() {
  const { data, loading, error, refetch } = useQuery(GET_STATS)

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>

  const stats = (data as any)?.stats

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <Button onClick={() => refetch()}>Refresh</Button>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Total Words</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.totalWords || 0}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Mastered</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.masteredCount || 0}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Learning</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.learningCount || 0}</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Due for Review</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-3xl font-bold">{stats?.dueForReviewCount || 0}</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

