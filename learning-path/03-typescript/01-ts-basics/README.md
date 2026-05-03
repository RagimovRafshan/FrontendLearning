# Модуль 1: Основы TypeScript

## 📚 Теория

### 1.1 Что такое TypeScript?

**TypeScript** — это строго типизированный язык программирования, который является надмножеством JavaScript. Он добавляет статическую типизацию и другие возможности, которые помогают писать более надежный код.

**Преимущества TypeScript:**
- Статическая типизация (ошибки ловятся на этапе компиляции)
- Улучшенная автодополнение в IDE
- Лучшая документация кода через типы
- Рефакторинг становится проще и безопаснее
- Поддержка современных стандартов ES6+

```typescript
// JavaScript - ошибка будет только при выполнении
function add(a, b) {
  return a + b;
}
add(5, '10'); // "510" - неожиданное поведение!

// TypeScript - ошибка будет показана сразу
function add(a: number, b: number): number {
  return a + b;
}
add(5, '10'); // ❌ Ошибка: Argument of type 'string' is not assignable to parameter of type 'number'
```

### 1.2 Базовые типы данных

#### Примитивные типы

```typescript
// String - строки
let name: string = 'Alice';
let greeting: string = `Hello, ${name}`; // Шаблонные строки

// Number - числа (целые и дробные)
let age: number = 25;
let price: number = 19.99;
let hex: number = 0xf00d; // Шестнадцатеричная система
let binary: number = 0b1010; // Двоичная система
let octal: number = 0o744; // Восьмеричная система

// Boolean - логические значения
let isActive: boolean = true;
let isDisabled: boolean = false;

// Null и Undefined - специальные типы
let nullable: null = null;
let undefinedValue: undefined = undefined;

// Void - отсутствие значения (обычно для функций)
function logMessage(message: string): void {
  console.log(message);
  // Функция ничего не возвращает
}

// Never - значение, которое никогда не будет получено
function throwError(message: string): never {
  throw new Error(message);
}

function infiniteLoop(): never {
  while (true) {}
}
```

#### Специальные типы

```typescript
// Any - отключение типизации (использовать осторожно!)
let anything: any = 'string';
anything = 123;
anything = true; // ✅ Никаких ошибок, но теряем преимущества TS

// Unknown - безопасная версия any
let unknownValue: unknown = 'hello';
// unknownValue.length; // ❌ Ошибка: нужно проверить тип

if (typeof unknownValue === 'string') {
  console.log(unknownValue.length); // ✅ Теперь безопасно
}

// Object - любой объект (не примитив)
let obj: object = { name: 'Alice' };
// let primitive: object = 'string'; // ❌ Ошибка
```

### 1.3 Типизация массивов и кортежей

```typescript
// Массивы - несколько способов объявления
let numbers: number[] = [1, 2, 3, 4, 5];
let names: Array<string> = ['Alice', 'Bob', 'Charlie']; // Generic синтаксис

// Массив объектов
interface User {
  id: number;
  name: string;
}
let users: User[] = [
  { id: 1, name: 'Alice' },
  { id: 2, name: 'Bob' }
];

// Кортежи (Tuples) - массивы фиксированной длины с известными типами
let tuple: [string, number, boolean] = ['Alice', 25, true];
let first: string = tuple[0]; // 'Alice'
let second: number = tuple[1]; // 25

// Именованные кортежи (TS 4.0+)
type Person = [name: string, age: number, isActive?: boolean];
const person: Person = ['Alice', 25, true];

// Опциональные элементы кортежа
let optionalTuple: [string, number?] = ['hello']; // Второй элемент необязателен
```

### 1.4 Type Aliases (Псевдонимы типов)

```typescript
// Создание псевдонима для типа
type UserId = string | number;
type Point = {
  x: number;
  y: number;
};

// Использование
let userId: UserId = 'abc123';
userId = 456; // ✅ Тоже работает

let point: Point = { x: 10, y: 20 };

// Композиция типов
type Color = 'red' | 'green' | 'blue';
type Circle = {
  kind: 'circle';
  radius: number;
  color: Color;
};

type Rectangle = {
  kind: 'rectangle';
  width: number;
  height: number;
  color: Color;
};

type Shape = Circle | Rectangle;

function getArea(shape: Shape): number {
  if (shape.kind === 'circle') {
    return Math.PI * shape.radius ** 2;
  } else {
    return shape.width * shape.height;
  }
}
```

### 1.5 Union и Intersection типы

```typescript
// Union Types (ИЛИ) - значение может быть одного из нескольких типов
type Id = string | number;
let userId: Id = 'abc123';
userId = 456; // ✅

// Сужение типа (Type Narrowing)
function printId(id: Id) {
  if (typeof id === 'string') {
    console.log(id.toUpperCase()); // id здесь имеет тип string
  } else {
    console.log(id.toFixed(2)); // id здесь имеет тип number
  }
}

// Intersection Types (И) - комбинация всех свойств
type A = { a: number; b: string };
type B = { b: string; c: boolean };
type C = A & B; // { a: number; b: string; c: boolean }

let obj: C = {
  a: 1,
  b: 'hello',
  c: true
};

// Практический пример
type StyleProps = { color: string; fontSize: number };
type LayoutProps = { margin: number; padding: number };
type ComponentProps = StyleProps & LayoutProps;

const props: ComponentProps = {
  color: 'blue',
  fontSize: 16,
  margin: 10,
  padding: 5
};
```

