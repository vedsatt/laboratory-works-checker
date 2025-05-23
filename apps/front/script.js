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
        const response = await fetch(`labs/lab${labNumber}/description.md`); //`laboratory-works-checker/labs/lab${labNumber}/description.md`
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

function loadTasks(labNumber, variant) {
    const tasksContainer = document.getElementById('tasks-container');
    tasksContainer.innerHTML = '';
    
    if (!variant) {
        tasksContainer.innerHTML = '<p>Введите вариант, чтобы увидеть задания</p>';
        return;
    }
    
    const variantNumber = parseInt(variant, 10);
    if (isNaN(variantNumber)) {
        tasksContainer.innerHTML = '<p>Некорректный номер варианта</p>';
        return;
    }
    
    const labData = tasksByLab[labNumber];
    if (!labData) {
        tasksContainer.innerHTML = '<p>Задания для этой работы не найдены</p>';
        return;
    }
    
    let tasksHTML = '<h3>Ваши задания:</h3><ol>';
    
    for (let part = 1; part <= labData.parts; part++) {
        const partKey = `part${part}`;
        const taskText = labData.templates[partKey](variantNumber);
        tasksHTML += `<li>${taskText}</li>`;
    }
    
    tasksHTML += '</ol>';
    tasksContainer.innerHTML = tasksHTML;
}

function submitSolution() {
    const labNumber = document.getElementById('lab-number').textContent;
    const variant = document.getElementById('variant').value;
    const code = document.getElementById('code').value;
    
    if (!variant || !code) {
        alert('Пожалуйста, заполните все поля');
        return;
    }
    
    const isSuccess = Math.random() > 0.3; // Заглушка для демонстрации

    addToHistory(labNumber, isSuccess, code);
    document.getElementById('code').value = '';
}

function addToHistory(labNumber, isSuccess, code) {
    const historyItem = document.createElement('div');
    historyItem.className = `history-item ${isSuccess ? 'success' : 'error'}`;
    
    const now = new Date();
    const timestamp = now.toLocaleString();
    
    historyItem.innerHTML = `
        <div class="status">${isSuccess ? '✓ Решение верное' : '✗ Ошибка в решении'}</div>
        <div class="timestamp">${timestamp}</div>
        <pre class="code-snippet">${code.substring(0, 100)}${code.length > 100 ? '...' : ''}</pre>
    `;
    
    document.getElementById('history-list').prepend(historyItem);
    saveToLocalStorage(labNumber, isSuccess, code, timestamp);
}

function saveToLocalStorage(labNumber, isSuccess, code, timestamp) {
    const history = JSON.parse(localStorage.getItem(`lab_${labNumber}_history`) || '[]');
    
    const newHistory = [{
        success: isSuccess,
        code: code,
        timestamp: timestamp
    }, ...history.slice(0, 9)]; 

    localStorage.setItem(`lab_${labNumber}_history`, JSON.stringify(newHistory));
}

function loadSolutionHistory(labNumber) {
    const history = JSON.parse(localStorage.getItem(`lab_${labNumber}_history`) || []);
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
            <div class="timestamp">${item.timestamp}</div>
            <pre class="code-snippet">${item.code.substring(0, 100)}${item.code.length > 100 ? '...' : ''}</pre>
        `;
        
        historyList.appendChild(historyItem);
    });
}