import type { ReactNode } from "react"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { BookOpen, Inbox, GraduationCap, BarChart3, Search } from "lucide-react"

interface LayoutProps {
  children: ReactNode
  currentPage: string
  onPageChange: (page: string) => void
}

export function Layout({ children, currentPage, onPageChange }: LayoutProps) {
  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-6 max-w-7xl">
        <div className="mb-6">
          <h1 className="text-3xl font-bold mb-4">Изучение английского</h1>
          <Tabs value={currentPage} onValueChange={onPageChange}>
            <TabsList className="grid w-full grid-cols-5">
              <TabsTrigger value="cards" className="flex items-center gap-2">
                <BookOpen className="size-4" />
                Мои слова
              </TabsTrigger>
              <TabsTrigger value="inbox" className="flex items-center gap-2">
                <Inbox className="size-4" />
                Inbox
              </TabsTrigger>
              <TabsTrigger value="study" className="flex items-center gap-2">
                <GraduationCap className="size-4" />
                Изучение
              </TabsTrigger>
              <TabsTrigger value="stats" className="flex items-center gap-2">
                <BarChart3 className="size-4" />
                Статистика
              </TabsTrigger>
              <TabsTrigger value="dictionary" className="flex items-center gap-2">
                <Search className="size-4" />
                Словарь
              </TabsTrigger>
            </TabsList>
          </Tabs>
        </div>
        {children}
      </div>
    </div>
  )
}