### 1.6 Literal Types (Литеральные типы)

```typescript
// Строковые литералы
type Direction = 'up' | 'down' | 'left' | 'right';
function move(direction: Direction) {
  console.log(`Moving ${direction}`);
}
move('up'); // ✅
move('diagonal'); // ❌ Ошибка

// Числовые литералы
type DiceRoll = 1 | 2 | 3 | 4 | 5 | 6;
function rollDice(): DiceRoll {
  return Math.floor(Math.random() * 6) + 1 as DiceRoll;
}

// Булевы литералы
type YesNo = true | false;

// Комбинирование с union
type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';
type APIResponse<T> = {
  method: HttpMethod;
  status: number;
  data: T;
};

const response: APIResponse<string[]> = {
  method: 'GET',
  status: 200,
  data: ['item1', 'item2']
};
```

## 🎯 Практические задания

### Задание 1: Базовая типизация переменных
Создайте файл `src/variables.ts` и объявите переменные со следующими типами:
1. `userName` (string) - ваше имя
2. `userAge` (number) - ваш возраст
3. `isStudent` (boolean) - учитесь ли вы
4. `hobbies` (массив строк) - список хобби
5. `address` (объект с полями: city, street, building)
6. `nullableValue` (может быть string или null)

**Критерии выполнения:**
- Все переменные имеют явную типизацию
- Код компилируется без ошибок
- Использованы разные способы объявления типов

### Задание 2: Функции с типизацией
Создайте файл `src/functions.ts` и реализуйте функции:

```typescript
// 1. Функция сложения двух чисел
// 2. Функция приветствия, принимающая имя и возвращающая строку
// 3. Функция проверки возраста (возвращает boolean)
// 4. Функция, которая может принимать строку или число, а возвращает строку
// 5. Функция с опциональным параметром
// 6. Функция с параметром по умолчанию
// 7. Функция с rest-параметрами
```

**Пример реализации одной функции:**
```typescript
function greet(name: string, greeting: string = 'Hello'): string {
  return `${greeting}, ${name}!`;
}
```

### Задание 3: Работа с объектами и массивами
Создайте файл `src/data-structures.ts`:

1. Создайте интерфейс `Product` с полями: id, name, price, category, inStock
2. Создайте массив из минимум 5 продуктов
3. Создайте функцию `filterByCategory`, которая фильтрует продукты по категории
4. Создайте функцию `calculateTotalPrice`, которая считает общую стоимость товаров в корзине
5. Создайте тип `Cart` как массив объектов с продуктом и количеством

### Задание 4: Union и Intersection
Создайте файл `src/advanced-types.ts`:

1. Создайте тип `UserRole` как `'admin' | 'user' | 'guest'`
2. Создайте тип `User` с полями: id, name, email, role
3. Создайте тип `AdminUser` как пересечение `User` и `{ permissions: string[] }`
4. Напишите функцию `checkPermission`, которая проверяет права доступа
5. Используйте narrowing для обработки разных ролей

### Задание 5: Практический кейс - Система уведомлений
Создайте файл `src/notifications.ts`:

Реализуйте систему уведомлений с разными типами:
- `SuccessNotification` - содержит message и duration
- `ErrorNotification` - содержит message и errorCode
- `WarningNotification` - содержит message и warningLevel ('low' | 'medium' | 'high')

Создайте:
1. Union тип `Notification` из всех типов уведомлений
2. Функцию `showNotification`, которая принимает любой тип уведомления
3. Функцию `logNotification`, которая выводит разную информацию в зависимости от типа
4. Массив уведомлений и функцию для фильтрации по типу

## ❓ Проверочные вопросы

Ответьте на вопросы письменно в файле `answers.md`:

1. **В чем главное отличие TypeScript от JavaScript?**
2. **Когда стоит использовать тип `any`, а когда `unknown`? Приведите примеры.**
3. **Что такое type narrowing и какие способы сужения типа вы знаете?**
4. **В чем разница между `type` и `interface`? Когда что лучше использовать?**
5. **Объясните разницу между Union (`|`) и Intersection (`&`) типами. Приведите примеры использования.**
6. **Что такое кортежи и в каких ситуациях они полезны?**
7. **Как TypeScript обрабатывает опциональные параметры в функциях?**
8. **Что произойдет, если передать в функцию с типом `number` значение типа `string`? Почему?**
9. **Как типизировать функцию, которая ничего не возвращает? А которая никогда не завершает выполнение?**
10. **Зачем нужны литеральные типы и как они помогают в разработке?**

## ✅ Критерии выполнения модуля

- [ ] Все 5 практических заданий выполнены
- [ ] Код компилируется без ошибок (`tsc --noEmit`)
- [ ] На все проверочные вопросы даны развернутые ответы
- [ ] Код следует best practices (имена переменных, форматирование)
- [ ] Добавлены комментарии к сложным местам

## 🔗 Дополнительные ресурсы

- [Официальная документация TypeScript](https://www.typescriptlang.org/docs/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [TypeScript Playground](https://www.typescriptlang.org/play) - онлайн-песочница
- [Type Challenges](https://github.com/type-challenges/type-challenges) - задачи для практики
- [Книга "Programming TypeScript" by Boris Cherny](https://www.oreilly.com/library/view/programming-typescript/9781492037644/)

## 📝 Следующий шаг

После выполнения всех заданий переходите к следующему модулю: **02-interfaces-generics**
