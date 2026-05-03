#!/bin/bash

# ===== NEXT.JS 01-basics =====
cat > /workspace/learning-path/06-nextjs/01-nextjs-basics/README.md << 'EOF'
# Модуль 1: Основы Next.js

## 📚 Теория

### 1.1 Что такое Next.js?
Next.js - это React фреймворк для production-grade приложений с SSR, SSG и другими возможностями.

**Ключевые преимущества:**
- Server-Side Rendering (SSR)
- Static Site Generation (SSG)
- Incremental Static Regeneration (ISR)
- API Routes
- Image Optimization
- Built-in CSS/Sass support
- Automatic Code Splitting

### 1.2 App Router vs Pages Router

```jsx
// App Router (Next.js 13+)
// app/page.tsx
export default function Home() {
  return <main><h1>Home</h1></main>;
}

// app/about/page.tsx
export default function About() {
  return <h1>About</h1>;
}
```

### 1.3 Рендеринг в Next.js

```jsx
// Server Component (по умолчанию)
async function ProductPage({ params }) {
  const product = await fetchProduct(params.id);
  return <div>{product.name}</div>;
}

// Client Component
'use client';
import { useState } from 'react';

function Counter() {
  const [count, setCount] = useState(0);
  return <button onClick={() => setCount(c => c + 1)}>{count}</button>;
}
```

### 1.4 Data Fetching

```jsx
// В Server Component
async function Page() {
  const data = await fetch('https://api.example.com/data', {
    cache: 'force-cache' // SSG
    // cache: 'no-store' // SSR
  });
  return <div>{data.title}</div>;
}
```

---

## 🎯 Практические задания

### Задание 1: Первый проект (Уровень: Начинающий)
Создайте Next.js приложение:
- Установка через create-next-app
- Структура App Router
- Базовая навигация

### Задание 2: Многостраничный сайт (Уровень: Средний)
Сайт с маршрутами:
- Главная страница
- О компании
- Услуги
- Контакты
- Layout с навигацией

### Задание 3: Блог с SSG (Уровень: Продвинутый)
Блог со статической генерацией:
- Список постов из markdown
- Динамические маршруты /[slug]
- generateStaticParams

### Задание 4: E-commerce с ISR (Уровень: Продвинутый)
Магазин с товарами:
- Каталог товаров (SSG)
- Карточка товара (ISR)
- Корзина (Client Component)

### Задание 5: Dashboard с SSR (Уровень: Эксперт)
Личный кабинет:
- Защита маршрутов
- Данные пользователя (SSR)
- Real-time обновления

---

## ❓ Проверочные вопросы

1. Чем SSR отличается от SSG?
2. Когда использовать Client Components?
3. Что такое Streaming в Next.js?
4. Как работает ISR?
5. Что такое Server Actions?
6. Как оптимизировать изображения?
7. Для чего нужен metadata API?
8. Как реализовать SEO в Next.js?
9. Что такое Middleware?
10. Как деплоить Next.js приложение?

---

**Следующий модуль:** `02-advanced-routing` - Продвинутая маршрутизация
EOF

# ===== NEXT.JS 02-advanced-routing =====
cat > /workspace/learning-path/06-nextjs/02-advanced-routing/README.md << 'EOF'
# Модуль 2: Продвинутая маршрутизация

## 📚 Теория

### 2.1 File-based Routing

```
app/
├── layout.tsx
├── page.tsx
├── about/
│   └── page.tsx
├── blog/
│   ├── page.tsx
│   └── [slug]/
│       └── page.tsx
└── dashboard/
    ├── layout.tsx
    └── settings/
        └── page.tsx
```

### 2.2 Dynamic Routes

```jsx
// app/products/[id]/page.tsx
export default async function ProductPage({ params }) {
  const { id } = await params;
  return <div>Product {id}</div>;
}

// app/blog/[category]/[slug]/page.tsx
export default async function PostPage({ params }) {
  const { category, slug } = await params;
  // ...
}
```

### 2.3 Loading и Error States

```jsx
// loading.tsx
export default function Loading() {
  return <div>Загрузка...</div>;
}

// error.tsx
'use client';
export default function Error({ error, reset }) {
  return (
    <div>
      <h2>Ошибка: {error.message}</h2>
      <button onClick={reset}>Попробовать снова</button>
    </div>
  );
}
```

### 2.4 Not Found

```jsx
import { notFound } from 'next/navigation';

export default async function Page({ params }) {
  const data = await fetchData(params.id);
  if (!data) notFound();
  return <div>{data.title}</div>;
}
```

---

## 🎯 Практические задания

### Задание 1: Вложенные маршруты (Уровень: Средний)
Создайте структуру:
- Dashboard layout
- Вложенные разделы
- Активная навигация

### Задание 2: Parallel Routes (Уровень: Продвинутый)
Реализуйте:
- @modal для модальных окон
- Несколько views одновременно

### Задание 3: Intercepting Routes (Уровень: Продвинутый)
Сделайте:
- Модальное окно на том же URL
- Сохранение контекста

### Задание 4: Route Groups (Уровень: Средний)
Организуйте:
- (marketing) группа
- (shop) группа
- Разные layouts

### Задание 5: Полноценный портал (Уровень: Эксперт)
Комплексное решение:
- Все виды маршрутов
- Переходы между состояниями
- Оптимизация

---

**Следующий модуль:** `03-data-fetching` - Продвинутая работа с данными
EOF

# ===== NEXT.JS 03-data-fetching =====
cat > /workspace/learning-path/06-nextjs/03-data-fetching/README.md << 'EOF'
# Модуль 3: Продвинутая работа с данными

## 📚 Теория

### 3.1 Caching Strategies

```jsx
// Force Cache (SSG)
fetch('url', { cache: 'force-cache' });

// No Store (SSR)
fetch('url', { cache: 'no-store' });

// Revalidate (ISR)
fetch('url', { next: { revalidate: 3600 } });

// Tag-based revalidation
fetch('url', { next: { tags: ['products'] } });
// revalidateTag('products');
```

### 3.2 Server Actions

```jsx
// actions.ts
'use server';
export async function createTodo(formData: FormData) {
  const title = formData.get('title');
  await db.todo.create({ data: { title } });
  revalidatePath('/todos');
}

// Component
<form action={createTodo}>
  <input name="title" />
  <button type="submit">Add</button>
</form>
```

### 3.3 Optimistic Updates

```jsx
'use client';
import { useOptimistic } from 'react';

function TodoList({ todos }) {
  const [optimisticTodos, addOptimistic] = useOptimistic(
    todos,
    (state, newTodo) => [...state, newTodo]
  );
  
  // ...
}
```

---

## 🎯 Практические задания

### Задание 1: CRUD с Server Actions (Уровень: Средний)
Реализуйте:
- Создание записи
- Редактирование
- Удаление
- Валидация Zod

### Задание 2: Поиск с debounce (Уровень: Средний)
Сделайте:
- Поисковая строка
- Debounce запросов
- Pagination

### Задание 3: Форма с прогрессом (Уровень: Продвинутый)
Multi-step форма:
- Сохранение черновика
- Валидация каждого шага
- Отправка всех данных

### Задание 4: Real-time dashboard (Уровень: Продвинутый)
Дашборд с:
- SSE или WebSocket
- Оптимистичные обновления
- Background sync

### Задание 5: E-commerce платформа (Уровень: Эксперт)
Полноценный магазин:
- Каталог с фильтрами
- Корзина
- Оформление заказа
- Интеграция платежей

---

**Следующий модуль:** `04-authentication` - Аутентификация и авторизация
EOF

echo "Next.js модули созданы!"
