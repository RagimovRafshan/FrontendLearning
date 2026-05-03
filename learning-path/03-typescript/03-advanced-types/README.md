# Модуль 3: Продвинутые типы в TypeScript

## 📚 Теория

### 3.1 Type Guards и Type Predicates

**Type Guards** — это выражения, которые позволяют сузить тип значения в определенной области видимости.

#### typeof и instanceof guards

```typescript
// typeof guard
function printId(id: number | string) {
  if (typeof id === 'string') {
    // Здесь id имеет тип string
    console.log(id.toUpperCase());
  } else {
    // Здесь id имеет тип number
    console.log(id.toFixed(2));
  }
}

// instanceof guard
class Dog {
  bark() {
    console.log('Woof!');
  }
}

class Cat {
  meow() {
    console.log('Meow!');
  }
}

function makeSound(animal: Dog | Cat) {
  if (animal instanceof Dog) {
    animal.bark(); // animal имеет тип Dog
  } else {
    animal.meow(); // animal имеет тип Cat
  }
}
```

#### User-defined type guards

```typescript
// Type predicate - функция, которая возвращает boolean и указывает тип
interface Fish {
  swim(): void;
}

interface Bird {
  fly(): void;
}

function isFish(pet: Fish | Bird): pet is Fish {
  return (pet as Fish).swim !== undefined;
}

const myPet: Fish | Bird = { swim: () => console.log('Swimming') };

if (isFish(myPet)) {
  myPet.swim(); // ✅ myPet имеет тип Fish
} else {
  myPet.fly(); // ✅ myPet имеет тип Bird
}

// Type guard с проверкой свойства
interface Admin {
  role: 'admin';
  permissions: string[];
}

interface User {
  role: 'user';
  name: string;
}

function isAdmin(user: Admin | User): user is Admin {
  return user.role === 'admin';
}

function handleUser(user: Admin | User) {
  if (isAdmin(user)) {
    console.log(user.permissions); // user имеет тип Admin
  } else {
    console.log(user.name); // user имеет тип User
  }
}
```

#### In keyword

```typescript
interface Circle {
  kind: 'circle';
  radius: number;
}

interface Square {
  kind: 'square';
  side: number;
}

type Shape = Circle | Square;

function getArea(shape: Shape): number {
  if ('radius' in shape) {
    // shape имеет тип Circle
    return Math.PI * shape.radius ** 2;
  } else {
    // shape имеет тип Square
    return shape.side ** 2;
  }
}
```

### 3.2 Discriminated Unions

**Discriminated Unions** (также известные как tagged unions) — это паттерн, использующий общее поле для различения типов.

```typescript
interface SuccessResponse {
  status: 'success';
  data: unknown;
}

interface ErrorResponse {
  status: 'error';
  error: string;
  code: number;
}

interface LoadingResponse {
  status: 'loading';
  progress: number;
}

type ApiResponse = SuccessResponse | ErrorResponse | LoadingResponse;

function handleResponse(response: ApiResponse) {
  switch (response.status) {
    case 'success':
      console.log('Data:', response.data); // response имеет тип SuccessResponse
      break;
    case 'error':
      console.log('Error:', response.error, 'Code:', response.code);
      break;
    case 'loading':
      console.log('Progress:', response.progress);
      break;
  }
}

// Практический пример: состояние компонента
type ComponentState<T> =
  | { state: 'idle' }
  | { state: 'loading'; progress?: number }
  | { state: 'success'; data: T }
  | { state: 'error'; error: Error };

function renderComponent<T>(componentState: ComponentState<T>): string {
  switch (componentState.state) {
    case 'idle':
      return 'Ready to load';
    case 'loading':
      return `Loading... ${componentState.progress ?? 0}%`;
    case 'success':
      return `Data: ${JSON.stringify(componentState.data)}`;
    case 'error':
      return `Error: ${componentState.error.message}`;
  }
}
```

### 3.3 Mapped Types

**Mapped Types** позволяют создавать новые типы на основе существующих путем трансформации свойств.

```typescript
// Базовый mapped type
type Readonly<T> = {
  readonly [P in keyof T]: T[P];
};

interface Person {
  name: string;
  age: number;
}

type ReadonlyPerson = Readonly<Person>;
// { readonly name: string; readonly age: number; }

// Partial - делает все свойства опциональными
type Partial<T> = {
  [P in keyof T]?: T[P];
};

// Required - делает все свойства обязательными
type Required<T> = {
  [P in keyof T]-?: T[P];
};

// Pick - выбирает определенные свойства
type Pick<T, K extends keyof T> = {
  [P in K]: T[P];
};

// Omit - исключает определенные свойства
type Omit<T, K extends keyof T> = {
  [P in Exclude<keyof T, K>]: T[P];
};

// Custom mapped types с модификаторами
type Mutable<T> = {
  -readonly [P in keyof T]: T[P]; // Убираем readonly
};

type Nullable<T> = {
  [P in keyof T]: T[P] | null; // Добавляем null к каждому свойству
};

// Мэппинг с изменением ключей
type Getters<T> = {
  [P in keyof T as `get${Capitalize<string & P>}`]: () => T[P];
};

interface User {
  id: number;
  name: string;
}

type UserGetters = Getters<User>;
// { getId: () => number; getName: () => string; }

// Фильтрация свойств через as
type RemoveNever<T> = {
  [P in keyof T as T[P] extends never ? never : P]: T[P];
};

interface Data {
  a: string;
  b: never;
  c: number;
}

type CleanData = RemoveNever<Data>;
// { a: string; c: number; }
```

