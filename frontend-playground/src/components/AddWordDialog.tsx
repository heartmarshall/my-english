import { useState, useEffect } from "react"
import { useMutation } from "@apollo/client/react"
import { Loader2Icon } from "lucide-react"

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
import { CREATE_WORD, GET_WORDS } from "@/graphql/queries"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { PartOfSpeech } from "@/gql/graphql"

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
  initialText?: string // Если открываем из SmartSearch с предустановленным текстом
  initialData?: Partial<WordFormData> // Если есть данные из внешнего словаря
}

export function AddWordDialog({ 
  open, 
  onOpenChange, 
  initialText = "",
  initialData 
}: AddWordDialogProps) {
  const [formData, setFormData] = useState<WordFormData>(INITIAL_DATA)

  // Сброс или инициализация формы при открытии
  useEffect(() => {
    if (open) {
      setFormData({
        ...INITIAL_DATA,
        text: initialText || "",
        ...initialData
      })
    }
  }, [open, initialText, initialData])

  const [createWord, { loading }] = useMutation(CREATE_WORD, {
    refetchQueries: [GET_WORDS], // Обновляем список после создания
    onCompleted: () => {
      onOpenChange(false)
      setFormData(INITIAL_DATA)
    },
    onError: (err: Error) => {
      alert(`Error creating word: ${err.message}`)
    }
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    createWord({
      variables: {
        input: {
          text: formData.text,
          transcription: formData.transcription || null,
          meanings: [
            {
              partOfSpeech: formData.partOfSpeech,
              definitionEn: formData.definition || null,
              translationRu: formData.translation, // GraphQL ждет string (из твоей схемы) или [string]
              // В текущей схеме mutations (CreateMeaningInput) translationRu это String!
              // Но в query это [String!]. Проверь generated types.
              // Судя по твоим тестам e2e: translationRu: "привет" (String)
              examples: formData.example ? [{ sentenceEn: formData.example }] : [],
              tags: []
            }
          ]
        }
      }
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add New Word</DialogTitle>
        </DialogHeader>
        
        <form onSubmit={handleSubmit} className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="text" className="text-right">
              Word
            </Label>
            <Input
              id="text"
              value={formData.text}
              onChange={(e) => setFormData({ ...formData, text: e.target.value })}
              className="col-span-3"
              placeholder="e.g. Serendipity"
              required
            />
          </div>
          
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="transcription" className="text-right">
              Transcription
            </Label>
            <Input
              id="transcription"
              value={formData.transcription}
              onChange={(e) => setFormData({ ...formData, transcription: e.target.value })}
              className="col-span-3 font-mono"
              placeholder="e.g. /ˌser.ənˈdɪp.ə.ti/"
            />
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="pos" className="text-right">
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

          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="translation" className="text-right">
              Translation
            </Label>
            <Input
              id="translation"
              value={formData.translation}
              onChange={(e) => setFormData({ ...formData, translation: e.target.value })}
              className="col-span-3"
              placeholder="Russian translation"
              required
            />
          </div>

          <div className="grid grid-cols-4 items-start gap-4">
            <Label htmlFor="definition" className="text-right pt-2">
              Definition
            </Label>
            <Input // Или Textarea если есть
              id="definition"
              value={formData.definition}
              onChange={(e) => setFormData({ ...formData, definition: e.target.value })}
              className="col-span-3"
              placeholder="English definition (optional)"
            />
          </div>

          <div className="grid grid-cols-4 items-start gap-4">
            <Label htmlFor="example" className="text-right pt-2">
              Example
            </Label>
            <Input
              id="example"
              value={formData.example}
              onChange={(e) => setFormData({ ...formData, example: e.target.value })}
              className="col-span-3"
              placeholder="Example sentence (optional)"
            />
          </div>

          <DialogFooter>
            <Button type="submit" disabled={loading}>
              {loading && <Loader2Icon className="mr-2 h-4 w-4 animate-spin" />}
              Save Word
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}