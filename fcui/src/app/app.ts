import { Component, effect, inject, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SharedInputService } from './shared-input-service';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { FileService, UploadResponse } from './file-service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, ReactiveFormsModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('fileserver UI');
  sharedInputService = inject(SharedInputService);
  fileService = inject(FileService)
  uploading = signal(false);
  progress = signal(0);

  filterInput = new FormGroup({
    name: new FormControl(''),
  });

  constructor() {
    effect(() => {
      const isUploading = this.sharedInputService.getUploadingStatus();
      this.uploading.set(isUploading);
    });
  }

  handleSubmit() {
    this.sharedInputService.updateInput(this.filterInput.value.name || '');
    this.filterInput.reset();
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files?.length) return;

    const file = input.files[0];

    // Reset input so same file can be selected again
    input.value = '';

    // Upload immediately
    this.sharedInputService.updateUploadingStatus(true);

    this.fileService.upload(file).subscribe({
      next: (status) => {
        this.progress.set(status.percent);
        if (status.complete) {
          this.sharedInputService.updateUploadingStatus(false);
          console.log('Uploaded:', status.response);
        }
      },
      error: (err) => {
        console.error('Upload failed:', err);
        this.sharedInputService.updateUploadingStatus(false);
      }
    });
  }
}
