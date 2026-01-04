import { useState } from "react"
import { gql } from "@apollo/client"
import { useQuery, useMutation } from "@apollo/client/react"
import { Input } from "@/components/ui/input"
import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useDebounce } from "@/hooks/use-debounce"
import { Search, BookOpen, Plus, BookPlus } from "lucide-react"
import { toast } from "sonner"

const SEARCH_QUERY = gql`
  query Search($query: String!) {
    search(query: $query) {
      id
      text
      senses {
        id
        partOfSpeech
        definition
        translations {
          text
        }
        examples {
          id
          sentenceEn
          sentenceRu
        }
      }
      pronunciations {
        audioUrl
        transcription
        region
      }
    }
  }
`

const CREATE_CARD_FROM_SENSE_MUTATION = gql`
  mutation CreateCardFromSense($input: CreateCardInput!) {
    createCard(input: $input) {
      id
      customText
      progress {
        status
      }
    }
  }
`

const CREATE_LEXEME_MUTATION = gql`
  mutation CreateLexeme($input: CreateLexemeInput!) {
    createLexeme(input: $input) {
      id
      text
      senses {
        id
        partOfSpeech
        definition
        translations {
          text
        }
      }
    }
  }
`

interface Lexeme {
  id: string
  text: string
  senses: Array<{
    id: string
    partOfSpeech: string
    definition: string
    translations: Array<{ text: string }>
    examples: Array<{
      id: string
      sentenceEn: string
      sentenceRu?: string | null
    }>
  }>
  pronunciations: Array<{
    audioUrl: string
    transcription?: string | null
    region: string
  }>
}