### 3.4 Template Literal Types

**Template Literal Types** позволяют создавать строковые типы на основе шаблонных строк.

```typescript
// Базовое использование
type Greeting = `Hello, ${string}`;
const greeting: Greeting = 'Hello, World'; // ✅
const invalid: Greeting = 'Hi'; // ❌

// Комбинирование union типов
type EventName = 'click' | 'hover' | 'focus';
type EventHandler = `on${Capitalize<EventName>}Handler`;
// 'onClickHandler' | 'onHoverHandler' | 'onFocusHandler'

const handler: EventHandler = 'onClickHandler'; // ✅

// Несколько union типов
type Alignment = 'top' | 'middle' | 'bottom';
type Position = 'left' | 'center' | 'right';
type ClassNames = `align-${Alignment}-${Position}`;
// 'align-top-left' | 'align-top-center' | ... (9 комбинаций)

// Инференс из шаблонных строк
type ParseUrl<T extends string> = T extends `${string}://${infer Host}/${infer Path}`
  ? { host: Host; path: Path }
  : { host: string; path: string };

type MyUrl = ParseUrl<'https://example.com/api/users'>;
// { host: 'example.com'; path: 'api/users' }

// Рекурсивные шаблонные типы
type ToUpperCase<T extends string> = T extends `${infer First}${infer Rest}`
  ? `${Uppercase<First>}${ToUpperCase<Rest>}`
  : T;

type Result = ToUpperCase<'hello'>; // 'HELLO'
```

### 3.5 Conditional Types

**Conditional Types** позволяют выбирать тип на основе условия.

```typescript
// Базовый синтаксис
type IsString<T> = T extends string ? true : false;

type A = IsString<'hello'>; // true
type B = IsString<42>; // false

// Distributive conditional types (распределительные)
type ToArray<T> = T extends any ? T[] : never;

type StrOrNum = ToArray<string | number>;
// string[] | number[] (не (string | number)[])

// Multiple conditions
type TypeOf<T> = 
  T extends string ? 'string' :
  T extends number ? 'number' :
  T extends boolean ? 'boolean' :
  T extends undefined ? 'undefined' :
  T extends Function ? 'function' :
  'object';

type A1 = TypeOf<'hello'>; // 'string'
type A2 = TypeOf<42>; // 'number'
type A3 = TypeOf<true>; // 'boolean'

// Infer в conditional types
type ReturnType<T> = T extends (...args: any[]) => infer R ? R : never;
type PromiseType<T> = T extends Promise<infer U> ? U : never;
type ArrayItemType<T> = T extends (infer U)[] ? U : never;

function getString(): string {
  return 'hello';
}

type Str = ReturnType<typeof getString>; // string
type Num = PromiseType<Promise<number>>; // number
type Item = ArrayItemType<string[]>; // string

// Complex infer example
type FirstArgument<T> = T extends (arg: infer A, ...args: any[]) => any ? A : never;

function createUser(name: string, age: number): void {}
type NameType = FirstArgument<typeof createUser>; // string
```

### 3.6 Advanced Utility Types

```typescript
// DeepPartial - рекурсивно делает все свойства опциональными
type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

interface Config {
  database: {
    host: string;
    port: number;
    credentials: {
      username: string;
      password: string;
    };
  };
}

type PartialConfig = DeepPartial<Config>;
// Все свойства становятся опциональными на всех уровнях

// DeepReadonly - рекурсивно делает все свойства readonly
type DeepReadonly<T> = {
  readonly [P in keyof T]: T[P] extends object ? DeepReadonly<T[P]> : T[P];
};

// Flatten - объединяет вложенные свойства
type Flatten<T> = T extends object
  ? { [K in keyof T]: T[K] }
  : T;

// Merge - объединяет два типа
type Merge<T, U> = {
  [K in keyof T | keyof U]: K extends keyof T
    ? K extends keyof U
      ? T[K] | U[K]
      : T[K]
    : K extends keyof U
    ? U[K]
    : never;
};

// TupleToObject - преобразует кортеж в объект
type TupleToObject<T extends readonly PropertyKey[]> = {
  [K in T[number]]: K;
};

const colors = ['red', 'green', 'blue'] as const;
type Colors = TupleToObject<typeof colors>;
// { red: 'red'; green: 'green'; blue: 'blue'; }

