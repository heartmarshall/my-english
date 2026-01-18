import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client/react';
import { GET_DICTIONARY, DELETE_WORD } from '../graphql/queries';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Badge } from '../components/ui/badge';
import { useDebounce } from '../hooks/use-debounce';
import { Search, Trash2, BookOpen } from 'lucide-react';

export function Dictionary() {
  const [search, setSearch] = useState('');
  const [hasCard, setHasCard] = useState<boolean | undefined>(undefined);
  const debouncedSearch = useDebounce(search, 500);

  const { data, loading, error, refetch } = useQuery(GET_DICTIONARY, {
    variables: {
      filter: {
        search: debouncedSearch || undefined,
        hasCard: hasCard,
        limit: 50,
        sortBy: 'UPDATED_AT',
        sortDir: 'DESC',
      },
    },
  });

  const [deleteWord] = useMutation(DELETE_WORD, {
    onCompleted: () => {
      refetch();
    },
  });

  const handleDelete = async (id: string) => {
    if (confirm('Удалить это слово?')) {
      try {
        await deleteWord({ variables: { id } });
      } catch (err) {
        alert('Ошибка при удалении: ' + (err as Error).message);
      }
    }
  };

  if (loading) return <div className="p-6">Загрузка...</div>;
  if (error) return <div className="p-6 text-red-500">Ошибка: {error.message}</div>;

  const entries = data?.dictionary || [];

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Словарь</h1>
        <p className="text-muted-foreground mt-2">Поиск и управление словами</p>
      </div>

      <div className="flex gap-4 items-center">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
          <Input
            placeholder="Поиск слов..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>
        <div className="flex gap-2">
          <Button
            variant={hasCard === undefined ? 'default' : 'outline'}
            onClick={() => setHasCard(undefined)}
          >
            Все
          </Button>
          <Button
            variant={hasCard === true ? 'default' : 'outline'}
            onClick={() => setHasCard(true)}
          >
            С карточками
          </Button>
          <Button
            variant={hasCard === false ? 'default' : 'outline'}
            onClick={() => setHasCard(false)}
          >
            Без карточек
          </Button>
        </div>
      </div>

      <div className="space-y-4">
        {entries.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center text-muted-foreground">
              <BookOpen className="size-12 mx-auto mb-4 opacity-50" />
              <p>Слова не найдены</p>
            </CardContent>
          </Card>
        ) : (
          entries.map((entry: any) => (
            <Card key={entry.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-2xl">{entry.text}</CardTitle>
                    {entry.card && (
                      <CardDescription className="mt-1">
                        Статус: {entry.card.status} | Интервал: {entry.card.intervalDays} дн.
                      </CardDescription>
                    )}
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(entry.id)}
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                {entry.senses && entry.senses.length > 0 && (
                  <div className="space-y-2">
                    {entry.senses.map((sense) => (
                      <div key={sense.id} className="border-l-2 border-primary pl-4">
                        {sense.partOfSpeech && (
                          <Badge variant="outline" className="mb-2">
                            {sense.partOfSpeech}
                          </Badge>
                        )}
                        {sense.definition && (
                          <p className="text-sm text-muted-foreground mb-1">
                            {sense.definition}
                          </p>
                        )}
                        {sense.translations && sense.translations.length > 0 && (
                          <p className="font-medium">
                            {sense.translations.map((t) => t.text).join(', ')}
                          </p>
                        )}
                        {sense.examples && sense.examples.length > 0 && (
                          <div className="mt-2 space-y-1">
                            {sense.examples.map((ex) => (
                              <p key={ex.id} className="text-sm italic text-muted-foreground">
                                "{ex.sentence}"
                                {ex.translation && ` — ${ex.translation}`}
                              </p>
                            ))}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
                {entry.images && entry.images.length > 0 && (
                  <div className="flex gap-2 flex-wrap">
                    {entry.images.map((img) => (
                      <img
                        key={img.id}
                        src={img.url}
                        alt={img.caption || ''}
                        className="h-20 w-20 object-cover rounded"
                      />
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  );
}

