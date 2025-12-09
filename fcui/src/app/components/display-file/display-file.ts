import { Component, input } from '@angular/core';
import { FileMeta } from '../../models/file-meta.model';

@Component({
  selector: 'app-display-file',
  imports: [],
  templateUrl: './display-file.html',
  styleUrl: './display-file.css',
  standalone: true,
})
export class DisplayFile {
    fileMeta = input.required<FileMeta>();
}
