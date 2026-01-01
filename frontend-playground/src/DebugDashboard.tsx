import { Tabs, TabsContent, TabsList, TabsTrigger } from './components/ui/tabs'
import WordsTab from './components/WordsTab'
import InboxTab from './components/InboxTab'
import StudyTab from './components/StudyTab'
import StatsTab from './components/StatsTab'
import SuggestTab from './components/SuggestTab'

export default function DebugDashboard() {
  return (
    <div className="min-h-screen bg-gray-50 p-4">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Admin Dashboard - My English</h1>
        
        <Tabs defaultValue="words" className="w-full">
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="words">Words</TabsTrigger>
            <TabsTrigger value="inbox">Inbox</TabsTrigger>
            <TabsTrigger value="study">Study</TabsTrigger>
            <TabsTrigger value="stats">Stats</TabsTrigger>
            <TabsTrigger value="suggest">Suggest</TabsTrigger>
          </TabsList>

          <TabsContent value="words" className="mt-4">
            <WordsTab />
          </TabsContent>

          <TabsContent value="inbox" className="mt-4">
            <InboxTab />
          </TabsContent>

          <TabsContent value="study" className="mt-4">
            <StudyTab />
          </TabsContent>

          <TabsContent value="stats" className="mt-4">
            <StatsTab />
          </TabsContent>

          <TabsContent value="suggest" className="mt-4">
            <SuggestTab />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}

