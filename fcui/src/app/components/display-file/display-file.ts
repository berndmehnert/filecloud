import { Component, computed, effect, inject, input, signal } from '@angular/core';
import { FileMeta } from '../../models/file-meta.model';
import { DatePipe } from '@angular/common';
import { FileSizePipe } from '../../file-size-pipe';
import { FileService } from '../../file-service';
import {
  NgbModal, NgbDropdownToggle,
  NgbDropdownMenu,
  NgbDropdownItem,
  NgbDropdownButtonItem,
  NgbDropdown,
} from '@ng-bootstrap/ng-bootstrap';
import { ConfirmDialog } from '../confirm-dialog/confirm-dialog';
import { ConfirmDialogResultModel } from '../../models/confirm-dialog-result.model';
import { HttpClient, HttpStatusCode } from '@angular/common/http';

export interface ThumbnailState {
  status: 'loading' | 'ready' | 'pending' | 'processing' | 'failed' | 'not-found';
  url?: string;
}

@Component({
  selector: 'app-display-file',
  imports: [DatePipe, FileSizePipe, NgbDropdown, NgbDropdownToggle, NgbDropdownMenu, NgbDropdownItem, NgbDropdownButtonItem],
  templateUrl: './display-file.html',
  styleUrl: './display-file.css',
  standalone: true,
})
export class DisplayFile {
  fileMeta = input.required<FileMeta>();
  fileService = inject(FileService);
  private modalService = inject(NgbModal);
  private http = inject(HttpClient);

  state = signal<ThumbnailState>({ status: 'loading' });

  constructor() {
    effect(() => {
      if (this.fileMeta().mime.startsWith('image/')) {
        this.loadThumbnail(this.fileMeta().id);
      }
    });
  }

  async handleClick() {
    const modalRef = this.modalService.open(ConfirmDialog, {
      centered: true,
      backdrop: 'static'
    });

    modalRef.componentInstance.title.set('Download');
    modalRef.componentInstance.message.set('Do you want to download ' + this.fileMeta().filename + '?');
    modalRef.componentInstance.confirmText.set('Yes');
    modalRef.componentInstance.confirmBtnClass.set('btn-primary');
    modalRef.componentInstance.showWarning.set(false);

    try {
      const result: ConfirmDialogResultModel = await modalRef.result.then(res => res);
      if (result.confirmed) {
        this.fileService.download(this.fileMeta().id).subscribe(blob => {
          this.downloadBlob(blob, this.fileMeta().filename);
        });
      }// Perform delete action here
    } catch (e) {
      console.log('Cancelled');
    }
  }

  handleDownloadClick() {
    this.fileService.download(this.fileMeta().id).subscribe(blob => {
      this.downloadBlob(blob, this.fileMeta().filename);
    });
  }

  private loadThumbnail(id: number) {
    this.state.set({ status: 'loading' });

    this.http.get(this.fileService.getThumbnailUrl(this.fileMeta().id), {
      observe: 'response',
      responseType: 'blob'
    }).subscribe({
      next: (response) => {
        if (response.status === HttpStatusCode.Accepted) {
          // 202 - Check the JSON status
          this.parseStatusResponse(response.body!);
        } else {
          // 200 - It's an image
          const objectUrl = URL.createObjectURL(response.body!);
          this.state.set({ status: 'ready', url: objectUrl });
        }
      },
      error: (err) => {
        if (err.status === HttpStatusCode.NotFound) {
          this.state.set({ status: 'not-found' });
        } else if (err.status === HttpStatusCode.InternalServerError) {
          this.state.set({ status: 'failed' });
        } else {
          this.state.set({ status: 'failed' });
        }
      }
    });
  }

  private async parseStatusResponse(blob: Blob) {
    try {
      const text = await blob.text();
      const json = JSON.parse(text);

      if (json.status === 'pending') {
        this.state.set({ status: 'pending' });
        this.scheduleRetry();
      } else if (json.status === 'processing') {
        this.state.set({ status: 'processing' });
        this.scheduleRetry();
      }
    } catch {
      this.state.set({ status: 'failed' });
    }
  }

  private scheduleRetry() {
    setTimeout(() => {
      this.loadThumbnail(this.fileMeta().id);
    }, 3000);
  }


  private downloadBlob(data: Blob, filename = 'file.bin') {
    const url = URL.createObjectURL(data);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  }
}
