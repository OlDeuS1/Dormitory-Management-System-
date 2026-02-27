const express = require('express');
const app = express();
const roomRoutes = require('./routes/rooms');

app.use(express.json());
app.use('/rooms', roomRoutes);

app.get('/health', (req, res) => res.json({ status: 'ok', service: 'room-service' }));

const PORT = process.env.PORT || 3002;
app.listen(PORT, () => console.log(`Room Service running on port ${PORT}`));
