const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();

const STUDENT_SERVICE_URL = process.env.STUDENT_SERVICE_URL || 'http://localhost:3001';
const ROOM_SERVICE_URL    = process.env.ROOM_SERVICE_URL    || 'http://localhost:3002';
const BOOKING_SERVICE_URL = process.env.BOOKING_SERVICE_URL || 'http://localhost:3003';

// Route /students/* -> student-service
app.use('/students', createProxyMiddleware({ target: STUDENT_SERVICE_URL, changeOrigin: true }));

// Route /rooms/* -> room-service
app.use('/rooms', createProxyMiddleware({ target: ROOM_SERVICE_URL, changeOrigin: true }));

// Route /bookings/* -> booking-service
app.use('/bookings', createProxyMiddleware({ target: BOOKING_SERVICE_URL, changeOrigin: true }));

app.get('/health', (req, res) => res.json({ status: 'ok', service: 'api-gateway' }));

app.use((req, res) => res.status(404).json({ error: 'Route not found' }));

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log(`API Gateway running on port ${PORT}`));
