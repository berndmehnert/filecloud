import { Component, computed, inject, input, signal } from '@angular/core';
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
  thumbnailUrl = computed(() => this.fileService.getThumbnailUrl(this.fileMeta().id));

  async handleClick() {
    const modalRef = this.modalService.open(ConfirmDialog, {
      centered: true,
      backdrop: 'static'
    });

    // ✅ Now you can use .set() because they're regular signals
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

  private downloadBlob(data: Blob, filename = 'file.bin') {
    const url = URL.createObjectURL(data);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);        // optional but safer for some browsers
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  }
}
