# Модуль 2: Интерфейсы и Дженерики в TypeScript

## 📚 Теория

### 2.1 Интерфейсы (Interfaces)

**Интерфейс** — это контракт, который описывает форму объекта. Он определяет, какие свойства и методы должны быть у объекта, но не реализует их.

```typescript
// Базовый интерфейс
interface User {
  id: number;
  name: string;
  email: string;
  age?: number; // Опциональное свойство
  readonly createdAt: Date; // Только для чтения
}

// Использование интерфейса
const user: User = {
  id: 1,
  name: 'Alice',
  email: 'alice@example.com',
  createdAt: new Date()
};

// user.createdAt = new Date(); // ❌ Ошибка: нельзя изменить readonly свойство
```

#### Расширение интерфейсов

```typescript
// Наследование интерфейсов
interface BaseUser {
  id: number;
  name: string;
}

interface Employee extends BaseUser {
  position: string;
  department: string;
  salary: number;
}

const employee: Employee = {
  id: 1,
  name: 'Bob',
  position: 'Developer',
  department: 'Engineering',
  salary: 5000
};

// Множественное наследование
interface Identifiable {
  id: number;
}

interface Timestampable {
  createdAt: Date;
  updatedAt: Date;
}

interface Post extends Identifiable, Timestampable {
  title: string;
  content: string;
  authorId: number;
}

const post: Post = {
  id: 1,
  createdAt: new Date(),
  updatedAt: new Date(),
  title: 'My First Post',
  content: 'Hello, world!',
  authorId: 1
};
```

#### Интерфейсы для функций

```typescript
// Интерфейс функции
interface SearchFunc {
  (source: string, subString: string): boolean;
}

const mySearch: SearchFunc = (src, sub) => {
  return src.includes(sub);
};

// Интерфейс с параметрами разной арности
interface MultiFunc {
  (a: number): number;
  (a: number, b: number): number;
  (a: number, b?: number): number {
    if (b !== undefined) {
      return a + b;
    }
    return a * 2;
  }
}
```

#### Интерфейсы для классов

```typescript
// Интерфейс для класса
interface ClockInterface {
  currentTime: Date;
  setTime(date: Date): void;
}

class Clock implements ClockInterface {
  currentTime: Date;
  
  constructor() {
    this.currentTime = new Date();
  }
  
  setTime(date: Date): void {
    this.currentTime = date;
  }
  
  displayTime(): string {
    return this.currentTime.toLocaleTimeString();
  }
}

// Реализация нескольких интерфейсов
interface Flyable {
  fly(): void;
}

interface Swimmable {
  swim(): void;
}

class Duck implements Flyable, Swimmable {
  fly(): void {
    console.log('Duck is flying');
  }
  
  swim(): void {
    console.log('Duck is swimming');
  }
}
```

### 2.2 Generics (Дженерики/Обобщения)

**Дженерики** позволяют создавать компоненты, которые работают с различными типами данных, сохраняя при этом типизацию.

#### Базовые дженерики

```typescript
// Функция без дженериков - теряем информацию о типе
function identity(arg: any): any {
  return arg;
}
const result = identity('hello'); // result имеет тип any

// Функция с дженериками - сохраняем тип
function identity<T>(arg: T): T {
  return arg;
}

const str = identity<string>('hello'); // str имеет тип string
const num = identity<number>(42); // num имеет тип number
const bool = identity(true); // TS сам выводит тип: bool имеет тип boolean

// Дженерик-переменные могут называться как угодно, но T - общепринятое соглашение
// T - Type, K/V - Key/Value, U/V - другие типы
```

#### Дженерики в интерфейсах

```typescript
// Дженерик интерфейс
interface Box<T> {
  value: T;
  getValue(): T;
  setValue(value: T): void;
}

class StringBox implements Box<string> {
  value: string;
  
  constructor(value: string) {
    this.value = value;
  }
  
  getValue(): string {
    return this.value;
  }
  
  setValue(value: string): void {
    this.value = value;
  }
}

// Дженерик интерфейс с несколькими параметрами
interface Pair<K, V> {
  key: K;
  value: V;
}

const userPair: Pair<number, string> = {
  key: 1,
  value: 'Alice'
};

const configPair: Pair<string, boolean> = {
  key: 'debug',
  value: true
};
```

#### Дженерики в классах

```typescript
// Дженерик класс
class Stack<T> {
  private items: T[] = [];
  
  push(item: T): void {
    this.items.push(item);
  }
  
  pop(): T | undefined {
    return this.items.pop();
  }
  
  peek(): T | undefined {
    return this.items[this.items.length - 1];
  }
  
  isEmpty(): boolean {
    return this.items.length === 0;
  }
  
  size(): number {
    return this.items.length;
  }
}

const numberStack = new Stack<number>();
numberStack.push(1);
numberStack.push(2);
numberStack.push(3);
console.log(numberStack.pop()); // 3

const stringStack = new Stack<string>();
stringStack.push('hello');
stringStack.push('world');
console.log(stringStack.pop()); // 'world'
```

