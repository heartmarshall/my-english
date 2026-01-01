import { useState, useEffect, useRef } from "react"
import { useLazyQuery } from "@apollo/client/react"
import { SearchIcon, PlusCircleIcon, Loader2Icon } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { GET_SUGGEST } from "@/graphql/queries"
import { useDebounce } from "@/hooks/use-debounce"

interface SmartSearchProps {
  onSelectLocal: (wordId: string) => void
  onSelectDictionary: (suggestion: any) => void
}

export function SmartSearch({ onSelectLocal, onSelectDictionary }: SmartSearchProps) {
  const [query, setQuery] = useState("")
  const [isOpen, setIsOpen] = useState(false)
  const wrapperRef = useRef<HTMLDivElement>(null)
  
  const debouncedQuery = useDebounce(query, 300)
  
  const [getSuggestions, { data, loading }] = useLazyQuery(GET_SUGGEST) as [
    any,
    { data: any; loading: boolean }
  ]

  useEffect(() => {
    if (debouncedQuery.trim().length > 1) {
      getSuggestions({ variables: { query: debouncedQuery } })
      setIsOpen(true)
    } else {
      setIsOpen(false)
    }
  }, [debouncedQuery, getSuggestions])

  // Закрытие при клике вне компонента
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [])

  const suggestions = data?.suggest || []

  return (
    <div ref={wrapperRef} className="relative w-full max-w-xl">
      <div className="relative">
        <SearchIcon className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search or add a word..."
          className="pl-9 h-10 w-full bg-background/50 border-muted-foreground/20 focus:bg-background"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => {
            if (suggestions.length > 0) setIsOpen(true)
          }}
        />
        {loading && (
          <div className="absolute right-3 top-3">
            <Loader2Icon className="h-4 w-4 animate-spin text-muted-foreground" />
          </div>
        )}
      </div>

      {isOpen && (suggestions.length > 0 || query.length > 1) && (
        <div className="absolute top-full mt-1 w-full bg-popover rounded-md border shadow-md overflow-hidden z-50 animate-in fade-in zoom-in-95 duration-100">
          <div className="max-h-[300px] overflow-y-auto p-1">
            {suggestions.length === 0 && !loading ? (
              <div 
                className="p-2 text-sm text-muted-foreground flex items-center gap-2 cursor-pointer hover:bg-accent rounded-sm"
                onClick={() => {
                  onSelectDictionary({ text: query, translations: [] })
                  setIsOpen(false)
                  setQuery("")
                }}
              >
                <PlusCircleIcon className="h-4 w-4" />
                Add "{query}" as new word
              </div>
            ) : (
              suggestions.map((item: any, index: number) => {
                const isLocal = item.origin === "LOCAL"
                
                return (
                  <div
                    key={index}
                    className="flex items-center justify-between px-3 py-2 text-sm rounded-sm cursor-pointer hover:bg-accent hover:text-accent-foreground group"
                    onClick={() => {
                      if (isLocal && item.existingWordId) {
                        onSelectLocal(item.existingWordId)
                      } else {
                        onSelectDictionary(item)
                      }
                      setIsOpen(false)
                      setQuery("")
                    }}
                  >
                    <div className="flex flex-col gap-0.5">
                      <div className="flex items-center gap-2 font-medium">
                        {item.text}
                        {item.transcription && (
                          <span className="text-xs text-muted-foreground font-normal font-mono">
                            {item.transcription}
                          </span>
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground line-clamp-1">
                        {item.translations.join(", ")}
                      </div>
                    </div>

                    <div className="flex items-center">
                      {isLocal ? (
                        <Badge variant="secondary" className="h-5 px-1.5 text-[10px]">
                          My Dictionary
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="h-5 px-1.5 text-[10px] text-muted-foreground group-hover:bg-background">
                          <PlusCircleIcon className="mr-1 h-3 w-3" />
                          Add
                        </Badge>
                      )}
                    </div>
                  </div>
                )
              })
            )}
          </div>
        </div>
      )}
    </div>
  )
}