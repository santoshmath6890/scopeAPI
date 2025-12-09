import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AttackProtectionOverviewComponent } from './attack-protection-overview.component';

describe('AttackProtectionOverviewComponent', () => {
  let component: AttackProtectionOverviewComponent;
  let fixture: ComponentFixture<AttackProtectionOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AttackProtectionOverviewComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AttackProtectionOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
