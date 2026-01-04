import { useState } from "react"
import { gql } from "@apollo/client"
import { useQuery, useMutation } from "@apollo/client/react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Plus, Trash2, BookPlus } from "lucide-react"
import { toast } from "sonner"

const INBOX_ITEMS_QUERY = gql`
  query InboxItems {
    inboxItems {
      id
      text
      context
      createdAt
    }
  }
`

const ADD_TO_INBOX_MUTATION = gql`
  mutation AddToInbox($text: String!, $context: String) {
    addToInbox(text: $text, context: $context) {
      id
      text
      context
      createdAt
    }
  }
`

const DELETE_INBOX_ITEM_MUTATION = gql`
  mutation DeleteInboxItem($id: ID!) {
    deleteInboxItem(id: $id)
  }
`

const CONVERT_INBOX_TO_CARD_MUTATION = gql`
  mutation ConvertInboxToCard($inboxId: ID!, $input: CreateCardInput!) {
    convertInboxToCard(inboxId: $inboxId, input: $input) {
      id
      customText
      progress {
        status
      }
    }
  }
`

interface InboxItem {
  id: string
  text: string
  context?: string | null
  createdAt: string
}

export function InboxPage() {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [isConvertDialogOpen, setIsConvertDialogOpen] = useState<string | null>(null)
  const [newItemForm, setNewItemForm] = useState({
    text: "",
    context: "",
  })
  const [convertForm, setConvertForm] = useState({
    translations: "",
    transcription: "",
    note: "",
    tags: "",
  })

  const { data, loading, refetch } = useQuery<{ inboxItems: InboxItem[] }>(INBOX_ITEMS_QUERY)

  const [addToInbox, { loading: adding }] = useMutation(ADD_TO_INBOX_MUTATION, {
    onCompleted: () => {
      toast.success("Добавлено в inbox!")
      setIsDialogOpen(false)
      setNewItemForm({ text: "", context: "" })
      refetch()
    },
    onError: (error) => {
      toast.error("Ошибка: " + error.message)
    },
  })

  const [deleteInboxItem] = useMutation(DELETE_INBOX_ITEM_MUTATION, {
    onCompleted: () => {
      toast.success("Удалено из inbox")
      refetch()
    },
    onError: (error) => {
      toast.error("Ошибка: " + error.message)
    },
  })

  const [convertToCard, { loading: converting }] = useMutation(CONVERT_INBOX_TO_CARD_MUTATION, {
    onCompleted: () => {
      toast.success("Карточка создана из inbox!")
      setIsConvertDialogOpen(null)
      setConvertForm({ translations: "", transcription: "", note: "", tags: "" })
      refetch()
    },
    onError: (error) => {
      toast.error("Ошибка: " + error.message)
    },
  })

  const handleAddToInbox = () => {
    if (!newItemForm.text.trim()) {
      toast.error("Введите текст")
      return
    }

    addToInbox({
      variables: {
        text: newItemForm.text.trim(),
        context: newItemForm.context.trim() || undefined,
      },
    })
  }

  const handleDelete = (id: string) => {
    if (confirm("Удалить из inbox?")) {
      deleteInboxItem({ variables: { id } })
    }
  }

  const handleConvert = (inboxId: string) => {
    const translations = convertForm.translations
      .split(",")
      .map((t) => t.trim())
      .filter((t) => t.length > 0)

    if (translations.length === 0) {
      toast.error("Введите хотя бы один перевод")
      return
    }

    const tags = convertForm.tags
      .split(",")
      .map((t) => t.trim())
      .filter((t) => t.length > 0)

    const inboxItem = data?.inboxItems.find((item) => item.id === inboxId)
    if (!inboxItem) return

    convertToCard({
      variables: {
        inboxId,
        input: {
          customText: inboxItem.text,
          translations,
          transcription: convertForm.transcription.trim() || undefined,
          note: convertForm.note.trim() || undefined,
          tags: tags.length > 0 ? tags : undefined,
        },
      },
    })
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString("ru-RU", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold mb-2">Inbox</h2>
          <p className="text-muted-foreground">
            Временное хранилище для новых слов и фраз
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="size-4" />
              Добавить в inbox
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Добавить в inbox</DialogTitle>
              <DialogDescription>
                Сохраните слово или фразу для последующей обработки
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid gap-2">
                <Label htmlFor="text">Текст *</Label>
                <Input
                  id="text"
                  placeholder="Например: hello"
                  value={newItemForm.text}
                  onChange={(e) =>
                    setNewItemForm({ ...newItemForm, text: e.target.value })
                  }
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="context">Контекст</Label>
                <Textarea
                  id="context"
                  placeholder="Где вы встретили это слово..."
                  value={newItemForm.context}
                  onChange={(e) =>
                    setNewItemForm({ ...newItemForm, context: e.target.value })
                  }
                  rows={3}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDialogOpen(false)} disabled={adding}>
                Отмена
              </Button>
              <Button onClick={handleAddToInbox} disabled={adding}>
                {adding ? "Добавление..." : "Добавить"}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground">Загрузка...</p>
        </div>
      ) : !data?.inboxItems || data.inboxItems.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground">Inbox пуст</p>
        </div>
      ) : (
        <div className="border rounded-lg overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Текст</TableHead>
                <TableHead>Контекст</TableHead>
                <TableHead className="w-[150px]">Дата</TableHead>
                <TableHead className="w-[200px]">Действия</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.inboxItems.map((item) => (
                <TableRow key={item.id}>
                  <TableCell className="font-medium">{item.text}</TableCell>
                  <TableCell>
                    {item.context ? (
                      <p className="text-sm text-muted-foreground">{item.context}</p>
                    ) : (
                      <span className="text-muted-foreground text-sm">—</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <span className="text-sm text-muted-foreground">
                      {formatDate(item.createdAt)}
                    </span>
                  </TableCell>
                  <TableCell>
                    <div className="flex gap-2">
                      <Dialog
                        open={isConvertDialogOpen === item.id}
                        onOpenChange={(open) =>
                          setIsConvertDialogOpen(open ? item.id : null)
                        }
                      >
                        <DialogTrigger asChild>
                          <Button variant="outline" size="sm">
                            <BookPlus className="size-4" />
                            В карточку
                          </Button>
                        </DialogTrigger>
                        <DialogContent>
                          <DialogHeader>
                            <DialogTitle>Создать карточку из "{item.text}"</DialogTitle>
                            <DialogDescription>
                              Заполните информацию для карточки
                            </DialogDescription>
                          </DialogHeader>
                          <div className="grid gap-4 py-4">
                            <div className="grid gap-2">
                              <Label>Слово</Label>
                              <Input value={item.text} disabled />
                            </div>
                            <div className="grid gap-2">
                              <Label htmlFor="convert-translations">Переводы *</Label>
                              <Input
                                id="convert-translations"
                                placeholder="Через запятую: привет, здравствуй"
                                value={convertForm.translations}
                                onChange={(e) =>
                                  setConvertForm({
                                    ...convertForm,
                                    translations: e.target.value,
                                  })
                                }
                              />
                            </div>
                            <div className="grid gap-2">
                              <Label htmlFor="convert-transcription">Транскрипция</Label>
                              <Input
                                id="convert-transcription"
                                placeholder="/həˈloʊ/"
                                value={convertForm.transcription}
                                onChange={(e) =>
                                  setConvertForm({
                                    ...convertForm,
                                    transcription: e.target.value,
                                  })
                                }
                              />
                            </div>
                            <div className="grid gap-2">
                              <Label htmlFor="convert-tags">Теги</Label>
                              <Input
                                id="convert-tags"
                                placeholder="Через запятую"
                                value={convertForm.tags}
                                onChange={(e) =>
                                  setConvertForm({ ...convertForm, tags: e.target.value })
                                }
                              />
                            </div>
                            <div className="grid gap-2">
                              <Label htmlFor="convert-note">Заметка</Label>
                              <Textarea
                                id="convert-note"
                                value={convertForm.note}
                                onChange={(e) =>
                                  setConvertForm({ ...convertForm, note: e.target.value })
                                }
                                rows={3}
                              />
                            </div>
                          </div>
                          <DialogFooter>
                            <Button
                              variant="outline"
                              onClick={() => setIsConvertDialogOpen(null)}
                              disabled={converting}
                            >
                              Отмена
                            </Button>
                            <Button
                              onClick={() => handleConvert(item.id)}
                              disabled={converting}
                            >
                              {converting ? "Создание..." : "Создать"}
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDelete(item.id)}
                      >
                        <Trash2 className="size-4 text-destructive" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}

