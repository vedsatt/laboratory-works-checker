const md = window.markdownit({
  html: true,
  linkify: true,
  typographer: true,
});

document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const labNumber = urlParams.get('lab') || '1';

    if (!labNumber || labNumber < 1 || labNumber > 11) {
        window.location.href = 'index.html';
        return;
    }

    document.getElementById('lab-number').textContent = labNumber;
    document.getElementById('lab-title').textContent = `Лабораторная работа №${labNumber}`;

    loadLabDescription(labNumber);

    document.getElementById('submit-btn').addEventListener('click', submitSolution);

    loadSolutionHistory(labNumber);
});

async function loadLabDescription(labNumber) {
    try {
        const response = await fetch(`../../labs/lab${labNumber}/description.md`); //`laboratory-works-checker/labs/lab${labNumber}/description.md`
        if (!response.ok) throw new Error('Описание не найдено');
        
        const mdContent = await response.text();
        const contentBeforeAssignment = mdContent.split('## Задание')[0]; 
        const htmlContent = md.render(contentBeforeAssignment); // Рендерим только нужную часть

        document.getElementById('lab-description').innerHTML = htmlContent;
    } catch (error) {
        document.getElementById('lab-description').innerHTML = `
            <p>Не удалось загрузить описание работы: ${error.message}</p>
        `;
    }
}

async function loadTasks(labNumber) {
    console.log('--- Начало выполнения loadTasks() ---');
    
    console.log('1. Получен номер лабораторной:', labNumber);
    
    const variantInput = document.getElementById('variant').value;
    console.log('2. Получен ввод варианта:', variantInput);
    
    const variant = parseInt(variantInput) || 0;
    console.log('3. Числовое значение варианта:', variant);
    
    const tasksContainer = document.getElementById('tasks-container');
    console.log('4. Контейнер заданий получен:', tasksContainer);
    
    try {
        console.log('5. Пытаюсь загрузить файл с описанием...');
        const filePath = `../../labs/lab${labNumber}/description.md`;
        console.log('6. Путь к файлу:', filePath);
        
        const response = await fetch(filePath);
        console.log('7. Ответ от fetch:', response);
        
        if (!response.ok) {
            console.error('8. Ошибка: файл не найден или ошибка сервера');
            throw new Error('Файл с описанием не найден');
        }
        
        const mdContent = await response.text();
        console.log('9. Содержимое файла получено, длина:', mdContent.length);
        console.log('10. Первые 100 символов:', mdContent.substring(0, 100));
        
        // Разделяем на части
        const parts = mdContent.split('### Часть').slice(1);
        console.log('11. Найдено частей:', parts.length);
        
        tasksContainer.innerHTML = '';
        console.log('12. Контейнер заданий очищен');
        
        if (parts.length === 0) {
            console.log('13. Нет частей для отображения');
            tasksContainer.innerHTML = '<p>Задания не найдены</p>';
            return;
        }

        parts.forEach((part, partIndex) => {
            console.log(`14. Обработка части ${partIndex + 1}`);
            
            const items = part.split(/\n\d+\./).slice(1); // разделяем часть на пункты
            console.log(`15. В части ${partIndex + 1} найдено пунктов:`, items.length);
            
            if (items.length === 0) {
                console.log(`16. Часть ${partIndex + 1} не содержит пунктов`);
                return;
            }

            const selectedIndex = variant % items.length; // вычисляем индекс выбранного пункта
            console.log(`17. Для варианта ${variant} выбран индекс:`, selectedIndex);
            
            const selectedItem = items[selectedIndex].trim();
            console.log(`18. Выбранный пункт:`, selectedItem);
            
            // Создаем элемент для отображения задания
            const partElement = document.createElement('div');
            partElement.className = 'task-part';
            
            const partTitle = document.createElement('h3');
            partTitle.textContent = `Часть ${partIndex + 1}`;
            
            const taskContent = document.createElement('div');
            console.log('19. Пытаюсь отрендерить markdown...');
            
            if (!window.markdownit) {
                console.error('20. Ошибка: markdownit не загружен!');
                taskContent.textContent = `1.${selectedItem}`;
            } else {
                const md = window.markdownit();
                taskContent.innerHTML = md.render(`${selectedItem}`);
                console.log('21. Markdown успешно отрендерен');
            }
            
            partElement.appendChild(partTitle);
            partElement.appendChild(taskContent);
            tasksContainer.appendChild(partElement);
            console.log(`22. Часть ${partIndex + 1} добавлена в контейнер`);
        });
        
        console.log('23. Все части успешно обработаны');
        
    } catch (error) {
        console.error('24. Произошла ошибка:', error);
        tasksContainer.innerHTML = `<p>Ошибка загрузки заданий: ${error.message}</p>`;
    }
    
    console.log('--- Конец выполнения loadTasks() ---');
}


