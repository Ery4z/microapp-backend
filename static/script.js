let currentGroupIds = [];

document.addEventListener('DOMContentLoaded', function() {
    fetchGroups();
    setInterval(fetchGroups, 20000); // Poll for new groups every 20 seconds
    setInterval(updateSensorData, 5000); // Update sensor values every 5 seconds
});

function fetchGroups() {
    fetch('http://localhost:1323/groups')
        .then(response => response.json())
        .then(groups => {
            currentGroupIds = groups.map(group => group.groupId);
            renderGroups(groups);
        })
        .catch(error => console.error('Error fetching groups:', error));
}

function renderGroups(groups) {
    const groupsContainer = document.getElementById('groups');
    groups.forEach(group => {
        let groupElement = document.getElementById(`group-${group.groupId}`);
        if (!groupElement) {
            groupElement = document.createElement('div');
            groupElement.id = `group-${group.groupId}`;
            groupElement.className = 'group';
            groupElement.innerHTML = `
                <h2>${group.name}</h2>
                <ul id="sensors-${group.groupId}"></ul>
            `;
            groupsContainer.appendChild(groupElement);
        }
    });
}
function updateSensorData() {
    if (currentGroupIds.length === 0) {
        console.log("No group IDs to fetch sensor data for.");
        return; // No groups to update
    }

    fetch('http://localhost:1323/groups/sensors', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ groupIds: currentGroupIds })
    })
    .then(response => response.json())
    .then(groupSensorsData => {
        console.log("Fetched sensor data:", groupSensorsData);

        // Iterate over each group in the object
        Object.entries(groupSensorsData).forEach(([groupId, sensors]) => {
            updateGroupSensors(groupId, sensors);
        });
    })
    .catch(error => console.error('Error fetching group sensors:', error));
}

function updateGroupSensors(groupId, sensors) {
    const sensorsList = document.getElementById(`sensors-${groupId}`);
    if (!sensorsList) {
        console.error(`No sensors list element found for group ID ${groupId}`);
        return; // Group not rendered yet
    }

    sensorsList.innerHTML = ''; // Clear existing sensor data
    sensors.forEach(sensor => {
        const sensorItem = document.createElement('li');
        sensorItem.textContent = `Sensor ID: ${sensor.sensorId}, Last Measure: ${sensor.data}`;
        sensorsList.appendChild(sensorItem);
    });
}

