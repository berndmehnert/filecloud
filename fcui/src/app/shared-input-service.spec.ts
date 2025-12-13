import { TestBed } from '@angular/core/testing';

import { SharedInputService } from './shared-input-service';

describe('SharedInputService', () => {
  let service: SharedInputService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SharedInputService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
