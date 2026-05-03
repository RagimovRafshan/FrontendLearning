# Модуль 3: Продвинутая работа с данными в Next.js

## 📚 Теория

### 3.1 Caching Strategies

```jsx
// Force Cache (SSG) - кэшируется навсегда
fetch('url', { cache: 'force-cache' });

// No Store (SSR) - каждый запрос на сервер
fetch('url', { cache: 'no-store' });

// Revalidate (ISR) - обновление по таймеру
fetch('url', { next: { revalidate: 3600 } }); // 1 час

// Tag-based revalidation
fetch('url', { next: { tags: ['products'] } });
// В Server Action: revalidateTag('products');
```

### 3.2 Server Actions

```jsx
// app/actions/todos.ts
'use server';
import { revalidatePath } from 'next/cache';

export async function createTodo(formData: FormData) {
  const title = formData.get('title') as string;
  await db.todo.create({ data: { title } });
  revalidatePath('/todos');
}

export async function deleteTodo(id: number) {
  await db.todo.delete({ where: { id } });
  revalidatePath('/todos');
}
```

### 3.3 Optimistic Updates

```jsx
'use client';
import { useOptimistic } from 'react';

function TodoList({ todos }: { todos: Todo[] }) {
  const [optimisticTodos, addOptimistic] = useOptimistic(
    todos,
    (state, newTodo: Todo) => [...state, newTodo]
  );
  
  return (
    <ul>
      {optimisticTodos.map(todo => <li key={todo.id}>{todo.title}</li>)}
    </ul>
  );
}
```

---

## 🎯 Практические задания

### Задание 1: CRUD с Server Actions (Уровень: Средний)
Создайте приложение заметок:
- Создание, чтение, обновление, удаление
- Валидация с Zod
- Реалидейт путей

### Задание 2: Поиск с фильтрами (Уровень: Средний)
Каталог товаров:
- Поисковая строка с debounce
- Фильтры по категориям
- Пагинация
- URL синхронизация

### Задание 3: Форма с черновиками (Уровень: Продвинутый)
Multi-step форма заказа:
- Сохранение черновика в БД
- Валидация каждого шага
- Восстановление при возврате

### Задание 4: Real-time дашборд (Уровень: Продвинутый)
Дашборд аналитики:
- SSE для обновлений
- Оптимистичные обновления
- Фоновая синхронизация

### Задание 5: E-commerce платформа (Уровень: Эксперт)
Полноценный магазин:
- Каталог с ISR
- Корзина (Client)
- Оформление заказа (Server Actions)
- Интеграция платежей (Stripe)

---

## ❓ Проверочные вопросы

1. Когда использовать SSR vs SSG vs ISR?
2. Как работает тег-based ревалидация?
3. Что такое Server Actions?
4. Как передать formData в Server Action?
5. Для чего нужен useOptimistic?
6. Как отменить действие при ошибке?
7. Что такое progressive enhancement?
8. Как обрабатывать ошибки в Server Actions?
9. Можно ли вызывать Server Actions из useEffect?
10. Как протестировать Server Actions?

---

**Следующий модуль:** `04-authentication` - Аутентификация и авторизация
