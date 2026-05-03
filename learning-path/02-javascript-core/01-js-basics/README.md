# 📜 Основы JavaScript

## 📚 Теория

### Что такое JavaScript?

JavaScript — это язык программирования, который делает веб-страницы интерактивными. Это язык браузеров, но также используется на сервере (Node.js).

### Переменные

```javascript
// var - устаревший способ (не рекомендуется)
var oldVariable = 'old';

// let - изменяемая переменная
let count = 5;
count = 10; // можно изменить

// const - константа (нельзя переназначить)
const PI = 3.14159;
// PI = 3; // Ошибка!

// Hoisting (поднятие)
console.log(hoisted); // undefined
var hoisted = 'value';

console.log(letHoisted); // Ошибка! Cannot access before initialization
let letHoisted = 'value';
```

### Типы данных

```javascript
// Примитивные типы
const string = 'Hello';           // String
const number = 42;                // Number
const bigint = 123n;              // BigInt
const boolean = true;             // Boolean
const nullValue = null;           // Null
const undefinedValue = undefined; // Undefined
const symbol = Symbol('id');      // Symbol

// Ссылочный тип
const object = { key: 'value' };  // Object
const array = [1, 2, 3];          // Array (тоже Object)
const func = () => {};            // Function (тоже Object)

// Проверка типа
typeof 'hello';        // 'string'
typeof 42;             // 'number'
typeof true;           // 'boolean'
typeof undefined;      // 'undefined'
typeof null;           // 'object' (баг JS!)
typeof {};             // 'object'
typeof [];             // 'object'
typeof function(){};   // 'function'

Array.isArray([]);     // true
```

### Операторы

```javascript
// Арифметические
5 + 3;    // 8
5 - 3;    // 2
5 * 3;    // 15
5 / 3;    // 1.666...
5 % 3;    // 2 (остаток)
5 ** 2;   // 25 (возведение в степень)
++count;  // инкремент
--count;  // декремент

// Строковые
'Hello' + ' ' + 'World';  // 'Hello World'

// Сравнения
5 == '5';   // true (нестрогое, с приведением типов)
5 === '5';  // false (строгое)
5 != '5';   // false
5 !== '5';  // true
5 > 3;      // true
5 < 3;      // false
5 >= 5;     // true
5 <= 4;     // false

// Логические
true && false;  // false (И)
true || false;  // true (ИЛИ)
!true;          // false (НЕ)

// Nullish coalescing
const value = null ?? 'default';  // 'default'
const value2 = 0 ?? 'default';    // 0

// Optional chaining
const name = user?.profile?.name;
```

### Условия

```javascript
// if/else
const age = 18;

if (age < 13) {
    console.log('Ребенок');
} else if (age < 20) {
    console.log('Подросток');
} else {
    console.log('Взрослый');
}

// Тернарный оператор
const status = age >= 18 ? 'adult' : 'minor';

// Switch
const day = 'Monday';

switch (day) {
    case 'Monday':
        console.log('Понедельник');
        break;
    case 'Tuesday':
        console.log('Вторник');
        break;
    default:
        console.log('Другой день');
}

// Truthy/Falsy значения
// Falsy: false, 0, -0, 0n, '', null, undefined, NaN
// Остальное - Truthy
```

### Циклы

```javascript
// for
for (let i = 0; i < 5; i++) {
    console.log(i);
}

// while
let i = 0;
while (i < 5) {
    console.log(i);
    i++;
}

// do...while
do {
    console.log(i);
    i++;
} while (i < 5);

// for...of (для итерируемых объектов)
const arr = [1, 2, 3];
for (const item of arr) {
    console.log(item);
}

// for...in (для ключей объекта)
const obj = { a: 1, b: 2 };
for (const key in obj) {
    console.log(key, obj[key]);
}

// break и continue
for (let i = 0; i < 10; i++) {
    if (i === 3) continue; // пропустить 3
    if (i === 7) break;    // остановить на 7
    console.log(i);
}
```

---

## 🎯 Практическое задание

### Задание 1: Калькулятор
Создайте калькулятор с функциями: сложение, вычитание, умножение, деление.

### Задание 2: Конвертер температур
Напишите функции для конвертации между Celsius и Fahrenheit.

### Задание 3: Простой анализатор строки
Создайте функции для: подсчета символов, слов, проверки палиндрома.

---

## ❓ Проверочные вопросы

1. В чем разница между let, const и var?
2. Какие типы данных в JavaScript вы знаете?
3. Что такое hoisting?
4. В чем разница между == и ===?
5. Какие falsy значения вы знаете?
6. Как работает optional chaining?
7. Что делает оператор ??
8. Когда использовать for...of vs for...in?

---

**Следующий шаг:** `02-functions-scope` - функции и область видимости!
