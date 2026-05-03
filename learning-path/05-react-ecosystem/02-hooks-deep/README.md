# Модуль 2: Углубленное изучение хуков React

## 📚 Теория

### 2.1 useState - продвинутые паттерны

```jsx
// Функциональное обновление состояния
function Counter() {
  const [count, setCount] = useState(0);
  
  // Правильно: используем предыдущее состояние
  const increment = () => {
    setCount(prev => prev + 1);
  };
  
  // Неправильно: может привести к гонкам
  const badIncrement = () => {
    setCount(count + 1);
  };
}

// Инициализация ленивым способом
function ExpensiveComponent() {
  const [data, setData] = useState(() => {
    return computeExpensiveValue();
  });
}
```

### 2.2 useEffect - глубокая работа

```jsx
import { useEffect, useState } from 'react';

function DataFetcher({ url }) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    let isMounted = true;
    
    const fetchData = async () => {
      try {
        const response = await fetch(url);
        const result = await response.json();
        
        if (isMounted) {
          setData(result);
          setLoading(false);
        }
      } catch (err) {
        if (isMounted) {
          setError(err);
          setLoading(false);
        }
      }
    };

    fetchData();

    // Cleanup функция
    return () => {
      isMounted = false;
    };
  }, [url]);

  if (loading) return <div>Загрузка...</div>;
  if (error) return <div>Ошибка: {error.message}</div>;
  return <div>{JSON.stringify(data)}</div>;
}
```

### 2.3 useContext - глобальное состояние

```jsx
// Создание контекста
const ThemeContext = createContext('light');

// Провайдер
function App() {
  const [theme, setTheme] = useState('dark');
  
  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      <Toolbar />
    </ThemeContext.Provider>
  );
}

// Потребление контекста
function ThemedButton() {
  const { theme, setTheme } = useContext(ThemeContext);
  
  return (
    <button 
      style={{ background: theme === 'dark' ? '#333' : '#fff' }}
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
    >
      Текущая тема: {theme}
    </button>
  );
}
```

### 2.4 useReducer - сложная логика состояния

```jsx
function reducer(state, action) {
  switch (action.type) {
    case 'ADD':
      return { ...state, items: [...state.items, action.payload] };
    case 'REMOVE':
      return { 
        ...state, 
        items: state.items.filter(item => item.id !== action.payload) 
      };
    case 'UPDATE':
      return {
        ...state,
        items: state.items.map(item =>
          item.id === action.payload.id ? { ...item, ...action.payload } : item
        )
      };
    default:
      throw new Error(`Неизвестный тип: ${action.type}`);
  }
}

function TodoApp() {
  const [state, dispatch] = useReducer(reducer, { items: [], filter: 'all' });

  const addTodo = (text) => {
    dispatch({ 
      type: 'ADD', 
      payload: { id: Date.now(), text, completed: false } 
    });
  };

  return (
    <div>
      {state.items.map(todo => (
        <TodoItem key={todo.id} todo={todo} dispatch={dispatch} />
      ))}
    </div>
  );
}
```

### 2.5 useMemo и useCallback - оптимизация

```jsx
import { useMemo, useCallback } from 'react';

function ExpensiveList({ items, filter }) {
  // Мемоизация вычислений
  const filteredItems = useMemo(() => {
    console.log('Фильтрация...');
    return items.filter(item => item.category === filter);
  }, [items, filter]);

  // Мемоизация функции
  const handleItemClick = useCallback((id) => {
    console.log('Клик по элементу:', id);
  }, []);

  return (
    <ul>
      {filteredItems.map(item => (
        <li key={item.id} onClick={() => handleItemClick(item.id)}>
          {item.name}
        </li>
      ))}
    </ul>
  );
}
```

### 2.6 Custom Hooks - создание своих хуков

```jsx
// useLocalStorage
function useLocalStorage(key, initialValue) {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.error(error);
      return initialValue;
    }
  });

  const setValue = (value) => {
    try {
      setStoredValue(value);
      window.localStorage.setItem(key, JSON.stringify(value));
    } catch (error) {
      console.error(error);
    }
  };

  return [storedValue, setValue];
}

// useFetch
function useFetch(url) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(url);
        const result = await response.json();
        setData(result);
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [url]);

  return { data, loading, error };
}

// useDebounce
function useDebounce(value, delay) {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}
```

