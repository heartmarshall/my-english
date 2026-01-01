import { useState, useEffect } from "react"
import { useQuery, useMutation } from "@apollo/client/react" // Важно: импорт из @apollo/client
import { 
  Trash2Icon, 
  Volume2Icon, 
  XCircleIcon, 
  PlusIcon,
  Loader2Icon
} from "lucide-react"
import { toast } from "sonner"

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetFooter,
} from "@/components/ui/sheet"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

import { GET_WORD, UPDATE_WORD, DELETE_WORD, GET_WORDS } from "@/graphql/queries"
import { PartOfSpeech } from "@/gql/graphql"

interface WordCardSheetProps {
  wordId: string | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function WordCardSheet({ wordId, open, onOpenChange }: WordCardSheetProps) {
  const [isEditing, setIsEditing] = useState(false)
  
  // -- Queries & Mutations --
  const { data, loading, error, refetch } = useQuery(GET_WORD, {
    variables: { id: wordId || "" },
    skip: !wordId,
    fetchPolicy: "network-only"
  }) as { data: any; loading: boolean; error: any; refetch: () => void }

  const [updateWord, { loading: updating }] = useMutation(UPDATE_WORD, {
    onCompleted: () => {
      setIsEditing(false)
      toast.success("Word updated successfully")
      refetch()
    },
    onError: (err) => {
      toast.error(`Failed to update: ${err.message}`)
    }
  })

  const [deleteWord, { loading: deleting }] = useMutation(DELETE_WORD, {
    variables: { id: wordId || "" },
    update(cache) {
      cache.modify({
        fields: {
          words(existingWords = {}) {
            return existingWords
          }
        }
      })
    },
    refetchQueries: [GET_WORDS],
    onCompleted: () => {
      toast.success("Word deleted")
      onOpenChange(false)
    },
    onError: (err) => {
      toast.error(`Failed to delete: ${err.message}`)
    }
  })

  // -- Local State for Editing --
  const [formData, setFormData] = useState<any>(null)

  // Инициализация формы при загрузке данных
  useEffect(() => {
    if (data?.word) {
      setFormData({
        text: data.word.text,
        transcription: data.word.transcription || "",
        audioUrl: data.word.audioUrl || "",
        sourceContext: "", 
        meanings: data.word.meanings?.map((m: any) => ({
          id: m.id, 
          partOfSpeech: m.partOfSpeech,
          definitionEn: m.definitionEn || "",
          // Обработка переводов: берем первый из массива или пустую строку
          translations: (m.translationRu && m.translationRu.length > 0) ? m.translationRu : [""], 
          cefrLevel: m.cefrLevel || "",
          imageUrl: m.imageUrl || "",
          tags: m.tags?.map((t: any) => t.name) || [],
          examples: m.examples?.map((e: any) => ({
            sentenceEn: e.sentenceEn,
            sentenceRu: e.sentenceRu || "",
            sourceName: e.sourceName || null
          })) || []
        })) || []
      })
    }
  }, [data, isEditing]) // Reset form when entering edit mode

  // -- Handlers --

  const handleSave = () => {
    if (!wordId || !formData) return

    updateWord({
      variables: {
        id: wordId,
        input: {
          text: formData.text,
          transcription: formData.transcription || null,
          audioUrl: formData.audioUrl || null,
          meanings: formData.meanings.map((m: any) => ({
            partOfSpeech: m.partOfSpeech,
            definitionEn: m.definitionEn || null,
            // Бэкенд ожидает translationRu как String! (согласно твоей схеме mutation CreateMeaningInput)
            // Но мы храним массив в UI. Берем первый элемент.
            translationRu: m.translations[0] || "translation", 
            cefrLevel: m.cefrLevel || null,
            imageUrl: m.imageUrl || null,
            tags: m.tags,
            examples: m.examples.map((e: any) => ({
              sentenceEn: e.sentenceEn,
              sentenceRu: e.sentenceRu || null,
              sourceName: e.sourceName || null
            }))
          }))
        }
      }
    })
  }

  const handleDelete = () => {
    if (confirm("Are you sure you want to delete this word? This cannot be undone.")) {
      deleteWord()
    }
  }

  // Вспомогательные функции для изменения стейта формы
  const updateMeaning = (index: number, field: string, value: any) => {
    const newMeanings = [...formData.meanings]
    newMeanings[index] = { ...newMeanings[index], [field]: value }
    setFormData({ ...formData, meanings: newMeanings })
  }

  const addMeaning = () => {
    setFormData({
      ...formData,
      meanings: [
        ...formData.meanings,
        {
          partOfSpeech: "NOUN",
          definitionEn: "",
          translations: [""],
          examples: [],
          tags: []
        }
      ]
    })
  }

  const removeMeaning = (index: number) => {
    const newMeanings = [...formData.meanings]
    newMeanings.splice(index, 1)
    setFormData({ ...formData, meanings: newMeanings })
  }

  if (!open) return null

  const word = data?.word

  return (
    <Sheet open={open} onOpenChange={(val) => {
      if (!val) setIsEditing(false) // Сбрасываем режим редактирования при закрытии
      onOpenChange(val)
    }}>
      <SheetContent className="w-full sm:max-w-2xl flex flex-col h-full p-0">
        {loading ? (
          <div className="p-6 flex items-center justify-center h-full">
            <Loader2Icon className="animate-spin text-muted-foreground" />
          </div>
        ) : error ? (
          <div className="p-6 text-destructive">Error loading word</div>
        ) : !word ? (
          <div className="p-6">Word not found</div>
        ) : (
          <>
            {/* --- HEADER --- */}
            <SheetHeader className="px-6 py-4 border-b flex-shrink-0 bg-muted/5">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  {/* SheetTitle всегда должен быть для доступности */}
                  <SheetTitle className={isEditing ? "sr-only" : "text-3xl font-bold"}>
                    {word.text}
                  </SheetTitle>
                  {isEditing && (
                    <Input 
                      value={formData?.text} 
                      onChange={e => setFormData({...formData, text: e.target.value})}
                      className="text-2xl font-bold h-auto px-2 py-1 -ml-2 w-full mb-2"
                      aria-label="Word text"
                    />
                  )}
                  
                  <div className="flex items-center gap-2 text-muted-foreground">
                    {isEditing ? (
                      <Input 
                        value={formData?.transcription} 
                        onChange={e => setFormData({...formData, transcription: e.target.value})}
                        placeholder="Transcription"
                        className="h-8 w-40"
                      />
                    ) : (
                      <span className="font-mono text-sm">{word.transcription}</span>
                    )}
                    {!isEditing && word.audioUrl && (
                      <Volume2Icon className="size-4 cursor-pointer hover:text-primary transition-colors" />
                    )}
                  </div>
                </div>

                <div className="flex gap-2 ml-4">
                  {!isEditing && (
                    <>
                      <Button variant="outline" size="sm" onClick={() => setIsEditing(true)}>
                        Edit
                      </Button>
                      <Button variant="ghost" size="icon" className="text-destructive hover:text-destructive hover:bg-destructive/10" onClick={handleDelete} disabled={deleting}>
                        <Trash2Icon className="size-4" />
                      </Button>
                    </>
                  )}
                </div>
              </div>
            </SheetHeader>

            {/* --- BODY --- */}
            <ScrollArea className="flex-1 px-6 py-6">
              <div className="flex flex-col gap-8 pb-10">
                
                {isEditing ? (
                  // === EDIT MODE ===
                  <div className="flex flex-col gap-6">
                    {formData?.meanings.map((meaning: any, mIndex: number) => (
                      <div key={mIndex} className="p-4 border rounded-lg bg-card relative group shadow-sm">
                        <Button 
                          variant="ghost" 
                          size="icon" 
                          className="absolute top-2 right-2 text-muted-foreground hover:text-destructive"
                          onClick={() => removeMeaning(mIndex)}
                        >
                          <XCircleIcon className="size-4" />
                        </Button>

                        <div className="grid gap-4">
                          <div className="grid grid-cols-3 gap-4">
                            <div className="col-span-1">
                              <Label className="text-xs mb-1.5 block text-muted-foreground">Part of Speech</Label>
                              <Select
                                value={meaning.partOfSpeech}
                                onValueChange={(val) => updateMeaning(mIndex, 'partOfSpeech', val)}
                              >
                                <SelectTrigger className="h-9">
                                  <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                  {Object.values(PartOfSpeech).map(pos => (
                                    <SelectItem key={pos} value={pos}>{pos}</SelectItem>
                                  ))}
                                </SelectContent>
                              </Select>
                            </div>
                            <div className="col-span-2">
                              <Label className="text-xs mb-1.5 block text-muted-foreground">Translation (Primary)</Label>
                              <Input 
                                value={meaning.translations[0]} 
                                onChange={e => {
                                  const newTrans = [...meaning.translations]
                                  newTrans[0] = e.target.value
                                  updateMeaning(mIndex, 'translations', newTrans)
                                }}
                                placeholder="Translation"
                              />
                            </div>
                          </div>

                          <div>
                            <Label className="text-xs mb-1.5 block text-muted-foreground">Definition (EN)</Label>
                            <Input 
                              value={meaning.definitionEn} 
                              onChange={e => updateMeaning(mIndex, 'definitionEn', e.target.value)}
                              placeholder="Definition in English"
                            />
                          </div>
                        </div>
                      </div>
                    ))}
                    
                    <Button variant="outline" className="w-full border-dashed text-muted-foreground" onClick={addMeaning}>
                      <PlusIcon className="mr-2 size-4" /> Add Meaning
                    </Button>
                  </div>
                ) : (
                  // === VIEW MODE ===
                  word.meanings?.map((meaning: any) => (
                    <div key={meaning.id} className="flex flex-col gap-3 group border-b pb-6 last:border-0">
                      <div className="flex items-center gap-2">
                        <Badge variant="outline" className="uppercase text-[10px] tracking-wider font-semibold">
                          {meaning.partOfSpeech}
                        </Badge>
                        {meaning.cefrLevel && (
                          <Badge variant="secondary" className="text-[10px]">
                            {meaning.cefrLevel}
                          </Badge>
                        )}
                        <div className="ml-auto text-xs text-muted-foreground flex items-center gap-2">
                          <span>Reviews: {meaning.reviewCount}</span>
                          <span className="w-px h-3 bg-border"></span>
                          <span className={`font-medium uppercase text-[10px] ${
                            meaning.status === 'NEW' ? 'text-blue-500' : 
                            meaning.status === 'LEARNING' ? 'text-yellow-500' :
                            meaning.status === 'MASTERED' ? 'text-green-500' : 'text-orange-500'
                          }`}>{meaning.status}</span>
                        </div>
                      </div>

                      <div className="text-xl font-medium text-foreground">
                        {meaning.translationRu?.join(", ")}
                      </div>

                      {meaning.definitionEn && (
                        <div className="text-sm text-muted-foreground italic">
                          "{meaning.definitionEn}"
                        </div>
                      )}

                      {meaning.examples && meaning.examples.length > 0 && (
                        <div className="mt-2 flex flex-col gap-2 bg-muted/40 rounded-lg p-3">
                          {meaning.examples.map((ex: any) => (
                            <div key={ex.id} className="text-sm grid gap-0.5">
                              <span className="text-foreground font-medium">{ex.sentenceEn}</span>
                              {ex.sentenceRu && <span className="text-muted-foreground text-xs">{ex.sentenceRu}</span>}
                            </div>
                          ))}
                        </div>
                      )}
                      
                      {meaning.tags && meaning.tags.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-1">
                          {meaning.tags.map((tag: any) => (
                            <Badge key={tag.id} variant="secondary" className="text-[10px] px-1.5 h-5 bg-secondary/50 text-secondary-foreground">
                              #{tag.name}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </div>
                  ))
                )}
              </div>
            </ScrollArea>

            {/* --- FOOTER (Edit Mode) --- */}
            {isEditing && (
              <SheetFooter className="px-6 py-4 border-t bg-muted/5 sm:justify-between">
                <Button variant="ghost" onClick={() => setIsEditing(false)}>
                  Cancel
                </Button>
                <Button onClick={handleSave} disabled={updating}>
                  {updating && <Loader2Icon className="mr-2 h-4 w-4 animate-spin" />}
                  Save Changes
                </Button>
              </SheetFooter>
            )}
          </>
        )}
      </SheetContent>
    </Sheet>
  )
}