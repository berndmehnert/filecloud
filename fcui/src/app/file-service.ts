import { HttpClient, HttpEventType } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { filter, map, Observable } from 'rxjs';
import { FileMetaPage } from './models/file-meta.model';

export interface UploadResponse {
  id: number;
  filename: string;
  size: number;
  checksum: string;
}

export interface UploadProgress {
  percent: number;
  complete: boolean;
  response?: UploadResponse;
}

@Injectable({
  providedIn: 'root',
})
export class FileService {
  private http = inject(HttpClient);
  private base = 'http://localhost:8080';
  
  list(filter : string): Observable<FileMetaPage> {
    return this.http.get<FileMetaPage>(`${this.base}/api/files?q=${filter}`);
  }

  download(id: number): Observable<Blob> {
    return this.http.get<Blob>(`${this.base}/files/${id}`, { responseType: 'blob' as 'json' });
  }

  upload(file: File): Observable<UploadProgress> {
    const formData = new FormData();
    formData.append('file', file);

    return this.http.post<UploadResponse>(this.base + '/upload', formData, {
      reportProgress: true,
      observe: 'events'
    }).pipe(
      filter(event => 
        event.type === HttpEventType.UploadProgress ||
        event.type === HttpEventType.Response
      ),
      map(event => {
        if (event.type === HttpEventType.UploadProgress) {
          return {
            percent: event.total ? Math.round(100 * event.loaded / event.total) : 0,
            complete: false
          };
        } else {
          return {
            percent: 100,
            complete: true,
            response: event.body as UploadResponse
          };
        }
      })
    );
  }
  getThumbnailUrl(fileId: number): string {
    console.log(`${this.base}/files/${fileId}/thumbnail`);
    return `${this.base}/files/${fileId}/thumbnail`;
  }
}