document.getElementById('variant').addEventListener('input', function() {
    const labNumber = new URLSearchParams(window.location.search).get('lab') || '1';
    const variant = this.value.trim();
    
    if (variant) {
        loadTasks(labNumber, variant);
    } else {
        document.getElementById('tasks-container').innerHTML = '<p>Введите вариант, чтобы увидеть задания</p>';
    }
});

// Конфигурация сервера (должна быть объявлена в начале файла)
const SERVER_CONFIG = {
    host: window.location.hostname || '127.0.0.1',
    port: window.location.port || '80',
    apiPath: '/api/v1/submit'
};

// Логирование в консоль с меткой времени
function log(message, data = null) {
    const timestamp = new Date().toISOString();
    console.log(`[${timestamp}] ${message}`);
    if (data) {
        console.log('Data:', JSON.stringify(data, null, 2));
    }
}

// Функция для преобразования кода в строку с экранированными символами
function escapeCodeString(code) {
    log('Начало преобразования кода в строку');
    try {
        const escaped = code
            .replace(/\\/g, '\\\\')
            .replace(/"/g, '\\"')
            .replace(/\n/g, '\\n')
            .replace(/\r/g, '\\r')
            .replace(/\t/g, '\\t');
        log('Код успешно преобразован в строку');
        return escaped;
    } catch (error) {
        log('Ошибка при преобразовании кода в строку', { error: error.message });
        throw error;
    }
}

async function submitSolution() {
    const labNumber = new URLSearchParams(window.location.search).get('lab') || '1';
    const variantInput = document.getElementById('variant').value;
    const variant = parseInt(variantInput) || 0;
    const codeEditor = document.getElementById('code');
    const code = codeEditor.value;
    
    if (!code.trim()) {
        alert('Пожалуйста, введите код для отправки');
        return;
    }

    if (!variantInput) {
        alert('Пожалуйста, укажите вариант');
        return;
    }

    try {
        // Собираем данные для отправки
        const requestData = {
            id: Date.now(),
            lab_number: parseInt(labNumber),
            code: escapeCodeString(code),
            tasks: {},
            task: 0
        };

        // Получаем информацию о вариантах для каждой части
        const filePath = `../../labs/lab${labNumber}/description.md`;
        const descriptionResponse = await fetch(filePath);
        
        if (!descriptionResponse.ok) throw new Error('Файл с описанием не найден');
        
        const mdContent = await descriptionResponse.text();
        const parts = mdContent.split('### Часть').slice(1);
        
        // Заполняем Tasks
        parts.forEach((part, partIndex) => {
            const items = part.split(/\n\d+\./).slice(1);
            if (items.length > 0) {
                const selectedIndex = variant % items.length;
                requestData.tasks[`task${partIndex + 1}`] = selectedIndex + 1;
            }
        });

        // Если нужно отправить только текущую часть
        const currentTaskElement = document.querySelector('input[name="current-task"]:checked');
        if (currentTaskElement) {
            requestData.task = parseInt(currentTaskElement.value);
        }

        log('Подготовка данных для отправки', requestData);

        // Отправляем запрос на сервер
        //const serverUrl = `http://${SERVER_CONFIG.host}:${SERVER_CONFIG.port}${SERVER_CONFIG.apiPath}`;
        const serverUrl = `http://localhost:80/api/v1/submit`
        const serverResponse = await fetch(serverUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData)
        });

        if (!serverResponse.ok) {
            const errorData = await serverResponse.json();
            throw new Error(errorData.message || 'Ошибка сервера');
        }

        const result = await serverResponse.json();
        log('Ответ сервера', result);
        alert('Решение успешно отправлено на проверку!');
        
        // Обновляем историю решений
        loadSolutionHistory(labNumber);

    } catch (error) {
        log('Ошибка при отправке решения', { error: error.message });
        alert(`Ошибка: ${error.message}`);
    }
}

// Добавляем обработчик события после загрузки DOM
document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('submit-btn').addEventListener('click', submitSolution);
});