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
