# 🔧 Функции и Область Видимости

## Теория

### Объявление функций

```javascript
// Function Declaration (поднимается)
function add(a, b) {
    return a + b;
}

// Function Expression (не поднимается)
const subtract = function(a, b) {
    return a - b;
};

// Arrow Function (стрелочная)
const multiply = (a, b) => a * b;

// Со многими параметрами
const divide = (a, b = 1) => a / b; // значение по умолчанию

// Rest параметры
function sum(...numbers) {
    return numbers.reduce((acc, n) => acc + n, 0);
}
```

### Область видимости (Scope)

```javascript
// Глобальная область
const global = 'global';

function outer() {
    // Внешняя функция
    const outerVar = 'outer';
    
    function inner() {
        // Внутренняя функция
        const innerVar = 'inner';
        console.log(outerVar); // работает
    }
}

// Блочная область (let, const)
{
    let blockLet = 'let';
    const blockConst = 'const';
    var blockVar = 'var';
}
console.log(blockVar); // работает
console.log(blockLet); // ошибка
```

### Замыкания (Closures)

```javascript
function createCounter() {
    let count = 0;
    
    return function() {
        count++;
        return count;
    };
}

const counter = createCounter();
counter(); // 1
counter(); // 2
```

## Практика
1. Создайте функцию debounce
2. Реализуйте curry функцию
3. Создайте мемоизацию

**Следующий шаг:** `03-dom-events`
