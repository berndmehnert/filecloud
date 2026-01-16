import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class SharedInputService {
    // Signal to hold the input value
  readonly inputValue = signal('');
  readonly isUploadingFile = signal(false);

  updateUploadingStatus(isUploading: boolean): void {
    this.isUploadingFile.set(isUploading);
  }

  getUploadingStatus(): boolean {
    return this.isUploadingFile();      
  }

  updateInput(value: string): void {
    this.inputValue.set(value);
  }

  getInput(): string {
    return this.inputValue();
  }
}
