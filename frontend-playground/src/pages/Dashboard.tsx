import { useQuery } from '@apollo/client/react';
import { GET_DASHBOARD_STATS } from '../graphql/queries';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card';
import { BookOpen, BookCheck, Clock, TrendingUp, CheckCircle2, Circle } from 'lucide-react';

export function Dashboard() {
  const { data, loading, error } = useQuery(GET_DASHBOARD_STATS);

  if (loading) return <div className="p-6">Загрузка...</div>;
  if (error) return <div className="p-6 text-red-500">Ошибка: {error.message}</div>;

  const stats = data?.dashboardStats;

  const statCards = [
    {
      title: 'Всего слов',
      value: stats?.totalWords || 0,
      icon: BookOpen,
      color: 'text-blue-500',
    },
    {
      title: 'Карточек',
      value: stats?.totalCards || 0,
      icon: BookCheck,
      color: 'text-green-500',
    },
    {
      title: 'К изучению',
      value: stats?.newCards || 0,
      icon: Circle,
      color: 'text-gray-500',
    },
    {
      title: 'Изучаются',
      value: stats?.learningCards || 0,
      icon: TrendingUp,
      color: 'text-yellow-500',
    },
    {
      title: 'На повтор',
      value: stats?.reviewCards || 0,
      icon: Clock,
      color: 'text-orange-500',
    },
    {
      title: 'Изучено',
      value: stats?.masteredCards || 0,
      icon: CheckCircle2,
      color: 'text-purple-500',
    },
    {
      title: 'Сегодня',
      value: stats?.dueToday || 0,
      icon: Clock,
      color: 'text-red-500',
    },
  ];

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Панель управления</h1>
        <p className="text-muted-foreground mt-2">Общая статистика вашего обучения</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {statCards.map((stat) => {
          const Icon = stat.icon;
          return (
            <Card key={stat.title}>
              <CardHeader>
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.title}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-3">
                  <Icon className={`${stat.color} size-8`} />
                  <div className="text-3xl font-bold">{stat.value}</div>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </div>
  );
}

