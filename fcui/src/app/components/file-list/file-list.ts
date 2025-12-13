import { Component, effect, inject, signal } from "@angular/core";
import { DisplayFile } from "../display-file/display-file";
import { FileService } from "../../file-service";
import { FileMeta } from "../../models/file-meta.model";
import { SharedInputService } from "../../shared-input-service";


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

  constructor() {
    effect(() => {
      const searchTerm = this.sharedInputService.getInput();
      console.log('trigger', ' ', searchTerm);
      this.fs.list(searchTerm).subscribe(page => {
        this.filteredFiles.set(page.items);
      });
    });
  }
}
