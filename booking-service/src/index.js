const express = require('express');
const app = express();
const bookingRoutes = require('./routes/bookings');

app.use(express.json());
app.use('/bookings', bookingRoutes);

app.get('/health', (req, res) => res.json({ status: 'ok', service: 'booking-service' }));

const PORT = process.env.PORT || 3003;
app.listen(PORT, () => console.log(`Booking Service running on port ${PORT}`));
