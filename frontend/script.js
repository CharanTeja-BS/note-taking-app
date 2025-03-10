const API_URL = "https://note-taking-app-8n69.onrender.com/";

// DOM Elements
const noteForm = document.getElementById("note-form");
const noteTitle = document.getElementById("note-title");
const noteContent = document.getElementById("note-content");
const searchInput = document.getElementById("search");
const searchButton = document.getElementById("search-button");
const notesList = document.getElementById("notes-list");
const darkModeToggle = document.getElementById("dark-mode-toggle");

let editingNoteId = null; // Track the note being edited

// Fetch all notes
async function fetchNotes() {
    const response = await fetch(`${API_URL}/notes`);
    const notes = await response.json();
    displayNotes(notes);
}

// Display notes
function displayNotes(notes) {
    notesList.innerHTML = "";
    notes.forEach(note => {
        if (!note.ID) {
            console.error("Note ID is undefined:", note);
            return;
        }
        const noteElement = document.createElement("div");
        noteElement.className = "note";
        noteElement.innerHTML = `
            <h3>${note.title}</h3>
            <div class="content">${marked.parse(note.content)}</div>
            <div class="actions">
                <button class="edit" onclick="editNote(${note.ID}, '${note.title}', '${note.content}')">Edit</button>
                <button class="delete" onclick="deleteNote(${note.ID})">Delete</button>
            </div>
        `;
        notesList.appendChild(noteElement);
    });
}

// Add or update a note
noteForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const newNote = {
        title: noteTitle.value,
        content: noteContent.value,
    };

    if (editingNoteId) {
        // Update existing note
        await fetch(`${API_URL}/notes/${editingNoteId}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(newNote),
        });
        editingNoteId = null; // Reset editing state
        document.querySelector("#note-form button[type='submit']").textContent = "Save Note";
    } else {
        // Create new note
        await fetch(`${API_URL}/notes`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(newNote),
        });
    }

    noteForm.reset();
    fetchNotes();
});

// Edit a note
function editNote(id, title, content) {
    editingNoteId = id;
    noteTitle.value = title;
    noteContent.value = content;
    document.querySelector("#note-form button[type='submit']").textContent = "Update Note";
}

// Delete a note
async function deleteNote(id) {
    console.log("Deleting note with ID:", id); // Debugging
    if (!id) {
        console.error("Note ID is undefined");
        return;
    }
    await fetch(`${API_URL}/notes/${id}`, { method: "DELETE" });
    fetchNotes();
}

// Search notes
searchButton.addEventListener("click", async () => {
    const query = searchInput.value;
    const response = await fetch(`${API_URL}/notes?q=${query}`);
    const notes = await response.json();
    displayNotes(notes);
});

// Dark mode toggle
darkModeToggle.addEventListener("click", () => {
    document.body.classList.toggle("dark-mode");
    document.body.setAttribute(
        "data-theme",
        document.body.classList.contains("dark-mode") ? "dark" : "light"
    );
    darkModeToggle.textContent = document.body.classList.contains("dark-mode") ? "☀️" : "🌙";
});

// Initial fetch
fetchNotes();