#### Ограничения дженериков (Constraints)

```typescript
// Ограничение: T должен иметь свойство length
interface Lengthwise {
  length: number;
}

function logLength<T extends Lengthwise>(arg: T): T {
  console.log(`Length: ${arg.length}`);
  return arg;
}

logLength('hello'); // ✅ Length: 5
logLength([1, 2, 3]); // ✅ Length: 3
logLength({ length: 10 }); // ✅ Length: 10
// logLength(42); // ❌ Ошибка: number не имеет свойства length

// Ограничение с несколькими типами
function merge<T extends object, U extends object>(obj1: T, obj2: U): T & U {
  return { ...obj1, ...obj2 };
}

const merged = merge({ name: 'Alice' }, { age: 25 });
// merged имеет тип: { name: string } & { age: number }
console.log(merged.name); // 'Alice'
console.log(merged.age); // 25
```

#### Дженерики с ключами объекта (keyof)

```typescript
// keyof получает все ключи объекта как union тип
interface Person {
  name: string;
  age: number;
  email: string;
}

type PersonKeys = keyof Person; // 'name' | 'age' | 'email'

// Функция, которая принимает ключ объекта
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}

const person: Person = {
  name: 'Alice',
  age: 25,
  email: 'alice@example.com'
};

const name = getProperty(person, 'name'); // type: string
const age = getProperty(person, 'age'); // type: number
// const invalid = getProperty(person, 'invalid'); // ❌ Ошибка
```

#### Utility Types с дженериками

```typescript
// Partial - делает все свойства опциональными
interface Todo {
  title: string;
  description: string;
  completed: boolean;
}

function updateTodo(todo: Todo, fieldsToUpdate: Partial<Todo>): Todo {
  return { ...todo, ...fieldsToUpdate };
}

const todo1: Todo = {
  title: 'Learn TS',
  description: 'Study generics',
  completed: false
};

const updated = updateTodo(todo1, { completed: true });

// Required - делает все свойства обязательными
interface OptionalProps {
  name?: string;
  age?: number;
}

const required: Required<OptionalProps> = {
  name: 'Alice',
  age: 25
};

// Pick - выбирает определенные свойства
interface Article {
  id: number;
  title: string;
  content: string;
  author: string;
  publishedAt: Date;
}

type ArticlePreview = Pick<Article, 'id' | 'title' | 'author'>;

const preview: ArticlePreview = {
  id: 1,
  title: 'TypeScript Guide',
  author: 'Alice'
};

// Omit - исключает определенные свойства
type ArticleWithoutContent = Omit<Article, 'content'>;

// Record - создает объект с определенными ключами и значением
type UserRole = 'admin' | 'user' | 'guest';
interface Permissions {
  canRead: boolean;
  canWrite: boolean;
  canDelete: boolean;
}

const rolePermissions: Record<UserRole, Permissions> = {
  admin: { canRead: true, canWrite: true, canDelete: true },
  user: { canRead: true, canWrite: true, canDelete: false },
  guest: { canRead: true, canWrite: false, canDelete: false }
};

// Exclude и Extract
type Union = 'success' | 'error' | 'warning' | 'info';
type SuccessOrError = Extract<Union, 'success' | 'error'>; // 'success' | 'error'
type WarningOrInfo = Exclude<Union, 'success' | 'error'>; // 'warning' | 'info'
```

### 2.3 Продвинутые паттерны с дженериками

```typescript
// Условные типы (Conditional Types)
type IsString<T> = T extends string ? true : false;

type A = IsString<'hello'>; // true
type B = IsString<42>; // false

// Infer keyword для извлечения типов
type ReturnType<T> = T extends (...args: any[]) => infer R ? R : never;

function getString(): string {
  return 'hello';
}

type Str = ReturnType<typeof getString>; // string

// Мэппинг типов (Mapped Types)
type Readonly<T> = {
  readonly [P in keyof T]: T[P];
};

interface Mutable {
  name: string;
  age: number;
}

type Immutable = Readonly<Mutable>;
// { readonly name: string; readonly age: number; }

// Шаблонные литеральные типы (Template Literal Types)
type EventName = 'click' | 'hover' | 'focus';
type EventHandler = `on${Capitalize<EventName>}Handler`;
// 'onClickHandler' | 'onHoverHandler' | 'onFocusHandler'
```

## 🎯 Практические задания

### Задание 1: Создание системы интерфейсов
Создайте файл `src/interfaces.ts`:

1. Создайте базовый интерфейс `BaseEntity` с полями: `id`, `createdAt`, `updatedAt`
2. Создайте интерфейс `User`, расширяющий `BaseEntity`, с дополнительными полями: `name`, `email`, `role`
3. Создайте интерфейс `Product`, расширяющий `BaseEntity`, с полями: `name`, `price`, `category`, `stock`
4. Создайте интерфейс `Order`, содержащий: `user`, массив `products`, `total`, `status`
5. Реализуйте функцию `createOrder`, которая принимает пользователя и продукты, возвращает заказ

