# Go Dropbox

A simple file upload/download service built with:

- **Backend**: Go
- **Frontend**: Vanilla JavaScript
- **Styling**: CSS

## Features

- File upload with progress bar
- File listing with type-specific icons
- Download/delete functionality
- Responsive UI

## How to Run

1. Install Go (https://golang.org/dl/)
2. Clone this repository:
   ```bash
   git clone https://github.com/anik4good/go_dropbox.git
   ```
3. Navigate to project directory:
   ```bash
   cd go_dropbox
   ```
4. Run the server:
   ```bash
   go run main.go
   ```
5. Open http://localhost:8080 in your browser

## Project Structure

- `main.go` - Go backend server
- `static/` - Frontend files
  - `script.js` - Client-side JavaScript
  - `style.css` - Styling
