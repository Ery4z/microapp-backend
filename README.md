# API Documentation

## Base URL
`http://<your-server-address>:1323`

## Endpoints

### 1. POST /data
- **Description**: Receives and stores data from sensors.
- **Parameters**:
  - `sensorId`: Unique ID of the sensor.
  - `groupId`: Unique ID of the group.
  - `dataType`: Type of data (e.g., string, int, float, date).
  - `dataUnit`: Unit of the data, if relevant.
  - `dataInfo`: Information about the data.
  - `data`: The actual data being sent.
- **Response**: Confirmation message with sensor ID, data type, and data.

### 2. POST /groups
- **Description**: Creates a new group.
- **Body**:
  - JSON object containing `groupId`, `name`, and `description`.
- **Response**: JSON object of the created group.

### 3. PUT /groups/:id
- **Description**: Updates an existing group.
- **Parameters**:
  - `id`: Group ID.
- **Body**:
  - JSON object with updated `name` and `description`.
- **Response**: JSON object of the updated group.

### 4. GET /groups/:id
- **Description**: Retrieves a specific group by ID.
- **Parameters**:
  - `id`: Group ID.
- **Response**: JSON object containing group details.

### 5. GET /groups
- **Description**: Lists all groups.
- **Response**: Array of JSON objects, each representing a group.

### 6. GET /groups/:groupId/sensors
- **Description**: Retrieves sensors for a specific group.
- **Parameters**:
  - `groupId`: Group ID.
- **Response**: JSON array of sensor data.

### 7. POST /groups/sensors
- **Description**: Retrieves sensor data for multiple groups.
- **Body**:
  - JSON object with an array of `groupIds`.
- **Response**: JSON object containing sensor data for the specified groups.
