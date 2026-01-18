import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client/react';
import { GET_INBOX_ITEMS, ADD_TO_INBOX, DELETE_INBOX_ITEM, CONVERT_INBOX_TO_WORD, FETCH_SUGGESTIONS } from '../graphql/queries';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Textarea } from '../components/ui/textarea';
import { Label } from '../components/ui/label';
import { Badge } from '../components/ui/badge';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '../components/ui/dialog';
import { Trash2, Plus, ArrowRight } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export function Inbox() {
  const navigate = useNavigate();
  const [newText, setNewText] = useState('');
  const [newContext, setNewContext] = useState('');
  const [convertDialogOpen, setConvertDialogOpen] = useState<string | null>(null);
  const [convertWord, setConvertWord] = useState('');
  const [convertCreateCard, setConvertCreateCard] = useState(true);

  const { data, loading, error, refetch } = useQuery(GET_INBOX_ITEMS);

  const [addToInbox] = useMutation(ADD_TO_INBOX, {
    onCompleted: () => {
      setNewText('');
      setNewContext('');
      refetch();
    },
  });

  const [deleteInboxItem] = useMutation(DELETE_INBOX_ITEM, {
    onCompleted: () => {
      refetch();
    },
  });

  const [convertInboxToWord] = useMutation(CONVERT_INBOX_TO_WORD, {
    onCompleted: () => {
      setConvertDialogOpen(null);
      setConvertWord('');
      refetch();
      navigate('/dictionary');
    },
  });

  const { data: suggestionsData, loading: suggestionsLoading } = useQuery(FETCH_SUGGESTIONS, {
    variables: {
      text: convertWord,
      sources: ['freedict'],
    },
    skip: !convertWord || convertWord.length < 2 || !convertDialogOpen,
  });

  const handleAdd = async () => {
    if (!newText.trim()) {
      alert('Введите текст');
      return;
    }

    try {
      await addToInbox({
        variables: {
          text: newText.trim(),
          context: newContext.trim() || undefined,
        },
      });
    } catch (err) {
      alert('Ошибка при добавлении: ' + (err as Error).message);
    }
  };

  const handleDelete = async (id: string) => {
    if (confirm('Удалить этот элемент?')) {
      try {
        await deleteInboxItem({ variables: { id } });
      } catch (err) {
        alert('Ошибка при удалении: ' + (err as Error).message);
      }
    }
  };

  const handleConvert = async (inboxId: string) => {
    if (!convertWord.trim()) {
      alert('Введите слово');
      return;
    }

    const suggestions = suggestionsData?.fetchSuggestions || [];
    if (suggestions.length === 0 || !suggestions[0]?.senses || suggestions[0].senses.length === 0) {
      alert('Не найдены подсказки для этого слова. Попробуйте добавить слово вручную.');
      return;
    }

    const firstSuggestion = suggestions[0];
    const firstSense = firstSuggestion.senses[0];

    try {
      await convertInboxToWord({
        variables: {
          inboxId,
          input: {
            text: convertWord.trim(),
            senses: [
              {
                definition: firstSense.definition || undefined,
                partOfSpeech: firstSense.partOfSpeech || undefined,
                sourceSlug: firstSuggestion.sourceSlug,
                translations: firstSense.translations.map((t) => ({
                  text: t,
                  sourceSlug: firstSuggestion.sourceSlug,
                })),
                examples: [
                  ...firstSense.examples.map((ex) => ({
                    sentence: ex.sentence,
                    translation: ex.translation || undefined,
                    sourceSlug: firstSuggestion.sourceSlug,
                  })),
                ],
              },
            ],
            images: [],
            pronunciations: [],
            createCard: convertCreateCard,
          },
        },
      });
    } catch (err) {
      alert('Ошибка при конвертации: ' + (err as Error).message);
    }
  };

  if (loading) return <div className="p-6">Загрузка...</div>;
  if (error) return <div className="p-6 text-red-500">Ошибка: {error.message}</div>;

  const items = data?.inboxItems || [];

  return (
    <div className="p-6 space-y-6 max-w-4xl mx-auto">
      <div>
        <h1 className="text-3xl font-bold">Входящие</h1>
        <p className="text-muted-foreground mt-2">Быстрое сохранение слов и контекста</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Добавить в Inbox</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label htmlFor="text">Текст</Label>
            <Input
              id="text"
              value={newText}
              onChange={(e) => setNewText(e.target.value)}
              placeholder="Введите слово или фразу..."
              className="mt-1"
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault();
                  handleAdd();
                }
              }}
            />
          </div>
          <div>
            <Label htmlFor="context">Контекст (опционально)</Label>
            <Textarea
              id="context"
              value={newContext}
              onChange={(e) => setNewContext(e.target.value)}
              placeholder="Где вы встретили это слово?"
              className="mt-1"
            />
          </div>
          <Button onClick={handleAdd} disabled={!newText.trim()}>
            <Plus className="size-4 mr-2" />
            Добавить
          </Button>
        </CardContent>
      </Card>

      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Элементы ({items.length})</h2>
        {items.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center text-muted-foreground">
              <p>Inbox пуст</p>
            </CardContent>
          </Card>
        ) : (
          items.map((item: any) => (
            <Card key={item.id}>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle>{item.text}</CardTitle>
                    {item.context && (
                      <CardDescription className="mt-2">{item.context}</CardDescription>
                    )}
                    <CardDescription className="mt-1">
                      {new Date(item.createdAt).toLocaleString('ru-RU')}
                    </CardDescription>
                  </div>
                  <div className="flex gap-2">
                    <Dialog
                      open={convertDialogOpen === item.id}
                      onOpenChange={(open) => {
                        setConvertDialogOpen(open ? item.id : null);
                        if (open) {
                          setConvertWord(item.text);
                        }
                      }}
                    >
                      <DialogTrigger asChild>
                        <Button variant="outline" size="sm">
                          <ArrowRight className="size-4 mr-2" />
                          Конвертировать
                        </Button>
                      </DialogTrigger>
                      <DialogContent>
                        <DialogHeader>
                          <DialogTitle>Конвертировать в слово</DialogTitle>
                        </DialogHeader>
                        <div className="space-y-4">
                          <div>
                            <Label htmlFor="convertWord">Слово</Label>
                            <Input
                              id="convertWord"
                              value={convertWord}
                              onChange={(e) => setConvertWord(e.target.value)}
                              className="mt-1"
                            />
                          </div>
                          {suggestionsLoading && (
                            <p className="text-sm text-muted-foreground">Загрузка подсказок...</p>
                          )}
                          {suggestionsData?.fetchSuggestions && suggestionsData.fetchSuggestions.length > 0 && (
                            <div className="space-y-2">
                              <p className="text-sm font-medium">Найдены подсказки:</p>
                              {suggestionsData.fetchSuggestions[0].senses.map((sense, idx) => (
                                <div key={idx} className="border rounded p-2 text-sm">
                                  {sense.definition && <p className="text-muted-foreground">{sense.definition}</p>}
                                  {sense.translations.length > 0 && (
                                    <p className="font-medium">{sense.translations.join(', ')}</p>
                                  )}
                                </div>
                              ))}
                            </div>
                          )}
                          <div className="flex items-center gap-2">
                            <input
                              type="checkbox"
                              id="convertCreateCard"
                              checked={convertCreateCard}
                              onChange={(e) => setConvertCreateCard(e.target.checked)}
                            />
                            <Label htmlFor="convertCreateCard" className="cursor-pointer">
                              Создать карточку для изучения
                            </Label>
                          </div>
                          <Button
                            onClick={() => handleConvert(item.id)}
                            disabled={!convertWord.trim() || suggestionsLoading}
                          >
                            Конвертировать
                          </Button>
                        </div>
                      </DialogContent>
                    </Dialog>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDelete(item.id)}
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>
            </Card>
          ))
        )}
      </div>
    </div>
  );
}

