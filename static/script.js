document.addEventListener('DOMContentLoaded', function() {
    fetchGroups();
    // Set an interval for automatic update every 30 seconds
    setInterval(fetchGroups, 1000);
});

function fetchGroups() {
    fetch('http://localhost:1323/groups')
        .then(response => response.json())
        .then(groups => updateGroups(groups))
        .catch(error => console.error('Error fetching groups:', error));
}

function updateGroups(groups) {
    const groupsContainer = document.getElementById('groups');
    // Create a Map of existing groups for quick access
    const existingGroupsMap = new Map();
    document.querySelectorAll('.group').forEach(groupDiv => {
        existingGroupsMap.set(groupDiv.id, groupDiv);
    });

    // Iterate over fetched groups and update the DOM accordingly
    groups.forEach(group => {
        let groupElement = existingGroupsMap.get(`group-${group.groupId}`);
        if (!groupElement) {
            // If the group doesn't exist on the page, create it
            groupElement = document.createElement('div');
            groupElement.id = `group-${group.groupId}`;
            groupElement.className = 'group';
            groupElement.innerHTML = `
                <h2>${group.name}</h2>
                <ul id="sensors-${group.groupId}"></ul>
            `;
            groupsContainer.appendChild(groupElement);
        } else {
            // If it exists, update the group name if necessary
            const groupNameElement = groupElement.querySelector('h2');
            if (groupNameElement.textContent !== group.name) {
                groupNameElement.textContent = group.name;
            }
            // Remove the group from the Map so we know it's been processed
            existingGroupsMap.delete(groupElement.id);
        }
        fetchSensors(group.groupId);
    });

    // Any groups left in the Map do not exist in the fetched groups and should be removed
    existingGroupsMap.forEach(groupDiv => groupDiv.remove());
}

function fetchSensors(groupId) {
    fetch(`http://localhost:1323/groups/${groupId}/sensors`)
        .then(response => response.json())
        .then(sensors => updateSensors(groupId, sensors))
        .catch(error => console.error(`Error fetching sensors for group ${groupId}:`, error));
}

function updateSensors(groupId, sensors) {
    const sensorsListId = `sensors-${groupId}`;
    const sensorsList = document.getElementById(sensorsListId);

    // Create a Map of existing sensor list items for quick access
    const existingSensorsMap = new Map();
    sensorsList.querySelectorAll('li').forEach(sensorLi => {
        existingSensorsMap.set(sensorLi.id, sensorLi);
    });

    // Iterate over fetched sensors and update the DOM accordingly
    sensors.forEach(sensor => {
        let sensorItem = existingSensorsMap.get(`sensor-${sensor.sensorId}`);
        if (!sensorItem) {
            // If the sensor item doesn't exist, create it
            sensorItem = document.createElement('li');
            sensorItem.id = `sensor-${sensor.sensorId}`;
            sensorItem.textContent = `Sensor ID: ${sensor.sensorId}, Last Measure: ${sensor.data}`;
            sensorsList.appendChild(sensorItem);
        } else {
            // If it exists, update the sensor data
            sensorItem.textContent = `Sensor ID: ${sensor.sensorId}, Last Measure: ${sensor.data}`;
            // Remove the sensor from the Map so we know it's been processed
            existingSensorsMap.delete(sensorItem.id);
        }
    });

    // Any sensors left in the Map do not exist in the fetched sensors and should be removed
    existingSensorsMap.forEach(sensorLi => sensorLi.remove());
}
