import { useEffect, useState } from "react";
import { format } from "date-fns";
import {
  Box,
  Button,
  Card,
  CardContent,
  Checkbox,
  Container,
  IconButton,
  List,
  ListItem,
  ListItemText,
  TextField,
  Typography,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import SaveIcon from "@mui/icons-material/Save";

const API_BASE = "http://localhost:8080";

function App() {
  const [todos, setTodos] = useState([]);
  const [text, setText] = useState("");
  const [editingId, setEditingId] = useState(null);
  const [editText, setEditText] = useState("");

  const fetchTodos = async () => {
    const res = await fetch(`${API_BASE}/todos`);
    const data = await res.json();
    setTodos(data);
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  const addTodo = async (e) => {
    e.preventDefault();
    if (!text.trim()) return;

    await fetch(`${API_BASE}/todos`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ text }),
    });

    setText("");
    fetchTodos();
  };

  const toggleTodo = async (id) => {
    await fetch(`${API_BASE}/todos/${id}/toggle`, {
      method: "PUT",
    });
    fetchTodos();
  };

  const deleteTodo = async (id) => {
    await fetch(`${API_BASE}/todos/${id}`, {
      method: "DELETE",
    });
    fetchTodos();
  };

  const startEdit = (todo) => {
    setEditingId(todo.id);
    setEditText(todo.text);
  };

  const saveEdit = async (todo) => {
    await fetch(`${API_BASE}/todos/${todo.id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        text: editText,
        done: todo.done,
      }),
    });

    setEditingId(null);
    setEditText("");
    fetchTodos();
  };

  function formatDate(dateString) {
    return format(new Date(dateString), "MMM d, HH:mm");
  }

  return (
    <Box
      sx={{
        minHeight: "100vh",
        backgroundColor: "#f5f6f8",
        py: 6,
      }}
    >
      <Container
        maxWidth="md"
        sx={{
          maxWidth: "720px !important",
        }}
      >
        <Card sx={{ borderRadius: 3 }}>
          <CardContent sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom>
              ToDo List
            </Typography>

            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              React + Go + PostgreSQL
            </Typography>

            <Box
              component="form"
              onSubmit={addTodo}
              sx={{ display: "flex", gap: 1, mb: 3 }}
            >
              <TextField
                fullWidth
                label="Add a new task"
                value={text}
                onChange={(e) => setText(e.target.value)}
              />
              <Button type="submit" variant="contained" sx={{ px: 3 }}>
                Add
              </Button>
            </Box>

            <List sx={{ mt: 2 }}>
              {todos.map((todo) => (
                <ListItem
                  key={todo.id}
                  sx={{
                    mb: 2,
                    px: 2.5,
                    py: 1.8,
                    borderRadius: 2,
                    border: "1px solid #eee",
                    backgroundColor: "#fafafa",
                    boxShadow: "0 1px 2px rgba(0,0,0,0.06)",
                    display: "flex",
                    alignItems: "flex-start",
                    gap: 1.5,
                  }}
                >
                  <Checkbox
                    checked={todo.done}
                    onChange={() => toggleTodo(todo.id)}
                    color="success"
                    sx={{ mt: 0.5 }}
                  />

                  <Box sx={{ flex: 1, minWidth: 0 }}>
                    {editingId === todo.id ? (
                      <TextField
                        fullWidth
                        size="small"
                        value={editText}
                        onChange={(e) => setEditText(e.target.value)}
                      />
                    ) : (
                      <ListItemText
                        primary={
                          <Typography
                            sx={{
                              wordBreak: "break-word",
                              overflowWrap: "anywhere",
                              lineHeight: 1.45,
                            }}
                          >
                            {todo.text}
                          </Typography>
                        }
                        secondary={
                          <Typography
                            variant="body2"
                            color="text.secondary"
                            sx={{ mt: 0.5 }}
                          >
                            Created {formatDate(todo.created_at)} • Updated{" "}
                            {formatDate(todo.updated_at)}
                          </Typography>
                        }
                        sx={{ m: 0 }}
                      />
                    )}
                  </Box>

                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      gap: 0.5,
                      ml: 1,
                      flexShrink: 0,
                    }}
                  >
                    <IconButton onClick={() => deleteTodo(todo.id)} color="error">
                      <DeleteIcon />
                    </IconButton>

                    {editingId === todo.id ? (
                      <IconButton onClick={() => saveEdit(todo)} color="primary">
                        <SaveIcon />
                      </IconButton>
                    ) : (
                      <IconButton onClick={() => startEdit(todo)} color="primary">
                        <EditIcon />
                      </IconButton>
                    )}
                  </Box>
                </ListItem>
              ))}
            </List>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
}

export default App;