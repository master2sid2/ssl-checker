document.addEventListener('DOMContentLoaded', function() {
    // Скрытие сообщения через 10 секунд
    setTimeout(function() {
        var errorMessage = document.querySelector('.error-message');
        if (errorMessage) {
            errorMessage.style.display = 'none';
        }
    }, 5000); // 10000 миллисекунд = 10 секунд
});

document.addEventListener('DOMContentLoaded', function() {
    const table = document.getElementById('domainTable');
    const headers = table.querySelectorAll('th[data-sort]');
    let sortDirection = 1;
    let currentSortColumn = null;

    headers.forEach(header => {
        header.addEventListener('click', function() {
            const sortType = this.getAttribute('data-sort');
            sortTableByColumn(table, sortType);
            updateSortDirection(this);
        });
    });

    function sortTableByColumn(table, sortType) {
        const tbody = table.querySelector('tbody');
        const rows = Array.from(tbody.querySelectorAll('tr'));

        rows.sort((rowA, rowB) => {
            const cellA = rowA.querySelector(`td:nth-child(${getIndexBySortType(sortType)})`).innerText.trim();
            const cellB = rowB.querySelector(`td:nth-child(${getIndexBySortType(sortType)})`).innerText.trim();

            if (sortType === 'validUntil') {
                return compareDates(cellA, cellB) * sortDirection;
            } else if (sortType === 'daysLeft') {
                return (parseInt(cellA) - parseInt(cellB)) * sortDirection;
            } else {
                return cellA.localeCompare(cellB) * sortDirection;
            }
        });

        rows.forEach(row => tbody.appendChild(row)); // Переставляем строки в таблице
    }

    function updateSortDirection(header) {
        if (currentSortColumn && currentSortColumn !== header) {
            currentSortColumn.querySelector('.sort-arrow').innerText = '';
        }

        if (header === currentSortColumn) {
            sortDirection *= -1; // Меняем направление сортировки
        } else {
            sortDirection = 1;
        }

        currentSortColumn = header;
        header.querySelector('.sort-arrow').innerText = sortDirection === 1 ? '▲' : '▼';
    }

    function getIndexBySortType(sortType) {
        switch (sortType) {
            case 'domain':
                return 1;
            case 'validUntil':
                return 3;
            case 'daysLeft':
                return 4;
            default:
                return 0;
        }
    }

    function compareDates(dateA, dateB) {
        const a = new Date(dateA);
        const b = new Date(dateB);
        return a - b;
    }
});