### 2.7 useRef и DOM манипуляции

```jsx
function FocusInput() {
  const inputRef = useRef(null);

  useEffect(() => {
    inputRef.current.focus();
  }, []);

  return <input ref={inputRef} placeholder="Автофокус" />;
}

// Хранение мутабельных значений
function Timer() {
  const intervalRef = useRef(null);
  const [count, setCount] = useState(0);

  const startTimer = () => {
    intervalRef.current = setInterval(() => {
      setCount(prev => prev + 1);
    }, 1000);
  };

  const stopTimer = () => {
    clearInterval(intervalRef.current);
  };

  return (
    <div>
      <p>Секунд: {count}</p>
      <button onClick={startTimer}>Старт</button>
      <button onClick={stopTimer}>Стоп</button>
    </div>
  );
}
```

---

## 🎯 Практические задания

### Задание 1: Кастомный хук useToggle (Уровень: Начинающий)
Создайте хук `useToggle`, который:
- Принимает начальное значение (boolean)
- Возвращает текущее значение и функцию переключения
- Имеет функции `setTrue`, `setFalse`, `toggle`

**Требования:**
- Типизация на TypeScript
- Пример использования в компоненте

### Задание 2: Форма с валидацией на useReducer (Уровень: Средний)
Создайте компонент формы регистрации с использованием `useReducer`:
- Поля: email, password, confirmPassword
- Действия: CHANGE_FIELD, VALIDATE, SUBMIT
- Валидация каждого поля
- Отображение ошибок

**Требования:**
- Чистая архитектура редюсера
- Разделение логики валидации
- Контролируемые компоненты

### Задание 3: Поиск с debounce (Уровень: Средний)
Создайте компонент поиска товаров:
- Поле ввода с debounce (500мс)
- Запрос к API при изменении поискового запроса
- Отображение результатов
- Обработка состояний загрузки и ошибки

**Требования:**
- Используйте кастомный хук `useDebounce`
- Отмена предыдущего запроса при новом вводе
- Мемоизация результатов

### Задание 4: Корзина товаров с useContext + useReducer (Уровень: Продвинутый)
Создайте систему управления корзиной:
- Глобальный контекст CartContext
- Операции: добавить, удалить, изменить количество, очистить
- Подсчет общей суммы
- Сохранение в localStorage

**Требования:**
- Кастомный хук `useCart`
- Персистентность данных
- Оптимизация рендеров

### Задание 5: Бесконечный скролл (Уровень: Продвинутый)
Создайте компонент с бесконечной прокруткой:
- Загрузка данных порциями по 20 элементов
- Использование Intersection Observer
- Индикатор загрузки
- Обработка ошибок сети

**Требования:**
- Кастомный хук `useInfiniteScroll`
- Виртуализация списка (опционально)
- Пагинация на сервере (эмуляция)

---

## ❓ Проверочные вопросы

1. **В чем разница между `useMemo` и `useCallback`?**
2. **Когда `useEffect` выполняется чаще одного раза?**
3. **Почему важно возвращать cleanup функцию из `useEffect`?**
4. **Как избежать лишних рендеров при использовании контекста?**
5. **Что такое "stale closure" и как с ним бороться?**
6. **Когда стоит использовать `useReducer` вместо `useState`?**
7. **Как работает ленивая инициализация в `useState`?**
8. **Почему нельзя вызывать хуки внутри условий?**
9. **Как создать свой кастомный хук? Приведите пример.**
10. **В чем разница между `ref` и `state`?**

---

## ✅ Критерии выполнения

- [ ] Все 5 заданий выполнены
- [ ] Кастомные хуки переиспользуемы
- [ ] Правильная очистка эффектов
- [ ] Оптимизация производительности
- [ ] Код типизирован (TypeScript)
- [ ] Ответы на вопросы даны правильно

---

## 📖 Дополнительные ресурсы

- [React Hooks Documentation](https://react.dev/reference/react)
- [Use Hooks](https://usehooks.com/)
- [React Hooks Best Practices](https://kentcdodds.com/blog/react-hooks-best-practices)

---

**Следующий модуль:** `03-state-management` - Управление состоянием приложения
