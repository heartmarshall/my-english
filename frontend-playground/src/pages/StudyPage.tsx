import { useState, useEffect, useMemo } from "react"
import { gql } from "@apollo/client"
import { useQuery, useMutation } from "@apollo/client/react"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { CheckCircle2, XCircle, Clock } from "lucide-react"
import { toast } from "sonner"

const STUDY_QUEUE_QUERY = gql`
  query StudyQueue($filter: StudyFilter) {
    studyQueue(filter: $filter) {
      id
      customText
      customTranscription
      customTranslations
      sense {
        id
        lexeme {
          id
          text
        }
        partOfSpeech
        definition
        translations {
          text
        }
      }
      progress {
        status
        reviewCount
        nextReviewAt
      }
      tags {
        id
        name
      }
    }
  }
`

const MY_CARDS_FOR_STUDY_QUERY = gql`
  query MyCardsForStudy($filter: CardsFilter) {
    myCards(filter: $filter, limit: 1000) {
      id
      customText
      customTranscription
      customTranslations
      sense {
        id
        lexeme {
          id
          text
        }
        partOfSpeech
        definition
        translations {
          text
        }
      }
      progress {
        status
        reviewCount
        nextReviewAt
      }
      tags {
        id
        name
      }
    }
  }
`

const REVIEW_CARD_MUTATION = gql`
  mutation ReviewCard($cardId: ID!, $grade: Int!, $timeTakenMs: Int) {
    reviewCard(cardId: $cardId, grade: $grade, timeTakenMs: $timeTakenMs) {
      card {
        id
        progress {
          status
          reviewCount
        }
      }
      nextReviewInDays
      statusChanged
    }
  }
`

interface StudyCard {
  id: string
  customText?: string | null
  customTranscription?: string | null
  customTranslations?: string[] | null
  sense?: {
    lexeme?: {
      text: string
    } | null
    partOfSpeech?: string | null
    definition?: string | null
    translations?: Array<{ text: string }> | null
  } | null
  progress: {
    status: string
    reviewCount: number
    nextReviewAt?: string | null
  }
  tags: Array<{ id: string; name: string }>
}

interface ReviewCardData {
  reviewCard: {
    card: {
      id: string
      progress: {
        status: string
        reviewCount: number
      }
    }
    nextReviewInDays: number
    statusChanged: boolean
  }
}

