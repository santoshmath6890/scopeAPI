import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ThreatDetectionOverviewComponent } from './components/threat-detection-overview/threat-detection-overview.component';
import { authGuard } from '../../core/guards/auth.guard';

const routes: Routes = [
  { path: '', component: ThreatDetectionOverviewComponent, canActivate: [authGuard] }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ThreatDetectionRoutingModule { }
