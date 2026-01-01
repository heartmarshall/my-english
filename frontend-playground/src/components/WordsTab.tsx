import { useState } from 'react'
import { useQuery, useMutation } from '@apollo/client/react'
import { GET_WORDS, CREATE_WORD, UPDATE_WORD, DELETE_WORD } from '../graphql/queries'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Textarea } from './ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from './ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/table'
import { Badge } from './ui/badge'

interface MeaningInput {
  partOfSpeech: string
  definitionEn?: string
  translationRu: string
  imageUrl?: string
  examples: Array<{ sentenceEn: string; sentenceRu?: string; sourceName?: string }>
  tags: string[]
}

export default function WordsTab() {
  const [searchText, setSearchText] = useState('')
  const [selectedWordId, setSelectedWordId] = useState<string | null>(null)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  
  const { data, loading, error, refetch } = useQuery(GET_WORDS, {
    variables: {
      filter: searchText ? { search: searchText } : undefined,
      first: 50,
    },
  })

  const [createWord] = useMutation(CREATE_WORD, {
    onCompleted: () => {
      refetch()
      setIsCreateDialogOpen(false)
      resetForm()
    },
  })

  const [updateWord] = useMutation(UPDATE_WORD, {
    onCompleted: () => {
      refetch()
      setIsEditDialogOpen(false)
      setSelectedWordId(null)
      resetForm()
    },
  })

  const [deleteWord] = useMutation(DELETE_WORD, {
    onCompleted: () => {
      refetch()
    },
  })

  const [formData, setFormData] = useState({
    text: '',
    transcription: '',
    audioUrl: '',
    meanings: [] as MeaningInput[],
  })

  const resetForm = () => {
    setFormData({
      text: '',
      transcription: '',
      audioUrl: '',
      meanings: [],
    })
  }

  const handleCreateWord = () => {
    createWord({
      variables: {
        input: {
          text: formData.text,
                          transcription: formData.transcription || undefined,
                          audioUrl: formData.audioUrl || undefined,
          meanings: formData.meanings.map(m => ({
            partOfSpeech: m.partOfSpeech,
            definitionEn: m.definitionEn || undefined,
            translationRu: m.translationRu,
            imageUrl: m.imageUrl || undefined,
            examples: m.examples.map(e => ({
              sentenceEn: e.sentenceEn,
              sentenceRu: e.sentenceRu || undefined,
              sourceName: e.sourceName || undefined,
            })),
            tags: m.tags,
          })),
        },
      },
    })
  }

  const handleUpdateWord = () => {
    if (!selectedWordId) return
    updateWord({
      variables: {
        id: selectedWordId,
        input: {
          text: formData.text,
                          transcription: formData.transcription || undefined,
                          audioUrl: formData.audioUrl || undefined,
          meanings: formData.meanings.map(m => ({
            partOfSpeech: m.partOfSpeech,
            definitionEn: m.definitionEn || undefined,
            translationRu: m.translationRu,
            imageUrl: m.imageUrl || undefined,
            examples: m.examples.map(e => ({
              sentenceEn: e.sentenceEn,
              sentenceRu: e.sentenceRu || undefined,
              sourceName: e.sourceName || undefined,
            })),
            tags: m.tags,
          })),
        },
      },
    })
  }

  const handleDeleteWord = (id: string) => {
    if (confirm('Are you sure you want to delete this word?')) {
      deleteWord({ variables: { id } })
    }
  }

  const handleEditWord = (word: any) => {
    setSelectedWordId(word.id)
    setFormData({
      text: word.text,
      transcription: word.transcription || '',
      audioUrl: word.audioUrl || '',
      meanings: word.meanings?.map((m: any) => ({
        partOfSpeech: m.partOfSpeech,
        definitionEn: m.definitionEn || '',
        translationRu: Array.isArray(m.translationRu) ? m.translationRu[0] : m.translationRu || '',
        imageUrl: m.imageUrl || '',
        examples: m.examples?.map((e: any) => ({
          sentenceEn: e.sentenceEn,
          sentenceRu: e.sentenceRu || '',
          sourceName: e.sourceName || '',
        })) || [],
        tags: m.tags?.map((t: any) => t.name) || [],
      })) || [],
    })
    setIsEditDialogOpen(true)
  }

  const addMeaning = () => {
    setFormData({
      ...formData,
      meanings: [
        ...formData.meanings,
        {
          partOfSpeech: 'NOUN',
          translationRu: '',
          examples: [],
          tags: [],
        },
      ],
    })
  }

  const removeMeaning = (index: number) => {
    setFormData({
      ...formData,
      meanings: formData.meanings.filter((_, i) => i !== index),
    })
  }

  const updateMeaning = (index: number, field: string, value: any) => {
    const newMeanings = [...formData.meanings]
    newMeanings[index] = { ...newMeanings[index], [field]: value }
    setFormData({ ...formData, meanings: newMeanings })
  }

  const addExample = (meaningIndex: number) => {
    const newMeanings = [...formData.meanings]
    newMeanings[meaningIndex].examples.push({
      sentenceEn: '',
      sentenceRu: '',
      sourceName: '',
    })
    setFormData({ ...formData, meanings: newMeanings })
  }

  const removeExample = (meaningIndex: number, exampleIndex: number) => {
    const newMeanings = [...formData.meanings]
    newMeanings[meaningIndex].examples = newMeanings[meaningIndex].examples.filter(
      (_, i) => i !== exampleIndex
    )
    setFormData({ ...formData, meanings: newMeanings })
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>

  const words = (data as any)?.words?.edges?.map((e: any) => e.node) || []

  return (
    <div className="space-y-4">
      <div className="flex gap-4">
        <Input
          placeholder="Search words..."
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          className="flex-1"
        />
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={() => resetForm()}>Create Word</Button>
          </DialogTrigger>
          <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Create Word</DialogTitle>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label>Text *</Label>
                <Input
                  value={formData.text}
                  onChange={(e) => setFormData({ ...formData, text: e.target.value })}
                />
              </div>
              <div>
                <Label>Transcription</Label>
                <Input
                  value={formData.transcription}
                  onChange={(e) => setFormData({ ...formData, transcription: e.target.value })}
                />
              </div>
              <div>
                <Label>Audio URL</Label>
                <Input
                  value={formData.audioUrl}
                  onChange={(e) => setFormData({ ...formData, audioUrl: e.target.value })}
                />
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <Label>Meanings</Label>
                  <Button type="button" size="sm" onClick={addMeaning}>
                    Add Meaning
                  </Button>
                </div>
                {formData.meanings.map((meaning, mIndex) => (
                  <Card key={mIndex} className="mb-4">
                    <CardHeader>
                      <div className="flex justify-between items-center">
                        <CardTitle className="text-sm">Meaning {mIndex + 1}</CardTitle>
                        <Button
                          type="button"
                          size="sm"
                          variant="destructive"
                          onClick={() => removeMeaning(mIndex)}
                        >
                          Remove
                        </Button>
                      </div>
                    </CardHeader>
                    <CardContent className="space-y-2">
                      <div>
                        <Label>Part of Speech *</Label>
                        <select
                          className="w-full p-2 border rounded"
                          value={meaning.partOfSpeech}
                          onChange={(e) => updateMeaning(mIndex, 'partOfSpeech', e.target.value)}
                        >
                          <option value="NOUN">NOUN</option>
                          <option value="VERB">VERB</option>
                          <option value="ADJECTIVE">ADJECTIVE</option>
                          <option value="ADVERB">ADVERB</option>
                          <option value="OTHER">OTHER</option>
                        </select>
                      </div>
                      <div>
                        <Label>Translation RU *</Label>
                        <Input
                          value={meaning.translationRu}
                          onChange={(e) => updateMeaning(mIndex, 'translationRu', e.target.value)}
                        />
                      </div>
                      <div>
                        <Label>Definition EN</Label>
                        <Textarea
                          value={meaning.definitionEn || ''}
                          onChange={(e) => updateMeaning(mIndex, 'definitionEn', e.target.value)}
                        />
                      </div>
                      <div>
                        <Label>Image URL</Label>
                        <Input
                          value={meaning.imageUrl || ''}
                          onChange={(e) => updateMeaning(mIndex, 'imageUrl', e.target.value)}
                        />
                      </div>
                      <div>
                        <div className="flex justify-between items-center mb-2">
                          <Label>Examples</Label>
                          <Button
                            type="button"
                            size="sm"
                            onClick={() => addExample(mIndex)}
                          >
                            Add Example
                          </Button>
                        </div>
                        {meaning.examples.map((example, eIndex) => (
                          <div key={eIndex} className="mb-2 p-2 border rounded">
                            <div className="flex justify-end mb-2">
                              <Button
                                type="button"
                                size="sm"
                                variant="destructive"
                                onClick={() => removeExample(mIndex, eIndex)}
                              >
                                Remove
                              </Button>
                            </div>
                            <div className="space-y-2">
                              <Input
                                placeholder="English sentence *"
                                value={example.sentenceEn}
                                onChange={(e) => {
                                  const newMeanings = [...formData.meanings]
                                  newMeanings[mIndex].examples[eIndex].sentenceEn = e.target.value
                                  setFormData({ ...formData, meanings: newMeanings })
                                }}
                              />
                              <Input
                                placeholder="Russian sentence"
                                value={example.sentenceRu || ''}
                                onChange={(e) => {
                                  const newMeanings = [...formData.meanings]
                                  newMeanings[mIndex].examples[eIndex].sentenceRu = e.target.value
                                  setFormData({ ...formData, meanings: newMeanings })
                                }}
                              />
                              <select
                                className="w-full p-2 border rounded"
                                value={example.sourceName || ''}
                                onChange={(e) => {
                                  const newMeanings = [...formData.meanings]
                                  newMeanings[mIndex].examples[eIndex].sourceName = e.target.value || undefined
                                  setFormData({ ...formData, meanings: newMeanings })
                                }}
                              >
                                <option value="">No source</option>
                                <option value="FILM">FILM</option>
                                <option value="BOOK">BOOK</option>
                                <option value="CHAT">CHAT</option>
                                <option value="VIDEO">VIDEO</option>
                                <option value="PODCAST">PODCAST</option>
                              </select>
                            </div>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>

              <div className="flex gap-2">
                <Button onClick={handleCreateWord}>Create</Button>
                <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                  Cancel
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Words ({(data as any)?.words?.totalCount || 0})</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Text</TableHead>
                <TableHead>Transcription</TableHead>
                <TableHead>Meanings</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {words.map((word: any) => (
                <TableRow key={word.id}>
                  <TableCell className="font-medium">{word.text}</TableCell>
                  <TableCell>{word.transcription || '-'}</TableCell>
                  <TableCell>
                    {word.meanings?.length || 0} meaning(s)
                  </TableCell>
                  <TableCell>
                    {word.meanings?.map((m: any, i: number) => (
                      <Badge key={i} className="mr-1">{m.status}</Badge>
                    ))}
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleEditWord(word)}
                      >
                        Edit
                      </Button>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() => handleDeleteWord(word.id)}
                      >
                        Delete
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Word</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label>Text *</Label>
              <Input
                value={formData.text}
                onChange={(e) => setFormData({ ...formData, text: e.target.value })}
              />
            </div>
            <div>
              <Label>Transcription</Label>
              <Input
                value={formData.transcription}
                onChange={(e) => setFormData({ ...formData, transcription: e.target.value })}
              />
            </div>
            <div>
              <Label>Audio URL</Label>
              <Input
                value={formData.audioUrl}
                onChange={(e) => setFormData({ ...formData, audioUrl: e.target.value })}
              />
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <Label>Meanings</Label>
                <Button type="button" size="sm" onClick={addMeaning}>
                  Add Meaning
                </Button>
              </div>
              {formData.meanings.map((meaning, mIndex) => (
                <Card key={mIndex} className="mb-4">
                  <CardHeader>
                    <div className="flex justify-between items-center">
                      <CardTitle className="text-sm">Meaning {mIndex + 1}</CardTitle>
                      <Button
                        type="button"
                        size="sm"
                        variant="destructive"
                        onClick={() => removeMeaning(mIndex)}
                      >
                        Remove
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent className="space-y-2">
                    <div>
                      <Label>Part of Speech *</Label>
                      <select
                        className="w-full p-2 border rounded"
                        value={meaning.partOfSpeech}
                        onChange={(e) => updateMeaning(mIndex, 'partOfSpeech', e.target.value)}
                      >
                        <option value="NOUN">NOUN</option>
                        <option value="VERB">VERB</option>
                        <option value="ADJECTIVE">ADJECTIVE</option>
                        <option value="ADVERB">ADVERB</option>
                        <option value="OTHER">OTHER</option>
                      </select>
                    </div>
                    <div>
                      <Label>Translation RU *</Label>
                      <Input
                        value={meaning.translationRu}
                        onChange={(e) => updateMeaning(mIndex, 'translationRu', e.target.value)}
                      />
                    </div>
                    <div>
                      <Label>Definition EN</Label>
                      <Textarea
                        value={meaning.definitionEn || ''}
                        onChange={(e) => updateMeaning(mIndex, 'definitionEn', e.target.value)}
                      />
                    </div>
                    <div>
                      <Label>Image URL</Label>
                      <Input
                        value={meaning.imageUrl || ''}
                        onChange={(e) => updateMeaning(mIndex, 'imageUrl', e.target.value)}
                      />
                    </div>
                    <div>
                      <div className="flex justify-between items-center mb-2">
                        <Label>Examples</Label>
                        <Button
                          type="button"
                          size="sm"
                          onClick={() => addExample(mIndex)}
                        >
                          Add Example
                        </Button>
                      </div>
                      {meaning.examples.map((example, eIndex) => (
                        <div key={eIndex} className="mb-2 p-2 border rounded">
                          <div className="flex justify-end mb-2">
                            <Button
                              type="button"
                              size="sm"
                              variant="destructive"
                              onClick={() => removeExample(mIndex, eIndex)}
                            >
                              Remove
                            </Button>
                          </div>
                          <div className="space-y-2">
                            <Input
                              placeholder="English sentence *"
                              value={example.sentenceEn}
                              onChange={(e) => {
                                const newMeanings = [...formData.meanings]
                                newMeanings[mIndex].examples[eIndex].sentenceEn = e.target.value
                                setFormData({ ...formData, meanings: newMeanings })
                              }}
                            />
                            <Input
                              placeholder="Russian sentence"
                              value={example.sentenceRu || ''}
                              onChange={(e) => {
                                const newMeanings = [...formData.meanings]
                                newMeanings[mIndex].examples[eIndex].sentenceRu = e.target.value
                                setFormData({ ...formData, meanings: newMeanings })
                              }}
                            />
                            <select
                              className="w-full p-2 border rounded"
                              value={example.sourceName || ''}
                              onChange={(e) => {
                                const newMeanings = [...formData.meanings]
                                newMeanings[mIndex].examples[eIndex].sourceName = e.target.value || undefined
                                setFormData({ ...formData, meanings: newMeanings })
                              }}
                            >
                              <option value="">No source</option>
                              <option value="FILM">FILM</option>
                              <option value="BOOK">BOOK</option>
                              <option value="CHAT">CHAT</option>
                              <option value="VIDEO">VIDEO</option>
                              <option value="PODCAST">PODCAST</option>
                            </select>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>

            <div className="flex gap-2">
              <Button onClick={handleUpdateWord}>Update</Button>
              <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
                Cancel
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

