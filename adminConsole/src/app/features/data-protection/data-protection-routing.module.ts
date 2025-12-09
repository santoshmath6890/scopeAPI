import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DataProtectionOverviewComponent } from './components/data-protection-overview/data-protection-overview.component';
import { authGuard } from '../../core/guards/auth.guard';

const routes: Routes = [
  { path: '', component: DataProtectionOverviewComponent, canActivate: [authGuard] }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DataProtectionRoutingModule { }
