import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ApiDiscoveryRoutingModule } from './api-discovery-routing.module';
import { ApiDiscoveryOverviewComponent } from './components/api-discovery-overview/api-discovery-overview.component';


@NgModule({
  declarations: [
    ApiDiscoveryOverviewComponent
  ],
  imports: [
    CommonModule,
    ApiDiscoveryRoutingModule
  ]
})
export class ApiDiscoveryModule { }
