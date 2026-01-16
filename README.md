A little Fileserver/Filecloud project in go, with a separate frontend in Angular/Node.js.

![Status](https://img.shields.io/badge/Status-Work_in_Progress-orange)

### 📋 Fileserver API Reference 

**Base URL:** `/v1`

| Method | Endpoint | Description | Request Body / Params |
| :--- | :--- | :--- | :--- |
| **GET** | `/files` | List all files (metadata only) | `?q='some text'`, `?limit=20`, TODO cursor |
| **POST** | `/files` | Upload a new file | `multipart/form-data` |
| **DELETE** | `/files/{id}` | Permanently delete a file | Not implemented yet! |
| **GET** | `/files/{id}/content` | Download the binary file content | - |
| **GET** | `/files/{id}/thumbnail` | Get the thumbnail image | - |

<br>

### 🛠 Usage Examples

**Upload a file:**
```bash
curl -X POST http://localhost:8080/v1/files  -F "file=@/path/to/image.jpg"
```
**Download a file:**
```bash
curl -OJ http://localhost:8080/v1/files/123/content
```

### 🚧 Roadmap

**Core Features**
- [x] **File Upload:** Implementation of file upload (frontend and backend)
- [x] **File Download:** Content delivery with proper headers plus usage in UI
- [ ] **File Deletion:** Implementation of handler for `DELETE /v1/files/{id}` plus application in UI
- [ ] **Pagination:** Cursor-based handling for the file list endpoint

**Media Handling**
- [x] **Image Thumbnails:** Auto-generation and display in UI
- [ ] **File Icons:** Fallback icons for non-image file types (PDF, Doc, etc.)

**Security & Quality**
- [ ] **Auth:** User Authentication and Authorization
- [ ] **Testing:** Unit and Integration tests
- [ ] **Documentation:** Setup instructions for running the full stack