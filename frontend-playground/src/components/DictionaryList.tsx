import { useState, useMemo } from "react"
import { useQuery } from "@apollo/client/react"
import { SearchIcon, Loader2Icon, PlusIcon, ArrowUpDownIcon, Settings2Icon } from "lucide-react"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuCheckboxItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useDebounce } from "@/hooks/use-debounce"
import { GET_WORDS } from "@/graphql/queries"
import { LearningStatus } from "@/gql/graphql"

interface DictionaryListProps {
  onWordClick: (wordId: string) => void
  onAddWord: () => void
}

type SortField = "text" | "translation" | "partOfSpeech" | "status" | "createdAt"
type SortDirection = "asc" | "desc"

interface ColumnVisibility {
  translation: boolean
  partOfSpeech: boolean
  status: boolean
  tags: boolean
  createdAt: boolean
  meanings: boolean
}

export function DictionaryList({ onWordClick, onAddWord }: DictionaryListProps) {
  const [search, setSearch] = useState("")
  const [status, setStatus] = useState<string>("ALL")
  const [sortField, setSortField] = useState<SortField>("createdAt")
  const [sortDirection, setSortDirection] = useState<SortDirection>("desc")
  const [columnVisibility, setColumnVisibility] = useState<ColumnVisibility>({
    translation: true,
    partOfSpeech: true,
    status: true,
    tags: true,
    createdAt: true,
    meanings: true,
  })
  
  const debouncedSearch = useDebounce(search, 500)

  const filter: any = {}
  if (debouncedSearch) filter.search = debouncedSearch
  if (status !== "ALL") filter.status = status as LearningStatus

  // Передаем null вместо пустого объекта, если фильтр пустой
  const filterValue = Object.keys(filter).length > 0 ? filter : null

  const { data, loading, error, fetchMore } = useQuery(GET_WORDS, {
    variables: { 
      filter: filterValue,
      first: 50 
    },
    notifyOnNetworkStatusChange: true
  }) as { data: any; loading: boolean; error: any; fetchMore: any }

  const words = data?.words?.edges?.map((edge: any) => edge.node) || []
  const pageInfo = data?.words?.pageInfo
  const totalCount = data?.words?.totalCount || 0

  // Сортировка слов
  const sortedWords = useMemo(() => {
    const sorted = [...words]
    sorted.sort((a, b) => {
      let aValue: any
      let bValue: any

      switch (sortField) {
        case "text":
          aValue = a.text?.toLowerCase() || ""
          bValue = b.text?.toLowerCase() || ""
          break
        case "translation":
          const aTrans = a.meanings?.[0]?.translationRu?.[0] || a.meanings?.[0]?.translationRu || ""
          const bTrans = b.meanings?.[0]?.translationRu?.[0] || b.meanings?.[0]?.translationRu || ""
          aValue = (typeof aTrans === "string" ? aTrans : aTrans[0] || "").toLowerCase()
          bValue = (typeof bTrans === "string" ? bTrans : bTrans[0] || "").toLowerCase()
          break
        case "partOfSpeech":
          aValue = a.meanings?.[0]?.partOfSpeech || ""
          bValue = b.meanings?.[0]?.partOfSpeech || ""
          break
        case "status":
          aValue = a.meanings?.[0]?.status || ""
          bValue = b.meanings?.[0]?.status || ""
          break
        case "createdAt":
          aValue = new Date(a.createdAt).getTime()
          bValue = new Date(b.createdAt).getTime()
          break
        default:
          return 0
      }

      if (aValue < bValue) return sortDirection === "asc" ? -1 : 1
      if (aValue > bValue) return sortDirection === "asc" ? 1 : -1
      return 0
    })
    return sorted
  }, [words, sortField, sortDirection])

  const handleLoadMore = () => {
    if (pageInfo?.hasNextPage) {
      fetchMore({
        variables: {
          after: pageInfo.endCursor
        },
        updateQuery: (prev: any, { fetchMoreResult }: { fetchMoreResult: any }) => {
          if (!fetchMoreResult) return prev
          return {
            words: {
              ...fetchMoreResult.words,
              edges: [
                ...prev.words.edges,
                ...fetchMoreResult.words.edges
              ]
            }
          }
        }
      })
    }
  }

  const getStatusColorClass = (status: LearningStatus) => {
    switch (status) {
      case LearningStatus.New: return "bg-blue-500/15 text-blue-700 hover:bg-blue-500/25 border-blue-200"; 
      case LearningStatus.Learning: return "bg-yellow-500/15 text-yellow-700 hover:bg-yellow-500/25 border-yellow-200";
      case LearningStatus.Review: return "bg-orange-500/15 text-orange-700 hover:bg-orange-500/25 border-orange-200";
      case LearningStatus.Mastered: return "bg-green-500/15 text-green-700 hover:bg-green-500/25 border-green-200";
      default: return "";
    }
  }

  const handleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDirection(sortDirection === "asc" ? "desc" : "asc")
    } else {
      setSortField(field)
      setSortDirection("asc")
    }
  }

  const toggleColumnVisibility = (column: keyof ColumnVisibility) => {
    setColumnVisibility(prev => ({
      ...prev,
      [column]: !prev[column]
    }))
  }

  return (
    <div className="flex flex-col h-full gap-4">
      {/* --- Toolbar --- */}
      <div className="flex flex-col sm:flex-row gap-4 justify-between items-end sm:items-center bg-card p-4 rounded-lg border shadow-sm">
        <div className="relative w-full sm:w-80">
          <SearchIcon className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Filter words..."
            className="pl-9 bg-background"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        
        <div className="flex items-center gap-4 w-full sm:w-auto">
          <Select value={status} onValueChange={setStatus}>
            <SelectTrigger className="w-[160px] bg-background">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">All Statuses</SelectItem>
              <SelectItem value={LearningStatus.New}>New</SelectItem>
              <SelectItem value={LearningStatus.Learning}>Learning</SelectItem>
              <SelectItem value={LearningStatus.Review}>Review</SelectItem>
              <SelectItem value={LearningStatus.Mastered}>Mastered</SelectItem>
            </SelectContent>
          </Select>

          <Button onClick={onAddWord} className="bg-primary">
            <PlusIcon className="h-4 w-4 mr-2" />
            Add Word
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="icon">
                <Settings2Icon className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-48">
              <DropdownMenuLabel>Column Visibility</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuCheckboxItem
                checked={columnVisibility.translation}
                onCheckedChange={() => toggleColumnVisibility("translation")}
              >
                Translation
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={columnVisibility.meanings}
                onCheckedChange={() => toggleColumnVisibility("meanings")}
              >
                Meanings
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={columnVisibility.partOfSpeech}
                onCheckedChange={() => toggleColumnVisibility("partOfSpeech")}
              >
                Part of Speech
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={columnVisibility.status}
                onCheckedChange={() => toggleColumnVisibility("status")}
              >
                Status
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={columnVisibility.tags}
                onCheckedChange={() => toggleColumnVisibility("tags")}
              >
                Tags
              </DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItem
                checked={columnVisibility.createdAt}
                onCheckedChange={() => toggleColumnVisibility("createdAt")}
              >
                Created
              </DropdownMenuCheckboxItem>
            </DropdownMenuContent>
          </DropdownMenu>
          
          <div className="text-sm text-muted-foreground whitespace-nowrap min-w-[80px] text-right">
            {totalCount} words
          </div>
        </div>
      </div>

      {/* --- Table --- */}
      <div className="rounded-md border bg-card shadow-sm overflow-hidden flex-1 flex flex-col">
        {loading && words.length === 0 ? (
          <div className="flex justify-center items-center py-20 h-64">
            <Loader2Icon className="animate-spin size-8 text-primary/50" />
          </div>
        ) : error ? (
          <div className="text-destructive py-10 text-center">Error loading words</div>
        ) : words.length === 0 ? (
          <div className="text-center py-20 text-muted-foreground">
            No words found.
          </div>
        ) : (
          <div className="overflow-auto">
            <Table>
              <TableHeader className="bg-muted/50">
                <TableRow>
                  <TableHead className="w-[200px]">
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-8 -ml-3 hover:bg-transparent"
                      onClick={() => handleSort("text")}
                    >
                      Word
                      <ArrowUpDownIcon className="ml-2 h-4 w-4" />
                    </Button>
                  </TableHead>
                  {columnVisibility.translation && (
                    <TableHead className="w-[250px]">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 -ml-3 hover:bg-transparent"
                        onClick={() => handleSort("translation")}
                      >
                        Translation
                        <ArrowUpDownIcon className="ml-2 h-4 w-4" />
                      </Button>
                    </TableHead>
                  )}
                  {columnVisibility.meanings && (
                    <TableHead className="w-[200px]">Meanings</TableHead>
                  )}
                  {columnVisibility.partOfSpeech && (
                    <TableHead className="w-[100px]">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 -ml-3 hover:bg-transparent"
                        onClick={() => handleSort("partOfSpeech")}
                      >
                        Part of Speech
                        <ArrowUpDownIcon className="ml-2 h-4 w-4" />
                      </Button>
                    </TableHead>
                  )}
                  {columnVisibility.status && (
                    <TableHead className="w-[120px]">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 -ml-3 hover:bg-transparent"
                        onClick={() => handleSort("status")}
                      >
                        Status
                        <ArrowUpDownIcon className="ml-2 h-4 w-4" />
                      </Button>
                    </TableHead>
                  )}
                  {columnVisibility.tags && (
                    <TableHead>Tags</TableHead>
                  )}
                  {columnVisibility.createdAt && (
                    <TableHead className="w-[150px]">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 -ml-3 hover:bg-transparent"
                        onClick={() => handleSort("createdAt")}
                      >
                        Created
                        <ArrowUpDownIcon className="ml-2 h-4 w-4" />
                      </Button>
                    </TableHead>
                  )}
                </TableRow>
              </TableHeader>
              <TableBody>
                {sortedWords.map((word: any) => {
                  const mainMeaning = word.meanings?.[0];
                  const translation = mainMeaning?.translationRu?.[0] || mainMeaning?.translationRu || "—";
                  
                  return (
                    <TableRow 
                      key={word.id} 
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => onWordClick(word.id)}
                    >
                      <TableCell className="font-medium">
                        <div className="flex flex-col">
                          <span className="text-base">{word.text}</span>
                          {word.transcription && (
                            <span className="text-xs text-muted-foreground font-mono">{word.transcription}</span>
                          )}
                        </div>
                      </TableCell>
                      {columnVisibility.translation && (
                        <TableCell className="text-sm">
                          {translation}
                        </TableCell>
                      )}
                      {columnVisibility.meanings && (
                        <TableCell>
                          <div className="flex flex-col gap-1">
                            {word.meanings?.map((meaning: any, idx: number) => {
                              const meaningTranslation = meaning.translationRu?.[0] || meaning.translationRu || "—";
                              return (
                                <div key={meaning.id} className="text-xs">
                                  <span className="text-muted-foreground">
                                    {idx + 1}. {typeof meaningTranslation === "string" ? meaningTranslation : meaningTranslation[0] || "—"}
                                  </span>
                                  {meaning.partOfSpeech && (
                                    <Badge variant="outline" className="ml-2 text-[9px] uppercase text-muted-foreground font-normal">
                                      {meaning.partOfSpeech}
                                    </Badge>
                                  )}
                                </div>
                              )
                            })}
                          </div>
                        </TableCell>
                      )}
                      {columnVisibility.partOfSpeech && (
                        <TableCell>
                          {mainMeaning && (
                            <Badge variant="outline" className="text-[10px] uppercase text-muted-foreground font-normal">
                              {mainMeaning.partOfSpeech}
                            </Badge>
                          )}
                        </TableCell>
                      )}
                      {columnVisibility.status && (
                        <TableCell>
                          {mainMeaning && (
                            <Badge 
                              variant="outline" 
                              className={`text-[10px] uppercase border-0 ${getStatusColorClass(mainMeaning.status)}`}
                            >
                              {mainMeaning.status}
                            </Badge>
                          )}
                        </TableCell>
                      )}
                      {columnVisibility.tags && (
                        <TableCell>
                          <div className="flex gap-1 flex-wrap">
                            {mainMeaning?.tags?.slice(0, 2).map((tag: any) => (
                              <span key={tag.id} className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded-sm">
                                #{tag.name}
                              </span>
                            ))}
                            {mainMeaning?.tags?.length > 2 && (
                              <span className="text-xs text-muted-foreground pl-1">
                                +{mainMeaning.tags.length - 2}
                              </span>
                            )}
                          </div>
                        </TableCell>
                      )}
                      {columnVisibility.createdAt && (
                        <TableCell className="text-xs text-muted-foreground">
                          {new Date(word.createdAt).toLocaleDateString()}
                        </TableCell>
                      )}
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          </div>
        )}
      </div>

      {/* --- Pagination --- */}
      {pageInfo?.hasNextPage && (
        <div className="flex justify-center py-2">
          <Button 
            variant="outline" 
            onClick={handleLoadMore} 
            disabled={loading}
            className="w-full sm:w-auto"
          >
            {loading && <Loader2Icon className="mr-2 h-4 w-4 animate-spin" />}
            Load More Words
          </Button>
        </div>
      )}
    </div>
  )
}