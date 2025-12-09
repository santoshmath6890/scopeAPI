import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DashboardOverviewComponent } from './components/dashboard-overview/dashboard-overview.component';
import { authGuard } from '../../core/guards/auth.guard';

const routes: Routes = [
  { path: '', component: DashboardOverviewComponent, canActivate: [authGuard] }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DashboardRoutingModule { }