export function DictionaryPage() {
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedSense, setSelectedSense] = useState<{ senseId: string; translations: string[] } | null>(null)
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [isAddWordDialogOpen, setIsAddWordDialogOpen] = useState(false)
  const [cardForm, setCardForm] = useState({
    note: "",
    tags: "",
  })
  const [newWordForm, setNewWordForm] = useState({
    text: "",
    partOfSpeech: "NOUN",
    definition: "",
    translations: "",
    transcription: "",
    exampleEn: "",
    exampleRu: "",
  })
  const debouncedSearch = useDebounce(searchQuery, 500)

  const { data, loading, error } = useQuery<{ search: Lexeme[] }>(SEARCH_QUERY, {
    variables: { query: debouncedSearch },
    skip: !debouncedSearch || debouncedSearch.length < 2,
  })

  const [createCard, { loading: creating }] = useMutation(CREATE_CARD_FROM_SENSE_MUTATION, {
    onCompleted: () => {
      toast.success("Слово добавлено в мои слова!")
      setIsDialogOpen(false)
      setCardForm({ note: "", tags: "" })
      setSelectedSense(null)
    },
    onError: (error) => {
      toast.error("Ошибка: " + error.message)
    },
  })

  const [createLexeme, { loading: creatingLexeme }] = useMutation(CREATE_LEXEME_MUTATION, {
    onCompleted: () => {
      toast.success("Слово добавлено в словарь!")
      setIsAddWordDialogOpen(false)
      setNewWordForm({
        text: "",
        partOfSpeech: "NOUN",
        definition: "",
        translations: "",
        transcription: "",
        exampleEn: "",
        exampleRu: "",
      })
      // Обновляем результаты поиска, если слово совпадает с текущим запросом
      if (debouncedSearch && debouncedSearch.toLowerCase() === newWordForm.text.toLowerCase()) {
        // Можно добавить refetch здесь, если нужно
      }
    },
    onError: (error) => {
      toast.error("Ошибка: " + error.message)
    },
  })

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
  }

  const handleAddToCards = (senseId: string, translations: string[]) => {
    setSelectedSense({ senseId, translations })
    setIsDialogOpen(true)
  }

  const handleCreateCard = () => {
    if (!selectedSense) return

    const tags = cardForm.tags
      .split(",")
      .map((t) => t.trim())
      .filter((t) => t.length > 0)

    createCard({
      variables: {
        input: {
          senseId: selectedSense.senseId,
          translations: selectedSense.translations,
          note: cardForm.note.trim() || undefined,
          tags: tags.length > 0 ? tags : undefined,
        },
      },
    })
  }

  const handleCreateLexeme = () => {
    if (!newWordForm.text.trim()) {
      toast.error("Введите слово")
      return
    }
    if (!newWordForm.definition.trim()) {
      toast.error("Введите определение")
      return
    }

    const translations = newWordForm.translations
      .split(",")
      .map((t) => t.trim())
      .filter((t) => t.length > 0)

    if (translations.length === 0) {
      toast.error("Введите хотя бы один перевод")
      return
    }

    const examples = []
    if (newWordForm.exampleEn.trim()) {
      examples.push({
        sentenceEn: newWordForm.exampleEn.trim(),
        sentenceRu: newWordForm.exampleRu.trim() || undefined,
      })
    }

    createLexeme({
      variables: {
        input: {
          text: newWordForm.text.trim(),
          partOfSpeech: newWordForm.partOfSpeech,
          definition: newWordForm.definition.trim(),
          translations,
          transcription: newWordForm.transcription.trim() || undefined,
          examples: examples.length > 0 ? examples : undefined,
        },
      },
    })
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold mb-2">Словарь</h2>
          <p className="text-muted-foreground">
            Поиск по глобальному словарю
          </p>
        </div>
        <Dialog open={isAddWordDialogOpen} onOpenChange={setIsAddWordDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <BookPlus className="size-4" />
              Добавить слово
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Добавить слово в словарь</DialogTitle>
              <DialogDescription>
                Создайте новую запись в глобальном словаре
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="word-text">Слово *</Label>
                <Input
                  id="word-text"
                  placeholder="Например: hello"
                  value={newWordForm.text}
                  onChange={(e) =>
                    setNewWordForm({ ...newWordForm, text: e.target.value })
                  }
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="word-pos">Часть речи *</Label>
                <Select
                  value={newWordForm.partOfSpeech}
                  onValueChange={(value) =>
                    setNewWordForm({ ...newWordForm, partOfSpeech: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="NOUN">Существительное</SelectItem>
                    <SelectItem value="VERB">Глагол</SelectItem>
                    <SelectItem value="ADJECTIVE">Прилагательное</SelectItem>
                    <SelectItem value="ADVERB">Наречие</SelectItem>
                    <SelectItem value="PRONOUN">Местоимение</SelectItem>
                    <SelectItem value="PREPOSITION">Предлог</SelectItem>
                    <SelectItem value="CONJUNCTION">Союз</SelectItem>
                    <SelectItem value="INTERJECTION">Междометие</SelectItem>
                    <SelectItem value="PHRASE">Фраза</SelectItem>
                    <SelectItem value="IDIOM">Идиома</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="word-definition">Определение *</Label>
                <Textarea
                  id="word-definition"
                  placeholder="Определение слова..."
                  value={newWordForm.definition}
                  onChange={(e) =>
                    setNewWordForm({ ...newWordForm, definition: e.target.value })
                  }
                  rows={3}
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="word-translations">Переводы *</Label>
                <Input
                  id="word-translations"
                  placeholder="Через запятую: привет, здравствуй"
                  value={newWordForm.translations}
                  onChange={(e) =>
                    setNewWordForm({ ...newWordForm, translations: e.target.value })
                  }
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="word-transcription">Транскрипция</Label>
                <Input
                  id="word-transcription"
                  placeholder="/həˈloʊ/"
                  value={newWordForm.transcription}
                  onChange={(e) =>
                    setNewWordForm({ ...newWordForm, transcription: e.target.value })
                  }
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="word-example-en">Пример (английский)</Label>
                <Input
                  id="word-example-en"
                  placeholder="Hello, how are you?"
                  value={newWordForm.exampleEn}
                  onChange={(e) =>
                    setNewWordForm({ ...newWordForm, exampleEn: e.target.value })
                  }
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="word-example-ru">Пример (русский)</Label>
                <Input
                  id="word-example-ru"
                  placeholder="Привет, как дела?"
                  value={newWordForm.exampleRu}
                  onChange={(e) =>
                    setNewWordForm({ ...newWordForm, exampleRu: e.target.value })
                  }
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsAddWordDialogOpen(false)}
                disabled={creatingLexeme}
              >
                Отмена
              </Button>
              <Button onClick={handleCreateLexeme} disabled={creatingLexeme}>
                {creatingLexeme ? "Добавление..." : "Добавить"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <form onSubmit={handleSearch} className="mb-6">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground size-4" />
          <Input
            placeholder="Введите слово для поиска..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
      </form>

      {loading && (
        <div className="text-center py-12">
          <p className="text-muted-foreground">Поиск...</p>
        </div>
      )}

      {error && (
        <div className="text-center py-12">
          <p className="text-destructive">Ошибка поиска: {error.message}</p>
        </div>
      )}

      {!debouncedSearch || debouncedSearch.length < 2 ? (
        <div className="text-center py-12">
          <BookOpen className="size-12 text-muted-foreground mx-auto mb-4" />
          <p className="text-muted-foreground">
            Введите минимум 2 символа для поиска
          </p>
        </div>
      ) : data?.search && data.search.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground">Ничего не найдено</p>
        </div>
      ) : (
        <div className="space-y-4">
          {data?.search.map((lexeme) => (
            <Card key={lexeme.id} className="p-6">
              <div className="mb-4">
                <h3 className="text-2xl font-bold mb-2">{lexeme.text}</h3>
                {lexeme.pronunciations.length > 0 && (
                  <div className="flex flex-wrap gap-2 mb-2">
                    {lexeme.pronunciations.map((pron, idx) => (
                      <div key={idx} className="flex items-center gap-2">
                        {pron.transcription && (
                          <span className="text-sm text-muted-foreground">
                            [{pron.transcription}]
                          </span>
                        )}
                        <Badge variant="outline" className="text-xs">
                          {pron.region}
                        </Badge>
                        {pron.audioUrl && (
                          <audio controls className="h-6">
                            <source src={pron.audioUrl} />
                          </audio>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="space-y-4">
                {lexeme.senses.map((sense) => (
                  <div key={sense.id} className="border-l-2 border-primary pl-4">
                    <div className="flex items-center justify-between mb-2">
                      <Badge variant="secondary">{sense.partOfSpeech}</Badge>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() =>
                          handleAddToCards(
                            sense.id,
                            sense.translations.map((t) => t.text)
                          )
                        }
                      >
                        <Plus className="size-4" />
                        Добавить в мои слова
                      </Button>
                    </div>
                    <p className="font-medium mb-2">{sense.definition}</p>
                    {sense.translations.length > 0 && (
                      <div className="mb-2">
                        <p className="text-sm text-muted-foreground mb-1">Переводы:</p>
                        <div className="flex flex-wrap gap-2">
                          {sense.translations.map((t, idx) => (
                            <Badge key={idx} variant="outline">
                              {t.text}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                    {sense.examples.length > 0 && (
                      <div>
                        <p className="text-sm text-muted-foreground mb-1">Примеры:</p>
                        <div className="space-y-1">
                          {sense.examples.map((ex) => (
                            <div key={ex.id} className="text-sm italic">
                              "{ex.sentenceEn}"
                              {ex.sentenceRu && (
                                <span className="text-muted-foreground block">
                                  "{ex.sentenceRu}"
                                </span>
                              )}
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </Card>
          ))}
        </div>
      )}

      {/* Диалог добавления слова */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Добавить слово в мои слова</DialogTitle>
            <DialogDescription>
              Создайте карточку из словаря
            </DialogDescription>
          </DialogHeader>
          {selectedSense && (
            <>
              <div className="grid gap-4 py-4">
                <div className="grid gap-2">
                  <Label>Переводы</Label>
                  <div className="flex flex-wrap gap-2 p-2 border rounded-md">
                    {selectedSense.translations.map((t, idx) => (
                      <Badge key={idx} variant="secondary">
                        {t}
                      </Badge>
                    ))}
                  </div>
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="note">Заметка</Label>
                  <Textarea
                    id="note"
                    placeholder="Дополнительная информация..."
                    value={cardForm.note}
                    onChange={(e) =>
                      setCardForm({ ...cardForm, note: e.target.value })
                    }
                    rows={3}
                  />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor="tags">Теги</Label>
                  <Input
                    id="tags"
                    placeholder="Через запятую"
                    value={cardForm.tags}
                    onChange={(e) =>
                      setCardForm({ ...cardForm, tags: e.target.value })
                    }
                  />
                </div>
              </div>
              <DialogFooter>
                <Button
                  variant="outline"
                  onClick={() => setIsDialogOpen(false)}
                  disabled={creating}
                >
                  Отмена
                </Button>
                <Button onClick={handleCreateCard} disabled={creating}>
                  {creating ? "Добавление..." : "Добавить"}
                </Button>
              </DialogFooter>
            </>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

