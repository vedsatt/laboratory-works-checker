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

// Конфигурация сервера по умолчанию
const SERVER_CONFIG = {
    host: process.env.SERVER_HOST || '127.0.0.1',
    port: process.env.SERVER_PORT || '80',
    apiPath: '/api/v1/submit'
};

async function submitSolution(labNumber) {
    //const labNumber = document.getElementById('labNumber').textContent;
    const variant = document.getElementById('variant').value.trim();
    const code = document.getElementById('code').value.trim();
    
    if (!variant) {
        alert('Пожалуйста, введите вариант');
        return;
    }
    
    if (!code) {
        alert('Пожалуйста, введите код решения');
        return;
    }
    
    try {
        // Формируем данные для отправки
        const requestData = {
            id: 1, // сгенерировать уникальный айди
            lab_number: parseInt(labNumber),
            code: code, // код-строка под json
            tasks: {} // мапа: ключ task№, значение - вариация таски
        };
        
        // Отправляем запрос на сервер
        const response = await sendToServer(requestData);
        
        // Обрабатываем ответ
        const isSuccess = response.ResMsg === 'OK';
        const errorMessage = isSuccess ? '' : response.ResMsg;
        
        addToHistory(labNumber, isSuccess, code, errorMessage);
        
        // Показываем результат пользователю
        const resultMessage = isSuccess 
            ? '<div class="success-message">✓ Решение верное!</div>'
            : `<div class="error-message">✗ Ошибка в решении: ${errorMessage}</div>
               <div class="error-note">Строки, начинающиеся с #>, а также описание уведомлений и ошибок не учитываются при проверке ответа</div>`;
        
        const historyList = document.getElementById('history-list');
        historyList.insertAdjacentHTML('afterbegin', resultMessage);
        
        // Очищаем поле кода, если решение верное
        if (isSuccess) {
            document.getElementById('code').value = '';
        }
    } catch (error) {
        console.error('Ошибка при отправке решения:', error);
        alert('Произошла ошибка при проверке решения. Пожалуйста, попробуйте позже.');
    }
}

async function sendToServer(data) {
    const url = `http://${SERVER_CONFIG.host}:${SERVER_CONFIG.port}$/api/v1/submit`;
    
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    });
    
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
}

function addToHistory(labNumber, isSuccess, code, errorMessage = '') {
    const historyItem = document.createElement('div');
    historyItem.className = `history-item ${isSuccess ? 'success' : 'error'}`;
    
    const now = new Date();
    const timestamp = now.toLocaleString();
    
    historyItem.innerHTML = `
        <div class="status">${isSuccess ? '✓ Решение верное' : '✗ Ошибка в решении'}</div>
        ${errorMessage ? `<div class="error-details">${errorMessage}</div>` : ''}
        <div class="timestamp">${timestamp}</div>
        <pre class="code-snippet">${code.substring(0, 100)}${code.length > 100 ? '...' : ''}</pre>
    `;
    
    document.getElementById('history-list').prepend(historyItem);
    saveToLocalStorage(labNumber, isSuccess, code, timestamp, errorMessage);
}

function saveToLocalStorage(labNumber, isSuccess, code, timestamp, errorMessage = '') {
    const history = JSON.parse(localStorage.getItem(`lab_${labNumber}_history`) || '[]');
    
    const newHistory = [{
        success: isSuccess,
        code: code,
        timestamp: timestamp,
        errorMessage: errorMessage
    }, ...history.slice(0, 9)]; 

    localStorage.setItem(`lab_${labNumber}_history`, JSON.stringify(newHistory));
}

function loadSolutionHistory(labNumber) {
    const history = JSON.parse(localStorage.getItem(`lab_${labNumber}_history`) || '[]');
    const historyList = document.getElementById('history-list');
    
    historyList.innerHTML = '';
    
    if (history.length === 0) {
        historyList.innerHTML = '<p>Здесь будет отображаться история ваших решений</p>';
        return;
    }
    
    history.forEach(item => {
        const historyItem = document.createElement('div');
        historyItem.className = `history-item ${item.success ? 'success' : 'error'}`;
        
        historyItem.innerHTML = `
            <div class="status">${item.success ? '✓ Решение верное' : '✗ Ошибка в решении'}</div>
            ${item.errorMessage ? `<div class="error-details">${item.errorMessage}</div>` : ''}
            <div class="timestamp">${item.timestamp}</div>
            <pre class="code-snippet">${item.code.substring(0, 100)}${item.code.length > 100 ? '...' : ''}</pre>
        `;
        
        historyList.appendChild(historyItem);
    });
}