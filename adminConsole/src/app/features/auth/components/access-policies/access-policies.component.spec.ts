import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AccessPoliciesComponent } from './access-policies.component';

describe('AccessPoliciesComponent', () => {
  let component: AccessPoliciesComponent;
  let fixture: ComponentFixture<AccessPoliciesComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [AccessPoliciesComponent]
    });
    fixture = TestBed.createComponent(AccessPoliciesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
