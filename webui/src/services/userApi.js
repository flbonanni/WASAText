// webui/src/services/userApi.js
import axios from 'axios'

// Axios è già configurato in main.js con l'interceptor che aggiunge
// `Authorization: Bearer <userId>` a tutte le richieste.

export function getCurrentUser() {
  // GET /me legge l'header e restituisce { id, name, avatarUrl, email, ... }
  return axios.get('http://localhost:8080/me')
}