import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DisplayFile } from './display-file';

describe('DisplayFile', () => {
  let component: DisplayFile;
  let fixture: ComponentFixture<DisplayFile>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DisplayFile]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DisplayFile);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
