import { Badge } from "@/components/ui/badge"
import { Card } from "@/components/ui/card"
import { Volume2Icon } from "lucide-react"

interface WordListCardProps {
  word: any // Type from GraphQL
  onClick: () => void
}

export function WordListCard({ word, onClick }: WordListCardProps) {
  // Получаем основной перевод (первый попавшийся из meanings)
  const mainTranslation = word.meanings?.[0]?.translationRu?.[0] || word.meanings?.[0]?.translationRu || "";
  const mainMeaning = word.meanings?.[0];

  return (
    <Card 
      className="p-4 cursor-pointer hover:border-primary/50 transition-colors group relative overflow-hidden"
      onClick={onClick}
    >
      <div className="flex justify-between items-start mb-2">
        <h3 className="font-bold text-lg">{word.text}</h3>
        {word.audioUrl && (
          <Volume2Icon className="size-4 text-muted-foreground opacity-50 group-hover:opacity-100" />
        )}
      </div>
      
      <div className="text-sm text-muted-foreground mb-3 font-mono">
        {word.transcription}
      </div>

      <div className="text-sm font-medium line-clamp-2 min-h-[1.25rem]">
        {mainTranslation}
      </div>

      <div className="flex gap-2 mt-4 items-center">
        {mainMeaning && (
          <Badge variant="secondary" className="text-[10px] px-1.5 h-5 uppercase">
            {mainMeaning.partOfSpeech}
          </Badge>
        )}
        
        {/* Status dot */}
        {mainMeaning && (
          <div className="ml-auto flex items-center gap-1.5 text-[10px] text-muted-foreground">
            <div className={`size-2 rounded-full ${getStatusColor(mainMeaning.status)}`} />
            <span className="uppercase">{mainMeaning.status}</span>
          </div>
        )}
      </div>
    </Card>
  )
}

function getStatusColor(status: string) {
  switch (status) {
    case 'NEW': return 'bg-blue-500';
    case 'LEARNING': return 'bg-yellow-500';
    case 'REVIEW': return 'bg-orange-500';
    case 'MASTERED': return 'bg-green-500';
    default: return 'bg-gray-500';
  }
}