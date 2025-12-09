import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DataProtectionOverviewComponent } from './data-protection-overview.component';

describe('DataProtectionOverviewComponent', () => {
  let component: DataProtectionOverviewComponent;
  let fixture: ComponentFixture<DataProtectionOverviewComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [DataProtectionOverviewComponent]
    });
    fixture = TestBed.createComponent(DataProtectionOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
