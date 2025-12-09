import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AttackProtectionOverviewComponent } from './components/attack-protection-overview/attack-protection-overview.component';

const routes: Routes = [
  { path: '', component: AttackProtectionOverviewComponent }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AttackProtectionRoutingModule { }
