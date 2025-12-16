import { Component, inject, input, signal } from '@angular/core';
import { NgbActiveModal, NgbModal } from '@ng-bootstrap/ng-bootstrap';

@Component({
  standalone: true,
  selector: 'app-confirm-dialog',
  imports: [],
  templateUrl: './confirm-dialog.html',
  styleUrl: './confirm-dialog.css',
})
export class ConfirmDialog {
  activeModal = inject(NgbActiveModal);

  // ✅ Use regular signals (writable) instead of input()
  title = signal<string>('Confirm');
  message = signal<string>('Are you sure?');
  confirmText = signal<string>('Confirm');
  confirmBtnClass = signal<string>('btn-primary');
  showWarning = signal<boolean>(false);

  onConfirm() {
    this.activeModal.close({ confirmed: true, timestamp: Date.now() });
  }
}
