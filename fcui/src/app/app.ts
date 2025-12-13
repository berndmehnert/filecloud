import { Component, inject, signal } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { SharedInputService } from './shared-input-service';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, RouterLink, ReactiveFormsModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('fileserver UI');
  sharedInputService = inject(SharedInputService);
  filterInput = new FormGroup({
    name: new FormControl(''),
  });

  handleSubmit() {
    if (this.filterInput.value.name && this.filterInput.value.name.length > 0) {
      this.sharedInputService.inputValue.set(this.filterInput.value.name);
    }
    this.filterInput.reset();
  }
}
