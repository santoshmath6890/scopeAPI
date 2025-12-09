import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ApiDiscoveryOverviewComponent } from './components/api-discovery-overview/api-discovery-overview.component';
import { authGuard } from '../../core/guards/auth.guard';

const routes: Routes = [
  { path: '', component: ApiDiscoveryOverviewComponent, canActivate: [authGuard] }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ApiDiscoveryRoutingModule { }
