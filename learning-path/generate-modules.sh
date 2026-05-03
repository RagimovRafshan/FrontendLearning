#!/bin/bash

# Генерация модуля React 03-state-management
cat > /workspace/learning-path/05-react-ecosystem/03-state-management/README.md << 'EOF'
# Модуль 3: Управление состоянием приложения

## 📚 Теория

### 3.1 Проблемы управления состоянием
- Prop drilling (передача пропсов через множество компонентов)
- Глобальное состояние vs локальное
- Синхронизация состояния между компонентами

### 3.2 Redux Toolkit - современный Redux

```jsx
import { createSlice, configureStore } from '@reduxjs/toolkit';

// Создание слайса
const counterSlice = createSlice({
  name: 'counter',
  initialState: { value: 0 },
  reducers: {
    increment: (state) => { state.value += 1; },
    decrement: (state) => { state.value -= 1; },
    incrementByAmount: (state, action) => { state.value += action.payload; }
  }
});

// Настройка store
const store = configureStore({
  reducer: { counter: counterSlice.reducer }
});

// Использование в компоненте
import { useSelector, useDispatch } from 'react-redux';

function Counter() {
  const count = useSelector((state) => state.counter.value);
  const dispatch = useDispatch();
  
  return (
    <div>
      <p>{count}</p>
      <button onClick={() => dispatch(increment())}>+</button>
    </div>
  );
}
```

### 3.3 Zustand - легкая альтернатива

```jsx
import { create } from 'zustand';

// Создание store
const useStore = create((set) => ({
  bears: 0,
  increasePopulation: () => set((state) => ({ bears: state.bears + 1 })),
  removeAllBears: () => set({ bears: 0 })
}));

// Использование
function BearDisplay() {
  const bears = useStore((state) => state.bears);
  return <h1>{bears} bears around here</h1>;
}
```

### 3.4 React Query (TanStack Query) - серверное состояние

```jsx
import { useQuery, useMutation, QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <UserProfile />
    </QueryClientProvider>
  );
}

function UserProfile() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['user', userId],
    queryFn: () => fetch(`/api/users/${userId}`).then(res => res.json())
  });
  
  const mutation = useMutation({
    mutationFn: (newData) => fetch('/api/users', { method: 'POST', body: newData }),
    onSuccess: () => {
      queryClient.invalidateQueries(['users']);
    }
  });
}
```

---

## 🎯 Практические задания

### Задание 1: Todo на Redux Toolkit (Уровень: Средний)
Создайте приложение списка задач с Redux Toolkit:
- Слайс с CRUD операциями
- Фильтрация (все/активные/завершенные)
- Подсчет статистики

### Задание 2: Корзина на Zustand (Уровень: Средний)
Реализуйте корзину товаров:
- Добавление/удаление товаров
- Изменение количества
- Подсчет общей суммы
- Персистентность в localStorage

### Задание 3: Блог с React Query (Уровень: Продвинутый)
Создайте блог с загрузкой постов:
- Пагинация
- Кэширование
- Оптимистичные обновления
- Background refetch

### Задание 4: E-commerce приложение (Уровень: Продвинутый)
Полноценное приложение с:
- Redux Toolkit для глобального состояния
- React Query для серверных данных
- Авторизация
- Избранное
- История заказов

### Задание 5: Сравнительный анализ (Уровень: Эксперт)
Реализуйте одно и то же приложение на:
- Redux Toolkit
- Zustand
- Context API
Сравните объем кода, производительность, DX.

---

## ❓ Проверочные вопросы

1. Когда использовать глобальное состояние, а когда локальное?
2. В чем преимущества Redux Toolkit перед классическим Redux?
3. Как работает селектор в Redux?
4. Что такое нормализация состояния?
5. Чем отличается серверное состояние от клиентского?
6. Как работает кэширование в React Query?
7. Когда стоит выбрать Zustand вместо Redux?
8. Что такое middleware в Redux?
9. Как оптимизировать рендеры при использовании Redux?
10. Что такое optimistic updates?

---

## ✅ Критерии выполнения

- [ ] Все задания выполнены
- [ ] Понимание различий между решениями
- [ ] Правильная архитектура store
- [ ] Обработка ошибок и loading состояний
- [ ] Типизация TypeScript

---

**Следующий модуль:** `04-routing` - Маршрутизация в React
EOF

# Генерация модуля React 04-routing
cat > /workspace/learning-path/05-react-ecosystem/04-routing/README.md << 'EOF'
# Модуль 4: Маршрутизация в React

## 📚 Теория

### 4.1 React Router v6

```jsx
import { BrowserRouter, Routes, Route, Link, useParams, useNavigate } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <nav>
        <Link to="/">Главная</Link>
        <Link to="/about">О нас</Link>
        <Link to="/users">Пользователи</Link>
      </nav>
      
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/users/:id" element={<UserDetail />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}

// Параметры маршрута
function UserDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  
  return (
    <div>
      <h1>Пользователь {id}</h1>
      <button onClick={() => navigate(-1)}>Назад</button>
    </div>
  );
}
```

### 4.2 Защищенные маршруты

```jsx
function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  
  return children;
}

// Использование
<Route 
  path="/dashboard" 
  element={
    <ProtectedRoute>
      <Dashboard />
    </ProtectedRoute>
  } 
/>
```

### 4.3 Ленивая загрузка маршрутов

```jsx
import { Suspense, lazy } from 'react';

const About = lazy(() => import('./pages/About'));
const Dashboard = lazy(() => import('./pages/Dashboard'));

function App() {
  return (
    <Suspense fallback={<LoadingSpinner />}>
      <Routes>
        <Route path="/about" element={<About />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </Suspense>
  );
}
```

---

## 🎯 Практические задания

### Задание 1: Многостраничный сайт (Уровень: Начинающий)
Создайте сайт с маршрутами:
- Главная
- О компании
- Контакты
- Страница 404

### Задание 2: Блог с динамическими маршрутами (Уровень: Средний)
- Список постов `/posts`
- Детальная страница `/posts/:id`
- Фильтрация по категориям `/posts/category/:name`

### Задание 3: Личный кабинет с защитой (Уровень: Продвинутый)
- Регистрация/Вход
- Защищенные маршруты
- Восстановление пароля
- Редактирование профиля

### Задание 4: E-commerce с вложенными маршрутами (Уровень: Продвинутый)
- Каталог товаров
- Карточка товара
- Корзина
- Оформление заказа (multi-step form)
- История заказов

### Задание 5: Dashboard приложение (Уровень: Эксперт)
- Боковая навигация
- Вложенные маршруты
- Lazy loading
- Сохранение состояния при переходе
- Bread crumbs

---

## ❓ Проверочные вопросы

1. Чем отличается `<Link>` от `<a>` тега?
2. Как передать состояние между маршрутами?
3. Что такое вложенные маршруты?
4. Как реализовать защиту маршрутов?
5. Для чего нужен `useNavigate`?
6. Как получить query параметры?
7. Что такое code splitting в контексте роутинга?
8. Как обработать несуществующий маршрут?
9. Можно ли иметь несколько `<BrowserRouter>`?
10. Как программно перейти на другой маршрут?

---

**Следующий модуль:** `05-performance` - Оптимизация производительности
EOF

echo "Модули React созданы!"
