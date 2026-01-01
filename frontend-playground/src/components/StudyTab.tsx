import { useState } from 'react'
import { useQuery, useMutation } from '@apollo/client/react'
import { GET_STUDY_QUEUE, REVIEW_MEANING } from '../graphql/queries'
import { Button } from './ui/button'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Label } from './ui/label'

export default function StudyTab() {
  const [limit, setLimit] = useState(10)
  const [selectedMeaning, setSelectedMeaning] = useState<any>(null)
  const [showAnswer, setShowAnswer] = useState(false)

  const { data, loading, error, refetch } = useQuery(GET_STUDY_QUEUE, {
    variables: { limit },
  })

  const [reviewMeaning] = useMutation(REVIEW_MEANING, {
    onCompleted: () => {
      refetch()
      setSelectedMeaning(null)
      setShowAnswer(false)
    },
  })

  const handleReview = (meaningId: string, grade: number) => {
    reviewMeaning({
      variables: {
        meaningId,
        grade,
      },
    })
  }

  const handleStartReview = (meaning: any) => {
    setSelectedMeaning(meaning)
    setShowAnswer(false)
  }

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error: {error.message}</div>

  const studyQueue = (data as any)?.studyQueue || []

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle>Study Queue Settings</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4">
            <Label>Limit:</Label>
            <input
              type="number"
              min="1"
              max="100"
              value={limit}
              onChange={(e) => setLimit(parseInt(e.target.value) || 10)}
              className="w-20 p-2 border rounded"
            />
            <Button onClick={() => refetch()}>Refresh</Button>
          </div>
        </CardContent>
      </Card>

      {selectedMeaning ? (
        <Card>
          <CardHeader>
            <CardTitle>Review Meaning</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label>Word ID:</Label>
              <p className="font-mono text-sm">{selectedMeaning.wordId}</p>
            </div>
            <div>
              <Label>Part of Speech:</Label>
              <Badge>{selectedMeaning.partOfSpeech}</Badge>
            </div>
            <div>
              <Label>Status:</Label>
              <Badge>{selectedMeaning.status}</Badge>
            </div>
            <div>
              <Label>Review Count:</Label>
              <p>{selectedMeaning.reviewCount}</p>
            </div>
            {selectedMeaning.definitionEn && (
              <div>
                <Label>Definition:</Label>
                <p>{selectedMeaning.definitionEn}</p>
              </div>
            )}
            {!showAnswer ? (
              <Button onClick={() => setShowAnswer(true)}>Show Answer</Button>
            ) : (
              <div className="space-y-4">
                <div>
                  <Label>Translation:</Label>
                  <p className="font-bold">
                    {Array.isArray(selectedMeaning.translationRu)
                      ? selectedMeaning.translationRu.join(', ')
                      : selectedMeaning.translationRu}
                  </p>
                </div>
                {selectedMeaning.examples && selectedMeaning.examples.length > 0 && (
                  <div>
                    <Label>Examples:</Label>
                    <ul className="list-disc list-inside space-y-1">
                      {selectedMeaning.examples.map((ex: any, i: number) => (
                        <li key={i}>
                          <p>{ex.sentenceEn}</p>
                          {ex.sentenceRu && <p className="text-gray-600">{ex.sentenceRu}</p>}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
                <div>
                  <Label>Rate your knowledge (1-5):</Label>
                  <div className="flex gap-2 mt-2">
                    {[1, 2, 3, 4, 5].map((grade) => (
                      <Button
                        key={grade}
                        onClick={() => handleReview(selectedMeaning.id, grade)}
                        variant={grade <= 2 ? 'destructive' : grade === 3 ? 'outline' : 'default'}
                      >
                        {grade}
                      </Button>
                    ))}
                  </div>
                </div>
                <Button variant="outline" onClick={() => handleStartReview(selectedMeaning)}>
                  Skip
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardHeader>
            <CardTitle>Study Queue ({studyQueue.length} items)</CardTitle>
          </CardHeader>
          <CardContent>
            {studyQueue.length === 0 ? (
              <p>No items in study queue</p>
            ) : (
              <div className="space-y-2">
                {studyQueue.map((meaning: any) => (
                  <Card key={meaning.id}>
                    <CardContent className="pt-4">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <p className="font-mono text-sm">Word ID: {meaning.wordId}</p>
                          <p className="font-mono text-sm">Meaning ID: {meaning.id}</p>
                          <div className="flex gap-2 mt-2">
                            <Badge>{meaning.partOfSpeech}</Badge>
                            <Badge>{meaning.status}</Badge>
                            <Badge variant="outline">Reviews: {meaning.reviewCount}</Badge>
                          </div>
                          {meaning.nextReviewAt && (
                            <p className="text-sm text-gray-600 mt-1">
                              Next review: {new Date(meaning.nextReviewAt).toLocaleString()}
                            </p>
                          )}
                        </div>
                        <Button onClick={() => handleStartReview(meaning)}>Start Review</Button>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}

