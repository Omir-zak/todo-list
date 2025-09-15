import { GetTasks, AddTask, DeleteTask, ToggleTask, GetCombinedFilteredTasks } from '../wailsjs/go/backend/App.js';

let todos = [];
let currentFilters = {
    status: 'all',
    date: 'all',
    sort: '',
    ascending: false
};
let taskToDelete = null;

// Инициализация приложения
document.addEventListener('DOMContentLoaded', async () => {
    await loadTodos();
    setupEventListeners();
    updateStats();
});

// Настройка обработчиков событий
function setupEventListeners() {
    // Форма добавления задачи
    const addBtn = document.getElementById('add-btn');
    const todoTitle = document.getElementById('todo-title');
    
    addBtn.addEventListener('click', addTodo);
    todoTitle.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            addTodo();
        }
    });

    // Фильтры и сортировка
    const statusFilter = document.getElementById('status-filter');
    const dateFilter = document.getElementById('date-filter');
    const sortFilter = document.getElementById('sort-filter');
    const sortOrder = document.getElementById('sort-order');

    statusFilter.addEventListener('change', (e) => {
        currentFilters.status = e.target.value;
        applyFilters();
    });

    dateFilter.addEventListener('change', (e) => {
        currentFilters.date = e.target.value;
        applyFilters();
    });

    sortFilter.addEventListener('change', (e) => {
        currentFilters.sort = e.target.value;
        applyFilters();
    });

    sortOrder.addEventListener('click', () => {
        currentFilters.ascending = !currentFilters.ascending;
        sortOrder.classList.toggle('ascending', currentFilters.ascending);
        sortOrder.textContent = currentFilters.ascending ? '↑' : '↓';
        applyFilters();
    });

    // Модальное окно
    const deleteModal = document.getElementById('delete-modal');
    const confirmDelete = document.getElementById('confirm-delete');
    const cancelDelete = document.getElementById('cancel-delete');

    confirmDelete.addEventListener('click', async () => {
        if (taskToDelete) {
            await deleteTodo(taskToDelete);
            taskToDelete = null;
        }
        hideModal();
    });

    cancelDelete.addEventListener('click', () => {
        taskToDelete = null;
        hideModal();
    });

    // Закрытие модального окна по клику вне его
    deleteModal.addEventListener('click', (e) => {
        if (e.target === deleteModal) {
            taskToDelete = null;
            hideModal();
        }
    });
}

// Загрузка задач
async function loadTodos() {
    try {
        todos = await GetTasks();
        renderTodos();
        updateStats();
    } catch (error) {
        console.error('Ошибка загрузки задач:', error);
        showNotification('Не удалось загрузить задачи', 'error');
    }
}

// Применение фильтров
async function applyFilters() {
    try {
        const filteredTodos = await GetCombinedFilteredTasks(
            currentFilters.status,
            currentFilters.date,
            currentFilters.sort,
            currentFilters.ascending
        );
        todos = filteredTodos;
        renderTodos();
        updateStats();
    } catch (error) {
        console.error('Ошибка применения фильтров:', error);
        showNotification('Ошибка фильтрации задач', 'error');
    }
}

// Отображение задач
function renderTodos() {
    const todoList = document.getElementById('todo-list');

    if (todos.length === 0) {
        todoList.innerHTML = '<div class="empty-message">Нет задач для отображения</div>';
        return;
    }

    todoList.innerHTML = todos.map(todo => {
        const dueDate = todo.due_date ? new Date(todo.due_date) : null;
        const createdDate = new Date(todo.created_at);
        const now = new Date();
        const isOverdue = dueDate && dueDate < now && !todo.completed;

        return `
            <div class="todo-item ${todo.completed ? 'completed' : ''} priority-${todo.priority}" data-id="${todo.id}">
                <div class="todo-content">
                    <input type="checkbox" ${todo.completed ? 'checked' : ''} 
                           onchange="toggleTodo(${todo.id})">
                    <div class="todo-text">
                        <h3>${escapeHtml(todo.title)}</h3>
                        ${todo.description ? `<p>${escapeHtml(todo.description)}</p>` : ''}
                        <div class="todo-meta">
                            <span class="priority priority-${todo.priority}">
                                ${getPriorityText(todo.priority)}
                            </span>
                            ${dueDate ? `
                                <span class="due-date ${isOverdue ? 'overdue' : ''}">
                                    ${isOverdue ? '⚠️ ' : '📅 '}${formatDate(dueDate)}
                                </span>
                            ` : ''}
                            <span class="created-date">
                                Создано: ${formatDate(createdDate, false)}
                            </span>
                        </div>
                    </div>
                </div>
                <div class="todo-actions">
                    <button class="delete-btn" onclick="showDeleteModal(${todo.id})">
                        🗑️ Удалить
                    </button>
                </div>
            </div>
        `;
    }).join('');
}

