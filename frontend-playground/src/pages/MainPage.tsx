import { useState } from "react"
import { DictionaryList } from "@/components/DictionaryList"
import { WordCardSheet } from "@/components/WordCardSheet"
import { AddWordDialog } from "@/components/AddWordDialog"
import { SmartSearch } from "@/components/SmartSearch"
import { BookMarkedIcon } from "lucide-react"

export function MainPage() {
  const [selectedWordId, setSelectedWordId] = useState<string | null>(null)
  
  // State for Add Dialog
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [addDialogInitialText, setAddDialogInitialText] = useState("")
  const [addDialogInitialData, setAddDialogInitialData] = useState<any>({})

  // Handlers for SmartSearch
  const handleSelectLocal = (wordId: string) => {
    setSelectedWordId(wordId)
  }

  const handleSelectDictionary = (suggestion: any) => {
    setAddDialogInitialText(suggestion.text)
    setAddDialogInitialData({
      transcription: suggestion.transcription,
      translation: suggestion.translations?.[0] || "",
      // Можно добавить маппинг других полей, если они приходят из suggest
    })
    setIsAddDialogOpen(true)
  }

  return (
    <div className="min-h-screen bg-background">
      {/* --- HEADER --- */}
      <header className="sticky top-0 z-40 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between gap-4">
          <div className="flex items-center gap-2 font-bold text-xl hidden md:flex text-primary">
            <BookMarkedIcon className="h-6 w-6" />
            <span>My English</span>
          </div>

          <div className="flex-1 max-w-xl mx-auto">
            <SmartSearch 
              onSelectLocal={handleSelectLocal}
              onSelectDictionary={handleSelectDictionary}
            />
          </div>

          {/* Spacer for right side actions if needed, or Profile menu */}
          <div className="w-10 md:w-auto" />
        </div>
      </header>

      {/* --- CONTENT --- */}
      <main className="container mx-auto px-4 py-6">
        <DictionaryList 
          onWordClick={handleSelectLocal}
          onAddWord={() => setIsAddDialogOpen(true)}
        />
      </main>

      {/* --- DIALOGS & SHEETS --- */}
      
      {/* 1. Детальный просмотр / Редактирование */}
      {/* Важно: DictionaryList тоже может открывать этот Sheet. 
          Сейчас DictionaryList управляет своим состоянием selectedWordId внутри себя.
          
          Правильнее поднять состояние selectedWordId сюда (в MainPage), 
          чтобы поиск тоже мог его открывать.
          
          Но чтобы не переписывать DictionaryList полностью прямо сейчас, 
          мы можем рендерить WordCardSheet здесь ТОЛЬКО для поиска, 
          а в DictionaryList пусть живет свой экземпляр (или отрефакторим).
          
          ЛУЧШИЙ ВАРИАНТ: Удалить WordCardSheet из DictionaryList и передать 
          ему пропс onSelectWord. Давай сделаем это правильно.
      */}
      <WordCardSheet 
        wordId={selectedWordId} 
        open={!!selectedWordId} 
        onOpenChange={(open) => !open && setSelectedWordId(null)}
      />

      {/* 2. Добавление слова */}
      <AddWordDialog 
        open={isAddDialogOpen} 
        onOpenChange={setIsAddDialogOpen}
        initialText={addDialogInitialText}
        initialData={addDialogInitialData}
      />
    </div>
  )
}