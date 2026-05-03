# Модуль 2: Продвинутая маршрутизация в Next.js

## 📚 Теория

### 2.1 File-based Routing (App Router)

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
// loading.tsx - автоматически показывается при навигации
export default function Loading() {
  return <div className="spinner">Загрузка...</div>;
}

// error.tsx - обработка ошибок
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
Создайте структуру Dashboard:
- Общий layout для раздела
- Вложенные страницы настроек
- Активная навигация

### Задание 2: Parallel Routes (Уровень: Продвинутый)
Реализуйте:
- @modal для модальных окон
- Несколько views одновременно на одной странице

### Задание 3: Intercepting Routes (Уровень: Продвинутый)
Сделайте:
- Модальное окно открывается на том же URL
- При refresh страница показывает полноценный маршрут

### Задание 4: Route Groups (Уровень: Средний)
Организуйте проект:
- (marketing) группа с одним layout
- (shop) группа с другим layout
- Разделение логики без влияния на URL

### Задание 5: Полноценный портал (Уровень: Эксперт)
Комплексное решение:
- Все виды маршрутов
- Переходы между состояниями
- Оптимизация загрузки

---

## ❓ Проверочные вопросы

1. Как создать динамический маршрут?
2. Что делает файл loading.tsx?
3. Как обработать ошибку на уровне маршрута?
4. Что такое parallel routes?
5. Для чего нужны route groups?
6. Как работает intercepting routes?
7. Что такое catch-all routes?
8. Как передать данные между маршрутами?
9. Когда использовать notFound()?
10. Как сделать редирект?

---

**Следующий модуль:** `03-data-fetching` - Продвинутая работа с данными
