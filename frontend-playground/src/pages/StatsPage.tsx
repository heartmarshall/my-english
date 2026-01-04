import { gql } from "@apollo/client"
import { useQuery } from "@apollo/client/react"
import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { BookOpen, CheckCircle2, Clock, Target } from "lucide-react"

const STATS_QUERY = gql`
  query Stats {
    stats {
      totalCards
      masteredCount
      learningCount
      dueCount
    }
  }
`

interface StatsData {
  stats: {
    totalCards: number
    masteredCount: number
    learningCount: number
    dueCount: number
  }
}

export function StatsPage() {
  const { data, loading, error } = useQuery<StatsData>(STATS_QUERY)

  if (loading) {
    return (
      <div className="text-center py-12">
        <p className="text-muted-foreground">Загрузка статистики...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <p className="text-destructive">Ошибка загрузки статистики</p>
      </div>
    )
  }

  const stats = data?.stats
  if (!stats) return null

  const masteredPercent = stats.totalCards > 0 
    ? Math.round((stats.masteredCount / stats.totalCards) * 100) 
    : 0

  const learningPercent = stats.totalCards > 0
    ? Math.round((stats.learningCount / stats.totalCards) * 100)
    : 0

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Статистика</h2>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <Card className="p-6">
          <div className="flex items-center justify-between mb-2">
            <BookOpen className="size-8 text-primary" />
            <Badge variant="secondary" className="text-lg">
              {stats.totalCards}
            </Badge>
          </div>
          <h3 className="font-semibold">Всего карточек</h3>
          <p className="text-sm text-muted-foreground mt-1">
            Все слова в вашей коллекции
          </p>
        </Card>

        <Card className="p-6">
          <div className="flex items-center justify-between mb-2">
            <CheckCircle2 className="size-8 text-green-500" />
            <Badge variant="secondary" className="text-lg">
              {stats.masteredCount}
            </Badge>
          </div>
          <h3 className="font-semibold">Изучено</h3>
          <p className="text-sm text-muted-foreground mt-1">
            {masteredPercent}% от общего числа
          </p>
        </Card>

        <Card className="p-6">
          <div className="flex items-center justify-between mb-2">
            <Clock className="size-8 text-yellow-500" />
            <Badge variant="secondary" className="text-lg">
              {stats.learningCount}
            </Badge>
          </div>
          <h3 className="font-semibold">Изучаю</h3>
          <p className="text-sm text-muted-foreground mt-1">
            {learningPercent}% от общего числа
          </p>
        </Card>

        <Card className="p-6">
          <div className="flex items-center justify-between mb-2">
            <Target className="size-8 text-orange-500" />
            <Badge variant="secondary" className="text-lg">
              {stats.dueCount}
            </Badge>
          </div>
          <h3 className="font-semibold">К повторению</h3>
          <p className="text-sm text-muted-foreground mt-1">
            Готовы к изучению
          </p>
        </Card>
      </div>

      <Card className="p-6">
        <h3 className="font-semibold mb-4">Прогресс изучения</h3>
        <div className="space-y-4">
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span>Изучено</span>
              <span>{masteredPercent}%</span>
            </div>
            <div className="w-full bg-muted rounded-full h-3">
              <div
                className="bg-green-500 h-3 rounded-full transition-all"
                style={{ width: `${masteredPercent}%` }}
              />
            </div>
          </div>
          <div>
            <div className="flex justify-between text-sm mb-1">
              <span>Изучаю</span>
              <span>{learningPercent}%</span>
            </div>
            <div className="w-full bg-muted rounded-full h-3">
              <div
                className="bg-yellow-500 h-3 rounded-full transition-all"
                style={{ width: `${learningPercent}%` }}
              />
            </div>
          </div>
        </div>
      </Card>
    </div>
  )
}

