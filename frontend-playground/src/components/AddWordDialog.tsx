import { useState, useEffect } from "react"
import { useMutation, useLazyQuery } from "@apollo/client/react"
import { Loader2Icon, Wand2Icon, SparklesIcon } from "lucide-react"
import { toast } from "sonner"

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { CREATE_WORD, GET_WORDS, GET_SUGGEST } from "@/graphql/queries"
import { PartOfSpeech } from "@/gql/graphql"
import { cn } from "@/lib/utils"

// Типы для формы
interface WordFormData {
  text: string
  transcription: string
  translation: string
  definition: string
  partOfSpeech: string
  example: string
}

const INITIAL_DATA: WordFormData = {
  text: "",
  transcription: "",
  translation: "",
  definition: "",
  partOfSpeech: "NOUN",
  example: "",
}

interface AddWordDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  initialText?: string 
  initialData?: Partial<WordFormData>
}

export function AddWordDialog({ 
  open, 
  onOpenChange, 
  initialText = "",
  initialData 
}: AddWordDialogProps) {
  const [formData, setFormData] = useState<WordFormData>(INITIAL_DATA)
  
  // Кэш для данных, полученных из API, чтобы не дергать его каждый раз
  const [cachedSuggestion, setCachedSuggestion] = useState<any>(null)

  // -- GraphQL --
  // Запрос подсказок (включая данные из внешнего словаря)
  const [fetchSuggestion, { loading: fetchLoading }] = useLazyQuery(GET_SUGGEST, {
    fetchPolicy: "network-only"
  })

  // Мутация создания слова
  const [createWord, { loading: createLoading }] = useMutation(CREATE_WORD, {
    refetchQueries: [GET_WORDS],
    onCompleted: () => {
      onOpenChange(false)
      setFormData(INITIAL_DATA)
      setCachedSuggestion(null)
      toast.success("Word created successfully")
    },
    onError: (err: any) => {
      toast.error(`Error creating word: ${err.message}`)
    }
  })

  // -- Effects --
  
  // Инициализация при открытии
  useEffect(() => {
    if (open) {
      setFormData({
        ...INITIAL_DATA,
        text: initialText || "",
        ...initialData
      })
      setCachedSuggestion(null) // Сброс кэша при открытии нового диалога
    }
  }, [open, initialText, initialData])

  // Сбрасываем кэш, если пользователь изменил само слово в инпуте
  useEffect(() => {
    if (cachedSuggestion && cachedSuggestion.text.toLowerCase() !== formData.text.toLowerCase()) {
      setCachedSuggestion(null)
    }
  }, [formData.text, cachedSuggestion])

  // -- Logic --

  // Функция получения данных (из кэша или API)
  const getWordData = async () => {
    if (!formData.text.trim()) {
      toast.error("Please enter a word first")
      return null
    }

    // 1. Если есть кэш и слово совпадает, возвращаем кэш
    if (cachedSuggestion && cachedSuggestion.text.toLowerCase() === formData.text.toLowerCase()) {
      return cachedSuggestion
    }

    // 2. Иначе делаем запрос
    try {
      const { data } = await fetchSuggestion({ 
        variables: { query: formData.text } 
      })

      // Ищем точное совпадение
      const match = (data as any)?.suggest?.find((s: any) => 
        s.text.toLowerCase() === formData.text.toLowerCase()
      )

      if (match) {
        setCachedSuggestion(match)
        return match
      } else {
        toast.info("No detailed data found for this word")
        return null
      }
    } catch (error) {
      toast.error("Failed to fetch word data")
      return null
    }
  }

  // Универсальная функция заполнения
  const handleAutofill = async (field?: keyof WordFormData | 'all') => {
    const data = await getWordData()
    if (!data) return

    setFormData(prev => {
      const newData = { ...prev }
      let filledCount = 0

      // Маппинг данных из API в форму
      const apiValues: Partial<WordFormData> = {
        transcription: data.transcription || "",
        translation: data.translations?.[0] || "",
        definition: data.definition || "", // Используем полученное определение
        // example: data.examples?.[0] || "", // TODO: Добавить поддержку примеров в API Suggestion
      }

      // Функция обновления конкретного поля
      const updateField = (key: keyof WordFormData) => {
        const value = apiValues[key]
        if (value) {
          // Если режим 'all' - заполняем только если поле было пустым
          // Если режим конкретного поля (клик по палочке) - перезаписываем всегда
          if (field === 'all') {
            if (!prev[key]) {
              newData[key] = value
              filledCount++
            }
          } else if (field === key) {
            newData[key] = value
            filledCount++
          }
        }
      }

      if (field === 'all') {
        updateField('transcription')
        updateField('translation')
        updateField('definition')
        updateField('example')
        
        if (filledCount > 0) toast.success(`Auto-filled ${filledCount} fields`)
        else toast.info("No new data to auto-fill")
        
      } else if (field) {
        updateField(field)
        if (filledCount > 0) toast.success(`Updated ${field}`)
        else toast.info(`No data found for ${field}`)
      }

      return newData
    })
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    // Валидация перевода (бэкенд требует)
    if (!formData.translation.trim()) {
      toast.error("Translation is required")
      return
    }

    createWord({
      variables: {
        input: {
          text: formData.text,
          transcription: formData.transcription || null,
          meanings: [
            {
              partOfSpeech: formData.partOfSpeech as PartOfSpeech,
              definitionEn: formData.definition || null,
              translationRu: formData.translation,
              examples: formData.example ? [{ sentenceEn: formData.example }] : [],
              tags: []
            }
          ]
        }
      }
    })
  }

  const isLoading = createLoading || fetchLoading

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[550px]">
        <DialogHeader>
          <DialogTitle>Add New Word</DialogTitle>
        </DialogHeader>
        
        <form onSubmit={handleSubmit} className="grid gap-5 py-4">
          
          {/* Main Word Input + Global Autofill */}
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="text" className="text-right font-bold">
              Word
            </Label>
            <div className="col-span-3 flex gap-2">
              <Input
                id="text"
                value={formData.text}
                onChange={(e) => setFormData({ ...formData, text: e.target.value })}
                placeholder="e.g. Serendipity"
                required
                className="text-lg font-medium"
              />
              <Button 
                type="button" 
                variant="secondary"
                onClick={() => handleAutofill('all')}
                disabled={fetchLoading || !formData.text}
                title="Autofill all empty fields from dictionary"
                className="shrink-0 text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-950/30 hover:bg-indigo-100 dark:hover:bg-indigo-900/50 border-indigo-200 dark:border-indigo-800"
              >
                {fetchLoading ? (
                  <Loader2Icon className="h-4 w-4 animate-spin" />
                ) : (
                  <>
                    <SparklesIcon className="h-4 w-4 mr-2" />
                    Fill All
                  </>
                )}
              </Button>
            </div>
          </div>
          
          <div className="border-t my-1" />

          {/* Fields with individual autofill */}
          
          <InputWithAction
            id="transcription"
            label="Transcription"
            value={formData.transcription}
            onChange={(val) => setFormData({...formData, transcription: val})}
            onAutofill={() => handleAutofill('transcription')}
            placeholder="e.g. /ˌser.ənˈdɪp.ə.ti/"
            className="font-mono"
            isLoading={fetchLoading}
          />

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="pos" className="text-right text-muted-foreground">
              Part of Speech
            </Label>
            <div className="col-span-3">
              <Select 
                value={formData.partOfSpeech} 
                onValueChange={(val) => setFormData({ ...formData, partOfSpeech: val })}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select part of speech" />
                </SelectTrigger>
                <SelectContent>
                  {Object.values(PartOfSpeech).map((pos) => (
                    <SelectItem key={pos} value={pos}>
                      {pos}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          <InputWithAction
            id="translation"
            label="Translation"
            value={formData.translation}
            onChange={(val) => setFormData({...formData, translation: val})}
            onAutofill={() => handleAutofill('translation')}
            placeholder="Russian translation"
            required
            isLoading={fetchLoading}
          />

          <InputWithAction
            id="definition"
            label="Definition"
            value={formData.definition}
            onChange={(val) => setFormData({...formData, definition: val})}
            onAutofill={() => handleAutofill('definition')} 
            placeholder="English definition (optional)"
            isLoading={fetchLoading}
          />

          <InputWithAction
            id="example"
            label="Example"
            value={formData.example}
            onChange={(val) => setFormData({...formData, example: val})}
            onAutofill={() => handleAutofill('example')}
            placeholder="Example sentence (optional)"
            isLoading={fetchLoading}
          />

          <DialogFooter className="mt-4">
            <Button type="submit" disabled={isLoading} className="w-full sm:w-auto">
              {createLoading && <Loader2Icon className="mr-2 h-4 w-4 animate-spin" />}
              Save to Dictionary
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

// --- Helper Component ---

interface InputWithActionProps {
  id: string
  label: string
  value: string
  onChange: (value: string) => void
  onAutofill: () => void
  placeholder?: string
  required?: boolean
  className?: string
  isLoading?: boolean
}

function InputWithAction({ 
  id, label, value, onChange, onAutofill, placeholder, required, className, isLoading 
}: InputWithActionProps) {
  return (
    <div className="grid grid-cols-4 items-center gap-4 group">
      <Label htmlFor={id} className="text-right text-muted-foreground group-focus-within:text-foreground transition-colors">
        {label}
      </Label>
      <div className="col-span-3 relative">
        <Input
          id={id}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className={cn("pr-9", className)} // Отступ справа для кнопки
          placeholder={placeholder}
          required={required}
        />
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="absolute right-0 top-0 h-full w-9 text-muted-foreground hover:text-primary rounded-l-none"
          onClick={onAutofill}
          disabled={isLoading}
          title={`Autofill ${label} from dictionary`}
        >
          {isLoading ? (
            <Loader2Icon className="h-3.5 w-3.5 animate-spin" />
          ) : (
            <Wand2Icon className="h-3.5 w-3.5" />
          )}
        </Button>
      </div>
    </div>
  )
}