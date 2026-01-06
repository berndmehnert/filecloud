A little Fileserver/Filecloud project in go, with a separate frontend in Angular/nodeJS. 

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

## TODOS
- write tests
- complete file upload (done)
- complete cursor handling
- implement file deleting
- provide thumbnails for file display, for the moment only for images (done)
- display image thumbnails and for different file types appropriate icons  (partially done)
- add users + authentication + authorization 
- extend this document appropriately, in particular, include how everything is run