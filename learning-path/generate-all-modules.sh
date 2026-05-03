#!/bin/bash

# ===== REACT 05-performance =====
cat > /workspace/learning-path/05-react-ecosystem/05-performance/README.md << 'EOF'
# Модуль 5: Оптимизация производительности React

## 📚 Теория

### 5.1 Мемоизация компонентов

```jsx
import { memo, useMemo, useCallback } from 'react';

// Мемоизация компонента
const ExpensiveComponent = memo(({ data, onAction }) => {
  return <div>{/* рендеринг */}</div>;
});

// Мемоизация вычислений
const filtered = useMemo(() => {
  return items.filter(item => item.active);
}, [items]);

// Мемоизация функций
const handleClick = useCallback(() => {
  doSomething(id);
}, [id]);
```

### 5.2 Code Splitting

```jsx
import { lazy, Suspense } from 'react';

const Dashboard = lazy(() => import('./Dashboard'));

function App() {
  return (
    <Suspense fallback={<Spinner />}>
      <Dashboard />
    </Suspense>
  );
}
```

### 5.3 Виртуализация списков

```jsx
import { FixedSizeList } from 'react-window';

function VirtualizedList({ items }) {
  return (
    <FixedSizeList
      height={600}
      itemCount={items.length}
      itemSize={50}
    >
      {({ index, style }) => (
        <div style={style}>{items[index]}</div>
      )}
    </FixedSizeList>
  );
}
```

### 5.4 React DevTools Profiler

Используйте Profiler для выявления узких мест:
- Длительные рендеры
- Лишние ре-рендеры
- Оптимизация зависимостей useEffect

---

## 🎯 Практические задания

### Задание 1: Оптимизация списка (Уровень: Средний)
Оптимизируйте рендеринг большого списка (1000+ элементов):
- Используйте виртуализацию
- Добавьте мемоизацию
- Сравните FPS до/после

### Задание 2: Lazy Loading изображений (Уровень: Средний)
Реализуйте ленивую загрузку изображений:
- Intersection Observer
- Placeholder при загрузке
- Прогрессивное улучшение

### Задание 3: Оптимизация формы (Уровень: Продвинутый)
Оптимизируйте форму с множеством полей:
- Разделение на подкомпоненты
- Мемоизация обработчиков
- Debounce валидации

### Задание 4: Анализ с React Profiler (Уровень: Продвинутый)
Проанализируйте существующее приложение:
- Найдите компоненты с лишними рендерами
- Исправьте проблемы
- Задокументируйте улучшения

### Задание 5: Полная оптимизация приложения (Уровень: Эксперт)
Проведите полную оптимизацию:
- Code splitting по маршрутам
- Lazy loading компонентов
- Мемоизация где необходимо
- Webpack bundle analysis

---

## ❓ Проверочные вопросы

1. Когда использовать `memo`?
2. В чем разница между `useMemo` и `useCallback`?
3. Что такое virtual DOM и как он работает?
4. Как работает code splitting?
5. Что такое windowing/virtualization?
6. Как измерить производительность React приложения?
7. Что такое Fiber архитектура?
8. Когда оптимизация вредна?
9. Как избежать лишних рендеров?
10. Что такое Concurrent Mode?

---

**Следующий модуль:** `06-testing-react` - Тестирование React приложений
EOF

# ===== REACT 06-testing-react =====
cat > /workspace/learning-path/05-react-ecosystem/06-testing-react/README.md << 'EOF'
# Модуль 6: Тестирование React приложений

## 📚 Теория

### 6.1 Jest + React Testing Library

```jsx
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

test('отображает приветствие', () => {
  render(<Welcome name="Анна" />);
  expect(screen.getByText('Привет, Анна!')).toBeInTheDocument();
});

test('обработка клика', async () => {
  const handleClick = jest.fn();
  render(<button onClick={handleClick}>Click</button>);
  
  await userEvent.click(screen.getByRole('button'));
  expect(handleClick).toHaveBeenCalledTimes(1);
});
```

