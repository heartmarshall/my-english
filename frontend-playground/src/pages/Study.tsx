import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client/react';
import { GET_STUDY_QUEUE, REVIEW_CARD } from '../graphql/queries';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Badge } from '../components/ui/badge';
import { RotateCcw, X, CheckCircle, Zap } from 'lucide-react';

type ReviewGrade = 'AGAIN' | 'HARD' | 'GOOD' | 'EASY';

export function Study() {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [showAnswer, setShowAnswer] = useState(false);
  const [startTime, setStartTime] = useState(Date.now());

  const { data, loading, error, refetch } = useQuery(GET_STUDY_QUEUE, {
    variables: { limit: 20 },
  });

  const [reviewCard] = useMutation(REVIEW_CARD, {
    onCompleted: () => {
      refetch();
      setShowAnswer(false);
      setStartTime(Date.now());
    },
  });

  if (loading) return <div className="p-6">Загрузка...</div>;
  if (error) return <div className="p-6 text-red-500">Ошибка: {error.message}</div>;

  const entries = data?.studyQueue || [];

  if (entries.length === 0) {
    return (
      <div className="p-6">
        <Card>
          <CardContent className="py-12 text-center">
            <CheckCircle className="size-12 mx-auto mb-4 text-green-500" />
            <h2 className="text-2xl font-bold mb-2">Отлично!</h2>
            <p className="text-muted-foreground">Нет карточек для повторения на данный момент.</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const currentEntry = entries[currentIndex];
  const card = currentEntry?.card;

  const handleGrade = async (grade: ReviewGrade) => {
    if (!card) return;

    const timeTakenMs = Date.now() - startTime;

    try {
      await reviewCard({
        variables: {
          cardId: card.id,
          grade,
          timeTakenMs,
        },
      });

      if (currentIndex < entries.length - 1) {
        setCurrentIndex(currentIndex + 1);
      } else {
        setCurrentIndex(0);
      }
    } catch (err) {
      alert('Ошибка при сохранении ответа: ' + (err as Error).message);
    }
  };

  if (!currentEntry) {
    return <div className="p-6">Карточка не найдена</div>;
  }

  return (
    <div className="p-6 max-w-2xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Изучение</h1>
          <p className="text-muted-foreground mt-2">
            Карточка {currentIndex + 1} из {entries.length}
          </p>
        </div>
        {card && (
          <Badge variant="outline">
            {card.status} | Интервал: {card.intervalDays} дн.
          </Badge>
        )}
      </div>

      <Card className="min-h-[400px]">
        <CardHeader>
          <CardTitle className="text-4xl text-center">{currentEntry.text}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          {!showAnswer ? (
            <div className="text-center py-12">
              <Button onClick={() => setShowAnswer(true)} size="lg">
                Показать ответ
              </Button>
            </div>
          ) : (
            <div className="space-y-6">
              {currentEntry.senses && currentEntry.senses.length > 0 && (
                <div className="space-y-4">
                  {currentEntry.senses.map((sense) => (
                    <div key={sense.id} className="border-l-2 border-primary pl-4 space-y-2">
                      {sense.partOfSpeech && (
                        <Badge variant="outline">{sense.partOfSpeech}</Badge>
                      )}
                      {sense.definition && (
                        <p className="text-muted-foreground">{sense.definition}</p>
                      )}
                      {sense.translations && sense.translations.length > 0 && (
                        <p className="text-xl font-semibold">
                          {sense.translations.map((t) => t.text).join(', ')}
                        </p>
                      )}
                      {sense.examples && sense.examples.length > 0 && (
                        <div className="space-y-2 mt-4">
                          {sense.examples.map((ex) => (
                            <div key={ex.id} className="text-sm">
                              <p className="italic">{ex.sentence}</p>
                              {ex.translation && (
                                <p className="text-muted-foreground mt-1">{ex.translation}</p>
                              )}
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}

              {currentEntry.images && currentEntry.images.length > 0 && (
                <div className="flex gap-2 flex-wrap justify-center">
                  {currentEntry.images.map((img) => (
                    <img
                      key={img.id}
                      src={img.url}
                      alt={img.caption || ''}
                      className="h-32 w-32 object-cover rounded"
                    />
                  ))}
                </div>
              )}

              {currentEntry.pronunciations && currentEntry.pronunciations.length > 0 && (
                <div className="space-y-2">
                  {currentEntry.pronunciations.map((pron) => (
                    <div key={pron.id} className="flex items-center gap-2">
                      {pron.audioUrl && (
                        <audio controls className="flex-1">
                          <source src={pron.audioUrl} />
                        </audio>
                      )}
                      {pron.transcription && (
                        <span className="text-muted-foreground">[{pron.transcription}]</span>
                      )}
                    </div>
                  ))}
                </div>
              )}

              <div className="grid grid-cols-2 md:grid-cols-4 gap-2 pt-4 border-t">
                <Button
                  variant="destructive"
                  onClick={() => handleGrade('AGAIN')}
                  className="flex-col h-auto py-4"
                >
                  <RotateCcw className="size-5 mb-1" />
                  <span className="text-xs">Снова</span>
                </Button>
                <Button
                  variant="outline"
                  onClick={() => handleGrade('HARD')}
                  className="flex-col h-auto py-4"
                >
                  <X className="size-5 mb-1" />
                  <span className="text-xs">Сложно</span>
                </Button>
                <Button
                  variant="default"
                  onClick={() => handleGrade('GOOD')}
                  className="flex-col h-auto py-4"
                >
                  <CheckCircle className="size-5 mb-1" />
                  <span className="text-xs">Хорошо</span>
                </Button>
                <Button
                  variant="outline"
                  onClick={() => handleGrade('EASY')}
                  className="flex-col h-auto py-4"
                >
                  <Zap className="size-5 mb-1" />
                  <span className="text-xs">Легко</span>
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

