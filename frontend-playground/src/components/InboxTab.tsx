import { useState } from 'react'
import { useQuery, useMutation } from '@apollo/client/react'
import { GET_INBOX_ITEMS, ADD_TO_INBOX, DELETE_INBOX_ITEM, CONVERT_INBOX_ITEM } from '../graphql/queries'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Textarea } from './ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/table'

export default function InboxTab() {
  const [text, setText] = useState('')
  const [sourceContext, setSourceContext] = useState('')
  const [selectedInboxId, setSelectedInboxId] = useState<string | null>(null)
  const [isConvertDialogOpen, setIsConvertDialogOpen] = useState(false)
  const [convertFormData, setConvertFormData] = useState({
    text: '',
    transcription: '',
    audioUrl: '',
    meanings: [] as any[],
  })

  const { data, loading, error, refetch } = useQuery(GET_INBOX_ITEMS)

  const [addToInbox] = useMutation(ADD_TO_INBOX, {
    onCompleted: () => {
      refetch()
      setText('')
      setSourceContext('')
    },
  })

  const [deleteInboxItem] = useMutation(DELETE_INBOX_ITEM, {
    onCompleted: () => {
      refetch()
    },
  })

  const [convertInboxItem] = useMutation(CONVERT_INBOX_ITEM, {
    onCompleted: () => {
      refetch()
      setIsConvertDialogOpen(false)
      setSelectedInboxId(null)
      setConvertFormData({
        text: '',
        transcription: '',
        audioUrl: '',
        meanings: [],
      })
    },
  })

  const handleAddToInbox = () => {
    if (!text.trim()) return
    addToInbox({
      variables: {
        text: text.trim(),
        sourceContext: sourceContext.trim() || undefined,
      },
    })
  }

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this inbox item?')) {
      deleteInboxItem({ variables: { id } })
    }
  }

  const handleConvert = (item: any) => {
    setSelectedInboxId(item.id)
    setConvertFormData({
      text: item.text,
      transcription: '',
      audioUrl: '',
      meanings: [
        {
          partOfSpeech: 'NOUN',
          translationRu: '',
          examples: [],
          tags: [],
        },
      ],
    })
    setIsConvertDialogOpen(true)
  }

  const handleConvertSubmit = () => {
    if (!selectedInboxId) return
    convertInboxItem({
      variables: {
        inboxId: selectedInboxId,
        input: {
          text: convertFormData.text,
            transcription: convertFormData.transcription || undefined,
            audioUrl: convertFormData.audioUrl || undefined,
          meanings: convertFormData.meanings.map((m) => ({
            partOfSpeech: m.partOfSpeech,
            definitionEn: m.definitionEn || undefined,
            translationRu: m.translationRu,
            imageUrl: m.imageUrl || undefined,
            examples: m.examples.map((e: any) => ({
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

  const addMeaning = () => {
    setConvertFormData({
      ...convertFormData,
      meanings: [
        ...convertFormData.meanings,
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
    setConvertFormData({
      ...convertFormData,
      meanings: convertFormData.meanings.filter((_, i) => i !== index),
    })
  }

  const updateMeaning = (index: number, field: string, value: any) => {
    const newMeanings = [...convertFormData.meanings]
    newMeanings[index] = { ...newMeanings[index], [field]: value }
    setConvertFormData({ ...convertFormData, meanings: newMeanings })
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>

  const inboxItems = (data as any)?.inboxItems || []

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Add to Inbox</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label>Text *</Label>
            <Input
              value={text}
              onChange={(e) => setText(e.target.value)}
              placeholder="Enter word text"
            />
          </div>
          <div>
            <Label>Source Context</Label>
            <Textarea
              value={sourceContext}
              onChange={(e) => setSourceContext(e.target.value)}
              placeholder="e.g., Harry Potter, page 50"
            />
          </div>
          <Button onClick={handleAddToInbox}>Add to Inbox</Button>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Inbox Items ({inboxItems.length})</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Text</TableHead>
                <TableHead>Source Context</TableHead>
                <TableHead>Created At</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {inboxItems.map((item: any) => (
                <TableRow key={item.id}>
                  <TableCell className="font-medium">{item.text}</TableCell>
                  <TableCell>{item.sourceContext || '-'}</TableCell>
                  <TableCell>{new Date(item.createdAt).toLocaleString()}</TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Button size="sm" onClick={() => handleConvert(item)}>
                        Convert
                      </Button>
                      <Button
                        size="sm"
                        variant="destructive"
                        onClick={() => handleDelete(item.id)}
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

      <Dialog open={isConvertDialogOpen} onOpenChange={setIsConvertDialogOpen}>
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Convert Inbox Item to Word</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label>Text *</Label>
              <Input
                value={convertFormData.text}
                onChange={(e) =>
                  setConvertFormData({ ...convertFormData, text: e.target.value })
                }
              />
            </div>
            <div>
              <Label>Transcription</Label>
              <Input
                value={convertFormData.transcription}
                onChange={(e) =>
                  setConvertFormData({ ...convertFormData, transcription: e.target.value })
                }
              />
            </div>
            <div>
              <Label>Audio URL</Label>
              <Input
                value={convertFormData.audioUrl}
                onChange={(e) =>
                  setConvertFormData({ ...convertFormData, audioUrl: e.target.value })
                }
              />
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <Label>Meanings</Label>
                <Button type="button" size="sm" onClick={addMeaning}>
                  Add Meaning
                </Button>
              </div>
              {convertFormData.meanings.map((meaning, mIndex) => (
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
                  </CardContent>
                </Card>
              ))}
            </div>

            <div className="flex gap-2">
              <Button onClick={handleConvertSubmit}>Convert</Button>
              <Button variant="outline" onClick={() => setIsConvertDialogOpen(false)}>
                Cancel
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