### 6.2 Тестирование хуков

```jsx
import { renderHook, act } from '@testing-library/react-hooks';
import { useCounter } from './useCounter';

test('увеличивает счетчик', () => {
  const { result } = renderHook(() => useCounter());
  
  act(() => {
    result.current.increment();
  });
  
  expect(result.current.count).toBe(1);
});
```

### 6.3 Моки и снапшоты

```jsx
// Мок API
jest.mock('../api');
api.fetchUser.mockResolvedValue({ name: 'John' });

// Снапшот тесты
test('соответствует снапшоту', () => {
  const { container } = render(<Component />);
  expect(container).toMatchSnapshot();
});
```

---

## 🎯 Практические задания

### Задание 1: Тесты компонентов (Уровень: Начинающий)
Напишите тесты для:
- Отображение props
- Обработка событий
- Условный рендеринг

### Задание 2: Тесты форм (Уровень: Средний)
Протестируйте форму:
- Валидация полей
- Отправка данных
- Ошибки сервера

### Задание 3: Тесты хуков (Уровень: Средний)
Протестируйте кастомные хуки:
- useToggle
- useLocalStorage
- useFetch

### Задание 4: Integration тесты (Уровень: Продвинутый)
Напишите интеграционные тесты:
- Взаимодействие компонентов
- Роутинг
- Глобальное состояние

### Задание 5: E2E тесты с Cypress (Уровень: Эксперт)
Создайте E2E тесты:
- Критические пользовательские сценарии
- Тесты авторизации
- Тесты покупки товара

---

## ❓ Проверочные вопросы

1. Чем unit тесты отличаются от integration?
2. Что такое TDD?
3. Зачем нужны моки?
4. Когда использовать снапшот тесты?
5. Как тестировать асинхронный код?
6. Что такое arrange-act-assert паттерн?
7. Как тестировать Redux store?
8. Что такое test coverage?
9. Как отлаживать тесты?
10. Какие антипаттерны в тестировании знаете?

---

**Следующий модуль:** `07-project` - Финальный проект React
EOF

# ===== REACT 07-project =====
cat > /workspace/learning-path/05-react-ecosystem/07-project/README.md << 'EOF'
# Модуль 7: Финальный проект React

## 🎯 Задача: Интернет-магазин электроники

### Требования к проекту

**Функциональность:**
- Каталог товаров с фильтрами и поиском
- Карточка товара с галереей
- Корзина с управлением количеством
- Оформление заказа (multi-step form)
- Личный кабинет пользователя
- Избранное
- Отзывы о товарах
- Админ панель (CRUD товаров)

**Технические требования:**
- TypeScript
- Redux Toolkit + RTK Query
- React Router v6
- React Hook Form + Zod
- Axios interceptors
- Lazy loading
- Оптимизация производительности
- Unit и Integration тесты
- ESLint + Prettier
- Husky pre-commit hooks

**Дополнительно:**
- PWA (offline режим)
- Темная тема
- Мультиязычность (i18n)
- Docker контейнеризация

### Этапы реализации

1. **Неделя 1:** Настройка проекта, архитектура, UI Kit
2. **Неделя 2:** Каталог товаров, фильтры, поиск
3. **Неделя 3:** Корзина, оформление заказа
4. **Неделя 4:** Авторизация, личный кабинет
5. **Неделя 5:** Админ панель, тесты
6. **Неделя 6:** Оптимизация, деплой

### Критерии оценки

- [ ] Код типизирован TypeScript
- [ ] Архитектура соответствует best practices
- [ ] Все фичи реализованы
- [ ] Покрытие тестами > 70%
- [ ] Производительность (Lighthouse > 90)
- [ ] Адаптивный дизайн
- [ ] Деплой на Vercel/Netlify

---

**Поздравляем!** После завершения этого модуля вы готовы перейти к изучению Next.js!
EOF

echo "React модули созданы!"
