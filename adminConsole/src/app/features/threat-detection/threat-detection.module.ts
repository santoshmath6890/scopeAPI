import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ThreatDetectionRoutingModule } from './threat-detection-routing.module';
import { ThreatDetectionOverviewComponent } from './components/threat-detection-overview/threat-detection-overview.component';


@NgModule({
  declarations: [
    ThreatDetectionOverviewComponent
  ],
  imports: [
    CommonModule,
    ThreatDetectionRoutingModule
  ]
})
export class ThreatDetectionModule { }
