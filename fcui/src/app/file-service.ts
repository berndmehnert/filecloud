import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { FileMeta, FileMetaPage } from './models/file-meta.model';

@Injectable({
  providedIn: 'root',
})
export class FileService {
  private http = inject(HttpClient);
  private base = 'http://localhost:8080'; // set if backend is on same origin or prefix like 'http://localhost:8080'
  
  list(): Observable<FileMetaPage> {
    return this.http.get<FileMetaPage>(`${this.base}/api/files`);
  }

  download(id: number): string {
    return `${this.base}/files/${id}`;
  }
}