export function StudyPage() {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [showAnswer, setShowAnswer] = useState(false)
  const [startTime, setStartTime] = useState(Date.now())
  const [statusFilter, setStatusFilter] = useState<string>("all")

  const { data: studyQueueData, loading: loadingQueue, refetch: refetchQueue } = useQuery<{
    studyQueue: StudyCard[]
  }>(STUDY_QUEUE_QUERY, {
    variables: {
      filter:
        statusFilter === "all"
          ? undefined
          : statusFilter === "new"
            ? undefined // Для new используем myCards
            : {
                statuses:
                  statusFilter === "learning" ? ["LEARNING"] : ["REVIEW"],
              },
    },
    skip: statusFilter === "new", // Пропускаем для новых карточек
  })

  const { data: myCardsData, loading: loadingMyCards, refetch: refetchMyCards } = useQuery<{
    myCards: StudyCard[]
  }>(MY_CARDS_FOR_STUDY_QUERY, {
    variables: {
      filter:
        statusFilter === "all"
          ? { statuses: ["NEW"] }
          : statusFilter === "new"
            ? { statuses: ["NEW"] }
            : undefined,
    },
    skip: statusFilter !== "all" && statusFilter !== "new",
  })

  // Объединяем данные из обоих запросов с useMemo, чтобы избежать бесконечных циклов
  const combinedCards = useMemo(() => {
    const cards: StudyCard[] = []
    if (statusFilter === "all") {
      // Для "all" объединяем новые карточки и карточки из очереди
      const newCards = myCardsData?.myCards.filter((c) => c.progress.status === "NEW") || []
      const queueCards = studyQueueData?.studyQueue || []
      cards.push(...newCards, ...queueCards)
    } else if (statusFilter === "new") {
      cards.push(...(myCardsData?.myCards || []))
    } else {
      cards.push(...(studyQueueData?.studyQueue || []))
    }
    return cards
  }, [statusFilter, myCardsData?.myCards, studyQueueData?.studyQueue])

  const data = useMemo(() => ({ studyQueue: combinedCards }), [combinedCards])
  const loading = statusFilter === "all" 
    ? (loadingMyCards || loadingQueue)
    : statusFilter === "new"
    ? loadingMyCards
    : loadingQueue
  
  const refetch = async () => {
    if (statusFilter === "all") {
      await Promise.all([refetchMyCards(), refetchQueue()])
    } else if (statusFilter === "new") {
      await refetchMyCards()
    } else {
      await refetchQueue()
    }
  }

  // Сброс времени при смене карточки или фильтра
  useEffect(() => {
    setCurrentIndex(0)
    setShowAnswer(false)
    setStartTime(Date.now())
  }, [statusFilter])

  // Сброс времени при смене карточки
  useEffect(() => {
    if (combinedCards.length > 0 && combinedCards[currentIndex]) {
      setStartTime(Date.now())
      setShowAnswer(false)
    }
  }, [currentIndex, combinedCards.length])


  const [reviewCard, { loading: reviewing }] = useMutation<ReviewCardData>(REVIEW_CARD_MUTATION, {
    onCompleted: (mutationData) => {
      const result = mutationData.reviewCard
      if (result.statusChanged) {
        toast.success("Статус изменен!")
      }
      toast.success(`Следующее повторение через ${result.nextReviewInDays.toFixed(1)} дней`)
      setShowAnswer(false)
      if (data?.studyQueue && currentIndex < data.studyQueue.length - 1) {
        setCurrentIndex(currentIndex + 1)
      } else {
        refetch()
        setCurrentIndex(0)
      }
    },
    onError: (error) => {
      toast.error("Ошибка: " + error.message)
    },
  })

  const handleGrade = (grade: number) => {
    const currentCard = data?.studyQueue[currentIndex]
    if (!currentCard) return

    const timeTakenMs = Date.now() - startTime

    reviewCard({
      variables: {
        cardId: currentCard.id,
        grade,
        timeTakenMs,
      },
    })
  }

  const getCardText = (card: StudyCard) => {
    return card.customText || card.sense?.lexeme?.text || "Без названия"
  }

  const getCardTranslations = (card: StudyCard) => {
    if (card.customTranslations && card.customTranslations.length > 0) {
      return card.customTranslations
    }
    if (card.sense?.translations && card.sense.translations.length > 0) {
      return card.sense.translations.map((t) => t.text)
    }
    return []
  }

  const currentCard = data?.studyQueue?.[currentIndex]
  const progress = data?.studyQueue && data.studyQueue.length > 0
    ? ((currentIndex + 1) / data.studyQueue.length) * 100
    : 0

  return (
    <div className="max-w-2xl mx-auto">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-2xl font-bold">Изучение</h2>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="Фильтр по статусу" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Все карточки</SelectItem>
              <SelectItem value="new">Новые</SelectItem>
              <SelectItem value="learning">Изучаю</SelectItem>
              <SelectItem value="review">К повторению</SelectItem>
            </SelectContent>
          </Select>
        </div>
        {loading ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">Загрузка карточек для изучения...</p>
          </div>
        ) : !data?.studyQueue || data.studyQueue.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground mb-4">Нет карточек для изучения</p>
            <p className="text-sm text-muted-foreground">
              Добавьте карточки или дождитесь времени повторения
            </p>
          </div>
        ) : (
          <>
            <div className="flex items-center justify-between mb-2">
              <p className="text-muted-foreground">
                Карточка {currentIndex + 1} из {data.studyQueue.length}
              </p>
              <Badge variant="outline">
                Повторений: {currentCard?.progress.reviewCount || 0}
              </Badge>
            </div>
            <div className="w-full bg-muted rounded-full h-2">
              <div
                className="bg-primary h-2 rounded-full transition-all"
                style={{ width: `${progress}%` }}
              />
            </div>
          </>
        )}
      </div>

      {loading ? null : !data?.studyQueue || data.studyQueue.length === 0 ? null : currentCard ? (
        <Card className="p-8 min-h-[400px] flex flex-col">
          {!showAnswer ? (
            <>
              <div className="flex-1 flex items-center justify-center">
                <div className="text-center">
                  <h3 className="text-3xl font-bold mb-4">{getCardText(currentCard)}</h3>
                  {currentCard.customTranscription && (
                    <p className="text-xl text-muted-foreground mb-4">
                      [{currentCard.customTranscription}]
                    </p>
                  )}
                  {currentCard.sense?.partOfSpeech && (
                    <Badge variant="secondary" className="mb-4">
                      {currentCard.sense.partOfSpeech}
                    </Badge>
                  )}
                </div>
              </div>
              <div className="mt-auto">
                <Button
                  onClick={() => {
                    setShowAnswer(true)
                  }}
                  className="w-full"
                  size="lg"
                >
                  Показать ответ
                </Button>
              </div>
            </>
          ) : (
            <>
              <div className="flex-1 space-y-4">
                <div className="text-center mb-6">
                  <h3 className="text-3xl font-bold mb-4">{getCardText(currentCard)}</h3>
                  {currentCard.customTranscription && (
                    <p className="text-xl text-muted-foreground mb-2">
                      [{currentCard.customTranscription}]
                    </p>
                  )}
                </div>

                <div className="space-y-3">
                  <div>
                    <p className="text-sm text-muted-foreground mb-1">Переводы:</p>
                    <div className="space-y-1">
                      {getCardTranslations(currentCard).map((translation, idx) => (
                        <p key={idx} className="text-lg">• {translation}</p>
                      ))}
                    </div>
                  </div>

                  {currentCard.sense?.definition && (
                    <div>
                      <p className="text-sm text-muted-foreground mb-1">Определение:</p>
                      <p className="text-base">{currentCard.sense.definition}</p>
                    </div>
                  )}

                  {currentCard.tags.length > 0 && (
                    <div>
                      <p className="text-sm text-muted-foreground mb-2">Теги:</p>
                      <div className="flex flex-wrap gap-2">
                        {currentCard.tags.map((tag) => (
                          <Badge key={tag.id} variant="secondary">
                            {tag.name}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>

              <div className="mt-auto pt-6 border-t">
                <p className="text-sm text-muted-foreground mb-4 text-center">
                  Насколько хорошо вы знаете это слово?
                </p>
                <div className="grid grid-cols-5 gap-2">
                  <Button
                    variant="outline"
                    onClick={() => handleGrade(1)}
                    disabled={reviewing}
                    className="flex flex-col h-auto py-3"
                  >
                    <XCircle className="size-5 mb-1 text-red-500" />
                    <span className="text-xs">1</span>
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => handleGrade(2)}
                    disabled={reviewing}
                    className="flex flex-col h-auto py-3"
                  >
                    <XCircle className="size-5 mb-1 text-orange-500" />
                    <span className="text-xs">2</span>
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => handleGrade(3)}
                    disabled={reviewing}
                    className="flex flex-col h-auto py-3"
                  >
                    <Clock className="size-5 mb-1 text-yellow-500" />
                    <span className="text-xs">3</span>
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => handleGrade(4)}
                    disabled={reviewing}
                    className="flex flex-col h-auto py-3"
                  >
                    <CheckCircle2 className="size-5 mb-1 text-green-500" />
                    <span className="text-xs">4</span>
                  </Button>
                  <Button
                    variant="outline"
                    onClick={() => handleGrade(5)}
                    disabled={reviewing}
                    className="flex flex-col h-auto py-3"
                  >
                    <CheckCircle2 className="size-5 mb-1 text-green-600" />
                    <span className="text-xs">5</span>
                  </Button>
                </div>
              </div>
            </>
          )}
        </Card>
      ) : null}
    </div>
  )
}

