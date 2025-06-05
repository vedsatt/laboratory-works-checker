const md = window.markdownit({
  html: true,
  linkify: true,
  typographer: true,
});

function loadSolutionHistory() {
    const historyList = document.getElementById('history-list');
    const history = JSON.parse(localStorage.getItem('solutionsHistory') || '[]');

    historyList.innerHTML = '';

    if (history.length === 0) {
        historyList.innerHTML = '<p>История решений пуста.</p>';
        return;
    }

    // Сортируем по дате (новые сверху)
    history.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));

    history.forEach((item, index) => {
        const solutionElement = document.createElement('div');
        
        // Проверяем, содержит ли ответ сервера "OK"
        const isSuccess = item.serverResponse && 
                         item.serverResponse.res_msg && 
                         item.serverResponse.res_msg.includes("OK");
        
        solutionElement.className = `history-item ${isSuccess ? 'success' : 'error'}`;

        if (isSuccess) {
            // Для успешных решений - просто статус и дата
            solutionElement.innerHTML = `
                <div class="status">
                    ✅ Верно 
                    <small>${new Date(item.timestamp).toLocaleString()}</small>
                </div>
            `;
        } else {
            // Для ошибок - добавляем кнопку для просмотра деталей
            solutionElement.innerHTML = `
                <div class="status">
                    ❌ Ошибка
                    <small>${new Date(item.timestamp).toLocaleString()}</small>
                </div>
                <button class="toggle-details">Показать детали</button>
                <div class="details" style="display: none;">
                    <h4>Ответ сервера:</h4>
                    <pre class="server-response">${formatServerResponse(item.serverResponse)}</pre>
                </div>
            `;

            const toggleBtn = solutionElement.querySelector('.toggle-details');
            const details = solutionElement.querySelector('.details');

            toggleBtn.addEventListener('click', () => {
                details.style.display = details.style.display === 'none' ? 'block' : 'none';
                toggleBtn.textContent = details.style.display === 'none' 
                    ? 'Показать детали' 
                    : 'Скрыть детали';
            });
        }

        historyList.appendChild(solutionElement);
    });
}

// Функция для форматирования ответа сервера
function formatServerResponse(response) {
    if (!response || !response.res_msg) return 'Нет данных об ответе сервера';
    
    // Форматируем сообщение для лучшего отображения
    let formattedMsg = response.res_msg
        .replace(/\\n/g, '\n')  // Заменяем \n на переносы строк
        .replace(/\\r/g, '\r')  // Заменяем \r на возврат каретки
        .replace(/\\t/g, '\t'); // Заменяем \t на табуляцию
    
    // Удаляем лишние кавычки, если они есть
    if (formattedMsg.startsWith('"') && formattedMsg.endsWith('"')) {
        formattedMsg = formattedMsg.slice(1, -1);
    }
    
    return formattedMsg;
}

document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const labNumber = urlParams.get('lab') || '1';
    const codeEditor = CodeMirror.fromTextArea(document.getElementById('code'), {
        mode: 'text/x-c++src', // Важно: меняем режим с Python на C++
        theme: 'dracula',
        lineNumbers: true,
        indentUnit: 4,
        tabSize: 4,
        lineWrapping: true,
        matchBrackets: true,
        autoCloseBrackets: true,
        extraKeys: {
            "Ctrl-Space": "autocomplete", // Добавляем автодополнение
            "Tab": function(cm) {
                cm.replaceSelection("    ", "end");
            }
        }
    });

    window.codeEditor = codeEditor;

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
    const variantInput = document.getElementById('variant').value;
    const variant = parseInt(variantInput) || 0;
    const tasksContainer = document.getElementById('tasks-container');

    try {
        const filePath = `../../labs/lab${labNumber}/description.md`;
        const response = await fetch(filePath);
        if (!response.ok) {
            throw new Error('Файл с описанием не найден');
        }
        const mdContent = await response.text();    
        
        // Разделяем на части
        const parts = mdContent.split('### Часть').slice(1);
        tasksContainer.innerHTML = '';

        if (parts.length === 0) {
            tasksContainer.innerHTML = '<p>Задания не найдены</p>';
            return;
        }

        parts.forEach((part, partIndex) => {            
            const items = part.split(/\n\d+\./).slice(1); // разделяем часть на пункты            
            if (items.length === 0) {
                return;
            }

            const selectedIndex = variant % items.length; // вычисляем индекс выбранного пункта
            const selectedItem = items[selectedIndex].trim();    

            // Создаем элемент для отображения задания
            const partElement = document.createElement('div');
            partElement.className = 'task-part';            
            const partTitle = document.createElement('h3');
            partTitle.textContent = `Часть ${partIndex + 1}`;            
            const taskContent = document.createElement('div');            

            if (!window.markdownit) {
                taskContent.textContent = `1.${selectedItem}`;
            } else {
                const md = window.markdownit();
                taskContent.innerHTML = md.render(`${selectedItem}`);
            }  

            partElement.appendChild(partTitle);
            partElement.appendChild(taskContent);
            tasksContainer.appendChild(partElement);
        });
    } catch (error) {
        tasksContainer.innerHTML = `<p>Ошибка загрузки заданий: ${error.message}</p>`;
    }  
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
        // Экранируем все специальные символы для JSON
        return JSON.stringify(code).slice(1, -1);
    } catch (error) {
        log('Ошибка при преобразовании кода в строку', { error: error.message });
        throw error;
    }
}

function saveToHistory(solution) {
    const history = JSON.parse(localStorage.getItem('solutionsHistory') || '[]');
    history.push(solution);
    localStorage.setItem('solutionsHistory', JSON.stringify(history));
}

async function submitSolution() {
    const labNumber = new URLSearchParams(window.location.search).get('lab') || '1';
    const variantInput = document.getElementById('variant').value;
    const variant = parseInt(variantInput) || 0;
    const code = window.codeEditor.getValue();
    
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
            code: code, // Не экранируем здесь - это сделает JSON.stringify
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
        const serverUrl = `http://localhost:80/api/v1/submit`; // Исправлено submit вместо submit
        const serverResponse = await fetch(serverUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData) // JSON.stringify правильно экранирует все спецсимволы
        });

        if (!serverResponse.ok) {
            const errorData = await serverResponse.json();
            throw new Error(errorData.message || 'Ошибка сервера');
        }

        const result = await serverResponse.json();
        log('Ответ сервера', result);
        saveToHistory({
            code: code,
            serverResponse: result,
            status: serverResponse.ok ? 'success' : 'error',
            timestamp: new Date().toISOString()
        });        
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