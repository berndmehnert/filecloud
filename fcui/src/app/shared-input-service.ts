import { Injectable, signal } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class SharedInputService {
    // Signal to hold the input value
  readonly inputValue = signal('');

  updateInput(value: string): void {
    this.inputValue.set(value);
  }

  getInput(): string {
    return this.inputValue();
  }
}
