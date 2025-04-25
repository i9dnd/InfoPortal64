document.addEventListener('DOMContentLoaded', () => {
    fetchEvents(); 
    setInterval(fetchEvents, 5000); 
});

const eventList = document.getElementById('eventList');

async function fetchEvents() {
    const response = await fetch('/');
    const events = await response.json();
    renderEvents(events);
}

function renderEvents(events) {
    eventList.innerHTML = '';
    events.forEach(event => {
        const li = document.createElement('li');
        li.className = 'event-item';
        li.setAttribute('data-title', event.Title);
        li.innerHTML = `
            <strong>${event.Title}</strong> - ${event.Description}
            <a class="button edit" href="/edit/${event.ID}">Редактировать</a>
            <button class="button delete" data-id="${event.ID}">Удалить</button>
        `;
        eventList.appendChild(li);
    });

    attachDeleteEventListeners();
}

function attachDeleteEventListeners() {
    const deleteButtons = document.querySelectorAll('.delete');
    deleteButtons.forEach(button => {
        button.addEventListener('click', async () => {
            const eventId = button.getAttribute('data-id');
            const response = await fetch(`/delete/${eventId}`, { method: 'DELETE' });

            if (response.ok) {
                fetchEvents(); 
            } else {
                alert('Ошибка при удалении события.'); 
            }
        });
    });
}

function filterEvents() {
    const searchTerm = document.getElementById('search').value.toLowerCase();
    const items = document.querySelectorAll('.event-item');

    items.forEach(item => {
        const title = item.getAttribute('data-title').toLowerCase();
        item.style.display = title.includes(searchTerm) ? '' : 'none';
    });
}

document.getElementById('search').addEventListener('input', filterEvents);




