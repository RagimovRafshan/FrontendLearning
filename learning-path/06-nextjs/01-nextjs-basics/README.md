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
