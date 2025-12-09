import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ThreatDetectionOverviewComponent } from './threat-detection-overview.component';

describe('ThreatDetectionOverviewComponent', () => {
  let component: ThreatDetectionOverviewComponent;
  let fixture: ComponentFixture<ThreatDetectionOverviewComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ThreatDetectionOverviewComponent]
    });
    fixture = TestBed.createComponent(ThreatDetectionOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
