import { useState } from 'react';
import { useQuery, useMutation } from '@apollo/client/react';
import { FETCH_SUGGESTIONS, CREATE_WORD } from '../graphql/queries';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Textarea } from '../components/ui/textarea';
import { Label } from '../components/ui/label';
import { Badge } from '../components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../components/ui/tabs';
import { useDebounce } from '../hooks/use-debounce';
import { Loader2, Plus, X, Check, Trash2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

type Sense = {
  definition: string;
  partOfSpeech: string;
  translations: string[];
  examples: Array<{ sentence: string; translation: string }>;
};

type Image = {
  url: string;
  caption: string;
};

type Pronunciation = {
  audioUrl: string;
  transcription: string;
  region: string;
};

const PARTS_OF_SPEECH = [
  'NOUN',
  'VERB',
  'ADJECTIVE',
  'ADVERB',
  'PRONOUN',
  'PREPOSITION',
  'CONJUNCTION',
  'INTERJECTION',
  'PHRASE',
  'IDIOM',
  'OTHER',
] as const;

export function AddWord() {
  const navigate = useNavigate();
  const [word, setWord] = useState('');
  const [createCard, setCreateCard] = useState(true);
  const [senses, setSenses] = useState<Sense[]>([
    { definition: '', partOfSpeech: '', translations: [''], examples: [{ sentence: '', translation: '' }] },
  ]);
  const [images, setImages] = useState<Image[]>([{ url: '', caption: '' }]);
  const [pronunciations, setPronunciations] = useState<Pronunciation[]>([
    { audioUrl: '', transcription: '', region: '' },
  ]);

  const debouncedWord = useDebounce(word, 500);

  const { data: suggestionsData, loading: suggestionsLoading } = useQuery(FETCH_SUGGESTIONS, {
    variables: {
      text: debouncedWord,
      sources: ['freedict'],
    },
    skip: !debouncedWord || debouncedWord.length < 2,
  });

  const [createWord, { loading: creating }] = useMutation(CREATE_WORD, {
    onCompleted: () => {
      navigate('/dictionary');
    },
  });

  const suggestions = ((suggestionsData as any)?.fetchSuggestions as any[]) || [];
  const [selectedSenses, setSelectedSenses] = useState<Set<string>>(new Set());

  const handleToggleSense = (suggestionIndex: number, senseIndex: number, sourceSlug: string) => {
    const key = `${sourceSlug}-${suggestionIndex}-${senseIndex}`;
    const newSet = new Set(selectedSenses);
    if (newSet.has(key)) {
      newSet.delete(key);
    } else {
      newSet.add(key);
    }
    setSelectedSenses(newSet);
  };

  const addSense = () => {
    setSenses([...senses, { definition: '', partOfSpeech: '', translations: [''], examples: [{ sentence: '', translation: '' }] }]);
  };

  const removeSense = (index: number) => {
    setSenses(senses.filter((_, i) => i !== index));
  };

  const updateSense = (index: number, field: keyof Sense, value: any) => {
    const newSenses = [...senses];
    newSenses[index] = { ...newSenses[index], [field]: value };
    setSenses(newSenses);
  };

  const addTranslation = (senseIndex: number) => {
    const newSenses = [...senses];
    newSenses[senseIndex].translations.push('');
    setSenses(newSenses);
  };

  const removeTranslation = (senseIndex: number, transIndex: number) => {
    const newSenses = [...senses];
    newSenses[senseIndex].translations = newSenses[senseIndex].translations.filter((_, i) => i !== transIndex);
    setSenses(newSenses);
  };

  const updateTranslation = (senseIndex: number, transIndex: number, value: string) => {
    const newSenses = [...senses];
    newSenses[senseIndex].translations[transIndex] = value;
    setSenses(newSenses);
  };

  const addExample = (senseIndex: number) => {
    const newSenses = [...senses];
    newSenses[senseIndex].examples.push({ sentence: '', translation: '' });
    setSenses(newSenses);
  };

  const removeExample = (senseIndex: number, exIndex: number) => {
    const newSenses = [...senses];
    newSenses[senseIndex].examples = newSenses[senseIndex].examples.filter((_, i) => i !== exIndex);
    setSenses(newSenses);
  };

  const updateExample = (senseIndex: number, exIndex: number, field: 'sentence' | 'translation', value: string) => {
    const newSenses = [...senses];
    newSenses[senseIndex].examples[exIndex] = { ...newSenses[senseIndex].examples[exIndex], [field]: value };
    setSenses(newSenses);
  };

  const addImage = () => {
    setImages([...images, { url: '', caption: '' }]);
  };

  const removeImage = (index: number) => {
    setImages(images.filter((_, i) => i !== index));
  };

  const updateImage = (index: number, field: keyof Image, value: string) => {
    const newImages = [...images];
    newImages[index] = { ...newImages[index], [field]: value };
    setImages(newImages);
  };

  const addPronunciation = () => {
    setPronunciations([...pronunciations, { audioUrl: '', transcription: '', region: '' }]);
  };

  const removePronunciation = (index: number) => {
    setPronunciations(pronunciations.filter((_, i) => i !== index));
  };

  const updatePronunciation = (index: number, field: keyof Pronunciation, value: string) => {
    const newPronunciations = [...pronunciations];
    newPronunciations[index] = { ...newPronunciations[index], [field]: value };
    setPronunciations(newPronunciations);
  };

  const handleCreate = async () => {
    if (!word.trim()) {
      alert('Введите слово');
      return;
    }

    // Собираем senses из suggestions
    const suggestionSenses = suggestions.flatMap((suggestion: any, sIdx: number) =>
      suggestion.senses
        .map((sense: any, senseIdx: number) => {
          const key = `${suggestion.sourceSlug}-${sIdx}-${senseIdx}`;
          if (!selectedSenses.has(key)) return null;

          return {
            definition: sense.definition || undefined,
            partOfSpeech: sense.partOfSpeech || undefined,
            sourceSlug: suggestion.sourceSlug,
            translations: sense.translations.map((t: string) => ({
              text: t,
              sourceSlug: suggestion.sourceSlug,
            })),
            examples: sense.examples.map((ex: any) => ({
              sentence: ex.sentence,
              translation: ex.translation || undefined,
              sourceSlug: suggestion.sourceSlug,
            })),
          };
        })
        .filter(Boolean)
    );

    // Собираем senses из ручного ввода
    const manualSenses = senses
      .map((sense) => {
        const translations = sense.translations.filter((t) => t.trim());
        if (translations.length === 0) return null;

        return {
          definition: sense.definition.trim() || undefined,
          partOfSpeech: sense.partOfSpeech || undefined,
          sourceSlug: 'user',
          translations: translations.map((t) => ({
            text: t.trim(),
            sourceSlug: 'user',
          })),
          examples: sense.examples
            .filter((ex) => ex.sentence.trim())
            .map((ex) => ({
              sentence: ex.sentence.trim(),
              translation: ex.translation.trim() || undefined,
              sourceSlug: 'user',
            })),
        };
      })
      .filter(Boolean);

    const allSenses = [...suggestionSenses, ...manualSenses];

    if (allSenses.length === 0) {
      alert('Добавьте хотя бы одно значение слова (перевод)');
      return;
    }

    // Собираем images
    const validImages = images
      .filter((img) => img.url.trim())
      .map((img) => ({
        url: img.url.trim(),
        caption: img.caption.trim() || undefined,
        sourceSlug: 'user',
      }));

    // Собираем pronunciations
    const validPronunciations = pronunciations
      .filter((pron) => pron.audioUrl.trim())
      .map((pron) => ({
        audioUrl: pron.audioUrl.trim(),
        transcription: pron.transcription.trim() || undefined,
        region: pron.region.trim() || undefined,
        sourceSlug: 'user',
      }));

    try {
      await createWord({
        variables: {
          input: {
            text: word.trim(),
            senses: allSenses,
            images: validImages,
            pronunciations: validPronunciations,
            createCard,
          },
        },
      });
    } catch (err) {
      alert('Ошибка при создании: ' + (err as Error).message);
    }
  };

  return (
    <div className="p-6 space-y-6 max-w-5xl mx-auto">
      <div>
        <h1 className="text-3xl font-bold">Добавить слово</h1>
        <p className="text-muted-foreground mt-2">Создайте новое слово с полной информацией</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Основная информация</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label htmlFor="word">Слово *</Label>
            <Input
              id="word"
              value={word}
              onChange={(e) => setWord(e.target.value)}
              placeholder="Введите слово на английском..."
              className="mt-1"
            />
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="createCard"
              checked={createCard}
              onChange={(e) => setCreateCard(e.target.checked)}
            />
            <Label htmlFor="createCard" className="cursor-pointer">
              Создать карточку для изучения
            </Label>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="senses" className="w-full">
        <TabsList>
          <TabsTrigger value="senses">Значения слова</TabsTrigger>
          <TabsTrigger value="images">Изображения</TabsTrigger>
          <TabsTrigger value="pronunciations">Произношение</TabsTrigger>
          <TabsTrigger value="suggestions">Подсказки</TabsTrigger>
        </TabsList>

        <TabsContent value="senses" className="space-y-4">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Значения слова *</CardTitle>
                <Button variant="outline" size="sm" onClick={addSense}>
                  <Plus className="size-4 mr-2" />
                  Добавить значение
                </Button>
              </div>
              <CardDescription>Добавьте переводы и определения для слова</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {senses.map((sense, senseIdx) => (
                <Card key={senseIdx} className="border-2">
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-lg">Значение {senseIdx + 1}</CardTitle>
                      {senses.length > 1 && (
                        <Button variant="ghost" size="icon" onClick={() => removeSense(senseIdx)}>
                          <Trash2 className="size-4" />
                        </Button>
                      )}
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <Label>Часть речи</Label>
                      <Select
                        value={sense.partOfSpeech}
                        onValueChange={(value) => updateSense(senseIdx, 'partOfSpeech', value)}
                      >
                        <SelectTrigger className="mt-1 w-full">
                          <SelectValue placeholder="Выберите часть речи" />
                        </SelectTrigger>
                        <SelectContent>
                          {PARTS_OF_SPEECH.map((pos) => (
                            <SelectItem key={pos} value={pos}>
                              {pos}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>

                    <div>
                      <Label>Определение (опционально)</Label>
                      <Textarea
                        value={sense.definition}
                        onChange={(e) => updateSense(senseIdx, 'definition', e.target.value)}
                        placeholder="Определение слова..."
                        className="mt-1"
                        rows={3}
                      />
                    </div>

                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <Label>Переводы *</Label>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => addTranslation(senseIdx)}
                        >
                          <Plus className="size-4 mr-1" />
                          Добавить
                        </Button>
                      </div>
                      <div className="space-y-2">
                        {sense.translations.map((trans, transIdx) => (
                          <div key={transIdx} className="flex gap-2">
                            <Input
                              value={trans}
                              onChange={(e) => updateTranslation(senseIdx, transIdx, e.target.value)}
                              placeholder="Перевод..."
                              className="flex-1"
                            />
                            {sense.translations.length > 1 && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => removeTranslation(senseIdx, transIdx)}
                              >
                                <X className="size-4" />
                              </Button>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>

                    <div>
                      <div className="flex items-center justify-between mb-2">
                        <Label>Примеры использования (опционально)</Label>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => addExample(senseIdx)}
                        >
                          <Plus className="size-4 mr-1" />
                          Добавить
                        </Button>
                      </div>
                      <div className="space-y-2">
                        {sense.examples.map((ex, exIdx) => (
                          <div key={exIdx} className="space-y-2 border rounded p-3">
                            <div className="flex items-center justify-between">
                              <Label className="text-sm">Пример {exIdx + 1}</Label>
                              {sense.examples.length > 1 && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => removeExample(senseIdx, exIdx)}
                                >
                                  <X className="size-4" />
                                </Button>
                              )}
                            </div>
                            <Input
                              value={ex.sentence}
                              onChange={(e) => updateExample(senseIdx, exIdx, 'sentence', e.target.value)}
                              placeholder="Предложение на английском..."
                            />
                            <Input
                              value={ex.translation}
                              onChange={(e) => updateExample(senseIdx, exIdx, 'translation', e.target.value)}
                              placeholder="Перевод предложения (опционально)..."
                            />
                          </div>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="images" className="space-y-4">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Изображения</CardTitle>
                <Button variant="outline" size="sm" onClick={addImage}>
                  <Plus className="size-4 mr-2" />
                  Добавить изображение
                </Button>
              </div>
              <CardDescription>Добавьте изображения для визуализации слова</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {images.map((img, idx) => (
                <div key={idx} className="flex gap-2 items-start">
                  <div className="flex-1 space-y-2">
                    <Input
                      value={img.url}
                      onChange={(e) => updateImage(idx, 'url', e.target.value)}
                      placeholder="URL изображения..."
                    />
                    <Input
                      value={img.caption}
                      onChange={(e) => updateImage(idx, 'caption', e.target.value)}
                      placeholder="Подпись (опционально)..."
                    />
                  </div>
                  {images.length > 1 && (
                    <Button variant="ghost" size="icon" onClick={() => removeImage(idx)}>
                      <Trash2 className="size-4" />
                    </Button>
                  )}
                </div>
              ))}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="pronunciations" className="space-y-4">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Произношение</CardTitle>
                <Button variant="outline" size="sm" onClick={addPronunciation}>
                  <Plus className="size-4 mr-2" />
                  Добавить произношение
                </Button>
              </div>
              <CardDescription>Добавьте аудио и транскрипцию произношения</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {pronunciations.map((pron, idx) => (
                <div key={idx} className="flex gap-2 items-start">
                  <div className="flex-1 space-y-2">
                    <Input
                      value={pron.audioUrl}
                      onChange={(e) => updatePronunciation(idx, 'audioUrl', e.target.value)}
                      placeholder="URL аудио файла..."
                    />
                    <div className="grid grid-cols-2 gap-2">
                      <Input
                        value={pron.transcription}
                        onChange={(e) => updatePronunciation(idx, 'transcription', e.target.value)}
                        placeholder="Транскрипция (например, /həˈloʊ/)..."
                      />
                      <Input
                        value={pron.region}
                        onChange={(e) => updatePronunciation(idx, 'region', e.target.value)}
                        placeholder="Регион (US, UK, etc.)..."
                      />
                    </div>
                  </div>
                  {pronunciations.length > 1 && (
                    <Button variant="ghost" size="icon" onClick={() => removePronunciation(idx)}>
                      <Trash2 className="size-4" />
                    </Button>
                  )}
                </div>
              ))}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="suggestions" className="space-y-4">
          {suggestionsLoading && debouncedWord && (
            <Card>
              <CardContent className="py-12 text-center">
                <Loader2 className="size-8 mx-auto animate-spin text-muted-foreground" />
                <p className="text-muted-foreground mt-4">Загрузка подсказок...</p>
              </CardContent>
            </Card>
          )}

          {suggestions.length > 0 && (
            <div className="space-y-4">
              <h2 className="text-xl font-semibold">Подсказки из внешних источников</h2>
              <p className="text-sm text-muted-foreground">
                Выберите подсказки, которые хотите использовать. Они будут добавлены к вашим значениям слова.
              </p>
              {suggestions.map((suggestion: any, sIdx: number) => (
                <Card key={suggestion.sourceSlug}>
                  <CardHeader>
                    <CardTitle>{suggestion.sourceName}</CardTitle>
                    <CardDescription>Источник: {suggestion.sourceSlug}</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {suggestion.senses.map((sense: any, senseIdx: number) => {
                      const key = `${suggestion.sourceSlug}-${sIdx}-${senseIdx}`;
                      const isSelected = selectedSenses.has(key);

                      return (
                        <div
                          key={senseIdx}
                          className={`border rounded-lg p-4 cursor-pointer transition-colors ${
                            isSelected ? 'border-primary bg-primary/5' : 'border-border'
                          }`}
                          onClick={() => handleToggleSense(sIdx, senseIdx, suggestion.sourceSlug)}
                        >
                          <div className="flex items-start justify-between gap-4">
                            <div className="flex-1 space-y-2">
                              {sense.partOfSpeech && (
                                <Badge variant="outline">{sense.partOfSpeech}</Badge>
                              )}
                              {sense.definition && (
                                <p className="text-sm text-muted-foreground">{sense.definition}</p>
                              )}
                              {sense.translations.length > 0 && (
                                <p className="font-medium">
                                  {sense.translations.join(', ')}
                                </p>
                              )}
                              {sense.examples.length > 0 && (
                                <div className="space-y-1 mt-2">
                              {sense.examples.map((ex: any, exIdx: number) => (
                                <p key={exIdx} className="text-sm italic text-muted-foreground">
                                      "{ex.sentence}"
                                      {ex.translation && ` — ${ex.translation}`}
                                    </p>
                                  ))}
                                </div>
                              )}
                            </div>
                            <div className="flex items-center gap-2">
                              {isSelected ? (
                                <Check className="size-5 text-primary" />
                              ) : (
                                <Plus className="size-5 text-muted-foreground" />
                              )}
                            </div>
                          </div>
                        </div>
                      );
                    })}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          {!suggestionsLoading && debouncedWord && suggestions.length === 0 && (
            <Card>
              <CardContent className="py-12 text-center text-muted-foreground">
                <p>Подсказки не найдены для слова "{debouncedWord}"</p>
              </CardContent>
            </Card>
          )}

          {!debouncedWord && (
            <Card>
              <CardContent className="py-12 text-center text-muted-foreground">
                <p>Введите слово в поле выше, чтобы получить подсказки</p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>

      <div className="flex gap-4">
        <Button onClick={handleCreate} disabled={creating || !word.trim()}>
          {creating ? (
            <>
              <Loader2 className="size-4 animate-spin mr-2" />
              Создание...
            </>
          ) : (
            'Создать слово'
          )}
        </Button>
        <Button variant="outline" onClick={() => navigate('/dictionary')}>
          Отмена
        </Button>
      </div>
    </div>
  );
}
