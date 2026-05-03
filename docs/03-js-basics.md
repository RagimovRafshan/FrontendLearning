# Модуль 3: Основы JavaScript

## 📖 Теория

### Что такое JavaScript?
JavaScript — это язык программирования, который позволяет делать веб-страницы интерактивными. Это один из трех основных технологий веба (HTML, CSS, JS).

### Переменные и типы данных

```javascript
// Объявление переменных
let mutableVar = 'можно изменить'; // блочная область видимости
const CONSTANT = 'нельзя изменить'; // блочная область видимости
var oldVar = 'устаревший способ'; // функциональная область видимости (не использовать!)

// Примитивные типы
const stringType = 'строка';
const numberType = 42;
const booleanType = true;
const nullType = null;
const undefinedType = undefined;
const symbolType = Symbol('unique');
const bigIntType = 9007199254740991n;

// Объекты и массивы
const obj = { key: 'value' };
const arr = [1, 2, 3];
```

### Операторы

```javascript
// Арифметические
const sum = 5 + 3;        // 8
const diff = 5 - 3;       // 2
const product = 5 * 3;    // 15
const quotient = 5 / 3;   // 1.666...
const remainder = 5 % 3;  // 2
const power = 5 ** 3;     // 125

// Сравнения
5 == '5'   // true (нестрогое, с приведением типов)
5 === '5'  // false (строгое)
5 != '5'   // false
5 !== '5'  // true
5 > 3      // true
5 >= 5     // true

// Логические
true && false  // false (И)
true || false  // true (ИЛИ)
!true          // false (НЕ)
```

### Условные конструкции

```javascript
// if/else
const age = 18;

if (age < 18) {
    console.log('Младше 18');
} else if (age === 18) {
    console.log('Ровно 18');
} else {
    console.log('Старше 18');
}

// Тернарный оператор
const status = age >= 18 ? 'adult' : 'minor';

// switch
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
const arr = ['a', 'b', 'c'];
for (const item of arr) {
    console.log(item); // 'a', 'b', 'c'
}

// for...in (для ключей объекта)
const obj = { a: 1, b: 2 };
for (const key in obj) {
    console.log(key, obj[key]); // 'a' 1, 'b' 2
}
```

### Функции

```javascript
// Function Declaration
function greet(name) {
    return `Hello, ${name}!`;
}

// Function Expression
const greetExpr = function(name) {
    return `Hello, ${name}!`;
};

// Arrow Function (стрелочная функция)
const greetArrow = (name) => `Hello, ${name}!`;

// Параметры по умолчанию
function greetDefault(name = 'Guest') {
    return `Hello, ${name}!`;
}

// Rest параметры
function sum(...numbers) {
    return numbers.reduce((acc, num) => acc + num, 0);
}

// Деструктуризация параметров
function printUser({ name, age }) {
    console.log(`${name}, ${age}`);
}
```

### Область видимости и замыкания

```javascript
// Глобальная область
const globalVar = 'global';

function outer() {
    // Внешняя функция
    const outerVar = 'outer';
    
    function inner() {
        // Внутренняя функция имеет доступ к внешним переменным
        console.log(outerVar); // 'outer'
        console.log(globalVar); // 'global'
    }
    
    return inner;
}

const closure = outer();
closure(); // Замыкание сохраняет доступ к outerVar
```

### Объекты

```javascript
// Создание объекта
const user = {
    name: 'John',
    age: 30,
    greet() {
        console.log(`Hello, I'm ${this.name}`);
    }
};

// Доступ к свойствам
console.log(user.name); // точечная нотация
console.log(user['age']); // квадратные скобки

// Деструктуризация
const { name, age } = user;

// Spread оператор
const updatedUser = { ...user, age: 31 };

// Методы Object
Object.keys(user);      // ['name', 'age']
Object.values(user);    // ['John', 30]
Object.entries(user);   // [['name', 'John'], ['age', 30]]
```

### Массивы

```javascript
const arr = [1, 2, 3, 4, 5];

// Основные методы
arr.push(6);           // добавить в конец
arr.pop();             // удалить из конца
arr.unshift(0);        // добавить в начало
arr.shift();           // удалить из начала

// Методы высшего порядка
const doubled = arr.map(x => x * 2);
const evens = arr.filter(x => x % 2 === 0);
const sum = arr.reduce((acc, x) => acc + x, 0);
const found = arr.find(x => x > 3);
const index = arr.findIndex(x => x > 3);
const exists = arr.some(x => x > 3);
const all = arr.every(x => x > 0);

