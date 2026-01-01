# Frontend Playground - Debug Dashboard

Админ-панель для отладки и управления backend сервисом My English.

## Технологии

- React 19 + TypeScript
- Vite
- Apollo Client (GraphQL)
- Tailwind CSS
- Shadcn UI

## Быстрый старт

### 1. Установка зависимостей

```bash
npm install
```

### 2. Запуск dev сервера

```bash
npm run dev
```

Приложение будет доступно на `http://localhost:5173`

### 3. Убедитесь, что backend запущен

Backend должен быть доступен на `http://localhost:8080`

Прокси настроен в `vite.config.ts` для автоматического перенаправления запросов `/graphql` на backend.

## Функционал

### Словарь
- Просмотр всех слов из базы данных
- Таблица с информацией о словах, значениях, статусах и тегах
- Пагинация

### Очередь изучения
- Просмотр слов, готовых к повторению
- Интерактивное ревью с оценкой от 1 до 5
- Автоматическое обновление после ревью

### Статистика
- Общее количество слов
- Количество изученных слов
- Количество изучаемых слов
- Количество слов к повторению
- Автоматическое обновление каждые 5 секунд

### Добавление слов
- Диалог для добавления новых слов
- Поддержка транскрипции, переводов и определений

## Команды

```bash
# Запуск dev сервера
npm run dev

# Сборка для production
npm run build

# Просмотр production сборки
npm run preview

# Генерация GraphQL типов (если нужно)
npm run codegen
```

## Структура проекта

```
src/
  ├── components/          # React компоненты
  │   ├── ui/             # Shadcn UI компоненты
  │   ├── WordsList.tsx   # Список слов
  │   ├── StatsCard.tsx   # Статистика
  │   ├── StudyQueue.tsx  # Очередь изучения
  │   └── AddWordDialog.tsx # Диалог добавления слова
  ├── graphql/
  │   └── queries.ts      # GraphQL запросы и мутации
  ├── App.tsx             # Главный компонент
  ├── DebugDashboard.tsx  # Основная панель
  └── main.tsx            # Точка входа
```

## Настройка

### Изменение URL backend

Измените прокси в `vite.config.ts`:

```typescript
server: {
  proxy: {
    '/graphql': {
      target: 'http://localhost:8080', // Ваш backend URL
      changeOrigin: true,
    },
  },
}
```

Или измените URI в `src/main.tsx`:

```typescript
const client = new ApolloClient({
  uri: 'http://localhost:8080/graphql', // Прямой URL
  cache: new InMemoryCache(),
});
```
