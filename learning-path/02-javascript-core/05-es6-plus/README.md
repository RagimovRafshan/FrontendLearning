# 🚀 ES6+

## Теория

### Деструктуризация

```javascript
const { name, age } = user;
const [first, second] = array;
```

### Spread/Rest

```javascript
const newArr = [...oldArr, 4, 5];
const newObj = { ...oldObj, c: 3 };
```

### Модули

```javascript
// export
export const value = 1;
export default function() {}

// import
import value, { other } from './module.js';
```

### Template Literals

```javascript
const str = `Hello ${name}`;
```

## Практика
1. Рефакторинг кода на ES6+
2. Создание модулей

**Следующий шаг:** `06-oop-functional`