// Slice и splice
const sliced = arr.slice(1, 3);     // копия части массива
arr.splice(1, 2, 'a', 'b');         // изменить исходный массив

// Sort
const sorted = arr.sort((a, b) => a - b);
```

### ES6+ фичи

```javascript
// Шаблонные строки
const name = 'John';
const greeting = `Hello, ${name}!`;

// Деструктуризация
const [first, second, ...rest] = [1, 2, 3, 4, 5];
const { name: userName, age } = user;

// Spread/Rest операторы
const merged = [...arr1, ...arr2];
const cloned = { ...obj };

// Optional chaining
const city = user?.address?.city;

// Nullish coalescing
const value = input ?? 'default'; // только null/undefined

// Модули (import/export)
// export const PI = 3.14;
// import { PI } from './math.js';
```

### Асинхронность

```javascript
// Callback (устаревший подход)
setTimeout(() => {
    console.log('Через 1 секунду');
}, 1000);

// Promise
const promise = new Promise((resolve, reject) => {
    setTimeout(() => {
        resolve('Успех!');
        // reject('Ошибка!');
    }, 1000);
});

promise
    .then(result => console.log(result))
    .catch(error => console.error(error));

// async/await
async function fetchData() {
    try {
        const response = await fetch('/api/data');
        const data = await response.json();
        return data;
    } catch (error) {
        console.error(error);
    }
}

// Promise.all
const [result1, result2] = await Promise.all([
    fetch('/api/1'),
    fetch('/api/2')
]);
```

### DOM манипуляции

```javascript
// Выбор элементов
const elem = document.querySelector('.class');
const elems = document.querySelectorAll('.class');
const byId = document.getElementById('id');

// Изменение содержимого
elem.textContent = 'Текст';
elem.innerHTML = '<span>HTML</span>';

// Атрибуты
elem.setAttribute('data-id', '123');
elem.getAttribute('data-id');
elem.removeAttribute('data-id');

// Классы
elem.classList.add('active');
elem.classList.remove('active');
elem.classList.toggle('active');
elem.classList.contains('active');

// Стили
elem.style.color = 'red';
elem.style.cssText = 'color: red; font-size: 16px;';

// Создание и удаление
const newElem = document.createElement('div');
elem.appendChild(newElem);
newElem.remove();

// События
elem.addEventListener('click', (event) => {
    event.preventDefault();
    console.log('Клик!', event.target);
});

elem.removeEventListener('click', handler);
```

---

## ✅ Практическое задание

### Задание 1: Калькулятор

Создайте калькулятор с функциями:
- Сложение, вычитание, умножение, деление
- История операций
- Очистка истории

Используйте модульную структуру кода.

### Задание 2: Todo List

Реализуйте список задач:
- Добавление задач
- Отметка о выполнении
- Удаление задач
- Фильтрация (все/активные/завершенные)
- Сохранение в localStorage

### Задание 3: Работа с API

Создайте приложение, которое:
- Загружает данные с публичного API (например, JSONPlaceholder)
- Отображает список постов/пользователей
- Реализует поиск/фильтрацию
- Обрабатывает ошибки загрузки

---

## 📝 Проверочные вопросы

1. В чем разница между `let`, `const` и `var`?
2. Что такое замыкание? Приведите пример использования.
3. Чем отличаются `==` и `===`?
4. Какие методы массивов изменяют оригинальный массив?
5. Что такое Promise и какие у него состояния?
6. Как работает event bubbling и event capturing?
7. Что такое hoisting (поднятие)?

---

## 🎯 Критерии выполнения

- [ ] Код разбит на логические функции/модули
- [ ] Использованы современные возможности ES6+
- [ ] Обработаны краевые случаи и ошибки
- [ ] Код отформатирован (Prettier/ESLint)
- [ ] Нет утечек памяти (правильная работа с событиями)

---

## 🔗 Полезные ресурсы

- [MDN JavaScript](https://developer.mozilla.org/ru/docs/Web/JavaScript)
- [JavaScript.info](https://learn.javascript.ru/)
- [ES6 Features](https://github.com/lukehoban/es6features)
- [Promise MDN](https://developer.mozilla.org/ru/docs/Web/JavaScript/Reference/Global_Objects/Promise)

---

**Следующий модуль:** Продвинутый JavaScript →
