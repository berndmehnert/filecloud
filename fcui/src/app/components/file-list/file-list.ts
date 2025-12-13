import { Component, inject } from "@angular/core";
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
  filteredFiles: FileMeta[] = [];
  fs = inject(FileService);
  sharedInputService = inject(SharedInputService);
  
  constructor() { 
    this.fs.list('ern').subscribe(page => {
      this.filteredFiles = page.items;
    });
  }

  filterResults(query: string) {
  }
}
