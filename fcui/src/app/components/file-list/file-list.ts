import { Component, inject } from "@angular/core";
import { DisplayFile } from "../display-file/display-file";
import { FileService } from "../../file-service";
import { FileMeta } from "../../models/file-meta.model";


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
  
  constructor() { 
    this.fs.list().subscribe(page => {
      this.filteredFiles = page.items;
    });
  }

  filterResults(query: string) {
   // this.fs.search(query).subscribe(page => {
    //  this.filteredFiles = page.items;
    //});
  }
}
