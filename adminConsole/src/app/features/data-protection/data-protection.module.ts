import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { DataProtectionRoutingModule } from './data-protection-routing.module';
import { DataProtectionOverviewComponent } from './components/data-protection-overview/data-protection-overview.component';


@NgModule({
  declarations: [
    DataProtectionOverviewComponent
  ],
  imports: [
    CommonModule,
    DataProtectionRoutingModule
  ]
})
export class DataProtectionModule { }