// Добавление новой задачи
async function addTodo() {
    const titleInput = document.getElementById('todo-title');
    const descriptionInput = document.getElementById('todo-description');
    const priorityInput = document.getElementById('todo-priority');
    const dueDateInput = document.getElementById('todo-duedate');

    const title = titleInput.value.trim();
    const description = descriptionInput.value.trim();
    const priority = priorityInput.value;
    const dueDate = dueDateInput.value;

    if (!title) {
        showNotification('Введите название задачи', 'error');
        titleInput.focus();
        return;
    }

    try {
        await AddTask(title, description, priority, dueDate);
        
        // Очищаем форму
        titleInput.value = '';
        descriptionInput.value = '';
        priorityInput.value = 'medium';
        dueDateInput.value = '';

        await loadTodos();
        showNotification('Задача добавлена', 'success');
    } catch (error) {
        console.error('Ошибка добавления задачи:', error);
        showNotification('Не удалось добавить задачу', 'error');
    }
}

// Переключение статуса задачи
window.toggleTodo = async function(id) {
    try {
        await ToggleTask(id);
        await loadTodos();
        
        const task = todos.find(t => t.id === id);
        const message = task?.completed ? 'Задача выполнена' : 'Задача возвращена в активные';
        showNotification(message, 'success');
    } catch (error) {
        console.error('Ошибка изменения статуса:', error);
        showNotification('Не удалось изменить статус задачи', 'error');
    }
};

// Показать модальное окно подтверждения удаления
window.showDeleteModal = function(id) {
    taskToDelete = id;
    const modal = document.getElementById('delete-modal');
    modal.classList.add('show');
};

// Скрыть модальное окно
function hideModal() {
    const modal = document.getElementById('delete-modal');
    modal.classList.remove('show');
}

// Удаление задачи
async function deleteTodo(id) {
    try {
        await DeleteTask(id);
        await loadTodos();
        showNotification('Задача удалена', 'success');
    } catch (error) {
        console.error('Ошибка удаления задачи:', error);
        showNotification('Не удалось удалить задачу', 'error');
    }
}

// Обновление статистики
function updateStats() {
    const totalTasks = document.getElementById('total-tasks');
    const activeTasks = document.getElementById('active-tasks');
    const completedTasks = document.getElementById('completed-tasks');

    const total = todos.length;
    const active = todos.filter(todo => !todo.completed).length;
    const completed = todos.filter(todo => todo.completed).length;

    totalTasks.textContent = `Всего: ${total}`;
    activeTasks.textContent = `Активных: ${active}`;
    completedTasks.textContent = `Выполненных: ${completed}`;
}

// Вспомогательные функции
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function getPriorityText(priority) {
    const priorities = {
        'low': 'Низкий',
        'medium': 'Средний',
        'high': 'Высокий'
    };
    return priorities[priority] || 'Средний';
}

function formatDate(date, includeTime = true) {
    if (!date) return '';
    
    const options = {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    };
    
    if (includeTime) {
        options.hour = '2-digit';
        options.minute = '2-digit';
    };
    
    return new Intl.DateTimeFormat('ru-RU', options).format(date);
}

function showNotification(message, type) {
    // Удаляем предыдущие уведомления
    const existing = document.querySelector('.notification');
    if (existing) {
        existing.remove();
    }

    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;

    document.body.appendChild(notification);

    // Показываем уведомление с анимацией
    setTimeout(() => {
        notification.classList.add('show');
    }, 10);

    // Автоматически убираем уведомление через 4 секунды
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            if (notification.parentNode) {
                notification.remove();
            }
        }, 300);
    }, 4000);
}
