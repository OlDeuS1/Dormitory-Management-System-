const express = require('express');
const router = express.Router();

// In-memory store (replace with a real DB in production)
let students = [
  { id: 1, name: 'Alice Smith', email: 'alice@example.com', phone: '555-0101', roomId: null },
  { id: 2, name: 'Bob Jones',  email: 'bob@example.com',   phone: '555-0102', roomId: null },
];
let nextId = 3;

// GET /students
router.get('/', (req, res) => {
  res.json(students);
});

// GET /students/:id
router.get('/:id', (req, res) => {
  const student = students.find(s => s.id === parseInt(req.params.id));
  if (!student) return res.status(404).json({ error: 'Student not found' });
  res.json(student);
});

// POST /students
router.post('/', (req, res) => {
  const { name, email, phone } = req.body;
  if (!name || !email) {
    return res.status(400).json({ error: 'name and email are required' });
  }
  const student = { id: nextId++, name, email, phone: phone || '', roomId: null };
  students.push(student);
  res.status(201).json(student);
});

// PUT /students/:id
router.put('/:id', (req, res) => {
  const index = students.findIndex(s => s.id === parseInt(req.params.id));
  if (index === -1) return res.status(404).json({ error: 'Student not found' });
  students[index] = { ...students[index], ...req.body, id: students[index].id };
  res.json(students[index]);
});

// DELETE /students/:id
router.delete('/:id', (req, res) => {
  const index = students.findIndex(s => s.id === parseInt(req.params.id));
  if (index === -1) return res.status(404).json({ error: 'Student not found' });
  const removed = students.splice(index, 1)[0];
  res.json(removed);
});

module.exports = router;