### Задание 2: Дженерик контейнеры
Создайте файл `src/generic-containers.ts`:

1. Реализуйте дженерик класс `Cache<T>` с методами:
   - `set(key: string, value: T, ttl?: number): void`
   - `get(key: string): T | undefined`
   - `delete(key: string): boolean`
   - `clear(): void`
   
2. Реализуйте дженерик класс `Repository<T>` с методами:
   - `findById(id: number): Promise<T | null>`
   - `findAll(): Promise<T[]>`
   - `create(item: Omit<T, 'id'>): Promise<T>`
   - `update(id: number, item: Partial<T>): Promise<T | null>`
   - `delete(id: number): Promise<boolean>`

3. Создайте конкретные реализации: `UserRepository` и `ProductRepository`

### Задание 3: Utility Types на практике
Создайте файл `src/utility-types.ts`:

1. Создайте интерфейс `FormConfig` с полями: `title`, `fields`, `submitUrl`, `method`
2. Используя `Partial<T>`, создайте функцию для обновления конфигурации формы
3. Используя `Pick<T>`, создайте тип `FormSummary` с полями `title` и `fields`
4. Используя `Omit<T>`, создайте тип `CreateFormInput` без поля `id`
5. Используя `Record<K, V>`, создайте словарь ошибок валидации по полям формы

### Задание 4: Продвинутые дженерики
Создайте файл `src/advanced-generics.ts`:

1. Реализуйте функцию `groupBy<T, K extends keyof T>(array: T[], key: K): Record<string, T[]>`
   - Группирует массив объектов по указанному ключу
   
2. Реализуйте функцию `deepClone<T>(obj: T): T` для глубокого клонирования объектов

3. Создайте тип `DeepReadonly<T>`, который делает все свойства рекурсивно readonly

4. Реализуйте тип `Flatten<T>`, который объединяет все свойства вложенных объектов

### Задание 5: Практический кейс - API Client
Создайте файл `src/api-client.ts`:

Реализуйте типизированный API клиент:

```typescript
interface ApiResponse<T> {
  data: T;
  status: number;
  message: string;
}

interface ApiClient {
  get<T>(endpoint: string): Promise<ApiResponse<T>>;
  post<T>(endpoint: string, data: unknown): Promise<ApiResponse<T>>;
  put<T>(endpoint: string, data: unknown): Promise<ApiResponse<T>>;
  delete<T>(endpoint: string): Promise<ApiResponse<T>>;
}
```

Дополнительно:
1. Добавьте интерцепторы запросов и ответов
2. Реализуйте кэширование GET запросов
3. Добавьте автоматическую повторную попытку при ошибках
4. Создайте типизированные методы для конкретных ресурсов (users, products, orders)

## ❓ Проверочные вопросы

Ответьте на вопросы письменно в файле `answers.md`:

1. **В чем разница между `interface` и `type`? Когда лучше использовать каждый из них?**
2. **Что такое structural typing в TypeScript и как он отличается от nominal typing?**
3. **Как работает наследование интерфейсов? Можно ли реализовать множественное наследование?**
4. **Что такое дженерики и зачем они нужны? Приведите 3 примера использования.**
5. **Объясните синтаксис `<T extends SomeType>`. Что означает эта конструкция?**
6. **Что делает оператор `keyof`? Приведите пример практического использования.**
7. **Как работают utility types `Partial`, `Pick`, `Omit`, `Record`? Когда какой использовать?**
8. **Что такое conditional types? Как используется ключевое слово `infer`?**
9. **Как создать собственный utility type? Приведите пример.**
10. **В чем разница между `any`, `unknown` и `never`? Когда использовать каждый тип?**

## ✅ Критерии выполнения модуля

- [ ] Все 5 практических заданий выполнены
- [ ] Код компилируется без ошибок (`tsc --noEmit`)
- [ ] На все проверочные вопросы даны развернутые ответы
- [ ] Использованы дженерики там, где это уместно
- [ ] Код следует принципам DRY и SOLID
- [ ] Добавлены JSDoc комментарии к публичным API

## 🔗 Дополнительные ресурсы

- [TypeScript Interfaces](https://www.typescriptlang.org/docs/handbook/interfaces.html)
- [TypeScript Generics](https://www.typescriptlang.org/docs/handbook/2/generics.html)
- [Utility Types](https://www.typescriptlang.org/docs/handbook/utility-types.html)
- [Advanced Types](https://www.typescriptlang.org/docs/handbook/2/types-from-types.html)
- [Type Challenges](https://github.com/type-challenges/type-challenges) - сложные задачи на типы

## 📝 Следующий шаг

После выполнения всех заданий переходите к следующему модулю: **03-advanced-types**
