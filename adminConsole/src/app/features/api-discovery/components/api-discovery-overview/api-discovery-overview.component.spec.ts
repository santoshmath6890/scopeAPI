import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ApiDiscoveryOverviewComponent } from './api-discovery-overview.component';

describe('ApiDiscoveryOverviewComponent', () => {
  let component: ApiDiscoveryOverviewComponent;
  let fixture: ComponentFixture<ApiDiscoveryOverviewComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ApiDiscoveryOverviewComponent]
    });
    fixture = TestBed.createComponent(ApiDiscoveryOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
