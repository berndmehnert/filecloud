import { Component, inject, input } from '@angular/core';
import { FileMeta } from '../../models/file-meta.model';
import { DatePipe } from '@angular/common';
import { FileSizePipe } from '../../file-size-pipe';
import { FileService } from '../../file-service';

@Component({
  selector: 'app-display-file',
  imports: [DatePipe, FileSizePipe],
  templateUrl: './display-file.html',
  styleUrl: './display-file.css',
  standalone: true,
})
export class DisplayFile {
  fileMeta = input.required<FileMeta>();
  fileService = inject(FileService);

  handleClick() {
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
