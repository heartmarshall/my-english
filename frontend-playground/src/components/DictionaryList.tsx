import { useState } from "react"
import { useQuery } from "@apollo/client/react"
import { SearchIcon, Loader2Icon, MoreHorizontal } from "lucide-react"

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
import { useDebounce } from "@/hooks/use-debounce"
import { GET_WORDS } from "@/graphql/queries"
import { LearningStatus } from "@/gql/graphql"

interface DictionaryListProps {
  onWordClick: (wordId: string) => void
}

export function DictionaryList({ onWordClick }: DictionaryListProps) {
  const [search, setSearch] = useState("")
  const [status, setStatus] = useState<string>("ALL") 
  
  const debouncedSearch = useDebounce(search, 500)

  const filter: any = {}
  if (debouncedSearch) filter.search = debouncedSearch
  if (status !== "ALL") filter.status = status as LearningStatus

  const { data, loading, error, fetchMore } = useQuery(GET_WORDS, {
    variables: { 
      filter,
      first: 50 
    },
    notifyOnNetworkStatusChange: true
  }) as { data: any; loading: boolean; error: any; fetchMore: any }

  const words = data?.words?.edges?.map((edge: any) => edge.node) || []
  const pageInfo = data?.words?.pageInfo
  const totalCount = data?.words?.totalCount || 0

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
                  <TableHead className="w-[200px]">Word</TableHead>
                  <TableHead className="w-[250px]">Translation</TableHead>
                  <TableHead className="w-[100px]">Part of Speech</TableHead>
                  <TableHead className="w-[120px]">Status</TableHead>
                  <TableHead>Tags</TableHead>
                  <TableHead className="w-[150px]">Created</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {words.map((word: any) => {
                  // Берем основное значение (первое) для отображения в таблице
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
                      <TableCell className="text-sm">
                        {translation}
                      </TableCell>
                      <TableCell>
                        {mainMeaning && (
                          <Badge variant="outline" className="text-[10px] uppercase text-muted-foreground font-normal">
                            {mainMeaning.partOfSpeech}
                          </Badge>
                        )}
                      </TableCell>
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
                      <TableCell className="text-xs text-muted-foreground">
                        {new Date(word.createdAt).toLocaleDateString()}
                      </TableCell>
                      <TableCell>
                        <Button variant="ghost" size="icon" className="h-8 w-8">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </TableCell>
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