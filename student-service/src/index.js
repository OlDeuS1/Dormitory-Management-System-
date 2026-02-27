const express = require('express');
const app = express();
const studentRoutes = require('./routes/students');

app.use(express.json());
app.use('/students', studentRoutes);

app.get('/health', (req, res) => res.json({ status: 'ok', service: 'student-service' }));

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => console.log(`Student Service running on port ${PORT}`));
