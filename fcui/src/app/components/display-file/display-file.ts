import { Component, input } from '@angular/core';
import { FileMeta } from '../../models/file-meta.model';
import { DatePipe } from '@angular/common';
import { FileSizePipe } from '../../file-size-pipe';

@Component({
  selector: 'app-display-file',
  imports: [DatePipe, FileSizePipe],
  templateUrl: './display-file.html',
  styleUrl: './display-file.css',
  standalone: true,
})
export class DisplayFile {
    fileMeta = input.required<FileMeta>();

    handleClick() {
      console.log('File clicked:', this.fileMeta().filename);
    }
}
