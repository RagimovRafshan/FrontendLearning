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
