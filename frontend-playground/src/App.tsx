import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom';
import { Dashboard } from './pages/Dashboard';
import { Dictionary } from './pages/Dictionary';
import { AddWord } from './pages/AddWord';
import { Study } from './pages/Study';
import { Inbox } from './pages/Inbox';
import { Button } from './components/ui/button';
import { LayoutDashboard, BookOpen, Plus, GraduationCap, Inbox as InboxIcon } from 'lucide-react';
import { cn } from './lib/utils';

function Navigation() {
  const location = useLocation();

  const navItems = [
    { path: '/', label: 'Панель', icon: LayoutDashboard },
    { path: '/dictionary', label: 'Словарь', icon: BookOpen },
    { path: '/add-word', label: 'Добавить', icon: Plus },
    { path: '/study', label: 'Изучение', icon: GraduationCap },
    { path: '/inbox', label: 'Inbox', icon: InboxIcon },
  ];

  return (
    <nav className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center gap-4">
          <h1 className="text-xl font-bold mr-8">My English</h1>
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;
            return (
              <Link key={item.path} to={item.path}>
                <Button
                  variant={isActive ? 'default' : 'ghost'}
                  className={cn('gap-2', isActive && 'bg-primary text-primary-foreground')}
                >
                  <Icon className="size-4" />
                  {item.label}
                </Button>
              </Link>
            );
          })}
        </div>
      </div>
    </nav>
  );
}

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-background">
        <Navigation />
        <main>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/dictionary" element={<Dictionary />} />
            <Route path="/add-word" element={<AddWord />} />
            <Route path="/study" element={<Study />} />
            <Route path="/inbox" element={<Inbox />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;
