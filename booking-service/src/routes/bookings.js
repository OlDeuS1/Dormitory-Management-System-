const express = require('express');
const router = express.Router();

// In-memory store (replace with a real DB in production)
let bookings = [];
let nextId = 1;

// GET /bookings
router.get('/', (req, res) => {
  res.json(bookings);
});

// GET /bookings/:id
router.get('/:id', (req, res) => {
  const booking = bookings.find(b => b.id === parseInt(req.params.id));
  if (!booking) return res.status(404).json({ error: 'Booking not found' });
  res.json(booking);
});

// POST /bookings
router.post('/', (req, res) => {
  const { studentId, roomId, checkIn, checkOut } = req.body;
  if (!studentId || !roomId || !checkIn || !checkOut) {
    return res.status(400).json({ error: 'studentId, roomId, checkIn and checkOut are required' });
  }
  const booking = {
    id: nextId++,
    studentId: parseInt(studentId),
    roomId: parseInt(roomId),
    checkIn,
    checkOut,
    status: 'active',
    createdAt: new Date().toISOString(),
  };
  bookings.push(booking);
  res.status(201).json(booking);
});

// PUT /bookings/:id
router.put('/:id', (req, res) => {
  const index = bookings.findIndex(b => b.id === parseInt(req.params.id));
  if (index === -1) return res.status(404).json({ error: 'Booking not found' });
  bookings[index] = { ...bookings[index], ...req.body, id: bookings[index].id };
  res.json(bookings[index]);
});

// DELETE /bookings/:id  (cancel booking)
router.delete('/:id', (req, res) => {
  const index = bookings.findIndex(b => b.id === parseInt(req.params.id));
  if (index === -1) return res.status(404).json({ error: 'Booking not found' });
  bookings[index].status = 'cancelled';
  res.json(bookings[index]);
});

module.exports = router;