// Functions utilities
type NoInfer<T> = [T][T extends any ? 0 : never];
type Cast<T, U> = T extends U ? T : U;
type Writable<T> = { -readonly [P in keyof T]: T[P] };
```

## 🎯 Практические задания

### Задание 1: Type Guards и Narrowing
Создайте файл `src/type-guards.ts`:

1. Создайте интерфейсы для разных типов платежей: `CreditCardPayment`, `PayPalPayment`, `CryptoPayment`
2. Реализуйте type guard функции для каждого типа платежа
3. Создайте функцию `processPayment`, которая обрабатывает разные типы платежей по-разному
4. Используйте discriminated unions для представления состояния заказа

### Задание 2: Mapped Types практика
Создайте файл `src/mapped-types.ts`:

1. Реализуйте собственный `DeepRequired<T>` тип
2. Создайте тип `History<T>`, который добавляет ко всем свойствам суффикс `History` и меняет тип на массив значений
3. Реализуйте тип `Serialize<T>`, который преобразует все свойства в строки
4. Создайте тип `Asyncify<T>`, который делает все методы асинхронными (возвращают Promise)

### Задание 3: Template Literal Types
Создайте файл `src/template-literals.ts`:

1. Создайте тип `CSSProperty` для генерации названий CSS свойств (например, `margin-top`, `padding-bottom`)
2. Реализуйте тип `EventEmitter<T>` для генерации имен событий (`on${EventName}`, `${EventName}Changed`)
3. Создайте тип для парсинга URL路径 параметров: `/users/:id/posts/:postId` → `{ id: string; postId: string }`
4. Реализуйте систему типизированных API endpoints: `GET /api/${Resource}/${Id}`

### Задание 4: Conditional Types
Создайте файл `src/conditional-types.ts`:

1. Реализуйте тип `NonNullable<T>`, который удаляет null и undefined из union типа
2. Создайте тип `ExtractFunctionProps<T>`, который извлекает только функциональные свойства объекта
3. Реализуйте тип `ConstructorParameters<T>`, который извлекает параметры конструктора класса
4. Создайте тип `DeepOmit<T, Filter>`, который рекурсивно удаляет свойства из объекта

### Задание 5: Практический кейс - Типизированная форма
Создайте файл `src/typed-form.ts`:

Реализуйте полностью типизированную систему форм:

```typescript
// 1. Создайте тип FormField<T> для описания поля формы
// 2. Создайте тип FormConfig<T> для конфигурации всей формы
// 3. Реализуйте валидатор с типизированными ошибками
// 4. Создайте type-safe функцию для получения значений формы
// 5. Реализуйте систему условной валидации на основе значений других полей
```

Дополнительно:
- Добавьте поддержку nested объектов
- Реализуйте conditional fields (поле показывается только при определенном условии)
- Создайте тип для transformation функций

## ❓ Проверочные вопросы

Ответьте на вопросы письменно в файле `answers.md`:

1. **Что такое type guard? Какие виды type guards вы знаете?**
2. **Как работает type predicate и зачем он нужен? Приведите пример.**
3. **Что такое discriminated union? Почему это полезный паттерн?**
4. **Объясните, как работают mapped types. Приведите 3 примера использования.**
5. **В чем разница между `Partial<T>` и `DeepPartial<T>`? Когда использовать каждый?**
6. **Что такое template literal types? Как их можно использовать на практике?**
7. **Как работают conditional types? Объясните синтаксис `T extends U ? X : Y`.**
8. **Что делает ключевое слово `infer` в conditional types? Приведите примеры.**
9. **Что такое distributive conditional types? Как они отличаются от обычных conditional types?**
10. **Как создать собственный utility type? Опишите процесс шаг за шагом.**

## ✅ Критерии выполнения модуля

- [ ] Все 5 практических заданий выполнены
- [ ] Код компилируется без ошибок (`tsc --noEmit`)
- [ ] На все проверочные вопросы даны развернутые ответы
- [ ] Использованы продвинутые техники типизации
- [ ] Код хорошо документирован комментариями
- [ ] Реализованы edge cases и обработка ошибок

## 🔗 Дополнительные ресурсы

- [TypeScript Handbook - Advanced Types](https://www.typescriptlang.org/docs/handbook/2/types-from-types.html)
- [Type Guards and Differentiating Types](https://www.typescriptlang.org/docs/handbook/2/narrowing.html)
- [Mapped Types](https://www.typescriptlang.org/docs/handbook/2/mapped-types.html)
- [Template Literal Types](https://www.typescriptlang.org/docs/handbook/2/template-literal-types.html)
- [Conditional Types](https://www.typescriptlang.org/docs/handbook/2/conditional-types.html)
- [Type Challenges](https://github.com/type-challenges/type-challenges) - хардкорные задачи

## 📝 Следующий шаг

После выполнения всех заданий переходите к следующему модулю: **04-react-with-ts**
