import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { DashboardRoutingModule } from './dashboard-routing.module';
import { DashboardOverviewComponent } from './components/dashboard-overview/dashboard-overview.component';


@NgModule({
  declarations: [
    DashboardOverviewComponent
  ],
  imports: [
    CommonModule,
    DashboardRoutingModule
  ]
})
export class DashboardModule { }
