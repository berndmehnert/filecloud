import { Component, effect, inject, signal } from "@angular/core";
import { DisplayFile } from "../display-file/display-file";
import { FileService } from "../../file-service";
import { FileMeta } from "../../models/file-meta.model";
import { SharedInputService } from "../../shared-input-service";

type ThumbnailState = 
  | { status: 'loading' }
  | { status: 'ready'; url: SafeUrl }
  | { status: 'pending' }
  | { status: 'processing' }
  | { status: 'failed' }
  | { status: 'not-found' };

@Component({
  selector: 'app-file-list',
  imports: [DisplayFile],
  templateUrl: './file-list.html',
  styleUrl: './file-list.css',
  standalone: true,
})
export class FileList {
  filteredFiles = signal<FileMeta[]>([]);
  fs = inject(FileService);
  sharedInputService = inject(SharedInputService);
  currentSearchTerm = signal(''); 

  constructor() {
    effect(() => {
      this.currentSearchTerm.set(this.sharedInputService.getInput());
      this.updateFileList();  
    });
    effect(() => {
      const isuploaded = this.sharedInputService.getUploadingStatus();
      if (isuploaded === false) {
        this.updateFileList();
      }
    });
  }

  updateFileList() {
    this.fs.list(this.currentSearchTerm()).subscribe(page => {
      this.filteredFiles.set(page.items);
    });
  }
}
