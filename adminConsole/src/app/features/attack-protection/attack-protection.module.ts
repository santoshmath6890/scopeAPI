import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AttackProtectionRoutingModule } from './attack-protection-routing.module';
import { AttackProtectionOverviewComponent } from './components/attack-protection-overview/attack-protection-overview.component';

@NgModule({
  declarations: [
    AttackProtectionOverviewComponent
  ],
  imports: [
    CommonModule,
    AttackProtectionRoutingModule
  ]
})
export class AttackProtectionModule { }
