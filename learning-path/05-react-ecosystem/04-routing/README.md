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
