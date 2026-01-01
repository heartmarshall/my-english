import { useState } from 'react'
import { useLazyQuery } from '@apollo/client/react'
import { GET_SUGGEST } from '../graphql/queries'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/table'

export default function SuggestTab() {
  const [query, setQuery] = useState('')
  const [getSuggest, { data, loading, error }] = useLazyQuery(GET_SUGGEST)

  const handleSearch = () => {
    if (query.trim()) {
      getSuggest({ variables: { query: query.trim() } })
    }
  }

  const suggestions = (data as any)?.suggest || []

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Search Suggestions</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label>Query</Label>
            <div className="flex gap-2">
              <Input
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="Enter search query"
              />
              <Button onClick={handleSearch}>Search</Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {loading && <div>Loading...</div>}
      {error && <div>Error: {error.message}</div>}

      {suggestions.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Suggestions ({suggestions.length})</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Text</TableHead>
                  <TableHead>Transcription</TableHead>
                  <TableHead>Translations</TableHead>
                  <TableHead>Origin</TableHead>
                  <TableHead>Existing Word ID</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {suggestions.map((suggestion: any, index: number) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">{suggestion.text}</TableCell>
                    <TableCell>{suggestion.transcription || '-'}</TableCell>
                    <TableCell>
                      {suggestion.translations?.join(', ') || '-'}
                    </TableCell>
                    <TableCell>
                      <Badge>{suggestion.origin}</Badge>
                    </TableCell>
                    <TableCell>
                      {suggestion.existingWordId ? (
                        <span className="font-mono text-sm">{suggestion.existingWordId}</span>
                      ) : (
                        '-'
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      {!loading && !error && suggestions.length === 0 && query && (
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-gray-500">No suggestions found</p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

