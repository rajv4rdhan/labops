'use strict';

const el = {
  list: document.getElementById('note-list'),
  newBtn: document.getElementById('new-note'),
  editor: document.getElementById('editor'),
  empty: document.getElementById('empty'),
  title: document.getElementById('title'),
  body: document.getElementById('body'),
  save: document.getElementById('save'),
  del: document.getElementById('delete'),
  status: document.getElementById('status'),
};

const state = {
  notes: [],
  currentId: null,
};

async function request(method, url, payload) {
  const res = await fetch(url, {
    method,
    headers: payload ? { 'Content-Type': 'application/json' } : undefined,
    body: payload ? JSON.stringify(payload) : undefined,
  });
  if (res.status === 204) return null;
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`);
  return data;
}

const api = {
  list: () => request('GET', '/api/notes'),
  create: (note) => request('POST', '/api/notes', note),
  update: (id, note) => request('PUT', `/api/notes/${id}`, note),
  remove: (id) => request('DELETE', `/api/notes/${id}`),
};

function preview(note) {
  const text = (note.body || '').replace(/\s+/g, ' ').trim();
  if (text) return text;
  return new Date(note.updated_at).toLocaleString();
}

function renderList() {
  el.list.innerHTML = '';
  if (state.notes.length === 0) {
    const li = document.createElement('li');
    li.className = 'note-item';
    li.textContent = 'No notes yet.';
    li.style.cursor = 'default';
    el.list.appendChild(li);
    return;
  }

  for (const note of state.notes) {
    const li = document.createElement('li');
    li.className = 'note-item' + (note.id === state.currentId ? ' active' : '');
    li.tabIndex = 0;
    li.addEventListener('click', () => selectNote(note.id));
    li.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        selectNote(note.id);
      }
    });

    const title = document.createElement('div');
    title.className = 'note-title' + (note.title ? '' : ' is-empty');
    title.textContent = note.title || 'Untitled note';

    const meta = document.createElement('div');
    meta.className = 'note-meta';
    meta.textContent = preview(note);

    li.append(title, meta);
    el.list.appendChild(li);
  }
}

function showEditor(show) {
  el.editor.hidden = !show;
  el.empty.hidden = show;
}

function selectNote(id) {
  const note = state.notes.find((n) => n.id === id);
  if (!note) return;
  state.currentId = id;
  el.title.value = note.title;
  el.body.value = note.body;
  el.del.hidden = false;
  showEditor(true);
  renderList();
  setStatus('');
}

function newNote() {
  state.currentId = null;
  el.title.value = '';
  el.body.value = '';
  el.del.hidden = true;
  showEditor(true);
  renderList();
  el.title.focus();
  setStatus('new');
}

async function refresh() {
  state.notes = await api.list();
  renderList();
}

async function save() {
  const payload = { title: el.title.value, body: el.body.value };
  if (!payload.title.trim() && !payload.body.trim()) {
    setStatus('empty');
    return;
  }
  try {
    let note;
    if (state.currentId === null) {
      note = await api.create(payload);
      state.currentId = note.id;
      el.del.hidden = false;
    } else {
      note = await api.update(state.currentId, payload);
    }
    await refresh();
    setStatus('saved');
  } catch (err) {
    setStatus('error: ' + err.message);
  }
}

async function remove() {
  if (state.currentId === null) return;
  if (!confirm('Delete this note?')) return;
  try {
    await api.remove(state.currentId);
    state.currentId = null;
    el.title.value = '';
    el.body.value = '';
    el.del.hidden = true;
    showEditor(false);
    await refresh();
    setStatus('');
  } catch (err) {
    setStatus('error: ' + err.message);
  }
}

let statusTimer;
function setStatus(text) {
  el.status.textContent = text;
  clearTimeout(statusTimer);
  if (text === 'saved') {
    statusTimer = setTimeout(() => {
      el.status.textContent = '';
    }, 1500);
  }
}

el.newBtn.addEventListener('click', newNote);
el.save.addEventListener('click', save);
el.del.addEventListener('click', remove);

document.addEventListener('keydown', (e) => {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    e.preventDefault();
    save();
  }
});

(async function init() {
  try {
    await refresh();
  } catch (err) {
    setStatus('error: ' + err.message);
  }
})();
