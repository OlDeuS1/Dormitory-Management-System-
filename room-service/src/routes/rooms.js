const express = require('express');
const router = express.Router();

// In-memory store (replace with a real DB in production)
let rooms = [
  { id: 1, number: '101', type: 'single', capacity: 1, available: true },
  { id: 2, number: '102', type: 'double', capacity: 2, available: true },
  { id: 3, number: '201', type: 'double', capacity: 2, available: true },
];
let nextId = 4;

// GET /rooms
router.get('/', (req, res) => {
  const { available } = req.query;
  if (available !== undefined) {
    const filtered = rooms.filter(r => r.available === (available === 'true'));
    return res.json(filtered);
  }
  res.json(rooms);
});

// GET /rooms/:id
router.get('/:id', (req, res) => {
  const room = rooms.find(r => r.id === parseInt(req.params.id));
  if (!room) return res.status(404).json({ error: 'Room not found' });
  res.json(room);
});

// POST /rooms
router.post('/', (req, res) => {
  const { number, type, capacity } = req.body;
  if (!number || !type || !capacity) {
    return res.status(400).json({ error: 'number, type and capacity are required' });
  }
  const room = { id: nextId++, number, type, capacity: parseInt(capacity), available: true };
  rooms.push(room);
  res.status(201).json(room);
});

// PUT /rooms/:id
router.put('/:id', (req, res) => {
  const index = rooms.findIndex(r => r.id === parseInt(req.params.id));
  if (index === -1) return res.status(404).json({ error: 'Room not found' });
  rooms[index] = { ...rooms[index], ...req.body, id: rooms[index].id };
  res.json(rooms[index]);
});

// DELETE /rooms/:id
router.delete('/:id', (req, res) => {
  const index = rooms.findIndex(r => r.id === parseInt(req.params.id));
  if (index === -1) return res.status(404).json({ error: 'Room not found' });
  const removed = rooms.splice(index, 1)[0];
  res.json(removed);
});

module.exports = router;
